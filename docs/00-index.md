# Docker & Kubernetes 知识库索引

## 项目概述

> **版本**: 2025年3月 | **K8s 版本**: 1.33 | **Docker 版本**: 25.x+

本知识库包含 **96+ 篇技术文档** 和 **103+ 个代码示例**，覆盖 Docker、Kubernetes 及其云原生生态系统的完整技术栈。

## 快速导航

### Docker 基础 (01-docker)

| 主题 | 文档 | 示例 |
|------|------|------|
| 容器基础 | [容器基础概念](01-docker/01-container-basics/README.md) | - |
| 镜像构建 | [Dockerfile 最佳实践](01-docker/02-images/dockerfile-best-practices.md) | [Go 多阶段构建](examples/dockerfile/go-multi-stage.dockerfile) |
| 镜像优化 | [镜像优化策略](01-docker/02-images/image-optimization.md) | - |
| BuildKit | [BuildKit 高级用法](01-docker/02-images/buildkit-advanced.md) | - |
| 多平台构建 | [跨平台镜像构建](01-docker/02-images/multi-platform.md) | - |
| 容器运行时 | [containerd 详解](01-docker/03-runtime/containerd.md) | - |
| 容器网络 | [Docker 网络模型](01-docker/04-networking/README.md) | - |
| 容器存储 | [Docker 存储管理](01-docker/05-storage/README.md) | - |
| Compose | [Docker Compose](01-docker/06-compose/README.md) | [Compose 示例](examples/compose/docker-compose.yml) |
| 安全 | [容器安全](01-docker/07-security/README.md) | - |
| 调试 | [容器调试技巧](01-docker/08-debug/README.md) | - |

### Kubernetes 核心 (02-kubernetes-core)

| 主题 | 文档 | 示例 |
|------|------|------|
| K8s 架构 | [K8s 架构概览](02-kubernetes-core/01-architecture/README.md) | - |
| API Server | [API Server 深入解析](02-kubernetes-core/02-api-server/README.md) | - |
| etcd | [etcd 存储详解](02-kubernetes-core/03-etcd/README.md) | - |
| Scheduler | [调度器原理](02-kubernetes-core/04-scheduler/README.md) | - |
| Controller | [控制器模式](02-kubernetes-core/05-controller/README.md) | - |
| Kubelet | [Kubelet 详解](02-kubernetes-core/06-kubelet/README.md) | - |
| KubeProxy | [KubeProxy 网络代理](02-kubernetes-core/07-kube-proxy/README.md) | - |
| 1.33 新特性 | [K8s 1.33 新特性](02-kubernetes-core/08-1.33-features/README.md) | - |
| **Gateway API GA** | [Gateway API GA](02-kubernetes-core/09-gateway-api-ga/README.md) | - |

### Kubernetes 工作负载 (03-kubernetes)

| 主题 | 文档 | 示例 |
|------|------|------|
| Pod | [Pod 详解](03-kubernetes/01-pod/README.md) | [Pod 示例](examples/k8s-basic/pod.yaml) |
| **原生 Sidecar** | [原生 Sidecar 容器](03-kubernetes/01-pod/sidecar-native.md) | - |
| Deployment | [Deployment 管理](03-kubernetes/02-deployment/README.md) | [Deployment 示例](examples/k8s-basic/deployment.yaml) |
| StatefulSet | [StatefulSet 有状态应用](03-kubernetes/03-statefulset/README.md) | - |
| DaemonSet | [DaemonSet 守护进程](03-kubernetes/04-daemonset/README.md) | - |
| Job/CronJob | [批处理任务](03-kubernetes/05-jobs/README.md) | - |
| HPA/VPA | [自动扩缩容](03-kubernetes/06-hpa-vpa/README.md) | - |
| **Operators** | [Operator 开发](03-kubernetes/07-operators/README.md) | [CRD 示例](examples/operator/crd.yaml) |

### 云原生生态 (04-ecosystem)

| 主题 | 文档 | 示例 |
|------|------|------|
| **KubeVirt** | [KubeVirt 虚拟化](04-ecosystem/kubevirt/README.md) | [VM 示例](examples/kubevirt/vm.yaml) |
| **GitOps 高级** | [ArgoCD ApplicationSet](04-ecosystem/gitops-advanced/README.md) | [List Generator](examples/argocd-advanced/applicationset-list.yaml) |
| **零信任安全** | [零信任架构](04-ecosystem/zero-trust/README.md) | [SPIRE 部署](examples/zero-trust/spire-server.yaml) |
| Helm | [Helm 包管理](04-ecosystem/helm/README.md) | - |
| Kustomize | [Kustomize 配置管理](04-ecosystem/kustomize/README.md) | - |
| ArgoCD | [ArgoCD GitOps](04-ecosystem/argocd/README.md) | - |
| Flux | [Flux GitOps](04-ecosystem/flux/README.md) | - |
| Tekton | [Tekton CI/CD](04-ecosystem/tekton/README.md) | - |
| Jenkins | [Jenkins on K8s](04-ecosystem/jenkins/README.md) | - |
| **Knative** | [Knative Serverless](04-ecosystem/knative/README.md) | [Service 示例](examples/knative/service.yaml) |
| Service Mesh | [Istio 服务网格](04-ecosystem/istio/README.md) | - |
| **Dapr** | [Dapr 运行时](04-ecosystem/dapr/README.md) | [Workflow 示例](examples/dapr/workflow.go) |
| Cert Manager | [证书管理](04-ecosystem/cert-manager/README.md) | - |
| External DNS | [外部 DNS](04-ecosystem/external-dns/README.md) | - |
| Ingress | [Ingress 控制器](04-ecosystem/ingress/README.md) | - |
| Gateway API | [Gateway API](04-ecosystem/gateway-api/README.md) | - |

### 可观测性 (05-observability)

| 主题 | 文档 | 示例 |
|------|------|------|
| 监控基础 | [K8s 监控概述](05-observability/01-monitoring/README.md) | - |
| Prometheus | [Prometheus 监控](05-observability/02-prometheus/README.md) | - |
| Grafana | [Grafana 可视化](05-observability/03-grafana/README.md) | - |
| 日志 | [日志收集](05-observability/04-logging/README.md) | - |
| **OpenTelemetry** | [OpenTelemetry](05-observability/05-opentelemetry/README.md) | - |
| 追踪 | [分布式追踪](05-observability/06-tracing/README.md) | - |
| eBPF | [eBPF 可观测性](05-observability/07-ebpf/README.md) | - |

### 存储 (06-storage)

| 主题 | 文档 | 示例 |
|------|------|------|
| 存储基础 | [K8s 存储概述](06-storage/01-storage-basics/README.md) | - |
| CSI | [CSI 驱动](06-storage/02-csi/README.md) | - |
| Rook Ceph | [Rook Ceph 存储](06-storage/03-rook-ceph/README.md) | - |
| Longhorn | [Longhorn 存储](06-storage/04-longhorn/README.md) | - |

### 安全与合规 (07-security)

| 主题 | 文档 | 示例 |
|------|------|------|
| RBAC | [RBAC 权限管理](07-security/01-rbac/README.md) | - |
| Network Policy | [网络策略](07-security/02-network-policy/README.md) | - |
| Pod Security | [Pod 安全标准](07-security/03-pod-security/README.md) | - |
| OPA | [OPA 策略](07-security/04-opa/README.md) | - |
| Kyverno | [Kyverno 策略](07-security/05-kyverno/README.md) | - |
| **供应链安全** | [Supply Chain Security](07-security/06-supply-chain/README.md) | - |
| Secrets | [Secret 管理](07-security/07-secrets/README.md) | - |

### 多集群与联邦 (08-multicluster)

| 主题 | 文档 | 示例 |
|------|------|------|
| 多集群基础 | [多集群概述](08-multicluster/01-multicluster-basics/README.md) | - |
| Karmada | [Karmada 联邦](08-multicluster/02-karmada/README.md) | - |
| **Cluster API** | [Cluster API](08-multicluster/03-cluster-api/README.md) | [Cluster 示例](examples/cluster-api/cluster.yaml) |

### AI/ML 与边缘 (09-ai-edge)

| 主题 | 文档 | 示例 |
|------|------|------|
| Kubeflow | [Kubeflow ML平台](09-ai-edge/01-kubeflow/README.md) | - |
| KServe | [KServe 模型服务](09-ai-edge/02-kserve/README.md) | - |
| **Kubeflow MLOps** | [Kubeflow MLOps](09-ai-edge/03-kubeflow-mlops/README.md) | [Pipeline 示例](examples/kubeflow/pipeline.yaml) |
| **Local AI** | [本地LLM部署](04-ecosystem/local-ai/README.md) | [vLLM 部署](examples/local-llm/vllm-deployment.yaml) |
| KubeEdge | [KubeEdge 边缘计算](09-ai-edge/04-kubeedge/README.md) | - |

### 混沌工程 (10-chaos-engineering)

| 主题 | 文档 | 示例 |
|------|------|------|
| **Chaos Mesh** | [Chaos Mesh](10-chaos-engineering/01-chaos-mesh/README.md) | [实验示例](examples/chaos-mesh/pod-failure.yaml) |
| Litmus | [Litmus Chaos](10-chaos-engineering/02-litmus/README.md) | - |

### 平台工程 (11-platform-engineering)

| 主题 | 文档 | 示例 |
|------|------|------|
| **Backstage** | [Backstage IDP](11-platform-engineering/01-backstage/README.md) | - |
| Crossplane | [Crossplane 多云](11-platform-engineering/02-crossplane/README.md) | - |
| **FinOps** | [云成本管理](04-ecosystem/finops/README.md) | - |

### 新兴技术 (12-emerging)

| 主题 | 文档 | 示例 |
|------|------|------|
| **WebAssembly** | [WebAssembly + Spin](12-emerging/01-wasm-spin/README.md) | [WASM 部署](examples/wasm-spin/deployment.yaml) |
| **eBPF/Cilium** | [eBPF 网络](12-emerging/02-ebpf-cilium/README.md) | [Cilium 配置](examples/cilium/cilium-install.yaml) |
| Unikernel | [Unikernel](12-emerging/03-unikernel/README.md) | - |

## 2025 技术栈全景

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    2025 Kubernetes 技术栈全景                               │
├─────────────────────────────────────────────────────────────────────────────┤
│  基础设施层: K8s 1.33 | containerd 2.0 | CRI-O | KubeVirt | Cluster API     │
│  运行时层:   Docker 25.x | BuildKit 1.0 | WebAssembly | gVisor              │
│  网络层:     Cilium eBPF | Istio Ambient | Gateway API GA | DSR            │
│  存储层:     Rook Ceph | Longhorn | Volume Populator | SeaweedFS           │
│  可观测性:   OpenTelemetry GA | Prometheus 3.0 | Grafana | eBPF probes      │
│  GitOps:     ArgoCD 2.12 | ApplicationSet | Flux CD | Image Updater        │
│  Serverless: Knative 1.17 | KEDA | Dapr 1.16 | SpinKube                    │
│  AI/ML:      Kubeflow 1.11 | KServe | vLLM | Triton | Trainer 2.0          │
│  安全:       SPIFFE/SPIRE | Sigstore/Cosign | SLSA | Kyverno | Trivy        │
│  混沌工程:   Chaos Mesh 2.7 | Litmus 3.0                                   │
│  平台工程:   Backstage | Crossplane | Cluster API | IDP                    │
│  FinOps:     OpenCost | Kubecost | VPA | Spot instances                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 2025 重磅新特性

### Kubernetes 1.33 (代号: Octarine)

1. **原生 Sidecar 容器 GA** - 确定性启动顺序，Job 自动生命周期管理
2. **Gateway API GA** - Ingress 的继任者，标准化流量管理
3. **DRA (Dynamic Resource Allocation)** Beta - GPU/FPGA 动态分配
4. **原地 Pod 资源调整** - 不停机垂直扩缩容
5. **User Namespaces 默认启用** - 增强容器隔离

### 架构演进趋势

```
2024 → 2025 转变:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ingress → Gateway API GA
Sidecar Istio → Ambient Mesh GA
Prometheus-only → OpenTelemetry GA
Secrets → SPIFFE/SPIRE 零信任
VMware → KubeVirt 虚拟化
Docker → WebAssembly (边缘场景)
Manual scaling → KEDA + Dapr
单体 CI/CD → GitOps + Platform Engineering
云厂商绑定 → Crossplane + Cluster API 多云
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## 示例索引

```
examples/
├── argocd-advanced/      # ApplicationSet 多集群管理
├── chaos-mesh/           # 混沌工程实验
├── ci-cd/                # 持续集成/交付
├── cilium/               # Cilium eBPF 网络策略
├── cluster-api/          # Cluster API 集群配置
├── cosign/               # 镜像签名验证
├── dapr/                 # 分布式应用运行时
├── docker/               # Dockerfile 示例
├── finops/               # 成本管理
├── gateway-api/          # Gateway API 配置
├── go-client/            # Go K8s 客户端
├── helm-charts/          # Helm 图表
├── knative/              # Serverless 服务
├── kubernetes/           # K8s 基础资源
├── kubeflow/             # MLOps Pipeline
├── kubevirt/             # 虚拟机工作负载
├── local-llm/            # 本地大模型部署
├── microservices-demo/   # 微服务完整示例
├── opentelemetry/        # 可观测性埋点
├── operator/             # Operator 开发
├── vcluster/             # 虚拟集群
├── wasm-spin/            # WebAssembly on K8s
└── zero-trust/           # 零信任安全
```

## 使用指南

1. **初学者**: Docker 基础 → K8s 核心 → K8s 工作负载
2. **进阶用户**: GitOps → 服务网格 → 可观测性
3. **专家级**: Operators → 零信任 → 混沌工程

---
*最后更新: 2025年3月 | 文档数: 96+ | 示例数: 103+*
