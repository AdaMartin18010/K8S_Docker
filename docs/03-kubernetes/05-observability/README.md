# Kubernetes 可观测性

> 监控、日志和追踪

---

## 可观测性三大支柱

```
┌─────────────────────────────────────────────┐
│           Observability                     │
├─────────────┬───────────────┬───────────────┤
│   Metrics   │     Logs      │    Traces     │
│  (指标)     │    (日志)     │    (追踪)     │
├─────────────┼───────────────┼───────────────┤
• CPU/Memory  • stdout/stderr • 请求链路      │
• 响应时间    • 应用日志      • 调用关系      │
• 错误率      • 系统日志      • 延迟分析      │
└─────────────┴───────────────┴───────────────┘
```

---

## 监控方案

### Prometheus + Grafana

```yaml
# ServiceMonitor 配置
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: app-metrics
spec:
  selector:
    matchLabels:
      app: myapp
  endpoints:
    - port: metrics
      path: /metrics
      interval: 30s
```

---

## 日志方案

```
Pod stdout/stderr
       ↓
Node Log Agent (Fluent Bit/DaemonSet)
       ↓
Log Storage (Elasticsearch/Loki)
       ↓
Log Query (Kibana/Grafana)
```

---

## 推荐工具栈

| 功能 | 工具 |
|------|------|
| 监控 | Prometheus + Grafana |
| 日志 | Fluent Bit + Loki |
| 追踪 | OpenTelemetry + Jaeger |
