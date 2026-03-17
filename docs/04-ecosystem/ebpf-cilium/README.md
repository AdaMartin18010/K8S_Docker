# eBPF 与 Cilium 网络

> 基于 eBPF 的下一代 Kubernetes 网络方案

---

## 什么是 eBPF？

eBPF (Extended Berkeley Packet Filter) 是 Linux 内核技术，允许在无需修改内核代码或加载内核模块的情况下，安全地在内核中运行沙箱程序。

```
┌─────────────────────────────────────────────────────────────┐
│                      eBPF 架构                               │
├─────────────────────────────────────────────────────────────┤
│                    ┌───────────────┐                        │
│                    │ 用户空间程序   │                        │
│                    │  (Go/C/Rust)  │                        │
│                    └───────┬───────┘                        │
│                            │ libbpf/bpf syscall             │
│                            ↓                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                    Linux Kernel                       │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │         eBPF Verifier (安全检查)               │  │  │
│  │  │    • 循环检测    • 内存边界检查                │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  │                         ↓                            │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │         eBPF JIT Compiler                      │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  │                         ↓                            │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │  KPROBE     │  │  TRACEPOINT │  │   XDP       │  │  │
│  │  │  (函数追踪)  │  │  (内核事件)  │  │ (网卡层)     │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## Cilium: eBPF 驱动的 Kubernetes 网络

Cilium 是一个基于 eBPF 的开源项目，为 Kubernetes 提供网络、安全和可观测性。

### 为什么选择 Cilium？

| 特性 | 传统 CNI (iptables) | Cilium (eBPF) |
|------|---------------------|---------------|
| **性能** | O(n) 规则链遍历 | O(1) 直接查找 |
| **延迟** | 随规则增加而增加 | 恒定低延迟 |
| **策略粒度** | L3/L4 (IP/端口) | **L3-L7 (支持 HTTP/gRPC)** |
| **可观测性** | 有限 | **Hubble 实时流量可视化** |
| **扩展性** | 节点压力时性能下降 | **内核级处理，无压力** |

### 性能对比

```
网络延迟: Calico(iptables) 250μs → Cilium(eBPF) 60μs
CPU 使用: kube-proxy 100% → Cilium 35%
```

---

## 替换 kube-proxy

```bash
# 使用 Helm 安装，完全替换 kube-proxy
helm install cilium cilium/cilium \
  --namespace kube-system \
  --set kubeProxyReplacement=strict \
  --set hubble.enabled=true \
  --set hubble.relay.enabled=true \
  --set hubble.ui.enabled=true
```

---

## Hubble 可观测性

```bash
# 查看所有流量
hubble observe --namespace production

# 查看被丢弃的流量
hubble observe --verdict DROPPED
```

---

## 主流 CNI 对比 (2025)

| 特性 | Calico | Cilium | Flannel |
|------|--------|--------|---------|
| **数据平面** | eBPF/iptables | **eBPF** | VXLAN |
| **网络策略** | L3-L4 | **L3-L7** | 无 |
| **Service Mesh** | 无 | **内置** | 无 |
| **多集群** | 支持 | **Cluster Mesh** | 不支持 |

---

## 2025 趋势

- **GKE Dataplane V2**: Cilium 作为默认 CNI
- **Tetragon**: 基于 eBPF 的运行时安全
- **EBPF 成为标准**: 新一代可观测性和安全工具的基础
