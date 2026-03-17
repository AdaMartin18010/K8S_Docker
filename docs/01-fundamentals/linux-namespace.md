# Linux Namespace

> 容器隔离的核心技术

---

## 什么是 Namespace？

Namespace 是 Linux 内核提供的一种资源隔离机制，它让进程看到的是独立的系统资源视图。

---

## 七种 Namespace

| 类型 | 隔离资源 | 说明 |
|------|----------|------|
| **PID** | 进程 ID | 进程编号空间隔离 |
| **Network** | 网络设备 | 网络接口、IP、端口隔离 |
| **IPC** | 进程间通信 | 消息队列、共享内存、信号量 |
| **Mount** | 文件系统挂载 | 独立的挂载点视图 |
| **UTS** | 主机名和域名 | hostname 隔离 |
| **User** | 用户和组 ID | 用户权限隔离 |
| **Cgroup** | 控制组 | 资源限制视图隔离 (Linux 4.6+) |

---

## 实验：手动创建 Namespace

```bash
# 创建新的 PID Namespace
sudo unshare --fork --pid --mount-proc /bin/bash

# 在新 Namespace 中查看进程
ps aux
# 只能看到当前 Namespace 的进程

# 创建新的 Network Namespace
sudo ip netns add myns
sudo ip netns exec myns bash

# 在新 Namespace 中查看网络
ip addr
```

---

## Docker 中的 Namespace

```bash
# 查看容器的 Namespace
ls /proc/<pid>/ns/

# 输出：ipc  mnt  net  pid  user  uts
```

---

## Namespace 与容器

```
┌─────────────────────────────────────────┐
│  Container A (独立 Namespace)           │
│  ┌────────┐ ┌────────┐ ┌────────┐      │
│  │ PID ns │ │ NET ns │ │ MNT ns │      │
│  │ PID 1  │ │ eth0   │ │ /      │      │
│  └────────┘ └────────┘ └────────┘      │
├─────────────────────────────────────────┤
│  Container B (独立 Namespace)           │
│  ┌────────┐ ┌────────┐ ┌────────┐      │
│  │ PID ns │ │ NET ns │ │ MNT ns │      │
│  │ PID 1  │ │ eth0   │ │ /      │      │
│  └────────┘ └────────┘ └────────┘      │
├─────────────────────────────────────────┤
│         Linux Kernel (共享)             │
└─────────────────────────────────────────┘
```
