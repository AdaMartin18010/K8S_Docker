# Dapr - 分布式应用运行时

> 简化微服务开发的便携式运行时

---

## 什么是 Dapr？

Dapr (Distributed Application Runtime) 是一个可移植的、事件驱动的运行时，让任何开发者都能轻松构建弹性、分布式应用。

```
┌─────────────────────────────────────────────────────────────┐
│                    Dapr 架构                                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                  Dapr Building Blocks                 │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │  • Service Invocation        • State Management       │  │
│  │  • Publish/Subscribe         • Bindings               │  │
│  │  • Actors                    • Observability          │  │
│  │  • Secrets Management        • Configuration          │  │
│  │  • Workflow                  • Distributed Lock       │  │
│  │  • Conversation (AI)         • Cryptography           │  │
│  └──────────────────────────────────────────────────────┘  │
│                          │                                   │
│  ┌───────────────────────┼───────────────────────────────┐  │
│  │                       │                    Dapr Sidecar │  │
│  │  ┌─────────────┐  ┌───▼────┐  ┌─────────────┐          │  │
│  │  │   App A     │  │  Dapr  │  │   App B     │          │  │
│  │  │   (Go)      │◄─┤ Sidecar├─►│   (Python)  │          │  │
│  │  └─────────────┘  └───┬────┘  └─────────────┘          │  │
│  │                       │                                  │  │
│  └───────────────────────┼───────────────────────────────┘  │
│                          │                                   │
│  ┌───────────────────────▼───────────────────────────────┐  │
│  │              基础设施组件 (可插拔)                       │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐      │  │
│  │  │  Redis  │ │  Kafka  │ │  AWS S3 │ │  Azure  │      │  │
│  │  │ (State) │ │ (Pub/Sub)│ │ (Blob)  │ │  KeyVault│     │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘      │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                              │
│  优势:                                                       │
│  • 语言无关 (Go/Java/Python/.NET/JS...)                      │
│  • 平台无关 (本地/K8s/VM/边缘)                                │
│  • 基础设施抽象，易于切换                                     │
│  • 内置最佳实践 (mTLS、重试、熔断)                            │
└─────────────────────────────────────────────────────────────┘
```

---

## Dapr 2025 数据 (CNCF State of Dapr Report)

| 指标 | 数据 |
|------|------|
| **开发者生产力提升** | 96% 开发者节省时间，60% 提升 30%+ |
| **生产环境采用率** | 近半数团队在生产环境运行 Dapr |
| **多云策略优势** | AWS 使用率从 22% 增至 38% |
| **AI 集成** | LLM Conversation API 发布 |

---

## 核心构建块

### 1. 服务调用 (Service Invocation)

```yaml
# Dapr 服务调用 - 无需知道目标地址
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: service-invocation
spec:
  type: middleware.http.ratelimit
  version: v1
```

```go
// Go SDK 示例
client, _ := dapr.NewClient()
defer client.Close()

// 调用其他服务
resp, err := client.InvokeMethod(
    ctx,
    "order-service",    // 目标服务名
    "api/orders",       // 方法
    "post",             // HTTP 方法
    orderData,
)
```

### 2. 状态管理 (State Management)

```yaml
# Redis 状态存储组件
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: statestore
spec:
  type: state.redis
  version: v1
  metadata:
    - name: redisHost
      value: redis:6379
    - name: redisPassword
      secretKeyRef:
        name: redis-secret
        key: password
```

```go
// 保存状态
err := client.SaveState(ctx, "statestore", "key1", []byte("value"))

// 获取状态
resp, err := client.GetState(ctx, "statestore", "key1")

// 删除状态
err := client.DeleteState(ctx, "statestore", "key1")
```

### 3. 发布订阅 (Pub/Sub)

```yaml
# Kafka Pub/Sub 组件
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: pubsub
spec:
  type: pubsub.kafka
  version: v1
  metadata:
    - name: brokers
      value: kafka:9092
    - name: consumerGroup
      value: my-group
```

```go
// 发布消息
err := client.PublishEvent(ctx, "pubsub", "orders", orderData)

// 订阅消息 (HTTP 端点)
// POST /dapr/subscribe
[
  {
    "pubsubname": "pubsub",
    "topic": "orders",
    "route": "/orders"
  }
]
```

### 4. 工作流 (Workflow)

```go
// Dapr 工作流示例
func OrderProcessingWorkflow(ctx *workflow.WorkflowContext) (string, error) {
    var order Order
    if err := ctx.GetInput(&order); err != nil {
        return "", err
    }

    // 调用活动 - 支付处理
    var paymentResult PaymentResult
    err := ctx.CallActivity("ProcessPayment", order.Payment).Await(&paymentResult)
    if err != nil {
        return "", err
    }

    // 调用活动 - 库存预留
    var inventoryResult InventoryResult
    err = ctx.CallActivity("ReserveInventory", order.Items).Await(&inventoryResult)
    if err != nil {
        // 补偿 - 退款
        ctx.CallActivity("RefundPayment", order.Payment).Await(nil)
        return "", err
    }

    // 调用活动 - 发货
    ctx.CallActivity("ShipOrder", order).Await(nil)

    return "order processed", nil
}
```

---

## K8s 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
      annotations:
        dapr.io/enabled: "true"
        dapr.io/app-id: "myapp"
        dapr.io/app-port: "8080"
        dapr.io/config: "app-config"
    spec:
      containers:
        - name: myapp
          image: myapp:latest
          ports:
            - containerPort: 8080
```

---

## Dapr 1.16 (2025) 新特性

| 特性 | 说明 |
|------|------|
| **工作流性能提升** | 性能提升 40%，支持跨应用编排 |
| **Conversation API** | LLM 工具调用，AI 集成 |
| **Sentry JWT/OIDC** | 强化身份认证 |
| **.NET Roslyn 分析器** | 更好的开发体验 |
| **Java 补偿工作流** | Saga 模式支持 |

---

## Dapr vs Service Mesh

| 对比项 | Dapr | Service Mesh (Istio) |
|--------|------|----------------------|
| **主要目的** | 开发者抽象 | 网络基础设施 |
| **开发者参与** | 主动调用 API | 透明代理 |
| **功能范围** | 广泛 (状态、Pub/Sub、工作流) | 网络 (流量、安全) |
| **学习曲线** | 中 (需学习 API) | 低 (对应用透明) |
| **最佳场景** | 新应用开发 | 现有应用治理 |

---

## 2025 趋势

- **AI 集成**: Dapr Agents, LLM Conversation API
- **边缘计算**: 更小的资源占用，边缘优化
- **多运行时模式**: 进程内 + Sidecar 混合部署
- **工作流编排**: 复杂业务流程的标准方案
