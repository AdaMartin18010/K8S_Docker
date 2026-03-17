# 服务网格 (Service Mesh)

> 微服务通信的基础设施层 - Istio Ambient Mesh GA (2025)

---

## 什么是服务网格？

服务网格是处理服务间通信的基础设施层，将服务发现、负载均衡、故障恢复、指标收集等功能从应用代码中剥离出来。

```
┌─────────────────────────────────────────────────────────┐
│                      应用层                              │
│   ┌─────────┐   ┌─────────┐   ┌─────────┐              │
│   │ Service │   │ Service │   │ Service │              │
│   │    A    │   │    B    │   │    C    │              │
│   └────┬────┘   └────┬────┘   └────┬────┘              │
├────────┼─────────────┼─────────────┼────────────────────┤
│        │             │             │      服务网格层     │
│   ┌────▼────┐   ┌────▼────┐   ┌────▼────┐              │
│   │ Sidecar │   │ Sidecar │   │ Sidecar │              │
│   │ Envoy   │   │ Envoy   │   │ Envoy   │              │
│   └───┬────┘   └───┬────┘   └───┬────┘              │
│       └─────────────┴─────────────┘                    │
│              服务发现/负载均衡/安全/可观测性            │
└─────────────────────────────────────────────────────────┘
```

---

## 2025 重大更新: Istio Ambient Mesh GA

Istio 宣布 **Ambient Mesh GA** (2025年10月)，这是服务网格架构的重大演进。

### 传统 Sidecar vs Ambient Mesh

```
传统 Sidecar 模式:
┌─────────┐    ┌─────────┐    ┌─────────┐
│   App   │◄──►│  Envoy  │◄──►│  Network │
└─────────┘    │ Sidecar │    └─────────┘
               └─────────┘

Ambient Mesh 模式:
┌─────────┐    ┌─────────┐    ┌─────────┐
│   App   │◄──►│ ztunnel │◄──►│ Waypoint │
└─────────┘    │  (L4)   │    │  (L7)   │
               └─────────┘    └─────────┘
```

### Ambient Mesh 优势

| 指标 | Sidecar 模式 | Ambient Mesh |
|------|-------------|--------------|
| **CPU 开销** | 基准 | **-35%** |
| **延迟** | 基准 | **-15%** |
| **内存占用** | 高 (每 Pod) | **低 (每节点)** |
| **升级影响** | 需重启 Pod | **零停机** |
| **资源隔离** | Pod 级 | **节点级 L4 + 可选 L7** |

### Ambient Mesh 架构

```
┌─────────────────────────────────────────────────────────┐
│                    Ambient Mesh 架构                     │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌─────────┐      ┌─────────┐      ┌─────────┐         │
│  │   App   │◄────►│ ztunnel │◄────►│ Waypoint │        │
│  │ (无代理) │      │ (L4/节点)│      │ (L7/可选) │        │
│  └─────────┘      └────┬────┘      └─────────┘         │
│                        │                                 │
│                        ▼                                 │
│              ┌─────────────────────┐                    │
│              │    HBONE 隧道       │                    │
│              │  (HTTP-Based Overlay)│                    │
│              └─────────────────────┘                    │
│                                                          │
│  ztunnel: 轻量级 L4 代理 (每个节点)                        │
│  Waypoint: 可选 L7 代理 (按命名空间/服务)                  │
│  HBONE: 基于 HTTP/2 的安全隧道                            │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 启用 Ambient Mesh

```bash
# 安装 Istio with Ambient 模式
istioctl install --set profile=ambient --skip-confirmation

# 标记命名空间加入 Ambient Mesh
kubectl label namespace default istio.io/dataplane-mode=ambient

# 为特定服务启用 L7 处理 (Waypoint Proxy)
istioctl waypoint apply --namespace default --name reviews-svc-waypoint
kubectl label service reviews istio.io/use-waypoint=reviews-svc-waypoint
```

### Waypoint Proxy 配置

```yaml
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: reviews-waypoint
  namespace: default
  annotations:
    istio.io/for-service-account: reviews
spec:
  gatewayClassName: istio-waypoint
  listeners:
    - name: mesh
      port: 15008
      protocol: HBONE
```

---

## 主流服务网格对比 (2025)

| 特性 | Istio (Ambient) | Linkerd 3.x | Cilium Service Mesh | Consul Connect |
|------|-----------------|-------------|---------------------|----------------|
| **数据平面** | ztunnel + Waypoint | Linkerd2-proxy | eBPF + Envoy | Envoy |
| **Sidecar-less** | ✅ **是** | 否 | ✅ **是** | 否 |
| **资源消耗** | **低** | 低 | **最低** | 中 |
| **功能丰富度** | ★★★★★ | ★★★★ | ★★★★ | ★★★ |
| **易用性** | ★★★ | ★★★★★ | ★★★★ | ★★★★ |
| **eBPF 优化** | 否 | 否 | **是** | 否 |
| **Gateway API** | ✅ **原生** | 部分 | ✅ **原生** | 否 |

---

## Cilium Service Mesh (无 Sidecar)

Cilium 提供基于 eBPF 的无 Sidecar 服务网格方案。

```yaml
# 启用 Cilium Service Mesh
apiVersion: cilium.io/v2alpha1
kind: CiliumClusterwideEnvoyConfig
metadata:
  name: envoy-lb
spec:
  services:
    - name: reviews
      namespace: default
  resources:
    - "@type": type.googleapis.com/envoy.config.listener.v3.Listener
      name: envoy-lb-listener
      address:
        socket_address:
          address: "0.0.0.0"
          port_value: 8080
```

---

## Istio 核心功能

### 流量管理

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
    - reviews
  http:
    - route:
        - destination:
            host: reviews
            subset: v1
          weight: 75
        - destination:
            host: reviews
            subset: v2
          weight: 25
```

### 安全通信

```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
spec:
  mtls:
    mode: STRICT
```

---

## 选择建议 (2025)

| 场景 | 推荐 | 原因 |
|------|------|------|
| **性能优先** | Cilium Service Mesh | eBPF 加速，无 Sidecar |
| **功能丰富** | Istio Ambient | 成熟生态，Ambient 降低开销 |
| **简单轻量** | Linkerd | 最易用，资源占用低 |
| **混合云** | Consul Connect | 多云服务发现 |
| **已有 Cilium** | Cilium Service Mesh | 无需额外组件 |

---

## 2025 趋势

- ✅ **Ambient Mesh**: Istio 无 Sidecar 方案 GA
- ✅ **eBPF 服务网格**: Cilium 方案成熟
- 🔄 **Gateway API 统一**: 替代 Ingress + 部分服务网格功能
- 🔄 **AI 流量管理**: LLM 推理服务的特殊路由需求
- 🔄 **零信任默认**: mTLS 成为标配而非选配

---

## 参考

- [Istio Ambient Mesh 文档](https://istio.io/latest/docs/ambient/)
- [Cilium Service Mesh](https://docs.cilium.io/en/stable/network/servicemesh/)
- [Linkerd 文档](https://linkerd.io/)
