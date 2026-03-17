# Docker 完整指南

> **版本**: Docker 25.x+ | BuildKit 1.0 | 最后更新: 2025年3月

---

## 目录

- [1. 核心概念](#1-核心概念)
- [2. Dockerfile 最佳实践](#2-dockerfile-最佳实践)
- [3. 现代 Docker 工具链](#3-现代-docker-工具链)
- [4. Docker Compose](#4-docker-compose)
- [5. 安全与优化](#5-安全与优化)
- [6. 关联代码示例](#6-关联代码示例)

---

## 1. 核心概念

### 1.1 架构演进

```
传统架构 (Docker < 20.10)          现代架构 (Docker 25.x+)
┌─────────────┐                    ┌─────────────┐
│  Docker CLI  │                    │  Docker CLI  │
└──────┬──────┘                    └──────┬──────┘
       │                                  │
┌──────▼──────┐                    ┌──────▼──────┐
│  dockerd    │                    │  dockerd    │
│  (Legacy)   │                    │  + BuildKit │ ← 默认启用
└──────┬──────┘                    └──────┬──────┘
       │                                  │
┌──────▼──────┐                    ┌──────▼──────┐
│ containerd  │                    │ containerd  │
└──────┬──────┘                    └──────┬──────┘
       │                                  │
┌──────▼──────┐                    ┌──────▼──────┐
│    runc     │                    │    runc     │
└─────────────┘                    └─────────────┘
```

### 1.2 OCI 标准

- **Runtime Spec**: 容器运行时标准 (runc, crun)
- **Image Spec**: 镜像格式标准
- **Distribution Spec**: 镜像分发的标准

---

## 2. Dockerfile 最佳实践

### 2.1 现代 Dockerfile (2025标准)

```dockerfile
# syntax=docker/dockerfile:1.6  ← 启用最新语法

# ========== 构建阶段 ==========
FROM golang:1.22-alpine3.19 AS builder

WORKDIR /src

# 利用缓存：先复制依赖
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 复制源码并构建
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-w -s" -o /bin/app .

# ========== 安全扫描阶段 (可选) ==========
FROM builder AS security-scan
RUN --mount=type=cache,target=/tmp/trivy \
    command -v trivy && trivy fs --exit-code 0 /src || true

# ========== 生产阶段 ==========
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /bin/app /app

# 非 root 用户
USER 65534:65534

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app", "health"]

ENTRYPOINT ["/app"]
```

### 2.2 关键优化点

| 技术 | 说明 | 收益 |
|------|------|------|
| BuildKit 缓存挂载 | `--mount=type=cache` | 加速依赖下载 |
| 多阶段构建 | 分离构建和运行环境 | 减小镜像体积 |
| distroless 基础镜像 | 最小化攻击面 | 提高安全性 |
| 非 root 用户 | `USER nonroot` | 安全最佳实践 |
| 健康检查 | `HEALTHCHECK` | 运行时监控 |

### 2.3 常见反例

```dockerfile
# ❌ 错误示例
FROM golang:latest              # 使用 latest，不可重现
COPY . .                       # 未利用缓存
RUN go build                   # 每次都要下载依赖
RUN apt-get update && apt-get install -y curl vim  # 包含不必要工具
CMD ./app                      # shell 格式，信号处理有问题
```

---

## 3. 现代 Docker 工具链

### 3.1 Docker BuildKit 1.0 (2024)

```bash
# 启用 BuildKit（Docker 23+ 默认启用）
export DOCKER_BUILDKIT=1

# 使用缓存挂载
docker build \
  --cache-from type=registry,ref=myapp:buildcache \
  --cache-to type=registry,ref=myapp:buildcache,mode=max \
  -t myapp:latest .
```

### 3.2 Docker Scout (安全扫描)

```bash
# 快速概览
docker scout quickview myapp:latest

# 详细 CVE 分析
docker scout cves myapp:latest

# 与基础镜像对比
docker scout compare --to node:22-alpine myapp:latest

# CI/CD 集成
docker scout cves --only-severity critical,high --exit-code myapp:latest
```

### 3.3 Docker Compose Watch (开发)

```yaml
# docker-compose.yml
services:
  app:
    build: .
    develop:
      watch:
        - action: sync
          path: ./src
          target: /app/src
        - action: rebuild
          path: ./package.json
```

```bash
# 开发模式热重载
docker compose watch
```

---

## 4. Docker Compose

### 4.1 生产级配置

```yaml
version: "3.9"

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
      target: production

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

    healthcheck:
      test: ["CMD", "/app", "health"]
      interval: 30s
      timeout: 5s
      retries: 3

    security_opt:
      - no-new-privileges:true

    read_only: true
    user: "1000:1000"

  db:
    image: postgres:16-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password

volumes:
  postgres_data:
    driver: local

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

---

## 5. 安全与优化

### 5.1 安全清单

- [ ] 使用非 root 用户运行
- [ ] 使用只读根文件系统
- [ ] 删除所有 capabilities
- [ ] 使用 seccomp 配置文件
- [ ] 扫描镜像漏洞
- [ ] 不硬编码敏感信息
- [ ] 使用 .dockerignore

### 5.2 性能优化

```bash
# 并行构建
docker buildx build --parallel -t myapp .

# 多平台构建
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t myapp:latest \
  --push .
```

---

## 6. 关联代码示例

| 主题 | 文档位置 | 代码示例 |
|------|----------|----------|
| Dockerfile 最佳实践 | 2.1 节 | `examples/docker/basic/Dockerfile.good` |
| Dockerfile 反例 | 2.3 节 | `examples/docker/basic/Dockerfile.bad` |
| 多阶段构建 | 3.1 节 | `examples/docker/multi-stage/` |
| Docker Compose | 4 节 | `examples/docker/compose/` |
| 安全加固 | 5.1 节 | `examples/docker/security/` |

---

## 参考

- [Docker 官方最佳实践](https://docs.docker.com/develop/dev-best-practices/)
- [Docker Scout 文档](https://docs.docker.com/scout/)
- [BuildKit 高级功能](https://docs.docker.com/build/buildkit/)
