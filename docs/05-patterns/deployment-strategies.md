# 部署策略模式

> 金丝雀、蓝绿、滚动更新详解

---

## 策略对比

| 策略 | 风险 | 资源需求 | 回滚速度 | 适用场景 |
|------|------|----------|----------|----------|
| **滚动更新** | 中 | 低 | 慢 | 常规更新 |
| **金丝雀** | 低 | 中 | 快 | 关键业务 |
| **蓝绿** | 低 | 高 | 瞬时 | 金融支付 |

---

## 1. 滚动更新

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 0
```

---

## 2. 金丝雀发布

```yaml
# 稳定版本 (90%)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app-stable
spec:
  replicas: 9
  selector:
    matchLabels:
      app: web-app
      track: stable
---
# 金丝雀版本 (10%)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app-canary
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web-app
      track: canary
```

---

## 3. 蓝绿部署

```yaml
# 蓝色版本
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app-blue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-app
      color: blue
---
# 绿色版本
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app-green
spec:
  replicas: 0  # 初始为 0
  selector:
    matchLabels:
      app: web-app
      color: green
```

---

## 关联代码

- [examples/kubernetes/02-deployment-patterns/](../../examples/kubernetes/02-deployment-patterns/)
