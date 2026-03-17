# Docker 存储管理

> 容器数据持久化与存储驱动 (Docker 28.x)

---

## 存储类型

| 类型 | 说明 | 数据持久化 | 使用场景 |
|------|------|-----------|----------|
| **Volume** | Docker 管理的存储 | ✅ 是 | 数据库、持久化数据 |
| **Bind Mount** | 绑定宿主机路径 | ✅ 是 | 开发环境、配置文件 |
| **tmpfs** | 内存存储 | ❌ 否 | 敏感数据、临时文件 |

---

## Volume（推荐）

```
┌─────────────────────────────────────────────────────────────┐
│                        宿主机                                │
│                                                              │
│  ┌─────────────────┐      ┌─────────────────────────────┐  │
│  │   Container A   │      │   /var/lib/docker/volumes/  │  │
│  │   ┌─────────┐   │      │                             │  │
│  │   │ /data   │◄──┼──────┼─► my-volume/_data/          │  │
│  │   └─────────┘   │      │                             │  │
│  └────────┬────────┘      └─────────────────────────────┘  │
│           │                                                  │
│  ┌────────▼────────┐                                        │
│  │   Container B   │                                        │
│  │   ┌─────────┐   │                                        │
│  │   │ /backup │◄──┼────────────────────────────────────────┤
│  │   └─────────┘   │                                        │
│  └─────────────────┘                                        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 创建和使用 Volume

```bash
# 创建 volume
docker volume create my-data

# 查看 volumes
docker volume ls

# 查看详情
docker volume inspect my-data

# 运行容器并挂载
docker run -d \
  -v my-data:/app/data \
  --name myapp \
  myimage:latest

# 指定读写权限
docker run -d -v my-data:/app/data:ro nginx  # 只读

# 删除 volume
docker volume rm my-data

# 清理未使用的 volumes
docker volume prune
```

---

## Bind Mount

```bash
# 绑定宿主机目录到容器
docker run -d \
  -v /host/path:/container/path \
  nginx

# 相对路径（基于当前目录）
docker run -d \
  -v $(pwd)/html:/usr/share/nginx/html \
  nginx

# Windows 路径
docker run -d \
  -v E:/data:/data \
  nginx

# 只读绑定
docker run -d \
  -v /host/config:/etc/nginx:ro \
  nginx
```

### Volume vs Bind Mount

| 特性 | Volume | Bind Mount |
|------|--------|------------|
| 管理 | Docker 管理 | 用户管理 |
| 路径 | 无需知道具体路径 | 需指定完整路径 |
| 备份 | 易于备份迁移 | 自行管理 |
| 性能 | 原生性能 | 依赖文件系统 |
| 跨平台 | 一致行为 | 可能有差异 |

---

## tmpfs 挂载

```bash
# 内存存储（数据不写入磁盘）
docker run -d \
  --tmpfs /app/cache:noexec,nosuid,size=100m \
  myapp

# 特点：
# - 数据存储在内存中
# - 容器停止数据丢失
# - 无法共享
# - 高性能
```

---

## 多容器共享数据

```bash
# 创建共享 volume
docker volume create shared-data

# 数据生产者
docker run -d \
  --name producer \
  -v shared-data:/output \
  data-producer:latest

# 数据消费者
docker run -d \
  --name consumer \
  -v shared-data:/input \
  data-consumer:latest

# 备份容器
docker run --rm \
  -v shared-data:/data \
  -v $(pwd)/backup:/backup \
  alpine tar czf /backup/data.tar.gz -C /data .
```

---

## 存储驱动

| 驱动 | 说明 | 适用场景 |
|------|------|----------|
| **overlay2** | 默认，性能优秀 | 大多数场景 |
| **btrfs** | 快照功能 | 需要高级存储功能 |
| **zfs** | 企业级 | ZFS 文件系统环境 |
| **devicemapper** | 块设备 | 旧系统兼容 |

```bash
# 查看当前存储驱动
docker info | grep "Storage Driver"

# 修改存储驱动 (/etc/docker/daemon.json)
{
  "storage-driver": "overlay2"
}
```

---

## 数据备份与恢复

### 备份 Volume

```bash
# 方法1: 使用 tar
docker run --rm \
  -v my-volume:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/my-volume-backup.tar.gz -C /data .

# 方法2: 使用 cp
docker cp my-container:/app/data ./backup

# 方法3: 使用 volume backup 镜像
docker run --rm \
  -v my-volume:/source \
  -v $(pwd):/backup \
  busybox tar cvf /backup/backup.tar -C /source .
```

### 恢复 Volume

```bash
# 恢复备份到 volume
docker run --rm \
  -v my-volume:/data \
  -v $(pwd):/backup \
  alpine sh -c "cd /data && tar xzf /backup/my-volume-backup.tar.gz"
```

---

## Docker 28.x 存储新特性

| 特性 | 说明 |
|------|------|
| **Volume 快照** | 创建 volume 时间点快照 |
| **加密 Volume** | 静态数据加密 |
| **远程 Volume** | 支持 NFS、CIFS 等远程存储 |
| **性能监控** | 内置 volume I/O 监控 |

---

## 最佳实践

1. **优先使用 Volume**: 更好的可移植性和管理性
2. **避免在容器中存储数据**: 数据应在 volume 中
3. **定期备份**: 重要数据需要备份策略
4. **监控存储使用**: 防止磁盘空间耗尽
5. **使用只读挂载**: 配置文件等使用 :ro 挂载
6. **清理未使用数据**: 定期 `docker system prune`

---

## 常见问题

```bash
# Volume 权限问题
docker run -v my-data:/data --user $(id -u):$(id -g) myapp

# SELinux 问题
docker run -v /host/data:/container/data:Z myapp

# 磁盘空间不足
docker system df -v  # 查看空间使用
docker system prune   # 清理
```
