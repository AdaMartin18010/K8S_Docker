# 云原生设计模式

> 云原生架构模式与实践

---

## 本章内容

1. [Sidecar 模式](./sidecar-pattern.md)
2. [微服务模式](./microservices.md)
3. [部署模式](./deployment-strategies.md)
4. [微服务模式](./microservices.md)

---

## 设计模式分类

### 1. 容器设计模式

| 模式 | 说明 | 示例 |
|------|------|------|
| **Sidecar** | 与主容器共享生命周期 | 日志收集、监控代理 |
| **Ambassador** | 代理外部服务访问 | Redis 连接代理 |
| **Adapter** | 标准化接口 | 指标转换器 |
| **Init** | 初始化容器 | 数据库迁移、配置准备 |

### 2. 微服务模式

| 模式 | 说明 | 示例 |
|------|------|------|
| **服务发现** | 动态服务定位 | Kubernetes Service |
| **配置中心** | 外部化配置管理 | ConfigMap/Secret |
| **熔断器** | 故障隔离 | Hystrix, Resilience4j |
| **限流** | 流量控制 | Rate Limiter |

### 3. 部署模式

| 模式 | 说明 | 适用场景 |
|------|------|----------|
| **滚动更新** | 渐进式替换 | 大多数场景 |
| **金丝雀** | 小流量验证 | 关键业务 |
| **蓝绿** | 瞬时切换 | 金融支付 |
| **A/B 测试** | 流量分割 | 产品验证 |

---

## Sidecar 模式详解

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
    # 主应用容器
    - name: myapp
      image: myapp:v1.0.0

    # Sidecar: 日志收集
    - name: fluent-bit
      image: fluent/fluent-bit:latest
      volumeMounts:
        - name: logs
          mountPath: /var/log/myapp

    # Sidecar: 监控代理
    - name: prometheus-exporter
      image: prom/node-exporter:latest
      ports:
        - containerPort: 9100
```

---

## 部署模式对比

```
滚动更新 (Rolling Update):
v1 ──▶ v1+v2 ──▶ v2

金丝雀 (Canary):
v1 (90%) ──▶ v1(70%)+v2(30%) ──▶ v2 (100%)

蓝绿 (Blue-Green):
蓝(100%) ──▶ 蓝(100%)+绿(0%,就绪) ──▶ 切换 ──▶ 绿(100%)
```

---

## 关联代码

- [K8s 示例](../../examples/kubernetes/)
