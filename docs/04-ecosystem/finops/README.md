# FinOps - 云成本管理

> Kubernetes 云成本优化与管理

---

## 什么是 FinOps？

FinOps (Financial Operations) 是连接工程、财务和业务的实践，通过数据驱动决策优化云支出，实现速度与成本的平衡。

```
┌─────────────────────────────────────────────────────────────┐
│              FinOps 核心原则                                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. 团队协作                                                  │
│     工程、财务、业务团队共同参与成本决策                       │
│                                                              │
│  2. 可见性                                                    │
│     实时了解云支出，分配到团队/项目/功能                       │
│                                                              │
│  3. 优化                                                      │
│     持续优化资源使用，消除浪费                                │
│                                                              │
│  FinOps 生命周期:                                            │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐              │
│  │  Inform  │ →  │ Optimize │ →  │ Operate  │              │
│  │  (知情)  │    │  (优化)  │    │  (运营)  │              │
│  └──────────┘    └──────────┘    └──────────┘              │
│       ↑                              │                      │
│       └──────────────────────────────┘                      │
│              持续迭代                                         │
└─────────────────────────────────────────────────────────────┘
```

---

## K8s 成本驱动因素

| 成本类别 | 占比 | 优化策略 |
|----------|------|----------|
| **计算 (CPU/内存)** | 60-70% | 资源优化、Spot实例、自动扩缩容 |
| **存储** | 15-20% | 存储类选择、数据生命周期管理 |
| **网络** | 10-15% | 流量优化、区域选择 |
| **负载均衡** | 5-10% | 共享LB、Gateway API |

---

## OpenCost - 开源成本监控

CNCF 项目，Kubernetes 成本监控的行业标准。

### 安装

```bash
# 使用 Helm 安装 OpenCost
helm install opencost opencost/opencost \
  --namespace opencost \
  --create-namespace

# 访问 UI
kubectl port-forward -n opencost svc/opencost 9090:9090
```

### Prometheus 集成

```yaml
# OpenCost 导出 Prometheus 指标
# 查询示例: 命名空间成本
sum(opencost_container_memory_working_set_bytes) by (namespace)
*
scalar(avg(opencost_node_ram_cost))
```

---

## Kubecost - 企业级成本管理

基于 OpenCost 的商业产品，提供更丰富的功能。

### 核心功能

| 功能 | 说明 |
|------|------|
| **实时成本分配** | 按命名空间、Deployment、Pod、标签分解 |
| **成本优化建议** | 资源优化、闲置资源识别 |
| **预算告警** | 阈值告警、异常检测 |
| **Showback/Chargeback** | 成本分摊和计费 |
| **多集群聚合** | 统一视图管理多个集群 |

### 成本优化建议

```
┌─────────────────────────────────────────────────────────────┐
│              Kubecost 成本优化建议                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  🔴 高优先级                                                  │
│  • 缩减 overprovisioned 的 Pod (节省 $500/月)               │
│  • 删除未使用的 PVC (节省 $200/月)                           │
│  • 使用 Spot 实例 (节省 60% 计算成本)                        │
│                                                              │
│  🟡 中优先级                                                  │
│  • 调整 Request/Limit 比例                                   │
│  • 优化存储类选择                                            │
│  • 实施 HPA 自动扩缩容                                       │
│                                                              │
│  🟢 低优先级                                                  │
│  • 购买 Reserved Instances                                   │
│  • 实施 Cluster Autoscaler                                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 成本优化策略

### 1. 资源优化

```yaml
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
```

### 2. Spot/Preemptible 实例

```yaml
apiVersion: apps/v1
kind: Deployment
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
  minReplicas: 1
  maxReplicas: 100
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 50
          periodSeconds: 60
```

---

## 成本分摊模型

### 标签策略

```yaml
metadata:
  labels:
    app: frontend
    team: platform
    environment: production
    cost-center: cc-12345
    project: ecommerce
    owner: platform-team
```

### Prometheus 成本查询

```promql
# 按团队统计成本
sum by (team) (
  (
    rate(container_cpu_usage_seconds_total[1h]) * 0.05 +
    container_memory_working_set_bytes / 1024 / 1024 / 1024 * 0.01
  )
)

# 按环境统计成本
sum by (environment) (
  opencost_container_memory_working_set_bytes * 0.01 +
  opencost_container_cpu_usage * 0.05
)
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
- [ ] 购买 Reserved Instances/Savings Plans
- [ ] 优化存储类和生命周期

---

## 2025 趋势

- **AI 驱动的成本优化**: ML 预测资源需求，自动调整
- **实时成本可见性**: 秒级成本数据更新
- **Shift-Left FinOps**: 在 CI/CD 中集成成本估算
- **碳成本关联**: 成本与碳排放一起优化
