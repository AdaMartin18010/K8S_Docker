# Init Container 模式

> Pod 启动前的初始化工作

---

## 什么是 Init Container？

Init Container 是在主容器启动之前运行的容器，用于执行初始化任务。

**特点**:

- 按顺序执行
- 必须成功完成
- 不支持 `restartPolicy: Always`
- 资源限制独立于主容器

```
Pod Lifecycle:
┌─────────┐   ┌─────────┐   ┌─────────┐   ┌──────────┐
│ Pending │ → │ Init 0  │ → │ Init 1  │ → │ Running  │
└─────────┘   │ Success │   │ Success │   │ (Main)   │
              └─────────┘   └─────────┘   └──────────┘
```

---

## 常见用途

| 用途 | 说明 |
|------|------|
| **配置生成** | 从外部系统拉取配置 |
| **数据库迁移** | 执行 schema 更新 |
| **权限设置** | 修改文件系统权限 |
| **等待依赖** | 等待数据库/服务就绪 |
| **数据初始化** | 从远程下载数据 |

---

## 示例：等待数据库

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-with-db
spec:
  initContainers:
    - name: wait-for-db
      image: busybox:1.36
      command:
        - sh
        - -c
        - |
          until nc -z postgres 5432; do
            echo "Waiting for database..."
            sleep 2
          done
  containers:
    - name: myapp
      image: myapp:v1
```

---

## 示例：数据库迁移

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-with-migration
spec:
  initContainers:
    - name: db-migrate
      image: myapp:v1
      command: ["python", "manage.py", "migrate"]
      env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
  containers:
    - name: myapp
      image: myapp:v1
```

---

## Init vs Sidecar 对比

| 特性 | Init Container | Sidecar |
|------|---------------|---------|
| **执行时机** | 主容器之前 | 与主容器并行 |
| **生命周期** | 完成后退出 | 与 Pod 同生命周期 |
| **重启策略** | 不支持 Always | 支持 Always (v1.29+) |
| **用途** | 初始化任务 | 辅助服务 |

---

## 注意事项

1. **超时控制**: 长时间运行的 Init 会阻塞 Pod 启动
2. **失败处理**: Init 失败会导致 Pod 进入 CrashLoopBackOff
3. **资源计算**: Init Container 的资源限制取最大值
4. **存储共享**: 使用 emptyDir 或 PVC 与主容器共享数据
