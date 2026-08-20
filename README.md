# Personal Blog Lab

个人博客项目，用于学习 Vue、TypeScript、Go、PostgreSQL、Docker 和全栈部署。

线上地址：[zhangdongyu.vip](https://zhangdongyu.vip)

## 现有功能

- 游客查看已发布文章列表和文章详情
- Markdown 渲染、代码高亮、数学公式
- 游客发表评论，包含长度校验和简单频率限制
- 选中文章正文引用原文发表评论，点击引用返回原文位置
- 管理员登录、新建、编辑、删除文章
- 草稿/发布状态，Slug 自动生成和重复检查
- Markdown 编辑器、工具栏和实时预览
- 雪之下背景图每 15 秒轮播，并进行压缩、预加载、淡入淡出
- 全局点击碎纸效果和评论成功庆祝特效
- 移动端适配

## 技术栈

- 前端：Vue 3、TypeScript、Vite、Vue Router
- Markdown：markdown-it、highlight.js、KaTeX、DOMPurify
- 后端：Go、Gin、pgx
- 数据库：PostgreSQL 17
- 反向代理：Caddy 2
- 容器：Docker Compose

## 本地开发

准备 Node.js 22+、Go 1.24+ 和 Docker Desktop：

```bash
node --version
npm --version
go version
docker compose version
```

```bash
cp .env.example .env
docker compose up -d db
```

启动 API：

```bash
cd server
set -a
source ../.env
set +a
go run ./cmd/api
```

启动前端：

```bash
cd web
npm install
npm run dev
```

本地前端为 `http://localhost:5173`，API 健康检查为 `http://localhost:8080/api/health`，管理端为 `http://localhost:5173/admin/login`。

验证构建和测试：

```bash
npm --prefix web run build
cd server && GOCACHE=/tmp/blog-go-cache go test ./...
```

## 生产部署

生产地址：[zhangdongyu.vip](https://zhangdongyu.vip) 和 [www.zhangdongyu.vip](https://www.zhangdongyu.vip)。

### 1. 准备服务器

Ubuntu 云服务器安全组放行 `22`、`80`、`443`，不要开放 PostgreSQL 的 `5432`：

```bash
apt update
apt install -y ca-certificates curl git
curl -fsSL https://get.docker.com | sh
systemctl enable --now docker
mkdir -p /opt
cd /opt
git clone https://github.com/stevenzhangdongyu/Blog.git
cd Blog
```

DNS 添加 `www` 和 `@` 指向服务器公网 IP 的 A 记录。Caddy 会自动申请 HTTPS，并将根域名重定向到 `www`。

### 2. 配置环境变量

```bash
cp .env.production.example .env.production
vim .env.production
chmod 600 .env.production
```

```env
POSTGRES_DB=blog
POSTGRES_USER=blog
POSTGRES_PASSWORD=新的长随机密码
DATABASE_URL=postgres://blog:URL编码后的密码@db:5432/blog?sslmode=disable
ADMIN_USERNAME=admin
ADMIN_PASSWORD=管理端密码
COOKIE_SECURE=true
```

数据库密码在两个变量中必须相同。若密码含 `@`、`#`、`:`、`/` 等字符，要在 `DATABASE_URL` 中 URL 编码，例如 `#` 编码为 `%23`。不要提交 `.env.production`。

### 3. 启动和更新

```bash
docker compose --env-file .env.production -f compose.prod.yaml up -d --build
docker compose --env-file .env.production -f compose.prod.yaml ps
curl https://www.zhangdongyu.vip/api/health
```

后续发布新版本：

```bash
cd /opt/Blog
git pull origin main
docker compose --env-file .env.production -f compose.prod.yaml up -d --build
```

管理端地址为 `https://www.zhangdongyu.vip/admin/login`。不要随意执行 `docker compose down -v`，否则会删除数据库数据卷。

## 数据迁移和备份

Git 只同步代码和静态图片，文章、评论和管理员账号需要使用 PostgreSQL 导出导入。

本地项目目录执行：

```bash
docker compose exec -T db pg_dump -U blog -d blog --clean --if-exists > blog-production.sql
scp blog-production.sql root@服务器IP:/opt/Blog/
```

服务器导入前先备份生产数据库：

```bash
cd /opt/Blog
docker compose --env-file .env.production -f compose.prod.yaml exec -T db pg_dump -U blog -d blog > backup-before-import.sql
docker compose --env-file .env.production -f compose.prod.yaml exec -T db psql -U blog -d blog < blog-production.sql
```

数据库初始化脚本只在 PostgreSQL 数据卷第一次创建时执行；已有数据卷不会重复执行。新增字段时，应用启动会执行兼容迁移，但重要操作仍应先备份。

## 安全注意事项

- 不要提交 `.env`、`.env.production` 和数据库备份文件
- 使用强数据库密码和管理员密码
- 不要将 PostgreSQL 端口暴露到公网
- 定期备份生产数据库
- 生产服务器上不要随意使用 `git reset --hard`
