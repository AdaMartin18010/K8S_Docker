# Sidecar 模式

> 分离主容器与辅助功能的容器模式

---

## 模式定义

Sidecar 模式将辅助功能（日志收集、监控、配置同步等）放在独立的容器中，与主容器共享生命周期和资源。

```
┌─────────────────────────────────────────────┐
│                    Pod                       │
│  ┌───────────────────────────────────────┐  │
│  │            Main Container             │  │
│  │           (Application)               │  │
│  │            Port: 8080                 │  │
│  └───────────────┬───────────────────────┘  │
│                  │ localhost               │
│  ┌───────────────▼───────────────────────┐  │
│  │           Sidecar Container           │  │
│  │         (Log Shipper)                 │  │
│  │         (Config Watcher)              │  │
│  │         (Service Mesh Proxy)          │  │
│  └───────────────────────────────────────┘  │
│                                             │
│  Shared: Volume, Network Namespace          │
└─────────────────────────────────────────────┘
```

---

## 常见用途

| 用途 | Sidecar 容器 | 主容器关注点 |
|------|-------------|-------------|
| **日志收集** | Fluent Bit, Filebeat | 业务逻辑 |
| **配置同步** | Reloader, ConfigMap Watcher | 业务应用 |
| **服务网格** | Istio Proxy, Linkerd | 微服务 |
| **监控** | Prometheus Exporter | 应用指标 |

---

## K8s 1.29+ 原生 Sidecar

```yaml
apiVersion: v1
kind: Pod
spec:
  initContainers:
    - name: istio-proxy
      image: istio/proxyv2:1.20
      restartPolicy: Always  # 关键！使其成为原生 Sidecar
      securityContext:
        allowPrivilegeEscalation: false
  containers:
    - name: myapp
      image: myapp:v1
```

**优势**:
- Sidecar 先于主容器启动
- Sidecar 后于主容器终止
- 独立的重启策略

---

## 示例：日志收集 Sidecar

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-with-logging
spec:
  containers:
    # 主应用
    - name: myapp
      image: myapp:v1
      volumeMounts:
        - name: logs
          mountPath: /var/log/myapp
    
    # Sidecar: 日志收集
    - name: fluent-bit
      image: fluent/fluent-bit:latest
      volumeMounts:
        - name: logs
          mountPath: /var/log/myapp
          readOnly: true
        - name: fluent-config
          mountPath: /fluent-bit/etc/
  
  volumes:
    - name: logs
      emptyDir: {}
    - name: fluent-config
      configMap:
        name: fluent-config
```

---

## 注意事项

1. **资源限制**: 为 Sidecar 设置合理的资源限制
2. **生命周期**: 确保 Sidecar 正确处理信号
3. **日志管理**: Sidecar 自身的日志也需要管理
4. **调试难度**: 多容器 Pod 调试更复杂
