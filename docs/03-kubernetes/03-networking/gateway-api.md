# Gateway API

> Kubernetes 下一代流量管理标准

---

## Gateway API vs Ingress

```
┌─────────────────────────────────────────────────────────────┐
│              Ingress vs Gateway API                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Ingress (传统)                    Gateway API (新)           │
│  ────────────────────────────────────────────────────────   │
│  • 单一资源定义所有配置            • 多角色分离               │
│  • 实现特定注解                    • 可移植标准               │
│  • 功能有限                        • 功能丰富 (L4/L7/TLS)     │
│  • 不支持 TCP/UDP                  • 支持多协议               │
│                                                              │
│  角色模型:                                                    │
│  ┌─────────────┐                ┌─────────────┐              │
│  │  平台团队    │  管理 Gateway  │  ClusterAdmin│              │
│  │  (基础设施)  │                │              │              │
│  ├─────────────┤                ├─────────────┤              │
│  │  应用团队    │  定义路由规则   │  HTTPRoute   │              │
│  │  (路由)     │                │  TCPRoute    │              │
│  └─────────────┘                └─────────────┘              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Gateway API 资源

| 资源 | 用途 | 管理者 |
|------|------|--------|
| **GatewayClass** | 定义网关类型 | 基础设施团队 |
| **Gateway** | 网关实例，监听端口 | 平台团队 |
| **HTTPRoute** | HTTP 路由规则 | 应用团队 |
| **TCPRoute** | TCP 路由规则 | 应用团队 |
| **TLSRoute** | TLS 路由规则 | 应用团队 |
| **GRPCRoute** | gRPC 路由规则 | 应用团队 |
| **ReferenceGrant** | 跨命名空间引用授权 | 安全团队 |

---

## 基本配置

### GatewayClass

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: nginx
spec:
  controllerName: gateway.nginx.org/nginx-gateway-controller
  description: "NGINX Gateway Controller"
```

### Gateway

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: production-gateway
  namespace: infrastructure
spec:
  gatewayClassName: nginx
  listeners:
    - name: http
      protocol: HTTP
      port: 80
      allowedRoutes:
        namespaces:
          from: Selector
          selector:
            matchLabels:
              gateway-access: "true"
    - name: https
      protocol: HTTPS
      port: 443
      tls:
        mode: Terminate
        certificateRefs:
          - name: production-cert
      allowedRoutes:
        namespaces:
          from: All
```

### HTTPRoute

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: api-routes
  namespace: production
  labels:
    gateway-access: "true"
spec:
  parentRefs:
    - name: production-gateway
      namespace: infrastructure
      sectionName: https
  hostnames:
    - api.example.com
  rules:
    # 路径匹配
    - matches:
        - path:
            type: PathPrefix
            value: /v1/users
      backendRefs:
        - name: user-service
          port: 8080

    # Header 匹配
    - matches:
        - headers:
            - name: version
              value: v2
      backendRefs:
        - name: user-service-v2
          port: 8080

    # 权重分流 (金丝雀)
    - backendRefs:
        - name: api-stable
          port: 8080
          weight: 90
        - name: api-canary
          port: 8080
          weight: 10
```

---

## 高级功能

### 流量镜像

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: traffic-mirror
spec:
  parentRefs:
    - name: production-gateway
  rules:
    - backendRefs:
        - name: production-service
          port: 8080
      filters:
        - type: RequestMirror
          requestMirror:
            backendRef:
              name: staging-service
              port: 8080
```

### 请求改写

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: rewrite-route
spec:
  parentRefs:
    - name: production-gateway
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /api/v2
      filters:
        - type: URLRewrite
          urlRewrite:
            path:
              type: ReplacePrefixMatch
              replacePrefixMatch: /v2
      backendRefs:
        - name: backend-service
          port: 8080
```

### 跨命名路由

```yaml
# 授权应用团队使用 Gateway
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: allow-gateway-ref
  namespace: infrastructure
spec:
  from:
    - group: gateway.networking.k8s.io
      kind: HTTPRoute
      namespace: team-a
  to:
    - group: gateway.networking.k8s.io
      kind: Gateway
      name: production-gateway
```

---

## 实现对比

| 实现 | 特点 | 成熟度 |
|------|------|--------|
| **NGINX Gateway** | NGINX 官方，功能全面 | GA |
| **Traefik** | 云原生，易用 | GA |
| **Istio** | 服务网格集成 | GA |
| **Cilium** | eBPF 加速 | GA |
| **Envoy Gateway** | Envoy 官方 | GA |
| **Contour** | Envoy 基础 | GA |

---

## 从 Ingress 迁移

```bash
# 使用 ingress2gateway 工具
kubectl get ingress my-ingress -o yaml > ingress.yaml
ingress2gateway -f ingress.yaml -o gateway.yaml
kubectl apply -f gateway.yaml
```

---

## 2025 趋势

- **GA 稳定版**: Gateway API v1.1+ 已生产就绪
- **服务网格整合**: Istio、Cilium 原生支持
- **多集群网关**: 跨集群流量管理
- **Gateway API for Mesh**: 服务网格标准化
