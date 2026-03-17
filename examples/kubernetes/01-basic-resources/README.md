# Kubernetes 基础资源示例

本目录包含 Pod、Deployment、Service 等基础资源的最佳实践配置。

## 文件说明

| 文件 | 说明 |
|------|------|
| `pod-good.yaml` | ✅ Pod 最佳实践 |
| `pod-bad.yaml` | ❌ Pod 常见错误（反例）|
| `deployment-good.yaml` | ✅ Deployment + HPA + PDB |
| `service-good.yaml` | ✅ 各类 Service 配置 |

## 关键概念

### Pod 安全
- `runAsNonRoot: true` - 禁止 root
- `readOnlyRootFilesystem: true` - 只读根文件系统
- `allowPrivilegeEscalation: false` - 禁止提权
- `capabilities: drop: [ALL]` - 丢弃所有能力

### 健康检查
- **startupProbe** - 启动检查，防止慢启动被误判
- **livenessProbe** - 存活检查，失败则重启
- **readinessProbe** - 就绪检查，失败则从 Service 移除

### 资源管理
- **requests** - 调度保证
- **limits** - 运行限制
- **HPA** - 水平自动扩缩
- **PDB** - 中断预算，保证可用性

## 快速使用

```bash
# 应用配置
kubectl apply -f pod-good.yaml
kubectl apply -f deployment-good.yaml

# 查看状态
kubectl get pods -o wide
kubectl get deployments
kubectl get hpa
kubectl get pdb

# 查看详情
kubectl describe pod web-app-good
kubectl describe deployment web-app
```
