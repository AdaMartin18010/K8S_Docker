# K8s 部署策略

> 滚动更新、金丝雀、蓝绿部署与 A/B 测试 (2025)

---

## 策略对比

| 策略 | 复杂度 | 风险 | 适用场景 |
|------|--------|------|----------|
| **滚动更新** | 低 | 中 | 常规发布 |
| **金丝雀** | 中 | 低 | 关键业务 |
| **蓝绿** | 中 | 低 | 金融支付 |
| **A/B 测试** | 高 | 中 | 产品验证 |

---

## 1. 滚动更新（Rolling Update）

```
时间线:
v1(3) → v1(2)+v2(1) → v1(1)+v2(2) → v2(3)
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1          # 最大超出副本数
      maxUnavailable: 0    # 最大不可用副本数
  template:
    spec:
      containers:
      - name: app
        image: myapp:v2
```

```bash
# 执行滚动更新
kubectl set image deployment/myapp app=myapp:v2

# 查看进度
kubectl rollout status deployment/myapp

# 暂停更新
kubectl rollout pause deployment/myapp

# 恢复更新
kubectl rollout resume deployment/myapp

# 回滚
kubectl rollout undo deployment/myapp
kubectl rollout undo deployment/myapp --to-revision=2
```

---

## 2. 金丝雀部署（Canary）

```
流量分配:
Phase 1: v1(100%)
Phase 2: v1(90%) + v2(10%)
Phase 3: v1(50%) + v2(50%)
Phase 4: v2(100%)
```

### 使用 Argo Rollouts

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: myapp
spec:
  replicas: 10
  strategy:
    canary:
      canaryService: myapp-canary
      stableService: myapp-stable
      trafficRouting:
        nginx:
          stableIngress: myapp-ingress
          annotationPrefix: nginx.ingress.kubernetes.io
      steps:
      # 阶段1: 10% 流量，手动确认
      - setWeight: 10
      - pause: {duration: 10m}

      # 阶段2: 25% 流量
      - setWeight: 25
      - pause: {duration: 10m}

      # 阶段3: 50% 流量，检查指标
      - setWeight: 50
      - analysis:
          templates:
          - templateName: success-rate
          args:
          - name: service
            value: myapp-canary

      # 阶段4: 100% 流量
      - setWeight: 100
```

### 使用 Flagger

```yaml
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: myapp
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp
  service:
    port: 80
  analysis:
    interval: 30s
    threshold: 5
    maxWeight: 50
    stepWeight: 10
    metrics:
    - name: request-success-rate
      thresholdRange:
        min: 99
      interval: 1m
    - name: request-duration
      thresholdRange:
        max: 500
      interval: 1m
    webhooks:
    - name: load-test
      url: http://flagger-loadtester.test/
      timeout: 5s
      metadata:
        cmd: "hey -z 1m -q 10 -c 2 http://myapp-canary/"
```

---

## 3. 蓝绿部署（Blue-Green）

```
时间点切换:
蓝环境(100%) ←→ 切换 → 绿环境(100%)
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  selector:
    version: blue  # 切换到 green
  ports:
  - port: 80
    targetPort: 8080
---
# Blue 版本
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-blue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
      version: blue
  template:
    metadata:
      labels:
        app: myapp
        version: blue
    spec:
      containers:
      - name: app
        image: myapp:v1
---
# Green 版本
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-green
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
      version: green
  template:
    metadata:
      labels:
        app: myapp
        version: green
    spec:
      containers:
      - name: app
        image: myapp:v2
```

```bash
# 切换流量到 Green
kubectl patch service myapp -p '{"spec":{"selector":{"version":"green"}}}'

# 回滚到 Blue
kubectl patch service myapp -p '{"spec":{"selector":{"version":"blue"}}}'
```

---

## 4. A/B 测试

```
流量路由:
- Header X-Version: v2 → v2
- Cookie user_group: beta → v2
- 其他 → v1
```

### 使用 Istio

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: myapp
spec:
  hosts:
  - myapp.example.com
  http:
  # A/B 测试路由
  - match:
    - headers:
        x-canary:
          exact: "true"
    route:
    - destination:
        host: myapp
        subset: v2
      weight: 100

  # Cookie 路由
  - match:
    - cookies:
        user:
          exact: "beta"
    route:
    - destination:
        host: myapp
        subset: v2
      weight: 100

  # 默认路由
  - route:
    - destination:
        host: myapp
        subset: v1
      weight: 90
    - destination:
        host: myapp
        subset: v2
      weight: 10
```

### 使用 Gateway API

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: myapp-ab-test
spec:
  parentRefs:
  - name: production-gateway
  rules:
  # Header 路由
  - matches:
    - headers:
      - name: x-canary
        value: "true"
    backendRefs:
    - name: myapp-v2
      port: 80

  # Cookie 路由
  - matches:
    - headers:
      - name: cookie
        value: "user=beta"
    backendRefs:
    - name: myapp-v2
      port: 80

  # 默认路由（90/10 分割）
  - backendRefs:
    - name: myapp-v1
      port: 80
      weight: 90
    - name: myapp-v2
      port: 80
      weight: 10
```

---

## 5. 自动回滚策略

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AnalysisTemplate
metadata:
  name: success-rate
spec:
  metrics:
  - name: success-rate
    interval: 1m
    count: 3
    successCondition: result[0] >= 0.99
    failureLimit: 1
    provider:
      prometheus:
        address: http://prometheus:9090
        query: |
          sum(rate(http_requests_total{service="myapp-canary",status!~"5.."}[1m]))
          /
          sum(rate(http_requests_total{service="myapp-canary"}[1m]))
```

---

## 策略选择建议

| 场景 | 推荐策略 | 工具 |
|------|----------|------|
| 日常发布 | 滚动更新 | 原生 Deployment |
| 关键业务 | 金丝雀 | Argo Rollouts / Flagger |
| 金融支付 | 蓝绿 | Service Selector |
| 产品实验 | A/B 测试 | Istio / Gateway API |
| 快速回滚 | 蓝绿 | Service Selector |

---

## 监控检查清单

- [ ] 错误率（HTTP 5xx）< 1%
- [ ] P99 延迟 < 500ms
- [ ] CPU 使用率 < 80%
- [ ] 内存使用率 < 80%
- [ ] 自定义业务指标正常
