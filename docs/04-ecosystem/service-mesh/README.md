# 服务网格 (Service Mesh)

> 微服务通信的基础设施层

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
│   │Sidecar│   │Sidecar│   │Sidecar│              │
│   │ Envoy │   │ Envoy │   │ Envoy │              │
│   └───┬────┘   └───┬────┘   └───┬────┘              │
│       └─────────────┴─────────────┘                    │
│              服务发现/负载均衡/安全/可观测性            │
└─────────────────────────────────────────────────────────┘
```

---

## 主流服务网格对比

| 特性 | Istio | Linkerd | Consul Connect |
|------|-------|---------|----------------|
| **数据平面** | Envoy | Linkerd-proxy | Envoy |
| **资源消耗** | 高 | 低 | 中 |
| **功能丰富度** | ★★★★★ | ★★★★ | ★★★ |
| **易用性** | ★★★ | ★★★★★ | ★★★★ |
| **社区活跃** | 高 | 中 | 中 |

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

## 选择建议

| 场景 | 推荐 |
|------|------|
| 大型企业 | Istio |
| 资源受限 | Linkerd |
| 混合云 | Consul Connect |
| 简单需求 | K8s Service + Ingress |
