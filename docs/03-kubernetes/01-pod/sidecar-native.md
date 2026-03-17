# Kubernetes 原生 Sidecar 容器 (1.33 GA)

## 概述

Kubernetes 1.33 (2025年4月) 将 Sidecar 容器支持提升为 GA，彻底解决了传统 Sidecar 的启动顺序和生命周期管理问题。

## 核心改进

| 特性 | 传统 Sidecar | 原生 Sidecar (1.33+) |
|------|-------------|---------------------|
| 启动顺序 | 不确定 | 确定（initContainers 顺序） |
| Job 支持 | 阻塞 | 自动退出 |
| 探针 | 不支持 startup | 支持全部探针 |
| OOM 处理 | 独立 | 与主容器对齐 |

## 基础配置

```yaml
apiVersion: v1
kind: Pod
spec:
  initContainers:
  # 标准 init（先执行）
  - name: init-config
    image: busybox
    command: ['sh', '-c', 'echo "init"']

  # Sidecar 容器
  - name: envoy-sidecar
    image: envoy:v1.30
    restartPolicy: Always  # 关键！标记为 Sidecar
    startupProbe:
      httpGet:
        path: /ready
        port: 15021
      failureThreshold: 30
    resources:
      limits:
        cpu: 500m
        memory: 256Mi

  containers:
  # 主容器（Sidecar 就绪后才启动）
  - name: main-app
    image: myapp:v1
```

## Job 场景

```yaml
apiVersion: batch/v1
kind: Job
spec:
  template:
    spec:
      initContainers:
      - name: dapr-sidecar
        image: daprio/daprd:1.13
        restartPolicy: Always
        startupProbe:
          httpGet:
            path: /v1.0/healthz
            port: 3500
      containers:
      - name: task
        image: task:v1
        command: ['python', 'run.py']  # 完成后退出
        # Sidecar 会自动停止，Job 正常完成
```

## 多 Sidecar 生产配置

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      initContainers:
      - name: istio-proxy
        image: istio/proxyv2:1.21
        restartPolicy: Always
        securityContext:
          runAsUser: 1337
        startupProbe:
          httpGet:
            path: /healthz/ready
            port: 15021

      - name: otel-collector
        image: otel/opentelemetry-collector:0.98
        restartPolicy: Always
        startupProbe:
          httpGet:
            path: /
            port: 13133

      containers:
      - name: app
        image: myapp:v2
        env:
        - name: OTEL_EXPORTER_OTLP_ENDPOINT
          value: http://localhost:4317
```

## 迁移建议

| K8s 版本 | 建议 |
|---------|------|
| < 1.29 | 继续使用传统模式 |
| 1.29-1.32 | 启用 SidecarContainers 特性门控 |
| >= 1.33 | 直接使用，特性 GA |

原生 Sidecar 消除了启动竞态条件，简化了 Job 管理，是 1.33 最重要的特性之一。
