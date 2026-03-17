# Pod 与容器

> Kubernetes 最小部署单元详解

---

## Pod 基础

Pod 是 Kubernetes 中最小的部署单元，可以包含一个或多个容器。

### Pod 生命周期

```
Pending → Running → Succeeded/Failed
```

### Pod 配置示例

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: app
    image: myapp:v1
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
```

---

## 容器类型

| 类型 | 说明 | 使用场景 |
|------|------|----------|
| 主容器 | 运行业务应用 | 核心业务逻辑 |
| Sidecar | 辅助容器 | 日志收集、监控 |
| Init | 初始化容器 | 数据准备、配置 |

---

## 本章内容

- [原生 Sidecar](./sidecar-native.md)

---

## K8s 1.33+ Sidecar 新特性

- **确定性启动顺序**: Sidecar 在业务容器之前启动
- **自动终止**: Job 完成后 Sidecar 自动退出
- **健康检查**: Sidecar 支持独立健康检查

---

## 相关文档

- [K8s 工作负载](../02-workloads/)
- [设计模式](../../05-patterns/)
