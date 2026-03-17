# 云原生应用场景知识树

> 从业务场景出发，找到对应的技术解决方案

---

## 1. 微服务应用部署场景树

```
场景: 部署微服务应用到 Kubernetes
│
├─► 服务间通信需求?
│   │
│   ├─► 简单 HTTP ──► 方案
│   │   │
│   │   └─► K8s Service (ClusterIP) + Ingress
│   │
│   ├─► 需要 gRPC ──► 方案
│   │   │
│   │   ├─► 无高级需求 ──► K8s Service + Ingress (支持 gRPC)
│   │   │
│   │   └─► 需要负载均衡/熔断 ──► Service Mesh (Istio/Linkerd)
│   │
│   └─► 需要事件驱动 ──► 方案
│       │
│       ├─► 简单队列 ──► Redis Pub/Sub / NATS
│       │
│       └─► 企业级消息 ──► Kafka / RabbitMQ Operator
│
├─► 配置管理需求?
│   │
│   ├─► 简单配置 ──► ConfigMap + Secret
│   │
│   ├─► 多环境配置 ──► Kustomize (Overlay)
│   │
│   └─► 敏感配置 ──► External Secrets + Vault
│
├─► 可观测性需求?
│   │
│   ├─► 基础监控 ──► Prometheus + Grafana
│   │
│   ├─► 分布式追踪 ──► OpenTelemetry + Jaeger
│   │
│   └─► 全面可观测 ──► OpenTelemetry (Metrics/Logs/Traces) + Grafana
│
└─► CI/CD 需求?
    │
    ├─► 简单部署 ──► GitHub Actions + kubectl
    │
    └─► GitOps ──► ArgoCD + Argo Rollouts (渐进式交付)

推荐组合 (2025):
├── Service Mesh: Istio Ambient (通信)
├── Config: Kustomize + External Secrets
├── Observability: OpenTelemetry + Grafana
├── CI/CD: ArgoCD (GitOps)
└── Gateway: Gateway API (Cilium/Envoy)
```

---

## 2. AI/ML 模型服务化场景树

```
场景: 将 AI 模型部署为在线服务
│
├─► 模型类型?
│   │
│   ├─► 传统 ML (SKLearn/XGBoost)
│   │   │
│   │   └─► 方案 ──► KServe + SKLearn Runtime
│   │
│   ├─► 深度学习 (PyTorch/TensorFlow)
│   │   │
│   │   ├─► 中小模型 ──► KServe + Triton Inference Server
│   │   │
│   │   └─► 大模型 ──► KServe v0.16 + vLLM
│   │
│   └─► 大语言模型 (LLM)
│       │
│       ├─► 单卡推理 ──► KServe + vLLM (单节点)
│       │
│       ├─► 多卡并行 ──► KServe + vLLM (Tensor Parallelism)
│       │
│       └─► 超大模型 (405B+) ──► KServe v0.16 + Multi-Node Inference
│
├─► 流量模式?
│   │
│   ├─► 稳定流量 ──► HPA (基于 CPU/GPU 利用率)
│   │
│   ├─► 突发流量 ──► KEDA (基于队列长度/自定义指标)
│   │
│   └─► 成本敏感 ──► Scale-to-Zero (Knative/KServe)
│
├─► 延迟要求?
│   │
│   ├─► 低延迟 (<100ms) ──► GPU + Continuous Batching (vLLM)
│   │
│   ├─► 可接受延迟 ──► CPU + 量化模型
│   │
│   └─► 流式响应 ──► Streaming Tokens (vLLM + KServe)
│
└─► 高级特性?
    │
    ├─► A/B 测试 ──► KServe 流量分割
    │
    ├─► 模型组合 ──► KServe InferenceGraph
    │
    └─► 跨副本缓存 ──► LMCache (Distributed KV Cache)

推荐组合 (2025):
├── Serving: KServe v0.16
├── Runtime: vLLM (LLM) / Triton (DL) / SKLearn (ML)
├── Scaling: KEDA (基于 GPU 利用率或队列长度)
├── Gateway: Envoy AI Gateway (Token 限流、模型路由)
└── Observability: OpenTelemetry GenAI Semantic Conventions
```

---

## 3. 多集群管理场景树

```
场景: 管理多个 Kubernetes 集群
│
├─► 集群分布?
│   │
│   ├─► 同一云厂商多区域 ──► 方案
│   │   │
│   │   ├─► 简单管理 ──► Rancher / OpenLens
│   │   │
│   │   └─► GitOps 多集群 ──► ArgoCD ApplicationSet
│   │
│   ├─► 多云混合 ──► 方案
│   │   │
│   │   ├─► 应用分发 ──► Karmada / ArgoCD ApplicationSet
│   │   │
│   │   └─► 集群生命周期 ──► Cluster API
│   │
│   └─► 边缘 + 中心 ──► 方案
│       │
│       ├─► 云边协同 ──► KubeEdge / OpenYurt
│       │
│       └─► 纯边缘 ──► K3s + Rancher Fleet
│
├─► 网络连通?
│   │
│   ├─► 需要跨集群服务发现 ──► Submariner / Cilium Cluster Mesh
│   │
│   └─► 仅需应用访问 ──► 边缘 Ingress / Global Load Balancer
│
├─► 统一策略?
│   │
│   ├─► 安全策略 ──► OPA/Gatekeeper (全局策略)
│   │
│   └─► 网络策略 ──► Cilium Cluster Mesh (统一身份)
│
└─► 灾难恢复?
    │
    ├─► 应用级备份 ──► Velero
    │
    └─► 数据级备份 ──► 云厂商工具 / CSI Snapshot

推荐组合 (2025):
├── Cluster Lifecycle: Cluster API
├── App Distribution: ArgoCD ApplicationSet
├── Networking: Cilium Cluster Mesh (统一网络身份)
├── Policy: OPA (全局策略)
└── Backup: Velero + CSI Snapshots
```

---

## 4. 零信任安全场景树

```
场景: 实施 Kubernetes 零信任安全
│
├─► 镜像安全?
│   │
│   ├─► 镜像签名 ──► Cosign (Keyless Signing)
│   │
│   ├─► 镜像扫描 ──► Trivy (CI/CD) + Admission Controller
│   │
│   └─► SBOM ──► Syft 生成 + 存储
│
├─► 运行时安全?
│   │
│   ├─► 异常检测 ──► Falco (规则引擎)
│   │
│   ├─► 进程/网络监控 ──► Tetragon (eBPF)
│   │
│   └─► 沙箱隔离 ──► gVisor / Kata (敏感工作负载)
│
├─► 网络安全?
│   │
│   ├─► L3/L4 策略 ──► Cilium Network Policy
│   │
│   ├─► L7 策略 ──► Istio Authorization Policy
│   │
│   └─► mTLS ──► Istio / Cilium (透明加密)
│
├─► 准入控制?
│   │
│   ├─► 简单策略 ──► Kyverno
│   │
│   ├─► 复杂策略 ──► OPA Gatekeeper
│   │
│   └─► K8s 原生 ──► ValidatingAdmissionPolicy (v1.30+)
│
└─► 密钥管理?
    │
    ├─► 云厂商集成 ──► External Secrets + AWS/GCP/Azure Secret Manager
    │
    └─► 自托管 ──► External Secrets + Vault

推荐组合 (2025):
├── Image: Cosign + Trivy + Kyverno (准入验证)
├── Runtime: Falco + Tetragon (双层防护)
├── Network: Cilium (eBPF) + mTLS
├── Policy: Kyverno (K8s 策略) + OPA (全局策略)
└── Secrets: External Secrets Operator
```

---

## 5. 边缘计算场景树

```
场景: 在边缘部署 Kubernetes
│
├─► 边缘节点数量?
│   │
│   ├─► 少量 (<10) ──► K3s (单节点或嵌入式 HA)
│   │
│   ├─► 中等 (10-100) ──► K3s + System Upgrade Controller
│   │
│   └─► 大规模 (>100) ──► KubeEdge / OpenYurt (云边协同)
│
├─► 网络条件?
│   │
│   ├─► 稳定连接 ──► 标准 K3s
│   │
│   ├─► 间歇连接 ──► KubeEdge (离线自治) / K3s (边缘自治配置)
│   │
│   └─► 纯离线 ──► K3s (完全离线安装)
│
├─► 设备管理?
│   │
│   ├─► 简单设备 ──► K3s + Device Plugin
│   │
│   └─► 复杂设备孪生 ──► KubeEdge DeviceTwin
│
├─► 应用特性?
│   │
│   ├─► 数据处理 ──► K3s + Kafka / MQTT
│   │
│   ├─► AI 推理 ──► K3s + KServe (轻量模式)
│   │
│   └─► 视频分析 ──► K3s + GPU + Edge AI
│
└─► 运维管理?
    │
    ├─► 集中监控 ──► Prometheus Agent (Remote Write)
    │
    └─► 批量升级 ──► Rancher Fleet / System Upgrade Controller

推荐组合 (2025):
├── Distro: K3s (<100 节点) / KubeEdge (>100 节点)
├── Networking: Flannel (简单) / Cilium (高级)
├── Storage: Longhorn (边缘存储) / HostPath (临时)
├── Observability: OpenTelemetry Collector (边缘聚合)
└── Management: Rancher Fleet (多集群) / System Upgrade (单集群)
```

---

## 6. 成本优化场景树

```
场景: 优化 Kubernetes 云成本
│
├─► 计算成本优化?
│   │
│   ├─► 工作负载类型?
│   │   │
│   │   ├─► 可中断 ──► Spot/Preemptible Instances + Tolerations
│   │   │
│   │   ├─► 批处理 ──► Kueue (队列调度) + Spot
│   │   │
│   │   └─► 在线服务 ──► Spot (多可用区容错) + On-Demand (保底)
│   │
│   └─► 自动扩缩?
│       │
│       ├─► Pod 级 ──► HPA (CPU/内存) / KEDA (自定义指标)
│       │
│       └─► 节点级 ──► Cluster Autoscaler / Karpenter
│
├─► 存储成本优化?
│   │
│   ├─► 存储分层 ──► SSD (热) → HDD (温) → 对象存储 (冷)
│   │
│   ├─► 数据生命周期 ──► S3 Lifecycle Policy / MinIO ILM
│   │
│   └─► 精简配置 ──► Thin Provisioning (Ceph/Longhorn)
│
├─► 网络成本优化?
│   │
│   ├─► 流量控制 ──► Cilium (eBPF 高效转发)
│   │
│   ├─► 跨区域流量 ──► Topology Aware Routing
│   │
│   └─► 出口流量 ──► NAT Gateway 共享 / VPC Endpoints
│
└─► 成本可视化?
    │
    ├─► 开源 ──► OpenCost (免费) / Kubecost (社区版)
    │
    └─► 企业 ──► Kubecost Enterprise / CloudHealth

推荐组合 (2025):
├── Compute: Karpenter (智能节点配置) + Spot Instances
├── Autoscaling: HPA + Cluster Autoscaler
├── Storage: Ceph (分层存储) + Lifecycle Policies
├── Networking: Cilium (eBPF 减少开销)
└── Visibility: OpenCost (成本分摊) + Kubecost (优化建议)
```

---

## 7. 平台工程场景树

```
场景: 构建内部开发者平台 (IDP)
│
├─► 开发者门户?
│   │
│   ├─► 开源 ──► Backstage (Spotify) - 高度可定制
│   │
│   └─► 商业 ──► Port / Cortex / OpsLevel
│
├─► 基础设施供应?
│   │
│   ├─► 云资源 ──► Crossplane (K8s API 管理云资源)
│   │
│   ├─► 集群创建 ──► Cluster API
│   │
│   └─► 应用配置 ──► Helm / Kustomize + GitOps
│
├─► 黄金路径?
│   │
│   ├─► 服务模板 ──► Backstage Scaffolder + Cookiecutter
│   │
│   ├─► 部署模板 ──► ArgoCD ApplicationSet + Git Generators
│   │
│   └─► 文档 ──► Backstage TechDocs (Docs-as-Code)
│
├─► 开发者自助?
│   │
│   ├─► 环境申请 ──► Backstage Template + Crossplane
│   │
│   ├─► 数据库申请 ──► Backstage Template + CNPG (PostgreSQL)
│   │
│   └─► 密钥管理 ──► Backstage + External Secrets
│
└─► 平台可观测?
    │
    ├─► 平台健康 ──► Prometheus + Grafana
    │
    ├─► 开发者体验 ──► DORA 指标 (部署频率、恢复时间)
    │
    └─► 成本归因 ──► Kubecost (按团队分摊)

推荐组合 (2025):
├── Portal: Backstage (开源) / Port (商业)
├── Infra: Crossplane (云资源) + Cluster API (集群)
├── GitOps: ArgoCD (应用交付)
├── Policy: OPA (策略即代码)
├── Observability: Prometheus + Grafana (平台健康)
└── Cost: Kubecost (成本归因)
```

---

## 📊 场景-技术映射速查表

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      场景-技术映射速查表                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  业务场景              核心技术栈                                             │
│  ────────────────────────────────────────────────────────────────────────   │
│                                                                              │
│  电商微服务            Cilium + Istio Ambient + ArgoCD + OpenTelemetry      │
│  金融核心系统          Kata + Cilium + Tetragon + Vault                     │
│  AI 推理平台           KServe v0.16 + vLLM + KEDA + Envoy AI Gateway        │
│  IoT 边缘              KubeEdge + MQTT + InfluxDB                           │
│  多租户 SaaS           gVisor + Cilium + OPA + Backstage                    │
│  数据平台              Spark Operator + Argo Workflows + MinIO              │
│  DevOps 平台           Backstage + Crossplane + ArgoCD + DORA 指标          │
│  游戏后端              Agones + Karpenter + Cilium DSR                      │
│  视频处理              GPU Operator + Kueue + Rook-Ceph                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
