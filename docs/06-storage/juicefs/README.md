# JuiceFS - 高性能分布式文件系统

## 概述

JuiceFS 是一个云原生高性能分布式文件系统，采用数据与元数据分离架构，基于对象存储和数据库构建，提供 POSIX、HDFS、S3 等多种访问接口。

## 架构设计

```text
┌─────────────────────────────────────────────────────────────┐
│                        JuiceFS 架构                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│  │ POSIX    │  │ HDFS     │  │ S3       │  │ CSI      │     │
│  │ 客户端   │  │ SDK      │  │ Gateway  │  │ Driver   │     │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘     │
│       │             │             │             │           │
│       └─────────────┴──────┬──────┴─────────────┘           │
│                            │                                 │
│                   ┌────────┴────────┐                        │
│                   │  JuiceFS Client │                        │
│                   │  (FUSE 客户端)   │                        │
│                   └────────┬────────┘                        │
│                            │                                 │
│       ┌────────────────────┼────────────────────┐            │
│       ▼                    ▼                    ▼            │
│  ┌─────────┐        ┌─────────┐         ┌─────────┐         │
│  │ 对象存储 │        │ 元数据   │         │ 缓存    │         │
│  │ (Data)  │        │ 引擎     │         │ (Cache) │         │
│  │         │        │         │         │         │         │
│  │ S3      │        │ Redis   │         │ 本地磁盘 │         │
│  │ Azure   │◄──────►│ TiKV    │         │ 内存    │         │
│  │ GCS     │        │ MySQL   │         │         │         │
│  │ MinIO   │        │ SQLite  │         │         │         │
│  └─────────┘        └─────────┘         └─────────┘         │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## 核心特性

| 特性 | 描述 |
|------|------|
| 弹性扩展 | 支持千亿级文件，容量无上限 |
| 高性能 | 分布式缓存，聚合读带宽 1.2TB/s |
| POSIX 兼容 | 完整 POSIX 兼容，应用无需修改 |
| 多云支持 | 支持 AWS、Azure、GCP、阿里云等 |
| 数据安全 | 支持加密、备份、快照 |

## 安装部署

### 快速开始

```bash
# 1. 下载 JuiceFS
wget https://github.com/juicedata/juicefs/releases/download/v1.3.0/juicefs-1.3.0-linux-amd64.tar.gz
tar -xzf juicefs-1.3.0-linux-amd64.tar.gz
sudo mv juicefs /usr/local/bin/

# 2. 创建文件系统（使用 Redis 作为元数据引擎）
juicefs format \
  --storage s3 \
  --bucket https://mybucket.s3.amazonaws.com \
  --access-key YOUR_ACCESS_KEY \
  --secret-key YOUR_SECRET_KEY \
  redis://localhost:6379/1 \
  myjfs

# 3. 挂载文件系统
sudo juicefs mount -d redis://localhost:6379/1 /mnt/juicefs

# 4. 验证
df -h /mnt/juicefs
cp -r ~/data /mnt/juicefs/
```

### Kubernetes CSI 驱动

```bash
# 安装 CSI 驱动
helm repo add juicefs https://juicedata.github.io/juicefs-csi-driver/
helm repo update

helm install juicefs-csi-driver juicefs/juicefs-csi-driver \
  --namespace kube-system \
  --set storageClasses[0].name=juicefs \
  --set storageClasses[0].enabled=true
```

```yaml
# StorageClass
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: juicefs
provisioner: csi.juicefs.com
parameters:
  csi.storage.k8s.io/provisioner-secret-name: juicefs-secret
  csi.storage.k8s.io/provisioner-secret-namespace: kube-system
  csi.storage.k8s.io/node-publish-secret-name: juicefs-secret
  csi.storage.k8s.io/node-publish-secret-namespace: kube-system
reclaimPolicy: Retain
volumeBindingMode: Immediate
---
# Secret
apiVersion: v1
kind: Secret
metadata:
  name: juicefs-secret
  namespace: kube-system
type: Opaque
stringData:
  name: myjfs
  metaurl: redis://redis:6379/1
  storage: s3
  bucket: https://mybucket.s3.amazonaws.com
  access-key: YOUR_ACCESS_KEY
  secret-key: YOUR_SECRET_KEY
---
# PVC
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: juicefs-pvc
spec:
  accessModes:
  - ReadWriteMany
  storageClassName: juicefs
  resources:
    requests:
      storage: 100Gi
```

## 元数据引擎选择

| 引擎 | 规模 | 延迟 | 适用场景 |
|------|------|------|----------|
| Redis | < 1 亿文件 | 亚毫秒 | 小规模、高性能 |
| TiKV | < 100 亿文件 | 毫秒级 | 大规模生产 |
| MySQL/PostgreSQL | < 10 亿文件 | 毫秒级 | 企业级 |
| SQLite | < 100 万文件 | 亚毫秒 | 单节点测试 |

### TiKV 配置（推荐生产环境）

```bash
# 格式化使用 TiKV
juicefs format \
  --storage s3 \
  --bucket https://mybucket.s3.amazonaws.com \
  tikv://pd-host:2379/juicefs \
  myjfs
```

## 分布式缓存

### Cache Group 配置

```yaml
# cache-group.yaml
apiVersion: juicefs.io/v1
kind: CacheGroup
metadata:
  name: ai-training-cache
spec:
  secretName: juicefs-secret
  secretNamespace: kube-system

  # 缓存组配置
  cacheGroup:
    name: ai-cache

  # 工作节点选择器
  workerSelector:
    matchLabels:
      cache-node: "true"

  # 缓存配置
  options:
    cache-size: 500Gi
    cache-dir: /var/jfsCache
    free-space-ratio: 0.1

  # 节点配置
  workerTemplate:
    resources:
      requests:
        memory: 2Gi
        cpu: "1"
      limits:
        memory: 8Gi
        cpu: "4"
    volumes:
    - name: cache-dir
      hostPath:
        path: /var/jfsCache
        type: DirectoryOrCreate
```

## 性能优化

### 挂载参数调优

```bash
# 大文件读取优化（AI 训练）
juicefs mount -d \
  --cache-size=500000 \
  --cache-dir=/var/jfsCache \
  --free-space-ratio=0.1 \
  --writeback \
  --buffer-size=1024 \
  redis://localhost:6379/1 /mnt/juicefs

# 小文件优化
juicefs mount -d \
  --cache-size=100000 \
  --open-cache=1h \
  --attr-cache=1h \
  --entry-cache=1h \
  redis://localhost:6379/1 /mnt/juicefs
```

### 参数说明

| 参数 | 说明 | 推荐值 |
|------|------|--------|
| cache-size | 本地缓存大小(MB) | 500000 |
| writeback | 异步写入 | 高吞吐场景 |
| buffer-size | 缓冲区大小(MB) | 1024 |
| max-uploads | 并发上传数 | 50 |
| max-downloads | 并发下载数 | 50 |

## 监控与运维

### Prometheus 指标

```yaml
# juicefs-monitor.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: juicefs-metrics
data:
  metrics.json: |
    {
      "metrics": [
        {"name": "juicefs_used_space", "help": "Used space in bytes"},
        {"name": "juicefs_available_space", "help": "Available space"},
        {"name": "juicefs_inode_total", "help": "Total inodes"},
        {"name": "juicefs_inode_used", "help": "Used inodes"},
        {"name": "juicefs_ops_total", "help": "Total operations"},
        {"name": "juicefs_ops_durations_histogram", "help": "Operation latency"},
        {"name": "juicefs_blockcache_hit", "help": "Block cache hit ratio"},
        {"name": "juicefs_blockcache_miss", "help": "Block cache miss"}
      ]
    }
```

### 运维命令

```bash
# 查看文件系统状态
juicefs status redis://localhost:6379/1

# 扫描文件系统
juicefs gc redis://localhost:6379/1

# 一致性检查
juicefs fsck redis://localhost:6379/1

# 备份元数据
juicefs dump redis://localhost:6379/1 backup.json

# 恢复元数据
juicefs load redis://localhost:6379/1 backup.json

# 查看实时统计
juicefs stats /mnt/juicefs
```

## 企业版特性

- **多区域架构**: 跨 AZ 部署，1ms 元数据延迟
- **RDMA 支持**: 高性能网络传输
- **Quota 管理**: 用户/组级别配额
- **快照克隆**: 秒级快照和克隆
- **数据分层**: 热温冷数据自动分层

## 2025 新特性

- **Python SDK**: 原生 Python 支持
- **Windows 客户端**: 完整 Windows 支持
- **Apache Ranger**: 细粒度权限管理
- **百亿文件**: 单卷支持 500 亿文件
- **LRU 缓存**: 智能缓存淘汰

## 应用场景

| 场景 | 优势 |
|------|------|
| AI/ML 训练 | 高吞吐、POSIX 兼容 |
| 大数据 | HDFS 兼容、低成本 |
| Kubernetes | CSI 驱动、共享存储 |
| 备份归档 | 低成本对象存储 |
| 内容分发 | 全球缓存加速 |

## 相关资源

- [JuiceFS 官网](https://juicefs.com/)
- [GitHub](https://github.com/juicedata/juicefs)
- [CSI 驱动](https://github.com/juicedata/juicefs-csi-driver)
- [社区论坛](https://github.com/juicedata/juicefs/discussions)
