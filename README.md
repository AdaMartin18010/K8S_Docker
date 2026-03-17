# K8S_Docker - 云原生技术知识体系

> 基于主题的 Docker & Kubernetes 知识库 (2025 Edition)

---

## 📚 项目简介

本项目是一个全面、系统的 Docker 和 Kubernetes 知识库，涵盖从基础概念到生产实践的全方位内容。

- **版本对齐**: Docker 25.x+ / Kubernetes 1.32-1.33
- **技术栈**: BuildKit 1.0、containerd 2.0、Gateway API、原生 Sidecar、eBPF/Cilium
- **实例语言**: Go

---

## 📁 项目结构

```text
K8S_Docker/
├── docs/                    # 📖 主题式文档 (3.5MB+, 90+ 文件)
│   ├── 00-overview/         # 项目概览与学习指南
│   ├── 01-fundamentals/     # 容器基础理论
│   ├── 02-docker/           # Docker 容器技术
│   ├── 03-kubernetes/       # Kubernetes 核心 (含1.33新特性)
│   ├── 04-ecosystem/        # 云原生生态
│   ├── 05-tools/            # 工具链
│   ├── 06-practices/        # 工程实践
│   └── 99-appendix/         # 附录
├── examples/                # 💻 代码示例 (90+ 文件)
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
│   ├── gateway-api/         # Gateway API
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

---

## ✨ 2025 关键更新

### 新增特性

- ✅ **K8s 1.33**: Sidecar GA、用户命名空间、DRA结构化参数
- ✅ **containerd 2.0**: Wasm支持、Transfer Service、Nerdctl v2
- ✅ **WebAssembly**: 冷启动快100倍、边缘计算新选择
- ✅ **eBPF/Cilium**: 替代kube-proxy、L7网络策略
- ✅ **AI/ML**: GPU调度、Kubeflow、KServe推理
- ✅ **沙箱安全**: gVisor、Kata、Firecracker
- ✅ **混沌工程**: Chaos Mesh、Litmus、故障注入
- ✅ **平台工程**: Backstage、Crossplane、IDP
- ✅ **供应链安全**: Sigstore、Cosign、SLSA、SBOM
- ✅ **Gateway API**: K8s新流量管理标准
- ✅ **vCluster**: 轻量级多租户方案
- ✅ **FinOps**: 云成本管理、OpenCost、Kubecost
- ✅ **Dapr**: 分布式应用运行时
- ✅ **OpenTelemetry**: 可观测性标准
- ✅ **Operator**: Kubebuilder开发指南
- ✅ **Knative**: Serverless平台
- ✅ **本地AI**: Ollama、vLLM部署

### 移除 (已废弃)

- ❌ PodSecurityPolicy (K8s 1.25 移除)
- ❌ Dockershim (K8s 1.24 移除)
- ❌ CephFS/RBD 内置插件

---

## 📊 内容统计

| 类别 | 数量 | 大小 |
|------|------|------|
| 文档文件 | 90+ | 3.5+ MB |
| 代码示例 | 90+ | 250+ KB |
| 主题分类 | 8 大主题 | 45+ 子目录 |

---

## 📖 推荐阅读路径

```text
1. 基础篇 (01-fundamentals)
   ├── 容器技术概述
   ├── OCI 开放容器标准
   ├── containerd 运行时生态
   └── 容器 vs 虚拟机

2. Docker 篇 (02-docker)
   ├── Docker 核心概念
   ├── Dockerfile 最佳实践
   ├── 镜像安全扫描
   └── Docker Compose

3. Kubernetes 篇 (03-kubernetes)
   ├── 核心概念
   ├── 工作负载管理
   ├── 网络与服务 (Gateway API)
   ├── 存储
   ├── 可观测性
   ├── 运维实践
   ├── Operator 开发
   ├── 高级调度 (DRA/多集群/边缘)
   └── K8s 1.33 新特性

4. 生态篇 (04-ecosystem)
   ├── WebAssembly
   ├── eBPF/Cilium
   ├── AI/ML 工作负载
   ├── Dapr 分布式运行时
   ├── FinOps 成本管理
   ├── Knative Serverless
   ├── 本地 AI/LLM 部署
   ├── 平台工程 (Backstage)
   ├── 供应链安全
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
   └── 开发工具链
```

---

## 🔗 相关资源

- [命令速查表](docs/99-appendix/commands.md)
- [术语表](docs/99-appendix/glossary.md)
- [资源推荐](docs/99-appendix/resources.md)

---

## 📝 许可证

[LICENSE](LICENSE)
