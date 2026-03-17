# 容器运行时生态

> containerd、CRI-O 和其他 OCI 运行时

---

## 容器运行时层次

```
┌─────────────────────────────────────────────────────────────┐
│                  容器运行时架构层次                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  高层运行时 (High-Level Runtime)                      │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │  │
│  │  │   Docker    │  │  containerd │  │   CRI-O     │   │  │
│  │  │   Engine    │  │             │  │             │   │  │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘   │  │
│  │         │                │                │          │  │
│  │         └────────────────┴────────────────┘          │  │
│  │                      │                               │  │
│  │                      ↓ CRI (Container Runtime Interface)│
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  低层运行时 (Low-Level Runtime) / OCI Runtime         │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐    │  │
│  │  │  runc   │ │  youki  │ │ crun    │ │runsc    │    │  │
│  │  │  (Go)   │ │  (Rust) │ │  (C)    │ │(gVisor) │    │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘    │  │
│  │                                                      │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────────────────────┐│  │
│  │  │ kata    │ │Firecracker│ │   runwasi (Wasm)      ││  │
│  │  └─────────┘ └─────────┘ └─────────────────────────┘│  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                     Linux Kernel                      │  │
│  │              (cgroups, namespaces, eBPF)              │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## containerd 2.0 (2025)

containerd 2.0 是 K8s 1.33+ 的推荐运行时，带来重大改进。

### 新特性

| 特性 | 说明 |
|------|------|
| **Transfer Service** | 简化镜像拉取/推送 API |
| **Nerdctl v2.0** | 完整 Docker CLI 兼容 |
| **Sandbox API** | 统一沙箱容器支持 |
| **Wasm 支持** | 内置 WebAssembly 运行时支持 |
| **用户命名空间** | 完整支持 UID/GID 映射 |

### 配置文件

```toml
version = 3
root = "/var/lib/containerd"
state = "/run/containerd"

[grpc]
address = "/run/containerd/containerd.sock"

[plugins."io.containerd.grpc.v1.cri"]
sandbox_image = "registry.k8s.io/pause:3.9"

[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
runtime_type = "io.containerd.runc.v2"

[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runsc]
runtime_type = "io.containerd.runsc.v1"

[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runwasi]
runtime_type = "io.containerd.runwasi.v1"
```

---

## Nerdctl - Docker 兼容 CLI

```bash
# 使用体验与 Docker 完全一致
nerdctl run -d --name nginx -p 80:80 nginx:latest
nerdctl ps
nerdctl exec -it nginx sh
nerdctl build -t myapp .
nerdctl compose up -d
```

---

## 运行时对比 (2025)

| 特性 | Docker | containerd | CRI-O |
|------|--------|------------|-------|
| **K8s 默认** | 否 | **是** | OpenShift |
| **镜像构建** | 原生 | 需 BuildKit | 不支持 |
| **Docker CLI** | 原生 | nerdctl | crictl |
| **资源占用** | 高 | **低** | 低 |
| **安全** | 一般 | **好** | 好 |
| **Wasm 支持** | 实验 | **原生** | 计划 |

---

## OCI 运行时演进

| 运行时 | 语言 | 特点 | 适用场景 |
|--------|------|------|----------|
| **runc** | Go | OCI 参考实现 | 通用 |
| **youki** | Rust | 更快、更安全 | 安全敏感 |
| **crun** | C | 更快、更轻量 | 资源受限 |
| **runsc** | Go | gVisor 沙箱 | 多租户 |
| **kata** | Go | Kata Containers | 强隔离 |
| **runwasi** | Rust | WebAssembly | 边缘/Serverless |

---

## 2025 趋势

- **containerd 2.0**: 成为事实标准
- **youki**: Rust 运行时崛起
- **runwasi**: Wasm 工作负载增长
- **沙箱化**: 安全隔离成为标配
