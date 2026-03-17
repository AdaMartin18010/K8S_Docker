# 节点自动扩缩容

> K8s 集群弹性伸缩完整方案

---

## 扩缩容方案对比

| 工具 | 类型 | 云厂商 | 特点 |
|------|------|--------|------|
| **HPA** | Pod 水平扩缩 | 通用 | 内置，基于指标 |
| **VPA** | Pod 垂直扩缩 | 通用 | 调整资源请求/限制 |
| **KEDA** | 事件驱动扩缩 | 通用 | 基于事件源 |
| **CA** | 节点水平扩缩 | 通用 | Cluster Autoscaler |
| **Karpenter** | 节点自动配置 | AWS/其他 | 智能、快速 |

---

## Cluster Autoscaler

```
┌─────────────────────────────────────────────────────────────┐
│              Cluster Autoscaler 工作原理                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐                   │
│  │ Pod A   │   │ Pod B   │   │ Pod C   │  无法调度         │
│  │ Pending │   │ Pending │   │ Pending │                   │
│  └────┬────┘   └────┬────┘   └────┬────┘                   │
│       └─────────────┴─────────────┘                         │
│                    │                                         │
│                    ↓                                         │
│         ┌─────────────────┐                                 │
│         │  Cluster        │                                 │
│         │  Autoscaler     │  检测不可调度 Pod               │
│         │                 │  计算需要节点数                 │
│         └────────┬────────┘                                 │
│                  │                                           │
│                  ↓                                           │
│         ┌─────────────────┐                                 │
│         │  Cloud Provider │  创建新节点                     │
│         │  API (AWS/GCP)  │                                 │
│         └─────────────────┘                                 │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 安装 CA

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cluster-autoscaler
spec:
  template:
    spec:
      containers:
        - name: cluster-autoscaler
          image: registry.k8s.io/autoscaling/cluster-autoscaler:v1.30.0
          command:
            - ./cluster-autoscaler
            - --cloud-provider=aws
            - --node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/my-cluster
            - --balance-similar-node-groups
            - --skip-nodes-with-system-pods=false
```

---

## Karpenter (AWS)

下一代节点自动配置工具，比 CA 更快更智能。

```
┌─────────────────────────────────────────────────────────────┐
│              Karpenter vs Cluster Autoscaler                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Cluster Autoscaler          Karpenter                      │
│  ─────────────────           ─────────                      │
│  • 基于节点组 (ASG)           • 直接创建 EC2 实例            │
│  • 预配置实例类型             • 动态选择最优实例类型          │
│  • 分钟级扩容                 • 秒级扩容                     │
│  • 组内统一配置               • 按 Pod 需求配置              │
│                                                              │
│  扩容时间: 2-5 分钟           扩容时间: 15-30 秒             │
└─────────────────────────────────────────────────────────────┘
```

### Karpenter NodePool

```yaml
apiVersion: karpenter.sh/v1
kind: NodePool
metadata:
  name: default
spec:
  template:
    spec:
      requirements:
        - key: karpenter.sh/capacity-type
          operator: In
          values: ["spot", "on-demand"]
        - key: node.kubernetes.io/instance-type
          operator: In
          values: ["m6i.large", "m6i.xlarge", "m5.large"]
      nodeClassRef:
        name: default
  limits:
    cpu: 1000
    memory: 1000Gi
  disruption:
    consolidationPolicy: WhenEmpty
    consolidateAfter: 30s
    expireAfter: 720h
```

---

## KEDA - 事件驱动扩缩

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: queue-scaler
spec:
  scaleTargetRef:
    name: worker
  minReplicaCount: 0
  maxReplicaCount: 100
  triggers:
    - type: aws-sqs-queue
      metadata:
        queueURL: https://sqs.us-east-1.amazonaws.com/123456789/my-queue
        queueLength: "5"
      authenticationRef:
        name: aws-credentials
```

---

## 扩缩容最佳实践

```
1. HPA + CA 组合
   HPA 扩容 Pod → 资源不足 → CA 扩容节点

2. 使用 PriorityClass
   关键服务优先调度，低优先级任务可被驱逐

3. 过度配置 (Overprovisioning)
   运行低优先级占位 Pod，保持缓冲容量

4. 混合实例策略
   On-Demand 用于基线，Spot 用于弹性
```
