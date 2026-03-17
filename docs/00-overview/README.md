# Docker & Kubernetes 全景指南

> **文档版本**: 3.0  
> **最后更新**: 2025年3月  
> **适用版本**: Docker 25.x+, Kubernetes 1.30-1.32

---

## 📖 文档简介

本指南是一套完整的 Docker 与 Kubernetes 学习体系，涵盖从基础概念到生产实践的完整知识链。

### 为什么需要这份指南？

- **系统化学习路径**: 从入门到精通的渐进式学习路线
- **生产级内容**: 所有内容对齐 2025 年最新技术标准和最佳实践
- **理论与实践结合**: 每个主题都配有可运行的代码示例
- **主题化组织**: 按知识点层次组织，便于快速查阅

---

## 📂 文档结构

```
docs/
├── 00-overview/          # 本目录 - 概述与导读
├── 01-fundamentals/      # 基础概念 - 容器技术原理
├── 02-docker/           # Docker 专题 - 从入门到精通
├── 03-kubernetes/       # Kubernetes 专题 - 云原生编排
├── 04-cloud-native/     # 云原生生态 - 服务网格/GitOps/可观测性
├── 05-patterns/         # 设计模式 - 云原生架构模式
├── 06-practices/        # 实践指南 - 生产实战案例
└── 99-appendix/         # 附录 - 速查表/资源/术语表
```

---

## 🎯 学习路径

### 路径一：初学者 (4-6 周)

| 阶段 | 内容 | 预计时间 |
|------|------|----------|
| 第 1 周 | [01-fundamentals/](../01-fundamentals/) 容器基础 | 10h |
| 第 2 周 | [02-docker/](../02-docker/) Docker 核心 | 15h |
| 第 3 周 | [03-kubernetes/01-architecture/](../03-kubernetes/01-architecture/) K8s 架构 | 10h |
| 第 4 周 | [03-kubernetes/02-workloads/](../03-kubernetes/02-workloads/) 工作负载 | 15h |
| 第 5-6 周 | [06-practices/](../06-practices/) 实战项目 | 20h |

### 路径二：进阶工程师 (3-4 周)

| 阶段 | 内容 | 预计时间 |
|------|------|----------|
| 第 1 周 | [03-kubernetes/03-networking/](../03-kubernetes/03-networking/) 网络深入 | 12h |
| 第 2 周 | [03-kubernetes/04-storage/](../03-kubernetes/04-storage/) 存储管理 | 10h |
| 第 3 周 | [03-kubernetes/05-security/](../03-kubernetes/05-security/) 安全加固 | 12h |
| 第 4 周 | [04-cloud-native/](../04-cloud-native/) 云原生生态 | 15h |

### 路径三：架构师 (2-3 周)

| 阶段 | 内容 | 预计时间 |
|------|------|----------|
| 第 1 周 | [05-patterns/](../05-patterns/) 设计模式 | 15h |
| 第 2 周 | [03-kubernetes/06-operations/](../03-kubernetes/06-operations/) 运维实践 | 12h |
| 第 3 周 | [06-practices/](../06-practices/) 案例研究 | 10h |

---

## 📚 目录详解

### [01-fundamentals/](../01-fundamentals/) - 基础概念

容器技术的理论基础，包括：
- 容器 vs 虚拟机
- Linux Namespace 和 Cgroups
- OCI 标准规范
- 容器运行时原理

**目标读者**: 零基础或需要夯实基础的学习者

### [02-docker/](../02-docker/) - Docker 专题

Docker 的完整知识体系：
- **01-core-concepts/**: Docker 核心概念与架构
- **02-dockerfile/**: Dockerfile 编写艺术与最佳实践
- **03-compose/**: Docker Compose 多容器编排
- **04-security/**: 容器安全加固

**目标读者**: 需要使用 Docker 的开发者和运维人员

### [03-kubernetes/](../03-kubernetes/) - Kubernetes 专题

K8s 深度解析，分为 6 个模块：
- **01-architecture/**: 架构原理与核心组件
- **02-workloads/**: 工作负载资源 (Pod/Deployment/StatefulSet)
- **03-networking/**: 网络模型与服务发现
- **04-storage/**: 存储系统与数据管理
- **05-security/**: 安全体系与加固
- **06-operations/**: 运维与故障排查

**目标读者**: K8s 使用者、管理员、开发者

### [04-cloud-native/](../04-cloud-native/) - 云原生生态

现代云原生技术栈：
- **01-service-mesh/**: Istio/Linkerd 服务网格
- **02-gitops/**: ArgoCD/Flux GitOps 实践
- **03-observability/**: 可观测性三大支柱

**目标读者**: 云原生架构师、SRE

### [05-patterns/](../05-patterns/) - 设计模式

云原生架构模式：
- Sidecar 模式
- 微服务设计模式
- 部署模式 (金丝雀/蓝绿/滚动更新)
- 弹性模式 (熔断/重试/限流)

**目标读者**: 架构师、技术负责人

### [06-practices/](../06-practices/) - 实践指南

生产实战内容：
- 行业案例研究 (金融/电商/游戏)
- CI/CD 实践
- 性能调优
- 成本控制

**目标读者**: 需要落地实践的工程师

---

## 🔗 关联资源

### 代码示例

所有文档都关联了可运行的代码示例：

| 文档目录 | 代码位置 |
|----------|----------|
| `docs/02-docker/` | `examples/docker/` |
| `docs/03-kubernetes/` | `examples/kubernetes/` |
| `docs/03-kubernetes/02-workloads/` | `examples/go-client/` |

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

### 笔记建议

- 每个文档都包含**核心要点**总结
- 重要配置都有**YAML 示例**
- 常见错误都有**注意事项**提示

---

## 🆕 更新日志

### v3.0 (2025-03)
- 重构文档结构，主题化组织
- 添加 2025 年最新技术内容 (K8s 1.32, Docker 25.x)
- 整合 examples/ 代码库
- 删除过时内容 (PSP, Dockershim)

### v2.0 (2025-02)
- 更新技术内容至 2025 标准
- 新增 Docker BuildKit 1.0 内容
- 新增 Gateway API 内容

### v1.0 (2024)
- 初始版本

---

## 📄 版权声明

本文档采用 [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/) 协议。

欢迎贡献和反馈！
