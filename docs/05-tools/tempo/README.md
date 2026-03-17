# Grafana Tempo - 分布式追踪后端

## 概述

Grafana Tempo 是一个开源、易于使用且大规模的分布式追踪后端，与 Grafana 深度集成，支持 TraceQL 查询语言，能够高效存储和查询追踪数据。

## 核心特性

| 特性 | 描述 |
|------|------|
| 对象存储 | 支持 S3、GCS、Azure Blob、本地存储 |
| TraceQL | 专为追踪设计的查询语言 |
| 与指标关联 | 通过 Exemplars 连接指标和追踪 |
| 自动扩缩 | 支持水平扩展 |
| 多租户 | 内置租户隔离 |

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                     Tempo 架构                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐              │
│  │ Distributor│   │  Ingester │   │  Querier  │              │
│  │  (分发)   │───▶│  (写入)   │───▶│  (查询)   │              │
│  └──────────┘    └────┬─────┘    └────┬─────┘              │
│                       │               │                     │
│                       ▼               ▼                     │
│              ┌─────────────────────────────┐               │
│              │       Object Storage        │               │
│              │      (S3/GCS/Azure)         │               │
│              └─────────────────────────────┘               │
│                                                              │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐              │
│  │ Compactor│    │  Query   │    │  Gateway │              │
│  │ (压缩)   │    │ Frontend │    │ (网关)   │              │
│  └──────────┘    └──────────┘    └──────────┘              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## 安装部署

### Helm 安装

```bash
# 添加仓库
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# 安装 Tempo
helm install tempo grafana/tempo-distributed \
  --namespace monitoring \
  --create-namespace \
  --set global.storage.backend=s3 \
  --set global.storage.s3.bucket=tempo-traces \
  --set global.storage.s3.endpoint=s3.amazonaws.com \
  --set global.storage.s3.region=us-west-2
```

### 单节点模式（测试）

```yaml
# docker-compose.yaml
version: "3"
services:
  tempo:
    image: grafana/tempo:latest
    command: ["-config.file=/etc/tempo.yaml"]
    volumes:
    - ./tempo.yaml:/etc/tempo.yaml
    ports:
    - "3200:3200"  # Tempo HTTP
    - "4317:4317"  # OTLP gRPC
    - "4318:4318"  # OTLP HTTP
```

```yaml
# tempo.yaml
server:
  http_listen_port: 3200

distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: 0.0.0.0:4317
        http:
          endpoint: 0.0.0.0:4318

ingester:
  trace_idle_period: 10s
  max_block_bytes: 1048576
  max_block_duration: 5m

compactor:
  compaction:
    compaction_window: 1h
    max_block_bytes: 100_000_000
    block_retention: 168h  # 7天
    compacted_block_retention: 1h

storage:
  trace:
    backend: local
    local:
      path: /tmp/tempo
    wal:
      path: /tmp/tempo/wal

overrides:
  defaults:
    global:
      max_bytes_per_trace: 5000000
```

## TraceQL 查询

### 基本查询

```traceql
// 查询特定服务
{resource.service.name="checkout-service"}

// 查询错误追踪
{status=error}

// 查询特定操作
{name="GET /api/orders"}
```

### 高级查询

```traceql
// 延迟大于 1秒的 POST 请求
{name="POST /api/checkout" && duration > 1s}

// 特定 HTTP 状态码
{span.http.status_code=500}

// 多个条件组合
{resource.service.name="payment-service" && status=error && duration > 500ms}

// 范围查询（查找包含数据库查询的追踪）
{span.db.system="postgresql"} | select(status, span.db.statement)

// 聚合查询
{resource.service.name="api-gateway"} | count() by (span.http.route)
```

### 结构查询

```traceql
// 父子关系查询
{resource.service.name="frontend"} >> {resource.service.name="backend"}

// 后代查询
{resource.service.name="gateway"} > {span.db.system="mysql"}
```

## 与 Grafana 集成

```yaml
# Grafana 数据源配置
datasources:
- name: Tempo
  type: tempo
  url: http://tempo:3200
  access: proxy
  jsonData:
    httpMethod: GET
    nodeGraph:
      enabled: true
    search:
      hide: false
    traceQuery:
      timeShiftEnabled: true
      spanStartTimeShift: 1h
      spanEndTimeShift: 1h
    spanBar:
      type: Tag
      tag: http.path
    # 关联其他数据源
    tracesToLogs:
      datasourceUid: loki
      tags: ['pod', 'namespace']
      spanStartTimeShift: '-1h'
      spanEndTimeShift: '1h'
    tracesToMetrics:
      datasourceUid: prometheus
      spanStartTimeShift: '-1h'
      spanEndTimeShift: '1h'
```

## 2025 新特性

- **Tempo 2.9**: MCP Server 支持，AI 辅助调试
- **TraceQL Metrics**: 概率查询提示加速分析
- **eBPF Profiling**: 追踪到性能分析的关联
- **Streaming**: 流式查询结果

## 性能优化

```yaml
# 生产环境配置
ingester:
  flush_check_period: 10s
  flush_op_timeout: 5m
  trace_idle_period: 30s
  max_block_bytes: 1_000_000_000  # 1GB
  max_block_duration: 1h

compactor:
  compaction:
    chunk_size_bytes: 5_000_000
    flush_size_bytes: 30_000_000
    max_compaction_objects: 6000000
    max_block_bytes: 100_000_000_000  # 100GB

overrides:
  defaults:
    ingestion:
      rate_limit_bytes: 15000000
      burst_size_bytes: 20000000
    global:
      max_bytes_per_trace: 30000000
```

## 相关资源

- [Grafana Tempo 文档](https://grafana.com/docs/tempo/)
- [TraceQL 参考](https://grafana.com/docs/tempo/traceql/)
- [GitHub](https://github.com/grafana/tempo)
