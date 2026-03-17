# Docker & Kubernetes 全景指南

> **文档版本**: 3.0  
> **最后更新**: 2025年3月  
> **适用版本**: Docker 25.x+, Kubernetes 1.30-1.32

---

## 📖 文档简介

本指南是一套完整的 Docker 与 Kubernetes 学习体系，采用**主题化组织结构**，便于系统学习和快速查阅。

### 文档特点

- 🎯 **主题化组织**: 按知识点层次分类，便于定位
- 🔄 **2025 最新**: 对齐最新技术标准和最佳实践
- 🔗 **理论与实践结合**: 每个主题关联可运行代码示例
- 📚 **保留原始内容**: 归档内容可在 `docs-backup/` 查看

---

## 📂 文档结构

```
docs/
├── 00-overview/              # 概述与导读
│   ├── README.md             # 本文档
│   └── learning-path.md      # 学习路径指南
│
├── 01-fundamentals/          # 基础概念
│   ├── container-overview.md # 容器技术概述
│   ├── linux-namespace.md    # Linux Namespace
│   └── linux-cgroups.md      # Linux Cgroups
│
├── 02-docker/               # Docker 专题
│   ├── 01-core-concepts/    # Docker 核心概念
│   ├── 02-dockerfile/       # Dockerfile 指南
│   ├── 03-compose/          # Docker Compose
│   └── 04-security/         # 安全加固
│
├── 03-kubernetes/           # Kubernetes 专题
│   ├── 01-architecture/     # 架构原理
│   ├── 02-workloads/        # 工作负载
│   ├── 03-networking/       # 网络
│   ├── 04-storage/          # 存储
│   ├── 05-security/         # 安全
│   └── 06-operations/       # 运维
│
├── 04-cloud-native/         # 云原生生态
│   ├── 01-service-mesh/     # 服务网格
│   ├── 02-gitops/           # GitOps
│   └── 03-observability/    # 可观测性
│
├── 05-patterns/             # 设计模式
│   ├── sidecar-pattern.md
│   ├── microservices.md
│   └── deployment-strategies.md
│
├── 06-practices/            # 生产实践
│   ├── cicd-guide.md
│   ├── security-hardening.md
│   └── case-studies/
│
└── 99-appendix/             # 附录
    ├── cheatsheets.md       # 速查表
    ├── glossary.md          # 术语表
    └── resources.md         # 资源推荐
```

---

## 🚀 快速开始

### 选择你的学习路径

| 角色 | 推荐路径 | 预计时间 |
|------|----------|----------|
| **后端开发** | 01-fundamentals → 02-docker → 03-kubernetes/02-workloads | 5-6 周 |
| **DevOps 工程师** | 02-docker → 03-kubernetes → 04-cloud-native → 06-practices | 8-9 周 |
| **云原生架构师** | 01-fundamentals → 03-kubernetes → 04-cloud-native → 05-patterns | 10-12 周 |

详细学习路径请查看 [00-overview/learning-path.md](./00-overview/learning-path.md)

---

## 🔄 技术更新 (2025)

### Docker 25.x

- ✅ BuildKit 1.0 GA - 默认启用
- ✅ Docker Scout - 安全扫描集成
- ✅ Compose Watch - 开发模式热重载

### Kubernetes 1.30-1.32

- ✅ Sidecar 容器 GA (1.29+)
- ✅ 内存管理器 GA (1.32)
- ✅ QueueingHint 优化调度 (1.32)
- ✅ Gateway API v1.1 GA
- ❌ PodSecurityPolicy 已移除 (使用 Pod Security Standards)
- ❌ CephFS/RBD 内嵌插件已移除 (使用 CSI)

---

## 📚 关联资源

### 代码示例

| 文档 | 代码位置 |
|------|----------|
| `docs/02-docker/` | `examples/docker/` |
| `docs/03-kubernetes/` | `examples/kubernetes/` |
| `docs/05-patterns/` | `examples/kubernetes/02-deployment-patterns/` |
| `docs/06-practices/` | `examples/ci-cd/`, `examples/helm-charts/` |

### 外部资源

- [Docker 官方文档](https://docs.docker.com/)
- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [CNCF 云原生全景](https://landscape.cncf.io/)

---

## 📝 使用建议

1. **系统学习**: 按章节顺序阅读，配合代码实践
2. **快速查阅**: 使用 [99-appendix/cheatsheets.md](./99-appendix/cheatsheets.md) 速查表
3. **问题排查**: 参考 [06-practices/troubleshooting.md](./06-practices/troubleshooting.md)

---

## 📄 文档版本

### v3.0 (2025-03)
- 重构为主题化组织结构
- 对齐 2025 年最新技术标准
- 保留原始内容到 `docs-backup/`

### 归档内容

原始章节文件已归档到 `docs-backup/`:
- `chapter1_architecture.md` - 架构详解
- `chapter2_theory.md` - 数学理论
- `chapter3_patterns.md` - 设计模式
- `chapter4_cases.md` - 行业案例
- `chapter5_system_engineering.md` - 系统工程
- `chapter6_visualization.md` - 可视化

---

## 🤝 贡献

欢迎提交 Issue 和 PR 改进文档内容！
