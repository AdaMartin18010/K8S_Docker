# DragonflyDB - 现代 Redis 替代方案

## 概述

DragonflyDB 是一个现代化的内存数据存储系统，完全兼容 Redis 协议，采用多线程架构，在性能、资源效率和成本方面相比 Redis 有显著提升。

## 核心优势

| 特性 | DragonflyDB | Redis |
|------|-------------|-------|
| 线程模型 | 多线程 | 单线程 |
| 吞吐量 | 最高 4M QPS | 150K QPS |
| 内存效率 | 节省 ~30% | 标准 |
| 快照性能 | 无 Fork 增量 | Fork 阻塞 |
| 协议兼容 | Redis 6.2 (95-98%) | 原生 |

## 架构特点

### 多线程共享架构

- **无锁数据结构**: Dashtable 实现高效并发访问
- **CPU 亲和性**: 每个核心处理独立的键范围
- **零拷贝网络**: 优化网络 I/O 性能

### 内存优化

```
Redis 内存占用:    10GB
DragonflyDB 占用:  ~7GB (节省 30%)
```

### 快照机制对比

```
Redis BGSAVE:
- Fork 进程，内存翻倍
- 大实例时阻塞明显

Dragonfly:
- 增量快照，无 Fork
- 内存占用稳定
- 不影响在线服务
```

## 安装部署

### Docker 快速启动

```bash
# 单节点
docker run -d --name dragonfly \
  -p 6379:6379 \
  -v dragonfly-data:/data \
  docker.dragonflydb.io/dragonflydb/dragonfly:latest

# 带密码
docker run -d --name dragonfly \
  -p 6379:6379 \
  -e DFLY_requirepass=yourpassword \
  -v dragonfly-data:/data \
  docker.dragonflydb.io/dragonflydb/dragonfly:latest
```

### Kubernetes 部署

```yaml
# dragonfly-deployment.yaml
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
        image: docker.dragonflydb.io/dragonflydb/dragonfly:v1.25.0
        ports:
        - containerPort: 6379
          name: redis
        command:
        - dragonfly
        - --dir=/data
        - --dbfilename=dump.rdb
        - --maxmemory=4gb
        - --snapshot_cron=*/30 * * * *
        resources:
          requests:
            memory: "4Gi"
            cpu: "2"
          limits:
            memory: "8Gi"
            cpu: "4"
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
          storage: 100Gi
```

## 客户端连接

### Go 示例

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

func main() {
    // DragonflyDB 完全兼容 Redis 客户端
    rdb := redis.NewClient(&redis.Options{
        Addr:     "dragonfly:6379",
        Password: "",
        DB:       0,
        PoolSize: 100,
    })

    ctx := context.Background()

    // 基本操作
    err := rdb.Set(ctx, "key", "value", time.Hour).Err()
    if err != nil {
        panic(err)
    }

    val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key:", val)
}
```

## 相关资源

- [DragonflyDB 官网](https://www.dragonflydb.io/)
- [GitHub](https://github.com/dragonflydb/dragonfly)
