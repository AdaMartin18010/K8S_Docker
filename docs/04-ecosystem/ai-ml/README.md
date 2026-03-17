# AI/ML 与 Kubernetes

> 云原生 AI 工作负载管理

---

## K8s 上的 AI/ML 挑战

```
┌─────────────────────────────────────────────────────────────┐
│              AI/ML 工作负载特点 vs K8s 挑战                  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. GPU 资源管理                                             │
│     • 昂贵且稀缺                                             │
│     • 需要共享和调度                                         │
│     • 显存管理复杂                                           │
│                                                              │
│  2. 分布式训练                                               │
│     • Pod 间高速通信                                         │
│     • 检查点/容错                                            │
│     • 弹性伸缩                                               │
│                                                              │
│  3. 大规模作业调度                                           │
│     • 批处理队列                                             │
│     • 优先级和抢占                                            │
│     • 公平共享                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## GPU 调度方案

### NVIDIA GPU Operator

```yaml
apiVersion: nvidia.com/v1
kind: ClusterPolicy
metadata:
  name: gpu-cluster-policy
spec:
  operator:
    defaultRuntime: containerd
  driver:
    enabled: true
  toolkit:
    enabled: true
  devicePlugin:
    enabled: true
```

### GPU 共享 (时间分片)

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: time-slicing-config
data:
  any: |-
    version: v1
    sharing:
      timeSlicing:
        resources:
          - name: nvidia.com/gpu
            replicas: 4  # 1 个 GPU 分为 4 个虚拟 GPU
```

---

## Kubeflow 训练

```yaml
apiVersion: kubeflow.org/v1
kind: PyTorchJob
metadata:
  name: distributed-training
spec:
  pytorchReplicaSpecs:
    Master:
      replicas: 1
      template:
        spec:
          containers:
            - name: pytorch
              image: pytorch/pytorch:latest
              resources:
                limits:
                  nvidia.com/gpu: 2
    Worker:
      replicas: 3
      template:
        spec:
          containers:
            - name: pytorch
              image: pytorch/pytorch:latest
              resources:
                limits:
                  nvidia.com/gpu: 2
```

---

## KServe 推理服务

```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: llm-service
spec:
  predictor:
    model:
      modelFormat:
        name: huggingface
      storageUri: s3://models/llama-7b
    resources:
      limits:
        nvidia.com/gpu: 1
    minReplicas: 1
    maxReplicas: 10
```

---

## AI/ML 工具生态

| 工具 | 用途 |
|------|------|
| **Kubeflow** | ML 工作流平台 |
| **Ray** | 分布式 AI 框架 |
| **Volcano** | 批处理调度 |
| **KServe** | 模型推理服务 |
| **Triton** | NVIDIA 推理服务器 |

---

## 2025 趋势

- **LLM 训练**: 数千 GPU 的分布式训练
- **推理优化**: vLLM、TensorRT-LLM
- **MLOps**: 自动化 ML 生命周期
