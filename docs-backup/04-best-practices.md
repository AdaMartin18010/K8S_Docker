# 生产最佳实践

---

## Docker 最佳实践

### 镜像构建

| ✅ 应该 | ❌ 不应该 |
|--------|----------|
| 使用多阶段构建 | 单阶段包含所有构建工具 |
| 使用具体版本标签 | 使用 `latest` 标签 |
| 非 root 用户运行 | 以 root 运行 |
| 使用 .dockerignore | 复制不必要的文件 |
| 最小基础镜像 | 庞大的基础镜像 |
| 健康检查 | 无健康检查 |

### Dockerfile 模板

```dockerfile
# syntax=docker/dockerfile:1.6
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /bin/app /app
USER 65534:65534
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD ["/app", "health"]
ENTRYPOINT ["/app"]
```

---

## Kubernetes 最佳实践

### Pod 安全

```yaml
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    seccompProfile:
      type: RuntimeDefault
  containers:
    - name: app
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities:
          drop: [ALL]
```

### 资源管理

```yaml
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

### 健康检查

```yaml
startupProbe:
  httpGet:
    path: /health/startup
    port: 8080
  failureThreshold: 30
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  periodSeconds: 10
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  periodSeconds: 5
```

---

## 安全检查清单

### 容器安全

- [ ] 使用非 root 用户
- [ ] 只读根文件系统
- [ ] 删除所有 capabilities
- [ ] 使用 seccomp
- [ ] 扫描镜像漏洞
- [ ] 资源限制

### K8s 安全

- [ ] RBAC 最小权限
- [ ] NetworkPolicy
- [ ] Pod Security Standards
- [ ] Secret 加密
- [ ] API Server 审计日志

---

## 反例

详见 `examples/anti-patterns/` 目录，包含：

- Dockerfile 安全风险和性能问题
- Kubernetes 配置错误示例

---

## 参考

- [CIS Docker Benchmark](https://www.cisecurity.org/benchmark/docker)
- [CIS Kubernetes Benchmark](https://www.cisecurity.org/benchmark/kubernetes)
