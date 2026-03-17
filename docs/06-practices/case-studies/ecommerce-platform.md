# 案例：电商平台容器化改造

> 某大型电商平台的 K8s 迁移实践

---

## 背景

- **公司**: 某头部电商平台
- **规模**: 日均订单 1000 万+
- **挑战**: 大促期间流量波动大，传统架构扩容慢

---

## 架构演进

### 改造前 (单体架构)

```
┌─────────────────────────────────────┐
│           单体应用                   │
│  ┌─────────┐ ┌─────────┐ ┌────────┐ │
│  │   Web   │ │   API   │ │  Admin │ │
│  │   MVC   │ │   RPC   │ │   MVC  │ │
│  └────┬────┘ └────┬────┘ └───┬────┘ │
│       └────────────┴─────────┘      │
│              ┌──────────┐           │
│              │  MySQL   │           │
│              │ (Master) │           │
│              └──────────┘           │
└─────────────────────────────────────┘
```

**问题**:
- 代码耦合严重，发布风险高
- 数据库单点，扩展困难
- 大促扩容需要数小时

### 改造后 (微服务 + K8s)

```
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway                             │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼───────┐    ┌───────▼───────┐    ┌───────▼───────┐
│  User Service │    │ Product Svc   │    │  Order Svc    │
│   (10 pods)   │    │   (15 pods)   │    │   (20 pods)   │
└───────┬───────┘    └───────┬───────┘    └───────┬───────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
                    ┌─────────▼─────────┐
                    │   Message Bus     │
                    │   (Kafka Cluster) │
                    └───────────────────┘
```

---

## 关键改造点

### 1. 服务拆分

| 服务 | 职责 | Pod 数 (日常) | Pod 数 (大促) |
|------|------|--------------|--------------|
| user-service | 用户、认证 | 10 | 30 |
| product-service | 商品、库存 | 15 | 50 |
| order-service | 订单、购物车 | 20 | 80 |
| payment-service | 支付 | 10 | 40 |
| notification | 通知 | 5 | 20 |

### 2. 自动扩缩容

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: order-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: order-service
  minReplicas: 20
  maxReplicas: 200
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          averageUtilization: 70
          type: Utilization
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - type: Pods
          value: 20
          periodSeconds: 15
```

### 3. 缓存策略

```yaml
apiVersion: v1
kind: Deployment
metadata:
  name: product-service
spec:
  template:
    spec:
      containers:
        - name: product
          image: product-service:v2.1.0
          env:
            - name: CACHE_TYPE
              value: "redis_cluster"
            - name: CACHE_TTL
              value: "300"
```

---

## 效果

| 指标 | 改造前 | 改造后 | 提升 |
|------|--------|--------|------|
| 扩容时间 | 2 小时 | 30 秒 | 99.6% |
| 发布频率 | 1 周/次 | 1 天/次 | 7x |
| 系统可用性 | 99.9% | 99.99% | +0.09% |
| 资源利用率 | 15% | 45% | 3x |

---

## 经验教训

1. **数据库拆分要慎重**: 初期保留部分共享库，逐步迁移
2. **监控先行**: 未完善的监控导致几次故障
3. **限流熔断**: 大促期间未配置限流导致雪崩
4. **渐进式迁移**: 按模块逐步迁移，而非一次性切换
