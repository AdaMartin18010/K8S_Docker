# Linux Cgroups

> 容器资源限制的核心技术

---

## 什么是 Cgroups？

Cgroups (Control Groups) 是 Linux 内核提供的一种资源限制和统计机制。

---

## Cgroups v1 vs v2

| 特性 | Cgroups v1 | Cgroups v2 |
|------|------------|------------|
| 架构 | 多层级 (每个子系统独立) | 统一层级 |
| 控制器 | 独立挂载 | 统一挂载 |
| 根权限 | 需要 | 支持 delegation |
| 默认 | 旧系统 | 新系统 (RHEL 9, Ubuntu 22.04+) |

---

## Cgroups 控制器

| 控制器 | 说明 | 用途 |
|--------|------|------|
| **cpu** | CPU 时间分配 | 限制 CPU 使用率 |
| **memory** | 内存限制 | 限制内存使用 |
| **blkio** | 块设备 IO | 限制磁盘 IO |
| **pids** | 进程数 | 限制进程数量 |
| **devices** | 设备访问 | 控制设备访问权限 |
| **net_cls/net_prio** | 网络 | 标记和优先级 |

---

## 实验：手动设置 Cgroups

```bash
# 创建 cgroup
sudo mkdir /sys/fs/cgroup/memory/demo

# 设置内存限制 (100MB)
echo 104857600 | sudo tee /sys/fs/cgroup/memory/demo/memory.limit_in_bytes

# 将进程加入 cgroup
echo $$ | sudo tee /sys/fs/cgroup/memory/demo/tasks

# 查看统计
cat /sys/fs/cgroup/memory/demo/memory.usage_in_bytes
```

---

## Docker 资源限制

```bash
# CPU 限制
docker run --cpus=0.5 --memory=512m nginx:alpine

# 对应 Cgroups 配置
cat /sys/fs/cgroup/docker/<container_id>/memory.limit_in_bytes
```

---

## K8s 资源限制

```yaml
resources:
  requests:
    cpu: 100m      # 0.1 核
    memory: 128Mi
  limits:
    cpu: 500m      # 0.5 核
    memory: 512Mi
```
