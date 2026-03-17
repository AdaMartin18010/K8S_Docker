# OpenTelemetry

> 云原生可观测性标准

---

## 什么是 OpenTelemetry？

OpenTelemetry (OTel) 是 CNCF 项目，提供 vendor-neutral 的观测数据采集标准，统一了 Metrics、Traces、Logs 三大支柱。

```
┌─────────────────────────────────────────────────────────────┐
│              OpenTelemetry 架构                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  应用层                                                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │   Go     │ │   Java   │ │  Python  │ │  Node.js │       │
│  │   App    │ │   App    │ │   App    │ │   App    │       │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘       │
│       │            │            │            │              │
│  ┌────▼────────────▼────────────▼────────────▼─────┐       │
│  │              OpenTelemetry SDK                  │       │
│  │  • Auto-Instrumentation  • Manual Instrumentation│       │
│  │  • Resource Detection    • Context Propagation  │       │
│  └─────────────────────┬───────────────────────────┘       │
│                        │ OTLP                              │
│  ┌─────────────────────▼───────────────────────────┐       │
│  │           OpenTelemetry Collector               │       │
│  │  • Receive  • Process  • Export                │       │
│  │  • Batch    • Filter   • Transform             │       │
│  └─────────────────────┬───────────────────────────┘       │
│                        │                                   │
│  ┌─────────────────────▼───────────────────────────┐       │
│  │              后端存储 (可切换)                    │       │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐          │       │
│  │  │Prometheus│ │  Jaeger │ │  Grafana │          │       │
│  │  │(Metrics)│ │(Traces) │ │  (Logs)  │          │       │
│  │  └─────────┘ └─────────┘ └─────────┘          │       │
│  └─────────────────────────────────────────────────┘       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 2025 OpenTelemetry 里程碑

| 特性 | 状态 | 说明 |
|------|------|------|
| **Logs GA** | ✅ 已发布 | 日志信号正式发布 |
| **Semantic Conventions 1.0** | ✅ 已发布 | 标准化属性命名 |
| **Profiling** | 🔄 2025 中 | 性能分析信号 |
| **eBPF 集成** | 🔄 实验中 | 内核级数据采集 |
| **WASI 支持** | 🔄 实验中 | WebAssembly 集成 |

---

## 核心概念

### 三大支柱

```
┌─────────────────────────────────────────────────────────────┐
│              可观测性三大支柱                                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  METRICS                    TRACES                  LOGS    │
│  (指标)                    (追踪)                  (日志)   │
│  ─────────────────────────────────────────────────────────  │
│                                                              │
│  数值测量                    请求链路                 事件记录 │
│  • CPU 使用率               • 请求路径              • 错误日志 │
│  • 请求速率                 • 服务依赖              • 访问日志 │
│  • 延迟分布                 • 耗时分析              • 审计日志 │
│                                                              │
│  示例:                      示例:                   示例:     │
│  http_requests_total        Trace ID: abc123        ERROR... │
│  http_request_duration      Span: GET /api          WARN...  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 关键组件

| 组件 | 说明 |
|------|------|
| **Resource** | 描述产生遥测数据的实体 (Service、Host、K8s) |
| **Span** | 追踪中的单个操作单元 |
| **Trace** | 跨服务的请求链路 (由多个 Span 组成) |
| **Metric** | 可聚合的数值测量 |
| **Log Record** | 带时间戳的事件记录 |
| **Baggage** | 跨服务传播的上下文信息 |

---

## Go 应用集成

### 自动埋点

```go
package main

import (
    "context"
    "log"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func main() {
    ctx := context.Background()

    // 创建 OTLP Exporter
    exp, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint("localhost:4317"),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 创建 Tracer Provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName("my-service"),
            semconv.ServiceVersion("1.0.0"),
            semconv.DeploymentEnvironment("production"),
        )),
    )
    defer tp.Shutdown(ctx)

    otel.SetTracerProvider(tp)

    // 使用 Tracer
    tracer := otel.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "process-order")
    defer span.End()

    // 业务逻辑
    processOrder(ctx)
}
```

### HTTP 服务埋点

```go
import (
    "net/http"

    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
    // 包装 Handler 自动埋点
    handler := otelhttp.NewHandler(
        http.HandlerFunc(myHandler),
        "server-request",
    )

    http.Handle("/api/", handler)
    http.ListenAndServe(":8080", nil)
}
```

---

## K8s 部署

### OpenTelemetry Collector

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: my-collector
spec:
  mode: deployment
  config: |
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
          http:
            endpoint: 0.0.0.0:4318

    processors:
      batch:
        timeout: 1s
        send_batch_size: 1024
      resource:
        attributes:
          - key: environment
            value: production
            action: upsert

    exporters:
      prometheus:
        endpoint: 0.0.0.0:8889
      jaeger:
        endpoint: jaeger-collector:14250
        tls:
          insecure: true
      loki:
        endpoint: http://loki:3100/loki/api/v1/push

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [batch, resource]
          exporters: [jaeger]
        metrics:
          receivers: [otlp]
          processors: [batch]
          exporters: [prometheus]
        logs:
          receivers: [otlp]
          processors: [batch]
          exporters: [loki]
```

### 自动注入 Sidecar

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    metadata:
      annotations:
        instrumentation.opentelemetry.io/inject-sdk: "true"
    spec:
      containers:
        - name: myapp
          image: myapp:latest
```

---

## 与 Prometheus/Grafana 集成

```yaml
# 使用 OpenTelemetry 生成 Prometheus 指标
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  metricstransform:
    transforms:
      - include: http_server_duration
        match_type: regexp
        action: update
        operations:
          - action: aggregate_labels
            label_set: [http_method, http_status_code]
            aggregation_type: sum

exporters:
  prometheus:
    endpoint: 0.0.0.0:8889
    namespace: otel
    const_labels:
      app: myapp

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [metricstransform]
      exporters: [prometheus]
```

---

## 最佳实践

1. **使用语义约定**: 遵循 OpenTelemetry Semantic Conventions
2. **采样策略**: 生产环境使用概率采样或尾部采样
3. **资源属性**: 统一设置 Service Name、Version、Environment
4. **上下文传播**: 确保跨服务上下文正确传递
5. **避免高基数**: 不要在 metrics 中使用 trace-id 等唯一标签

---

## 2025 趋势

- **Profiling GA**: 性能分析成为第四支柱
- **eBPF 集成**: 无侵入式内核级观测
- **AI 可观测性**: LLM 应用专门优化
- **统一信号**: Metrics/Traces/Logs/Profiles 统一处理
