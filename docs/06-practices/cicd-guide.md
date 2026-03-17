# CI/CD 实践指南

> 容器化应用的持续集成与交付

---

## CI/CD 流程

```
┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
│  Code   │──▶│  Build  │──▶│  Test   │──▶│  Scan   │──▶│ Deploy  │
│  Commit │   │  Image  │   │  Unit   │   │ Security│   │ to K8s  │
└─────────┘   └─────────┘   └─────────┘   └─────────┘   └─────────┘
```

---

## GitHub Actions 工作流

```yaml
name: CI/CD
on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build
        uses: docker/build-push-action@v5
        with:
          context: .
          push: false
          tags: myapp:${{ github.sha }}

      - name: Scan
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: myapp:${{ github.sha }}

      - name: Push
        run: |
          docker push myapp:${{ github.sha }}

      - name: Deploy
        run: |
          kubectl set image deployment/myapp app=myapp:${{ github.sha }}
          kubectl rollout status deployment/myapp
```

---

## GitOps 模式

### ArgoCD 应用定义

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: myapp
spec:
  project: default
  source:
    repoURL: https://github.com/org/repo
    targetRevision: HEAD
    path: k8s
  destination:
    server: https://kubernetes.default.svc
    namespace: production
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

---

## 最佳实践

1. **镜像标签策略**: 使用 Git SHA 或语义化版本
2. **安全扫描**: 每次构建都进行漏洞扫描
3. **多环境**: 开发/测试/生产环境分离
4. **回滚**: 保留历史版本，支持快速回滚
