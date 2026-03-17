# Kubernetes 存储管理

本目录包含 PV、PVC、StorageClass 和 StatefulSet 的最佳实践。

## 存储类型对比

| 类型 | 访问模式 | 适用场景 | 性能 |
|------|---------|----------|------|
| **标准磁盘** | RWO | 普通应用 | 中等 |
| **SSD** | RWO | 数据库 | 高 |
| **本地 SSD** | RWO | 高性能缓存 | 极高 |
| **NFS** | RWX | 共享文件 | 中等 |
| **对象存储** | - | 备份、静态资源 | 低延迟 |

## 文件说明

| 文件 | 描述 |
|------|------|
| `statefulset-good.yaml` | StatefulSet 完整配置 |
| `pv-pvc-good.yaml` | PV/PVC/StorageClass 示例 |

## 关键概念

### Access Modes

- **RWO** (ReadWriteOnce): 单节点读写
- **ROX** (ReadOnlyMany): 多节点只读
- **RWX** (ReadWriteMany): 多节点读写
- **RWOP** (ReadWriteOncePod): 单 Pod 读写（1.22+）

### Reclaim Policy

- **Retain**: 保留数据，需要手动清理
- **Delete**: 自动删除 PV 和数据
- **Recycle**: 擦除后重用（已废弃）

## 使用指南

### 创建存储

```bash
# 创建 StorageClass
kubectl apply -f pv-pvc-good.yaml

# 创建 PVC
kubectl get pvc -n production

# 查看 PV
kubectl get pv

# 查看存储类
kubectl get storageclass
```

### 使用 StatefulSet

```bash
# 部署 StatefulSet
kubectl apply -f statefulset-good.yaml

# 查看 Pod（注意有序的命名）
kubectl get pods -n database
# NAME         READY   STATUS
# postgres-0   1/1     Running
# postgres-1   1/1     Running
# postgres-2   1/1     Running

# 查看 PVC（自动创建）
kubectl get pvc -n database
# NAME            STATUS   VOLUME
# data-postgres-0 Bound    pvc-xxx
# data-postgres-1 Bound    pvc-yyy
# data-postgres-2 Bound    pvc-zzz

# 扩缩容
kubectl scale statefulset postgres --replicas=5 -n database

# 滚动更新（灰度）
# 修改 updateStrategy.partition 为 2，只更新 postgres-2
kubectl patch statefulset postgres -n database -p '{"spec":{"updateStrategy":{"rollingUpdate":{"partition":2}}}}'
```

### 备份和恢复

```bash
# 手动备份
kubectl exec -it postgres-0 -n database -- pg_dump -U postgres myapp > backup.sql

# 恢复
kubectl exec -i postgres-0 -n database -- psql -U postgres myapp < backup.sql

# 使用 CronJob 自动备份（已在配置中定义）
kubectl get cronjob -n database
```

## 最佳实践

1. **使用 StorageClass**: 避免手动创建 PV
2. **延迟绑定**: `volumeBindingMode: WaitForFirstConsumer`
3. **允许扩容**: `allowVolumeExpansion: true`
4. **设置配额**: 防止资源耗尽
5. **定期备份**: 使用 CronJob 自动备份
6. **监控存储**: 监控使用量和性能

## 故障排查

```bash
# Pod 无法挂载卷
kubectl describe pod <pod-name>
kubectl logs <pod-name>

# PVC 无法绑定
kubectl describe pvc <pvc-name>

# 查看存储事件
kubectl get events --field-selector reason=FailedMount

# 检查节点磁盘
kubectl get nodes -o json | jq '.items[].status.capacity'
```
