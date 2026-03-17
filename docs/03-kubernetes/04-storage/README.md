# Kubernetes 存储管理

> PV、PVC、StorageClass 与 StatefulSet

---

## 本章内容

1. [存储概述](./storage-overview.md)
2. [PV 与 PVC](./pv-pvc.md)
3. [StorageClass](./storage-class.md)
4. [StatefulSet 存储](./statefulset-storage.md)
5. [存储快照](./snapshots.md)

---

## 存储类型对比

| 类型 | 访问模式 | 适用场景 |
|------|----------|----------|
| **emptyDir** | 单 Pod 读写 | 临时缓存 |
| **hostPath** | 单节点读写 | 节点日志 |
| **Local PV** | 单节点读写 | 高性能数据库 |
| **Network PV** | 多节点读写 | 共享存储 (NFS) |
| **Cloud PV** | 单节点读写 | 云盘 (EBS/GCE/阿里云盘) |

---

## StorageClass 配置

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/gce-pd
parameters:
  type: pd-ssd
  replication-type: regional
volumeBindingMode: WaitForFirstConsumer
allowVolumeExpansion: true
reclaimPolicy: Retain
```

---

## PVC 使用示例

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-data
spec:
  storageClassName: fast-ssd
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100Gi
```

---

## StatefulSet 存储模板

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres
  replicas: 3
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        storageClassName: fast-ssd
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 50Gi
```

---

## 关联代码

- [examples/kubernetes/05-storage/](../../examples/kubernetes/05-storage/)
