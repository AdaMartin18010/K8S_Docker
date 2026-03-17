# 可观测性体系

> 指标、日志、追踪三大支柱

---

## 可观测性三大支柱

```
┌─────────────────────────────────────────────────────────────┐
│                      可观测性体系                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐   │
│  │    Metrics    │  │     Logs      │  │    Traces     │   │
│  │    (指标)      │  │    (日志)      │  │    (追踪)      │   │
│  │               │  │               │  │               │   │
│  │ - Prometheus  │  │ - Fluent Bit  │  │ - Jaeger      │   │
│  │ - Grafana     │  │ - Loki        │  │ - Zipkin      │   │
│  │ - Datadog     │  │ - ELK Stack   │  │ - Tempo       │   │
│  └───────────────┘  └───────────────┘  └───────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 指标监控 (Metrics)

### Prometheus + Grafana

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: web-app
spec:
  selector:
    matchLabels:
      app: web-app
  endpoints:
    - port: metrics
      interval: 15s
```

### 关键指标

| 类别 | 指标 | 说明 |
|------|------|------|
| 资源 | CPU/Memory/Disk | 系统资源使用 |
| 应用 | QPS/Latency/Error Rate | 应用性能 |
| 业务 | Order Count/User Active | 业务指标 |

---

## 日志收集 (Logs)

### Fluent Bit 配置

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluent-bit-config
data:
  fluent-bit.conf: |
    [INPUT]
        Name tail
        Path /var/log/containers/*.log
        Parser docker

    [FILTER]
        Name kubernetes
        Match kube.*

    [OUTPUT]
        Name es
        Match kube.*
        Host elasticsearch
        Port 9200
```

---

## 分布式追踪 (Traces)

### OpenTelemetry

```yaml
apiVersion: v1
kind: Deployment
metadata:
  name: otel-collector
spec:
  template:
    spec:
      containers:
        - name: otel-collector
          image: otel/opentelemetry-collector:latest
          args: ["--config=/conf/collector-config.yaml"]
```

---

## 告警配置

```yaml
groups:
  - name: web-app
    rules:
      - alert: HighErrorRate
        expr: |
          sum(rate(http_requests_total{status=~"5.."}[5m]))
          / sum(rate(http_requests_total[5m])) > 0.01
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
```
