# K8S_Docker - 云原生技术知识体系

> 基于主题的 Docker & Kubernetes 知识库 (2025 Edition)

---

## 📚 项目简介

本项目是一个全面、系统的 Docker 和 Kubernetes 知识库，涵盖从基础概念到生产实践的全方位内容。

- **版本对齐**: Docker 28.x+ / Kubernetes 1.33-1.34
- **技术栈**: BuildKit 1.0、containerd 2.0、Gateway API v1.2+、eBPF/Cilium 1.18
- **实例语言**: Go

---

## 📁 项目结构

```text
K8S_Docker/
├── docs/                    # 📖 主题式文档 (130+ 文件)
│   ├── 00-overview/         # 项目概览与学习指南
│   ├── 01-fundamentals/     # 容器基础理论
│   ├── 02-docker/           # Docker 容器技术
│   ├── 03-kubernetes/       # Kubernetes 核心
│   ├── 04-ecosystem/        # 云原生生态
│   ├── 05-tools/            # 工具链
│   ├── 06-practices/        # 工程实践
│   └── 99-appendix/         # 附录与速查表
├── examples/                # 💻 代码示例 (210+ 文件)
│   ├── docker/              # Dockerfile最佳实践
│   ├── kubernetes/          # K8s 资源清单
│   ├── monitoring/          # 监控告警配置
│   ├── ci-cd/               # CI/CD 流水线
│   └── ...                  # 其他示例
└── docs-backup/             # 📦 归档的原始文档
```

---

## 🗺️ 文档导航

### 核心内容

| 主题 | 路径 | 说明 |
|------|------|------|
| **容器基础** | [docs/01-fundamentals/](docs/01-fundamentals/) | 容器技术概述、OCI标准、cgroups/namespace |
| **Docker** | [docs/02-docker/](docs/02-docker/) | Docker核心概念、Dockerfile、Compose、安全 |
| **Kubernetes** | [docs/03-kubernetes/](docs/03-kubernetes/) | K8s架构、工作负载、网络、存储、调度 |
| **云原生生态** | [docs/04-ecosystem/](docs/04-ecosystem/) | Dapr、Knative、eBPF、AI/ML、FinOps等 |
| **工具链** | [docs/05-tools/](docs/05-tools/) | Prometheus 3.0、OpenTelemetry、Dagger等 |
| **工程实践** | [docs/06-practices/](docs/06-practices/) | 混沌工程、性能测试、案例研究 |

### 附录

| 资源 | 路径 |
|------|------|
| **kubectl 速查表** | [docs/99-appendix/kubectl-cheatsheet.md](docs/99-appendix/kubectl-cheatsheet.md) |
| **Dockerfile 速查表** | [docs/99-appendix/dockerfile-cheatsheet.md](docs/99-appendix/dockerfile-cheatsheet.md) |
| **Helm 速查表** | [docs/99-appendix/helm-cheatsheet.md](docs/99-appendix/helm-cheatsheet.md) |
| **命令参考** | [docs/99-appendix/commands.md](docs/99-appendix/commands.md) |
| **术语表** | [docs/99-appendix/glossary.md](docs/99-appendix/glossary.md) |
| **知识图谱** | [docs/99-appendix/knowledge-graphs/](docs/99-appendix/knowledge-graphs/) |

---

## ✨ 2025 关键更新

### 新增特性

- ✅ **K8s 1.33/1.34**: Sidecar GA、用户命名空间默认启用、InPlacePodVerticalScaling Beta
- ✅ **containerd 2.0**: 用户命名空间、Transfer Service、NRI 默认启用
- ✅ **Gateway API v1.2/v1.3**: GA 稳定版，超时、重试、CORS、流量镜像
- ✅ **Cilium 1.18**: kube-proxy 替代 GA、Gateway API 支持、L2 Announcements
- ✅ **eBPF/Cilium**: 替代kube-proxy、L7网络策略
- ✅ **AI/ML**: KServe、vLLM、本地LLM部署
- ✅ **OpenTelemetry**: Logs GA、Metrics Stability 1.0
- ✅ **Prometheus 3.0**: OTLP 原生支持
- ✅ **混沌工程**: Chaos Mesh、故障注入
- ✅ **平台工程**: Backstage、Crossplane
- ✅ **供应链安全**: SLSA、Sigstore、Keyless Signing
- ✅ **FinOps**: 云成本管理

---

## 📖 推荐阅读路径

### 初学者路径

1. [容器基础概述](docs/01-fundamentals/container-overview.md)
2. [Docker 核心概念](docs/02-docker/01-core-concepts/)
3. [Dockerfile 最佳实践](docs/02-docker/02-dockerfile/best-practices.md)
4. [K8s 架构概览](docs/03-kubernetes/01-architecture/)
5. [K8s Pod 与 Sidecar](docs/03-kubernetes/01-pod/)

### 进阶路径

1. [K8s 高级调度](docs/03-kubernetes/08-advanced-scheduling/)
2. [eBPF 与 Cilium](docs/04-ecosystem/ebpf-cilium/)
3. [GitOps 实践](docs/04-ecosystem/gitops/)
4. [混沌工程](docs/06-practices/chaos-engineering/)

### 专家路径

1. [Crossplane 平台工程](docs/04-ecosystem/crossplane/)
2. [Knative Serverless](docs/04-ecosystem/knative/)
3. [Dapr 分布式运行时](docs/04-ecosystem/dapr/)
4. [供应链安全](docs/04-ecosystem/supply-chain-security/)

---

## 📊 内容统计

| 类别 | 数量 |
|------|------|
| 文档文件 | 130+ |
| 代码示例 | 210+ |
| 知识图谱 | 9 张 |
| 速查表 | 3 份 |
| 总大小 | 4.3 MB |

---

## 📝 许可证

[LICENSE](LICENSE)
