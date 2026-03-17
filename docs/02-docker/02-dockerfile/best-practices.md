# Dockerfile 最佳实践

> 编写生产级 Dockerfile 的完整指南

---

## 核心原则

### 1. 使用多阶段构建

```dockerfile
# 构建阶段
FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o /bin/app .

# 运行阶段
FROM alpine:3.19
COPY --from=builder /bin/app /app
CMD ["/app"]
```

### 2. 选择合适的基础镜像

| 镜像 | 大小 | 适用场景 |
|------|------|----------|
| `node:20` | 1.1GB | 开发环境 |
| `node:20-slim` | 240MB | 通用应用 |
| `node:20-alpine` | 180MB | 生产环境 |
| `gcr.io/distroless/static` | 50MB | 最小攻击面 |

### 3. 利用缓存优化构建顺序

```dockerfile
# 好：先复制依赖文件
COPY package.json pnpm-lock.yaml ./
RUN pnpm install

# 后复制源代码
COPY . .
RUN pnpm build
```

### 4. 使用非 root 用户

```dockerfile
# 创建用户
RUN adduser -D -u 1000 appuser
USER appuser
```

---

## BuildKit 高级特性

### 缓存挂载

```dockerfile
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
```

### 密钥挂载

```dockerfile
RUN --mount=type=secret,id=npmrc,target=/root/.npmrc \
    npm ci
```

### SSH 挂载

```dockerfile
RUN --mount=type=ssh \
    git clone git@github.com:org/private-repo.git
```

---

## 安全检查清单

- [ ] 使用非 root 用户
- [ ] 使用具体版本标签
- [ ] 最小化层数
- [ ] 使用 .dockerignore
- [ ] 扫描镜像漏洞
- [ ] 配置健康检查
