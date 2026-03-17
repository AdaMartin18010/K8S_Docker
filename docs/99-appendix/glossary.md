# 云原生术语表

> 快速查阅云原生领域核心概念和术语

---

## A

| 术语 | 英文 | 解释 |
|------|------|------|
| **Admission Controller** | 准入控制器 | K8s 中拦截 API 请求的组件，用于验证或修改资源对象 |
| **Affinity** | 亲和性 | Pod 调度偏好，让 Pod 倾向于部署在特定节点或与某些 Pod 共存 |
| **Annotation** | 注解 | K8s 资源对象的元数据，用于存储非标识信息 |
| **Anti-affinity** | 反亲和性 | Pod 调度偏好，让 Pod 避免部署在特定节点或与某些 Pod 共存 |
| **ArgoCD** | ArgoCD | 声明式 GitOps 持续交付工具 |

## B

| 术语 | 英文 | 解释 |
|------|------|------|
| **BuildKit** | BuildKit | Docker 的下一代镜像构建工具，支持并行构建和缓存优化 |
| **Backstage** | Backstage | Spotify 开源的开发者门户框架，用于构建内部开发者平台 |

## C

| 术语 | 英文 | 解释 |
|------|------|------|
| **CNI** | Container Network Interface | 容器网络接口，定义容器网络连接的标准 |
| **ConfigMap** | ConfigMap | K8s 中用于存储非敏感配置数据的对象 |
| **Container** | 容器 | 应用及其依赖的轻量级打包和运行单元 |
| **containerd** | containerd | 行业标准的容器运行时，Docker 和 K8s 的底层运行时 |
| **CSI** | Container Storage Interface | 容器存储接口，定义容器存储的标准 |
| **CRD** | Custom Resource Definition | 自定义资源定义，扩展 K8s API 的方式 |
| **Crossplane** | Crossplane | 云原生控制平面，使用 K8s API 管理云资源 |
| **Cilium** | Cilium | 基于 eBPF 的 K8s 网络和安全解决方案 |

## D

| 术语 | 英文 | 解释 |
|------|------|------|
| **DaemonSet** | DaemonSet | K8s 工作负载类型，确保每个节点运行一个 Pod 副本 |
| **Deployment** | Deployment | K8s 工作负载类型，管理无状态应用的声明式更新 |
| **Docker** | Docker | 最流行的容器平台，包括运行时和工具集 |
| **Dapr** | Dapr | 分布式应用运行时，简化微服务开发 |
| **DevOps** | DevOps | 开发和运维的协作文化与实践 |
| **DR** | Disaster Recovery | 灾难恢复，系统故障后的恢复能力 |
| **DRA** | Dynamic Resource Allocation | 动态资源分配，K8s 1.26+ 引入的 GPU 调度新机制 |

## E

| 术语 | 英文 | 解释 |
|------|------|------|
| **eBPF** | extended Berkeley Packet Filter | Linux 内核技术，用于安全监控和可观测性 |
| **Envoy** | Envoy | 高性能 C++ 分布式代理，Istio 的数据平面 |
| **ETCD** | etcd | 分布式键值存储，K8s 的集群状态存储 |
| **External Secret** | External Secret | 将外部密钥管理系统集成到 K8s 的工具 |

## F

| 术语 | 英文 | 解释 |
|------|------|------|
| **FinOps** | FinOps | 云成本管理实践，优化云资源支出 |
| **Falco** | Falco | 云原生运行时安全监控工具 |
| **Flux** | Flux | CNCF 的 GitOps 工具 |

## G

| 术语 | 英文 | 解释 |
|------|------|------|
| **GitOps** | GitOps | 以 Git 为唯一真相源的运维模式 |
| **Gateway API** | Gateway API | K8s 的下一代流量管理标准 |
| **gRPC** | gRPC | 高性能开源 RPC 框架 |
| **gVisor** | gVisor | Google 的沙箱容器运行时 |

## H

| 术语 | 英文 | 解释 |
|------|------|------|
| **HPA** | Horizontal Pod Autoscaler | 水平 Pod 自动扩缩容 |
| **Helm** | Helm | K8s 的包管理工具 |
| **Hubble** | Hubble | Cilium 的网络可观测性组件 |

## I

| 术语 | 英文 | 解释 |
|------|------|------|
| **Istio** | Istio | 最流行的服务网格实现 |
| **IDP** | Internal Developer Platform | 内部开发者平台 |
| **Ingress** | Ingress | K8s 传统的 HTTP 路由资源 |
| **Init Container** | Init Container | 在应用容器启动前运行的初始化容器 |

## J

| 术语 | 英文 | 解释 |
|------|------|------|
| **Jaeger** | Jaeger | 分布式追踪系统 |
| **Job** | Job | K8s 中运行一次性任务的工作负载 |

## K

| 术语 | 英文 | 解释 |
|------|------|------|
| **K8s** | Kubernetes | 容器编排平台的事实标准 |
| **K3s** | K3s | 轻量级 K8s 发行版 |
| **KEDA** | Kubernetes Event-driven Autoscaling | 基于事件的 K8s 自动扩缩容 |
| **Kind** | Kind | 使用 Docker 容器运行 K8s 集群的工具 |
| **Knative** | Knative | K8s 上的 Serverless 平台 |
| **KServe** | KServe | K8s 上的模型推理服务平台 |
| **Kubelet** | Kubelet | 每个 K8s 节点上运行的代理 |
| **Kubectl** | kubectl | K8s 命令行工具 |
| **Kustomize** | Kustomize | K8s 原生的配置定制工具 |
| **Kyverno** | Kyverno | K8s 原生策略引擎 |

## L

| 术语 | 英文 | 解释 |
|------|------|------|
| **Label** | Label | K8s 资源的键值对标识，用于选择和组织 |
| **LimitRange** | LimitRange | 命名空间级别的资源限制默认值 |
| **Linkerd** | Linkerd | 轻量级服务网格 |
| **Loki** | Loki | 水平可扩展的日志聚合系统 |

## M

| 术语 | 英文 | 解释 |
|------|------|------|
| **mTLS** | Mutual TLS | 双向 TLS 认证 |
| **Multi-tenancy** | 多租户 | 多个用户/团队共享集群资源 |

## N

| 术语 | 英文 | 解释 |
|------|------|------|
| **Namespace** | Namespace | K8s 中的虚拟集群，用于资源隔离 |
| **Node** | Node | K8s 集群中的工作机器 |
| **NetworkPolicy** | NetworkPolicy | K8s 网络策略，控制 Pod 间流量 |

## O

| 术语 | 英文 | 解释 |
|------|------|------|
| **OCI** | Open Container Initiative | 开放容器倡议，容器标准组织 |
| **OIDC** | OpenID Connect | 基于 OAuth 2.0 的身份认证协议 |
| **OpenTelemetry** | OpenTelemetry | 云原生可观测性标准 |
| **OPA** | Open Policy Agent | 通用策略引擎 |
| **Operator** | Operator | 封装运维知识的 K8s 控制器 |

## P

| 术语 | 英文 | 解释 |
|------|------|------|
| **PersistentVolume** | PersistentVolume | 集群级别的存储资源 |
| **PersistentVolumeClaim** | PersistentVolumeClaim | 用户对存储的请求 |
| **Pod** | Pod | K8s 中最小的可部署单元 |
| **PodDisruptionBudget** | PodDisruptionBudget | 确保升级期间的最小可用 Pod 数 |
| **Prometheus** | Prometheus | 云原生监控和告警系统 |

## Q

| 术语 | 英文 | 解释 |
|------|------|------|
| **QoS** | Quality of Service | K8s 中的服务质量等级 |
| **Quota** | ResourceQuota | 命名空间级别的资源配额限制 |

## R

| 术语 | 英文 | 解释 |
|------|------|------|
| **RBAC** | Role-Based Access Control | 基于角色的访问控制 |
| **ReplicaSet** | ReplicaSet | 确保指定数量的 Pod 副本运行 |
| **ResourceQuota** | ResourceQuota | 命名空间的资源配额 |
| **Runtime** | Runtime | 容器运行时，如 containerd、CRI-O |
| **RuntimeClass** | RuntimeClass | 选择不同容器运行时的机制 |

## S

| 术语 | 英文 | 解释 |
|------|------|------|
| **SaaS** | Software as a Service | 软件即服务 |
| **Secret** | Secret | K8s 中存储敏感数据的对象 |
| **SecurityContext** | SecurityContext | Pod/容器的安全设置 |
| **Service** | Service | K8s 中暴露应用的网络抽象 |
| **ServiceAccount** | ServiceAccount | Pod 访问 API 的身份 |
| **Service Mesh** | Service Mesh | 管理服务间通信的基础设施层 |
| **Sidecar** | Sidecar | 与应用容器共同运行的辅助容器 |
| **StatefulSet** | StatefulSet | 管理有状态应用的工作负载 |
| **StorageClass** | StorageClass | 存储类型定义 |
| **SLO** | Service Level Objective | 服务等级目标 |
| **SLI** | Service Level Indicator | 服务等级指标 |
| **SLA** | Service Level Agreement | 服务等级协议 |
| **SLSA** | Supply-chain Levels for Software Artifacts | 软件供应链安全等级框架 |
| **SBOM** | Software Bill of Materials | 软件物料清单 |
| **Sigstore** | Sigstore | 开源软件签名和验证生态 |

## T

| 术语 | 英文 | 解释 |
|------|------|------|
| **Taint** | Taint | 节点的污点，用于排斥 Pod |
| **Toleration** | Toleration | Pod 对污点的容忍 |
| **Trivy** | Trivy | 容器镜像安全扫描工具 |
| **Tempo** | Tempo | Grafana 的分布式追踪后端 |
| **Tetragon** | Tetragon | 基于 eBPF 的运行时安全监控 |

## V

| 术语 | 英文 | 解释 |
|------|------|------|
| **VPA** | Vertical Pod Autoscaler | 垂直 Pod 自动扩缩容 |
| **Volume** | Volume | Pod 可访问的存储 |
| **vCluster** | vCluster | 虚拟 K8s 集群 |
| **Velero** | Velero | K8s 备份恢复工具 |

## W

| 术语 | 英文 | 解释 |
|------|------|------|
| **Workload** | Workload | K8s 中的应用工作负载 |
| **Wasm** | WebAssembly | 用于 Web 的二进制指令格式 |
| **Webhook** | Webhook | HTTP 回调机制 |

---

## 常用缩写速查

| 缩写 | 全称 | 含义 |
|------|------|------|
| API | Application Programming Interface | 应用程序接口 |
| CRD | Custom Resource Definition | 自定义资源定义 |
| CSI | Container Storage Interface | 容器存储接口 |
| CNI | Container Network Interface | 容器网络接口 |
| CR | Custom Resource | 自定义资源 |
| CM | ConfigMap | 配置映射 |
| PVC | PersistentVolumeClaim | 持久卷声明 |
| PV | PersistentVolume | 持久卷 |
| SA | ServiceAccount | 服务账户 |
| NS | Namespace | 命名空间 |
| HPA | Horizontal Pod Autoscaler | 水平 Pod 自动扩缩容 |
| VPA | Vertical Pod Autoscaler | 垂直 Pod 自动扩缩容 |
| PDB | PodDisruptionBudget | Pod 中断预算 |
| QoS | Quality of Service | 服务质量 |
| RBAC | Role-Based Access Control | 基于角色的访问控制 |
| SRE | Site Reliability Engineering | 站点可靠性工程 |
| SLI | Service Level Indicator | 服务等级指标 |
| SLO | Service Level Objective | 服务等级目标 |
| SLA | Service Level Agreement | 服务等级协议 |
| MTTR | Mean Time To Recovery | 平均恢复时间 |
| MTBF | Mean Time Between Failures | 平均故障间隔时间 |

---

## 相关资源

- [CNCF Glossary](https://glossary.cncf.io/)
- [Kubernetes Documentation](https://kubernetes.io/docs/concepts/)
- [Docker Glossary](https://docs.docker.com/glossary/)
