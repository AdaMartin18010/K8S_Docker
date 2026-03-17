# Kubernetes 存储管理

> PV、PVC、StorageClass 与 StatefulSet 存储 (K8s 1.33/1.34)

---

## 本章内容

- [存储概述](./README.md)
- [存储卷](./volumes.md)
- [持久化存储](./volumes.md)

---

## 核心概念

```
┌─────────────────────────────────────────────────────────────────┐
│                     Kubernetes 存储体系                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐  │
│  │     Pod      │ ───► │     PVC      │ ───► │      PV      │  │
│  │              │      │ (请求存储)    │      │  (实际存储)   │  │
│  └──────────────┘      └──────┬───────┘      └──────┬───────┘  │
│                               │                     │          │
│                               └─────────────────────┘          │
│                                         │                      │
│                               ┌─────────▼─────────┐            │
│                               │   StorageClass    │            │
│                               │  (动态供给模板)    │            │
│                               └───────────────────┘            │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 存储类型对比

| 类型 | 访问模式 | 生命周期 | 适用场景 |
|------|----------|----------|----------|
| **emptyDir** | 单 Pod 读写 | Pod | 临时缓存、共享空间 |
| **hostPath** | 单节点读写 | Pod | 节点日志、监控数据 |
| **Local PV** | 单节点读写 | PVC | 高性能数据库 |
| **Network PV** | 多节点读写 | PVC | 共享存储 (NFS) |
| **Cloud PV** | 单节点读写 | PVC | 云盘 (EBS/GCE) |

---

## emptyDir - 临时存储

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cache-pod
spec:
  containers:
  - name: app
    image: nginx
    volumeMounts:
    - name: cache
      mountPath: /cache
  volumes:
  - name: cache
    emptyDir:
      medium: ""  # 空字符串=磁盘，Memory=内存
      sizeLimit: 1Gi
```

**特点**:

- Pod 创建时创建，删除时删除
- 同一 Pod 内多个容器共享
- 可指定内存作为介质（tmpfs）

---

## hostPath - 宿主机路径

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: node-exporter
spec:
  containers:
  - name: exporter
    image: prom/node-exporter
    volumeMounts:
    - name: sys
      mountPath: /host/sys
      readOnly: true
  volumes:
  - name: sys
    hostPath:
      path: /sys
      type: Directory  # DirectoryOrCreate, File, Socket
```

**⚠️ 注意**: 生产环境慎用，节点迁移时数据不跟随

---

## PV 与 PVC

### 静态 PV（管理员创建）

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-nfs
spec:
  capacity:
    storage: 100Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  storageClassName: nfs
  nfs:
    server: 192.168.1.100
    path: /exports/data
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-nfs
spec:
  accessModes:
    - ReadWriteMany
  storageClassName: nfs
  resources:
    requests:
      storage: 50Gi
```

### 动态 PVC（自动创建 PV）

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-dynamic
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: fast-ssd
  resources:
    requests:
      storage: 100Gi
```

---

## StorageClass

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/gce-pd
parameters:
  type: pd-ssd
  replication-type: regional
volumeBindingMode: WaitForFirstConsumer  # 延迟绑定到 Pod 调度节点
allowVolumeExpansion: true               # 允许扩容
reclaimPolicy: Delete                    # 删除 PVC 时删除 PV
mountOptions:
  - debug
```

### 常用 StorageClass 配置

| 云厂商 | Provisioner | 参数 |
|--------|-------------|------|
| AWS | ebs.csi.aws.com | type: gp3, iops: 10000 |
| GCP | pd.csi.storage.gke.io | type: pd-ssd, replication-type: regional |
| Azure | disk.csi.azure.com | skuName: Premium_LRS |
| 本地 | local-volume-provisioner | 本地磁盘 |

---

## StatefulSet 存储

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres
  replicas: 3
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:16
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      storageClassName: fast-ssd
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 100Gi
```

**生成的 PVC**: `data-postgres-0`, `data-postgres-1`, `data-postgres-2`

---

## 存储快照（K8s 1.33+）

```yaml
# 创建 VolumeSnapshotClass
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshotClass
metadata:
  name: csi-snapclass
driver: ebs.csi.aws.com
deletionPolicy: Delete
---
# 创建快照
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: mysql-snapshot
spec:
  volumeSnapshotClassName: csi-snapclass
  source:
    persistentVolumeClaimName: mysql-data
---
# 从快照恢复
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-data-restored
spec:
  storageClassName: gp3
  dataSource:
    name: mysql-snapshot
    kind: VolumeSnapshot
    apiGroup: snapshot.storage.k8s.io
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100Gi
```

---

## 存储扩容

```yaml
# 1. StorageClass 需要 allowVolumeExpansion: true
# 2. 直接修改 PVC
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
spec:
  resources:
    requests:
      storage: 200Gi  # 从 100Gi 扩容到 200Gi
```

```bash
# 在线扩容（无需重启 Pod）
kubectl patch pvc my-pvc -p '{"spec":{"resources":{"requests":{"storage":"200Gi"}}}}'
```

---

## K8s 1.33/1.34 存储新特性

| 特性 | 状态 | 说明 |
|------|------|------|
| **ReadWriteOncePod** | GA | 单 Pod 独占读写 |
| **VolumeAttributesClass** | Beta | 动态修改卷参数 |
| **跨命名空间快照** | Beta | 快照可跨 NS 恢复 |
| **FSGroup 策略优化** | GA | 更快的 chown |

---

## 最佳实践

1. **使用 StorageClass**: 避免手动管理 PV
2. **合理设置 reclaimPolicy**: 生产环境建议 Retain
3. **使用延迟绑定**: `volumeBindingMode: WaitForFirstConsumer`
4. **定期快照**: 重要数据启用自动快照
5. **监控存储**: 使用 kubelet_volume_stats 指标
6. **使用 ReadWriteOncePod**: 数据库等独占访问场景

---

## 故障排查

```bash
# 查看 PVC 状态
kubectl get pvc

# 查看 PV 详情
kubectl describe pv <pv-name>

# 查看存储事件
kubectl get events --field-selector reason=FailedMount

# 检查 CSI 驱动
kubectl get csidrivers

# 查看节点存储
kubectl get nodes -o json | jq '.items[].status.capacity'
```
