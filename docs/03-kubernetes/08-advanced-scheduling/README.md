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

## 调度优化

### 亲和性与反亲和性

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

### 污点与容忍

```yaml
tolerations:
- key: "dedicated"
  operator: "Equal"
  value: "gpu"
  effect: "NoSchedule"
```

---

## 相关文档

- [K8s 架构](../01-architecture/)
- [多集群](../../08-multicluster/)
