# Kubernetes 工作负载

> Pod、Deployment、StatefulSet、DaemonSet、Job 详解

---

## 本章内容

1. [Pod 详解](./pods.md)
2. [Deployment](./deployments.md)
3. [StatefulSet](./statefulsets.md)
4. [DaemonSet](./daemonsets.md)
5. [Job 和 CronJob](./jobs.md)
6. [HPA 自动扩缩容](./hpa.md)

---

## 工作负载对比

| 资源类型 | 适用场景 | 特点 | 示例 |
|----------|----------|------|------|
| **Pod** | 单个容器/紧密耦合容器 | 最小部署单元 | 调试、Sidecar |
| **Deployment** | 无状态应用 | 滚动更新、回滚 | Web 服务、API |
| **StatefulSet** | 有状态应用 | 稳定网络标识、有序部署 | 数据库、消息队列 |
| **DaemonSet** | 每个节点一个 Pod | 节点级服务 | 日志收集、监控 |
| **Job** | 一次性任务 | 完成即退出 | 数据迁移、批处理 |
| **CronJob** | 定时任务 | 基于时间调度 | 备份、清理 |

---

## Pod 完整配置示例 (2025 标准)

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: web-app
  namespace: production
  labels:
    app: web-app
    version: v1
spec:
  # 安全上下文
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 1000
    seccompProfile:
      type: RuntimeDefault
  
  # 初始化容器
  initContainers:
    - name: init-db
      image: busybox:1.36
      command: ['sh', '-c', 'until nc -z db 5432; do sleep 2; done']
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities:
          drop: [ALL]
  
  # 主容器
  containers:
    - name: app
      image: myapp:v1.0.0
      ports:
        - name: http
          containerPort: 8080
      
      # 环境变量
      env:
        - name: ENV
          value: production
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
      
      # 资源限制 (必需)
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
          port: http
        failureThreshold: 30
      
      livenessProbe:
        httpGet:
          path: /health/live
          port: http
        periodSeconds: 10
      
      readinessProbe:
        httpGet:
          path: /health/ready
          port: http
        periodSeconds: 5
      
      # 安全上下文
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities:
          drop: [ALL]
      
      # 卷挂载
      volumeMounts:
        - name: tmp
          mountPath: /tmp
        - name: config
          mountPath: /app/config
          readOnly: true
  
  # 卷定义
  volumes:
    - name: tmp
      emptyDir:
        sizeLimit: 100Mi
    - name: config
      configMap:
        name: app-config
  
  # 亲和性
  affinity:
    podAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
        - weight: 100
          podAffinityTerm:
            labelSelector:
              matchExpressions:
                - key: app
                  operator: In
                  values: [web-app]
            topologyKey: kubernetes.io/hostname
  
  # 容忍度
  tolerations:
    - key: dedicated
      operator: Equal
      value: web
      effect: NoSchedule
  
  # 终止优雅期
  terminationGracePeriodSeconds: 60
```

---

## Deployment 示例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
  labels:
    app: web-app
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
          ports:
            - containerPort: 8080
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
```

---

## 关联代码

- [examples/kubernetes/01-basic-resources/pod-good.yaml](../../examples/kubernetes/01-basic-resources/pod-good.yaml)
- [examples/kubernetes/01-basic-resources/deployment-good.yaml](../../examples/kubernetes/01-basic-resources/deployment-good.yaml)
