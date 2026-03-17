# DragonflyDB - 现代内存数据存储

> Redis 兼容、高性能、水平可扩展的内存数据库 (2025)

---

## 概述

DragonflyDB 是一个现代内存数据存储，完全兼容 Redis 和 Memcached 协议，但提供多线程架构和更好的资源效率。

| 特性 | DragonflyDB | Redis |
|------|-------------|-------|
| 架构 | 多线程 | 单线程 |
| 垂直扩展 | 优秀 | 受限 |
| 内存效率 | 更高 | 标准 |
| 快照性能 | 无阻塞 | fork 阻塞 |
| 协议兼容 | Redis + Memcached | Redis |

---

## 架构优势

```
┌─────────────────────────────────────────────────────────────────┐
│                     DragonflyDB 架构                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                   多线程处理层                             │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐      │  │
│  │  │ Thread 1│  │ Thread 2│  │ Thread 3│  │ Thread N│      │  │
│  │  │ (查询)   │  │ (查询)   │  │ (查询)   │  │ (查询)   │      │  │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘      │  │
│  │       └─────────────┴─────────────┴─────────────┘         │  │
│  │                         │                                  │  │
│  │  ┌──────────────────────▼────────────────────────┐        │  │
│  │  │              全局数据存储层                     │        │  │
│  │  │    无锁数据结构 (Dash Table, B+ Tree)         │        │  │
│  │  └───────────────────────────────────────────────┘        │  │
│  │                         │                                  │  │
│  │  ┌──────────────────────▼────────────────────────┐        │  │
│  │  │              持久化层                          │        │  │
│  │  │    增量快照 (非阻塞 RDB + AOF)                 │        │  │
│  │  └───────────────────────────────────────────────┘        │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Kubernetes 部署

### 单实例部署

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: dragonfly
spec:
  serviceName: dragonfly
  replicas: 1
  selector:
    matchLabels:
      app: dragonfly
  template:
    metadata:
      labels:
        app: dragonfly
    spec:
      containers:
      - name: dragonfly
        image: docker.dragonflydb.io/dragonflydb/dragonfly:v1.26.0
        ports:
        - containerPort: 6379
        args:
        - --dir=/data
        - --dbfilename=dump.rdb
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "4Gi"
            cpu: "2000m"
        volumeMounts:
        - name: data
          mountPath: /data
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
---
apiVersion: v1
kind: Service
metadata:
  name: dragonfly
spec:
  selector:
    app: dragonfly
  ports:
  - port: 6379
    targetPort: 6379
```

### 集群模式（使用 Operator）

```yaml
# 安装 Dragonfly Operator
kubectl apply -f https://raw.githubusercontent.com/dragonflydb/dragonfly-operator/main/manifests/dragonfly-operator.yaml

# 创建集群
apiVersion: dragonflydb.io/v1alpha1
kind: Dragonfly
metadata:
  name: dragonfly-cluster
spec:
  replicas: 3
  resources:
    requests:
      cpu: "500m"
      memory: "1Gi"
    limits:
      cpu: "2000m"
      memory: "4Gi"
  snapshot:
    enabled: true
    cron: "0 */6 * * *"  # 每6小时
```

---

## 客户端连接

```bash
# Redis CLI 兼容
redis-cli -h dragonfly.default.svc.cluster.local -p 6379

# Python
import redis
r = redis.Redis(host='dragonfly.default.svc.cluster.local', port=6379)
r.set('key', 'value')
print(r.get('key'))

# Go
import "github.com/redis/go-redis/v9"
rdb := redis.NewClient(&redis.Options{
    Addr: "dragonfly.default.svc.cluster.local:6379",
})
```

---

## 性能对比

| 场景 | Redis | DragonflyDB | 提升 |
|------|-------|-------------|------|
| GET 操作 | 100K ops/s | 1M ops/s | 10x |
| SET 操作 | 80K ops/s | 800K ops/s | 10x |
| 内存使用 (1M keys) | 100MB | 75MB | 25% |
| 快照时间 | 阻塞秒级 | 毫秒级 | 100x |

---

## 从 Redis 迁移

```bash
# 1. 导出 Redis 数据
redis-cli --rdb backup.rdb

# 2. 导入 Dragonfly
dragonfly --dir=/data --dbfilename=backup.rdb

# 3. 零停机迁移（双写）
# 应用同时写入 Redis 和 Dragonfly
# 验证后切换读取到 Dragonfly
```

---

## 监控

```yaml
# ServiceMonitor
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: dragonfly-metrics
spec:
  selector:
    matchLabels:
      app: dragonfly
  endpoints:
  - port: http
    path: /metrics
```

```promql
# 连接数
dragonfly_connected_clients

# 内存使用
dragonfly_used_memory

# 操作速率
rate(dragonfly_commands_processed_total[5m])
```
