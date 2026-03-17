# Docker Compose 指南

> 多容器应用编排

---

## 简介

Docker Compose 是 Docker 官方的多容器编排工具，适合：

- 本地开发环境
- 测试环境
- 小型生产部署

---

## Compose 文件结构

```yaml
version: "3.9"

services:
  web:
    build: .
    ports:
      - "80:8080"
    depends_on:
      - db
      - redis
    environment:
      - DATABASE_URL=postgres://db/myapp

  db:
    image: postgres:16-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

---

## 常用命令

```bash
# 启动服务
docker compose up -d

# 查看日志
docker compose logs -f

# 停止服务
docker compose down

# 停止并删除卷
docker compose down -v

# 构建镜像
docker compose build

# 重启服务
docker compose restart web
```

---

## 生产环境配置

```yaml
services:
  web:
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
      restart_policy:
        condition: any
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
```
