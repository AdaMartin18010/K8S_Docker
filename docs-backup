# 第一章：Docker 与 Kubernetes 体系结构深度解析

> 本章作为容器技术权威文档的核心章节，全面剖析 Docker 和 Kubernetes 的底层架构设计、核心组件交互机制以及生态体系。

---

## 目录

- [第一章：Docker 与 Kubernetes 体系结构深度解析](#第一章docker-与-kubernetes-体系结构深度解析)
  - [目录](#目录)
  - [1. Docker 核心架构](#1-docker-核心架构)
    - [1.1 Docker Engine 分层架构](#11-docker-engine-分层架构)
      - [核心组件职责](#核心组件职责)
    - [1.2 OCI 开放容器倡议规范](#12-oci-开放容器倡议规范)
      - [OCI Runtime Spec](#oci-runtime-spec)
      - [OCI Image Spec](#oci-image-spec)
    - [1.3 Linux 内核容器技术](#13-linux-内核容器技术)
      - [Namespace 隔离机制](#namespace-隔离机制)
      - [Cgroup v2 资源控制](#cgroup-v2-资源控制)
      - [Union Filesystem 对比](#union-filesystem-对比)
    - [1.4 容器镜像结构](#14-容器镜像结构)
    - [1.5 Docker 网络模型](#15-docker-网络模型)
      - [Bridge 网络数据包流程](#bridge-网络数据包流程)
  - [2. Kubernetes 核心架构](#2-kubernetes-核心架构)
    - [2.1 整体架构概览](#21-整体架构概览)
    - [2.2 核心组件详解](#22-核心组件详解)
      - [API Server (kube-apiserver)](#api-server-kube-apiserver)
      - [etcd - 集群状态存储](#etcd---集群状态存储)
      - [Scheduler 调度器](#scheduler-调度器)
    - [2.3 核心对象关系](#23-核心对象关系)
    - [2.4 控制器模式详解](#24-控制器模式详解)
  - [3. 接口标准：CRI/CNI/CSI](#3-接口标准cricnicsi)
    - [3.1 CRI (Container Runtime Interface)](#31-cri-container-runtime-interface)
      - [CRI 运行时对比](#cri-运行时对比)
    - [3.2 CNI (Container Network Interface)](#32-cni-container-network-interface)
      - [CNI 插件对比](#cni-插件对比)
    - [3.3 CSI (Container Storage Interface)](#33-csi-container-storage-interface)
  - [4. 网络与存储体系](#4-网络与存储体系)
    - [4.1 Service 实现机制](#41-service-实现机制)
    - [4.2 Ingress 控制器对比](#42-ingress-控制器对比)
    - [4.3 存储类型详解](#43-存储类型详解)
    - [4.4 PV/PVC 生命周期](#44-pvpvc-生命周期)
  - [5. 安全体系](#5-安全体系)
    - [5.1 认证机制](#51-认证机制)
    - [5.2 RBAC 授权模型](#52-rbac-授权模型)
    - [5.3 准入控制器](#53-准入控制器)
    - [5.4 Pod 安全标准](#54-pod-安全标准)
  - [6. 生态工具链](#6-生态工具链)
    - [6.1 包管理工具对比](#61-包管理工具对比)
    - [6.2 GitOps 工具对比](#62-gitops-工具对比)
    - [6.3 可观测性栈](#63-可观测性栈)
    - [6.4 服务网格对比](#64-服务网格对比)
    - [6.5 Serverless 平台对比](#65-serverless-平台对比)
  - [7. WebAssembly 与容器融合](#7-webassembly-与容器融合)
    - [7.1 WebAssembly 运行时对比](#71-webassembly-运行时对比)
    - [7.2 Wasm 与容器对比](#72-wasm-与容器对比)
    - [7.3 containerd-wasm-shim 架构](#73-containerd-wasm-shim-架构)
  - [8. Go 代码实战](#8-go-代码实战)
    - [8.1 使用 client-go 操作 K8s 资源](#81-使用-client-go-操作-k8s-资源)
    - [8.2 自定义 Controller 实现](#82-自定义-controller-实现)
    - [8.3 CRI 客户端调用](#83-cri-客户端调用)
    - [8.4 自定义 CNI 插件框架](#84-自定义-cni-插件框架)
    - [8.5 高级 client-go 操作](#85-高级-client-go-操作)
  - [9. 总结与最佳实践](#9-总结与最佳实践)
    - [9.1 架构设计原则](#91-架构设计原则)
    - [9.2 生产环境建议](#92-生产环境建议)
    - [9.3 性能优化要点](#93-性能优化要点)
  - [参考资源](#参考资源)

---

## 1. Docker 核心架构

### 1.1 Docker Engine 分层架构

Docker Engine 采用分层架构设计，将容器生命周期管理职责分离到不同组件：

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker CLI / API                          │
│              (docker run, docker build, etc.)               │
└──────────────────────┬──────────────────────────────────────┘
                       │ REST API / Unix Socket
┌──────────────────────▼──────────────────────────────────────┐
│                  Docker Daemon (dockerd)                     │
│         ┌─────────────────────────────────────┐              │
│         │   Image Management                  │              │
│         │   - Build (BuildKit)                │              │
│         │   - Pull/Push (Registry)            │              │
│         │   - Layer Caching                   │              │
│         └─────────────────────────────────────┘              │
└──────────────────────┬──────────────────────────────────────┘
                       │ containerd.sock (gRPC)
┌──────────────────────▼──────────────────────────────────────┐
│                    containerd                                │
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
│                      runc                                    │
│         ┌─────────────────────────────┐                     │
│         │   OCI Runtime Spec          │                     │
│         │   - Namespace Setup         │                     │
│         │   - Cgroup Configuration    │                     │
│         │   - Capabilities            │                     │
│         │   - Seccomp/AppArmor        │                     │
│         └─────────────────────────────┘                     │
└─────────────────────────────────────────────────────────────┘
                              │
┌──────────────────────▼──────────────────────────────────────┐
│              Linux Kernel (Namespace/Cgroup)                 │
└─────────────────────────────────────────────────────────────┘
```

#### 核心组件职责

| 组件 | 职责 | 通信方式 |
|------|------|----------|
| **dockerd** | 处理 Docker API 请求、镜像构建、网络管理 | REST API / Unix Socket |
| **containerd** | 容器生命周期管理、镜像存储、执行器管理 | gRPC (containerd.sock) |
| **containerd-shim** | 隔离容器进程与 containerd，支持 daemon 重启 | stdio / ttrpc |
| **runc** | 根据 OCI 规范创建和运行容器 | CLI / libcontainer |

### 1.2 OCI 开放容器倡议规范

OCI（Open Container Initiative）定义了容器运行时的标准化规范：

#### OCI Runtime Spec

```json
{
  "ociVersion": "1.0.2",
  "process": {
    "terminal": false,
    "user": {"uid": 0, "gid": 0},
    "args": ["sh", "-c", "echo hello"],
    "env": ["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],
    "cwd": "/",
    "capabilities": {
      "bounding": ["CAP_CHOWN", "CAP_DAC_OVERRIDE"],
      "effective": ["CAP_CHOWN", "CAP_DAC_OVERRIDE"],
      "permitted": ["CAP_CHOWN", "CAP_DAC_OVERRIDE"]
    }
  },
  "root": {
    "path": "rootfs",
    "readonly": true
  },
  "hostname": "mycontainer",
  "mounts": [
    {
      "destination": "/proc",
      "type": "proc",
      "source": "proc"
    },
    {
      "destination": "/sys",
      "type": "sysfs",
      "source": "sysfs",
      "options": ["nosuid", "noexec", "nodev", "ro"]
    }
  ],
  "linux": {
    "namespaces": [
      {"type": "pid"},
      {"type": "network"},
      {"type": "ipc"},
      {"type": "uts"},
      {"type": "mount"},
      {"type": "cgroup"}
    ],
    "cgroupsPath": "/docker/mycontainer",
    "resources": {
      "cpu": {
        "shares": 1024,
        "quota": 100000,
        "period": 100000
      },
      "memory": {
        "limit": 536870912,
        "swap": 536870912
      }
    }
  }
}
```

#### OCI Image Spec

```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
  "config": {
    "mediaType": "application/vnd.docker.container.image.v1+json",
    "size": 7023,
    "digest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  },
  "layers": [
    {
      "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
      "size": 32654,
      "digest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    }
  ]
}
```

### 1.3 Linux 内核容器技术

#### Namespace 隔离机制

| Namespace | 系统调用参数 | 隔离内容 | Go 常量 |
|-----------|-------------|----------|---------|
| **PID** | `CLONE_NEWPID` | 进程 ID | `syscall.CLONE_NEWPID` |
| **Network** | `CLONE_NEWNET` | 网络设备、端口、路由 | `syscall.CLONE_NEWNET` |
| **IPC** | `CLONE_NEWIPC` | 信号量、消息队列、共享内存 | `syscall.CLONE_NEWIPC` |
| **Mount** | `CLONE_NEWNS` | 挂载点 | `syscall.CLONE_NEWNS` |
| **UTS** | `CLONE_NEWUTS` | 主机名、域名 | `syscall.CLONE_NEWUTS` |
| **User** | `CLONE_NEWUSER` | 用户/组 ID | `syscall.CLONE_NEWUSER` |
| **Cgroup** | `CLONE_NEWCGROUP` | Cgroup 根目录 | `syscall.CLONE_NEWCGROUP` |
| **Time** | `CLONE_NEWTIME` | 系统时间（Linux 5.6+） | `syscall.CLONE_NEWTIME` |

#### Cgroup v2 资源控制

```
/sys/fs/cgroup/
├── cgroup.controllers          # 可用的控制器
├── cgroup.subtree_control      # 启用的子控制器
├── cgroup.procs                # 根 cgroup 的进程
├── cgroup.max.depth            # 最大嵌套深度
├── memory.max                  # 内存硬限制
├── memory.high                 # 内存软限制
├── memory.current              # 当前内存使用
├── cpu.max                     # CPU 配额 (quota period)
├── cpu.weight                  # CPU 权重 (1-10000)
├── io.max                      # IO 限速
├── pids.max                    # 最大进程数
└── [container_id]/
    ├── cgroup.procs
    ├── memory.max
    ├── cpu.max
    └── ...
```

#### Union Filesystem 对比

| 特性 | AUFS | OverlayFS | OverlayFS2 | Btrfs | ZFS |
|------|------|-----------|------------|-------|-----|
| **内核要求** | 需补丁 | 3.18+ | 4.0+ | 内置 | 需模块 |
| **层数限制** | 127 | 2 (lower+upper) | 500+ | 无限制 | 无限制 |
| **性能** | 中等 | 好 | 优秀 | 优秀 | 优秀 |
| **Docker 默认** | 旧版本 | 是 | 推荐 | 可选 | 可选 |
| **Copy-on-Write** | 是 | 是 | 是 | 是 | 是 |

### 1.4 容器镜像结构

```
┌─────────────────────────────────────────────────────────────┐
│                    Image Manifest                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  config: {digest, size, mediaType}                  │   │
│  │  layers: [{digest, size, mediaType}, ...]           │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│  Config JSON  │    │   Layer 1     │    │   Layer N     │
│  (metadata)   │    │  (base image) │    │  (app layer)  │
│               │    │               │    │               │
│ - architecture│    │ - filesystem  │    │ - filesystem  │
│ - os          │    │   changes     │    │   changes     │
│ - config.env  │    │               │    │               │
│ - config.cmd  │    │               │    │               │
│ - rootfs      │    │               │    │               │
│ - history     │    │               │    │               │
└───────────────┘    └───────────────┘    └───────────────┘
```

### 1.5 Docker 网络模型

```
┌─────────────────────────────────────────────────────────────┐
│                     Docker Network Types                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    bridge    │  │     host     │  │     none     │      │
│  │   (默认)      │  │              │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   overlay    │  │   macvlan    │  │    ipvlan    │      │
│  │  (Swarm/K8s) │  │  (物理网络)   │  │  (L2/L3)     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Bridge 网络数据包流程

```
Container A (172.17.0.2)          Container B (172.17.0.3)
        │                                  │
        │ veth0                    veth0   │
        └────┬────────────────────────┬────┘
             │      docker0 bridge   │
             │     (172.17.0.1/16)   │
             └──────────┬────────────┘
                        │
              ┌─────────▼──────────┐
              │   iptables NAT     │
              │  (MASQUERADE)      │
              └─────────┬──────────┘
                        │
              ┌─────────▼──────────┐
              │   Host eth0        │
              │  (外部网络)         │
              └────────────────────┘
```

---

## 2. Kubernetes 核心架构

### 2.1 整体架构概览

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           Kubernetes Cluster                             │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                        Control Plane                             │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │    │
│  │  │  API Server │  │    etcd     │  │      Scheduler          │  │    │
│  │  │  (kube-apiserver) │  │  (Key-Value)│  │  (kube-scheduler)       │  │    │
│  │  │             │  │             │  │                         │  │    │
│  │  │ • REST API  │  │ • Cluster   │  │ • Predicates            │  │    │
│  │  │ • AuthN/Z   │  │   State     │  │ • Priorities            │  │    │
│  │  │ • Admission │  │ • Watch API │  │ • Binding               │  │    │
│  │  │   Webhooks  │  │             │  │                         │  │    │
│  │  └─────────────┘  └─────────────┘  └─────────────────────────┘  │    │
│  │                                                                  │    │
│  │  ┌─────────────────────────────────────────────────────────┐    │    │
│  │  │              Controller Manager                          │    │    │
│  │  │  (kube-controller-manager)                               │    │    │
│  │  │  • Node Controller    • Deployment Controller            │    │    │
│  │  │  • ReplicaSet         • StatefulSet Controller           │    │    │
│  │  │  • EndpointSlice      • Job/CronJob Controller           │    │    │
│  │  │  • Service Account    • Namespace Controller             │    │    │
│  │  └─────────────────────────────────────────────────────────┘    │    │
│  │                                                                  │    │
│  │  ┌─────────────────────────────────────────────────────────┐    │    │
│  │  │              Cloud Controller Manager                    │    │    │
│  │  │  • Node Lifecycle     • Route Controller                 │    │    │
│  │  │  • Service Controller • Volume Controller                │    │    │
│  │  └─────────────────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                         Worker Nodes                             │    │
│  │                                                                  │    │
│  │  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────┐ │    │
│  │  │    Kubelet      │    │   Kube-proxy    │    │   Runtime   │ │    │
│  │  │                 │    │                 │    │             │ │    │
│  │  │ • Pod Lifecycle │    │ • Service Proxy │    │ • containerd│ │    │
│  │  │ • CRI Client    │    │ • iptables/IPVS │    │ • CRI-O     │ │    │
│  │  │ • CNI Client    │    │ • eBPF (可选)   │    │ • gVisor    │ │    │
│  │  │ • CSI Client    │    │                 │    │ • Kata      │ │    │
│  │  │ • Health Check  │    │                 │    │             │ │    │
│  │  └─────────────────┘    └─────────────────┘    └─────────────┘ │    │
│  │                                                                  │    │
│  │  ┌─────────────────────────────────────────────────────────┐    │    │
│  │  │                      Pod(s)                              │    │    │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐    │    │    │
│  │  │  │Container│  │Container│  │Container│  │Container│    │    │    │
│  │  │  │   1     │  │   2     │  │   3     │  │   N     │    │    │    │
│  │  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘    │    │    │
│  │  │  • Shared Network Namespace                              │    │    │
│  │  │  • Shared Storage Volumes                                │    │    │
│  │  │  • Shared IPC Namespace                                  │    │    │
│  │  └─────────────────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      Add-ons                                     │    │
│  │  • DNS (CoreDNS)  • Ingress Controller  • CNI Plugin            │    │
│  │  • Metrics Server • Dashboard           • CSI Driver            │    │
│  └─────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 核心组件详解

#### API Server (kube-apiserver)

```go
// API Server 核心功能
// 1. 暴露 RESTful API
// 2. 认证与授权
// 3. 准入控制
// 4. 数据验证
// 5. etcd 数据交互

type APIServer struct {
    // 认证链
    Authenticator authenticator.Request

    // 授权器
    Authorizer authorizer.Authorizer

    // 准入控制器
    AdmissionChain admission.Interface

    // REST 存储
    RESTStorageProvider RESTStorageProvider

    // etcd 客户端
    EtcdClient *etcdclient.Client
}

// 请求处理流程
// 1. HTTP Request -> Authentication -> Authorization -> Admission -> Validation -> etcd
// 2. Watch 机制: Long-polling / WebSocket / HTTP/2 Server Push
```

#### etcd - 集群状态存储

```
etcd 数据存储结构:
/registry/
├── apiextensions.k8s.io/
│   └── customresourcedefinitions/
├── apps/
│   ├── daemonsets/
│   ├── deployments/
│   ├── replicasets/
│   └── statefulsets/
├── batch/
│   ├── cronjobs/
│   └── jobs/
├── core/
│   ├── configmaps/
│   ├── endpoints/
│   ├── namespaces/
│   ├── nodes/
│   ├── persistentvolumeclaims/
│   ├── persistentvolumes/
│   ├── pods/
│   ├── secrets/
│   ├── serviceaccounts/
│   └── services/
├── networking.k8s.io/
│   ├── ingresses/
│   ├── ingressclasses/
│   └── networkpolicies/
├── rbac.authorization.k8s.io/
│   ├── clusterrolebindings/
│   ├── clusterroles/
│   ├── rolebindings/
│   └── roles/
└── storage.k8s.io/
    ├── csidrivers/
    ├── csinodes/
    ├── storageclasses/
    └── volumeattachments/
```

#### Scheduler 调度器

```
调度流程:
┌─────────────────────────────────────────────────────────────┐
│  1. Predicates (预选) - 排除不符合条件的节点                  │
├─────────────────────────────────────────────────────────────┤
│  • PodFitsResources      - 资源充足性检查                    │
│  • PodFitsHost           - 节点名称匹配                      │
│  • PodFitsHostPorts      - 端口冲突检查                      │
│  • PodMatchNodeSelector  - 节点选择器匹配                    │
│  • NoDiskConflict        - 磁盘冲突检查                      │
│  • NoVolumeZoneConflict  - 存储区域匹配                      │
│  • CheckNodeMemoryPressure - 内存压力检查                    │
│  • CheckNodeDiskPressure   - 磁盘压力检查                    │
│  • CheckNodeCondition    - 节点状态检查                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Priorities (优选) - 为节点打分                           │
├─────────────────────────────────────────────────────────────┤
│  • LeastRequested        - 资源利用率低优先 (权重 1)          │
│  • BalancedResourceAllocation - 资源均衡分配 (权重 1)         │
│  • SelectorSpreadPriority - 选择器分散 (权重 1)              │
│  • InterPodAffinityPriority - Pod 亲和性 (权重 1)            │
│  • NodeAffinityPriority  - 节点亲和性 (权重 1)               │
│  • TaintTolerationPriority - 容忍度匹配 (权重 1)             │
│  • ImageLocalityPriority - 镜像本地存在 (权重 1)             │
│  • ServiceSpreadingPriority - Service 分散 (权重 1)          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  3. Binding - 绑定 Pod 到选定节点                            │
│     创建 Binding 对象，由 API Server 写入 etcd               │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 核心对象关系

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Kubernetes Object Hierarchy                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Deployment (声明式应用管理)                                          │
│      │                                                               │
│      ├── manages ──► ReplicaSet (副本控制)                           │
│      │                  │                                            │
│      │                  ├── creates ──► Pod (最小调度单元)            │
│      │                  │                  │                         │
│      │                  │                  ├── contains ──► Container│
│      │                  │                  │                         │
│      │                  │                  ├── mounts ──► Volume     │
│      │                  │                  │                         │
│      │                  │                  └── uses ──► ConfigMap/Secret
│      │                  │                                            │
│      └── exposes ──► Service (服务发现)                              │
│                          │                                           │
│                          ├── routes ──► EndpointSlice (端点列表)     │
│                          │                                           │
│                          └── exposed ──► Ingress (外部访问)          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 2.4 控制器模式详解

| 控制器 | 用途 | 关键特性 |
|--------|------|----------|
| **Deployment** | 无状态应用 | 滚动更新、回滚、扩缩容 |
| **StatefulSet** | 有状态应用 | 稳定网络标识、有序部署、持久存储 |
| **DaemonSet** | 节点级服务 | 每节点一个 Pod、自动扩缩 |
| **Job** | 一次性任务 | 完成指定次数后停止 |
| **CronJob** | 定时任务 | 基于 cron 表达式调度 |
| **ReplicaSet** | 副本控制 | 维护指定数量的 Pod 副本 |

---

## 3. 接口标准：CRI/CNI/CSI

### 3.1 CRI (Container Runtime Interface)

```protobuf
// CRI Protocol Buffer 定义 (cri-api)
syntax = "proto3";

service RuntimeService {
    // Pod Sandbox 管理
    rpc RunPodSandbox(RunPodSandboxRequest) returns (RunPodSandboxResponse);
    rpc StopPodSandbox(StopPodSandboxRequest) returns (StopPodSandboxResponse);
    rpc RemovePodSandbox(RemovePodSandboxRequest) returns (RemovePodSandboxResponse);
    rpc PodSandboxStatus(PodSandboxStatusRequest) returns (PodSandboxStatusResponse);
    rpc ListPodSandbox(ListPodSandboxRequest) returns (ListPodSandboxResponse);

    // Container 管理
    rpc CreateContainer(CreateContainerRequest) returns (CreateContainerResponse);
    rpc StartContainer(StartContainerRequest) returns (StartContainerResponse);
    rpc StopContainer(StopContainerRequest) returns (StopContainerResponse);
    rpc RemoveContainer(RemoveContainerRequest) returns (RemoveContainerResponse);
    rpc ListContainers(ListContainersRequest) returns (ListContainersResponse);
    rpc ContainerStatus(ContainerStatusRequest) returns (ContainerStatusResponse);

    // 其他接口...
    rpc ExecSync(ExecSyncRequest) returns (ExecSyncResponse);
    rpc Exec(ExecRequest) returns (ExecResponse);
    rpc Attach(AttachRequest) returns (AttachResponse);
    rpc PortForward(PortForwardRequest) returns (PortForwardResponse);
}

service ImageService {
    rpc ListImages(ListImagesRequest) returns (ListImagesResponse);
    rpc ImageStatus(ImageStatusRequest) returns (ImageStatusResponse);
    rpc PullImage(PullImageRequest) returns (PullImageResponse);
    rpc RemoveImage(RemoveImageRequest) returns (RemoveImageResponse);
    rpc ImageFsInfo(ImageFsInfoRequest) returns (ImageFsInfoResponse);
}
```

#### CRI 运行时对比

| 特性 | containerd | CRI-O | Docker (已弃用) |
|------|-----------|-------|-----------------|
| **设计目标** | 通用容器运行时 | 专为 K8s 设计 | 通用容器平台 |
| **架构** | daemon + shim | 轻量级 daemon | 完整平台 |
| **镜像格式** | OCI + Docker | OCI | Docker |
| **性能** | 优秀 | 优秀 | 较好 |
| **资源占用** | 低 | 最低 | 较高 |
| **安全特性** | 完整 | 完整 | 完整 |
| **维护状态** | 活跃 | 活跃 | K8s 1.24+ 移除 |

### 3.2 CNI (Container Network Interface)

```go
// CNI 接口规范
type CNI interface {
    // 添加网络接口到容器
    AddNetwork(ctx context.Context, net *NetworkConfig,
               rt *RuntimeConf) (types.Result, error)

    // 从容器删除网络接口
    DelNetwork(ctx context.Context, net *NetworkConfig,
               rt *RuntimeConf) error

    // 检查网络配置
    CheckNetwork(ctx context.Context, net *NetworkConfig,
                 rt *RuntimeConf) error

    // 获取网络状态
    GetNetworkCachedResult(net *NetworkConfig,
                          rt *RuntimeConf) (types.Result, error)
}

// CNI 配置文件格式
{
    "cniVersion": "1.0.0",
    "name": "mynet",
    "type": "bridge",
    "bridge": "cni0",
    "isGateway": true,
    "ipMasq": true,
    "ipam": {
        "type": "host-local",
        "subnet": "10.22.0.0/16",
        "routes": [
            { "dst": "0.0.0.0/0" }
        ]
    }
}
```

#### CNI 插件对比

| 特性 | Calico | Cilium | Flannel | Weave Net |
|------|--------|--------|---------|-----------|
| **数据平面** | eBPF/iptables | eBPF | VXLAN/UDP | VXLAN |
| **网络策略** | 完整支持 | 完整支持 | 需配合 | 完整支持 |
| **性能** | 高 | 极高 | 中等 | 中等 |
| **可观测性** | 好 | 优秀 | 基础 | 好 |
| **服务网格集成** | 部分 | 原生 | 无 | 部分 |
| **加密** | WireGuard | WireGuard/IPsec | 无 | NaCl |
| **多集群** | 支持 | 支持 | 需配置 | 支持 |
| **学习曲线** | 中等 | 较陡 | 简单 | 简单 |

### 3.3 CSI (Container Storage Interface)

```protobuf
// CSI Protocol Buffer 定义
syntax = "proto3";

service Identity {
    rpc GetPluginInfo(GetPluginInfoRequest) returns (GetPluginInfoResponse);
    rpc GetPluginCapabilities(GetPluginCapabilitiesRequest) returns (GetPluginCapabilitiesResponse);
    rpc Probe(ProbeRequest) returns (ProbeResponse);
}

service Controller {
    rpc CreateVolume(CreateVolumeRequest) returns (CreateVolumeResponse);
    rpc DeleteVolume(DeleteVolumeRequest) returns (DeleteVolumeResponse);
    rpc ControllerPublishVolume(ControllerPublishVolumeRequest) returns (ControllerPublishVolumeResponse);
    rpc ControllerUnpublishVolume(ControllerUnpublishVolumeRequest) returns (ControllerUnpublishVolumeResponse);
    rpc ValidateVolumeCapabilities(ValidateVolumeCapabilitiesRequest) returns (ValidateVolumeCapabilitiesResponse);
    rpc ListVolumes(ListVolumesRequest) returns (ListVolumesResponse);
    rpc GetCapacity(GetCapacityRequest) returns (GetCapacityResponse);
    rpc ControllerGetCapabilities(ControllerGetCapabilitiesRequest) returns (ControllerGetCapabilitiesResponse);
    rpc CreateSnapshot(CreateSnapshotRequest) returns (CreateSnapshotResponse);
    rpc DeleteSnapshot(DeleteSnapshotRequest) returns (DeleteSnapshotResponse);
    rpc ListSnapshots(ListSnapshotsRequest) returns (ListSnapshotsResponse);
    rpc ControllerExpandVolume(ControllerExpandVolumeRequest) returns (ControllerExpandVolumeResponse);
}

service Node {
    rpc NodeStageVolume(NodeStageVolumeRequest) returns (NodeStageVolumeResponse);
    rpc NodeUnstageVolume(NodeUnstageVolumeRequest) returns (NodeUnstageVolumeResponse);
    rpc NodePublishVolume(NodePublishVolumeRequest) returns (NodePublishVolumeResponse);
    rpc NodeUnpublishVolume(NodeUnpublishVolumeRequest) returns (NodeUnpublishVolumeResponse);
    rpc NodeGetVolumeStats(NodeGetVolumeStatsRequest) returns (NodeGetVolumeStatsResponse);
    rpc NodeExpandVolume(NodeExpandVolumeRequest) returns (NodeExpandVolumeResponse);
    rpc NodeGetCapabilities(NodeGetCapabilitiesRequest) returns (NodeGetCapabilitiesResponse);
    rpc NodeGetInfo(NodeGetInfoRequest) returns (NodeGetInfoResponse);
}
```

---

## 4. 网络与存储体系

### 4.1 Service 实现机制

```
┌─────────────────────────────────────────────────────────────────────┐
│                      Service 实现演进                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  userspace 模式 (已弃用)                                      │   │
│  │  kube-proxy ──► 用户空间代理 ──► iptables                    │   │
│  │  缺点: 性能差、用户态/内核态切换开销                           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                       │
│                              ▼                                       │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  iptables 模式 (默认)                                         │   │
│  │  kube-proxy ──► iptables rules ──► DNAT                     │   │
│  │  缺点: 规则数量 O(n²)，大规模集群性能下降                      │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                       │
│                              ▼                                       │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  IPVS 模式 (推荐用于大规模)                                    │   │
│  │  kube-proxy ──► IPVS rules ──► 内核负载均衡                  │   │
│  │  优点: O(1) 查找，支持多种调度算法                            │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                       │
│                              ▼                                       │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  eBPF 模式 (未来趋势)                                         │   │
│  │  Cilium/kube-proxy-replacement ──► eBPF programs            │   │
│  │  优点: 绕过 iptables，直接内核处理，性能最优                   │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 4.2 Ingress 控制器对比

| 特性 | NGINX Ingress | Traefik | HAProxy | Kong | Istio Gateway |
|------|---------------|---------|---------|------|---------------|
| **配置方式** | Ingress/CRD | Ingress/CRD | ConfigMap | CRD/Admin API | Gateway API |
| **动态重载** | 需 reload | 完全动态 | 需 reload | 完全动态 | 完全动态 |
| **性能** | 高 | 高 | 极高 | 高 | 高 |
| **SSL 终止** | 支持 | 支持 | 支持 | 支持 | 支持 |
| **mTLS** | 支持 | 支持 | 支持 | 支持 | 原生支持 |
| **金丝雀发布** | 需注解 | 原生支持 | 需配置 | 插件 | 原生支持 |
| **服务网格** | 需集成 | 需集成 | 需集成 | 需集成 | 原生 |
| **学习曲线** | 低 | 低 | 中等 | 中等 | 较陡 |

### 4.3 存储类型详解

```
┌─────────────────────────────────────────────────────────────────────┐
│                       Kubernetes Storage Types                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  临时存储 (Ephemeral)                                                 │
│  ├── emptyDir        - Pod 生命周期，节点本地存储                     │
│  ├── configMap       - 配置数据挂载                                   │
│  ├── secret          - 敏感数据挂载                                   │
│  └── downwardAPI     - Pod 元数据挂载                                 │
│                                                                      │
│  节点存储 (Node-local)                                                │
│  ├── hostPath        - 节点文件系统路径                               │
│  └── local           - 本地持久卷 (需调度约束)                         │
│                                                                      │
│  网络存储 (Network)                                                   │
│  ├── NFS             - 网络文件系统                                   │
│  ├── iSCSI           - 块存储协议                                     │
│  ├── Ceph/RBD        - 分布式块存储                                   │
│  ├── CephFS          - 分布式文件系统                                 │
│  ├── GlusterFS       - 分布式文件系统                                 │
│  └── AWS EBS/GCE PD  - 云提供商块存储                                 │
│                                                                      │
│  对象存储 (Object)                                                    │
│  └── S3/MinIO/GCS    - 通过 CSI 驱动或应用层访问                       │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 4.4 PV/PVC 生命周期

```
┌─────────────────────────────────────────────────────────────────────┐
│                    PV/PVC Lifecycle                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────┐       │
│  │  Admin  │────►│  Create │────►│   PV    │────►│ Available│      │
│  │         │     │   PV    │     │ (静态)   │     │          │      │
│  └─────────┘     └─────────┘     └─────────┘     └────┬────┘       │
│                                                       │              │
│  ┌─────────┐     ┌─────────┐     ┌─────────┐         │              │
│  │  User   │────►│  Create │────►│   PVC   │─────────┘              │
│  │         │     │   PVC   │     │         │                        │
│  └─────────┘     └─────────┘     └─────────┘                        │
│                                         │                            │
│                                         ▼                            │
│                              ┌─────────────────────┐                 │
│                              │   Dynamic Provisioning │               │
│                              │   (StorageClass)      │               │
│                              │   • Provisioner       │               │
│                              │   • Parameters        │               │
│                              │   • ReclaimPolicy     │               │
│                              │   • VolumeBindingMode │               │
│                              └─────────────────────┘                 │
│                                         │                            │
│                                         ▼                            │
│  ┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────┐       │
│  │  Pod    │────►│  Mount  │────►│  Bound  │◄────│   PV    │       │
│  │         │     │  Volume │     │  PVC    │     │ (动态)   │       │
│  └─────────┘     └─────────┘     └─────────┘     └─────────┘       │
│                                         │                            │
│                                         ▼                            │
│                              ┌─────────────────────┐                 │
│                              │   Reclaim Policy     │                 │
│                              │   • Retain (保留)    │                 │
│                              │   • Delete (删除)    │                 │
│                              │   • Recycle (回收)   │                 │
│                              └─────────────────────┘                 │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 5. 安全体系

### 5.1 认证机制

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Authentication Methods                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐     │
│  │   X.509 certs   │  │  ServiceAccount │  │    Webhook      │     │
│  │                 │  │                 │  │                 │     │
│  │ • Client certs  │  │ • JWT tokens    │  │ • External IdP  │     │
│  │ • CA signed     │  │ • Auto-mounted  │  │ • LDAP/OIDC     │     │
│  │ • kubeconfig    │  │ • In-cluster    │  │ • Custom auth   │     │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘     │
│                                                                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐     │
│  │   Token File    │  │    Bootstrap    │  │   Anonymous     │     │
│  │                 │  │    Tokens       │  │                 │     │
│  │ • Static tokens │  │ • kubeadm       │  │ • Disabled by   │     │
│  │ • CSV format    │  │ • Node join     │  │   default       │     │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.2 RBAC 授权模型

```yaml
# Role 定义 - 命名空间级别
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: production
  name: pod-reader
rules:
- apiGroups: [""]  # "" 表示 core API group
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  resourceNames: ["my-deployment"]  # 可选: 限定特定资源

---
# ClusterRole 定义 - 集群级别
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cluster-admin
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
- nonResourceURLs: ["/healthz", "/version"]
  verbs: ["get"]

---
# RoleBinding - 绑定 Role 到用户/组/ServiceAccount
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: production
subjects:
- kind: User
  name: alice
  apiGroup: rbac.authorization.k8s.io
- kind: ServiceAccount
  name: default
  namespace: production
- kind: Group
  name: developers
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

### 5.3 准入控制器

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Admission Controller Flow                         │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Request ──► Authentication ──► Authorization ──► Admission ──► etcd │
│                                                        │             │
│                                                        ▼             │
│                                           ┌─────────────────────┐    │
│                                           │   Mutating Webhooks  │    │
│                                           │   (修改请求)          │    │
│                                           │   • Defaulting       │    │
│                                           │   • Sidecar Injection│    │
│                                           └─────────────────────┘    │
│                                                        │             │
│                                                        ▼             │
│                                           ┌─────────────────────┐    │
│                                           │   Object Schema      │    │
│                                           │   Validation         │    │
│                                           └─────────────────────┘    │
│                                                        │             │
│                                                        ▼             │
│                                           ┌─────────────────────┐    │
│                                           │   Validating Webhooks│    │
│                                           │   (验证请求)          │    │
│                                           │   • Policy Check     │    │
│                                           │   • Custom Logic     │    │
│                                           └─────────────────────┘    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.4 Pod 安全标准

| 级别 | 描述 | 主要限制 |
|------|------|----------|
| **Privileged** | 无限制 | 允许所有配置 |
| **Baseline** | 最小限制 | 禁止特权容器、hostNetwork、hostPID 等 |
| **Restricted** | 严格限制 | 非 root、只读根文件系统、禁止特权升级等 |

```yaml
# Pod Security Standards - Restricted
apiVersion: v1
kind: Pod
metadata:
  name: restricted-pod
spec:
  securityContext:
    runAsNonRoot: true
    seccompProfile:
      type: RuntimeDefault
  containers:
  - name: app
    image: myapp:latest
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop: ["ALL"]
      runAsUser: 1000
      runAsGroup: 1000
    resources:
      limits:
        memory: "128Mi"
        cpu: "500m"
      requests:
        memory: "64Mi"
        cpu: "250m"
```

---

## 6. 生态工具链

### 6.1 包管理工具对比

| 特性 | Helm | Kustomize | Operator SDK |
|------|------|-----------|--------------|
| **设计哲学** | 模板化包管理 | 配置覆盖 | 自定义控制器 |
| **学习曲线** | 中等 | 低 | 较陡 |
| **模板引擎** | Go template | 无 (原生 YAML) | Go/Kubebuilder |
| **版本管理** | Chart 版本 | Git 版本 | Operator 版本 |
| **依赖管理** | 支持 | 不支持 | 不支持 |
| **适用场景** | 复杂应用部署 | 环境差异化配置 | 有状态应用管理 |

### 6.2 GitOps 工具对比

| 特性 | ArgoCD | Flux | Rancher Fleet |
|------|--------|------|---------------|
| **UI** | 丰富 Web UI | CLI 为主 | Rancher 集成 |
| **多集群** | 原生支持 | 支持 | 原生支持 |
| **应用类型** | 多种 | GitOps 原生 | 多种 |
| **同步策略** | 自动/手动 | 自动 | 自动 |
| **回滚** | 支持 | 支持 | 支持 |
| **通知** | 丰富 | 基础 | 基础 |
| **社区活跃度** | 高 | 高 | 中等 |

### 6.3 可观测性栈

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Observability Stack                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                      Metrics (指标)                          │   │
│  │  ┌───────────┐    ┌───────────┐    ┌───────────┐           │   │
│  │  │ Prometheus│───►│  Grafana  │    │ Alertmanager│          │   │
│  │  │           │    │           │    │            │           │   │
│  │  │ • Pull    │    │ • Dashboard│   │ • Routing  │           │   │
│  │  │ • TSDB    │    │ • Alerts  │    │ • Silencing│           │   │
│  │  │ • PromQL  │    │ • Explore │    │ • Inhibit  │           │   │
│  │  └───────────┘    └───────────┘    └───────────┘           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                      Logs (日志)                             │   │
│  │  ┌───────────┐    ┌───────────┐    ┌───────────┐           │   │
│  │  │ Fluentd/  │───►│Elasticsearch│──►│   Kibana  │           │   │
│  │  │ Fluent Bit│    │           │    │           │           │   │
│  │  │           │    │ • Search  │    │ • Visualize│          │   │
│  │  │ • Collect │    │ • Index   │    │ • Discover │          │   │
│  │  │ • Filter  │    │ • Scale   │    │ • Dashboard│          │   │
│  │  └───────────┘    └───────────┘    └───────────┘           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                      Traces (链路追踪)                       │   │
│  │  ┌───────────┐    ┌───────────┐    ┌───────────┐           │   │
│  │  │ OpenTelemetry│──►│   Jaeger  │    │   Tempo   │           │   │
│  │  │             │    │           │    │           │           │   │
│  │  │ • SDK     │    │ • UI      │    │ • Grafana │           │   │
│  │  │ • Collector│   │ • Storage │    │   Native  │           │   │
│  │  │ • OTLP    │    │ • Query   │    │           │           │   │
│  │  └───────────┘    └───────────┘    └───────────┘           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 6.4 服务网格对比

| 特性 | Istio | Linkerd | Consul Connect | AWS App Mesh |
|------|-------|---------|----------------|--------------|
| **数据平面** | Envoy | Linkerd-proxy | Envoy | Envoy |
| **控制平面** | istiod | linkerd-control-plane | Consul | AWS 托管 |
| **性能开销** | 中等 | 低 | 中等 | 中等 |
| **mTLS** | 支持 | 支持 | 支持 | 支持 |
| **流量管理** | 丰富 | 基础 | 中等 | 中等 |
| **可观测性** | 优秀 | 好 | 好 | 基础 |
| **多集群** | 支持 | 支持 | 支持 | 部分支持 |
| **VM 支持** | 支持 | 支持 | 原生 | AWS 原生 |
| **学习曲线** | 较陡 | 低 | 中等 | 低 |

### 6.5 Serverless 平台对比

| 特性 | Knative | OpenFaaS | Kubeless | Fission |
|------|---------|----------|----------|---------|
| **构建** | 内置 (Buildpacks) | 多种方式 | 内置 | 内置 |
| **自动扩缩** | 优秀 (KPA/HPA) | 基础 | 基础 | 基础 |
| **冷启动** | 低 (支持预热) | 中等 | 中等 | 低 |
| **事件驱动** | 原生支持 | 支持 | 支持 | 支持 |
| **服务网格** | 可选集成 | 可选 | 无 | 无 |
| **社区** | 活跃 (Google/Red Hat) | 活跃 | 较低 | 中等 |

---

## 7. WebAssembly 与容器融合

### 7.1 WebAssembly 运行时对比

| 特性 | WasmEdge | Wasmtime | WAMR | Wasmer |
|------|----------|----------|------|--------|
| **语言** | C++/Rust | Rust | C | Rust |
| **性能** | 高 | 高 | 中等 | 高 |
| **AOT 编译** | 支持 | 支持 | 支持 | 支持 |
| **WASI** | 完整 | 完整 | 完整 | 完整 |
| **扩展插件** | 丰富 | 中等 | 基础 | 中等 |
| **云原生集成** | 优秀 | 好 | 基础 | 好 |
| **容器支持** | containerd-shim | 实验性 | 无 | 实验性 |

### 7.2 Wasm 与容器对比

```
┌─────────────────────────────────────────────────────────────────────┐
│                  Wasm vs Container Comparison                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────┐        ┌─────────────────────┐            │
│  │      Container      │        │    WebAssembly      │            │
│  │                     │        │                     │            │
│  │  ┌───────────────┐  │        │  ┌───────────────┐  │            │
│  │  │   App Code    │  │        │  │   Wasm Module │  │            │
│  │  │   + Runtime   │  │        │  │   (.wasm)     │  │            │
│  │  └───────┬───────┘  │        │  └───────┬───────┘  │            │
│  │          │          │        │          │          │            │
│  │  ┌───────▼───────┐  │        │  ┌───────▼───────┐  │            │
│  │  │  Guest OS     │  │        │  │  Wasm Runtime │  │            │
│  │  │  (Libraries)  │  │        │  │  (沙箱)        │  │            │
│  │  └───────┬───────┘  │        │  └───────┬───────┘  │            │
│  │          │          │        │          │          │            │
│  │  ┌───────▼───────┐  │        │  ┌───────▼───────┐  │            │
│  │  │  Host OS      │  │        │  │  Host OS      │  │            │
│  │  │  (Kernel)     │  │        │  │  (Kernel)     │  │            │
│  │  └───────────────┘  │        │  └───────────────┘  │            │
│  └─────────────────────┘        └─────────────────────┘            │
│                                                                      │
│  特性对比:                                                            │
│  ┌──────────────┬─────────────────┬─────────────────┐               │
│  │    特性      │     容器        │    WebAssembly  │               │
│  ├──────────────┼─────────────────┼─────────────────┤               │
│  │ 启动时间     │ 秒级            │ 毫秒级          │               │
│  │ 镜像大小     │ MB-GB           │ KB-MB           │               │
│  │ 隔离级别     │ OS 级           │ 沙箱级          │               │
│  │ 安全边界     │ Namespace/Cgroup│ Capability-based│               │
│  │ 可移植性     │ 架构相关        │ 架构无关        │               │
│  │ 语言支持     │ 任意            │ 编译到 Wasm     │               │
│  │ 生态成熟度   │ 成熟            │ 快速发展        │               │
│  └──────────────┴─────────────────┴─────────────────┘               │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 7.3 containerd-wasm-shim 架构

```
┌─────────────────────────────────────────────────────────────────────┐
│              containerd + WebAssembly Integration                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                     containerd                               │   │
│  │  ┌─────────────────────────────────────────────────────┐    │   │
│  │  │              containerd-shim-v2                      │    │   │
│  │  │  (RuntimeClass: wasm)                                │    │   │
│  │  └────────────────────┬────────────────────────────────┘    │   │
│  │                       │                                      │   │
│  │  ┌────────────────────▼────────────────────────────────┐    │   │
│  │  │              containerd-wasm-shim                    │    │   │
│  │  │  • Spin (Fermyon)                                    │    │   │
│  │  │  • Slight (Deis Labs)                                │    │   │
│  │  │  • WasmEdge (CNCF)                                   │    │   │
│  │  │  • Wasmtime (Bytecode Alliance)                      │    │   │
│  │  └────────────────────┬────────────────────────────────┘    │   │
│  └───────────────────────┼──────────────────────────────────────┘   │
│                          │                                           │
│  ┌───────────────────────▼──────────────────────────────────────┐   │
│  │                    Wasm Runtime                               │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │   │
│  │  │  WasmEdge   │  │  Wasmtime   │  │  WAMR               │  │   │
│  │  │             │  │             │  │                     │  │   │
│  │  │ • WASI      │  │ • WASI      │  │ • WASI              │  │   │
│  │  │ • AOT/JIT   │  │ • Cranelift │  │ • Interpreter       │  │   │
│  │  │ • Plugins   │  │ • Wasmtime  │  │ • AoT               │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 8. Go 代码实战

### 8.1 使用 client-go 操作 K8s 资源

```go
package main

import (
 "context"
 "fmt"
 "path/filepath"

 corev1 "k8s.io/api/core/v1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/client-go/kubernetes"
 "k8s.io/client-go/tools/clientcmd"
 "k8s.io/client-go/util/homedir"
)

// K8sClient 封装 Kubernetes 客户端操作
type K8sClient struct {
 clientset *kubernetes.Clientset
}

// NewK8sClient 创建新的 K8s 客户端
func NewK8sClient(kubeconfigPath string) (*K8sClient, error) {
 // 如果没有提供 kubeconfig 路径，使用默认路径
 if kubeconfigPath == "" {
  if home := homedir.HomeDir(); home != "" {
   kubeconfigPath = filepath.Join(home, ".kube", "config")
  }
 }

 // 构建配置
 config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
 if err != nil {
  // 尝试使用 in-cluster 配置
  config, err = rest.InClusterConfig()
  if err != nil {
   return nil, fmt.Errorf("failed to create config: %w", err)
  }
 }

 // 创建 clientset
 clientset, err := kubernetes.NewForConfig(config)
 if err != nil {
  return nil, fmt.Errorf("failed to create clientset: %w", err)
 }

 return &K8sClient{clientset: clientset}, nil
}

// CreateNamespace 创建命名空间
func (c *K8sClient) CreateNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
 ns := &corev1.Namespace{
  ObjectMeta: metav1.ObjectMeta{
   Name: name,
   Labels: map[string]string{
    "managed-by": "go-client",
   },
  },
 }
 return c.clientset.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
}

// CreateDeployment 创建 Deployment
func (c *K8sClient) CreateDeployment(ctx context.Context, namespace, name string, replicas int32) error {
 deployment := &appsv1.Deployment{
  ObjectMeta: metav1.ObjectMeta{
   Name:      name,
   Namespace: namespace,
   Labels: map[string]string{
    "app": name,
   },
  },
  Spec: appsv1.DeploymentSpec{
   Replicas: &replicas,
   Selector: &metav1.LabelSelector{
    MatchLabels: map[string]string{
     "app": name,
    },
   },
   Template: corev1.PodTemplateSpec{
    ObjectMeta: metav1.ObjectMeta{
     Labels: map[string]string{
      "app": name,
     },
    },
    Spec: corev1.PodSpec{
     Containers: []corev1.Container{
      {
       Name:  "app",
       Image: "nginx:latest",
       Ports: []corev1.ContainerPort{
        {
         ContainerPort: 80,
        },
       },
       Resources: corev1.ResourceRequirements{
        Requests: corev1.ResourceList{
         corev1.ResourceCPU:    resource.MustParse("100m"),
         corev1.ResourceMemory: resource.MustParse("128Mi"),
        },
        Limits: corev1.ResourceList{
         corev1.ResourceCPU:    resource.MustParse("500m"),
         corev1.ResourceMemory: resource.MustParse("512Mi"),
        },
       },
      },
     },
    },
   },
  },
 }

 _, err := c.clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
 return err
}

// WatchPods 监视 Pod 变化
func (c *K8sClient) WatchPods(ctx context.Context, namespace string) error {
 watchInterface, err := c.clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{})
 if err != nil {
  return err
 }
 defer watchInterface.Stop()

 fmt.Printf("Watching pods in namespace: %s\n", namespace)
 for event := range watchInterface.ResultChan() {
  pod, ok := event.Object.(*corev1.Pod)
  if !ok {
   continue
  }
  fmt.Printf("Event: %s, Pod: %s/%s, Phase: %s\n",
   event.Type, pod.Namespace, pod.Name, pod.Status.Phase)
 }
 return nil
}

// ListPodsWithFieldSelector 使用字段选择器列出 Pod
func (c *K8sClient) ListPodsWithFieldSelector(ctx context.Context, namespace, fieldSelector string) (*corev1.PodList, error) {
 return c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
  FieldSelector: fieldSelector,
 })
}

// PatchDeployment 使用 Strategic Merge Patch 更新 Deployment
func (c *K8sClient) PatchDeployment(ctx context.Context, namespace, name string, replicas int32) error {
 patchData := fmt.Sprintf(`{"spec":{"replicas":%d}}`, replicas)
 _, err := c.clientset.AppsV1().Deployments(namespace).Patch(
  ctx, name, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{})
 return err
}

func main() {
 client, err := NewK8sClient("")
 if err != nil {
  panic(err)
 }

 ctx := context.Background()

 // 创建命名空间
 ns, err := client.CreateNamespace(ctx, "demo-namespace")
 if err != nil {
  fmt.Printf("Error creating namespace: %v\n", err)
 } else {
  fmt.Printf("Created namespace: %s\n", ns.Name)
 }

 // 创建 Deployment
 if err := client.CreateDeployment(ctx, "demo-namespace", "nginx-deployment", 3); err != nil {
  fmt.Printf("Error creating deployment: %v\n", err)
 } else {
  fmt.Println("Created deployment: nginx-deployment")
 }
}
```

### 8.2 自定义 Controller 实现

```go
package main

import (
 "context"
 "fmt"
 "time"

 "k8s.io/apimachinery/pkg/api/errors"
 "k8s.io/apimachinery/pkg/runtime"
 "k8s.io/client-go/tools/cache"
 "k8s.io/client-go/util/workqueue"
 "sigs.k8s.io/controller-runtime/pkg/controller"
 "sigs.k8s.io/controller-runtime/pkg/handler"
 "sigs.k8s.io/controller-runtime/pkg/manager"
 "sigs.k8s.io/controller-runtime/pkg/reconcile"
 "sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler 实现 reconcile.Reconciler 接口
type PodReconciler struct {
 client client.Client
 scheme *runtime.Scheme
}

// Reconcile 处理资源调和
func (r *PodReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
 // 获取 Pod
 pod := &corev1.Pod{}
 err := r.client.Get(ctx, req.NamespacedName, pod)
 if err != nil {
  if errors.IsNotFound(err) {
   // Pod 已被删除
   fmt.Printf("Pod %s/%s not found, may be deleted\n", req.Namespace, req.Name)
   return reconcile.Result{}, nil
  }
  return reconcile.Result{}, err
 }

 // 业务逻辑：例如添加标签
 if pod.Labels == nil {
  pod.Labels = make(map[string]string)
 }

 if _, exists := pod.Labels["managed-by"]; !exists {
  pod.Labels["managed-by"] = "custom-controller"
  if err := r.client.Update(ctx, pod); err != nil {
   return reconcile.Result{}, err
  }
  fmt.Printf("Updated pod %s/%s with label\n", req.Namespace, req.Name)
 }

 return reconcile.Result{}, nil
}

// 使用 client-go informer 实现的自定义控制器
type CustomController struct {
 informer  cache.SharedIndexInformer
 queue     workqueue.RateLimitingInterface
 workercnt int
}

// NewCustomController 创建自定义控制器
func NewCustomController(informer cache.SharedIndexInformer, workercnt int) *CustomController {
 queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

 // 添加事件处理器
 informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
  AddFunc: func(obj interface{}) {
   key, err := cache.MetaNamespaceKeyFunc(obj)
   if err == nil {
    queue.Add(key)
   }
  },
  UpdateFunc: func(oldObj, newObj interface{}) {
   key, err := cache.MetaNamespaceKeyFunc(newObj)
   if err == nil {
    queue.Add(key)
   }
  },
  DeleteFunc: func(obj interface{}) {
   key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
   if err == nil {
    queue.Add(key)
   }
  },
 })

 return &CustomController{
  informer:  informer,
  queue:     queue,
  workercnt: workercnt,
 }
}

// Run 启动控制器
func (c *CustomController) Run(stopCh <-chan struct{}) {
 defer c.queue.ShutDown()

 // 启动 informer
 go c.informer.Run(stopCh)

 // 等待缓存同步
 if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
  fmt.Println("Failed to sync cache")
  return
 }

 // 启动工作协程
 for i := 0; i < c.workercnt; i++ {
  go c.worker()
 }

 <-stopCh
}

// worker 处理队列中的项目
func (c *CustomController) worker() {
 for c.processNextItem() {
 }
}

// processNextItem 处理队列中的下一个项目
func (c *CustomController) processNextItem() bool {
 key, quit := c.queue.Get()
 if quit {
  return false
 }
 defer c.queue.Done(key)

 err := c.syncHandler(key.(string))
 if err != nil {
  c.queue.AddRateLimited(key)
  fmt.Printf("Error syncing %s: %v\n", key, err)
 } else {
  c.queue.Forget(key)
 }

 return true
}

// syncHandler 实际的业务逻辑
func (c *CustomController) syncHandler(key string) error {
 namespace, name, err := cache.SplitMetaNamespaceKey(key)
 if err != nil {
  return err
 }

 obj, exists, err := c.informer.GetIndexer().GetByKey(key)
 if err != nil {
  return err
 }

 if !exists {
  fmt.Printf("Resource %s/%s deleted\n", namespace, name)
  return nil
 }

 pod, ok := obj.(*corev1.Pod)
 if !ok {
  return fmt.Errorf("unexpected type")
 }

 fmt.Printf("Processing pod: %s/%s, Phase: %s\n",
  pod.Namespace, pod.Name, pod.Status.Phase)

 // 在这里添加自定义业务逻辑
 return nil
}

// 使用 controller-runtime 的高级实现
func setupManager() (manager.Manager, error) {
 // 创建 manager
 mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
  MetricsBindAddress: ":8080",
  LeaderElection:     true,
  LeaderElectionID:   "custom-controller.example.com",
 })
 if err != nil {
  return nil, err
 }

 // 创建 reconciler
 reconciler := &PodReconciler{
  client: mgr.GetClient(),
  scheme: mgr.GetScheme(),
 }

 // 创建控制器
 ctrl, err := controller.New("pod-controller", mgr, controller.Options{
  Reconciler: reconciler,
 })
 if err != nil {
  return nil, err
 }

 // 设置监听器
 if err := ctrl.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}); err != nil {
  return nil, err
 }

 return mgr, nil
}
```

### 8.3 CRI 客户端调用

```go
package main

import (
 "context"
 "fmt"
 "time"

 "google.golang.org/grpc"
 "google.golang.org/grpc/credentials/insecure"

 runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// CRIClient CRI 客户端封装
type CRIClient struct {
 runtimeClient runtimeapi.RuntimeServiceClient
 imageClient   runtimeapi.ImageServiceClient
 conn          *grpc.ClientConn
}

// NewCRIClient 创建 CRI 客户端
func NewCRIClient(endpoint string) (*CRIClient, error) {
 conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
 if err != nil {
  return nil, fmt.Errorf("failed to connect to CRI endpoint: %w", err)
 }

 return &CRIClient{
  runtimeClient: runtimeapi.NewRuntimeServiceClient(conn),
  imageClient:   runtimeapi.NewImageServiceClient(conn),
  conn:          conn,
 }, nil
}

// Close 关闭连接
func (c *CRIClient) Close() error {
 return c.conn.Close()
}

// GetRuntimeVersion 获取运行时版本
func (c *CRIClient) GetRuntimeVersion(ctx context.Context) (*runtimeapi.VersionResponse, error) {
 return c.runtimeClient.Version(ctx, &runtimeapi.VersionRequest{})
}

// ListContainers 列出容器
func (c *CRIClient) ListContainers(ctx context.Context, filter *runtimeapi.ContainerFilter) ([]*runtimeapi.Container, error) {
 resp, err := c.runtimeClient.ListContainers(ctx, &runtimeapi.ListContainersRequest{
  Filter: filter,
 })
 if err != nil {
  return nil, err
 }
 return resp.Containers, nil
}

// ListPods 列出 Pod Sandbox
func (c *CRIClient) ListPods(ctx context.Context, filter *runtimeapi.PodSandboxFilter) ([]*runtimeapi.PodSandbox, error) {
 resp, err := c.runtimeClient.ListPodSandbox(ctx, &runtimeapi.ListPodSandboxRequest{
  Filter: filter,
 })
 if err != nil {
  return nil, err
 }
 return resp.Items, nil
}

// RunPod 创建并运行 Pod Sandbox
func (c *CRIClient) RunPod(ctx context.Context, config *runtimeapi.PodSandboxConfig) (string, error) {
 resp, err := c.runtimeClient.RunPodSandbox(ctx, &runtimeapi.RunPodSandboxRequest{
  Config: config,
 })
 if err != nil {
  return "", err
 }
 return resp.PodSandboxId, nil
}

// CreateContainer 在 Pod 中创建容器
func (c *CRIClient) CreateContainer(ctx context.Context, podID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
 resp, err := c.runtimeClient.CreateContainer(ctx, &runtimeapi.CreateContainerRequest{
  PodSandboxId:  podID,
  Config:        config,
  SandboxConfig: sandboxConfig,
 })
 if err != nil {
  return "", err
 }
 return resp.ContainerId, nil
}

// StartContainer 启动容器
func (c *CRIClient) StartContainer(ctx context.Context, containerID string) error {
 _, err := c.runtimeClient.StartContainer(ctx, &runtimeapi.StartContainerRequest{
  ContainerId: containerID,
 })
 return err
}

// StopContainer 停止容器
func (c *CRIClient) StopContainer(ctx context.Context, containerID string, timeout int64) error {
 _, err := c.runtimeClient.StopContainer(ctx, &runtimeapi.StopContainerRequest{
  ContainerId: containerID,
  Timeout:     timeout,
 })
 return err
}

// RemoveContainer 删除容器
func (c *CRIClient) RemoveContainer(ctx context.Context, containerID string) error {
 _, err := c.runtimeClient.RemoveContainer(ctx, &runtimeapi.RemoveContainerRequest{
  ContainerId: containerID,
 })
 return err
}

// PullImage 拉取镜像
func (c *CRIClient) PullImage(ctx context.Context, imageRef string, auth *runtimeapi.AuthConfig) (string, error) {
 resp, err := c.imageClient.PullImage(ctx, &runtimeapi.PullImageRequest{
  Image: &runtimeapi.ImageSpec{
   Image: imageRef,
  },
  Auth: auth,
 })
 if err != nil {
  return "", err
 }
 return resp.ImageRef, nil
}

// ListImages 列出镜像
func (c *CRIClient) ListImages(ctx context.Context, filter *runtimeapi.ImageFilter) ([]*runtimeapi.Image, error) {
 resp, err := c.imageClient.ListImages(ctx, &runtimeapi.ListImagesRequest{
  Filter: filter,
 })
 if err != nil {
  return nil, err
 }
 return resp.Images, nil
}

// ExecSync 在容器中同步执行命令
func (c *CRIClient) ExecSync(ctx context.Context, containerID string, cmd []string, timeout time.Duration) (*runtimeapi.ExecSyncResponse, error) {
 return c.runtimeClient.ExecSync(ctx, &runtimeapi.ExecSyncRequest{
  ContainerId: containerID,
  Cmd:         cmd,
  Timeout:     int64(timeout.Seconds()),
 })
}

// GetContainerStats 获取容器统计信息
func (c *CRIClient) GetContainerStats(ctx context.Context, containerID string) (*runtimeapi.ContainerStats, error) {
 resp, err := c.runtimeClient.ContainerStats(ctx, &runtimeapi.ContainerStatsRequest{
  ContainerId: containerID,
 })
 if err != nil {
  return nil, err
 }
 return resp.Stats, nil
}

// 使用示例
func main() {
 // 默认 containerd CRI 端点
 client, err := NewCRIClient("unix:///run/containerd/containerd.sock")
 if err != nil {
  panic(err)
 }
 defer client.Close()

 ctx := context.Background()

 // 获取运行时版本
 version, err := client.GetRuntimeVersion(ctx)
 if err != nil {
  panic(err)
 }
 fmt.Printf("Runtime: %s, Version: %s, Runtime API: %s\n",
  version.RuntimeName, version.RuntimeVersion, version.RuntimeApiVersion)

 // 列出所有容器
 containers, err := client.ListContainers(ctx, nil)
 if err != nil {
  panic(err)
 }
 fmt.Printf("Found %d containers\n", len(containers))
 for _, c := range containers {
  fmt.Printf("  - %s: %s (%s)\n", c.Id[:12], c.Metadata.Name, c.State)
 }

 // 列出所有 Pod
 pods, err := client.ListPods(ctx, nil)
 if err != nil {
  panic(err)
 }
 fmt.Printf("Found %d pods\n", len(pods))
 for _, p := range pods {
  fmt.Printf("  - %s: %s (%s)\n", p.Id[:12], p.Metadata.Name, p.State)
 }

 // 列出镜像
 images, err := client.ListImages(ctx, nil)
 if err != nil {
  panic(err)
 }
 fmt.Printf("Found %d images\n", len(images))
}
```

### 8.4 自定义 CNI 插件框架

```go
package main

import (
 "encoding/json"
 "fmt"
 "net"
 "os"

 "github.com/containernetworking/cni/pkg/skel"
 "github.com/containernetworking/cni/pkg/types"
 current "github.com/containernetworking/cni/pkg/types/100"
 "github.com/containernetworking/cni/pkg/version"
 "github.com/vishvananda/netlink"
)

// PluginConf CNI 插件配置结构
type PluginConf struct {
 types.NetConf
 Bridge     string   `json:"bridge"`
 Gateway    string   `json:"gateway"`
 Subnet     string   `json:"subnet"`
 MTU        int      `json:"mtu"`
 IsGateway  bool     `json:"isGateway"`
 IPMasq     bool     `json:"ipMasq"`
}

// 主入口函数
func main() {
 skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, "Custom CNI Plugin")
}

// cmdAdd 处理 ADD 命令
func cmdAdd(args *skel.CmdArgs) error {
 // 解析配置
 conf, err := parseConfig(args.StdinData)
 if err != nil {
  return fmt.Errorf("failed to parse config: %v", err)
 }

 // 创建或获取网桥
 br, err := setupBridge(conf.Bridge, conf.MTU, conf.IsGateway, conf.Gateway)
 if err != nil {
  return fmt.Errorf("failed to setup bridge: %v", err)
 }

 // 创建 veth pair
 hostVeth, containerVeth, err := setupVeth(args.Netns, args.IfName, conf.MTU, br)
 if err != nil {
  return fmt.Errorf("failed to setup veth: %v", err)
 }

 // 分配 IP 地址
 ipConfig, err := allocateIP(conf.Subnet, containerVeth)
 if err != nil {
  return fmt.Errorf("failed to allocate IP: %v", err)
 }

 // 设置路由
 if err := setupRoutes(args.Netns, conf.Gateway); err != nil {
  return fmt.Errorf("failed to setup routes: %v", err)
 }

 // 设置 IP Masquerade
 if conf.IPMasq {
  if err := setupIPMasq(conf.Subnet); err != nil {
   return fmt.Errorf("failed to setup IP masq: %v", err)
  }
 }

 // 构建结果
 result := &current.Result{
  CNIVersion: current.ImplementedSpecVersion,
  Interfaces: []*current.Interface{
   {
    Name: br.Attrs().Name,
    Mac:  br.Attrs().HardwareAddr.String(),
   },
   {
    Name:    hostVeth.Attrs().Name,
    Mac:     hostVeth.Attrs().HardwareAddr.String(),
    Sandbox: "", // 主机端
   },
   {
    Name:    args.IfName,
    Mac:     containerVeth.Attrs().HardwareAddr.String(),
    Sandbox: args.Netns,
   },
  },
  IPs: []*current.IPConfig{ipConfig},
 }

 return types.PrintResult(result, conf.CNIVersion)
}

// cmdDel 处理 DEL 命令
func cmdDel(args *skel.CmdArgs) error {
 conf, err := parseConfig(args.StdinData)
 if err != nil {
  return err
 }

 // 删除 veth 接口
 if err := deleteVeth(args.Netns, args.IfName); err != nil {
  return fmt.Errorf("failed to delete veth: %v", err)
 }

 // 清理 IP Masquerade 规则
 if conf.IPMasq {
  if err := teardownIPMasq(conf.Subnet); err != nil {
   return fmt.Errorf("failed to teardown IP masq: %v", err)
  }
 }

 return nil
}

// cmdCheck 处理 CHECK 命令
func cmdCheck(args *skel.CmdArgs) error {
 conf, err := parseConfig(args.StdinData)
 if err != nil {
  return err
 }

 // 检查网桥是否存在
 _, err = netlink.LinkByName(conf.Bridge)
 if err != nil {
  return fmt.Errorf("bridge %s not found: %v", conf.Bridge, err)
 }

 // 检查容器接口是否存在
 // ...

 return nil
}

// parseConfig 解析 CNI 配置
func parseConfig(stdin []byte) (*PluginConf, error) {
 conf := &PluginConf{}
 if err := json.Unmarshal(stdin, conf); err != nil {
  return nil, fmt.Errorf("failed to parse config: %v", err)
 }

 // 设置默认值
 if conf.Bridge == "" {
  conf.Bridge = "cni0"
 }
 if conf.MTU == 0 {
  conf.MTU = 1500
 }

 return conf, nil
}

// setupBridge 创建或获取网桥
func setupBridge(name string, mtu int, isGateway bool, gateway string) (*netlink.Bridge, error) {
 // 尝试获取现有网桥
 br, err := netlink.LinkByName(name)
 if err == nil {
  return br.(*netlink.Bridge), nil
 }

 // 创建新网桥
 bridge := &netlink.Bridge{
  LinkAttrs: netlink.LinkAttrs{
   Name: name,
   MTU:  mtu,
  },
 }

 if err := netlink.LinkAdd(bridge); err != nil {
  return nil, fmt.Errorf("failed to create bridge: %v", err)
 }

 // 设置网桥 IP
 if isGateway && gateway != "" {
  gwIP, ipNet, err := net.ParseCIDR(gateway)
  if err != nil {
   return nil, fmt.Errorf("invalid gateway: %v", err)
  }
  addr := &netlink.Addr{
   IPNet: &net.IPNet{
    IP:   gwIP,
    Mask: ipNet.Mask,
   },
  }
  if err := netlink.AddrAdd(bridge, addr); err != nil {
   return nil, fmt.Errorf("failed to add address to bridge: %v", err)
  }
 }

 // 启动网桥
 if err := netlink.LinkSetUp(bridge); err != nil {
  return nil, fmt.Errorf("failed to set bridge up: %v", err)
 }

 return bridge, nil
}

// setupVeth 创建 veth pair
func setupVeth(netnsPath, ifName string, mtu int, bridge *netlink.Bridge) (*netlink.Veth, *netlink.Veth, error) {
 // 生成唯一名称
 hostName := fmt.Sprintf("veth%s", generateID(8))

 // 创建 veth pair
 veth := &netlink.Veth{
  LinkAttrs: netlink.LinkAttrs{
   Name: hostName,
   MTU:  mtu,
  },
  PeerName: ifName,
 }

 if err := netlink.LinkAdd(veth); err != nil {
  return nil, nil, fmt.Errorf("failed to create veth: %v", err)
 }

 // 获取 veth 接口
 hostVeth, err := netlink.LinkByName(hostName)
 if err != nil {
  return nil, nil, err
 }

 containerVeth, err := netlink.LinkByName(ifName)
 if err != nil {
  return nil, nil, err
 }

 // 将主机端 veth 连接到网桥
 if err := netlink.LinkSetMaster(hostVeth, bridge); err != nil {
  return nil, nil, fmt.Errorf("failed to attach veth to bridge: %v", err)
 }

 // 启动主机端 veth
 if err := netlink.LinkSetUp(hostVeth); err != nil {
  return nil, nil, err
 }

 // 将容器端 veth 移动到目标命名空间
 // 这里需要调用 setns 系统调用
 // 简化实现，实际需要使用 github.com/containernetworking/plugins/pkg/ns

 return hostVeth.(*netlink.Veth), containerVeth.(*netlink.Veth), nil
}

// allocateIP 分配 IP 地址
func allocateIP(subnet string, containerVeth netlink.Link) (*current.IPConfig, error) {
 _, ipNet, err := net.ParseCIDR(subnet)
 if err != nil {
  return nil, err
 }

 // 简单的 IP 分配逻辑（实际应使用 IPAM 插件）
 // 这里使用固定 IP 作为示例
 ip := make(net.IP, len(ipNet.IP))
 copy(ip, ipNet.IP)
 ip[3] = 100 // 示例: x.x.x.100

 ipConfig := &current.IPConfig{
  Address: net.IPNet{
   IP:   ip,
   Mask: ipNet.Mask,
  },
  Gateway: ipNet.IP,
  Interface: current.Int(2), // 指向容器接口
 }

 return ipConfig, nil
}

// setupRoutes 设置路由
func setupRoutes(netnsPath, gateway string) error {
 // 添加默认路由
 // 实际实现需要在目标命名空间中执行
 return nil
}

// setupIPMasq 设置 IP Masquerade
func setupIPMasq(subnet string) error {
 // 使用 iptables 设置 MASQUERADE
 // 简化实现
 return nil
}

// teardownIPMasq 清理 IP Masquerade
func teardownIPMasq(subnet string) error {
 return nil
}

// deleteVeth 删除 veth 接口
func deleteVeth(netnsPath, ifName string) error {
 // 在目标命名空间中查找并删除接口
 return nil
}

// generateID 生成随机 ID
func generateID(length int) string {
 const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
 b := make([]byte, length)
 for i := range b {
  b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
 }
 return string(b)
}
```

### 8.5 高级 client-go 操作

```go
package main

import (
 "context"
 "fmt"
 "time"

 "k8s.io/apimachinery/pkg/api/meta"
 "k8s.io/apimachinery/pkg/runtime/schema"
 "k8s.io/client-go/discovery"
 "k8s.io/client-go/dynamic"
 "k8s.io/client-go/informers"
 "k8s.io/client-go/metadata"
 "k8s.io/client-go/tools/cache"
 "k8s.io/client-go/tools/leaderelection"
 "k8s.io/client-go/tools/leaderelection/resourcelock"
)

// DynamicClient 动态客户端操作
type DynamicClient struct {
 client dynamic.Interface
}

// GetResource 使用动态客户端获取任意资源
func (d *DynamicClient) GetResource(ctx context.Context, gvr schema.GroupVersionResource,
 namespace, name string) (*unstructured.Unstructured, error) {
 return d.client.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

// ListResources 使用动态客户端列出资源
func (d *DynamicClient) ListResources(ctx context.Context, gvr schema.GroupVersionResource,
 namespace string, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
 if namespace == "" {
  return d.client.Resource(gvr).List(ctx, opts)
 }
 return d.client.Resource(gvr).Namespace(namespace).List(ctx, opts)
}

// CreateOrUpdateResource 创建或更新资源
func (d *DynamicClient) CreateOrUpdateResource(ctx context.Context, gvr schema.GroupVersionResource,
 namespace string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {

 existing, err := d.client.Resource(gvr).Namespace(namespace).Get(
  ctx, obj.GetName(), metav1.GetOptions{})

 if err != nil {
  if errors.IsNotFound(err) {
   return d.client.Resource(gvr).Namespace(namespace).Create(ctx, obj, metav1.CreateOptions{})
  }
  return nil, err
 }

 // 更新资源版本
 obj.SetResourceVersion(existing.GetResourceVersion())
 return d.client.Resource(gvr).Namespace(namespace).Update(ctx, obj, metav1.UpdateOptions{})
}

// DiscoveryClient API 发现和元数据操作
type DiscoveryClient struct {
 discoveryClient discovery.DiscoveryInterface
 mapper          meta.RESTMapper
}

// GetAPIResources 获取所有 API 资源
func (d *DiscoveryClient) GetAPIResources() ([]*metav1.APIResourceList, error) {
 return d.discoveryClient.ServerPreferredResources()
}

// GetGVKForResource 获取资源的 GVK
func (d *DiscoveryClient) GetGVKForResource(gvr schema.GroupVersionResource) (schema.GroupVersionKind, error) {
 mappings, err := d.mapper.RESTMappings(gvr.GroupResource())
 if err != nil {
  return schema.GroupVersionKind{}, err
 }
 if len(mappings) == 0 {
  return schema.GroupVersionKind{}, fmt.Errorf("no mappings found")
 }
 return mappings[0].GroupVersionKind, nil
}

// LeaderElection 领导者选举示例
func runWithLeaderElection(ctx context.Context, client kubernetes.Interface, runFunc func(ctx context.Context)) {
 // 创建锁
 lock := &resourcelock.LeaseLock{
  LeaseMeta: metav1.ObjectMeta{
   Name:      "my-controller-lock",
   Namespace: "kube-system",
  },
  Client: client.CoordinationV1(),
  LockConfig: resourcelock.ResourceLockConfig{
   Identity: "my-controller-1",
  },
 }

 // 配置领导者选举
 lec := leaderelection.LeaderElectionConfig{
  Lock:            lock,
  LeaseDuration:   15 * time.Second,
  RenewDeadline:   10 * time.Second,
  RetryPeriod:     2 * time.Second,
  ReleaseOnCancel: true,
  Callbacks: leaderelection.LeaderCallbacks{
   OnStartedLeading: func(ctx context.Context) {
    fmt.Println("Started leading")
    runFunc(ctx)
   },
   OnStoppedLeading: func() {
    fmt.Println("Stopped leading")
   },
   OnNewLeader: func(identity string) {
    fmt.Printf("New leader elected: %s\n", identity)
   },
  },
 }

 // 启动领导者选举
 le, err := leaderelection.NewLeaderElector(lec)
 if err != nil {
  panic(err)
 }

 le.Run(ctx)
}

// SharedInformerFactory 共享 Informer 工厂
type InformerFactory struct {
 factory informers.SharedInformerFactory
}

// NewInformerFactory 创建共享 Informer 工厂
func NewInformerFactory(client kubernetes.Interface, defaultResync time.Duration) *InformerFactory {
 return &InformerFactory{
  factory: informers.NewSharedInformerFactory(client, defaultResync),
 }
}

// AddPodHandler 添加 Pod 事件处理器
func (f *InformerFactory) AddPodHandler(handler cache.ResourceEventHandler) {
 podInformer := f.factory.Core().V1().Pods().Informer()
 podInformer.AddEventHandler(handler)
}

// Start 启动所有 Informer
func (f *InformerFactory) Start(stopCh <-chan struct{}) {
 f.factory.Start(stopCh)
}

// WaitForCacheSync 等待缓存同步
func (f *InformerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool {
 return f.factory.WaitForCacheSync(stopCh)
}

// MetadataClient 元数据客户端（轻量级）
type MetadataClient struct {
 client metadata.Interface
}

// ListPodsMetadata 列出 Pod 元数据（无 spec/status）
func (m *MetadataClient) ListPodsMetadata(ctx context.Context, namespace string) (*metav1.PartialObjectMetadataList, error) {
 return m.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
}

// Indexer 自定义索引器示例
func createIndexer(informer cache.SharedIndexInformer) {
 // 添加自定义索引器
 informer.AddIndexers(cache.Indexers{
  "nodeName": func(obj interface{}) ([]string, error) {
   pod, ok := obj.(*corev1.Pod)
   if !ok {
    return nil, fmt.Errorf("not a pod")
   }
   return []string{pod.Spec.NodeName}, nil
  },
 })
}

// 使用索引器查询
func queryByNodeName(informer cache.SharedIndexInformer, nodeName string) ([]interface{}, error) {
 return informer.GetIndexer().ByIndex("nodeName", nodeName)
}
```

---

## 9. 总结与最佳实践

### 9.1 架构设计原则

1. **单一职责原则**：每个组件只负责一个明确的职责
2. **可扩展性**：通过接口标准（CRI/CNI/CSI）支持多种实现
3. **声明式 API**：用户描述期望状态，系统负责调和
4. **最终一致性**：允许短暂不一致，保证最终状态正确
5. **松耦合**：组件间通过标准接口通信，降低依赖

### 9.2 生产环境建议

| 场景 | 推荐方案 |
|------|----------|
| **容器运行时** | containerd (默认) 或 CRI-O |
| **CNI 插件** | Calico (通用) / Cilium (高性能) |
| **Ingress** | NGINX (通用) / Traefik (云原生) |
| **存储** | CSI 驱动 + StorageClass |
| **监控** | Prometheus + Grafana |
| **日志** | Fluent Bit + Elasticsearch + Kibana |
| **GitOps** | ArgoCD |
| **服务网格** | Istio (功能丰富) / Linkerd (轻量) |

### 9.3 性能优化要点

1. **API Server**: 启用缓存、增加副本数、使用优先/公平队列
2. **etcd**: 使用 SSD、定期压缩、合理配置配额
3. **调度器**: 自定义调度策略、使用 Pod 亲和性/反亲和性
4. **Kubelet**: 调整并行镜像拉取数、优化 PLEG 轮询间隔
5. **网络**: 使用 IPVS/eBPF 模式、优化 CNI 插件配置

---

## 参考资源

- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [containerd 文档](https://containerd.io/docs/)
- [OCI 规范](https://opencontainers.org/)
- [CNI 规范](https://www.cni.dev/)
- [CSI 规范](https://github.com/container-storage-interface/spec)
- [client-go 文档](https://github.com/kubernetes/client-go)
- [Kubebuilder 文档](https://book.kubebuilder.io/)

---

*文档版本: 1.0*
*最后更新: 2024*
