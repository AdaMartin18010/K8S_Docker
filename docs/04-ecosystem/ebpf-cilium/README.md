# eBPF 与 Cilium 网络

> 基于 eBPF 的下一代 Kubernetes 网络方案 - Cilium 1.18 新特性

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
| **性能** | O(n) 规则链遍历 | **O(1) 哈希查找** |
| **延迟** | 随规则增加而增加 | **恒定低延迟** |
| **策略粒度** | L3/L4 (IP/端口) | **L3-L7 (支持 HTTP/gRPC)** |
| **可观测性** | 有限 | **Hubble 实时流量可视化** |
| **扩展性** | 节点压力时性能下降 | **内核级处理，无压力** |
| **Service 数量** | 200+ 服务性能下降 | **2万+ 服务稳定** |

### 性能对比

```
网络延迟: Calico(iptables) 250μs → Cilium(eBPF) 60μs
CPU 使用: kube-proxy 100% → Cilium 35%
服务查找: O(n) 规则链 → O(1) BPF 哈希表
```

---

## 替换 kube-proxy

Cilium 完全替换 kube-proxy，使用 eBPF 实现高效的 Service 负载均衡。

```bash
# 使用 Helm 安装，完全替换 kube-proxy
helm install cilium cilium/cilium \
  --namespace kube-system \
  --set kubeProxyReplacement=strict \
  --set k8sServiceHost=auto \
  --set k8sServicePort=6443 \
  --set hubble.enabled=true \
  --set hubble.relay.enabled=true \
  --set hubble.ui.enabled=true
```

### kube-proxy 替换对比

| 特性 | kube-proxy (iptables) | Cilium eBPF |
|------|----------------------|-------------|
| **Service lookup** | 顺序规则匹配 (O(n)) | 哈希表查找 (O(1)) |
| **East-west LB** | 每包 DNAT | Socket-level (connect()) |
| **规则更新** | 刷新重写所有链 | 原子 BPF map 更新 |
| **可观测性** | 无 | Hubble (L3/L4/L7) |
| **Masquerading** | iptables MASQUERADE | eBPF (bpf.masquerade) |

### Socket-Level 负载均衡

Cilium 在 `connect()` 系统调用时拦截，将 ClusterIP 直接转换为后端 Pod IP，**无需 per-packet DNAT**。

```
Pod A → connect(ClusterIP:80)
  → Cilium eBPF: ClusterIP → Pod IP (socket level)
  → 直接建立到 Pod IP 的连接
  → 无 SNAT，无 conntrack 条目
```

---

## Cilium 1.18 新特性

### L2 Announcements (Beta)

替代 MetalLB，为裸金属集群提供 LoadBalancer IP。

```yaml
apiVersion: cilium.io/v2alpha1
kind: CiliumLoadBalancerIPPool
metadata:
  name: homelab-pool
spec:
  blocks:
    - start: "10.10.30.20"
      stop: "10.10.30.99"
---
apiVersion: cilium.io/v2alpha1
kind: CiliumL2AnnouncementPolicy
metadata:
  name: homelab-l2
spec:
  interfaces:
    - ^eno.*
    - ^eth.*
    - ^enp.*
  nodeSelector:
    matchLabels:
      kubernetes.io/os: linux
  loadBalancerIPs: true
  externalIPs: true
```

### Gateway API 支持 (GA)

Cilium 1.18 完全支持 Gateway API v1.2 实验特性。

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: cilium-gateway
spec:
  gatewayClassName: cilium
  listeners:
    - name: http
      protocol: HTTP
      port: 80
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: api-route
spec:
  parentRefs:
    - name: cilium-gateway
  rules:
    - backendRefs:
        - name: api-service
          port: 80
```

### eBPF Host Routing

启用原生 eBPF 主机路由，绕过 veth pair 和主机网络栈。

```bash
helm upgrade cilium cilium/cilium \
  --namespace kube-system \
  --set bpf.hostRouting=true \
  --set bpf.masquerade=true
```

### DSR (Direct Server Return)

响应流量绕过负载均衡节点，直接返回客户端。

```bash
helm upgrade cilium cilium/cilium \
  --namespace kube-system \
  --set loadBalancer.mode=dsr
```

### Maglev 一致性哈希

后端变化时最小化活跃连接影响。

```bash
helm upgrade cilium cilium/cilium \
  --namespace kube-system \
  --set loadBalancer.algorithm=maglev \
  --set maglev.tableSize=65521
```

---

## Hubble 可观测性

Hubble 读取 eBPF 程序的流量数据，提供 L3-L7 可视化。

```bash
# 查看所有流量
hubble observe --namespace production

# 查看被丢弃的流量
hubble observe --verdict DROPPED

# 查看特定服务的 HTTP 流量
hubble observe --pod grafana --protocol http

# 按命名空间查看流量
hubble observe --from-namespace monitoring
```

### Hubble 指标导出

```yaml
hubble:
  enabled: true
  metrics:
    enabled:
      - dns:query
      - drop
      - tcp
      - flow
      - icmp
      - http
    serviceMonitor:
      enabled: true
```

---

## 网络策略 (L3-L7)

### L7 HTTP 策略

```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: http-restrict
spec:
  endpointSelector:
    matchLabels:
      app: frontend
  ingress:
    - fromEndpoints:
        - matchLabels:
            app: backend
      toPorts:
        - ports:
            - port: "80"
              protocol: TCP
          rules:
            http:
              - method: GET
                path: "/api/.*"
              - method: POST
                path: "/api/users"
```

---

## 验证 kube-proxy 替换

```bash
# 确认没有 kube-proxy pods
kubectl get pods -n kube-system -l k8s-app=kube-proxy
# (空结果表示成功替换)

# 检查 Cilium 状态
cilium status
# KubeProxyReplacement: True [eth0, Direct Routing]

# 查看 BPF maps
cilium bpf lb list
cilium bpf ipcache list
```

---

## 主流 CNI 对比 (2025)

| 特性 | Calico | Cilium | Flannel |
|------|--------|--------|---------|
| **数据平面** | eBPF/iptables | **eBPF** | VXLAN |
| **网络策略** | L3-L4 | **L3-L7** | 无 |
| **Service Mesh** | 无 | **内置 (Envoy)** | 无 |
| **多集群** | 支持 | **Cluster Mesh** | 不支持 |
| **kube-proxy 替代** | 否 | **是** | 否 |
| **Gateway API** | 否 | **是** | 否 |
| **Observability** | 有限 | **Hubble** | 无 |

---

## Cilium 替代方案整合

| 功能 | 传统方案 | Cilium 单一方案 |
|------|----------|----------------|
| CNI | Calico/Flannel | Cilium |
| Service Proxy | kube-proxy | Cilium eBPF |
| Network Policy | Calico | Cilium |
| LoadBalancer | MetalLB | Cilium L2 Announcements |
| Ingress | NGINX/Traefik | Cilium Gateway API |
| Observability | Prometheus Exporters | Hubble |
| Service Mesh | Istio | Cilium Service Mesh |

**优势**: 更少的组件 = 更少的升级、更少的资源消耗、更少的故障点

---

## 2025 趋势

- ✅ **Cilium 1.18**: Gateway API GA，L2 Announcements Beta
- ✅ **GKE Dataplane V2**: Cilium 作为默认 CNI
- 🔄 **Tetragon**: 基于 eBPF 的运行时安全
- 🔄 **EBPF 成为标准**: 新一代可观测性和安全工具的基础
- 🔄 **Cilium Service Mesh**: 无 Sidecar 服务网格
- 🔄 **Cluster Mesh**: 多集群无边界网络

---

## 参考

- [Cilium 官方文档](https://docs.cilium.io/)
- [Cilium 1.18 Release Notes](https://github.com/cilium/cilium/releases)
- [eBPF 介绍](https://ebpf.io/)
- [Hubble 文档](https://docs.cilium.io/en/stable/observability/hubble/)
