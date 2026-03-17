# 云原生技术决策树

> 基于场景的决策路径，帮助快速选择合适的技术方案

---

## 1. 容器运行时选择决策树

```
开始: 需要选择容器运行时
│
├─► 用于 Kubernetes 生产环境?
│   │
│   ├─► 是 ───────────────────────► containerd 2.0 (推荐)
│   │
│   └─► 否 ──► 用于开发测试?
│       │
│       ├─► 是 ──► 需要构建镜像?
│       │   │
│       │   ├─► 是 ─────────────► Docker Desktop
│       │   │
│       │   └─► 否 ──► 需要无根?
│       │       │
│       │       ├─► 是 ─────────► Podman
│       │       │
│       │       └─► 否 ─────────► Docker
│       │
│       └─► 否 ──► 特殊需求?
│           │
│           ├─► 高安全多租户 ───► gVisor / Kata
│           │
│           └─► Serverless ─────► Firecracker
│
关键决策点:
• K8s 生产 → containerd (标准、轻量、安全)
• 开发测试 → Docker/Podman (易用)
• 高安全 → 沙箱运行时 (gVisor/Kata)
```

---

## 2. CNI 网络方案选择决策树

```
开始: 需要选择 Kubernetes 网络方案
│
├─► 是否有高级网络需求?
│   │
│   ├─► 是 ──► 需要 L7 网络策略?
│   │   │
│   │   ├─► 是 ──► 需要内置可观测性?
│   │   │   │
│   │   │   ├─► 是 ─────────────► Cilium + Hubble
│   │   │   │
│   │   │   └─► 否 ─────────────► Istio CNI
│   │   │
│   │   └─► 否 ──► 需要替代 kube-proxy?
│   │       │
│   │       ├─► 是 ─────────────► Cilium
│   │       │
│   │       └─► 否 ─────────────► Calico (eBPF 模式)
│   │
│   └─► 否 ──► 集群规模?
│       │
│       ├─► 大规模 (>1000 节点) ──► Calico / Cilium
│       │
│       ├─► 中等规模 ─────────────► Calico (标准)
│       │
│       └─► 小规模/测试 ──────────► Flannel
│
关键决策点:
• L7 策略 + 可观测性 → Cilium (2025 首选)
• 传统稳定 → Calico
• 简单最小化 → Flannel
```

---

## 3. 服务网格选择决策树

```
开始: 需要服务网格能力
│
├─► 是否已有 Cilium 网络?
│   │
│   ├─► 是 ──► 需求简单?
│   │   │
│   │   ├─► 是 ────────────────► Cilium Service Mesh
│   │   │
│   │   └─► 否 ────────────────► Istio + Cilium
│   │
│   └─► 否 ──► 资源敏感?
│       │
│       ├─► 是 ──► 功能需求?
│       │   │
│       │   ├─► 基础功能 ──────► Linkerd
│       │   │
│       │   └─► 完整功能 ──────► Istio Ambient
│       │
│       └─► 否 ──► 功能完整度?
│           │
│           ├─► 最全功能 ──────► Istio Ambient (2025 推荐)
│           │
│           ├─► 平衡选择 ──────► Linkerd
│           │
│           └─► VM 混合 ───────► Consul
│
2025 推荐路径:
新部署 → Istio Ambient Mesh (无 Sidecar、低资源)
```

---

## 4. 存储方案选择决策树

```
开始: 需要选择存储方案
│
├─► 存储类型?
│   │
│   ├─► 块存储 ──► 云厂商?
│   │   │
│   │   ├─► AWS ──────────────► EBS + EBS CSI
│   │   ├─► GCP ──────────────► Persistent Disk
│   │   ├─► Azure ────────────► Managed Disks
│   │   └─► 裸金属/私有云 ────► Longhorn / OpenEBS / Ceph
│   │
│   ├─► 文件存储 ──► 共享需求?
│   │   │
│   │   ├─► 跨多节点读 ───────► NFS / EFS
│   │   ├─► 高性能 ───────────► Lustre / GPFS
│   │   └─► 简单共享 ─────────► NFS CSI / Longhorn
│   │
│   ├─► 对象存储 ──► 使用场景?
│   │   │
│   │   ├─► 备份/日志 ────────► S3 / MinIO
│   │   ├─► ML 数据集 ────────► S3 + S3 CSI
│   │   └─► 内部部署 ─────────► MinIO / Ceph RGW
│   │
│   └─► 数据库存储 ──► 数据库类型?
│       │
│       ├─► PostgreSQL ───────► Rook-Ceph / 本地 SSD
│       ├─► MySQL ────────────► Rook-Ceph / 本地 SSD
│       ├─► NoSQL ────────────► 本地存储 / Ceph
│       └─► 分布式 DB ────────► 原生存储 (TiDB/Cassandra)
│
生产推荐:
• 通用: Rook-Ceph (块+文件+对象)
• 简单: Longhorn
• 云原生: 云厂商 CSI
```

---

## 5. GitOps 工具选择决策树

```
开始: 需要 GitOps 方案
│
├─► 团队规模?
│   │
│   ├─► 小型团队 (<10人) ──► 复杂度?
│   │   │
│   │   ├─► 简单 ───────────► Flux
│   │   │
│   │   └─► 需要 UI ────────► ArgoCD
│   │
│   ├─► 中型团队 ──► 特殊需求?
│   │   │
│   │   ├─► 渐进交付 ───────► ArgoCD + Argo Rollouts
│   │   │
│   │   ├─► 多集群管理 ─────► ArgoCD / Rancher Fleet
│   │   │
│   │   └─► 纯 GitOps ──────► Flux
│   │
│   └─► 大型企业 ──► 核心需求?
│       │
│       ├─► 完整功能 + UI ──► ArgoCD Enterprise
│       │
│       ├─► 多集群 + 边缘 ──► Rancher Fleet
│       │
│       └─► CI/CD 一体化 ───► GitLab + Flux
│
2025 推荐:
• 通用: ArgoCD (功能最全、生态最好)
• Git原生: Flux
• 多集群: ArgoCD 或 Fleet
```

---

## 6. 可观测性方案选择决策树

```
开始: 需要可观测性方案
│
├─► 数据采集 ──► 现有埋点?
│   │
│   ├─► 已有 ──► 格式?
│   │   │
│   │   ├─► Prometheus ──────► OTel Collector (转换)
│   │   ├─► Jaeger ──────────► OTel Collector (转换)
│   │   └─► 云厂商 ──────────► OTel Collector (导出)
│   │
│   └─► 新建 ──► 语言?
│       │
│       ├─► Java/Go/Python ──► OpenTelemetry SDK
│       ├─► .NET/Node.js ────► OpenTelemetry SDK
│       └─► 其他 ────────────► OTel Auto-Instrumentation
│
├─► 存储后端 ──► 预算?
│   │
│   ├─► 有限 ──► 自托管
│   │   │
│   │   ├─► Metrics ────────► Prometheus 3.0
│   │   ├─► Logs ───────────► Loki / VictoriaLogs
│   │   ├─► Traces ─────────► Jaeger / Tempo
│   │   └─► All-in-One ─────► SigNoz / Grafana Stack
│   │
│   └─► 充足 ──► SaaS
│       │
│       ├─► 全功能 ─────────► Datadog / Dynatrace
│       └─► 成本敏感 ───────► Honeycomb / Grafana Cloud
│
├─► 可视化 ──► 已有方案?
│   │
│   ├─► Prometheus ─────────► Grafana
│   ├─► Datadog ────────────► Datadog UI
│   └─► OpenTelemetry ──────► Grafana / Jaeger UI
│
2025 推荐方案:
OpenTelemetry (采集) → Prometheus 3.0 (指标) + Loki (日志) + Tempo (追踪) → Grafana (可视化)
```

---

## 7. AI/ML 平台选择决策树

```
开始: 需要 AI/ML 平台
│
├─► 主要场景?
│   │
│   ├─► 模型训练 ──► 规模?
│   │   │
│   │   ├─► 单机 ─────────────► Jupyter + MLflow
│   │   ├─► 分布式 ───────────► Kubeflow + Ray
│   │   └─► 大规模分布式 ─────► Kubeflow + Volcano + Ray
│   │
│   ├─► 模型推理 ──► 模型类型?
│   │   │
│   │   ├─► 传统 ML ──────────► KServe + SKLearn
│   │   ├─► 深度学习 ─────────► KServe + Triton
│   │   ├─► LLM ──────────────► KServe v0.16 + vLLM
│   │   └─► 多模型管理 ───────► KServe + Model Registry
│   │
│   ├─► MLOps 全流程 ──► 完整度?
│   │   │
│   │   ├─► 完整平台 ─────────► Kubeflow
│   │   └─► 轻量组合 ─────────► MLflow + KServe
│   │
│   └─► 实验管理 ─────────────► MLflow / Weights & Biases
│
2025 推荐组合:
Kubeflow (训练) + KServe v0.16 (推理) + MLflow (模型管理) + Ray (分布式)
```

---

## 8. 安全方案选择决策树

```
开始: 需要加强 K8s 安全
│
├─► 哪个层面?
│   │
│   ├─► 镜像安全 ──► 阶段?
│   │   │
│   │   ├─► 构建时 ───────────► Cosign 签名 + SLSA
│   │   ├─► 扫描 ─────────────► Trivy / Grype + SBOM
│   │   └─► 准入控制 ─────────► Kyverno/OPA 镜像验证
│   │
│   ├─► 运行时安全 ──► 威胁类型?
│   │   │
│   │   ├─► 异常行为 ─────────► Falco
│   │   ├─► 网络攻击 ─────────► Tetragon + Cilium
│   │   ├─► 逃逸防护 ─────────► gVisor / Kata
│   │   └─► 零信任 ───────────► Istio mTLS + SPIFFE
│   │
│   ├─► 策略治理 ──► 复杂度?
│   │   │
│   │   ├─► 简单 ─────────────► Kyverno
│   │   ├─► 复杂 ─────────────► OPA Gatekeeper
│   │   └─► K8s 原生 ─────────► ValidatingAdmissionPolicy
│   │
│   └─► 密钥管理 ─────────────► External Secrets + Vault
│
2025 安全栈推荐:
Cosign (镜像签名) + Trivy (扫描) + Kyverno (策略) + Falco (运行时) + Tetragon (eBPF)
```

---

## 9. 边缘计算方案选择决策树

```
开始: 需要边缘 K8s 方案
│
├─► 边缘节点数量?
│   │
│   ├─► 少量 (<10) ──► 连接性?
│   │   │
│   │   ├─► 稳定 ─────────────► K3s 单节点
│   │   └─► 不稳定 ───────────► K3s + 自治配置
│   │
│   ├─► 中等 (10-100) ──► 管理需求?
│   │   │
│   │   ├─► 云边协同 ─────────► KubeEdge
│   │   ├─► 纯边缘 ───────────► K3s + System Upgrade
│   │   └─► 国内优化 ─────────► SuperEdge
│   │
│   └─► 大规模 (>100) ──► 云平台?
│       │
│       ├─► 阿里云 ───────────► OpenYurt
│       ├─► AWS ─────────────► EKS Anywhere
│       └─► 多云 ────────────► KubeEdge / OpenYurt
│
关键决策点:
• 简单边缘 → K3s
• 云边协同 → KubeEdge
• 大规模管理 → OpenYurt / KubeEdge
```

---

## 📋 快速决策参考表

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          快速决策参考                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  场景                          推荐方案                                      │
│  ────────────────────────────────────────────────────────────────────────   │
│                                                                              │
│  新建 K8s 集群                 Cilium + containerd 2.0                      │
│  服务网格升级                  Istio Ambient Mesh                           │
│  GitOps 初建                   ArgoCD                                       │
│  可观测性现代化                OpenTelemetry + Grafana Stack                │
│  AI 推理服务                   KServe v0.16 + vLLM                          │
│  边缘部署                      K3s / KubeEdge                               │
│  多集群管理                    Cluster API + ArgoCD                         │
│  平台工程                      Backstage + Crossplane                       │
│  供应链安全                    Cosign + SLSA + Kyverno                      │
│  成本优化                      Kubecost + Spot instances                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
