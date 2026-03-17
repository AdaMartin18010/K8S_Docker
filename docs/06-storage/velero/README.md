# Velero - Kubernetes 备份与灾难恢复

## 概述

Velero 是 VMware 开源的 Kubernetes 备份、恢复和迁移工具。它可以备份 Kubernetes 集群资源和持久卷（PV），支持灾难恢复、集群迁移和数据保护。2025 年，Velero 1.17 版本进一步优化了 CSI 快照支持，成为 Kubernetes 数据保护的事实标准。

## 架构原理

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Velero 架构                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Velero Server                                  │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│  │  │   Backup     │  │   Restore    │  │  Schedule    │              │   │
│  │  │  Controller  │  │  Controller  │  │  Controller  │              │   │
│  │  └──────┬───────┘  └──────────────┘  └──────────────┘              │   │
│  │         │                                                          │   │
│  │  ┌──────▼───────────────────────────────────────────────────────┐  │   │
│  │  │              Backup Storage Location (S3/MinIO)               │  │   │
│  │  │  • Kubernetes 资源 (YAML)                                     │  │   │
│  │  │  • Volume 快照元数据                                          │  │   │
│  │  │  • Kopia/Restic 备份数据                                      │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 安装 Velero

### 安装 CLI

```bash
curl -L https://github.com/vmware-tanzu/velero/releases/download/v1.17.0/velero-v1.17.0-linux-amd64.tar.gz -o velero.tar.gz
tar -xzf velero.tar.gz
sudo mv velero-v1.17.0-linux-amd64/velero /usr/local/bin/
velero version
```

### 安装服务端

```bash
# 准备 AWS 凭证
cat > aws-credentials <<EOF
[default]
aws_access_key_id=<AWS_ACCESS_KEY_ID>
aws_secret_access_key=<AWS_SECRET_ACCESS_KEY>
EOF

# 安装 Velero
velero install \
  --provider aws \
  --plugins velero/velero-plugin-for-aws:v1.13.2 \
  --bucket velero-backups \
  --backup-location-config region=us-east-1 \
  --snapshot-location-config region=us-east-1 \
  --secret-file aws-credentials \
  --use-volume-snapshots=true \
  --use-node-agent \
  --features=EnableCSI

# 验证安装
kubectl get pods -n velero
```

## 基础使用示例

### 命名空间备份

```bash
# 备份整个命名空间
velero backup create nginx-backup \
  --include-namespaces nginx \
  --wait

# 包含 PV 快照
velero backup create nginx-backup-with-pv \
  --include-namespaces nginx \
  --snapshot-volumes \
  --wait

# 使用数据移动器（跨集群恢复）
velero backup create nginx-backup-data-mover \
  --include-namespaces nginx \
  --snapshot-move-data=true \
  --wait
```

### 定时备份

```yaml
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: daily-backup
  namespace: velero
spec:
  schedule: "0 2 * * *"  # 每天凌晨 2 点
  template:
    includedNamespaces:
    - production
    snapshotVolumes: true
    ttl: 720h0m0s  # 保留 30 天
```

### 备份恢复

```bash
# 列出备份
velero backup get

# 查看备份详情
velero backup describe nginx-backup --details

# 恢复整个命名空间
velero restore create --from-backup nginx-backup --wait

# 恢复到不同命名空间
velero restore create --from-backup nginx-backup \
  --namespace-mappings nginx:nginx-restore \
  --wait
```

## 集群迁移

```bash
# 1. 在源集群创建备份
velero backup create migration-backup \
  --include-namespaces app \
  --snapshot-move-data=true \
  --wait

# 2. 在目标集群安装 Velero
velero install \
  --provider aws \
  --bucket velero-backups \
  --backup-location-config region=us-east-1 \
  --secret-file aws-credentials

# 3. 在目标集群同步备份
velero backup sync

# 4. 在目标集群恢复
velero restore create --from-backup migration-backup --wait
```

## 故障排查

```bash
# 查看备份日志
velero backup logs nginx-backup

# 查看 Velero 服务器日志
kubectl logs -n velero deployment/velero

# 检查 VolumeSnapshot
kubectl get volumesnapshot -n nginx

# 调试备份失败
velero backup describe nginx-backup --details
```

## 总结

| 场景 | 推荐方案 |
|------|----------|
| 日常备份 | Schedule + CSI 快照 |
| 灾难恢复 | 完整备份 + 定期演练 |
| 集群迁移 | 数据移动器 + 跨集群恢复 |
| 长期归档 | Kopia 文件级备份 |

Velero 是 Kubernetes 数据保护的标准工具。
