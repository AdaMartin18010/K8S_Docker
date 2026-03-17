# Docker 架构详解

> Docker 从 CLI 到 runc 的完整调用链

---

## 架构总览

```
┌─────────────────────────────────────────────────────────────┐
│                      Docker Client                           │
│              (docker build, docker run, ...)                │
└──────────────────────┬──────────────────────────────────────┘
                       │ REST API / Unix Socket
┌──────────────────────▼──────────────────────────────────────┐
│                   Docker Daemon (dockerd)                    │
│  ┌─────────────────────────────────────────────────────┐   │
│  │   Image Management                                  │   │
│  │   - Build (BuildKit)                               │   │
│  │   - Pull/Push (Registry)                           │   │
│  │   - Layer Caching                                  │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │   Container Management                              │   │
│  │   - Lifecycle Management                           │   │
│  │   - Network Management                             │   │
│  │   - Volume Management                              │   │
│  └─────────────────────────────────────────────────────┘   │
└──────────────────────┬──────────────────────────────────────┘
                       │ containerd.sock (gRPC)
┌──────────────────────▼──────────────────────────────────────┐
│                      containerd                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Metadata   │  │  Content    │  │      Snapshot       │  │
│  │  (BoltDB)   │  │  (Images)   │  │     (Layers)        │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Task      │  │   Runtime   │  │      Diff           │  │
│  │ (Lifecycle) │  │  (Shim v2)  │  │   (Layer Ops)       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└──────────────────────┬──────────────────────────────────────┘
                       │ OCI Runtime (runc/crun)
┌──────────────────────▼──────────────────────────────────────┐
│                        runc                                  │
│         ┌─────────────────────────────┐                     │
│         │   OCI Runtime Spec          │                     │
│         │   - Namespace Setup         │                     │
│         │   - Cgroup Configuration    │                     │
│         │   - Capabilities            │                     │
│         └─────────────────────────────┘                     │
└─────────────────────────────────────────────────────────────┘
                              │
┌──────────────────────▼──────────────────────────────────────┐
│              Linux Kernel (Namespace/Cgroup)                 │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心组件职责

| 组件 | 职责 | 通信方式 |
|------|------|----------|
| **dockerd** | 处理 Docker API 请求、镜像构建、网络管理 | REST API / Unix Socket |
| **containerd** | 容器生命周期管理、镜像存储、执行器管理 | gRPC (containerd.sock) |
| **containerd-shim** | 隔离容器进程与 containerd，支持 daemon 重启 | stdio / ttrpc |
| **runc** | 根据 OCI 规范创建和运行容器 | CLI / libcontainer |

---

## BuildKit 架构 (Docker 20.10+)

```
Docker Client
     │
     ▼
BuildKit Daemon (LLB)
     │
     ├──▶ Solver (并发求解)
     │
     ├──▶ Cache Manager (缓存管理)
     │
     └──▶ Exporter (输出)
              │
              ├──▶ Docker Image
              ├──▶ OCI Image
              └──▶ Local Filesystem
```

---

## 命令调用链示例

```bash
# docker run 调用链
docker run nginx:alpine
    │
    ├──▶ Docker Client
    │       └── 发送 HTTP 请求到 dockerd
    │
    ├──▶ Docker Daemon
    │       ├── 检查本地镜像
    │       ├── 如不存在，从 Registry 拉取
    │       └── 调用 containerd 创建容器
    │
    ├──▶ containerd
    │       ├── 创建 containerd-shim
    │       └── 准备 OCI Runtime Spec
    │
    ├──▶ containerd-shim
    │       └── 调用 runc 启动容器
    │
    └──▶ runc
            ├── 设置 Namespace
            ├── 设置 Cgroups
            └── exec 容器进程
```

---

## 关键文件路径

| 组件 | 配置/数据路径 |
|------|--------------|
| dockerd | `/var/lib/docker/` |
| containerd | `/var/lib/containerd/` |
| containerd.sock | `/run/containerd/containerd.sock` |
| runc | `/usr/bin/runc` |
