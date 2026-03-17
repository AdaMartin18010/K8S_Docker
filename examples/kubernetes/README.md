# Kubernetes 示例

> K8s 资源清单与配置示例

---

## 目录结构

```
kubernetes/
├── 01-basic-resources/       # 基础资源
├── 02-deployment-patterns/   # 部署模式
├── 03-config-management/     # 配置管理
├── 04-networking/            # 网络配置
├── 05-storage/               # 存储配置
├── 06-security/              # 安全配置
├── 07-observability/         # 可观测性
├── 09-advanced-scheduling/   # 高级调度
├── 12-sidecar-native/        # 原生 Sidecar
├── 13-inplace-resizing/      # 原地资源调整
└── 14-user-namespaces/       # 用户命名空间
```

---

## 快速开始

```bash
# 部署基础资源
kubectl apply -f 01-basic-resources/

# 查看部署状态
kubectl get all

# 清理
kubectl delete -f 01-basic-resources/
```

---

## 示例类型

| 示例 | 说明 |
|------|------|
| Pod | 基础 Pod 配置 |
| Deployment | 无状态应用部署 |
| StatefulSet | 有状态应用部署 |
| Service | 服务暴露 |
| ConfigMap/Secret | 配置管理 |
| Ingress | 外部访问 |
| NetworkPolicy | 网络策略 |
| RBAC | 权限控制 |

---

## 相关文档

- [K8s 工作负载](../../docs/03-kubernetes/02-workloads/)
- [kubectl 速查表](../../docs/99-appendix/kubectl-cheatsheet.md)
