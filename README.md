# Personal Blog Lab

一个用于学习 Vue、TypeScript、Go 和 PostgreSQL 的个人博客项目。当前仓库只有最小可运行骨架，文章、认证和评论功能刻意留给你实现。

## 技术栈

- Web: Vue 3 + TypeScript + Vite
- API: Go + Gin
- Database: PostgreSQL
- Local infrastructure: Docker Compose

## 目录

```text
.
├── web/                 # 游客端；以后也可以在这里加入管理端路由
├── server/              # Go API
│   ├── cmd/api/         # 程序入口
│   └── internal/        # 未来的业务模块
├── docs/                # 需求和学习记录
├── db/init/             # PostgreSQL 首次初始化 SQL
├── deploy/              # Caddy 生产配置
├── compose.yaml         # 本地 PostgreSQL
└── compose.prod.yaml    # 生产环境服务编排
```

## 准备环境

需要安装 Node.js 22+、Go 1.24+ 和 Docker Desktop。安装后检查：

```bash
node --version
npm --version
go version
docker compose version
```

## 启动

```bash
cp .env.example .env
docker compose up -d db
```

另开两个终端：

```bash
cd server
set -a
source ../.env
set +a
go run ./cmd/api
```

```bash
cd web
npm install
npm run dev
```

浏览器打开 `http://localhost:5173`，API 健康检查地址为 `http://localhost:8080/api/health`。

## 你的第一个里程碑

只实现“已发布文章列表”这一条纵向功能：

1. 自己设计 `articles` 表和数据库迁移。
2. 在 Go 中实现 `GET /api/articles`。
3. 在 Vue 中请求接口并展示标题与摘要。
4. 为 API 添加至少一个测试。

先不要实现登录、评论、Redis 或搜索。完成这个闭环后，再决定下一步。

更具体的边界见 [docs/ROADMAP.md](docs/ROADMAP.md)。

## 生产部署

生产环境使用 Caddy 提供 HTTPS，并将 `/api` 转发给 Go API。参见
`.env.production.example`、`compose.prod.yaml` 和 `deploy/Caddyfile`。

不要将 `.env.production` 提交到 Git。生产数据库首次创建时会自动执行
`db/init/001_create_articles.sql`。
