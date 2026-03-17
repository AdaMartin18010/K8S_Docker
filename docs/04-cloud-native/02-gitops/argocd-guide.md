# ArgoCD GitOps 指南

> 使用 ArgoCD 实现声明式持续交付

---

## 什么是 GitOps？

GitOps 是一种使用 Git 作为唯一事实来源的运维模式：

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│   Git   │────▶│ ArgoCD  │────▶│   K8s   │
│ 仓库    │     │ 控制器  │     │ 集群    │
└─────────┘     └─────────┘     └─────────┘
     │                              │
     │◀─────────────────────────────┘
     │     自动同步/健康检查
```

---

## ArgoCD 核心概念

| 概念 | 说明 |
|------|------|
| **Application** | 关联 Git 仓库和 K8s 目标的 CRD |
| **AppProject** | 应用分组和权限控制 |
| **Repository** | Git 仓库配置 |
| **Cluster** | 目标 K8s 集群 |

---

## 安装 ArgoCD

```bash
# 创建命名空间
kubectl create namespace argocd

# 安装 ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# 暴露服务
kubectl port-forward svc/argocd-server -n argocd 8080:443

# 获取初始密码
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

---

## 创建应用

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: myapp
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/org/repo.git
    targetRevision: HEAD
    path: k8s/overlays/production
  destination:
    server: https://kubernetes.default.svc
    namespace: production
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

---

## 目录结构 (Kustomize)

```
k8s/
├── base/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── kustomization.yaml
└── overlays/
    ├── development/
    │   ├── replica-patch.yaml
    │   └── kustomization.yaml
    └── production/
        ├── replica-patch.yaml
        └── kustomization.yaml
```

---

## CLI 操作

```bash
# 登录
argocd login localhost:8080

# 添加仓库
argocd repo add https://github.com/org/repo.git

# 同步应用
argocd app sync myapp

# 查看应用状态
argocd app get myapp

# 查看应用差异
argocd app diff myapp
```
