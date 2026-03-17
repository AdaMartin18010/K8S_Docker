# Grafana Alloy - OpenTelemetry 收集器

## 概述

Grafana Alloy 是 Grafana Labs 提供的 OpenTelemetry Collector 发行版，统一收集指标、日志、追踪和性能分析数据，支持 Prometheus 和 OpenTelemetry 双管道。

## 核心特性

| 特性 | 描述 |
|------|------|
| 统一收集 | 指标、日志、追踪、性能分析一体化 |
| 双管道支持 | 同时支持 Prometheus 和 OTel 协议 |
| GitOps 兼容 | 支持从 Git、S3、HTTP 加载配置 |
| 集群模式 | 支持水平扩展和高可用 |
| 安全集成 | 支持 Vault、K8s Secret 管理 |

## 架构组件

```
┌─────────────────────────────────────────────────────────────────┐
│                      Grafana Alloy 架构                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐           │
│  │  Receiver   │   │  Processor  │   │   Exporter  │           │
│  │   (接收)     │ → │   (处理)     │ → │   (导出)     │           │
│  └─────────────┘   └─────────────┘   └─────────────┘           │
│                                                                 │
│  接收器: OTLP, Prometheus, Kafka, Journald, CloudWatch          │
│  处理器: Batch, Filter, Transform, Tail-based Sampling          │
│  导出器: Prometheus, Tempo, Loki, Mimir, CloudWatch             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 安装部署

### Kubernetes 安装

```bash
# Helm 安装
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# 安装 Alloy
helm install alloy grafana/alloy \
  --namespace monitoring \
  --create-namespace \
  --set controller.type=daemonset \
  --set alloy.configMap.create=true \
  --set alloy.configMap.key=config.alloy
```

### 配置示例

```alloy
// config.alloy - 完整可观测性管道

// ========== 接收器 ==========
// OTLP 接收器
otelcol.receiver.otlp "default" {
  grpc {
    endpoint = "0.0.0.0:4317"
  }
  http {
    endpoint = "0.0.0.0:4318"
  }

  output {
    metrics = [otelcol.processor.batch.default.input]
    logs    = [otelcol.processor.batch.default.input]
    traces  = [otelcol.processor.batch.default.input]
  }
}

// Prometheus 抓取
prometheus.scrape "k8s" {
  targets = prometheus.kubernetes.targets {
    role = "pod"
  }
  forward_to = [prometheus.remote_write.mimir.receiver]
}

// ========== 处理器 ==========
otelcol.processor.batch "default" {
  timeout = "1s"
  send_batch_size = 1024

  output {
    metrics = [otelcol.exporter.prometheus.mimir.input]
    logs    = [otelcol.exporter.loki.default.input]
    traces  = [otelcol.processor.tail_sampling.default.input]
  }
}

// 尾部采样（Tail-based Sampling）
otelcol.processor.tail_sampling "default" {
  decision_wait = "10s"
  num_traces = 100000
  expected_new_traces_per_sec = 1000

  policy {
    name = "errors"
    type = "status_code"
    status_code { status_codes = [ERROR] }
  }

  policy {
    name = "latency"
    type = "latency"
    latency { threshold_ms = 1000 }
  }

  output {
    traces = [otelcol.exporter.otlp.tempo.input]
  }
}

// ========== 导出器 ==========
otelcol.exporter.prometheus "mimir" {
  forward_to = [prometheus.remote_write.mimir.receiver]
}

prometheus.remote_write "mimir" {
  endpoint {
    url = "http://mimir.monitoring.svc:9009/api/v1/push"
    headers = { "X-Scope-OrgID" = "tenant-1" }
  }
}

otelcol.exporter.loki "default" {
  forward_to = [loki.write.default.receiver]
}

loki.write "default" {
  endpoint {
    url = "http://loki.monitoring.svc:3100/loki/api/v1/push"
  }
}

otelcol.exporter.otlp "tempo" {
  client {
    endpoint = "tempo.monitoring.svc:4317"
    tls { insecure = true }
  }
}
```

## 集群模式配置

```alloy
// 集群模式配置
prometheus.operator.podmonitors "default" {
  forward_to = [prometheus.remote_write.mimir.receiver]
  clustering {
    enabled = true
  }
}

// 自动分片抓取目标
prometheus.scrape "k8s" {
  targets = prometheus.kubernetes.targets {
    role = "pod"
  }
  clustering {
    enabled = true
  }
  forward_to = [prometheus.remote_write.mimir.receiver]
}
```

## 与 Grafana LGTM 集成

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Alloy     │────▶│   Mimir     │────▶│  Grafana    │
│  (收集器)    │     │  (指标)     │     │  (可视化)    │
└──────┬──────┘     └─────────────┘     └─────────────┘
       │
       ├──────────▶│   Tempo     │────▶│  Grafana    │
       │           │  (追踪)     │     │  (可视化)    │
       │           └─────────────┘     └─────────────┘
       │
       └──────────▶│   Loki      │────▶│  Grafana    │
                   │  (日志)     │     │  (可视化)    │
                   └─────────────┘     └─────────────┘
```

## 性能调优

```alloy
// 高吞吐量配置
otelcol.processor.batch "high_throughput" {
  timeout = "200ms"
  send_batch_size = 8192
  send_batch_max_size = 10240
}

// 内存限制
otelcol.processor.memory_limiter "default" {
  limit_mib = 4096
  spike_limit_mib = 512
  check_interval = "5s"
}
```

## 2025 新特性

- **集群模式 GA**: 生产级集群支持
- **Beyla eBPF**: 零代码自动仪器捐赠给 OpenTelemetry
- **v1 稳定版**: OpenTelemetry Collector v1 正式发布
- **Profiling 支持**: 持续性能分析集成

## 相关资源

- [Grafana Alloy 文档](https://grafana.com/docs/alloy/)
- [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)
- [Beyla eBPF 自动仪器](https://grafana.com/docs/beyla/)
