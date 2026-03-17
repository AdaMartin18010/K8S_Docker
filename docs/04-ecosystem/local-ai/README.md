# 本地 AI/LLM 部署

> 在 Kubernetes 上部署本地大语言模型

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

| 工具 | 特点 | 适用场景 |
|------|------|----------|
| **Ollama** | 易用，模型管理方便 | 开发、测试、原型 |
| **vLLM** | 高吞吐，PagedAttention | 生产环境、高并发 |
| **llama.cpp** | 纯 C++，硬件支持广泛 | 嵌入式、边缘设备 |
| **LocalAI** | API 兼容 OpenAI | API 迁移、多模态 |
| **LM Studio** | GUI 友好 | 非技术用户 |

---

## Ollama

### 安装

```bash
# macOS/Linux
curl -fsSL https://ollama.com/install.sh | sh

# Docker
docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
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
```

```bash
# 创建自定义模型
ollama create my-k8s-assistant -f Modelfile
ollama run my-k8s-assistant
```

---

## vLLM - 生产级推理引擎

### 核心特性

vLLM 使用 PagedAttention 技术，通过非连续内存块管理 KV Cache，显著提高并发性能。

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
          ports:
            - containerPort: 8000
          resources:
            limits:
              nvidia.com/gpu: 2
              memory: 80Gi
```

---

## KServe + vLLM

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
        resources:
          limits:
            nvidia.com/gpu: 1
            memory: 24Gi
```

---

## 2025 工具选择建议

| 场景 | 推荐工具 | 原因 |
|------|----------|------|
| **开发测试** | Ollama | 简单易用，快速迭代 |
| **生产部署** | vLLM | 高吞吐，低延迟 |
| **边缘设备** | llama.cpp | 硬件要求低，纯 CPU 运行 |
| **API 迁移** | LocalAI | OpenAI API 兼容 |
| **多模态** | LocalAI | 支持图像、音频 |

---

## 性能对比 (Llama 3.1 8B)

| 工具 | 延迟 (P50) | 吞吐 (tokens/s) | GPU 显存 |
|------|-----------|-----------------|----------|
| **Ollama** | 50ms | 80 | 6GB |
| **vLLM** | 20ms | 200 | 6GB |
| **llama.cpp** | 100ms | 40 (CPU) | N/A |
