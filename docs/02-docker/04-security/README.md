# Docker 安全加固

> 生产级容器安全指南

---

## 安全层次

```
┌─────────────────────────────────────────┐
│  1. 镜像安全                             │
│     - 基础镜像选择                       │
│     - 漏洞扫描                           │
├─────────────────────────────────────────┤
│  2. 构建安全                             │
│     - 多阶段构建                         │
│     - 密钥管理                           │
├─────────────────────────────────────────┤
│  3. 运行时安全                           │
│     - 非 root 用户                       │
│     - 只读文件系统                       │
│     - Capabilities 限制                  │
│     - Seccomp                            │
├─────────────────────────────────────────┤
│  4. 网络安全                             │
│     - 网络隔离                           │
│     - TLS 加密                           │
└─────────────────────────────────────────┘
```

---

## 运行时安全配置

```dockerfile
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app /app
USER 65534:65534
EXPOSE 8080
HEALTHCHECK CMD ["/app", "health"]
ENTRYPOINT ["/app"]
```

```bash
# 运行安全配置
docker run \
  --read-only \
  --security-opt=no-new-privileges:true \
  --cap-drop=ALL \
  --security-opt seccomp=profile.json \
  --user 65534:65534 \
  myapp:v1.0.0
```

---

## 安全扫描

```bash
# Trivy 扫描
trivy image myapp:v1.0.0

# Docker Scout
docker scout cves myapp:v1.0.0
```

---

## 关联代码

- [examples/docker/security/](../../examples/docker/security/)
