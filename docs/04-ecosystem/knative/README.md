# Knative - Kubernetes Serverless 平台

> 基于 Kubernetes 的 Serverless 应用框架

---

## 什么是 Knative？

Knative 是一个开源的 Kubernetes 平台，用于构建、部署和管理 Serverless 工作负载。它由 Google 发起，现为 CNCF 项目，支持超过 50 家公司贡献。

```
┌─────────────────────────────────────────────────────────────┐
│              Knative 架构                                    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                 Knative Serving                       │  │
│  │  • 自动扩缩容 (0 → n)                                   │  │
│  │  • 流量管理 (金丝雀、蓝绿部署)                          │  │
│  │  • 服务路由                                            │  │
│  │  • 自动版本管理                                         │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                Knative Eventing                       │  │
│  │  • 事件源 (Sources)                                     │  │
│  │  • 事件通道 (Channels)                                  │  │
│  │  • 事件订阅 (Subscriptions)                             │  │
│  │  • 事件代理 (Broker/Trigger)                            │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                Knative Functions                      │  │
│  │  • 函数即服务 (FaaS)                                    │  │
│  │  • 多语言支持 (Go/Python/Node.js/Quarkus)               │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Knative vs 传统 K8s

| 特性 | 传统 Kubernetes | Knative |
|------|-----------------|---------|
| **扩缩容** | 手动或 HPA | 自动 (0→n) |
| **流量管理** | 手动配置 | 内置流量分割 |
| **版本管理** | 手动 | 自动版本管理 |
| **事件驱动** | 无 | 内置 Eventing |
| **冷启动** | 不适用 | 优化快速启动 |

---

## 安装 Knative

```bash
# 安装 Knative Serving
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.13.1/serving-crds.yaml
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.13.1/serving-core.yaml

# 安装网络层 (Kourier)
kubectl apply -f https://github.com/knative/net-kourier/releases/download/knative-v1.13.0/kourier.yaml
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'

# 安装 Knative Eventing
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.13.1/eventing-crds.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.13.1/eventing-core.yaml

# 安装 CLI
curl -L https://github.com/knative/client/releases/download/knative-v1.13.0/kn-linux-amd64 -o kn
chmod +x kn && sudo mv kn /usr/local/bin/
```

---

## Knative Serving

### 基础服务

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: hello-world
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/minScale: "0"  # 可缩容到 0
        autoscaling.knative.dev/maxScale: "100"
    spec:
      containers:
        - image: gcr.io/knative-samples/helloworld-go
          ports:
            - containerPort: 8080
          env:
            - name: TARGET
              value: "World"
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 256Mi
```

### 流量管理

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: hello-world
spec:
  template:
    metadata:
      name: hello-world-v2  # 版本名
    spec:
      containers:
        - image: gcr.io/knative-samples/helloworld-go:v2
          env:
            - name: TARGET
              value: "v2"
  traffic:
    - tag: current
      revisionName: hello-world-v1
      percent: 90
    - tag: candidate
      revisionName: hello-world-v2
      percent: 10
    - tag: latest
      latestRevision: true
      percent: 0
```

### 自动扩缩容配置

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: auto-scale-demo
spec:
  template:
    metadata:
      annotations:
        # 扩缩容窗口
        autoscaling.knative.dev/window: "60s"
        # 并发目标
        autoscaling.knative.dev/target-concurrency: "10"
        # 并发限制
        autoscaling.knative.dev/scale-to-zero-pod-retention-period: "0s"
        # 指标类型
        autoscaling.knative.dev/metric: "concurrency"
        # 目标利用率
        autoscaling.knative.dev/target-utilization-percentage: "70"
    spec:
      containers:
        - image: myapp:latest
```

---

## Knative Eventing

### 事件源 (Source)

```yaml
# Kafka 事件源
apiVersion: sources.knative.dev/v1
kind: KafkaSource
metadata:
  name: kafka-source
spec:
  consumerGroup: knative-group
  bootstrapServers:
    - kafka:9092
  topics:
    - my-topic
  sink:
    ref:
      apiVersion: serving.knative.dev/v1
      kind: Service
      name: event-consumer
```

### 事件代理 (Broker/Trigger)

```yaml
# 创建 Broker
apiVersion: eventing.knative.dev/v1
kind: Broker
metadata:
  name: default
  namespace: default
---
# 创建 Trigger
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: order-trigger
spec:
  broker: default
  filter:
    attributes:
      type: order.created
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1
      kind: Service
      name: order-processor
```

---

## Knative Functions

### 创建函数

```bash
# 创建 Go 函数
kn func create --language go --template http my-function

# 创建 Python 函数
kn func create --language python --template http my-python-func

# 构建并部署
kn func build --registry myregistry.io
kn func deploy
```

### 函数示例 (Go)

```go
package function

import (
    "context"
    "net/http"
)

// Handle 处理 HTTP 请求
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
    // 获取 CloudEvent 数据 (如果使用事件触发)
    // event := cloudevents.FromContext(ctx)

    res.Write([]byte("Hello, World!"))
}
```

---

## 与 KServe 集成

```yaml
# KServe InferenceService (基于 Knative)
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: llm-service
spec:
  predictor:
    minReplicas: 0
    maxReplicas: 10
    scaleMetric: concurrency
    scaleTarget: 1
    containers:
      - name: predictor
        image: vllm/vllm-openai:latest
        args:
          - --model
          - meta-llama/Llama-2-7b
        resources:
          limits:
            nvidia.com/gpu: 1
            memory: 16Gi
```

---

## 2025 应用场景

| 场景 | 优势 |
|------|------|
| **Web 应用** | 自动扩缩容，节省成本 |
| **API 网关** | 流量分割，A/B 测试 |
| **事件处理** | 事件驱动，解耦架构 |
| **AI/ML 推理** | 冷启动优化，GPU 自动扩缩容 |
| **批处理任务** | 按需启动，用完即停 |

---

## 云服务集成

| 云厂商 | 服务 | 基于 Knative |
|--------|------|--------------|
| **Google** | Cloud Run | ✅ |
| **IBM** | Cloud Code Engine | ✅ |
| **Red Hat** | OpenShift Serverless | ✅ |
| **阿里云** | Knative on ACK | ✅ |
