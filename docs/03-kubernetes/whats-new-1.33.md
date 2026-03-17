# Kubernetes 1.33 新特性 (2025)

> 代号: Octarine - 最新版本深度解析

---

## 核心更新概览

| 特性 | 阶段 | 说明 |
|------|------|------|
| **Sidecar 容器 GA** | Stable | 原生 Sidecar 支持 |
| **用户命名空间** | Stable | 默认启用 UID/GID 隔离 |
| **DRA 结构化参数** | Beta | 动态资源分配新实现 |
| **VolumeGroupSnapshot** | Beta | 多卷一致性快照 |
| **Pod 级别资源限制** | Alpha | 容器组资源配额 |
| **InPlacePodVerticalScaling** | Beta | 原地垂直扩容 |

---

## 1. Sidecar 容器 GA

K8s 1.33 中 Sidecar 容器正式 GA，成为生产级特性。

```yaml
apiVersion: v1
kind: Pod
spec:
  initContainers:
    - name: istio-proxy
      image: istio/proxyv2:1.22
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

## 2. 用户命名空间 (User Namespaces)

containerd 2.0+ 配合 K8s 1.33 默认启用用户命名空间隔离。

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
```

---

## 3. DRA 结构化参数 (Beta)

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

## 4. VolumeGroupSnapshot (Beta)

多卷一致性快照，适用于数据库等多卷应用。

```yaml
apiVersion: groupsnapshot.storage.k8s.io/v1beta1
kind: VolumeGroupSnapshot
metadata:
  name: db-backup
spec:
  volumeGroupSnapshotClassName: csi-groupsnapshot
  selector:
    matchLabels:
      app: postgres
      tier: database
```

---

## 5. 原地 Pod 垂直扩容 (Beta)

Pod 运行中调整资源限制，无需重建。

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
# 运行时调整
kubectl patch pod mypod --patch '{
  "spec":{
    "containers":[{
      "name":"app",
      "resources":{"requests":{"cpu":"4"}}
    }]
  }
}'
```

---

## 6. 其他重要更新

### Windows 节点优雅关闭 (Stable)

Windows Pod 现在支持优雅关闭，会执行 pre-stop hooks。

### PV 回收策略修复 (Stable)

引入 finalizers 确保 PV 回收策略总是被正确执行。

### 自定义 kubectl debug 配置 (Stable)

```bash
kubectl debug node/mynode --profile=netadmin --image=nicolaka/netshoot
```

---

## 升级建议

```
前置条件:
  - containerd >= 2.0
  - runc >= 1.2
  - etcd >= 3.5.16

升级顺序:
  1. 升级控制平面到 1.33
  2. 升级工作节点到 1.33
  3. 升级 containerd 到 2.0
  4. 启用新特性门控

注意:
  - 用户命名空间需要容器运行时支持
  - DRA 新实现与旧版本不兼容
  - 测试环境验证后再生产升级
```
