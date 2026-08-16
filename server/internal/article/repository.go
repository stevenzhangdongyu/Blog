package article

import (
	"database/sql"
	"fmt"
	"os"
	// 隐式导入 pgx 的 stdlib 驱动以注册 PostgreSQL 数据库驱动
	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetDBConnection() (*sql.DB, error) {
	// 1. 数据库连接字符串 (DSN)
	// 格式：postgres://用户名:密码@主机地址:端口/数据库名?sslmode=disable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// 2. 打开数据库连接（此时并未真正与数据库建立 TCP 连接）
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %w", err)
	}

	// 3. 验证数据库连接状态（Ping 会真正发起网络握手）
	if err := db.Ping(); err != nil {
		db.Close() // 确保在出错时关闭数据库连接
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	fmt.Println("成功连接到 PostgreSQL 数据库！")
	return db, nil
}
func GetArticles(db *sql.DB) ([]Article, error) {
	// 1. 执行 SQL 查询所有记录（按 ID 降序排序）
	query := `SELECT id, title, slug, summary, status FROM articles where status = 'published' ORDER BY id DESC;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询文章列表失败: %w", err)
	}
	// 必须确保关闭 rows，释放连接池资源
	defer rows.Close()

	// 2. 初始化一个切片用来存放多条记录
	// 建议初始化为空切片 []Article{}，这样查无数据时序列化为 JSON 会是 [] 而不是 null
	articles := []Article{}

	// 3. 循环遍历结果集中的每一行
	for rows.Next() {
		var a Article
		// 依次将当前行的各列映射到结构体字段中
		if err := rows.Scan(&a.ID, &a.Title, &a.Slug, &a.Summary, &a.Status); err != nil {
			return nil, fmt.Errorf("读取数据行失败: %w", err)
		}
		articles = append(articles, a)
	}

	// 4. 检查遍历过程中是否有错
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历数据失败: %w", err)
	}

	return articles, nil
}

func GetArticleBySlug(db *sql.DB, slug string) (*Article, error) {
	query := `
		SELECT id, title, slug, summary, content, status
		FROM articles
		WHERE slug = $1 AND status = 'published';
	`
	row := db.QueryRow(query, slug)

	var a Article
	if err := row.Scan(&a.ID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没有找到对应的文章
		}
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}

	return &a, nil
}
