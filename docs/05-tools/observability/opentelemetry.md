# OpenTelemetry

> 云原生可观测性标准 - 2025 最新进展

---

## 什么是 OpenTelemetry？

OpenTelemetry (OTel) 是 CNCF 项目（活跃度仅次于 Kubernetes），提供 vendor-neutral 的观测数据采集标准，统一了 Metrics、Traces、Logs、Profiles 四大支柱。

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
| **Metrics Stability 1.0** | ✅ 已发布 | 指标稳定性保障 |
| **Semantic Conventions 1.0** | ✅ 已发布 | 标准化属性命名 |
| **Collector v1.0** | 🔄 2025 年 | Collector 稳定版 |
| **Profiling** | 🔄 2025 年 | 性能分析信号 |
| **eBPF 集成** | 🔄 实验中 | 内核级数据采集 |
| **Declarative Config 1.0** | ✅ 已发布 | JSON Schema 配置 |
| **GenAI Semantic Conventions** | ✅ 已发布 | AI 应用观测标准 |

---

## OpenTelemetry 与 Prometheus 3.0 集成

Prometheus 3.0 原生支持 OTLP 接收。

```yaml
# Prometheus 3.0 配置
otlp:
  promote_resource_attributes:
    - service.name
    - service.namespace
    - deployment.environment
```

### OTLP 接收端点

```
Prometheus 3.0 默认端点:
  - /api/v1/otlp/v1/metrics

无需 exporter 转换，直接接收 OTLP 格式！
```

---

## 核心概念

### 四大支柱

```
┌─────────────────────────────────────────────────────────────┐
│              可观测性四大支柱                                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  METRICS        TRACES         LOGS         PROFILES        │
│  (指标)        (追踪)          (日志)        (分析)         │
│  ─────────────────────────────────────────────────────────  │
│                                                              │
│  数值测量        请求链路        事件记录       性能分析      │
│  • CPU 使用率    • 请求路径      • 错误日志     • CPU 火焰图  │
│  • 请求速率      • 服务依赖      • 访问日志     • 内存分配    │
│  • 延迟分布      • 耗时分析      • 审计日志     • 锁竞争      │
│                                                              │
│  示例:          示例:          示例:         示例:          │
│  http_requests_  Trace ID:     ERROR...      pprof 采样     │
│  total          abc123                                        │
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
    semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
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

## GenAI 可观测性 (2025 新特性)

### LLM 应用埋点

```python
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

# 配置 Tracer
trace.set_tracer_provider(TracerProvider())
tracer = trace.get_tracer(__name__)

# OTLP Exporter
otlp_exporter = OTLPSpanExporter(endpoint="localhost:4317", insecure=True)
span_processor = BatchSpanProcessor(otlp_exporter)
trace.get_tracer_provider().add_span_processor(span_processor)

# LLM 调用埋点
with tracer.start_as_current_span("llm.chat") as span:
    span.set_attribute("gen_ai.system", "openai")
    span.set_attribute("gen_ai.request.model", "gpt-4")
    span.set_attribute("gen_ai.usage.input_tokens", 150)
    span.set_attribute("gen_ai.usage.output_tokens", 50)

    response = openai_client.chat.completions.create(...)
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
      # 2025 新特性: 基于属性的路由
      routing:
        from_attribute: http.route
        table:
          - value: /api/users
            exporters: [prometheus_users]

    exporters:
      prometheus:
        endpoint: 0.0.0.0:8889
      jaeger:
        endpoint: jaeger-collector:14250
        tls:
          insecure: true
      loki:
        endpoint: http://loki:3100/loki/api/v1/push
      # OTLP 直接发送到 Prometheus 3.0
      otlp/prometheus:
        endpoint: prometheus:4317
        tls:
          insecure: true

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [batch, resource]
          exporters: [jaeger]
        metrics:
          receivers: [otlp]
          processors: [batch]
          exporters: [otlp/prometheus]
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

## 声明式配置 (2025 新特性)

OpenTelemetry Declarative Configuration JSON Schema 1.0 发布。

```json
{
  "$schema": "https://opentelemetry.io/schemas/1.0.0/configuration",
  "resource": {
    "attributes": {
      "service.name": "my-service",
      "service.version": "1.0.0"
    }
  },
  "tracer_provider": {
    "processors": [
      {
        "batch": {
          "exporter": {
            "otlp": {
              "endpoint": "http://localhost:4317"
            }
          }
        }
      }
    ]
  }
}
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

## eBPF 集成 (实验中)

无侵入式内核级数据采集。

```yaml
# OpenTelemetry eBPF Agent
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: otel-ebpf-agent
spec:
  template:
    spec:
      containers:
        - name: ebpf-agent
          image: otel/opentelemetry-ebpf-agent:latest
          securityContext:
            privileged: true
          volumeMounts:
            - name: sys-kernel
              mountPath: /sys/kernel
            - name: debugfs
              mountPath: /sys/kernel/debug
      volumes:
        - name: sys-kernel
          hostPath:
            path: /sys/kernel
        - name: debugfs
          hostPath:
            path: /sys/kernel/debug
```

---

## 最佳实践

1. **使用语义约定**: 遵循 OpenTelemetry Semantic Conventions v1.26+
2. **采样策略**: 生产环境使用概率采样或尾部采样
3. **资源属性**: 统一设置 Service Name、Version、Environment
4. **上下文传播**: 确保跨服务上下文正确传递
5. **避免高基数**: 不要在 metrics 中使用 trace-id 等唯一标签
6. **Collector 扩展**: 大规模部署使用 Collector Gateway 模式
7. **GenAI 观测**: 对 LLM 应用使用专门的语义约定

---

## 2025 趋势

- ✅ **Logs GA**: 日志信号正式发布
- ✅ **Metrics Stability 1.0**: 指标稳定性保障
- ✅ **Collector v1.0**: 预计 2025 年发布
- ✅ **GenAI 语义约定**: LLM 应用观测标准
- 🔄 **Profiling GA**: 性能分析成为第四支柱
- 🔄 **eBPF 集成**: 无侵入式内核级观测
- 🔄 **Declarative Config**: 声明式配置标准
- 🔄 **统一信号**: Metrics/Traces/Logs/Profiles 统一处理

---

## 参考

- [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)
- [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)
- [Semantic Conventions](https://opentelemetry.io/docs/concepts/semantic-conventions/)
- [GenAI Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/)
- [Prometheus 3.0 OTLP](https://prometheus.io/docs/guides/opentelemetry/)
