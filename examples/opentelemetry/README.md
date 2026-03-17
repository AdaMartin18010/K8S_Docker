# OpenTelemetry 示例

> 云原生可观测性埋点

---

## 架构

```
┌─────────┐    ┌─────────────┐    ┌─────────────┐
│ 应用    │───▶│  OpenTelemetry│───▶│ 后端存储    │
│ (埋点)  │    │  Collector   │    │ (Prometheus/
└─────────┘    └─────────────┘    │  Jaeger/    │
                                   │  Grafana)   │
                                   └─────────────┘
```

---

## 语言示例

### Go 应用

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("my-app")

func handleRequest(ctx context.Context) {
    ctx, span := tracer.Start(ctx, "handle-request")
    defer span.End()
    
    // 业务逻辑
}
```

---

## 目录

- [Go 应用示例](./go-app/)

---

## 相关文档

- [OpenTelemetry 指南](../../docs/05-tools/observability/opentelemetry.md)
