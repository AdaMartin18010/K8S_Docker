# Prometheus 3.0 - 云原生监控新时代

## 概述

Prometheus 3.0 是 Prometheus 项目自 2.0（2017年）以来最重要的主版本发布。它带来了全新的 UI、Remote Write 2.0、原生 OTLP 支持、UTF-8 指标名支持等重大改进，进一步巩固了其作为云原生监控标准地位。

> **关键数据**: Prometheus 3.0 Remote Write 2.0 减少了 60% 的网络流量，90% 的内存分配和 70% 的 CPU 使用。

## 主要新特性

### 1. 全新 UI（默认启用）

```
新 UI 特性：
├── PromLens 风格的查询树视图
├── 现代化的界面设计
├── Explain 标签页解释 PromQL 查询
├── 改进的指标浏览器
└── UTF-8 指标名支持
```

```bash
# 临时切换回旧 UI
prometheus --enable-feature=old-ui
```

### 2. Remote Write 2.0

```yaml
# prometheus.yml
remote_write:
  - url: "https://mimir.example.com/api/v1/push"
    remote_write_relabel_configs:
      - source_labels: [__name__]
        regex: 'go_.*'
        action: drop
    # Remote Write 2.0 自动启用
    # 特性：
    # - 字符串驻留减少带宽 60%
    # - 支持元数据、exemplars、created timestamps
    # - 更好的部分写入处理
```

**性能提升对比**:

| 指标 | Remote Write 1.0 | Remote Write 2.0 |
|------|-----------------|-----------------|
| 网络流量 | 基准 | -60% |
| 内存分配 | 基准 | -90% |
| CPU 使用 | 基准 | -70% |

### 3. OTLP 原生接收

```yaml
# prometheus.yml - 直接接收 OTLP 指标
otlp:
  # 启用 OTLP 接收端点
  promote_resource_attributes:
    - service.name
    - service.namespace
    - service.instance.id
    - deployment.environment

  # 转换策略（处理 OTEL 命名规范）
  translation_strategy: NoUTF8EscapingWithSuffixes

# 或者使用命令行标志
# --enable-feature=otlp-write-receiver
```

```bash
# OTLP 端点
# gRPC: localhost:4317
# HTTP: localhost:4318/v1/metrics
```

### 4. UTF-8 指标名支持

```yaml
# prometheus.yml
# 默认启用 UTF-8 支持
# 允许在指标名和标签名中使用 UTF-8 字符

# 查询时需要引号
curl 'http://localhost:9090/api/v1/query?query={"my.metric.with.dots"}'

# PromQL 中使用反引号
{"http.request.duration", status="200"}
```

### 5. Native Histograms（实验性）

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'app'
    static_configs:
      - targets: ['localhost:8080']
    # 启用 Native Histograms 抓取
    enable_native_histograms: true

# 或者全局启用
# --enable-feature=native-histograms
```

**Native Histograms 优势**:

- 指数增长的桶边界，无需预定义
- 更高分辨率的数据
- 减少配置复杂性

## 升级指南

### 从 2.x 升级到 3.0

```bash
# 1. 检查破坏性变更
# 以下特性标志已移除（现在默认启用）：
# - --enable-feature=promql-at-modifier
# - --enable-feature=promql-negative-offset
# - --enable-feature=remote-write-receiver
# - --enable-feature=no-scrape-default-port
# - --enable-feature=new-service-discovery-manager

# 2. 数据兼容性
# Prometheus 3.0 可以读取 2.x 的数据文件
# 但 2.x 无法读取 3.0 的数据文件

# 3. 滚动升级步骤
# - 保留旧的 Prometheus 实例运行
# - 启动新的 3.0 实例
# - 逐步切换查询流量
# - 确认稳定后关闭旧实例
```

## OpenTelemetry 集成最佳实践

```yaml
# prometheus.yml - 完整的 OTLP 配置
otlp:
  promote_resource_attributes:
    - service.name
    - service.namespace
    - service.instance.id
    - deployment.environment
    - k8s.cluster.name
    - k8s.namespace.name
    - k8s.pod.name
  translation_strategy: NoUTF8EscapingWithSuffixes

# 存储配置
storage:
  tsdb:
    # 保留策略
    retention.time: 30d
    retention.size: 100GB
    # OTLP 数据压缩
    out_of_order_time_window: 5m

scrape_configs:
  # 抓取 OpenTelemetry Collector
  - job_name: 'otel-collector'
    static_configs:
      - targets: ['otel-collector:8889']
    # 收集 Collector 自身的指标
    metrics_path: /metrics
```

### OpenTelemetry Collector 配置

```yaml
# otel-collector-config.yaml
exporters:
  prometheusremotewrite:
    endpoint: http://prometheus:9090/api/v1/write
    # 或者直接使用 OTLP
  otlp:
    endpoint: prometheus:4317
    tls:
      insecure: true

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
```

## Prometheus Agent 模式

```bash
# 轻量级仅抓取模式（无查询、无本地存储）
prometheus --enable-feature=agent \
  --config.file=prometheus.yml \
  --remote-write-receiver

# 适用于：
# - 边缘集群
# - 大量边缘节点
# - 远程写入集中存储
```

## 与 Thanos/Cortex/Mimir 集成

```yaml
# prometheus.yml
remote_write:
  - url: "https://thanos-receive.example.com/api/v1/receive"
    queue_config:
      capacity: 10000
      max_samples_per_send: 2000
      max_shards: 200
    metadata_config:
      send: true
      max_samples_per_send: 500
    # 启用 exemplars 支持
    send_exemplars: true
    # 启用 native histograms
    send_native_histograms: true
```

## 性能优化

### WAL 压缩

```yaml
# prometheus.yml
storage:
  tsdb:
    # WAL 压缩（减少磁盘使用）
    wal_compression: true
    # 内存中保留的块
    retention.size: 50GB
```

### 样本限制

```yaml
scrape_configs:
  - job_name: 'kubernetes-pods'
    sample_limit: 10000
    label_limit: 50
    label_name_length_limit: 1024
    label_value_length_limit: 2048
```

## 监控 Prometheus 自身

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
    # 抓取自身的 metrics
```

### 关键告警

```yaml
# prometheus-alerts.yml
groups:
  - name: prometheus
    rules:
      - alert: PrometheusTargetMissing
        expr: up == 0
        for: 5m
        labels:
          severity: critical

      - alert: PrometheusHighMemoryUsage
        expr: |
          process_resident_memory_bytes /
          (container_spec_memory_limit_bytes > 0 or 1e18) > 0.8
        for: 10m
        labels:
          severity: warning

      - alert: PrometheusTSDBHighCompactionLoad
        expr: prometheus_tsdb_compactions_triggered_total > 100
        for: 30m
        labels:
          severity: warning
```

## 总结

| 特性 | Prometheus 2.x | Prometheus 3.0 |
|------|---------------|----------------|
| UI | 旧版 | 全新 React 架构 |
| Remote Write | 1.0 | 2.0 (60% 带宽节省) |
| OTLP | 需要 Collector | 原生支持 |
| UTF-8 | 不支持 | 默认启用 |
| Native Histograms | 实验性 | 更完善 |
| Agent 模式 | 特性标志 | 稳定版 |

Prometheus 3.0 是云原生监控的里程碑版本，特别是在 OpenTelemetry 生态集成方面有重大突破。
