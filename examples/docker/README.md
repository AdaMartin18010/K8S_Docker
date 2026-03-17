# Docker 示例

> Dockerfile 最佳实践示例

---

## 目录结构

```
docker/
├── basic/          # 基础 Dockerfile 示例
├── multi-stage/    # 多阶段构建示例
├── compose/        # Docker Compose 示例
└── security/       # 安全加固示例
```

---

## 基础示例

### Go 应用多阶段构建

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

---

## 安全最佳实践

1. 使用非 root 用户
2. 使用多阶段构建减小镜像
3. 定期扫描漏洞
4. 使用 distroless 镜像

---

## 相关文档

- [Dockerfile 最佳实践](../../docs/02-docker/02-dockerfile/best-practices.md)
- [Docker 安全](../../docs/02-docker/04-security/)
