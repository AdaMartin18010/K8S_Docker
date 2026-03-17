# Grafana Mimir - Prometheus 长期存储

## 概述

Grafana Mimir 是一个水平可扩展、高可用、多租户的 Prometheus 兼容指标后端，支持多年的指标保留，专为大规模部署设计。

## 核心特性

| 特性 | 描述 |
|------|------|
| 水平扩展 | 支持数百万活跃时间序列 |
| 多租户 | 内置租户隔离和限制 |
| 长期存储 | 基于对象存储，支持多年保留 |
| 查询分片 | 自动查询并行化 |
| 全局聚合 | 跨集群指标聚合 |

## Mimir 3.0 新架构

```
┌─────────────────────────────────────────────────────────────────┐
│                    Mimir 3.0 解耦架构                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ┌─────────────┐         ┌─────────────┐                      │
│   │  Writers    │◀───────▶│   Kafka     │                      │
│   │  (写入)     │         │  (消息队列)  │                      │
│   └──────┬──────┘         └──────┬──────┘                      │
│          │                       │                              │
│          │         ┌─────────────▼──────┐                      │
│          │         │     Readers        │                      │
│          │         │    (异步读取)       │                      │
│          │         └──────────┬─────────┘                      │
│          │                    │                                │
│   ┌──────▼──────┐    ┌────────▼────────┐                      │
│   │   Ingester  │    │  Object Storage │                      │
│   │   (内存缓存) │    │   (S3/GCS等)     │                      │
│   └──────┬──────┘    └─────────────────┘                      │
│          │                                                      │
│   ┌──────▼──────────────────────────┐                          │
│   │        Query Engine (MQE)       │                          │
│   │      流式查询，内存减少92%       │                          │
│   └─────────────────────────────────┘                          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 安装部署

### Helm 安装

```bash
# 添加仓库
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# 安装 Mimir
helm install mimir grafana/mimir-distributed \
  --namespace monitoring \
  --create-namespace \
  --set global.storage.backend=s3 \
  --set global.storage.s3.bucket=mimir-metrics \
  --set global.storage.s3.endpoint=s3.amazonaws.com \
  --set global.storage.s3.region=us-west-2
```

### 配置示例

```yaml
# mimir.yaml
# 多租户配置
multitenancy_enabled: true

# 限制配置
limits:
  ingestion_rate: 100000
  ingestion_burst_size: 200000
  max_global_series_per_user: 2000000
  max_label_names_per_series: 30
  max_label_value_length: 1024
  max_label_name_length: 1024
  max_metadata_length: 1024
  compactor_blocks_retention_period: 1y

# Ingester 配置
ingester:
  ring:
    replication_factor: 3
    kvstore:
      store: memberlist

# 存储配置
blocks_storage:
  backend: s3
  s3:
    bucket_name: mimir-blocks
    region: us-west-2
  tsdb:
    dir: /data/tsdb
    retention_period: 24h
    ship_interval: 1m

# Compactor 配置
compactor:
  data_dir: /data/compactor
  compaction_interval: 30m
  compaction_concurrency: 4
  cleanup_interval: 15m

# 查询配置
frontend:
  parallelize_shardable_queries: true
  cache_results: true
  results_cache:
    backend: memcached
    memcached:
      addresses: mimir-memcached:11211
      max_item_size: 5MB

# 查询调度
query_scheduler:
  max_outstanding_requests_per_tenant: 800

#  ruler（告警规则）
ruler:
  enable_api: true
  rule_path: /data/rules
  ring:
    kvstore:
      store: memberlist
  alertmanager_url: http://alertmanager:9093
```

## Prometheus Remote Write

```yaml
# prometheus.yaml
remote_write:
- url: http://mimir.monitoring.svc:8080/api/v1/push
  headers:
    X-Scope-OrgID: tenant-production
  queue_config:
    capacity: 50000
    max_samples_per_send: 1000
    max_shards: 200
    min_shards: 10
    batch_send_deadline: 5s
  metadata_config:
    send: true
    max_samples_per_send: 500
```

## 下采样（Downsampling）

```yaml
# 降低长期存储成本
compactor:
  downsampling_enabled: true
  downsampling:
    - from: 0
      to: 720h      # 30天内保持原精度
      resolution: 1m
    - from: 720h
      to: 2160h     # 30-90天下采样到5分钟
      resolution: 5m
    - from: 2160h
      to: 8760h     # 90天-1年下采样到1小时
      resolution: 1h
```

## 与 Grafana 集成

```yaml
# Grafana 数据源配置
datasources:
- name: Mimir
  type: prometheus
  url: http://mimir-query-frontend:8080/prometheus
  jsonData:
    httpHeaderName1: X-Scope-OrgID
    alertmanagerUid: alertmanager
    manageAlerts: true
    prometheusType: mimir
    prometheusVersion: 2.15.0
    cacheLevel: 'High'
    incrementalQuerying: true
  secureJsonData:
    httpHeaderValue1: tenant-production
```

## 2025 新特性

- **Mimir 3.0**: 解耦读写路径，Kafka 异步摄取
- **MQE 查询引擎**: 流式查询，内存使用减少92%
- **成本降低**: 资源使用减少15%，吞吐量提升
- **OTel 集成**: 原生 OpenTelemetry 指标支持

## 性能对比

| 指标 | Mimir | Thanos | VictoriaMetrics |
|------|-------|--------|-----------------|
| 摄取性能 | 500K-2M samples/sec | 依赖 Prometheus | 1M+ samples/sec |
| 查询延迟 | 30-80ms | 50-100ms | 20-50ms |
| 压缩比 | 2-3x | 2-4x | 10x |
| 多租户 | 原生支持 | 需额外配置 | 有限支持 |

## 相关资源

- [Grafana Mimir 文档](https://grafana.com/docs/mimir/)
- [GitHub](https://github.com/grafana/mimir)
- [配置参考](https://grafana.com/docs/mimir/latest/configure/)
