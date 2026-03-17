# 可观测性工具

> 监控、日志、追踪完整解决方案

---

## 三大支柱

| 支柱 | 工具 | 数据类型 | 用途 |
|------|------|----------|------|
| **Metrics** | Prometheus | 时序数据 | 告警、趋势分析 |
| **Logs** | Loki/Fluentd | 文本日志 | 故障排查、审计 |
| **Traces** | Jaeger/Tempo | 请求链路 | 性能分析、依赖梳理 |

---

## OpenTelemetry 统一标准

```
┌─────────┐     ┌─────────────┐     ┌─────────────┐
│ 应用    │────▶│  Collector  │────▶│  后端存储   │
│ (SDK)   │     │  (OTel)     │     │ (Prometheus/
└─────────┘     └─────────────┘     │  Jaeger)    │
                                     └─────────────┘
```

---

## 工具链

| 功能 | 工具 | 说明 |
|------|------|------|
| 指标收集 | [Prometheus 3.0](../prometheus-3/) | 时序数据库 |
| 日志收集 | [Grafana Alloy](../grafana-alloy/) | OpenTelemetry 收集器 |
| 分布式追踪 | [Tempo](../tempo/) | 轻量级追踪后端 |
| 可视化 | Grafana | 统一仪表板 |

---

## 本章内容

- [OpenTelemetry](./opentelemetry.md)

---

## 最佳实践

1. **统一埋点**：使用 OpenTelemetry SDK
2. **关联数据**：TraceID 串联 Metrics/Logs/Traces
3. **合理采样**：生产环境配置自适应采样
4. **成本控制**：设置数据保留策略

---

## 相关文档

- [K8s 可观测性](../../03-kubernetes/05-observability/)
