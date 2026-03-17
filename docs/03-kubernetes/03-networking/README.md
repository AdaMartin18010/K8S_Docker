# Kubernetes 网络

> K8s 网络模型、CNI 与 Gateway API (1.33/1.34)

---

## K8s 网络模型

```
┌─────────────────────────────────────────────────────────────────┐
│                     Kubernetes 网络模型                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1. Pod 内所有容器共享网络命名空间 (localhost 互通)              │
│   2. Pod 之间可以不经过 NAT 直接通信                             │
│   3. 节点可以与 Pod 直接通信                                     │
│   4. 每个 Pod 有独立的 IP (扁平网络)                             │
│                                                                  │
│  ┌──────────┐         ┌──────────┐         ┌──────────┐        │
│  │   Pod    │◄───────►│   Pod    │◄───────►│   Pod    │        │
│  │ 10.0.1.2 │         │ 10.0.1.3 │         │ 10.0.2.4 │        │
│  └────┬─────┘         └────┬─────┘         └────┬─────┘        │
│       │                    │                    │              │
│       └────────────────────┼────────────────────┘              │
│                            │                                   │
│                    ┌───────▼───────┐                           │
│                    │    CNI        │                           │
│                    │  (Flannel/    │                           │
│                    │   Calico/     │                           │
│                    │   Cilium)     │                           │
│                    └───────────────┘                           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## CNI 对比

| CNI | 特点 | 性能 | 适用场景 |
|-----|------|------|----------|
| **Cilium** | eBPF, L7策略, Hubble | ⭐⭐⭐⭐⭐ | 大规模生产环境 |
| **Calico** | BGP, iptables/eBPF | ⭐⭐⭐⭐ | 传统网络环境 |
| **Flannel** | 简单, VXLAN | ⭐⭐⭐ | 小规模测试 |
| **Weave** | 加密, 去中心化 | ⭐⭐⭐ | 多云环境 |

---

## Service 类型

```yaml
# ClusterIP (集群内部访问)
apiVersion: v1
kind: Service
metadata:
  name: backend
spec:
  type: ClusterIP
  selector:
    app: backend
  ports:
  - port: 80
    targetPort: 8080
---
# NodePort (节点端口暴露)
apiVersion: v1
kind: Service
metadata:
  name: web
spec:
  type: NodePort
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 8080
    nodePort: 30080  # 30000-32767
---
# LoadBalancer (云厂商负载均衡)
apiVersion: v1
kind: Service
metadata:
  name: api
spec:
  type: LoadBalancer
  selector:
    app: api
  ports:
  - port: 443
    targetPort: 8443
---
# ExternalName (DNS CNAME)
apiVersion: v1
kind: Service
metadata:
  name: external-db
spec:
  type: ExternalName
  externalName: db.example.com
```

---

## Ingress vs Gateway API

| 特性 | Ingress | Gateway API |
|------|---------|-------------|
| 标准 | 1.0 GA | v1.2 GA (2025) |
| 功能 | 基础路由 | 高级流量管理 |
| 扩展性 | 有限 | 强 (路由多路复用) |
| 多租户 | 弱 | 强 (Gateway/Route 分离) |

### Gateway API 示例

```yaml
# Gateway
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: production-gateway
spec:
  gatewayClassName: cilium
  listeners:
  - name: http
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: Same
---
# HTTPRoute
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: api-route
spec:
  parentRefs:
  - name: production-gateway
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /api
    backendRefs:
    - name: api-service
      port: 80
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: frontend-service
      port: 80
```

---

## 高级路由（Gateway API v1.2+）

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: advanced-route
spec:
  parentRefs:
  - name: production-gateway
  hostnames:
  - api.example.com
  rules:
  # 请求头路由
  - matches:
    - headers:
      - name: x-canary
        value: "true"
    backendRefs:
    - name: api-canary
      port: 80
  # 路径路由 + 重写
  - matches:
    - path:
        type: PathPrefix
        value: /v1
    filters:
    - type: URLRewrite
      urlRewrite:
        path:
          type: ReplacePrefixMatch
          replacePrefixMatch: /
    backendRefs:
    - name: api-v1
      port: 80
  # 流量分割
  - backendRefs:
    - name: api-stable
      port: 80
      weight: 90
    - name: api-canary
      port: 80
      weight: 10
```

---

## 网络策略

```yaml
# 隔离命名空间
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
---
# 允许 frontend 访问 backend
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: backend-policy
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
---
# Cilium L7 策略
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: api-l7-policy
spec:
  endpointSelector:
    matchLabels:
      app: api
  ingress:
  - fromEndpoints:
    - matchLabels:
        app: web
    toPorts:
    - ports:
      - port: "80"
        protocol: TCP
      rules:
        http:
        - method: GET
          path: "/api/.*"
```

---

## DNS 与服务发现

```bash
# Pod DNS 格式
<pod-ip>.<namespace>.pod.cluster.local
# 例如: 10-0-1-2.default.pod.cluster.local

# Service DNS 格式
<service>.<namespace>.svc.cluster.local
# 例如: mysql.database.svc.cluster.local

# 集群外 DNS
<service>.<namespace>.svc.<cluster-domain>
```

```yaml
# CoreDNS 自定义配置
apiVersion: v1
kind: ConfigMap
metadata:
  name: coredns-custom
  namespace: kube-system
data:
  example.server: |
    example.com {
      forward . 192.168.1.1
    }
```

---

## K8s 1.33/1.34 网络新特性

| 特性 | 状态 | 说明 |
|------|------|------|
| **Gateway API v1.2** | GA | 超时、重试、CORS |
| **Pod 级网络隔离** | Beta | 单个 Pod 网络策略 |
| **Service 内部流量策略** | GA | 拓扑感知路由 |
| **EndpointSlice** | GA | 替代 Endpoints |

---

## 故障排查

```bash
# 测试 Pod 连通性
kubectl run debug --rm -i --restart=Never --image=nicolaka/netshoot -- ping <pod-ip>

# 查看 DNS 解析
kubectl run debug --rm -i --restart=Never --image=busybox -- nslookup kubernetes.default

# 查看网络策略
kubectl get networkpolicy

# 查看 EndpointSlice
kubectl get endpointslices

# Cilium 排查
kubectl exec -n kube-system cilium-<xxx> -- cilium endpoint list
kubectl exec -n kube-system cilium-<xxx> -- cilium policy get
```

---

## 相关文档

- [Gateway API 详解](gateway-api.md)
- [Cilium eBPF 网络](docs/04-ecosystem/ebpf-cilium/)
