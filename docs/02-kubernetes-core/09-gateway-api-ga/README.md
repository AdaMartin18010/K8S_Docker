# Gateway API GA - 下一代 K8s 入口管理

## 概述

Kubernetes Gateway API 于 2025 年 10 月达到 GA (Generally Available)，正式成为 Ingress 的继任者，提供更灵活、更强大的流量管理能力。

> **重要里程碑**: Gateway API GA 标志着 Ingress 即将退役，Ingress NGINX 项目已进入维护模式。

## Gateway API vs Ingress

### 架构对比

**Ingress 模式**:

- 单一资源混合配置（路由 + TLS + 基础设施）
- 实现依赖注解，可移植性差
- 缺乏角色分离（平台 vs 应用团队）

**Gateway API 模式**:

- 资源分离（GatewayClass → Gateway → Route）
- 标准化配置，实现无关
- 清晰的 RBAC 边界

### 核心资源

| 资源 | 作用 | 管理者 |
|------|------|--------|
| GatewayClass | 控制器声明 | 集群管理员 |
| Gateway | 入口点配置 | 平台团队 |
| HTTPRoute | L7 路由规则 | 应用团队 |
| TCPRoute/TLSRoute | L4 路由 | 应用团队 |

## 基础配置示例

### 1. GatewayClass

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: envoy
spec:
  controllerName: gateway.envoyproxy.io/gatewayclass-controller
  description: "Envoy Gateway"
```

### 2. Gateway

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: production
  namespace: ingress
spec:
  gatewayClassName: envoy
  listeners:
  - name: http
    protocol: HTTP
    port: 80
    allowedRoutes:
      namespaces:
        from: All
  - name: https
    protocol: HTTPS
    port: 443
    tls:
      mode: Terminate
      certificateRefs:
      - name: prod-cert
    allowedRoutes:
      namespaces:
        from: Selector
        selector:
          matchLabels:
            gateway-access: "true"
```

### 3. HTTPRoute

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: api
  namespace: app
spec:
  parentRefs:
  - name: production
    namespace: ingress
  hostnames:
  - api.example.com
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /v1
    backendRefs:
    - name: api-v1
      port: 8080
      weight: 90
    - name: api-v2
      port: 8080
      weight: 10
```

## Ingress 迁移指南

### 使用 ingress2gateway 工具

```bash
# 安装
curl -L https://github.com/kubernetes-sigs/ingress2gateway/releases/latest/download/ingress2gateway_$(uname -s)_$(uname -m).tar.gz | tar -xz

# 转换
./ingress2gateway convert --all-namespaces --providers=ingress-nginx
```

### 关键差异

| 特性 | Ingress | Gateway API |
|------|---------|-------------|
| 注解配置 | 必需 | 无需 |
| 跨 Namespace | 不支持 | 原生支持 |
| 流量分割 | 依赖注解 | 原生 weight |
| TLS 配置 | 混合 | 分离到 Gateway |
| 协议支持 | HTTP/HTTPS | HTTP/TCP/TLS/UDP |

## 实现选择

| 实现 | 特点 |
|------|------|
| Envoy Gateway | CNCF 官方，功能全面 |
| Istio | 服务网格集成 |
| Cilium | eBPF 高性能 |
| NGINX Gateway Fabric | 平滑迁移 |

## 2025 最佳实践

1. **新集群直接使用 Gateway API**
2. **存量集群渐进迁移**（并行运行）
3. **使用 GatewayClass 参数化配置**
4. **利用 HTTPRoute 权重实现金丝雀发布**

Gateway API GA 标志着 Kubernetes 流量管理的现代化统一标准。
