# 沙箱化容器运行时

> 增强容器隔离的安全技术

---

## 为什么需要沙箱？

传统容器共享主机内核，存在"容器逃逸"风险。沙箱技术提供更强的隔离级别。

```
┌─────────────────────────────────────────────────────────────┐
│                  隔离级别对比                                │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  传统容器                      沙箱容器                       │
│  ┌──────────────┐            ┌────────────────────────────┐ │
│  │   App A      │            │   ┌────────────────────┐   │ │
│  │   Bin/Lib    │            │   │   App A            │   │ │
│  ├──────────────┤            │   │   Bin/Lib          │   │ │
│  │   Host Kernel│            │   ├────────────────────┤   │ │
│  │   (Shared)   │            │   │   MicroVM Kernel   │   │ │
│  ├──────────────┤            │   ├────────────────────┤   │ │
│  │   Hardware   │            │   │   VMM (KVM)        │   │ │
│  └──────────────┘            │   ├────────────────────┤   │ │
│                              │   │   Host Kernel      │   │ │
│  隔离级别: 进程级               │   ├────────────────────┤   │ │
│  启动时间: ~1s                  │   │   Hardware         │   │ │
│                              │   └────────────────────┘   │ │
│                              └────────────────────────────┘ │
│                                                              │
│                              隔离级别: 虚拟机级                │
│                              启动时间: ~100ms                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 主流沙箱技术

| 技术 | 原理 | 性能 | 适用场景 |
|------|------|------|----------|
| **gVisor** | 用户空间内核 (Sentry) | 中 | 多租户、不可信代码 |
| **Kata Containers** | 轻量级 VM | 高 | 金融、政务、隔离敏感 |
| **Firecracker** | MicroVM (AWS) | 高 | Serverless、边缘计算 |
| **Quark** | 用户空间内核 | 高 | 高安全场景 |

---

## gVisor

Google 开发的沙箱运行时，使用用户空间内核拦截和验证系统调用。

```
┌─────────────────────────────────────────────────────────────┐
│                      gVisor 架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                    Sentry (用户空间内核)               │  │
│  │   • 实现 Linux 系统调用                              │  │
│  │   • 隔离应用与主机内核                                │  │
│  │   • Go 语言编写，内存安全                             │  │
│  └───────────────────────┬──────────────────────────────┘  │
│                          │ 9P / 自定义协议                 │
│  ┌───────────────────────▼──────────────────────────────┐  │
│  │                    Gofer (文件代理)                   │  │
│  │   • 处理所有文件系统操作                              │  │
│  │   • 双重验证路径访问                                  │  │
│  └───────────────────────┬──────────────────────────────┘  │
│                          │                                 │
│  ┌───────────────────────▼──────────────────────────────┐  │
│  │                    Host Kernel                        │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 使用 gVisor

```bash
# 安装 runsc
wget https://storage.googleapis.com/gvisor/releases/release/latest/x86_64/runsc
chmod +x runsc
mv runsc /usr/local/bin/

# 配置 containerd
cat <<EOF >> /etc/containerd/config.toml
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runsc]
  runtime_type = "io.containerd.runsc.v1"
EOF

# 在 K8s 中使用
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: gvisor
handler: runsc
---
apiVersion: v1
kind: Pod
spec:
  runtimeClassName: gvisor
  containers:
    - name: app
      image: myapp
```

---

## Kata Containers

使用轻量级虚拟机提供强隔离，每个 Pod 有自己的内核。

```
┌─────────────────────────────────────────────────────────────┐
│                   Kata Containers 架构                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Pod A                    Pod B                    Pod C    │
│  ┌─────────────────┐     ┌─────────────────┐     ┌────────┐ │
│  │ ┌─────────────┐ │     │ ┌─────────────┐ │     │┌──────┐│ │
│  │ │  Container  │ │     │ │  Container  │ │     ││Container││
│  │ ├─────────────┤ │     │ ├─────────────┤ │     │├──────┤│ │
│  │ │ Kata Kernel │ │     │ │ Kata Kernel │ │     ││Kata  ││ │
│  │ ├─────────────┤ │     │ ├─────────────┤ │     ││Kernel││ │
│  │ │    QEMU     │ │     │ │ Cloud Hyper │ │     ││Firecracker│
│  │ └─────────────┘ │     │ └─────────────┘ │     │└──────┘│ │
│  │    MicroVM       │     │    MicroVM       │     │ MicroVM│ │
│  └─────────────────┘     └─────────────────┘     └────────┘ │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Kata Runtime (containerd-shim)          │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                    Host Kernel                        │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 使用 Kata

```yaml
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: kata
handler: kata
---
apiVersion: v1
kind: Pod
metadata:
  name: secure-app
spec:
  runtimeClassName: kata
  containers:
    - name: app
      image: nginx
      resources:
        limits:
          memory: "1Gi"
          cpu: "1"
```

---

## Firecracker

AWS 开源的 MicroVM，专为 Serverless 和容器优化。

**特点**:

- 启动时间: < 125ms
- 内存占用: < 15MB
- 专为高并发设计

---

## 选择建议

| 场景 | 推荐方案 | 原因 |
|------|----------|------|
| **多租户 SaaS** | gVisor | 零信任，拦截所有系统调用 |
| **金融/政务** | Kata | 合规要求，强隔离 |
| **Serverless** | Firecracker | 快速启动，高密度 |
| **边缘计算** | Kata + Firecracker | 轻量 + 安全 |

---

## 用户命名空间 (User Namespaces)

K8s 1.33 + containerd 2.0 支持用户命名空间隔离。

```yaml
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
        allowPrivilegeEscalation: false
        runAsNonRoot: true
        seccompProfile:
          type: RuntimeDefault
```

**效果**: 容器内的 root 映射为主机的非特权用户。

---

## 安全运行时对比 (2025)

| 特性 | runc | gVisor | Kata | Firecracker |
|------|------|--------|------|-------------|
| **隔离级别** | 进程 | 系统调用 | VM | MicroVM |
| **启动时间** | 1s | 1.5s | 2s | 125ms |
| **内存开销** | 低 | 中 | 高 | 低 |
| **内核共享** | 是 | 否 | 否 | 否 |
| **适用场景** | 通用 | 多租户 | 高安全 | Serverless |
