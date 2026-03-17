# NATS - 云原生消息系统

## 概述

NATS (Neural Autonomic Transport System) 是一个高性能、轻量级、开源的消息系统，专为构建现代分布式系统设计。作为 CNCF 项目，NATS 支持从边缘设备到大规模云部署的各种场景。

## 核心特性

### 轻量级架构

- **单二进制**: 单个二进制文件，无外部依赖
- **低资源**: 典型部署 <20MB 内存
- **高性能**: 每秒数百万消息，亚毫秒延迟
- **多语言**: 40+ 客户端语言支持

### 统一消息模式

| 模式 | 描述 | 用例 |
|------|------|------|
| Pub/Sub | 发布订阅 | 事件通知、广播 |
| Queue Groups | 队列组 | 负载均衡、工作池 |
| Request/Reply | 请求响应 | RPC、同步调用 |
| JetStream | 流式处理 | 持久化、重放 |
| Key-Value | 键值存储 | 配置管理、状态 |
| Object Store | 对象存储 | 文件存储 |

## JetStream 流处理

### 核心概念

```
Stream: 持久化、时间排序的消息列表
  └── Subject: 消息主题（如 orders.created）
  └── Consumer: 消费游标，按策略消费消息
      ├── Push Consumer: 服务器推送消息
      └── Pull Consumer: 客户端主动拉取
```

### 流配置

```bash
# 创建流
nats stream add ORDERS \
  --subjects="orders.*" \
  --storage=memory \
  --replicas=3 \
  --retention=limits \
  --discard=old \
  --max-msgs=1000000 \
  --max-bytes=1GB \
  --max-age=7d

# 创建消费者
nats consumer add ORDERS ORDER_PROCESSOR \
  --filter="orders.created" \
  --deliver=all \
  --ack=explicit \
  --max-deliver=3 \
  --replay=instant
```

### Go SDK 示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/nats-io/nats.go"
)

func main() {
    // 连接 NATS
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        panic(err)
    }
    defer nc.Close()

    // JetStream 上下文
    js, err := nc.JetStream()
    if err != nil {
        panic(err)
    }

    // 创建流
    _, err = js.AddStream(&nats.StreamConfig{
        Name:     "ORDERS",
        Subjects: []string{"orders.*"},
        Storage:  nats.MemoryStorage,
        Replicas: 3,
    })
    if err != nil {
        panic(err)
    }

    // 发布消息
    _, err = js.Publish("orders.created", []byte(`{"id": "123", "total": 99.99}`))
    if err != nil {
        panic(err)
    }

    // 订阅（Pull 模式）
    sub, err := js.PullSubscribe("orders.*", "PROCESSOR")
    if err != nil {
        panic(err)
    }

    // 拉取消息
    msgs, err := sub.Fetch(10)
    if err != nil {
        panic(err)
    }

    for _, msg := range msgs {
        fmt.Printf("Received: %s\n", string(msg.Data))
        msg.Ack()
    }
}
```

## Kubernetes 部署

### Helm 安装

```bash
# 添加仓库
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update

# 安装 JetStream 集群
helm install nats nats/nats \
  --namespace nats --create-namespace \
  --set config.jetstream.enabled=true \
  --set config.jetstream.memoryStore.enabled=true \
  --set config.jetstream.memoryStore.maxSize=1Gi \
  --set config.jetstream.fileStore.enabled=true \
  --set config.jetstream.fileStore.size=10Gi \
  --set config.jetstream.fileStore.storageClassName=standard \
  --set config.cluster.enabled=true \
  --set config.cluster.replicas=3 \
  --set podTemplate.resources.requests.memory=512Mi \
  --set podTemplate.resources.limits.memory=1Gi

# 验证
kubectl get pods -n nats
kubectl exec -n nats nats-0 -- nats server list
```

### 高级配置

```yaml
# values-production.yaml
config:
  jetstream:
    enabled: true
    domain: "hub"
    encrypt: true  # 启用加密
    uniqueTag: "az:aws:us-east-1a"

  leafnodes:
    enabled: true
    port: 7422
    remotes:
    - url: tls://nats-leaf.example.com:7422
      credentials:
        secret:
          name: leaf-credentials
          key: creds

  gateway:
    enabled: true
    name: "us-east-1"
    rejectUnknownCluster: true
    urls:
    - nats://nats-0.nats.nats.svc:7222
    - nats://nats-1.nats.nats.svc:7222
    - nats://nats-2.nats.nats.svc:7222

  websocket:
    enabled: true
    port: 8080
    noTLS: false

podTemplate:
  resources:
    requests:
      cpu: "1"
      memory: "2Gi"
    limits:
      cpu: "2"
      memory: "4Gi"
```

## 边缘与联邦部署

### Leaf Nodes（边缘节点）

```yaml
# 边缘集群 NATS 配置
leafnodes {
    remotes = [
        {
            url: "tls://connect.ngs.global:7422"
            credentials: "/etc/nats/leaf.creds"

            # 只同步特定主题
            subject_filters: ["orders.us-east.>"]

            # 压缩传输
            compression: {
                mode: s2_better
            }
        }
    ]
}

jetstream {
    store_dir: "/data/jetstream"
    max_memory_store: 256MB
    max_file_store: 1GB

    # 边缘域
    domain: "edge-facility-1"
}
```

### Supercluster（超集群）

```hcl
# 集群 A: us-east-1
gateway {
    name: "us-east-1"
    listen: "0.0.0.0:7222"

    gateways: [
        {name: "us-west-2", urls: ["nats://gw-us-west-2:7222"]},
        {name: "eu-west-1", urls: ["nats://gw-eu-west-1:7222"]}
    ]
}

# 集群 B: us-west-2
gateway {
    name: "us-west-2"
    listen: "0.0.0.0:7222"

    gateways: [
        {name: "us-east-1", urls: ["nats://gw-us-east-1:7222"]},
        {name: "eu-west-1", urls: ["nats://gw-eu-west-1:7222"]}
    ]
}
```

## 安全与认证

### NKey 认证

```bash
# 生成 NKey
nats-server -gen-nkeys

# 生成用户凭证
nsc add account PRODUCTION
nsc add user --account PRODUCTION service-user
nsc generate creds --account PRODUCTION service-user > service.creds
```

### JWT 令牌

```go
// 使用 JWT 连接
nc, err := nats.Connect(
    nats.DefaultURL,
    nats.UserJWTAndSeed(
        func() (string, error) { return loadJWT() },
        func(string) ([]byte, error) { return loadSeed() },
    ),
)
```

### TLS 配置

```hcl
tls {
    cert_file: "/etc/nats/server.pem"
    key_file: "/etc/nats/server-key.pem"
    ca_file: "/etc/nats/ca.pem"

    # 双向 TLS
    verify: true
    verify_and_map: true
}

# JetStream 加密
jetstream {
    store_dir: "/data/jetstream"
    cipher: "AES-GCM"
    key: "${JETSTREAM_KEY}"
}
```

## 与 wasmCloud 集成

```yaml
# wasmcloud-nats.yaml
apiVersion: k8s.wasmcloud.dev/v1alpha1
kind: WasmCloudHostConfig
metadata:
  name: wasmcloud-host
spec:
  natsAddress: nats://nats-headless.nats.svc.cluster.local:4222
  natsLeafImage: nats:2.10
  lattice: production

  # 使用 JetStream 作为后端
  jsDomain: "wasmcloud"
```

```rust
// Rust 中使用 NATS 能力
use wasmcloud::messaging::*;

struct Component;

impl Guest for Component {
    fn handle_message(message: Message) -> Result<(), String> {
        // 处理来自 NATS 的消息
        let subject = message.subject;
        let body = String::from_utf8_lossy(&message.body);

        println!("Received on {}: {}", subject, body);

        // 回复消息
        publish(&PublishRequest {
            subject: format!("{}.reply", subject),
            reply_to: None,
            body: b"Processed".to_vec(),
        })
    }
}
```

## 监控与可观测性

### Prometheus 指标

```yaml
# nats-metrics.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nats-metrics
data:
  prometheus.yaml: |
    global:
      scrape_interval: 15s
    scrape_configs:
    - job_name: 'nats'
      static_configs:
      - targets: ['nats:7777']
      metrics_path: /metrics
```

### OpenTelemetry 追踪

```go
// 启用消息追踪
js.Publish("orders.created", data,
    nats.MsgId("unique-id"),
    nats.ExpectLastSubjectSequence(seq),
    nats.OpenTelemetryContext(ctx), // OTel 传播
)
```

## 2025 新特性

- **MQTT 支持**: 原生 MQTT 协议支持，边缘 IoT 集成
- **WebSocket 优化**: 改进的浏览器/Web 客户端支持
- **消息压缩**: S2 压缩算法支持
- **分层存储**: 热/温/冷数据自动分层
- **Exactly-Once**: 更强的一致性保证

## 性能基准

| 场景 | 吞吐量 | 延迟 |
|------|--------|------|
| Core Pub/Sub | 2M+ msg/s | <100μs |
| JetStream (内存) | 400K msg/s | <1ms |
| JetStream (磁盘) | 100K msg/s | 1-5ms |
| 跨集群复制 | 50K msg/s | <50ms |

## 相关资源

- [NATS 官网](https://nats.io/)
- [NATS 文档](https://docs.nats.io/)
- [JetStream 文档](https://docs.nats.io/nats-concepts/jetstream)
- [NATS Go 客户端](https://github.com/nats-io/nats.go)
