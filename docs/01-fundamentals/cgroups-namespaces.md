# Cgroups 与 Namespaces

> Linux 容器核心技术

---

## Namespaces (隔离)

| Namespace | 隔离资源 | 系统调用参数 |
|-----------|----------|--------------|
| **PID** | 进程 ID | CLONE_NEWPID |
| **Network** | 网络设备 | CLONE_NEWNET |
| **Mount** | 挂载点 | CLONE_NEWNS |
| **UTS** | 主机名 | CLONE_NEWUTS |
| **IPC** | 进程通信 | CLONE_NEWIPC |
| **User** | 用户/组 | CLONE_NEWUSER |
| **Cgroup** | Cgroup 根 | CLONE_NEWCGROUP |

---

## Cgroups (资源限制)

| 子系统 | 功能 |
|--------|------|
| **cpu** | CPU 使用限制 |
| **memory** | 内存限制 |
| **blkio** | 块设备 I/O |
| **pids** | 进程数量 |

---

## 实际演示

```bash
# 查看容器 namespace
ls -la /proc/<pid>/ns/

# 查看 cgroup
ls -la /proc/<pid>/cgroup

# 进入容器 namespace
nsenter -t <pid> -n ip addr
```

---

## Docker 中的应用

```bash
# Docker 创建容器时设置
docker run -m 512m --cpus=1.5 nginx

# 对应 cgroup 限制
cat /sys/fs/cgroup/memory/docker/<id>/memory.limit_in_bytes
cat /sys/fs/cgroup/cpu/docker/<id>/cpu.cfs_quota_us
```

---

## 为什么需要两者？

| 技术 | 作用 |
|------|------|
| **Namespaces** | 隔离视图 (我能看到什么) |
| **Cgroups** | 限制资源 (我能用多少) |

两者结合才能实现完整的容器隔离。
