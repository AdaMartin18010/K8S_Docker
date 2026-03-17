# AI/ML 与 Kubernetes

> 云原生 AI 工作负载管理 - KServe v0.15/v0.16 新特性

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
│  3. LLM 推理服务                                             │
│     • 长连接流式响应                                         │
│     • KV Cache 管理                                          │
│     • 高并发 Token 流                                        │
│     • OpenAI 兼容 API                                        │
│                                                              │
│  4. 大规模作业调度                                           │
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
  mig:
    strategy: single  # 或 mixed
```

### GPU 共享 (时间分片)

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: time-slicing-config
  namespace: nvidia-device-plugin
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

### KServe v0.15 新特性

| 特性 | 说明 |
|------|------|
| **Envoy AI Gateway 集成** | 统一 AI 流量入口 |
| **Multi-Node Inference** | 多节点分布式推理 |
| **LLM Autoscaler with KEDA** | 基于 LLM 指标的自动扩缩容 |
| **Distributed KV Cache** | LMCache 集成，跨副本共享 KV Cache |

### KServe v0.16 新特性

| 特性 | 说明 |
|------|------|
| **LLMInferenceService** | 专为 LLM 优化的 CRD |
| **OpenAI 兼容 API** | `/v1/chat/completions` 端点 |
| **Streaming 响应** | Token 流式传输 |

---

## KServe 基础推理服务

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

## KServe + vLLM (生产级)

```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: llama-vllm
spec:
  predictor:
    minReplicas: 0
    maxReplicas: 3
    scaleMetric: concurrency
    scaleTarget: 1
    containers:
      - name: kserve-container
        image: vllm/vllm-openai:latest
        args:
          - --model
          - meta-llama/Llama-3.1-8B-Instruct
          - --max-model-len
          - "4096"
          - --gpu-memory-utilization
          - "0.9"
        resources:
          limits:
            nvidia.com/gpu: 1
            memory: 24Gi
```

---

## KServe v0.15: Multi-Node Inference

支持单节点无法容纳的大模型 (如 Llama 3.1 405B)。

```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: huggingface-llama3-405b
spec:
  predictor:
    model:
      modelFormat:
        name: huggingface
      storageUri: pvc://llama-3-405b-pvc/hf/405b
    workerSpec:
      pipelineParallelSize: 2  # 流水线并行
      tensorParallelSize: 4    # 张量并行
```

---

## KServe v0.15: LLM Autoscaler with KEDA

基于 vLLM 指标的智能扩缩容。

```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: huggingface-llama3-keda
  annotations:
    serving.kserve.io/autoscalerClass: "keda"
    sidecar.opentelemetry.io/inject: "huggingface-llama3-keda"
spec:
  predictor:
    model:
      modelFormat:
        name: huggingface
      args:
        - --model_name=llama3
        - --model_id=meta-llama/meta-llama-3-70b
    minReplicas: 1
    maxReplicas: 5
    autoScaling:
      metrics:
        - type: PodMetric
          podmetric:
            metric:
              backend: "opentelemetry"
              metricNames:
                - vllm:num_requests_running
              query: "vllm:num_requests_running"
            target:
              type: Value
              value: "4"
```

---

## KServe v0.15: Distributed KV Cache with LMCache

跨副本共享 KV Cache，减少 TTFT (Time To First Token)。

```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: huggingface-llama3-lmcache
spec:
  predictor:
    minReplicas: 2
    model:
      modelFormat:
        name: huggingface
      args:
        - --model_name=llama3
        - --model_id=meta-llama/meta-llama-3-70b
        - --kv-transfer-config
        - '{"kv_connector":"LMCacheConnectorV1", "kv_role":"kv_both"}'
        - --enable-chunked-prefill
```

---

## KServe v0.16: LLMInferenceService

专为 LLM 优化的 CRD，提供 OpenAI 兼容 API。

```yaml
apiVersion: serving.kserve.io/v1alpha1
kind: LLMInferenceService
metadata:
  name: llm-service
spec:
  predictor:
    model:
      modelFormat:
        name: huggingface
      runtime: vllm
      storageUri: s3://models/llama-3-8b
    resources:
      limits:
        nvidia.com/gpu: 1
    minReplicas: 1
    maxReplicas: 5
  # 自动暴露 OpenAI 兼容端点:
  # /v1/chat/completions
  # /v1/completions
  # /v1/models
```

### 调用示例

```bash
# OpenAI 兼容 API 调用
curl http://llm-service.default.svc.cluster.local/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama-3-8b",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ],
    "stream": true
  }'
```

---

## KServe + Envoy AI Gateway

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: ai-gateway
spec:
  gatewayClassName: envoy-ai-gateway
  listeners:
    - name: http
      protocol: HTTP
      port: 80
---
apiVersion: aigateway.envoyproxy.io/v1alpha1
kind: AIServiceBackend
metadata:
  name: llm-backend
spec:
  schema:
    name: OpenAI
  backendRef:
    name: llm-service
    port: 80
```

---

## AI/ML 工具生态

| 工具 | 用途 | 2025 状态 |
|------|------|----------|
| **Kubeflow** | ML 工作流平台 | v1.9 发布 |
| **KServe** | 模型推理服务 | v0.16 发布 |
| **Ray** | 分布式 AI 框架 | v2.40+ |
| **Volcano** | 批处理调度 | 生产就绪 |
| **Triton** | NVIDIA 推理服务器 | v2.50+ |
| **vLLM** | 高吞吐 LLM 推理 | v0.6+ |
| **TensorRT-LLM** | NVIDIA 优化推理 | v0.16+ |
| **LMCache** | 分布式 KV Cache | v0.1+ |

---

## 2025 趋势

- ✅ **LLM 推理优化**: vLLM、TensorRT-LLM 成为标准
- ✅ **KServe v0.16**: LLMInferenceService，OpenAI 兼容 API
- ✅ **Multi-Node Inference**: 405B+ 参数模型支持
- ✅ **KV Cache 共享**: LMCache 降低推理成本
- 🔄 **MLOps**: 自动化 ML 生命周期
- 🔄 **AI Gateway**: 统一 AI 流量管理
- 🔄 **Spot 实例训练**: 成本优化

---

## 参考

- [KServe 官方文档](https://kserve.github.io/website/)
- [KServe v0.15 发布说明](https://kserve.github.io/website/blog/kserve-0.15-release/)
- [vLLM 文档](https://docs.vllm.ai/)
- [Envoy AI Gateway](https://aigateway.envoyproxy.io/)
