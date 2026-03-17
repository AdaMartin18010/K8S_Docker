# Dockerfile 指南

> 编写生产级 Dockerfile 的完整指南

---

## 本章内容

1. [Dockerfile 基础语法](./syntax.md)
2. [最佳实践](./best-practices.md)
3. [多阶段构建](./multi-stage.md)
4. [BuildKit 高级特性](./buildkit-features.md)
5. [安全加固](./security-hardening.md)

---

## Dockerfile 示例模板

### Go 应用模板 (2025 标准)

```dockerfile
# syntax=docker/dockerfile:1.6
FROM golang:1.22-alpine3.19 AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

# 利用缓存
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 构建
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-w -s" -o /bin/app .

# 生产镜像
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/app /app
USER 65534:65534
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD ["/app", "health"]
ENTRYPOINT ["/app"]
```

### Node.js 应用模板

```dockerfile
# syntax=docker/dockerfile:1.6
FROM node:22-alpine AS base
RUN corepack enable && corepack prepare pnpm@latest --activate

FROM base AS deps
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile

FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN pnpm build && pnpm prune --prod

FROM node:22-alpine AS production
RUN apk add --no-cache dumb-init
WORKDIR /app
COPY --from=builder --chown=node:node /app/dist ./dist
COPY --from=builder --chown=node:node /app/node_modules ./node_modules
USER node
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s CMD node -e "require('http').get('http://localhost:3000/health')"
ENTRYPOINT ["dumb-init", "--"]
CMD ["node", "dist/main.js"]
```

---

## 关键指令速查

| 指令 | 说明 | 示例 |
|------|------|------|
| `FROM` | 基础镜像 | `FROM golang:1.22-alpine` |
| `RUN` | 执行命令 | `RUN go build -o app .` |
| `COPY` | 复制文件 | `COPY --chown=user:group src/ dst/` |
| `ADD` | 复制(支持URL) | `ADD https://example.com/file.tar.gz /` |
| `WORKDIR` | 工作目录 | `WORKDIR /app` |
| `USER` | 切换用户 | `USER 65534:65534` |
| `EXPOSE` | 暴露端口 | `EXPOSE 8080` |
| `ENTRYPOINT` | 入口点 | `ENTRYPOINT ["/app"]` |
| `CMD` | 默认命令 | `CMD ["--help"]` |
| `HEALTHCHECK` | 健康检查 | `HEALTHCHECK CMD curl -f /health` |

---

## 关联代码

- [examples/docker/basic/Dockerfile.good](../../examples/docker/basic/Dockerfile.good)
- [examples/docker/basic/Dockerfile.bad](../../examples/docker/basic/Dockerfile.bad)
- [examples/docker/multi-stage/Dockerfile](../../examples/docker/multi-stage/Dockerfile)
