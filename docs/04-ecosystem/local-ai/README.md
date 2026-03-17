# 本地 AI/LLM 部署

> 在 Kubernetes 上部署本地大语言模型 - 2025 最新实践

---

## 为什么本地部署 LLM？

| 优势 | 说明 |
|------|------|
| **数据隐私** | 数据不离开本地环境，满足合规要求 |
| **成本控制** | 无 API 调用费用，长期使用更经济 |
| **低延迟** | 本地推理，无需网络往返 |
| **离线可用** | 无需互联网连接，边缘场景可用 |
| **完全控制** | 自定义模型、参数、量化策略 |

---

## 主流本地 LLM 工具

| 工具 | 特点 | 适用场景 | 2025 版本 |
|------|------|----------|----------|
| **Ollama** | 易用，模型管理方便 | 开发、测试、原型 | v0.5+ |
| **vLLM** | 高吞吐，PagedAttention | 生产环境、高并发 | v0.6+ |
| **llama.cpp** | 纯 C++，硬件支持广泛 | 嵌入式、边缘设备 | b4000+ |
| **LocalAI** | API 兼容 OpenAI | API 迁移、多模态 | v2.20+ |
| **LM Studio** | GUI 友好 | 非技术用户 | v0.3+ |
| **TensorRT-LLM** | NVIDIA 优化 | NVIDIA GPU 生产环境 | v0.16+ |

---

## Ollama

### 安装

```bash
# macOS/Linux
curl -fsSL https://ollama.com/install.sh | sh

# Docker
docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama

# Kubernetes
kubectl apply -f https://raw.githubusercontent.com/ollama/ollama/main/kubernetes.yaml
```

### 基础使用

```bash
# 拉取并运行模型
ollama run llama3.2

# 运行特定版本
ollama run llama3.2:70b

# 列出本地模型
ollama list

# 删除模型
ollama rm llama3.2
```

### 自定义模型 (Modelfile)

```dockerfile
# Modelfile
FROM llama3.2

# 设置系统提示词
SYSTEM "You are a helpful AI assistant specialized in Kubernetes."

# 参数调整
PARAMETER temperature 0.7
PARAMETER top_p 0.9
PARAMETER top_k 40
PARAMETER num_ctx 4096
PARAMETER num_gpu 50  # GPU 层数
```

```bash
# 创建自定义模型
ollama create my-k8s-assistant -f Modelfile
ollama run my-k8s-assistant
```

### Kubernetes 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ollama
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ollama
  template:
    metadata:
      labels:
        app: ollama
    spec:
      containers:
        - name: ollama
          image: ollama/ollama:latest
          ports:
            - containerPort: 11434
          volumeMounts:
            - name: ollama-data
              mountPath: /root/.ollama
          resources:
            limits:
              nvidia.com/gpu: 1
      volumes:
        - name: ollama-data
          persistentVolumeClaim:
            claimName: ollama-pvc
```

---

## vLLM - 生产级推理引擎

### 核心特性

vLLM 使用 **PagedAttention** 技术，通过非连续内存块管理 KV Cache，显著提高并发性能。

- **Continuous Batching**: 动态批处理
- **PagedAttention**: 高效 KV Cache 管理
- **Streaming**: 流式 Token 生成
- **OpenAI API 兼容**: 无缝迁移

### Kubernetes 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vllm-llama
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vllm-llama
  template:
    metadata:
      labels:
        app: vllm-llama
    spec:
      nodeSelector:
        nvidia.com/gpu.product: NVIDIA-A100
      containers:
        - name: vllm
          image: vllm/vllm-openai:latest
          args:
            - --model
            - meta-llama/Llama-3.1-70B-Instruct
            - --tensor-parallel-size
            - "2"
            - --max-model-len
            - "4096"
            - --gpu-memory-utilization
            - "0.9"
            - --enable-chunked-prefill
            - --kv-transfer-config
            - '{"kv_connector":"LMCacheConnectorV1", "kv_role":"kv_both"}'
          ports:
            - containerPort: 8000
          resources:
            limits:
              nvidia.com/gpu: 2
              memory: 80Gi
```

---

## KServe + vLLM (推荐生产方案)

```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: llama-vllm
  annotations:
    serving.kserve.io/autoscalerClass: "keda"
spec:
  predictor:
    minReplicas: 0
    maxReplicas: 5
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

## TensorRT-LLM (NVIDIA 优化)

### 模型转换

```bash
# 使用 NVIDIA 容器转换模型
docker run --gpus all -it \
  -v $(pwd)/models:/models \
  nvcr.io/nvidia/tensorrt:24.12-py3 \
  trtllm-build \
  --checkpoint_dir /models/llama-3.1-8b \
  --output_dir /models/llama-3.1-8b-trt \
  --gemm_plugin float16
```

### Kubernetes 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: triton-trt-llm
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: triton
          image: nvcr.io/nvidia/tritonserver:24.12-trtllm-python-py3
          args:
            - tritonserver
            - --model-repository=/models
          ports:
            - containerPort: 8000
          resources:
            limits:
              nvidia.com/gpu: 1
```

---

## 多模态模型部署

### LocalAI 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: localai
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: localai
          image: localai/localai:latest-aio-gpu-nvidia-cuda-12
          env:
            - name: MODELS_PATH
              value: /models
          ports:
            - containerPort: 8080
          volumeMounts:
            - name: models
              mountPath: /models
          resources:
            limits:
              nvidia.com/gpu: 1
      volumes:
        - name: models
          persistentVolumeClaim:
            claimName: models-pvc
```

---

## 2025 工具选择建议

| 场景 | 推荐工具 | 原因 |
|------|----------|------|
| **开发测试** | Ollama | 简单易用，快速迭代 |
| **生产部署 (通用)** | vLLM | 高吞吐，低延迟，PagedAttention |
| **NVIDIA GPU 生产** | TensorRT-LLM | 极致性能优化 |
| **边缘设备** | llama.cpp | 硬件要求低，纯 CPU 运行 |
| **API 迁移** | vLLM/LocalAI | OpenAI API 兼容 |
| **多模态** | LocalAI | 支持图像、音频 |
| **大模型 (405B+)** | KServe Multi-Node | 分布式多节点推理 |

---

## 性能对比 (Llama 3.1 8B, A100)

| 工具 | 延迟 (P50) | 吞吐 (tokens/s) | GPU 显存 | 并发 |
|------|-----------|-----------------|----------|------|
| **Ollama** | 50ms | 80 | 6GB | 低 |
| **vLLM** | 15ms | 300+ | 6GB | 高 |
| **TensorRT-LLM** | 12ms | 400+ | 6GB | 高 |
| **llama.cpp (CPU)** | 100ms | 40 | N/A | 低 |

---

## 优化技巧

### 量化

```bash
# vLLM INT8 量化
vllm serve meta-llama/Llama-3.1-8B --quantization awq

# vLLM GPTQ 量化
vllm serve meta-llama/Llama-3.1-8B --quantization gptq
```

### 批处理优化

```yaml
args:
  - --max-num-seqs  # 最大序列数
  - "256"
  - --max-model-len
  - "8192"
  - --gpu-memory-utilization
  - "0.95"
```

### 前缀缓存

```yaml
args:
  - --enable-prefix-caching  # 启用前缀缓存
```

---

## 参考

- [Ollama 文档](https://github.com/ollama/ollama)
- [vLLM 文档](https://docs.vllm.ai/)
- [TensorRT-LLM 文档](https://nvidia.github.io/TensorRT-LLM/)
- [LocalAI 文档](https://localai.io/)
