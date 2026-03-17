# 多集群与联邦

> 跨 Kubernetes 集群的资源管理

---

## 多集群架构模式

```
┌─────────────────────────────────────────────────────────────┐
│                  多集群架构模式                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  模式 1: 集中式控制平面                                      │
│  ┌─────────────────┐                                         │
│  │  Global Control │                                         │
│  │     Plane       │                                         │
│  └────────┬────────┘                                         │
│           │ pushes configs                                    │
│     ┌─────┴─────┐                                            │
│     ↓         ↓                                              │
│  ┌──────┐  ┌──────┐                                          │
│  │K8s A │  │K8s B │                                          │
│  └──────┘  └──────┘                                          │
│                                                              │
│  模式 2: 联邦式 (Karmada/OCM)                                │
│  ┌─────────────────┐                                         │
│  │   Karmada       │                                         │
│  │   Control Plane │                                         │
│  └────────┬────────┘                                         │
│           │ schedules workloads                              │
│     ┌─────┴─────┐                                            │
│     ↓         ↓                                              │
│  ┌──────┐  ┌──────┐                                          │
│  │K8s A │  │K8s B │  (Member Clusters)                       │
│  └──────┘  └──────┘                                          │
└─────────────────────────────────────────────────────────────┘
```

---

## Karmada

Kubernetes Armada，多云原生应用管理平台。

```yaml
apiVersion: policy.karmada.io/v1alpha1
kind: PropagationPolicy
metadata:
  name: nginx-propagation
spec:
  resourceSelectors:
    - apiVersion: apps/v1
      kind: Deployment
      name: nginx
  placement:
    clusterAffinity:
      clusterNames:
        - cluster-beijing
        - cluster-shanghai
    replicaScheduling:
      replicaDivisionPreference: Weighted
      replicaSchedulingType: Divided
      weightPreference:
        staticWeightList:
          - targetCluster:
              clusterNames: [cluster-beijing]
            weight: 60
          - targetCluster:
              clusterNames: [cluster-shanghai]
            weight: 40
```

---

## Open Cluster Management (OCM)

CNCF 沙箱项目，多集群管理标准。

```yaml
apiVersion: cluster.open-cluster-management.io/v1
kind: ManagedCluster
metadata:
  name: cluster-west
  labels:
    region: us-west
    env: production
spec:
  hubAcceptsClient: true
---
apiVersion: work.open-cluster-management.io/v1
kind: ManifestWork
metadata:
  name: deploy-app
  namespace: cluster-west
spec:
  workload:
    manifests:
      - apiVersion: apps/v1
        kind: Deployment
        metadata:
          name: myapp
        spec:
          replicas: 3
```

---

## 服务跨集群通信

### Submariner

```
Cluster A                    Cluster B
┌──────────────┐             ┌──────────────┐
│  Service A   │◄───────────►│  Service B   │
│  10.0.1.10   │   Submariner│  10.1.2.20   │
└──────┬───────┘             └──────┬───────┘
       │                            │
       └────────────┬───────────────┘
                    │
              ┌─────▼──────┐
              │  Broker    │
              │  (共享)     │
              └────────────┘
```

### Cilium Cluster Mesh

```yaml
apiVersion: cilium.io/v2
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: allow-cross-cluster
spec:
  endpointSelector:
    matchLabels:
      app: frontend
  egress:
    - toEndpoints:
        - matchLabels:
            app: backend
            io.cilium.k8s.policy.cluster: cluster-b
```

---

## 工具对比

| 工具 | 类型 | 特点 |
|------|------|------|
| **Karmada** | 联邦调度 | 原生 K8s API 兼容 |
| **OCM** | 集群管理 | 解耦管理平面 |
| **Rancher** | 管理平台 | 完整 UI，企业级 |
| **Fleet** | GitOps | Rancher 生态 |
| **ArgoCD** | GitOps | 应用级多集群 |
| **Submariner** | 网络 | 跨集群网络 |
| **Cilium CM** | 网络 + 安全 | 服务网格级 |

---

## 使用场景

| 场景 | 推荐方案 |
|------|----------|
| **灾难恢复** | 主备集群 + Global Load Balancer |
| **就近访问** | 地域分散 + Karmada 调度 |
| **云厂商锁定避免** | 多云联邦 |
| **边缘计算** | 中心-边缘架构 (KubeEdge) |
