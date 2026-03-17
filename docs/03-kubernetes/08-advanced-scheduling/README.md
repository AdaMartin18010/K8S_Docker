# 高级调度

> Kubernetes 高级调度与资源管理

---

## 本章内容

- [DRA - 动态资源分配](./dra.md)
- [边缘计算](./edge-computing.md)
- [多集群管理](./multi-cluster.md)
- [节点自动扩缩容](./node-autoscaling.md)
- [虚拟集群](./vcluster.md)

---

## 调度优化策略

### 1. 亲和性与反亲和性

```yaml
affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchLabels:
            app: myapp
        topologyKey: kubernetes.io/hostname
```

### 2. 污点与容忍

```yaml
tolerations:
- key: "dedicated"
  operator: "Equal"
  value: "gpu"
  effect: "NoSchedule"
```

### 3. 拓扑分布约束

```yaml
topologySpreadConstraints:
- maxSkew: 1
  topologyKey: topology.kubernetes.io/zone
  whenUnsatisfiable: ScheduleAnyway
  labelSelector:
    matchLabels:
      app: myapp
```

---

## 调度器扩展

| 扩展点 | 用途 |
|--------|------|
| QueueSort | 自定义排序 |
| PreFilter | 前置过滤 |
| Filter | 节点筛选 |
| PostFilter | 后备处理 |
| Score | 节点评分 |

---

## 相关文档

- [K8s 架构](../01-architecture/)
- [多集群管理](../../08-multicluster/)
