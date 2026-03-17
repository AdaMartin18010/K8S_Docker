# Kubernetes 完整指南

> **版本**: Kubernetes 1.30-1.32 | 最后更新: 2025年3月

---

## 目录

1. [核心架构](#1-核心架构)
2. [基础资源](#2-基础资源)
3. [部署策略](#3-部署策略)
4. [存储管理](#4-存储管理)
5. [安全](#5-安全)
6. [可观测性](#6-可观测性)
7. [2025年重要更新](#7-2025年重要更新)
8. [关联代码示例](#8-关联代码示例)

---

## 1. 核心架构

### 1.1 控制平面组件

```
┌─────────────────────────────────────────────────────────────┐
│                     控制平面 (Control Plane)                  │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ API Server  │  │   etcd      │  │ Controller Manager  │ │
│  │  (kube-     │  │  (存储)      │  │   (控制器)          │ │
│  │ apiserver)  │  │             │  │                     │ │
│  └──────┬──────┘  └─────────────┘  └─────────────────────┘ │
│         │                                                   │
│  ┌──────┴────────────────────────────────────────────────┐ │
│  │              Scheduler (调度器)                        │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────┼───────────────────────────────┐
│                     数据平面 (Data Plane)                   │
├─────────────────────────────┼───────────────────────────────┤
│  ┌──────────────────────────┴───────────────────────────┐  │
│  │                    Worker Node                        │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌───────────────┐ │  │
│  │  │   Kubelet   │  │  Kube-proxy │  │    Runtime    │ │  │
│  │  │ (节点代理)  │  │  (网络代理) │  │(containerd/   │ │  │
│  │  │             │  │             │  │ CRI-O)        │ │  │
│  │  └─────────────┘  └─────────────┘  └───────────────┘ │  │
│  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 1.30-1.32 重要新特性

| 版本 | 特性 | 说明 |
|------|------|------|
| 1.30 | Sidecar 容器 GA | 原生支持 Sidecar 生命周期管理 |
| 1.31 | CephFS/RBD 移除 | 必须使用 CSI 驱动 |
| 1.32 | 内存管理器 GA | 为 Guaranteed QoS Pod 分配独占内存 |
| 1.32 | QueueingHint | 调度器吞吐量优化 |

---

## 2. 基础资源

### 2.1 Pod 最佳实践 (2025版)

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: web-app
spec:
  securityContext:
    runAsNonRoot: true        # 强制非 root
    runAsUser: 1000
    fsGroup: 1000
    seccompProfile:
      type: RuntimeDefault    # 默认 seccomp

  containers:
    - name: app
      image: myapp:v1.0.0
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities:
          drop: [ALL]         # 删除所有能力

      resources:
        requests:
          cpu: 100m
          memory: 128Mi
        limits:
          cpu: 500m
          memory: 512Mi

      # 健康检查
      startupProbe:
        httpGet:
          path: /health/startup
          port: 8080
        failureThreshold: 30

      livenessProbe:
        httpGet:
          path: /health/live
          port: 8080
        periodSeconds: 10

      readinessProbe:
        httpGet:
          path: /health/ready
          port: 8080
        periodSeconds: 5

      volumeMounts:
        - name: tmp
          mountPath: /tmp

  volumes:
    - name: tmp
      emptyDir:
        sizeLimit: 100Mi
```

### 2.2 Deployment + HPA + PDB

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 0
  selector:
    matchLabels:
      app: web-app
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
        - name: app
          image: myapp:v1.0.0
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
---
# 自动扩缩容
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: web-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: web-app
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
---
# 中断预算
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: web-app-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: web-app
```

### 2.3 新版 Sidecar 容器 (1.29+)

```yaml
apiVersion: v1
kind: Pod
spec:
  initContainers:
    - name: init-myservice
      image: busybox:1.36
      restartPolicy: Always  # Sidecar 特性
      command: ['sh', '-c', 'echo init']

  containers:
    - name: myapp
      image: myapp:v1

    # Sidecar 容器
    - name: nginx-sidecar
      image: nginx:alpine
      restartPolicy: Always  # 与主容器独立重启
```

---

## 3. 部署策略

### 3.1 金丝雀发布

```yaml
# 稳定版本 (90% 流量)
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
  template:
    metadata:
      labels:
        app: web-app
        track: stable
    spec:
      containers:
        - name: app
          image: myapp:v1.0.0
---
# 金丝雀版本 (10% 流量)
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
  template:
    metadata:
      labels:
        app: web-app
        track: canary
    spec:
      containers:
        - name: app
          image: myapp:v2.0.0
```

---

## 4. 存储管理

### 4.1 StorageClass 配置

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: premium-ssd
provisioner: kubernetes.io/gce-pd
parameters:
  type: pd-ssd
  replication-type: regional
volumeBindingMode: WaitForFirstConsumer  # 延迟绑定
allowVolumeExpansion: true
reclaimPolicy: Retain
```

### 4.2 StatefulSet

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres-headless
  replicas: 3
  podManagementPolicy: OrderedReady
  updateStrategy:
    type: RollingUpdate
    rollingUpdate:
      partition: 0
  selector:
    matchLabels:
      app: postgres
  template:
    spec:
      containers:
        - name: postgres
          image: postgres:16-alpine
          volumeMounts:
            - name: data
              mountPath: /var/lib/postgresql
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        storageClassName: premium-ssd
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 50Gi
```

---

## 5. 安全

### 5.1 Pod Security Standards (替代 PSP)

```yaml
# 命名空间级别应用
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/enforce-version: latest
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

### 5.2 RBAC 最小权限

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: web-app-sa
automountServiceAccountToken: false
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: web-app-role
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get", "list"]
    resourceNames: ["web-app-config"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: web-app-binding
subjects:
  - kind: ServiceAccount
    name: web-app-sa
roleRef:
  kind: Role
  name: web-app-role
  apiGroup: rbac.authorization.k8s.io
```

### 5.3 NetworkPolicy

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
spec:
  podSelector: {}
  policyTypes:
    - Ingress
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-app
spec:
  podSelector:
    matchLabels:
      app: web-app
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
      ports:
        - protocol: TCP
          port: 8080
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: postgres
      ports:
        - protocol: TCP
          port: 5432
```

---

## 6. 可观测性

### 6.1 ServiceMonitor

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: web-app-metrics
spec:
  namespaceSelector:
    matchNames: [production]
  selector:
    matchLabels:
      app: web-app
  endpoints:
    - port: metrics
      path: /metrics
      interval: 15s
```

---

## 7. 2025年重要更新

### 7.1 已废弃/移除功能

| 功能 | 版本 | 状态 | 替代方案 |
|------|------|------|----------|
| PodSecurityPolicy | 1.25+ | ❌ 已移除 | Pod Security Standards |
| Dockershim | 1.24+ | ❌ 已移除 | containerd/CRI-O |
| CephFS/RBD 内嵌 | 1.31+ | ❌ 已移除 | CSI 驱动 |
| `topology.kubernetes.io/zone` | 1.17+ | ✅ 标准 | 替代旧标签 |

### 7.2 Gateway API (替代 Ingress)

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: example-gateway
spec:
  gatewayClassName: nginx
  listeners:
    - name: http
      protocol: HTTP
      port: 80
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: example-route
spec:
  parentRefs:
    - name: example-gateway
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /api
      backendRefs:
        - name: api-service
          port: 80
```

---

## 8. 关联代码示例

| 主题 | 代码示例 |
|------|----------|
| Pod 配置 | `examples/kubernetes/01-basic-resources/pod-good.yaml` |
| Deployment | `examples/kubernetes/01-basic-resources/deployment-good.yaml` |
| 部署模式 | `examples/kubernetes/02-deployment-patterns/` |
| 存储 | `examples/kubernetes/05-storage/` |
| 安全 | `examples/kubernetes/06-security/` |
| 可观测性 | `examples/kubernetes/07-observability/` |

---

## 参考

- [Kubernetes 1.32 Release Notes](https://kubernetes.io/blog/2024/12/11/kubernetes-1-32-release/)
- [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [Gateway API](https://gateway-api.sigs.k8s.io/)
