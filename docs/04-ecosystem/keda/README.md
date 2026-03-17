# KEDA - Kubernetes 事件驱动自动扩缩容

## 概述

KEDA (Kubernetes Event-driven Autoscaling) 是 CNCF Incubating 项目，专为事件驱动的工作负载提供强大的自动扩缩容能力。与传统的基于 CPU/内存的 HPA 不同，KEDA 可以根据外部事件源（如 Kafka、RabbitMQ、数据库等）的指标进行扩缩容，甚至支持缩容到零。

> **关键数据**: KEDA 支持 70+ 种事件源 scaler，与 Azure Functions、AWS Lambda 等 Serverless 平台无缝集成。

## 架构原理

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              KEDA 架构                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      KEDA Operator                                  │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                   ScaledObject Controller                     │  │   │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │  │   │
│  │  │  │   Scaler 1   │  │   Scaler 2   │  │   Scaler N   │        │  │   │
│  │  │  │  (Kafka)     │  │  (RabbitMQ)  │  │  (PostgreSQL)│        │  │   │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘        │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                              │                                      │   │
│  │  ┌───────────────────────────▼───────────────────────────────────┐  │   │
│  │  │              Metrics Server (external.metrics.k8s.io)          │  │   │
│  │  └───────────────────────────┬───────────────────────────────────┘  │   │
│  └──────────────────────────────┼──────────────────────────────────────┘   │
│                                 │                                           │
│  ┌──────────────────────────────▼──────────────────────────────────────┐   │
│  │                     Kubernetes HPA                                  │   │
│  │                  (基于外部指标扩缩容)                                │   │
│  └──────────────────────────────┬──────────────────────────────────────┘   │
│                                 │                                           │
│                                 ▼                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Target Deployment                               │   │
│  │              (0 → N 副本，基于事件负载自动调整)                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 支持的 Scaler（事件源）

| 类别 | Scaler | 描述 |
|------|--------|------|
| **消息队列** | kafka | Apache Kafka 消费者组延迟 |
| | rabbitmq | RabbitMQ 队列长度 |
| | aws-sqs | AWS SQS 队列深度 |
| | azure-servicebus | Azure Service Bus |
| | gcp-pubsub | Google Cloud Pub/Sub |
| | pulsar | Apache Pulsar |
| **数据库** | postgresql | PostgreSQL 查询结果 |
| | mysql | MySQL 查询结果 |
| | mongodb | MongoDB 查询结果 |
| | redis | Redis 列表长度 |
| **云原生** | kubernetes-workload | 基于其他 workload 的 Pod 数量 |
| | metrics-api | 通用指标 API |
| | prometheus | Prometheus 指标 |
| **定时任务** | cron | 基于时间的扩缩容 |
| **其他** | github-runner | GitHub Actions Runner |
| | liiklus | Liiklus 流处理 |

## 安装 KEDA

### Helm 安装

```bash
# 添加 Helm 仓库
helm repo add kedacore https://kedacore.github.io/charts
helm repo update

# 安装 KEDA
helm install keda kedacore/keda \
  --namespace keda \
  --create-namespace \
  --set resources.operator.requests.cpu=100m \
  --set resources.operator.requests.memory=128Mi
```

### 验证安装

```bash
# 检查 KEDA Pod
kubectl get pods -n keda

# 查看 ScaledObject CRD
kubectl get crd scaledobjects.keda.sh

# 检查 metrics server
kubectl get apiservice v1beta1.external.metrics.k8s.io
```

## 基础使用示例

### 1. Kafka 事件驱动扩缩容

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: kafka-scaled-app
  namespace: default
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: kafka-consumer

  # 轮询间隔（秒）
  pollingInterval: 30

  # 冷却时间（秒）- 触发缩容前等待时间
  cooldownPeriod: 300

  # 最小/最大副本数
  minReplicaCount: 0  # 可以缩容到 0
  maxReplicaCount: 50

  triggers:
  - type: kafka
    metadata:
      bootstrapServers: kafka:9092
      consumerGroup: my-consumer-group
      topic: orders
      # 每个副本处理的消息数
      lagThreshold: "100"
      # 激活阈值（低于此值开始缩容）
      activationLagThreshold: "10"
    authenticationRef:
      name: kafka-trigger-auth
---
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: kafka-trigger-auth
  namespace: default
spec:
  secretTargetRef:
  - parameter: sasl
    name: kafka-secret
    key: sasl
  - parameter: username
    name: kafka-secret
    key: username
  - parameter: password
    name: kafka-secret
    key: password
```

### 2. RabbitMQ 队列深度扩缩容

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: rabbitmq-scaled-app
spec:
  scaleTargetRef:
    name: order-processor
  pollingInterval: 10
  cooldownPeriod: 30
  minReplicaCount: 1
  maxReplicaCount: 100

  # 高级扩缩容行为
  advanced:
    horizontalPodAutoscalerConfig:
      behavior:
        scaleDown:
          stabilizationWindowSeconds: 300
          policies:
          - type: Percent
            value: 50
            periodSeconds: 60
          - type: Pods
            value: 2
            periodSeconds: 60
          selectPolicy: Min
        scaleUp:
          stabilizationWindowSeconds: 0
          policies:
          - type: Percent
            value: 100
            periodSeconds: 30
          - type: Pods
            value: 10
            periodSeconds: 30
          selectPolicy: Max

  triggers:
  - type: rabbitmq
    metadata:
      protocol: amqp
      queueName: orders
      # 每个 Pod 处理的队列长度
      queueLength: "50"
    authenticationRef:
      name: rabbitmq-trigger-auth
```

### 3. PostgreSQL 查询驱动扩缩容

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: postgres-scaled-app
spec:
  scaleTargetRef:
    name: job-processor
  minReplicaCount: 0
  maxReplicaCount: 20
  triggers:
  - type: postgresql
    metadata:
      host: postgres.default.svc.cluster.local
      port: "5432"
      dbName: jobs
      sslmode: disable
      # 查询语句
      query: "SELECT COUNT(*) FROM pending_jobs WHERE status='pending'"
      # 目标值：查询结果除以 targetQueryValue 得到副本数
      targetQueryValue: "10"
    authenticationRef:
      name: postgres-trigger-auth
---
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: postgres-trigger-auth
spec:
  secretTargetRef:
  - parameter: connectionString
    name: postgres-secret
    key: connection-string
```

### 4. 多触发器组合

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: multi-trigger-app
spec:
  scaleTargetRef:
    name: adaptive-processor
  minReplicaCount: 2
  maxReplicaCount: 100

  triggers:
  # 触发器 1：Kafka 延迟
  - type: kafka
    name: kafka-trigger
    metadata:
      bootstrapServers: kafka:9092
      consumerGroup: my-group
      topic: events
      lagThreshold: "100"

  # 触发器 2：CPU 使用率
  - type: cpu
    name: cpu-trigger
    metadata:
      type: Utilization
      value: "70"

  # 触发器 3：定时预热
  - type: cron
    name: cron-trigger
    metadata:
      timezone: Asia/Shanghai
      start: 0 8 * * 1-5    # 工作日早8点
      end: 0 20 * * 1-5     # 工作日晚8点
      desiredReplicas: "10"
```

## 身份验证配置

### Secret 引用方式

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: kafka-secret
type: Opaque
stringData:
  sasl: "plaintext"
  username: "kafka-user"
  password: "kafka-password"
---
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: kafka-auth
spec:
  secretTargetRef:
  - parameter: sasl
    name: kafka-secret
    key: sasl
  - parameter: username
    name: kafka-secret
    key: username
  - parameter: password
    name: kafka-secret
    key: password
```

### Workload Identity（推荐）

```yaml
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: aws-irsa-auth
spec:
  podIdentity:
    provider: aws
    identityOwner: keda-operator  # 或 workload
---
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: sqs-scaled-app
spec:
  scaleTargetRef:
    name: sqs-processor
  triggers:
  - type: aws-sqs-queue
    authenticationRef:
      name: aws-irsa-auth
    metadata:
      queueURL: https://sqs.us-east-1.amazonaws.com/123456789/my-queue
      queueLength: "10"
      awsRegion: us-east-1
```

## 高级配置

### 缩容到零（Scale to Zero）

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: scale-to-zero-app
spec:
  scaleTargetRef:
    name: event-processor
  minReplicaCount: 0  # 关键：允许缩容到 0
  maxReplicaCount: 50
  triggers:
  - type: kafka
    metadata:
      bootstrapServers: kafka:9092
      consumerGroup: my-group
      topic: events
      lagThreshold: "100"
      # 激活阈值：低于此值开始缩容到 0
      activationLagThreshold: "10"
```

**缩容到零的注意事项**:

1. 应用必须能处理冷启动延迟
2. 第一个事件到达后会有启动时间
3. 建议使用 `activationLagThreshold` 避免频繁启停

### 空闲工作负载检测

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: idle-detection-app
spec:
  scaleTargetRef:
    name: worker
  idleReplicaCount: 0  # 空闲时副本数
  minReplicaCount: 1   # 非空闲时最小副本
  maxReplicaCount: 20
  triggers:
  - type: metrics-api
    metadata:
      targetValue: "1"
      url: "http://api-server/metrics/active-tasks"
      valueLocation: "data.active_count"
```

## 监控和告警

### Prometheus Metrics

```yaml
# ServiceMonitor 配置
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: keda-metrics
  namespace: keda
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: keda-operator
  endpoints:
  - port: metrics
    interval: 30s
```

### 关键指标

```promql
# 当前副本数
keda_scaled_object_current_replicas{scaledObject="my-app"}

# 期望副本数
keda_scaled_object_desired_replicas{scaledObject="my-app"}

# Scaler 错误率
rate(keda_scaler_errors_total[5m])

# 外部指标值
keda_scaler_metrics_value{scaledObject="my-app"}

# 扩缩容耗时
histogram_quantile(0.95, rate(keda_scaler_scaling_duration_bucket[5m]))
```

### 告警规则

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: keda-alerts
spec:
  groups:
  - name: keda
    rules:
    - alert: KedaScalerErrors
      expr: rate(keda_scaler_errors_total[5m]) > 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "KEDA scaler errors detected"

    - alert: KedaScalingLag
      expr: |
        keda_scaled_object_desired_replicas
        - keda_scaled_object_current_replicas > 5
      for: 10m
      labels:
        severity: critical
      annotations:
        summary: "KEDA scaling lag detected"
```

## 故障排查

```bash
# 查看 ScaledObject 状态
kubectl get scaledobjects
kubectl describe scaledobject <name>

# 查看生成的 HPA
kubectl get hpa -n <namespace>

# 查看 KEDA operator 日志
kubectl logs -n keda deployment/keda-operator

# 查看 metrics server 日志
kubectl logs -n keda deployment/keda-metrics-apiserver

# 检查外部指标
kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/s0-kafka-my-topic" | jq .

# 手动触发扩缩容测试
kubectl patch scaledobject my-app --type merge -p '{"spec":{"minReplicaCount":5}}'
```

## 最佳实践

1. **合理设置 pollingInterval**
   - 太短：增加外部系统负载
   - 太长：扩缩容延迟增加
   - 推荐：10-30 秒

2. **配置 cooldownPeriod**
   - 防止频繁扩缩容（flapping）
   - 根据应用启动时间调整
   - 推荐：300 秒（5 分钟）

3. **使用 activationLagThreshold**
   - 避免微量消息导致扩容
   - 特别适用于缩容到零的场景

4. **资源限制**
   - 为被扩容的 Pod 设置适当的 resource requests
   - 确保集群有足够资源承载最大副本数

5. **多触发器策略**
   - 事件驱动 + CPU/内存兜底
   - Cron 触发器用于预热

## 总结

| 场景 | 推荐 Scaler | 配置要点 |
|------|-------------|----------|
| 消息队列处理 | kafka/rabbitmq | lagThreshold, activationLagThreshold |
| 定时任务预热 | cron | timezone, start/end |
| 数据库任务队列 | postgresql/mysql | targetQueryValue |
| 混合负载 | 多触发器 | 组合事件和指标触发器 |
| Serverless | kafka + minReplicaCount=0 | activationLagThreshold |

KEDA 是实现真正的云原生事件驱动架构的关键组件，特别适合处理异步任务、流式数据和 Serverless 工作负载。
