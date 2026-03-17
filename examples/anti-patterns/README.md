# 反例集合 (Anti-Patterns)

本目录收集 Docker 和 Kubernetes 的常见错误配置，帮助开发者识别和避免这些问题。

## 目录结构

```
anti-patterns/
├── docker/
│   ├── Dockerfile.security-risks    # 安全风险示例
│   └── Dockerfile.performance-issues # 性能问题示例
└── kubernetes/
    ├── pod.security-risks.yaml      # Pod 安全反例
    └── deployment.anti-patterns.yaml # Deployment 反例
```

## Docker 反例

### 安全风险

| 风险 | 说明 | 解决方案 |
|------|------|----------|
| 使用 root 用户 | 容器以 root 运行，被攻击后获得主机 root 权限 | 使用 `USER` 指令指定非 root 用户 |
| 硬编码密码 | 敏感信息写在 Dockerfile 中 | 使用 Build Args 或运行时注入 |
| 使用 latest 标签 | 构建不可重现，可能引入意外变更 | 使用具体版本标签 |
| 特权容器 | 拥有主机全部权限 | 避免使用 `privileged` |
| 挂载 Docker Socket | 可以控制主机 Docker | 使用远程 Docker API 或避免挂载 |

### 性能问题

| 问题 | 说明 | 解决方案 |
|------|------|----------|
| 单层镜像 | 包含所有构建依赖 | 使用多阶段构建 |
| 没有 .dockerignore | 上下文过大，构建慢 | 添加 .dockerignore 文件 |
| 缓存失效 | 不合理的指令顺序 | 将不常变更的指令放在前面 |

## Kubernetes 反例

### 安全风险

| 风险 | 说明 | 解决方案 |
|------|------|----------|
| privileged: true | 特权容器 | 使用安全上下文限制 |
| hostNetwork/hostPID | 共享主机命名空间 | 保持隔离 |
| 没有资源限制 | 资源耗尽攻击 | 设置 resources.limits |
| 默认 Service Account | 权限过大 | 创建专用 SA 并限制权限 |

### 稳定性问题

| 问题 | 说明 | 解决方案 |
|------|------|----------|
| 没有健康检查 | 故障 Pod 不会被处理 | 添加 liveness/readiness 探针 |
| 没有 Pod 中断预算 | 升级时可能导致全部不可用 | 添加 PDB |
| 硬编码配置 | 难以管理 | 使用 ConfigMap/Secret |

## 如何使用

1. **学习**：查看反例，理解为什么是错误的
2. **识别**：在自己的项目中检查是否存在类似问题
3. **修复**：参考对应的最佳实践示例进行修复

## 相关资源

- `../docker/basic/Dockerfile.good` - Dockerfile 最佳实践
- `../kubernetes/01-basic-resources/` - K8s 资源最佳实践
- `../kubernetes/06-security/` - 安全配置示例
