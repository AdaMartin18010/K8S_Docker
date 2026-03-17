# Flagger - 渐进式交付 Kubernetes Operator

## 概述

Flagger 是一个渐进式交付工具，自动化 Kubernetes 应用的发布流程。它通过逐渐将流量转移到新版本，同时测量指标和运行一致性测试，降低在生产环境中引入新软件版本的风险。

## 核心特性

| 特性 | 描述 |
|------|------|
| 金丝雀发布 | 渐进式流量切换，支持会话亲和性 |
| A/B 测试 | 基于 HTTP 头和 cookie 的流量路由 |
| 蓝绿部署 | 流量切换和镜像测试 |
| 自动回滚 | 基于指标阈值自动回滚 |
| 多网格支持 | Istio、Linkerd、Kuma 等服务网格 |

## 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                     Flagger 架构                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │  Primary    │    │   Canary    │    │  Load       │         │
│  │  (稳定版)    │    │   (新版本)   │    │  Tester     │         │
│  └──────┬──────┘    └──────┬──────┘    └─────────────┘         │
│         │                  │                                     │
│         └──────────┬───────┘                                     │
│                    ▼                                            │
│         ┌─────────────────────┐                                │
│         │   VirtualService    │                                │
│         │   (流量分割)         │                                │
│         │  Primary: 90%       │                                │
│         │  Canary: 10%        │                                │
│         └─────────────────────┘                                │
│                    │                                            │
│                    ▼                                            │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Prometheus                            │   │
│  │  - 请求成功率 (request-success-rate)                     │   │
│  │  - 请求延迟 (request-duration)                           │   │
│  │  - 自定义指标                                            │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   Flagger Controller                     │   │
│  │  1. 检测 Deployment 变更                                  │   │
│  │  2. 创建 Canary 资源                                      │   │
│  │  3. 渐进式流量切换                                        │   │
│  │  4. 指标分析                                              │   │
│  │  5. 晋升或回滚                                            │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 安装部署

### Istio + Flagger 安装

```bash
# 安装 Istio
istioctl install --set profile=default -y

# 安装 Flagger
helm repo add flagger https://flagger.app
helm repo update

helm install flagger flagger/flagger \
  --namespace istio-system \
  --set meshProvider=istio \
  --set metricsServer=http://prometheus:9090

# 安装负载测试工具
helm install flagger-loadtester flagger/loadtester \
  --namespace test
```

### Gateway API 支持

```bash
# 安装 Gateway API CRDs
kubectl apply --server-side -k "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v1.4.0"

# 安装 Flagger (Gateway API 模式)
helm upgrade -i flagger flagger/flagger \
  --namespace flagger-system \
  --create-namespace \
  --set meshProvider=gatewayapi:v1 \
  --set metricsServer=http://prometheus.monitoring:9090
```

## Canary 配置

### 基本金丝雀发布

```yaml
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: my-app
  namespace: production
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  service:
    port: 80
    targetPort: 8080
    gateways:
    - mesh
    hosts:
    - my-app
  analysis:
    interval: 30s          # 每 30 秒检查一次
    threshold: 5           # 5 次失败后回滚
    maxWeight: 50          # 最大金丝雀流量 50%
    stepWeight: 10         # 每次增加 10%
    metrics:
    - name: request-success-rate
      thresholdRange:
        min: 99            # 成功率至少 99%
      interval: 1m
    - name: request-duration
      thresholdRange:
        max: 500           # P99 延迟小于 500ms
      interval: 1m
    webhooks:
    - name: load-test
      url: http://flagger-loadtester.test/
      timeout: 5s
      metadata:
        cmd: "hey -z 1m -q 10 -c 2 http://my-app-canary.production:80/"
```

### A/B 测试

```yaml
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: my-app
  namespace: production
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  service:
    port: 80
    targetPort: 8080
  analysis:
    interval: 1m
    threshold: 5
    iterations: 10
    match:
    - headers:
        x-canary:
          exact: "insider"
    - cookies:
        canary:
          exact: "always"
    metrics:
    - name: request-success-rate
      thresholdRange:
        min: 99
    webhooks:
    - name: load-test
      url: http://flagger-loadtester.test/
      metadata:
        cmd: "hey -z 1m -q 10 -c 2 -H 'x-canary: insider' http://my-app.production:80/"
```

### 蓝绿部署

```yaml
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: my-app
  namespace: production
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  service:
    port: 80
    targetPort: 8080
  analysis:
    interval: 30s
    threshold: 2
    iterations: 1          # 蓝绿只需一次迭代
    mirror: true           # 镜像流量到新版本
    mirrorWeight: 100
    metrics:
    - name: request-success-rate
      thresholdRange:
        min: 99
  promote:
    enable: true
    analysis:
      interval: 30s
      threshold: 2
```

## Gateway API 集成

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: external-gateway
  namespace: istio-ingress
spec:
  gatewayClassName: istio
  listeners:
  - name: default
    hostname: "*.example.com"
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: All
---
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: my-app
  namespace: production
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  service:
    port: 80
    targetPort: 8080
    gatewayRefs:
    - name: external-gateway
      namespace: istio-ingress
  analysis:
    interval: 30s
    threshold: 5
    maxWeight: 50
    stepWeight: 10
    metrics:
    - name: request-success-rate
      thresholdRange:
        min: 99
```

## 自定义指标

```yaml
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: my-app
  namespace: production
spec:
  analysis:
    metrics:
    - name: "5xx errors"
      templateRef:
        name: 5xx-error-rate
        namespace: flagger
      thresholdRange:
        max: 10
    - name: "custom metric"
      templateRef:
        name: custom-metric
      thresholdRange:
        min: 100
---
apiVersion: flagger.app/v1beta1
kind: MetricTemplate
metadata:
  name: 5xx-error-rate
  namespace: flagger
spec:
  provider:
    type: prometheus
    address: http://prometheus:9090
  query: |
    sum(rate(http_requests_total{status=~"5.."}[1m]))
    /
    sum(rate(http_requests_total[1m])) * 100
```

## 告警通知

```yaml
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: my-app
  namespace: production
spec:
  analysis:
    alerts:
    - name: "slack"
      severity: info
      providerRef:
        name: slack
        namespace: flagger
---
apiVersion: flagger.app/v1beta1
kind: AlertProvider
metadata:
  name: slack
  namespace: flagger
spec:
  type: slack
  channel: deployments
  username: flagger
  webhook: https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX
```

## 相关资源

- [Flagger 官网](https://flagger.app/)
- [GitHub](https://github.com/fluxcd/flagger)
- [Flux 文档](https://fluxcd.io/flagger/)
