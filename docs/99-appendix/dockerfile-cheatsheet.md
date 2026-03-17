# Dockerfile 速查表

> Docker 镜像构建最佳实践与指令参考 (Docker 28.x)

---

## 基础指令

```dockerfile
# 基础镜像
FROM ubuntu:24.04
FROM node:20-alpine
FROM golang:1.23 AS builder

# 标签
LABEL maintainer="user@example.com"
LABEL version="1.0"

# 工作目录
WORKDIR /app

# 复制文件
COPY . .
COPY file.txt /dest/
COPY --from=builder /app/build /app/

# 添加文件（支持URL和tar）
ADD https://example.com/file.tar.gz /dest/

# 环境变量
ENV NODE_ENV=production
ENV PORT=8080

# 构建参数
ARG VERSION=1.0
RUN echo $VERSION

# 执行命令
RUN apt-get update && apt-get install -y curl
RUN --mount=type=cache,target=/var/cache/apt apt-get update

# 暴露端口
EXPOSE 8080

# 数据卷
VOLUME ["/data"]

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# 入口点
ENTRYPOINT ["/app/start.sh"]
CMD ["--port", "8080"]

# 用户
USER 1000:1000
```

---

## 多阶段构建

```dockerfile
# 构建阶段
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

---

## 语言特定模板

### Node.js

```dockerfile
FROM node:20-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:20-alpine
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
USER node
EXPOSE 3000
CMD ["node", "index.js"]
```

### Python

```dockerfile
FROM python:3.12-slim AS builder
WORKDIR /app
RUN pip install --user --no-cache-dir -r requirements.txt

FROM python:3.12-slim
WORKDIR /app
COPY --from=builder /root/.local /root/.local
COPY . .
ENV PATH=/root/.local/bin:$PATH
EXPOSE 8000
CMD ["python", "app.py"]
```

### Java

```dockerfile
FROM maven:3.9-eclipse-temurin-21 AS build
WORKDIR /app
COPY pom.xml .
RUN mvn dependency:go-offline
COPY src ./src
RUN mvn package -DskipTests

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=build /app/target/*.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
```

---

## BuildKit 特性

```dockerfile
# 缓存挂载
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install -y python3

# 密钥挂载（不进入镜像层）
RUN --mount=type=secret,id=npmrc,target=/root/.npmrc \
    npm ci

# SSH 代理
RUN --mount=type=ssh git clone git@github.com:user/repo.git

# 绑定挂载
RUN --mount=type=bind,source=./config,target=/config \
    cp /config/app.conf /etc/
```

---

## 安全最佳实践

```dockerfile
# 使用非 root 用户
RUN useradd -m -s /bin/bash appuser
USER appuser

# 使用 distroless 镜像
FROM gcr.io/distroless/nodejs20-debian12
COPY --chown=nonroot:nonroot . /app
USER nonroot

# 最小权限
RUN chmod 755 /app && chmod 644 /app/config

# 扫描镜像
# docker scan myimage:latest
```

---

## 常用命令

```bash
# 构建
docker build -t myapp:latest .
docker build -t myapp:latest -f Dockerfile.prod .
docker build --no-cache -t myapp:latest .
docker build --target builder -t myapp:builder .

# 多平台构建
docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 -t myapp:latest --push .

# 缓存
docker build --cache-from myapp:cache -t myapp:latest .
docker build --build-arg VERSION=2.0 -t myapp:latest .

# 查看历史
docker history myapp:latest

# 导出导入
docker save myapp:latest > myapp.tar
docker load < myapp.tar
```
