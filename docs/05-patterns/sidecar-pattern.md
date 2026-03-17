# Sidecar 模式

> 为主应用容器提供辅助功能的容器模式

---

## 模式概述

Sidecar 模式将辅助功能从主应用容器中分离，作为一个独立的容器运行：

```
┌─────────────────────────────────────────┐
│                  Pod                    │
│  ┌─────────────────────────────────┐   │
│  │        Main Container           │   │
│  │       (Application)             │   │
│  │                                 │   │
│  │  - Business Logic               │   │
│  │  - API Endpoints                │   │
│  └─────────────────────────────────┘   │
│                  │                      │
│  ┌───────────────┼───────────────────┐  │
│  │               │                   │  │
│  │  ┌────────────▼────────────┐      │  │
│  │  │    Sidecar Container    │      │  │
│  │  │                         │      │  │
│  │  │  - Logging Agent        │      │  │
│  │  │  - Monitoring Agent     │      │  │
│  │  │  - Config Watcher       │      │  │
│  │  └─────────────────────────┘      │  │
│  │            Shared Volume          │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

---

## 常见 Sidecar 场景

| 场景 | Sidecar 功能 | 代表工具 |
|------|-------------|----------|
| **日志收集** | 收集应用日志并转发 | Fluent Bit, Logstash |
| **监控代理** | 暴露应用指标 | Prometheus Exporter |
| **配置热更新** | 监听配置变化并通知应用 | Consul Template |
| **服务网格** | 处理服务间通信 | Istio Envoy |
| **数据同步** | 同步数据到远程存储 | Rclone |

---

## K8s 1.29+ Sidecar 容器

原生支持 Sidecar 生命周期管理：

```yaml
apiVersion: v1
kind: Pod
spec:
  initContainers:
    - name: init-myservice
      image: busybox:1.36
      restartPolicy: Always  # Sidecar 特性
  
  containers:
    - name: myapp
      image: myapp:v1.0.0
    
    # Sidecar: 日志收集
    - name: fluent-bit
      image: fluent/fluent-bit:latest
      restartPolicy: Always  # 与主容器独立重启
      volumeMounts:
        - name: logs
          mountPath: /var/log/myapp
```

---

## 经典示例

### 日志收集 Sidecar

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: web-app
spec:
  containers:
    # 主应用
    - name: web
      image: nginx:alpine
      volumeMounts:
        - name: logs
          mountPath: /var/log/nginx
    
    # Sidecar: 日志收集
    - name: fluent-bit
      image: fluent/fluent-bit:latest
      volumeMounts:
        - name: logs
          mountPath: /var/log/nginx
        - name: fluent-bit-config
          mountPath: /fluent-bit/etc/
  
  volumes:
    - name: logs
      emptyDir: {}
    - name: fluent-bit-config
      configMap:
        name: fluent-bit-config
```

---

## Sidecar vs Init Container

| 特性 | Init Container | Sidecar Container |
|------|---------------|-------------------|
| 生命周期 | Pod 启动前运行，完成后退出 | 与主容器并行运行 |
| 重启策略 | 失败后重启 | Always (K8s 1.29+) |
| 用途 | 初始化任务 | 辅助功能 |
| 数量 | 可有多个，顺序执行 | 可有多个，并行运行 |
