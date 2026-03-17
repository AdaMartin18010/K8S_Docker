# GitOps 实践指南

> 基于 Git 的持续交付模式

---

## 什么是 GitOps？

GitOps 是一种实现持续交付的方法，它将 Git 仓库作为唯一可信源，通过自动化工具同步集群状态。

```
┌───────────┐     push      ┌───────────┐
│  Developer│──────────────→│   Git     │
│           │               │  (Source) │
└───────────┘               └─────┬─────┘
                                  │
                                  │ webhook
                                  ↓
                          ┌───────────────┐
                          │  GitOps Agent │
                          │  (ArgoCD/Flux)│
                          └───────┬───────┘
                                  │ apply
                                  ↓
                          ┌───────────────┐
                          │  Kubernetes   │
                          │   Cluster     │
                          └───────────────┘
```

---

## GitOps 原则

1. **声明式配置**: 使用 YAML 描述期望状态
2. **版本控制**: 所有配置存储在 Git 中
3. **自动化同步**: 自动应用配置变更
4. **持续协调**: 持续对比实际状态与期望状态

---

## 主流工具

| 工具 | 特点 | 适用场景 |
|------|------|----------|
| **ArgoCD** | UI 丰富、功能强大 | 企业级应用 |
| **Flux** | 云原生、轻量级 | GitOps 原生 |
| **Rancher Fleet** | 多集群管理 | 大规模集群 |

---

## ArgoCD 核心概念

### Application

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: myapp
spec:
  project: default
  source:
    repoURL: https://github.com/org/repo.git
    targetRevision: HEAD
    path: k8s/overlays/prod
  destination:
    server: https://kubernetes.default.svc
    namespace: default
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

### App of Apps 模式

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: apps
spec:
  source:
    path: apps
  # 该目录下的每个 Application 都会被部署
```

---

## 推荐目录结构

```
gitops-repo/
├── apps/                    # ArgoCD Application 定义
├── base/                    # Kustomize 基础资源
│   ├── deployment.yaml
│   ├── service.yaml
│   └── kustomization.yaml
├── overlays/                # 环境配置
│   ├── dev/
│   ├── staging/
│   └── prod/
└── README.md
```

---

## 最佳实践

1. **单一可信源**: 使用一个 Git 仓库管理所有配置
2. **环境隔离**: 使用目录或分支区分环境
3. **权限控制**: Git 权限与 K8s RBAC 结合
4. **审计追踪**: 所有变更可追溯
5. **灾难恢复**: Git 仓库即备份
