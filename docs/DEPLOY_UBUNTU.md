# Ubuntu Production Deployment

## Architecture

```text
Internet -> Caddy (HTTPS) -> Vue Web / Go API -> PostgreSQL
```

Only Caddy exposes ports `80` and `443`. PostgreSQL and the Go API stay inside the Docker network.

## 1. DNS, Security Group, and ICP

Create DNS A records:

```text
@       -> SERVER_PUBLIC_IP
www     -> SERVER_PUBLIC_IP
```

In the Alibaba Cloud security group, allow inbound ports `22`, `80`, and `443`. Do not expose `5432` or `8080`.

For mainland China servers, public domain access normally requires ICP filing. A WeChat `non-compliance icp filtering` page indicates the domain is not compliant for access through a mainland server.

## 2. Push Production Files

On the development machine:

```bash
git add .
git commit -m "Add production Docker deployment"
git push
```

Never commit `.env.production`; it contains production secrets. If generated files were tracked by mistake:

```bash
git rm --cached server/api web/tsconfig.app.tsbuildinfo web/tsconfig.node.tsbuildinfo
git commit -m "Remove generated build artifacts"
git push
```

## 3. Install Docker

```bash
ssh root@SERVER_PUBLIC_IP
apt update
apt install -y docker.io docker-compose-v2 git curl
systemctl enable --now docker
docker --version
docker compose version
```

## 4. Configure Docker Image Acceleration

Get the dedicated accelerator URL from Alibaba Cloud ACR:

```text
Container Registry ACR -> Image Tools -> Image Accelerator
```

Create `/etc/docker/daemon.json`:

```json
{
  "registry-mirrors": ["https://YOUR_ACCELERATOR.mirror.aliyuncs.com"],
  "dns": ["223.5.5.5", "223.6.6.6"]
}
```

Restart Docker:

```bash
systemctl daemon-reload
systemctl restart docker
docker info
```

If an accelerator times out or lacks an official tag, explicitly pull a verified proxy image and tag it locally:

```bash
docker pull docker.m.daocloud.io/library/golang:1.25-alpine
docker tag docker.m.daocloud.io/library/golang:1.25-alpine golang:1.25-alpine
```

Use the same method for `node`, `postgres`, `caddy`, `nginx`, and `alpine` when needed. Third-party image mirrors are supply-chain dependencies; assess them before production use.

## 5. Clone the Repository

```bash
mkdir -p /opt
cd /opt
git clone --depth 1 --single-branch --branch main \
  https://github.com/stevenzhangdongyu/Blog.git
cd Blog
```

If GitHub is slow:

```bash
git config --global http.version HTTP/1.1
```

Then retry the shallow clone.

## 6. Create Production Secrets

```bash
cd /opt/Blog
cp .env.production.example .env.production
openssl rand -hex 32
vim .env.production
```

Use the same generated password in both values:

```env
POSTGRES_DB=blog
POSTGRES_USER=blog
POSTGRES_PASSWORD=replace_with_random_hex_password
DATABASE_URL=postgres://blog:replace_with_random_hex_password@db:5432/blog?sslmode=disable
```

Protect the file:

```bash
chmod 600 .env.production
```

## 7. Speed Up Docker Builds

In `server/Dockerfile`:

```dockerfile
RUN GOPROXY=https://goproxy.cn,direct go mod download
```

In `web/Dockerfile`:

```dockerfile
RUN npm ci --registry=https://registry.npmmirror.com
```

Commit and push Dockerfile changes made on the server back to the main repository.

## 8. Build and Start

```bash
cd /opt/Blog
docker compose --env-file .env.production -f compose.prod.yaml up -d --build
```

The first build downloads base images and dependencies and can take several minutes.

Check services and logs:

```bash
docker compose --env-file .env.production -f compose.prod.yaml ps
docker compose --env-file .env.production -f compose.prod.yaml logs --tail=100
docker compose --env-file .env.production -f compose.prod.yaml logs -f caddy
```

Expected services are `db`, `api`, `web`, and `caddy`; the database should be `healthy`.

## 9. Verify

After DNS propagation and Caddy certificate issuance:

```bash
curl https://www.zhangdongyu.vip/api/health
```

Expected response:

```json
{"status":"ok"}
```

Open `https://www.zhangdongyu.vip` in a browser.

## 10. Import Local Articles

`db/init/001_create_articles.sql` creates the production table on first startup, but does not copy local data.

Export locally:

```bash
docker compose exec -T db pg_dump -U blog -d blog \
  --data-only --table=articles > articles.sql
```

Upload and import:

```bash
scp articles.sql root@SERVER_PUBLIC_IP:/opt/Blog/

cd /opt/Blog
docker compose --env-file .env.production -f compose.prod.yaml \
  exec -T db psql -U blog -d blog < articles.sql
```

## 11. Update

```bash
cd /opt/Blog
git pull --ff-only
docker compose --env-file .env.production -f compose.prod.yaml up -d --build
```

## Troubleshooting

For Docker Hub timeouts, check `/etc/docker/daemon.json`, restart Docker, and retry. For slow Go builds, check memory with `free -h`; a small server may need temporary swap. For certificate errors, verify public DNS for both `@` and `www`, security-group ports `80/443`, and Caddy logs.
