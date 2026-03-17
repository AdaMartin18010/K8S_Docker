# eBPF 与 Cilium - 下一代 Kubernetes 网络

## 概述

eBPF (Extended Berkeley Packet Filter) 是 Linux 内核的革命性技术，允许在内核空间安全地运行沙盒程序。Cilium 是基于 eBPF 的 Kubernetes CNI，提供高性能网络、安全和可观测性。

> **2025 关键数据**: Cilium 已成为 GKE Dataplane V2 的默认 CNI，被 AWS、Azure 广泛支持，社区超过 23,000 开发者。

## 为什么替换 Kube-Proxy

### 传统 Kube-Proxy 的瓶颈

```
Kube-Proxy iptables 模式问题:
┌─────────────────────────────────────────────────────────────────┐
│  服务数量增长 → iptables 规则线性增长 → 包处理延迟增加            │
│                                                                 │
│  1000 个服务 = 10000+ 条 iptables 规则                          │
│  每个包遍历 O(n) 复杂度                                          │
│  更新规则时全表刷新，导致延迟尖峰                                │
└─────────────────────────────────────────────────────────────────┘
```

**性能对比**:
| 指标 | iptables | IPVS | eBPF (Cilium) |
|------|----------|------|---------------|
| 查找复杂度 | O(n) | O(1) | O(1) |
| 服务更新延迟 | 高 | 中 | 极低 |
| 连接跟踪 | conntrack | conntrack | eBPF map |
| 后端变化影响 | 全表刷新 | 部分刷新 | 增量更新 |

### Cilium eBPF 优势

- **O(1) 服务查找**: eBPF map 哈希表查找，与服务数量无关
- **无 conntrack**: 绕过内核 conntrack，减少资源消耗
- **DSR 支持**: Direct Server Return 模式降低延迟
- **Maglev 一致性哈希**: 后端变化时最小化连接中断

## Cilium 安装与配置

### 全新集群安装（无 kube-proxy）

```bash
# kubeadm 初始化时跳过 kube-proxy
kubeadm init --skip-phases=addon/kube-proxy \
  --pod-network-cidr=10.244.0.0/16

# 安装 Cilium
helm repo add cilium https://helm.cilium.io/
helm repo update

helm install cilium cilium/cilium \
  --namespace kube-system \
  --set kubeProxyReplacement=strict \
  --set k8sServiceHost=API_SERVER_IP \
  --set k8sServicePort=6443 \
  --set hostServices.enabled=true \
  --set externalIPs.enabled=true \
  --set nodePort.enabled=true \
  --set bpf.masquerade=true \
  --set ipam.mode=kubernetes
```

### 存量集群迁移

```bash
# 第1步：并行安装 Cilium
helm install cilium cilium/cilium \
  --namespace kube-system \
  --set kubeProxyReplacement=false

# 等待就绪
cilium status --wait

# 第2步：启用 kube-proxy 替代
helm upgrade cilium cilium/cilium \
  --namespace kube-system \
  --reuse-values \
  --set kubeProxyReplacement=true \
  --set k8sServiceHost=$(kubectl get endpoints kubernetes -o jsonpath='{.subsets[0].addresses[0].ip}') \
  --set k8sServicePort=6443

# 第3步：删除 kube-proxy
kubectl -n kube-system delete ds kube-proxy
kubectl -n kube-system delete cm kube-proxy

# 第4步：清理 iptables 规则
kubectl -n kube-system exec ds/cilium -- cilium cleanup -f
```

## 高级网络功能

### 1. Direct Server Return (DSR)

```yaml
# values.yaml
kubeProxyReplacement: "true"
loadBalancer:
  mode: dsr        # SNAT | DSR | Hybrid
  algorithm: maglev

# DSR 模式工作原理:
# 入站: Client -> LB Node -> Backend Pod
# 出站: Backend Pod -> Client (绕过 LB)
```

**DSR 优势**:
- 响应流量不经过负载均衡器，降低延迟
- 减少负载均衡器带宽压力
- 特别适合大流量场景

### 2. Bandwidth Manager

```yaml
# 启用带宽管理
bandwidthManager:
  enabled: true

# 应用带宽限制
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: limited-bandwidth
  annotations:
    kubernetes.io/ingress-bandwidth: 100M
    kubernetes.io/egress-bandwidth: 50M
spec:
  containers:
  - name: app
    image: myapp:v1
```

### 3. eBPF Host Routing

```yaml
# 启用 XDP 加速
bpf:
  hostRouting: true
  preallocateMaps: true

# 验证 XDP 状态
# ip link show eth0 | grep xdp
```

## 网络策略（L3-L7）

### L3/L4 策略

```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: allow-frontend-to-backend
  namespace: production
spec:
  endpointSelector:
    matchLabels:
      app: backend
  ingress:
  - fromEndpoints:
    - matchLabels:
        app: frontend
        namespace: production
    toPorts:
    - ports:
      - port: "8080"
        protocol: TCP
```

### L7 HTTP 策略

```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: api-rate-limit
  namespace: production
spec:
  endpointSelector:
    matchLabels:
      app: api-gateway
  ingress:
  - fromEndpoints:
    - matchLabels:
        app: frontend
    toPorts:
    - ports:
      - port: "80"
        protocol: TCP
      rules:
        http:
        - method: GET
          path: "/api/v1/users/.*"
        - method: POST
          path: "/api/v1/orders"
          headers:
          - name: Content-Type
            value: application/json
```

### DNS 感知策略

```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: egress-dns
  namespace: production
spec:
  endpointSelector:
    matchLabels:
      app: microservice
  egress:
  # 允许 DNS 查询
  - toEndpoints:
    - matchLabels:
        k8s:io.kubernetes.pod.namespace: kube-system
        k8s-app: kube-dns
    toPorts:
    - ports:
      - port: "53"
        protocol: UDP
      rules:
        dns:
        - matchPattern: "*.amazonaws.com"
        - matchName: "api.stripe.com"
  
  # 基于 FQDN 的访问控制
  - toFQDNs:
    - matchName: "api.stripe.com"
    toPorts:
    - ports:
      - port: "443"
        protocol: TCP
```

## 可观测性（Hubble）

### Hubble 架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         Hubble UI                               │
│                    (流量可视化界面)                              │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                     Hubble Relay                                │
│              (集群级流量聚合)                                    │
└─────────────────────────────┬───────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼──────┐     ┌────────▼────────┐   ┌───────▼──────┐
│ Hubble Agent │     │  Hubble Agent   │   │ Hubble Agent │
│  (Node 1)    │     │    (Node 2)     │   │   (Node 3)   │
│ ┌──────────┐ │     │  ┌──────────┐   │   │ ┌──────────┐ │
│ │ eBPF Probes│ │    │  │ eBPF Probes│  │   │ │ eBPF Probes│ │
│ │ (kprobe/ │ │     │  │ (kprobe/ │   │   │ │ (kprobe/ │ │
│ │ tracepoint)│ │    │  │ tracepoint)│  │   │ │ tracepoint)│ │
│ └──────────┘ │     │  └──────────┘   │   │ └──────────┘ │
└──────────────┘     └─────────────────┘   └──────────────┘
```

### 启用 Hubble

```yaml
# Cilium values.yaml
hubble:
  enabled: true
  relay:
    enabled: true
  ui:
    enabled: true
  
  # 流量监控
  metrics:
    enabled:
    - dns:query
    - drop
    - tcp
    - flow
    - icmp
    - http
```

### Hubble CLI 使用

```bash
# 安装 hubble CLI
curl -L --remote-name-all https://github.com/cilium/hubble/releases/latest/download/hubble-linux-amd64.tar.gz
tar xzvf hubble-linux-amd64.tar.gz
sudo mv hubble /usr/local/bin/

# 端口转发
kubectl port-forward -n kube-system svc/hubble-relay 4245:4245

# 查看实时流量
hubble observe --follow

# 查看特定 namespace 流量
hubble observe --namespace production

# 查看丢弃的包
hubble observe --verdict DROPPED

# 查看 HTTP 流量详情
hubble observe --protocol http --namespace production
```

### 网络流量可视化

```bash
# 查看服务依赖图
hubble observe --service-dependency

# 查看流统计
hubble observe --statistic
```

## Cluster Mesh（多集群）

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          Cluster Mesh                                   │
│                                                                         │
│  ┌──────────────────────────┐      ┌──────────────────────────┐        │
│  │      Cluster 1 (EU)      │      │      Cluster 2 (US)      │        │
│  │  ┌──────────────────┐   │      │   ┌──────────────────┐   │        │
│  │  │ Cilium ClusterID: 1│   │◄────►│   │ Cilium ClusterID: 2│   │        │
│  │  │ PodCIDR: 10.1.0.0/16│  │      │   │ PodCIDR: 10.2.0.0/16│  │        │
│  │  │ Service: frontend   │  │      │   │ Service: backend   │  │        │
│  │  └──────────────────┘   │      │   └──────────────────┘   │        │
│  │           │              │      │            │              │        │
│  │     ┌─────┴─────┐        │      │      ┌─────┴─────┐        │        │
│  │     ▼           ▼        │      │      ▼           ▼        │        │
│  │  ┌──────┐   ┌──────┐    │      │   ┌──────┐   ┌──────┐    │        │
│  │  │ Pod 1│   │ Pod 2│    │      │   │ Pod 3│   │ Pod 4│    │        │
│  │  └──────┘   └──────┘    │      │   └──────┘   └──────┘    │        │
│  └──────────────────────────┘      └──────────────────────────┘        │
│                                                                         │
│  特性:                                                                   │
│  • 跨集群 Pod IP 可达性                                                  │
│  • 全局服务发现 (Global Services)                                         │
│  • 跨集群网络策略                                                          │
│  • 故障转移支持                                                            │
└─────────────────────────────────────────────────────────────────────────┘
```

### Cluster Mesh 配置

```yaml
# Cluster 1
cluster:
  name: cluster-eu
  id: 1

clustermesh:
  useAPIServer: true
  apiserver:
    service:
      type: LoadBalancer
```

```bash
# 获取 Cluster 1 的连接信息
cilium clustermesh status --context cluster-eu

# 在 Cluster 2 上连接 Cluster 1
cilium clustermesh connect \
  --destination-context cluster-eu \
  --source-context cluster-us
```

### 全局服务

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend
  namespace: production
  annotations:
    io.cilium/global-service: "true"  # 启用全局服务
    io.cilium/service-affinity: "local"  # 优先本地端点
spec:
  type: ClusterIP
  selector:
    app: backend
  ports:
  - port: 80
    targetPort: 8080
```

## 性能基准

### 吞吐量和延迟

| 测试场景 | Kube-Proxy iptables | Kube-Proxy IPVS | Cilium eBPF |
|----------|---------------------|-----------------|-------------|
| TCP 吞吐量 | 基准 | +15% | +30% |
| 延迟 (P99) | 基准 | -10% | -20% |
| CPU 使用率 | 基准 | -15% | -40% |
| 服务更新延迟 | 秒级 | 秒级 | 毫秒级 |

### 大规模集群表现

```bash
# 测试环境: 1000+ 服务, 5000+ Endpoint
# Cilium 性能:
# - 服务查找: O(1), 恒定时间
# - 规则更新: 增量更新，无全表刷新
# - 内存占用: 与端点数量线性关系，高效 eBPF map
```

## 故障排查

```bash
# 检查 Cilium 状态
cilium status

# 检查 eBPF 程序加载
kubectl -n kube-system exec ds/cilium -- cilium bpf lb list

# 查看服务后端
kubectl -n kube-system exec ds/cilium -- cilium service list

# 检查连接跟踪
kubectl -n kube-system exec ds/cilium -- cilium bpf ct list global

# 查看策略状态
kubectl -n kube-system exec ds/cilium -- cilium policy get

# 采集 sysdump（故障报告）
cilium sysdump
```

## 总结

| 场景 | Cilium 优势 |
|------|-------------|
| 大规模集群 | eBPF O(1) 查找，服务数量无影响 |
| 高性能要求 | DSR、XDP 加速、Maglev |
| 安全合规 | L3-L7 策略、DNS 感知、可观测 |
| 多集群 | Cluster Mesh、全局服务 |
| 云原生 | Gateway API、Service Mesh 集成 |

Cilium 通过 eBPF 技术重塑了 Kubernetes 网络，是 2025 年生产环境 CNI 的首选。
