# Kubernetes 网络

> 服务发现、负载均衡与网络策略

---

## 本章内容

1. [容器网络模型 (CNM)](./container-network-model.md)
2. [Service 服务发现](./services.md)
3. [Ingress 流量入口](./ingress.md)
4. [NetworkPolicy 网络策略](./network-policies.md)
5. [DNS 与 CoreDNS](./dns.md)

---

## K8s 网络核心原则

1. **每个 Pod 有独立 IP**: Pod IP 在集群内可达
2. **Pod 间直接通信**: 无需 NAT
3. **Service 虚拟 IP**: 提供稳定的服务端点

---

## Service 类型对比

| 类型 | 说明 | 使用场景 |
|------|------|----------|
| **ClusterIP** | 集群内部访问 | 微服务间通信 |
| **NodePort** | 暴露节点端口 | 开发测试 |
| **LoadBalancer** | 云负载均衡 | 生产环境 |
| **ExternalName** | 外部服务映射 | 外部依赖 |

---

## Service 配置示例

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-app
spec:
  type: ClusterIP
  selector:
    app: web-app
  ports:
    - name: http
      port: 80
      targetPort: 8080
      protocol: TCP
  sessionAffinity: None
```

---

## Ingress 配置 (新版)

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: web-app
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - app.example.com
      secretName: app-tls
  rules:
    - host: app.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: web-app
                port:
                  number: 80
```

---

## NetworkPolicy 最佳实践

```yaml
# 默认拒绝所有入站
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
spec:
  podSelector: {}
  policyTypes:
    - Ingress
---
# 允许特定流量
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-web
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
      ports:
        - protocol: TCP
          port: 8080
```

---

## 2025 年新趋势：Gateway API

Gateway API 是 Ingress 的下一代替代方案：

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: example-gateway
spec:
  gatewayClassName: nginx
  listeners:
    - name: http
      protocol: HTTP
      port: 80
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: example-route
spec:
  parentRefs:
    - name: example-gateway
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /api
      backendRefs:
        - name: api-service
          port: 80
```
