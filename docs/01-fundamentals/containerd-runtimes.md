# 容器运行时生态

> containerd 2.0、CRI-O 和其他 OCI 运行时 - 2025 最新

---

## 容器运行时架构层次

```
┌─────────────────────────────────────────────────────────────┐
│                  容器运行时架构层次                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  高层运行时 (High-Level Runtime)                      │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │  │
│  │  │   Docker    │  │  containerd │  │   CRI-O     │   │  │
│  │  │   Engine    │  │    2.0      │  │             │   │  │
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

## containerd 2.0 (2024年11月发布)

containerd 2.0 是 1.x 系列以来的第一个主版本，带来重大改进。

### containerd 2.0 新特性

| 特性 | 说明 |
|------|------|
| **NRI 默认启用** | Node Resource Interface 插件系统 |
| **Image Verifier 插件** | 镜像验证策略执行 |
| **Transfer Service** | 简化镜像拉取/推送 API |
| **Nerdctl v2.0** | 完整 Docker CLI 兼容 |
| **用户命名空间** | 完整支持 UID/GID 映射 |
| **igzip 支持** | Intel ISA-L 加速镜像解压 |
| **容器检查点** | CRIU 容器状态保存/恢复 |
| **CDI 默认启用** | Container Device Interface |

### 配置迁移

```bash
# 从 v2 配置迁移到 v3
containerd config migrate > /etc/containerd/config.toml.v3

# 验证配置
containerd config validate
```

### 配置文件 (v3)

```toml
version = 3
root = "/var/lib/containerd"
state = "/run/containerd"

[grpc]
address = "/run/containerd/containerd.sock"

[plugins.'io.containerd.cri.v1.images']
# 镜像验证插件
image_verifier_plugins = ['cosign']

[plugins.'io.containerd.cri.v1.images'.image_verifier_plugins.cosign]
type = 'cosign'
config_path = '/etc/containerd/cosign'

[plugins.'io.containerd.cri.v1.runtime']
sandbox_image = 'registry.k8s.io/pause:3.9'

[plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.runc]
runtime_type = 'io.containerd.runc.v2'

[plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.runsc]
runtime_type = 'io.containerd.runsc.v1'

[plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.runwasi]
runtime_type = 'io.containerd.runwasi.v1'
```

---

## NRI (Node Resource Interface)

containerd 2.0 默认启用 NRI，允许自定义容器配置。

```
NRI 插件类型:
├── 资源调整 (Resource Adjustment)
├── 设备注入 (Device Injection)
├── OCI 规范修改 (OCI Spec Modification)
└── 事件处理 (Event Handling)
```

### NRI 与 Kubernetes 集成

```yaml
# NRI 插件 Pod
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nri-plugin
spec:
  template:
    spec:
      hostNetwork: true
      containers:
        - name: plugin
          image: nri-sample-plugin:latest
          volumeMounts:
            - name: nri-socket
              mountPath: /var/run/nri
      volumes:
        - name: nri-socket
          hostPath:
            path: /var/run/nri
```

---

## 镜像验证插件

containerd 2.0 支持镜像验证插件，在镜像拉取时执行策略。

```toml
# /etc/containerd/config.toml
[plugins.'io.containerd.cri.v1.images']
image_verifier_plugins = ['cosign']

[plugins.'io.containerd.cri.v1.images'.image_verifier_plugins.cosign]
type = 'cosign'
config_path = '/etc/containerd/cosign'
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

# containerd 2.0 新特性
nerdctl checkpoint create mycontainer mycheckpoint
nerdctl restore mycheckpoint
```

---

## 运行时对比 (2025)

| 特性 | Docker | containerd 2.0 | CRI-O |
|------|--------|----------------|-------|
| **K8s 默认** | 否 | **是** | OpenShift |
| **镜像构建** | 原生 | 需 BuildKit | 不支持 |
| **Docker CLI** | 原生 | nerdctl v2 | crictl |
| **资源占用** | 高 | **低** | 低 |
| **安全** | 一般 | **好** | 好 |
| **Wasm 支持** | 实验 | **原生** | 计划 |
| **NRI** | 否 | **默认启用** | 否 |
| **镜像验证** | 需配置 | **内置插件** | 需配置 |

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

## 用户命名空间 (User Namespaces)

containerd 2.0 完全支持用户命名空间，配合 K8s 1.33+。

```bash
# 启用用户命名空间
echo "user.max_user_namespaces = 1048576" >> /etc/sysctl.conf
sysctl -p

# 验证
sysctl user.max_user_namespaces
```

```yaml
# Kubernetes Pod
apiVersion: v1
kind: Pod
metadata:
  name: userns-pod
spec:
  hostUsers: false  # 启用用户命名空间
  containers:
    - name: app
      image: nginx
      securityContext:
        runAsUser: 0
```

---

## 2025 趋势

- ✅ **containerd 2.0**: 成为事实标准，K8s 1.33+ 推荐
- ✅ **NRI 普及**: 插件生态系统发展
- ✅ **镜像验证**: 供应链安全内置支持
- 🔄 **youki**: Rust 运行时崛起
- 🔄 **runwasi**: Wasm 工作负载增长
- 🔄 **沙箱化**: 安全隔离成为标配

---

## 参考

- [containerd 2.0 Release Notes](https://github.com/containerd/containerd/releases/tag/v2.0.0)
- [NRI 文档](https://github.com/containerd/nri)
- [Nerdctl 文档](https://github.com/containerd/nerdctl)
