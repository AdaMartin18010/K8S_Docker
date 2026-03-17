# K8S_Docker - 云原生技术知识体系

> 基于主题的 Docker & Kubernetes 知识库 (2025 Edition)

---

## 📚 项目简介

本项目是一个全面、系统的 Docker 和 Kubernetes 知识库，涵盖从基础概念到生产实践的全方位内容。

- **版本对齐**: Docker 28.x+ / Kubernetes 1.33-1.34
- **技术栈**: BuildKit 1.0、containerd 2.0、Gateway API v1.2+、eBPF/Cilium 1.18
- **实例语言**: Go
- **知识图谱**: 9张多维认知图 + 3大速查表

---

## 📁 项目结构

```text
K8S_Docker/
├── docs/                    # 📖 主题式文档 (4.0MB+, 130+ 文件)
│   ├── 00-overview/         # 项目概览与学习指南
│   ├── 01-fundamentals/     # 容器基础理论
│   ├── 02-docker/           # Docker 容器技术
│   ├── 03-kubernetes/       # Kubernetes 核心 (含1.33/1.34新特性)
│   ├── 04-ecosystem/        # 云原生生态
│   ├── 05-tools/            # 工具链
│   ├── 06-practices/        # 工程实践
│   └── 99-appendix/         # 附录
├── examples/                # 💻 代码示例 (210+ 文件)
│   ├── docker/              # Dockerfile最佳实践
│   ├── kubernetes/          # K8s 资源清单
│   ├── go-client/           # client-go / eBPF 示例
│   ├── helm/                # Helm Charts
│   ├── ci-cd/               # CI/CD 流水线
│   ├── chaos-mesh/          # 混沌工程实验
│   ├── backstage/           # 开发者门户
│   ├── cosign/              # 供应链安全
│   ├── dapr/                # Dapr 组件
│   ├── finops/              # 成本管理
│   ├── gateway-api/         # Gateway API v1.2+
│   ├── knative/             # Serverless
│   ├── local-llm/           # 本地AI部署
│   ├── opentelemetry/       # 可观测性
│   ├── operator/            # Operator示例
│   ├── vcluster/            # 虚拟集群
│   ├── wasm/                # WebAssembly
│   └── patterns/            # 设计模式
└── docs-backup/             # 📦 归档的原始文档
```

---

## 🗺️ 文档导航

### 核心内容

| 主题 | 文档 | 示例 |
|------|------|------|
| **容器基础** | [docs/01-fundamentals/](docs/01-fundamentals/) | - |
| **Docker** | [docs/02-docker/](docs/02-docker/) | [examples/docker/](examples/docker/) |
| **Kubernetes** | [docs/03-kubernetes/](docs/03-kubernetes/) | [examples/kubernetes/](examples/kubernetes/) |
| **云原生生态** | [docs/04-ecosystem/](docs/04-ecosystem/) | [examples/](examples/) |
| **工程实践** | [docs/06-practices/](docs/06-practices/) | [examples/patterns/](examples/patterns/) |
| **CI/CD** | [docs/04-ecosystem/gitops/](docs/04-ecosystem/gitops/) | [examples/ci-cd/](examples/ci-cd/) |
| **可观测性** | [docs/03-kubernetes/observability/](docs/03-kubernetes/observability/) | [examples/monitoring/](examples/monitoring/) |

---

## ✨ 2025 关键更新

### 新增特性

- ✅ **K8s 1.33/1.34**: Sidecar GA、用户命名空间默认启用、InPlacePodVerticalScaling Beta
- ✅ **containerd 2.0**: 用户命名空间、Transfer Service、NRI 默认启用
- ✅ **Gateway API v1.2/v1.3**: GA 稳定版，超时、重试、CORS、流量镜像
- ✅ **Cilium 1.18**: kube-proxy 替代 GA、Gateway API 支持、L2 Announcements
- ✅ **WebAssembly**: 冷启动快100倍、边缘计算新选择
- ✅ **eBPF/Cilium**: 替代kube-proxy、L7网络策略
- ✅ **AI/ML**: KServe v0.16、vLLM、TensorRT-LLM、Multi-Node Inference
- ✅ **OpenTelemetry**: Logs GA、Metrics Stability 1.0、GenAI 语义约定
- ✅ **Prometheus 3.0**: OTLP 原生支持、Native Histograms
- ✅ **沙箱安全**: gVisor、Kata、Firecracker
- ✅ **混沌工程**: Chaos Mesh、Litmus、故障注入
- ✅ **平台工程**: Backstage、Crossplane、IDP
- ✅ **供应链安全**: SLSA 1.2、Sigstore 普及、Keyless Signing
- ✅ **vCluster**: 轻量级多租户方案
- ✅ **FinOps**: 云成本管理、OpenCost、Kubecost
- ✅ **Dapr**: 分布式应用运行时
- ✅ **Knative**: Serverless平台
- ✅ **本地AI**: Ollama、vLLM、TensorRT-LLM 部署

### 移除 (已废弃)

- ❌ PodSecurityPolicy (K8s 1.25 移除)
- ❌ Dockershim (K8s 1.24 移除)
- ❌ CephFS/RBD 内置插件
- ❌ Endpoints API (K8s 1.33 废弃)
- ❌ gitRepo Volume (K8s 1.33 移除)

---

## 📊 内容统计

| 类别 | 数量 | 大小 |
|------|------|------|
| 文档文件 | 130+ | 4.0+ MB |
| 代码示例 | 210+ | 350+ KB |
| 主题分类 | 8 大主题 | 55+ 子目录 |
| 知识图谱 | 9 张 | 260 KB |
| 速查表 | 3 份 | 15 KB |

---

## 📖 推荐阅读路径

```text
1. 基础篇 (01-fundamentals)
   ├── 容器技术概述
   ├── OCI 开放容器标准
   ├── containerd 2.0 运行时生态
   └── 容器 vs 虚拟机

2. Docker 篇 (02-docker)
   ├── Docker 核心概念
   ├── Dockerfile 最佳实践
   ├── 镜像安全扫描
   └── Docker Compose

3. Kubernetes 篇 (03-kubernetes)
   ├── 核心概念
   ├── 工作负载管理
   ├── 网络与服务 (Gateway API v1.2+)
   ├── 存储
   ├── 可观测性 (OpenTelemetry)
   ├── 运维实践
   ├── Operator 开发
   ├── 高级调度 (DRA/多集群/边缘)
   └── K8s 1.33/1.34 新特性

4. 生态篇 (04-ecosystem)
   ├── WebAssembly
   ├── eBPF/Cilium 1.18
   ├── AI/ML 工作负载 (KServe v0.16)
   ├── Dapr 分布式运行时
   ├── FinOps 成本管理
   ├── Knative Serverless
   ├── 本地 AI/LLM 部署
   ├── 平台工程 (Backstage)
   ├── 供应链安全 (SLSA 1.2)
   ├── Service Mesh
   └── GitOps

5. 实践篇 (06-practices)
   ├── 设计模式
   ├── 混沌工程
   ├── 案例研究
   ├── 性能基准
   ├── 故障排查
   └── 容量规划

6. 工具篇 (05-tools)
   ├── OpenTelemetry 可观测性
   ├── Prometheus 3.0
   └── 开发工具链
```

---

## 🧠 多维知识图谱 (新增)

> 全新升级的7种思维表征方式，帮助建立直观、系统的云原生认知体系

| 图谱类型 | 文件 | 用途 |
|----------|------|------|
| 🗺️ **思维导图** | [知识全景图](docs/99-appendix/knowledge-graphs/mindmap-cloudnative.md) | 全局技术栈概览 |
| 📊 **对比矩阵** | [技术选型对比](docs/99-appendix/knowledge-graphs/comparison-matrix.md) | 8大领域技术对比 |
| 🌲 **决策树** | [场景决策路径](docs/99-appendix/knowledge-graphs/decision-trees.md) | 9大场景选型指南 |
| 🏗️ **架构图** | [系统架构详解](docs/99-appendix/knowledge-graphs/architecture-systems.md) | 6大架构深度解析 |
| 🌳 **应用树** | [场景方案映射](docs/99-appendix/knowledge-graphs/application-tree.md) | 业务场景-技术映射 |
| 🎓 **学习路径** | [技能发展路线](docs/99-appendix/knowledge-graphs/learning-pathways.md) | 5大角色学习指南 |
| 🔗 **概念网络** | [技术关联关系](docs/99-appendix/knowledge-graphs/concept-network.md) | 概念依赖与演进 |

**快速入口**: [知识图谱总览](docs/99-appendix/knowledge-graphs/README.md)

---

## 📋 速查表 (新增)

| 速查表 | 文件 | 内容 |
|--------|------|------|
| **kubectl** | [kubectl-cheatsheet.md](docs/99-appendix/kubectl-cheatsheet.md) | 集群操作、排障、高级技巧 |
| **Dockerfile** | [dockerfile-cheatsheet.md](docs/99-appendix/dockerfile-cheatsheet.md) | 构建优化、安全最佳实践 |
| **Helm** | [helm-cheatsheet.md](docs/99-appendix/helm-cheatsheet.md) | 模板语法、Chart 开发 |

---

## 🔗 相关资源

### 附录

- [命令速查表](docs/99-appendix/commands.md)
- [kubectl 速查表](docs/99-appendix/kubectl-cheatsheet.md)
- [Dockerfile 速查表](docs/99-appendix/dockerfile-cheatsheet.md)
- [Helm 速查表](docs/99-appendix/helm-cheatsheet.md)
- [术语表](docs/99-appendix/glossary.md)
- [资源推荐](docs/99-appendix/resources.md)

### 知识图谱

- [知识图谱总览](docs/99-appendix/knowledge-graphs/)
- [知识全景图](docs/99-appendix/knowledge-graphs/mindmap-cloudnative.md)
- [技术选型对比](docs/99-appendix/knowledge-graphs/comparison-matrix.md)
- [场景决策路径](docs/99-appendix/knowledge-graphs/decision-trees.md)
- [维护计划](docs/99-appendix/knowledge-graphs/MAINTENANCE_PLAN.md)

### 示例精选

- [Prometheus 告警规则](examples/monitoring/prometheus-rules/)
- [GitHub Actions 可复用工作流](examples/ci-cd/github-actions/reusable-workflows/)
- [Gateway API v1.2 示例](examples/gateway-api/)
- [Sidecar 容器模式](examples/kubernetes/sidecar-native/)

---

## 📝 许可证

[LICENSE](LICENSE)
