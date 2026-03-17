# Istio 服务网格指南

> 使用 Istio 管理服务间通信

---

## 什么是服务网格？

服务网格是处理服务间通信的基础设施层，将服务治理能力下沉到基础设施：

```
┌─────────────────────────────────────────────┐
│           应用层 (Application Layer)          │
│     ┌─────────┐      ┌─────────┐           │
│     │ Service A│◀────▶│ Service B│           │
│     └────┬────┘      └────┬────┘           │
│          │                │                 │
├──────────┼────────────────┼─────────────────┤
│          │   服务网格层    │                 │
│     ┌────▼────┐      ┌────▼────┐           │
│     │Sidecar  │◀────▶│Sidecar  │           │
│     │(Envoy) │      │(Envoy) │           │
│     └────┬────┘      └────┬────┘           │
├──────────┼────────────────┼─────────────────┤
│          │   网络层       │                 │
│          └────────────────┘                 │
└─────────────────────────────────────────────┘
```

---

## Istio 核心功能

| 功能 | 说明 |
|------|------|
| **流量管理** | 路由、分流、熔断、重试 |
| **安全** | mTLS、认证、授权 |
| **可观测性** | 指标、日志、追踪 |

---

## 安装 Istio

```bash
# 下载 Istio
curl -L https://istio.io/downloadIstio | sh -
cd istio-*/

# 安装（demo 配置）
./bin/istioctl install --set profile=demo -y

# 启用自动 Sidecar 注入
kubectl label namespace default istio-injection=enabled
```

---

## 流量管理示例

### VirtualService - 路由规则

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
          weight: 50
        - destination:
            host: reviews
            subset: v2
          weight: 50
```

### DestinationRule - 目标规则

```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: reviews
spec:
  host: reviews
  subsets:
    - name: v1
      labels:
        version: v1
    - name: v2
      labels:
        version: v2
```

---

## 安全 - mTLS

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

## 可观测性

Istio 自动收集以下指标：
- 请求量
- 延迟分布
- 错误率
- 拓扑图
