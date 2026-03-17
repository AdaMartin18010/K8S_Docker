# Prometheus 监控实战

> Kubernetes 生产级监控方案

---

## 监控体系架构

```
┌─────────────────────────────────────────────────────────────┐
│              Kubernetes 监控架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  数据收集层                                                   │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ Prometheus  │ │  Node       │ │  kube-state │           │
│  │ Server      │ │  Exporter   │ │  metrics    │           │
│  │             │ │             │ │             │           │
│  └──────┬──────┘ └─────────────┘ └─────────────┘           │
│         │                                                    │
│  ┌──────┴───────────────────────────────────────────────┐   │
│  │              应用指标 (Application Metrics)            │   │
│  │  /metrics ──► Prometheus ◄── ServiceMonitor          │   │
│  └───────────────────────────────────────────────────────┘   │
│                                                              │
│  存储查询层                                                   │
│  ┌─────────────┐ ┌─────────────┐                           │
│  │ Prometheus  │ │  Thanos     │  (长期存储/高可用)         │
│  │ TSDB        │ │  Sidecar    │                           │
│  └──────┬──────┘ └──────┬──────┘                           │
│         │               │                                    │
│         └───────────────┴────────────────┐                   │
│                                          │                   │
│  展示告警层                               │                   │
│  ┌─────────────┐ ┌─────────────┐        │                   │
│  │  Grafana    │ │ Alertmanager│◄───────┘                   │
│  │  Dashboards │ │             │                            │
│  └─────────────┘ └─────────────┘                            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心指标类型

### 四金指标

| 指标 | 说明 | PromQL 示例 |
|------|------|-------------|
| **延迟 (Latency)** | 请求处理时间 | `histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))` |
| **流量 (Traffic)** | 请求数量 | `rate(http_requests_total[5m])` |
| **错误 (Errors)** | 错误率 | `rate(http_requests_total{status=~"5.."}[5m])` |
| **饱和度 (Saturation)** | 资源使用 | `1 - rate(node_cpu_seconds_total{mode="idle"}[5m])` |

---

## ServiceMonitor 配置

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: app-metrics
  namespace: monitoring
spec:
  namespaceSelector:
    matchNames:
      - production
  selector:
    matchLabels:
      app: myapp
      metrics: enabled
  endpoints:
    - port: metrics
      path: /metrics
      interval: 30s
      scrapeTimeout: 10s
      scheme: https
      tlsConfig:
        insecureSkipVerify: true
      bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
      metricRelabelings:
        - sourceLabels: [__name__]
          regex: 'go_.*'
          action: drop
```

---

## 关键告警规则

### Pod 级别

```yaml
groups:
  - name: pod-alerts
    rules:
      - alert: PodCrashLooping
        expr: |
          rate(kube_pod_container_status_restarts_total[10m]) > 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Pod {{ $labels.pod }} is crash looping"

      - alert: PodHighMemoryUsage
        expr: |
          (
            container_memory_working_set_bytes{container!=""}
            /
            container_spec_memory_limit_bytes{container!=""}
          ) > 0.85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Pod {{ $labels.pod }} memory usage > 85%"

      - alert: PodNotReady
        expr: |
          kube_pod_status_ready{condition="false"} == 1
        for: 15m
        labels:
          severity: warning
        annotations:
          summary: "Pod {{ $labels.pod }} not ready"
```

### 节点级别

```yaml
groups:
  - name: node-alerts
    rules:
      - alert: NodeHighCPU
        expr: |
          100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Node {{ $labels.instance }} CPU > 80%"

      - alert: NodeDiskPressure
        expr: |
          (
            node_filesystem_avail_bytes{mountpoint="/"}
            /
            node_filesystem_size_bytes{mountpoint="/"}
          ) < 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Node {{ $labels.instance }} disk < 10%"
```

---

## Grafana Dashboard

### Pod 仪表板 JSON 模型

```json
{
  "dashboard": {
    "title": "Kubernetes Pod Overview",
    "panels": [
      {
        "title": "CPU Usage",
        "type": "timeseries",
        "targets": [{
          "expr": "rate(container_cpu_usage_seconds_total{pod=~\"$pod\"}[5m])",
          "legendFormat": "{{ container }}"
        }]
      },
      {
        "title": "Memory Usage",
        "type": "timeseries",
        "targets": [{
          "expr": "container_memory_working_set_bytes{pod=~\"$pod\"}",
          "legendFormat": "{{ container }}"
        }]
      },
      {
        "title": "Network I/O",
        "type": "timeseries",
        "targets": [{
          "expr": "rate(container_network_receive_bytes_total{pod=~\"$pod\"}[5m])",
          "legendFormat": "RX"
        },{
          "expr": "rate(container_network_transmit_bytes_total{pod=~\"$pod\"}[5m])",
          "legendFormat": "TX"
        }]
      }
    ]
  }
}
```

---

## 自定义应用指标 (Go)

```go
package main

import (
 "net/http"
 "time"

 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promauto"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
 // 计数器
 requestsTotal = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "http_requests_total",
   Help: "Total HTTP requests",
  },
  []string{"method", "path", "status"},
 )

 // 直方图
 requestDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_request_duration_seconds",
   Help:    "HTTP request duration",
   Buckets: prometheus.DefBuckets,
  },
  []string{"method", "path"},
 )

 // 仪表盘
 activeConnections = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "active_connections",
   Help: "Number of active connections",
  },
 )
)

func main() {
 http.Handle("/metrics", promhttp.Handler())
 http.HandleFunc("/api/", handler)
 http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
 start := time.Now()
 activeConnections.Inc()
 defer activeConnections.Dec()

 // 处理请求
 status := "200"

 // 记录指标
 requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
 requestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
}
```

---

## 容量规划指标

| 指标 | 用途 | 建议阈值 |
|------|------|----------|
| 节点 CPU 使用率 | 扩容决策 | > 70% |
| 节点内存使用率 | 扩容决策 | > 80% |
| Pod 副本数 / HPA | 自动扩缩容 | 基于 QPS |
| ETCD 数据库大小 | 维护窗口 | > 2GB |
| API Server 延迟 | 性能评估 | P99 < 1s |
