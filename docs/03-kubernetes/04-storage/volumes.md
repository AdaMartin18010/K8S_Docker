# Kubernetes 存储卷

> Pod 数据持久化方案

---

## 卷类型

| 类型 | 用途 | 生命周期 |
|------|------|----------|
| **emptyDir** | 临时存储 | 随 Pod |
| **hostPath** | 节点文件 | 随节点 |
| **ConfigMap** | 配置注入 | 独立于 Pod |
| **Secret** | 敏感数据 | 独立于 Pod |
| **PVC** | 持久存储 | 独立于 Pod |

---

## emptyDir

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cache-pod
spec:
  containers:
    - name: app
      image: myapp
      volumeMounts:
        - name: cache
          mountPath: /cache
  volumes:
    - name: cache
      emptyDir:
        sizeLimit: 1Gi
```

**特点**:

- Pod 创建时创建
- Pod 删除时删除
- 可用于容器间共享

---

## hostPath

```yaml
volumes:
  - name: logs
    hostPath:
      path: /var/log/myapp
      type: DirectoryOrCreate
```

**警告**: 生产环境慎用，节点相关

---

## ConfigMap/Secret 挂载

```yaml
volumes:
  - name: config
    configMap:
      name: app-config
      items:
        - key: app.conf
          path: app.conf
  - name: secrets
    secret:
      secretName: app-secret
```

---

## 最佳实践

1. **敏感数据**: 使用 Secret + 内存挂载
2. **配置文件**: 使用 ConfigMap
3. **持久数据**: 使用 PVC
4. **缓存**: 使用 emptyDir
