# Kubernetes 可观测性

本目录包含监控、日志和追踪的配置示例。

## 三大支柱

| 支柱 | 工具 | 用途 |
|------|------|------|
| **指标** | Prometheus | 数值数据、趋势分析 |
| **日志** | Fluent Bit + Elasticsearch | 文本数据、故障排查 |
| **追踪** | Jaeger + OpenTelemetry | 请求链路、性能分析 |

## 文件说明

| 文件 | 描述 |
|------|------|
| `prometheus-monitoring.yaml` | Prometheus 监控配置 |
| `logging-stack.yaml` | ELK 日志栈配置 |

## 快速开始

### 部署监控

```bash
# 部署 Prometheus + Grafana
kubectl apply -f prometheus-monitoring.yaml

# 查看 ServiceMonitor
kubectl get servicemonitor -n monitoring

# 查看告警规则
kubectl get prometheusrules -n monitoring
```

### 部署日志栈

```bash
# 创建命名空间
kubectl create namespace logging

# 部署 Fluent Bit + Elasticsearch + Kibana
kubectl apply -f logging-stack.yaml

# 查看日志收集器
kubectl get daemonset -n logging

# 访问 Kibana
kubectl port-forward svc/kibana 5601:5601 -n logging
# 访问 http://localhost:5601
```

## 应用集成

### 暴露指标

```go
// 在应用中暴露 Prometheus 指标
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // 注册指标端点
    http.Handle("/metrics", promhttp.Handler())
    
    // 业务端点
    http.HandleFunc("/api", handler)
    
    http.ListenAndServe(":8080", nil)
}
```

### 结构化日志

```go
// 使用结构化日志
import (
    "go.uber.org/zap"
)

logger, _ := zap.NewProduction()
defer logger.Sync()

logger.Info("request processed",
    zap.String("method", "GET"),
    zap.String("path", "/api/users"),
    zap.Int("status", 200),
    zap.Duration("latency", time.Millisecond*50),
)
```

## 最佳实践

### 监控

1. **四大黄金指标**:
   - Latency（延迟）
   - Traffic（流量）
   - Errors（错误）
   - Saturation（饱和度）

2. **USE 方法**（资源）:
   - Utilization（使用率）
   - Saturation（饱和度）
   - Errors（错误）

3. **RED 方法**（服务）:
   - Rate（速率）
   - Errors（错误）
   - Duration（持续时间）

### 日志

1. **结构化日志**: 使用 JSON 格式
2. **日志级别**: DEBUG, INFO, WARN, ERROR
3. **上下文信息**: 包含时间、服务、追踪 ID
4. **敏感信息**: 不要记录密码、密钥

### 告警

1. **分层告警**:
   - P0: 立即处理（页面）
   - P1: 24 小时内处理
   - P2: 工作时间内处理

2. **告警原则**:
   - 可操作的
   - 有运行手册
   - 避免告警疲劳

## 查询示例

### Prometheus

```promql
# 查询 QPS
sum(rate(http_requests_total[5m])) by (status)

# 查询错误率
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))

# 查询 P95 延迟
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))
```

### Elasticsearch

```json
// 查询特定服务的错误日志
{
  "query": {
    "bool": {
      "must": [
        { "match": { "kubernetes.container.name": "web-app" }},
        { "match": { "level": "ERROR" }}
      ]
    }
  }
}
```
