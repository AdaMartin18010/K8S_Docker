# Docker 基础示例

本目录包含 Docker 基础最佳实践示例。

## 文件说明

| 文件 | 说明 |
|------|------|
| `Dockerfile.good` | ✅ 最佳实践示例 |
| `Dockerfile.bad` | ❌ 常见错误示例（反例）|
| `main.go` | 示例 Go 应用程序 |
| `go.mod` | Go 模块定义 |
| `.dockerignore` | Docker 构建忽略文件 |

## 构建与运行

### 构建镜像

```bash
# 使用 BuildKit 构建
docker build -f Dockerfile.good -t myapp:good .

# 不使用 BuildKit（对比）
DOCKER_BUILDKIT=0 docker build -f Dockerfile.good -t myapp:legacy .
```

### 运行容器

```bash
docker run -d -p 8080:8080 --name myapp myapp:good
```

### 查看差异

```bash
# 对比镜像大小
docker images myapp:good myapp:bad

# 查看构建历史
docker history myapp:good
```

## 最佳实践要点

1. **多阶段构建**：分离构建和运行环境
2. **最小基础镜像**：使用 distroless 或 alpine
3. **非 root 用户**：使用 USER 指令
4. **缓存优化**：合理排序 COPY 指令
5. **安全扫描**：集成 Trivy/Docker Scout
6. **健康检查**：配置 HEALTHCHECK
