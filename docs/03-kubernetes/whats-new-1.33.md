# Kubernetes 1.33 新特性 (2025)

> 代号: Octarine - 最新版本深度解析
>
> **发布日期**: 2025年4月23日 | **64项增强** (18 GA, 20 Beta, 26 Alpha)

---

## 核心更新概览

| 特性 | 阶段 | 说明 |
|------|------|------|
| **Sidecar 容器 GA** | Stable | 原生 Sidecar 支持，`restartPolicy: Always` |
| **用户命名空间** | Beta (默认启用) | UID/GID 隔离，容器逃逸防护 |
| **InPlacePodVerticalScaling** | Beta | 原地垂直扩容，无需重建 Pod |
| **DRA 结构化参数** | Beta | 动态资源分配新实现 |
| **Gateway API** | GA | 正式可用，v1.2+ 新特性 |
| **ServiceCIDR** | GA | 多 Service IP 范围支持 |
| **ClusterTrustBundle** | Beta | 集群级根证书管理 |
| **Ordered Namespace Deletion** | Alpha | 有序命名空间删除 |

---

## 1. Sidecar 容器 GA ⭐

K8s 1.33 中 Sidecar 容器正式 GA，成为生产级特性。

```yaml
apiVersion: v1
kind: Pod
spec:
  initContainers:
    - name: istio-proxy
      image: istio/proxyv2:1.24
      restartPolicy: Always  # 关键: 使其成为 Sidecar
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop: ["ALL"]
      resources:
        requests:
          cpu: 100m
          memory: 128Mi
  containers:
    - name: myapp
      image: myapp:v1
```

### Sidecar 生命周期

```
Pod 启动:
  1. 初始化容器按顺序启动
  2. Sidecar (restartPolicy: Always) 先于主容器启动
  3. 主容器启动

Pod 终止:
  1. 主容器收到 SIGTERM
  2. 主容器终止
  3. Sidecar 继续运行直到就绪探针失败
  4. Sidecar 终止
```

---

## 2. 用户命名空间 (User Namespaces) - Beta 默认启用

K8s 1.33 中用户命名空间进入 Beta 并**默认启用**，配合 containerd 2.0+ 提供更强的安全隔离。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: secure-pod
spec:
  hostUsers: false  # 启用用户命名空间
  containers:
    - name: app
      image: nginx
      securityContext:
        runAsUser: 0  # 容器内是 root
        runAsGroup: 0
        # 但主机上映射为非特权用户!
```

### 映射机制

```
容器内 UID          主机 UID
   0 (root)    →    100000
   1           →    100001
   65535       →    165535

优势:
- 容器逃逸后无 root 权限
- 可安全运行特权容器
- 符合多租户安全要求
```

---

## 3. 原地 Pod 垂直扩容 (InPlacePodVerticalScaling) - Beta ⭐

Pod 运行中调整资源限制，**无需重建**，状态ful应用零停机扩容。

```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    resize.kubernetes.io/resources: "InPlacePodVerticalScaling"
spec:
  containers:
    - name: app
      image: myapp
      resources:
        requests:
          cpu: "1"
          memory: 2Gi
        limits:
          cpu: "2"
          memory: 4Gi
```

```bash
# 运行时调整资源
kubectl patch pod mypod --patch '{
  "spec":{
    "containers":[{
      "name":"app",
      "resources":{"requests":{"cpu":"4"}}
    }]
  }
}'

# 查看资源调整状态
kubectl get pod mypod -o yaml | grep -A 10 resizeStatus
```

### 使用场景

- **数据库扩容**: MySQL/PostgreSQL 运行时增加内存
- **突发流量**: 快速增加 CPU 应对峰值
- **成本优化**: 空闲时降低资源请求

---

## 4. DRA 结构化参数 (Beta)

动态资源分配的核心重构，kube-scheduler 可直接模拟资源分配。

```yaml
apiVersion: resource.k8s.io/v1beta1
kind: ResourceClaim
metadata:
  name: gpu-claim
spec:
  resourceClassName: nvidia-gpu
  structuredParameters:
    apiVersion: gpu.resource.nvidia.com/v1alpha1
    kind: GpuParameters
    requests:
      count: 2
      selector:
        productName: "NVIDIA A100"
        memory: "40Gi"
    sharing:
      strategy: TimeSlicing
```

### 优势

- **调度器可预测**: 无需驱动即可验证资源可用性
- **支持自动扩缩容**: Cluster Autoscaler 可理解资源需求
- **更灵活的共享策略**: 时分复用、MIG 等

---

## 5. Service IP 扩展 (GA)

支持多个 Service CIDR，解决大型集群 IP 耗尽问题。

```yaml
apiVersion: networking.k8s.io/v1
kind: ServiceCIDR
metadata:
  name: additional-cidr
spec:
  cidr: "10.96.0.0/12"  # 额外的 Service IP 范围
```

---

## 6. ClusterTrustBundle (Beta)

集群级根证书管理，替代 per-namespace 的 kube-root-ca.crt。

```yaml
apiVersion: trust.certificates.k8s.io/v1beta1
kind: ClusterTrustBundle
metadata:
  name: internal-ca
spec:
  signerName: internal-ca.example.com
  trustBundle: |
    -----BEGIN CERTIFICATE-----
    MIICpDCCAYwCCQDU+pQ4nEHXqzANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDDAls
    ...
    -----END CERTIFICATE-----
```

---

## 7. 有序命名空间删除 (Alpha)

确保命名空间删除时按正确顺序清理资源，避免 NetworkPolicy 先于 Pod 删除导致的安全风险。

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: production
  annotations:
    deletion.kubernetes.io/order: "strict"  # 启用有序删除
```

---

## 8. 其他重要更新

### NetworkPolicy 日志 (Beta)

支持记录允许/拒绝的网络流量，便于调试。

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: logging-policy
  annotations:
    policy.networking.k8s.io/log-level: "info"  # 启用日志
```

### Job SuccessCriteriaMet 状态 (GA)

```bash
# kubectl get job 现在显示 SuccessCriteriaMet 状态
kubectl get job my-job
NAME     COMPLETIONS   DURATION   STATUS
my-job   1/1           5m30s      SuccessCriteriaMet
```

### Endpoints API 废弃

- **v1 Endpoints API 已废弃**，迁移到 **EndpointSlice**
- `status.nodeInfo.kubeProxyVersion` 字段已移除
- **gitRepo Volume 类型已移除**，使用 init-container 替代

---

## 9. 安全增强

| 特性 | 阶段 | 说明 |
|------|------|------|
| **PodSchedulingReadiness** | GA | 控制 Pod 是否参与调度 |
| **NodeLogQuery** | GA | 节点日志查询 |
| **ProcMount** | Beta | /proc 访问控制 |
| **ServiceAccountToken 改进** | GA | 令牌绑定到节点和生命周期 |

---

## 升级建议

```
前置条件:
  - containerd >= 2.0 (用户命名空间支持)
  - runc >= 1.2
  - etcd >= 3.5.21

升级顺序:
  1. 升级控制平面到 1.33
  2. 升级工作节点到 1.33
  3. 升级 containerd 到 2.0
  4. 验证并启用新特性

注意事项:
  - 从 1.32 → 1.33 → 1.34 逐步升级
  - 检查 matchLabelKeys 相关调度问题
  - DRA 新实现与旧版本不兼容
  - 用户命名空间需要容器运行时支持
```

---

## 参考

- [Kubernetes 1.33 Release Notes](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.33.md)
- [KEP-1287: In-place Pod Vertical Scaling](https://github.com/kubernetes/enhancements/issues/1287)
- [KEP-127: User Namespaces](https://github.com/kubernetes/enhancements/issues/127)
