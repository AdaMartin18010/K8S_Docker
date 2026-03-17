# Dockerfile 安全最佳实践

> 构建安全的容器镜像

---

## 镜像安全扫描

### Docker Scout (2025)

```bash
# 扫描镜像
docker scout quickview myapp:latest
docker scout cves myapp:latest

# 在 CI/CD 中使用
docker scout cves --format sarif --output report.sarif myapp:latest
```

### Trivy

```bash
# 镜像扫描
trivy image myapp:latest

# 文件系统扫描
trivy fs .
```

---

## Dockerfile 安全实践

### 1. 使用非 root 用户

```dockerfile
FROM node:20-alpine
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001
USER nextjs
```

### 2. 最小化基础镜像

```dockerfile
# 推荐
cr.io/distroless/static:nonroot
alpine:latest
deian:slim

# 避免
ubuntu:latest  # 过大
centos:latest  # 维护状态
```

### 3. 固定镜像版本

```dockerfile
# 避免
FROM node:latest

# 推荐
FROM node:20.11.1-alpine3.19
```

### 4. 使用多阶段构建

```dockerfile
# 构建阶段
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp

# 运行阶段
cr.io/distroless/static:nonroot
COPY --from=builder /app/myapp /myapp
USER nonroot:nonroot
ENTRYPOINT ["/myapp"]
```

### 5. 正确处理密钥

```dockerfile
# ❌ 错误：密钥留在镜像层
ARG API_KEY
ENV API_KEY=$API_KEY

# ✅ 正确：使用 BuildKit 密钥挂载
RUN --mount=type=secret,id=api_key \
    API_KEY=$(cat /run/secrets/api_key) \
    go build .
```

---

## 运行时安全

### 只读根文件系统

```bash
docker run --read-only -v /tmp:/tmp:rw myapp
```

### 安全选项

```bash
docker run \
  --security-opt=no-new-privileges:true \
  --cap-drop=ALL \
  --cap-add=NET_BIND_SERVICE \
  myapp
```

---

## 安全清单

- [ ] 使用非 root 用户运行
- [ ] 使用最小化基础镜像
- [ ] 固定镜像标签
- [ ] 多阶段构建减少攻击面
- [ ] 定期扫描镜像漏洞
- [ ] 使用 SBOM 追踪依赖
- [ ] 启用只读根文件系统
- [ ] 最小化容器权限
