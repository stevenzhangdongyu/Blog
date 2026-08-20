package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"golang.org/x/crypto/bcrypt"
	"personal-blog/server/internal/article"
)

type sessionStore struct {
	mu     sync.RWMutex
	tokens map[string]string
}

type commentRateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
}

func newCommentRateLimiter() *commentRateLimiter {
	return &commentRateLimiter{attempts: make(map[string][]time.Time)}
}

func (l *commentRateLimiter) allow(clientIP string) bool {
	const limit = 5
	windowStart := time.Now().Add(-10 * time.Minute)
	l.mu.Lock()
	defer l.mu.Unlock()

	recent := l.attempts[clientIP][:0]
	for _, attempt := range l.attempts[clientIP] {
		if attempt.After(windowStart) {
			recent = append(recent, attempt)
		}
	}
	if len(recent) >= limit {
		l.attempts[clientIP] = recent
		return false
	}
	l.attempts[clientIP] = append(recent, time.Now())
	return true
}

func newSessionStore() *sessionStore {
	return &sessionStore{tokens: make(map[string]string)}
}

func (s *sessionStore) create(username string) (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	token := hex.EncodeToString(buffer)
	s.mu.Lock()
	s.tokens[token] = username
	s.mu.Unlock()
	return token, nil
}

func (s *sessionStore) username(token string) (string, bool) {
	s.mu.RLock()
	username, ok := s.tokens[token]
	s.mu.RUnlock()
	return username, ok
}

func (s *sessionStore) delete(token string) {
	s.mu.Lock()
	delete(s.tokens, token)
	s.mu.Unlock()
}

func ensureAdmin(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS admins (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}
	username := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		return nil
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM admins)").Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO admins (username, password_hash) VALUES ($1, $2)", username, string(hash))
	return err
}

func ensureComments(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS comments (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
		author_name VARCHAR(50) NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT comments_author_name_not_blank CHECK (length(trim(author_name)) > 0),
		CONSTRAINT comments_content_not_blank CHECK (length(trim(content)) > 0)
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS comments_article_id_created_at_idx
		ON comments (article_id, created_at ASC)`)
	return err
}

func requireAdmin(store *sessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("admin_session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录已失效"})
			return
		}
		if username, ok := store.username(cookie); ok {
			c.Set("admin_username", username)
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录已失效"})
	}
}

type articleInput struct {
	Title   string `json:"title" binding:"required"`
	Slug    string `json:"slug"`
	Summary string `json:"summary"`
	Content string `json:"content" binding:"required"`
	Status  string `json:"status"`
}

type commentInput struct {
	AuthorName string `json:"authorName" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Website    string `json:"website"`
}

type comment struct {
	ID         int       `json:"id"`
	AuthorName string    `json:"authorName"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
}

var slugInvalidCharacters = regexp.MustCompile(`[^a-z0-9]+`)

func slugFromTitle(title string) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	parts := make([]string, 0)
	var plain strings.Builder
	flushPlain := func() {
		if plain.Len() > 0 {
			parts = append(parts, plain.String())
			plain.Reset()
		}
	}
	for _, char := range title {
		if !unicode.Is(unicode.Han, char) {
			plain.WriteRune(char)
			continue
		}
		flushPlain()
		values := pinyin.Pinyin(string(char), args)
		if len(values) > 0 && len(values[0]) > 0 {
			parts = append(parts, values[0][0])
		}
	}
	flushPlain()
	slug := slugInvalidCharacters.ReplaceAllString(strings.ToLower(strings.Join(parts, "-")), "-")
	return strings.Trim(slug, "-")
}

func availableSlug(db *sql.DB, base, articleID string) (string, error) {
	if base == "" {
		base = "article"
	}
	for suffix := 1; ; suffix++ {
		candidate := base
		if suffix > 1 {
			candidate += "-" + strconv.Itoa(suffix)
		}
		query := "SELECT EXISTS (SELECT 1 FROM articles WHERE slug = $1)"
		args := []any{candidate}
		if articleID != "" {
			query = "SELECT EXISTS (SELECT 1 FROM articles WHERE slug = $1 AND id <> $2)"
			args = append(args, articleID)
		}
		var exists bool
		if err := db.QueryRow(query, args...).Scan(&exists); err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
}

func main() {
	db, err := article.GetDBConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := ensureAdmin(db); err != nil {
		panic(err)
	}
	if err := ensureComments(db); err != nil {
		panic(err)
	}

	store := newSessionStore()
	commentLimiter := newCommentRateLimiter()
	secureCookie := os.Getenv("COOKIE_SECURE") != "false"
	router := gin.Default()

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/api/articles", func(c *gin.Context) {
		articles, err := article.GetArticles(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"articles": articles})
	})
	router.GET("/api/articles/:slug", func(c *gin.Context) {
		found, err := article.GetArticleBySlug(db, c.Param("slug"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败"})
			return
		}
		if found == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		c.JSON(http.StatusOK, found)
	})
	router.GET("/api/articles/:slug/comments", getComments(db))
	router.POST("/api/articles/:slug/comments", createComment(db, commentLimiter))

	admin := router.Group("/api/admin")
	admin.POST("/login", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if c.ShouldBindJSON(&input) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
			return
		}
		var hash string
		err := db.QueryRow("SELECT password_hash FROM admins WHERE username = $1", input.Username).Scan(&hash)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		token, err := store.create(input.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建会话失败"})
			return
		}
		c.SetCookie("admin_session", token, 86400, "/", "", secureCookie, true)
		c.JSON(http.StatusOK, gin.H{"username": input.Username})
	})
	admin.POST("/logout", func(c *gin.Context) {
		if token, err := c.Cookie("admin_session"); err == nil {
			store.delete(token)
		}
		c.SetCookie("admin_session", "", -1, "/", "", secureCookie, true)
		c.Status(http.StatusNoContent)
	})
	admin.GET("/session", func(c *gin.Context) {
		cookie, err := c.Cookie("admin_session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}
		username, ok := store.username(cookie)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"username": username})
	})

	protected := admin.Group("/articles", requireAdmin(store))
	protected.GET("", func(c *gin.Context) {
		rows, err := db.Query(`SELECT id, title, slug, summary, content, status FROM articles ORDER BY id DESC`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败"})
			return
		}
		defer rows.Close()
		items := []article.Article{}
		for rows.Next() {
			var item article.Article
			if err := rows.Scan(&item.ID, &item.Title, &item.Slug, &item.Summary, &item.Content, &item.Status); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文章失败"})
				return
			}
			items = append(items, item)
		}
		c.JSON(http.StatusOK, gin.H{"articles": items})
	})
	protected.POST("", saveArticle(db, false))
	protected.PUT("/:id", saveArticle(db, true))
	protected.DELETE("/:id", deleteArticle(db))

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

func saveArticle(db *sql.DB, update bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input articleInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式必须是 JSON"})
			return
		}
		input.Title = strings.TrimSpace(input.Title)
		input.Slug = strings.TrimSpace(input.Slug)
		if input.Title == "" || strings.TrimSpace(input.Content) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "标题和正文不能为空"})
			return
		}
		if input.Slug == "" {
			generated, err := availableSlug(db, slugFromTitle(input.Title), c.Param("id"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 slug 失败"})
				return
			}
			input.Slug = generated
		}
		if input.Status != "published" {
			input.Status = "draft"
		}
		var err error
		if update {
			query := `UPDATE articles SET title=$1, slug=$2, summary=$3, content=$4, status=$5, updated_at=CURRENT_TIMESTAMP, published_at=CASE WHEN $5::VARCHAR='published' THEN COALESCE(published_at, CURRENT_TIMESTAMP) ELSE NULL END WHERE id=$6`
			log.Printf("executing SQL (update article): %s", query)
			_, err = db.Exec(query, input.Title, input.Slug, input.Summary, input.Content, input.Status, c.Param("id"))
		} else {
			query := `INSERT INTO articles (title, slug, summary, content, status, published_at) VALUES ($1,$2,$3,$4,$5,CASE WHEN $5::VARCHAR='published' THEN CURRENT_TIMESTAMP END)`
			log.Printf("executing SQL (create article): %s", query)
			_, err = db.Exec(query, input.Title, input.Slug, input.Summary, input.Content, input.Status)
		}
		if err != nil {
			log.Printf("save article failed: %v", err)
			if strings.Contains(err.Error(), "duplicate key") {
				c.JSON(http.StatusConflict, gin.H{"error": "slug 已存在，请换一个 slug"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文章失败"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func deleteArticle(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := db.Exec("DELETE FROM articles WHERE id = $1", c.Param("id")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func getComments(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT c.id, c.author_name, c.content, c.created_at
			FROM comments c
			JOIN articles a ON a.id = c.article_id
			WHERE a.slug = $1 AND a.status = 'published'
			ORDER BY c.created_at ASC`, c.Param("slug"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询评论失败"})
			return
		}
		defer rows.Close()

		comments := []comment{}
		for rows.Next() {
			var item comment
			if err := rows.Scan(&item.ID, &item.AuthorName, &item.Content, &item.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "读取评论失败"})
				return
			}
			comments = append(comments, item)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "读取评论失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"comments": comments})
	}
}

func createComment(db *sql.DB, limiter *commentRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "评论过于频繁，请十分钟后再试"})
			return
		}
		var input commentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写昵称和评论内容"})
			return
		}
		if strings.TrimSpace(input.Website) != "" {
			c.Status(http.StatusNoContent)
			return
		}
		input.AuthorName = strings.TrimSpace(input.AuthorName)
		input.Content = strings.TrimSpace(input.Content)
		if input.AuthorName == "" || input.Content == "" || len([]rune(input.AuthorName)) > 50 || len([]rune(input.Content)) > 1000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "昵称限 50 字，评论限 1000 字，且不能为空"})
			return
		}

		var created comment
		err := db.QueryRow(`
			INSERT INTO comments (article_id, author_name, content)
			SELECT id, $2, $3 FROM articles WHERE slug = $1 AND status = 'published'
			RETURNING id, author_name, content, created_at`, c.Param("slug"), input.AuthorName, input.Content).
			Scan(&created.ID, &created.AuthorName, &created.Content, &created.CreatedAt)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在或尚未发布"})
			return
		}
		if err != nil {
			log.Printf("create comment failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "发表评论失败"})
			return
		}
		c.JSON(http.StatusCreated, created)
	}
}
