# Docker Compose 生产部署指南

> 使用 Docker Compose 部署生产环境

---

## 生产配置要点

### 1. 资源限制

```yaml
version: "3.9"

services:
  app:
    image: myapp:v1.0.0
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
      restart_policy:
        condition: any
        delay: 5s
        max_attempts: 3
      update_config:
        parallelism: 1
        delay: 10s
        failure_action: rollback
```

### 2. 健康检查

```yaml
services:
  app:
    image: myapp:v1.0.0
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
```

### 3. 安全配置

```yaml
services:
  app:
    image: myapp:v1.0.0
    read_only: true
    user: "1000:1000"
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
    tmpfs:
      - /tmp:noexec,nosuid,size=100m
```

---

## 多环境配置

### 基础配置 docker-compose.yml

```yaml
version: "3.9"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=production
  
  db:
    image: postgres:16-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password

volumes:
  postgres_data:

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

### 开发环境覆盖 docker-compose.override.yml

```yaml
version: "3.9"

services:
  app:
    build:
      target: development
    volumes:
      - .:/app
    environment:
      - ENV=development
      - DEBUG=true
  
  db:
    ports:
      - "5432:5432"
```

### 生产环境配置 docker-compose.prod.yml

```yaml
version: "3.9"

services:
  app:
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
    healthcheck:
      test: ["CMD", "/app", "health"]
      interval: 30s
```

---

## 部署命令

```bash
# 开发环境
docker compose up -d

# 生产环境
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f

# 滚动更新
docker compose up -d --no-deps --build app

# 回滚
docker compose down
docker compose up -d
```
