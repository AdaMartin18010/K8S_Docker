# Kubernetes 核心

> Kubernetes 1.33-1.34 完整指南

---

## 本章内容

| 章节 | 说明 |
|------|------|
| [架构原理](./01-architecture/) | K8s 核心组件与架构 |
| [Pod](./01-pod/) | Pod 与原生 Sidecar |
| [工作负载](./02-workloads/) | Deployment、StatefulSet 等 |
| [网络](./03-networking/) | Service、Ingress、Gateway API |
| [存储](./04-storage/) | PV、PVC、StorageClass |
| [可观测性](./05-observability/) | 监控、日志、追踪 |
| [安全](./05-security/) | RBAC、NetworkPolicy |
| [运维](./06-operations/) | 排障、性能调优 |
| [Operator](./07-operators/) | 自定义资源与控制 |
| [高级调度](./08-advanced-scheduling/) | DRA、多集群、边缘计算 |
| [新特性](./whats-new-1.33.md) | K8s 1.33/1.34 更新 |

---

## 2025 关键更新

- **Sidecar GA**: 原生 Sidecar 容器支持
- **Gateway API GA**: 替代 Ingress 的标准 API
- **InPlacePodVerticalScaling**: 原地 Pod 资源调整
- **User Namespaces**: 增强容器隔离

---

## 学习路径

1. **入门**: 架构 → Pod → 工作负载 → 网络
2. **进阶**: 存储 → 安全 → 可观测性
3. **专家**: Operator → 高级调度 → 多集群
