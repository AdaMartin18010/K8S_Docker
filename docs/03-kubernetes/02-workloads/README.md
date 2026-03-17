# Kubernetes 工作负载

> Pod、Deployment、StatefulSet 等工作负载资源详解

---

## 本章内容

1. [Pod 深入理解](./pods-deep-dive.md)
2. [Deployment 部署](./deployments.md)
3. [StatefulSet 有状态应用](./statefulsets.md)
4. [DaemonSet 守护进程](./daemonsets.md)
5. [Job & CronJob 批处理](./jobs.md)

---

## 工作负载对比

| 资源类型 | 适用场景 | 特点 |
|----------|----------|------|
| **Pod** | 单实例应用 | 最基础的部署单元 |
| **Deployment** | 无状态应用 | 滚动更新、回滚 |
| **StatefulSet** | 有状态应用 | 稳定网络标识、持久存储 |
| **DaemonSet** | 节点级服务 | 每个节点运行一个 |
| **Job** | 一次性任务 | 完成后退出 |
| **CronJob** | 定时任务 | 基于时间调度 |

---

## Pod 配置模板 (2025 标准)

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: web-app
  labels:
    app: web-app
spec:
  # 安全上下文
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 1000
    seccompProfile:
      type: RuntimeDefault

  containers:
    - name: app
      image: myapp:v1.0.0

      # 资源限制
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

      # 安全设置
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities:
          drop: [ALL]

      volumeMounts:
        - name: tmp
          mountPath: /tmp

  volumes:
    - name: tmp
      emptyDir:
        sizeLimit: 100Mi
```

---

## Deployment 高级配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 3

  # 滚动更新策略
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 0

  # 进度控制
  progressDeadlineSeconds: 600
  minReadySeconds: 30
  revisionHistoryLimit: 10

  selector:
    matchLabels:
      app: web-app
  template:
    spec:
      terminationGracePeriodSeconds: 60
      containers:
        - name: app
          image: myapp:v1.0.0
          lifecycle:
            preStop:
              exec:
                command: ["/bin/sh", "-c", "sleep 15"]
```

---

## 关联代码

- [examples/kubernetes/01-basic-resources/](../../examples/kubernetes/01-basic-resources/)
- [examples/kubernetes/02-deployment-patterns/](../../examples/kubernetes/02-deployment-patterns/)
