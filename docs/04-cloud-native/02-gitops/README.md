# GitOps 实践

> 声明式持续交付

---

## 本章内容

- [ArgoCD 指南](./argocd-guide.md)

---

## GitOps 工具对比

| 特性 | ArgoCD | Flux |
|------|--------|------|
| UI | ✅ 完善 | ❌ 无原生 |
| 多集群 | ✅ 支持 | ✅ 支持 |
| 应用集 | ✅ ApplicationSet | ✅ Kustomize |
| 镜像更新 | ✅ 支持 | ✅ 支持 |

---

## ArgoCD vs Flux

- **ArgoCD**: 可视化强，企业级功能丰富
- **Flux**: 原生 GitOps，更轻量

---

## 相关文档

- [ArgoCD ApplicationSet](../../04-ecosystem/gitops-advanced/)
