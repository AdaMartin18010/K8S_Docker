# Gateway API

> Kubernetes 下一代流量管理标准 - v1.2/v1.3 新特性

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
| **BackendTLSPolicy** | 后端 TLS 配置 (v1.2+) | 安全团队 |

---

## Gateway API v1.2 新特性 (2024年10月)

### HTTPRoute 超时 (GA)

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: api-with-timeout
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
      port: 8080
    timeouts:
      request: 300ms      # 请求超时
      backendRequest: 200ms  # 后端连接超时
```

### 基础设施标签和注解 (GA)

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: meshed-gateway
spec:
  gatewayClassName: cilium
  listeners:
  - name: http-listener
    protocol: HTTP
    port: 80
  infrastructure:
    labels:
      istio-injection: enabled
    annotations:
      linkerd.io/inject: enabled
```

### HTTPRoute 重试 (Experimental)

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: retry-route
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
      port: 8080
    retry:
      codes: [500, 502, 503, 504]
      attempts: 3
      backoff: 500ms
```

### 百分比流量镜像 (Experimental)

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: mirror-route
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
            percent: 42   # 42% 流量镜像
            # 或使用 fraction 更精确控制
            # fraction:
            #   numerator: 1
            #   denominator: 10000
```

---

## Gateway API v1.3 新特性 (2025年6月)

### CORS 过滤器 (Experimental)

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: cors-route
spec:
  parentRefs:
    - name: production-gateway
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /api
    filters:
    - type: CORS
      cors:
        allowOrigins:
          - "https://app.example.com"
          - "https://admin.example.com"
        allowMethods:
          - GET
          - POST
          - PUT
          - DELETE
        allowHeaders:
          - Authorization
          - Content-Type
        allowCredentials: true
        maxAge: 86400
    backendRefs:
    - name: api-service
      port: 8080
```

### 重试预算 (Experimental)

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: retry-budget-route
spec:
  parentRefs:
    - name: production-gateway
  rules:
  - backendRefs:
    - name: api-service
      port: 8080
    retry:
      budget:
        percent: 20      # 最多 20% 请求可重试
        minRetries: 10   # 最少重试次数
```

### 后端协议支持 (GA)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: websocket-service
spec:
  selector:
    app: websocket-app
  ports:
    - name: http
      port: 80
      targetPort: 9376
      protocol: TCP
      appProtocol: kubernetes.io/ws  # WebSocket 支持
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: ws-route
spec:
  parentRefs:
    - name: production-gateway
  rules:
  - backendRefs:
    - name: websocket-service
      port: 80
```

支持的后端协议:

- `kubernetes.io/h2c` - HTTP/2 明文
- `kubernetes.io/ws` - WebSocket 明文
- `kubernetes.io/wss` - WebSocket over TLS

---

## 基本配置

### GatewayClass

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: cilium
spec:
  controllerName: io.cilium/gateway-controller
  description: "Cilium Gateway Controller"
```

### Gateway

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: production-gateway
  namespace: infrastructure
spec:
  gatewayClassName: cilium
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

### BackendTLSPolicy (v1.2+)

```yaml
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: BackendTLSPolicy
metadata:
  name: backend-tls
spec:
  targetRef:
    group: ""
    kind: Service
    name: backend-service
  tls:
    certificateRefs:
      - name: backend-cert
    subjectAltNames:
      - "backend.example.com"
      - "spiffe://cluster.local/ns/backend/sa/default"
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

| 实现 | 特点 | 成熟度 | v1.2 支持 |
|------|------|--------|----------|
| **Cilium** | eBPF 加速，kube-proxy 替代 | GA | ✅ 实验特性 |
| **Envoy Gateway** | Envoy 官方，AI Gateway 集成 | GA | ✅ 实验特性 |
| **Istio** | 服务网格集成 | GA | ✅ 实验特性 |
| **NGINX Gateway** | NGINX 官方 | GA | ✅ 标准特性 |
| **Traefik** | 云原生，易用 | GA | ✅ 实验特性 |
| **Kong** | API 管理集成 | GA | ✅ 实验特性 |

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

- ✅ **Gateway API v1.2 GA**: WebSocket、超时、重试等功能稳定
- ✅ **Gateway API v1.3**: CORS、重试预算、百分比镜像
- 🔄 **服务网格整合**: Istio、Cilium 原生支持 Gateway API for Mesh
- 🔄 **多集群网关**: 跨集群流量管理
- 🔄 **AI Gateway**: Envoy AI Gateway 与 KServe 集成
- 🔄 **gwctl CLI**: 独立仓库，更好的 Gateway 管理工具

---

## 参考

- [Gateway API 官方文档](https://gateway-api.sigs.k8s.io/)
- [Gateway API v1.2 发布说明](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.2.0)
- [Gateway API v1.3 发布说明](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.3.0)
- [Envoy AI Gateway](https://aigateway.envoyproxy.io/)
