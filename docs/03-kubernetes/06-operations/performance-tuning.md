# K8s 性能调优

> 提升集群性能与资源利用

---

## 资源优化

### Pod 资源设置

```yaml
resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "1000m"
    memory: "512Mi"
```

**黄金法则**:

- Request = 正常负载下的平均使用
- Limit = 峰值负载 + 20% 缓冲

### VPA 自动调整

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
    updateMode: "Auto"
```

---

## 调度优化

### Pod 亲和性

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

---

## 网络优化

1. **启用 IPVS 模式**: `kube-proxy --proxy-mode=ipvs`
2. **CoreDNS 缓存**: 配置 CoreDNS 缓存插件
3. **本地流量**: 使用 `externalTrafficPolicy: Local`

---

## 存储优化

1. **使用 SSD**: 需要 IOPS 时使用 SSD
2. **本地存储**: 临时数据使用 emptyDir
3. **存储类**: 选择合适的 provisioner

---

## etcd 优化

```yaml
# etcd 调优参数
--quota-backend-bytes=8589934592  # 8GB
defrag 定期整理
snapshot 定期备份
```
