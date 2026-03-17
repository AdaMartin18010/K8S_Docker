# GitOps 实践

> 声明式持续交付与基础设施管理

---

## 什么是 GitOps？

GitOps 是一种使用 Git 作为单一事实源的运维模式：

- **声明式**：系统状态在 Git 中定义
- **版本化**：所有变更可追溯
- **自动化**：自动同步与部署

---

## 工具对比

| 特性 | ArgoCD | Flux |
|------|--------|------|
| UI | ✅ 完善 | ❌ 无原生 |
| 多集群 | ✅ 支持 | ✅ 支持 |
| 应用集 | ✅ ApplicationSet | ✅ Kustomize |
| 镜像更新 | ✅ 支持 | ✅ 支持 |
| 通知 | ✅ 支持 | ✅ 支持 |

---

## ArgoCD 架构

```
┌─────────┐     ┌──────────┐     ┌─────────┐
│  Git    │────▶│ ArgoCD   │────▶│   K8s   │
│ (Source)│     │ (Sync)   │     │ (Target)│
└─────────┘     └──────────┘     └─────────┘
```

---

## 本章内容

- [ArgoCD 指南](./argocd-guide.md)

---

## GitOps 工作流

1. 开发者提交代码到 Git
2. CI 构建镜像并更新 Git
3. GitOps 工具检测到变更
4. 自动同步到 K8s 集群

---

## 相关文档

- [ArgoCD ApplicationSet](../../04-ecosystem/gitops-advanced/)
- [CI/CD 指南](../../06-practices/cicd-guide.md)
