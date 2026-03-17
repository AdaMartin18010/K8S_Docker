# Kueue - Kubernetes 批处理调度系统

## 概述

Kueue 是一个 CNCF Incubating 项目，为 Kubernetes 提供作业队列和批处理调度能力。它增强默认调度器，引入队列系统，支持资源配额、作业优先级和配额借用。

## 核心特性

| 特性 | 描述 |
|------|------|
| 作业队列 | 工作负载排队直到资源就绪 |
| 资源配额 | 按资源组和类型配置配额 |
| 配额借用 | Cohort 队列间借用空闲配额 |
| 抢占策略 | 基于优先级的工作负载抢占 |
| GPU 感知 | 支持 GPU 和异构资源调度 |

## 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                      Kueue 架构                                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   ClusterQueue                          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │   │
│  │  │  Flavor A   │  │  Flavor B   │  │  Flavor C   │     │   │
│  │  │ (A100 GPU)  │  │  (H100 GPU) │  │  (CPU)      │     │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘     │   │
│  │                        │                                │   │
│  │              ┌─────────┴─────────┐                      │   │
│  │              ▼                   ▼                      │   │
│  │  ┌─────────────────────┐ ┌─────────────────────┐       │   │
│  │  │     Cohort          │ │     Cohort          │       │   │
│  │  │  (队列组)           │ │  (队列组)           │       │   │
│  │  │                     │ │                     │       │   │
│  │  │  ┌───────────────┐  │ │  ┌───────────────┐  │       │   │
│  │  │  │ LocalQueue    │  │ │  │ LocalQueue    │  │       │   │
│  │  │  │ (team-a)      │  │ │  │ (team-b)      │  │       │   │
│  │  │  └───────────────┘  │ │  └───────────────┘  │       │   │
│  │  │  ┌───────────────┐  │ │  ┌───────────────┐  │       │   │
│  │  │  │ LocalQueue    │  │ │  │ LocalQueue    │  │       │   │
│  │  │  │ (team-c)      │  │ │  │ (team-d)      │  │       │   │
│  │  │  └───────────────┘  │ │  └───────────────┘  │       │   │
│  │  └─────────────────────┘ └─────────────────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                      Admission                          │   │
│  │  1. 工作负载进入 LocalQueue                              │   │
│  │  2. ClusterQueue 检查配额                                │   │
│  │  3. 模拟调度 ( flavor 选择)                               │   │
│  │  4. 准入或等待                                           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  Kubernetes Workloads                    │   │
│  │  - Job, CronJob                                          │   │
│  │  - MPIJob, PyTorchJob                                    │   │
│  │  - RayJob, RayCluster                                    │   │
│  │  - Kubeflow Training                                     │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 安装部署

### Helm 安装

```bash
# 添加 Helm 仓库
helm repo add kueue https://kubernetes-sigs.github.io/kueue/
helm repo update

# 安装 Kueue
helm install kueue kueue/kueue \
  --namespace kueue-system \
  --create-namespace \
  --version 0.10.0

# 验证安装
kubectl get pods -n kueue-system
```

### kubectl 安装

```bash
# 安装最新版本
kubectl apply --server-side -f https://github.com/kubernetes-sigs/kueue/releases/download/v0.10.0/manifests.yaml

# 等待就绪
kubectl wait --for=condition=available deployment/kueue-controller-manager -n kueue-system
```

## 核心概念

### ResourceFlavor

```yaml
apiVersion: kueue.x-k8s.io/v1beta1
kind: ResourceFlavor
metadata:
  name: nvidia-a100
spec:
  nodeLabels:
    nvidia.com/gpu.family: ampere
    node.kubernetes.io/instance-type: p4d.24xlarge
  tolerations:
  - key: nvidia.com/gpu
    operator: Exists
    effect: NoSchedule
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: ResourceFlavor
metadata:
  name: nvidia-h100
spec:
  nodeLabels:
    nvidia.com/gpu.family: hopper
    node.kubernetes.io/instance-type: p5.48xlarge
```

### ClusterQueue

```yaml
apiVersion: kueue.x-k8s.io/v1beta1
kind: ClusterQueue
metadata:
  name: gpu-cluster-queue
spec:
  namespaceSelector: {}  # 所有命名空间
  resourceGroups:
  - coveredResources: ["cpu", "memory", "nvidia.com/gpu"]
    flavors:
    - name: nvidia-h100
      resources:
      - name: "cpu"
        nominalQuota: 1920  # 40 节点 × 48 CPU
      - name: "memory"
        nominalQuota: 7.5Ti
      - name: "nvidia.com/gpu"
        nominalQuota: 320   # 40 节点 × 8 GPU
    - name: nvidia-a100
      resources:
      - name: "cpu"
        nominalQuota: 960
      - name: "memory"
        nominalQuota: 3.8Ti
      - name: "nvidia.com/gpu"
        nominalQuota: 64
  preemption:
    withinClusterQueue: LowerPriority
  # 启用借用
  cohort: ai-training-cohort
```

### LocalQueue

```yaml
apiVersion: kueue.x-k8s.io/v1beta1
kind: LocalQueue
metadata:
  name: research-queue
  namespace: team-research
spec:
  clusterQueue: gpu-cluster-queue
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: LocalQueue
metadata:
  name: production-queue
  namespace: team-production
spec:
  clusterQueue: gpu-cluster-queue
```

### WorkloadPriorityClass

```yaml
apiVersion: kueue.x-k8s.io/v1beta1
kind: WorkloadPriorityClass
metadata:
  name: critical
value: 1000
preemptionPolicy: Never
description: "Critical workloads that should not be preempted"
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: WorkloadPriorityClass
metadata:
  name: high
value: 500
preemptionPolicy: PreemptLowerPriority
description: "High priority training jobs"
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: WorkloadPriorityClass
metadata:
  name: spot
value: 100
preemptionPolicy: PreemptLowerPriority
description: "Spot/preemptible workloads"
```

## 使用示例

### 基本 Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: training-job
  namespace: team-research
  labels:
    kueue.x-k8s.io/queue-name: research-queue
spec:
  parallelism: 4
  completions: 4
  template:
    spec:
      priorityClassName: high
      containers:
      - name: trainer
        image: pytorch/pytorch:latest
        resources:
          requests:
            cpu: 8
            memory: 32Gi
            nvidia.com/gpu: 2
          limits:
            nvidia.com/gpu: 2
        command:
        - python
        - train.py
        - --epochs=100
      restartPolicy: Never
```

### PyTorchJob (Kubeflow)

```yaml
apiVersion: kubeflow.org/v1
kind: PyTorchJob
metadata:
  name: distributed-training
  namespace: team-research
  labels:
    kueue.x-k8s.io/queue-name: research-queue
spec:
  pytorchReplicaSpecs:
    Master:
      replicas: 1
      template:
        spec:
          priorityClassName: critical
          containers:
          - name: pytorch
            image: pytorch/pytorch:latest
            resources:
              requests:
                cpu: 16
                memory: 128Gi
                nvidia.com/gpu: 8
            command:
            - python
            - -m
            - torch.distributed.run
            - --nproc_per_node=8
            - train.py
    Worker:
      replicas: 3
      template:
        spec:
          priorityClassName: critical
          containers:
          - name: pytorch
            image: pytorch/pytorch:latest
            resources:
              requests:
                cpu: 16
                memory: 128Gi
                nvidia.com/gpu: 8
```

### 抢占配置

```yaml
apiVersion: kueue.x-k8s.io/v1beta1
kind: ClusterQueue
metadata:
  name: production-queue
spec:
  preemption:
    withinClusterQueue: LowerPriority
    reclaimWithinCohort: Any
    borrowWithinCohort:
      policy: LowerPriority
      maxPriorityThreshold: 500
  resourceGroups:
  - coveredResources: ["nvidia.com/gpu"]
    flavors:
    - name: nvidia-h100
      resources:
      - name: "nvidia.com/gpu"
        nominalQuota: 100
        borrowingLimit: 50  # 可借用 50%
```

## MIG 配置

```yaml
apiVersion: kueue.x-k8s.io/v1beta1
kind: ResourceFlavor
metadata:
  name: nvidia-h100-mig
spec:
  nodeLabels:
    nvidia.com/gpu.product: NVIDIA-H100-80GB-HBM3
    nvidia.com/mig.strategy: single
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: ClusterQueue
metadata:
  name: inference-queue
spec:
  resourceGroups:
  - coveredResources: ["nvidia.com/mig-1g.10gb"]
    flavors:
    - name: nvidia-h100-mig
      resources:
      - name: "nvidia.com/mig-1g.10gb"
        nominalQuota: 280  # 40 节点 × 7 MIG slice
```

## 相关资源

- [Kueue 官网](https://kueue.sigs.k8s.io/)
- [GitHub](https://github.com/kubernetes-sigs/kueue)
- [CNCF 项目](https://www.cncf.io/projects/kueue/)
