# Docker & Kubernetes 全景指南

> **文档版本**: 2025 Edition
> **最后更新**: 2025年3月
> **适用版本**: Docker 28.x+, Kubernetes 1.33-1.34

---

## 📖 文档简介

本指南是一套完整的 Docker 与 Kubernetes 学习体系，涵盖从基础概念到生产实践的完整知识链。

### 为什么需要这份指南？

- **系统化学习路径**: 从入门到精通的渐进式学习路线
- **生产级内容**: 所有内容对齐 2025 年最新技术标准和最佳实践
- **理论与实践结合**: 每个主题都配有可运行的代码示例
- **主题化组织**: 按知识点层次组织，便于快速查阅

---

## 📂 实际文档结构

```
docs/
├── 00-overview/          # 本目录 - 概述与导读
│   ├── README.md
│   ├── learning-path.md
│   └── quickstart.md
├── 01-fundamentals/      # 基础概念 - 容器技术原理
│   ├── container-overview.md
│   ├── container-vs-vm.md
│   ├── oci-standard.md
│   ├── linux-namespace.md
│   ├── linux-cgroups.md
│   ├── cgroups-namespaces.md
│   └── containerd-runtimes.md
├── 02-docker/           # Docker 专题
│   ├── 01-core-concepts/
│   │   ├── docker-architecture.md
│   │   ├── networking.md
│   │   ├── storage.md
│   │   └── README.md
│   ├── 02-dockerfile/
│   │   ├── best-practices.md
│   │   ├── security.md
│   │   └── README.md
│   ├── 03-compose/
│   │   ├── production-guide.md
│   │   └── README.md
│   └── 04-security/
│       └── README.md
├── 03-kubernetes/       # Kubernetes 专题
│   ├── 01-architecture/
│   ├── 01-pod/
│   │   └── sidecar-native.md
│   ├── 02-workloads/
│   ├── 03-networking/
│   ├── 04-storage/
│   ├── 05-observability/
│   ├── 05-security/
│   ├── 06-operations/
│   ├── 07-operators/
│   ├── 08-advanced-scheduling/
│   └── whats-new-1.33.md
├── 04-ecosystem/        # 云原生生态
│   ├── ai-ml/, crossplane/, dapr/, dragonflydb/
│   ├── ebpf-cilium/, finops/, flagger/, gitops/
│   ├── keda/, keptn/, knative/, kubevirt/, kueue/
│   ├── local-ai/, nats/, openfeature/, platform-engineering/
│   ├── security/, service-mesh/, supply-chain-security/
│   └── wasmcloud/, webassembly/, zero-trust/
├── 05-tools/            # 工具链
├── 05-patterns/         # 设计模式
├── 06-practices/        # 工程实践
├── 06-storage/          # 存储
├── 07-security/         # 安全
├── 08-multicluster/     # 多集群
├── 09-ai-edge/          # AI与边缘
├── 12-emerging/         # 新兴技术
└── 99-appendix/         # 附录
```

---

## 🎯 学习路径

### 路径一：初学者 (4-6 周)

| 阶段 | 内容 | 预计时间 |
|------|------|----------|
| 第 1 周 | [容器基础](../01-fundamentals/container-overview.md) | 10h |
| 第 2 周 | [Docker 核心](../02-docker/01-core-concepts/) | 15h |
| 第 3 周 | [K8s 架构](../03-kubernetes/01-architecture/) | 10h |
| 第 4 周 | [工作负载](../03-kubernetes/02-workloads/) | 15h |
| 第 5-6 周 | [工程实践](../06-practices/) | 20h |

### 路径二：进阶工程师 (3-4 周)

| 阶段 | 内容 | 预计时间 |
|------|------|----------|
| 第 1 周 | [网络深入](../03-kubernetes/03-networking/) | 12h |
| 第 2 周 | [存储管理](../03-kubernetes/04-storage/) | 10h |
| 第 3 周 | [安全加固](../03-kubernetes/05-security/) | 12h |
| 第 4 周 | [云原生生态](../04-ecosystem/) | 15h |

### 路径三：架构师 (2-3 周)

| 阶段 | 内容 | 预计时间 |
|------|------|----------|
| 第 1 周 | [设计模式](../05-patterns/) | 15h |
| 第 2 周 | [运维实践](../03-kubernetes/06-operations/) | 12h |
| 第 3 周 | [案例研究](../06-practices/case-studies/) | 10h |

---

## 📚 目录详解

### [01-fundamentals/](../01-fundamentals/) - 基础概念

容器技术的理论基础，包括：

- [容器 vs 虚拟机](../01-fundamentals/container-vs-vm.md)
- [Linux Namespace](../01-fundamentals/linux-namespace.md)
- [Linux Cgroups](../01-fundamentals/linux-cgroups.md)
- [OCI 标准](../01-fundamentals/oci-standard.md)

**目标读者**: 零基础或需要夯实基础的学习者

### [02-docker/](../02-docker/) - Docker 专题

Docker 的完整知识体系：

- **核心概念**: [架构](../02-docker/01-core-concepts/docker-architecture.md)、[网络](../02-docker/01-core-concepts/networking.md)、[存储](../02-docker/01-core-concepts/storage.md)
- **Dockerfile**: [最佳实践](../02-docker/02-dockerfile/best-practices.md)、[安全](../02-docker/02-dockerfile/security.md)
- **Compose**: [生产指南](../02-docker/03-compose/production-guide.md)
- **安全**: [容器安全](../02-docker/04-security/)

**目标读者**: 需要使用 Docker 的开发者和运维人员

### [03-kubernetes/](../03-kubernetes/) - Kubernetes 专题

K8s 深度解析：

- **架构**: [核心组件](../03-kubernetes/01-architecture/)
- **Pod**: [原生 Sidecar](../03-kubernetes/01-pod/sidecar-native.md)
- **工作负载**: [Deployment/StatefulSet](../03-kubernetes/02-workloads/)
- **网络**: [Service/Ingress/Gateway API](../03-kubernetes/03-networking/)
- **存储**: [PV/PVC/StorageClass](../03-kubernetes/04-storage/)
- **可观测性**: [监控/日志/追踪](../03-kubernetes/05-observability/)
- **安全**: [RBAC/NetworkPolicy](../03-kubernetes/05-security/)
- **运维**: [排障/性能调优](../03-kubernetes/06-operations/)
- **新特性**: [K8s 1.33](../03-kubernetes/whats-new-1.33.md)

**目标读者**: K8s 使用者、管理员、开发者

### [04-ecosystem/](../04-ecosystem/) - 云原生生态

现代云原生技术栈：

- **服务网格**: [Istio/Cilium](../04-ecosystem/service-mesh/)
- **GitOps**: [ArgoCD](../04-ecosystem/gitops/)
- **可观测性**: [OpenTelemetry](../05-tools/observability/)
- **Serverless**: [Knative](../04-ecosystem/knative/)
- **运行时**: [Dapr](../04-ecosystem/dapr/)
- **AI/ML**: [本地LLM](../04-ecosystem/local-ai/)、[AI/ML平台](../04-ecosystem/ai-ml/)

**目标读者**: 云原生架构师、SRE

### [05-patterns/](../05-patterns/) - 设计模式

云原生架构模式：

- [Sidecar 模式](../05-patterns/sidecar-pattern.md)
- [微服务模式](../05-patterns/microservices.md)
- [部署模式](../05-patterns/)

**目标读者**: 架构师、技术负责人

### [06-practices/](../06-practices/) - 工程实践

生产实战内容：

- [混沌工程](../06-practices/chaos-engineering/)
- [CI/CD 指南](../06-practices/cicd-guide.md)
- [案例研究](../06-practices/case-studies/)

**目标读者**: 需要落地实践的工程师

---

## 🔗 关联资源

### 代码示例

| 文档目录 | 代码位置 |
|----------|----------|
| `docs/02-docker/` | `examples/docker/` |
| `docs/03-kubernetes/` | `examples/kubernetes/` |
| `docs/04-ecosystem/` | `examples/` 各子目录 |

### 速查表

- [kubectl 速查表](../99-appendix/kubectl-cheatsheet.md)
- [Dockerfile 速查表](../99-appendix/dockerfile-cheatsheet.md)
- [Helm 速查表](../99-appendix/helm-cheatsheet.md)

### 外部资源

- [Docker 官方文档](https://docs.docker.com/)
- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [CNCF 云原生全景](https://landscape.cncf.io/)

---

## 📝 使用建议

### 阅读方式

1. **顺序阅读**: 按章节顺序系统学习
2. **按需查阅**: 使用速查表快速定位知识点
3. **实践结合**: 每读完一章，动手实践配套代码

---

## 🆕 更新日志

### 2025 Edition (2025-03)

- 对齐 2025 技术栈 (K8s 1.33, Docker 28.x)
- 修复所有无效链接
- 删除虚假索引文件
- 补充缺失的核心文档
- 更新 Gateway API v1.2+ 内容
