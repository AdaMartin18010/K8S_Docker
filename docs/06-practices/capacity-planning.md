# 容量规划与成本优化

> Kubernetes 资源规划与成本控制

---

## 容量规划框架

```
┌─────────────────────────────────────────────────────────────┐
│              容量规划生命周期                                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   收集数据         分析趋势          预测需求         执行扩容 │
│  ┌──────┐       ┌──────┐         ┌──────┐        ┌──────┐   │
│  │Metrics│─────▶│Trend │────────▶│Forecast│──────▶│Action │   │
│  │ Logs │       │Analysis        │        │        │       │   │
│  └──────┘       └──────┘         └──────┘        └──────┘   │
│      ▲                                              │        │
│      └──────────────────────────────────────────────┘        │
│                    反馈优化                                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 资源使用基线

### 收集资源使用数据

```bash
# 获取历史资源使用 (需 Metrics Server)
kubectl top pods -A --containers

# 导出 CSV 用于分析
kubectl top pods -A --no-headers | \
  awk '{print $1","$2","$3","$4}' > resource_usage.csv

# 使用 Prometheus 查询
# CPU 使用率百分位
histogram_quantile(0.95,
  sum(rate(container_cpu_usage_seconds_total[5m])) by (pod, le)
)

# 内存使用率
container_memory_working_set_bytes / container_spec_memory_limit_bytes
```

### 基线分析

| 工作负载类型 | CPU Request | CPU Limit | 内存 Request | 内存 Limit |
|-------------|-------------|-----------|--------------|------------|
| **Web 前端** | 100m | 500m | 128Mi | 512Mi |
| **API 服务** | 200m | 1000m | 256Mi | 1Gi |
| **批处理任务** | 500m | 2000m | 512Mi | 2Gi |
| **缓存服务** | 500m | 1000m | 1Gi | 2Gi |
| **数据库** | 1000m | 4000m | 2Gi | 8Gi |

---

## 垂直扩容 (VPA)

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: myapp-vpa
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp
  updatePolicy:
    updateMode: "Auto"  # Off, Initial, Recreate, Auto
    minAllowed:
      cpu: 50m
      memory: 100Mi
    maxAllowed:
      cpu: 1000m
      memory: 2Gi
  resourcePolicy:
    containerPolicies:
      - containerName: '*'
        minAllowed:
          cpu: 50m
          memory: 100Mi
        maxAllowed:
          cpu: 1000m
          memory: 2Gi
        controlledResources: ["cpu", "memory"]
        controlledValues: RequestsAndLimits
```

---

## 成本优化策略

### 1. 资源优化

```bash
# 识别过度配置
kubectl resource-capacity --pods --util --sort cpu.util

# 推荐设置
kubectl apply -f - <<EOF
apiVersion: v1
kind: LimitRange
metadata:
  name: resource-limits
spec:
  limits:
    - default:
        cpu: 500m
        memory: 512Mi
      defaultRequest:
        cpu: 100m
        memory: 128Mi
      type: Container
EOF
```

### 2. Spot/Preemptible 实例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: spot-workload
spec:
  template:
    spec:
      nodeSelector:
        node.kubernetes.io/capacity: spot
      tolerations:
        - key: spot
          operator: Equal
          value: "true"
          effect: NoSchedule
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              podAffinityTerm:
                labelSelector:
                  matchLabels:
                    app: spot-workload
                topologyKey: kubernetes.io/hostname
```

### 3. 自动扩缩容

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: cost-optimized-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp
  minReplicas: 1  # 非工作时间最小化
  maxReplicas: 100
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Pods
      pods:
        metric:
          name: http_requests_per_second
        target:
          type: AverageValue
          averageValue: "100"
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300  # 5分钟稳定期
      policies:
        - type: Percent
          value: 50
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - type: Percent
          value: 100
          periodSeconds: 15
```

---

## 成本监控

### 标签策略

```yaml
metadata:
  labels:
    app: frontend
    team: platform
    environment: production
    cost-center: cc-12345
    project: ecommerce
```

### 成本分摊查询 (Prometheus)

```promql
# 按团队统计 CPU 使用
sum by (team) (
  rate(container_cpu_usage_seconds_total[1h])
)

# 按环境统计内存使用
sum by (environment) (
  container_memory_working_set_bytes
)

# 估算成本 (假设 $0.05/CPU-hour, $0.01/GB-hour)
(
  sum(rate(container_cpu_usage_seconds_total[1h])) * 0.05
) + (
  sum(container_memory_working_set_bytes) / 1024 / 1024 / 1024 * 0.01
)
```

---

## 容量预测

### 增长模型

```python
# 使用 Prometheus API 获取历史数据
# 进行时间序列预测

import requests
from sklearn.linear_model import LinearRegression
import numpy as np

# 获取 CPU 使用数据
response = requests.get(
    'http://prometheus:9090/api/v1/query_range',
    params={
        'query': 'rate(container_cpu_usage_seconds_total[5m])',
        'start': '2025-01-01T00:00:00Z',
        'end': '2025-03-17T00:00:00Z',
        'step': '1h'
    }
)

# 线性回归预测
data = response.json()['data']['result']
X = np.array(range(len(data))).reshape(-1, 1)
y = np.array([float(v[1]) for v in data[0]['values']])

model = LinearRegression()
model.fit(X, y)

# 预测 30 天后
future = np.array([[len(data) + 30 * 24]])  # 30 days * 24 hours
prediction = model.predict(future)
print(f"预测 30 天后 CPU 需求: {prediction[0]:.2f} cores")
```

---

## 成本优化清单

- [ ] 启用 HPA 自动扩缩容
- [ ] 配置 VPA 优化资源请求
- [ ] 使用 Spot/Preemptible 实例
- [ ] 设置资源配额和限制
- [ ] 定期清理未使用资源
- [ ] 启用集群自动扩缩容
- [ ] 实施标签策略便于成本分摊
- [ ] 监控和告警异常成本
