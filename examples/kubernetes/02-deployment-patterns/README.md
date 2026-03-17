# Kubernetes 部署模式

本目录包含常用的部署模式示例。

## 部署模式对比

| 模式 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| **滚动更新** | 简单、资源友好 | 回滚慢、同时运行两个版本 | 大多数场景 |
| **金丝雀** | 风险低、可监控 | 复杂、需要流量控制 | 关键业务 |
| **蓝绿** | 即时回滚、零停机 | 资源翻倍、数据同步复杂 | 金融支付 |
| **A/B 测试** | 精准控制、数据驱动 | 复杂、需要基础设施支持 | 产品验证 |

## 文件说明

| 文件 | 描述 |
|------|------|
| `rolling-update-advanced.yaml` | 高级滚动更新配置 |
| `canary-deployment.yaml` | 金丝雀发布配置 |
| `blue-green-deployment.yaml` | 蓝绿部署配置 |

## 使用指南

### 滚动更新

```bash
# 应用配置
kubectl apply -f rolling-update-advanced.yaml

# 更新镜像（触发滚动更新）
kubectl set image deployment/web-app-rolling web-app=myregistry/web-app:v2.0.0

# 监控更新进度
kubectl rollout status deployment/web-app-rolling

# 查看更新历史
kubectl rollout history deployment/web-app-rolling

# 回滚到上一个版本
kubectl rollout undo deployment/web-app-rolling

# 回滚到特定版本
kubectl rollout undo deployment/web-app-rolling --to-revision=2
```

### 金丝雀发布

```bash
# 部署稳定版本和金丝雀版本
kubectl apply -f canary-deployment.yaml

# 初始状态：90% v1, 10% v2

# 逐步增加金丝雀流量（修改 Ingress annotation）
kubectl annotate ingress web-app-canary nginx.ingress.kubernetes.io/canary-weight="30"

# 监控金丝雀指标
# - 错误率
# - 延迟
# - 资源使用

# 如果正常，继续增加
kubectl annotate ingress web-app-canary nginx.ingress.kubernetes.io/canary-weight="50"

# 最后全部切换
kubectl annotate ingress web-app-canary nginx.ingress.kubernetes.io/canary-weight="100"
# 或删除金丝雀 Ingress，修改稳定版 Deployment 镜像

# 如果有问题，立即回滚
kubectl annotate ingress web-app-canary nginx.ingress.kubernetes.io/canary-weight="0"
```

### 蓝绿部署

```bash
# 当前蓝色环境运行
kubectl apply -f blue-green-deployment.yaml

# 部署新版本（绿色）
kubectl scale deployment web-app-green --replicas=3

# 等待绿色就绪
kubectl rollout status deployment/web-app-green

# 测试绿色环境
curl http://web-app-preview.production.svc.cluster.local

# 切换流量（修改 Service selector）
kubectl patch service web-app-active -p '{"spec":{"selector":{"color":"green"}}}'

# 观察一段时间，如果正常，缩容蓝色
kubectl scale deployment web-app-blue --replicas=0

# 如果有问题，立即回滚
kubectl patch service web-app-active -p '{"spec":{"selector":{"color":"blue"}}}'
```

## 最佳实践

1. **始终设置 PDB** - 防止更新期间全部不可用
2. **配置健康检查** - 确保新 Pod 就绪才继续
3. **监控指标** - 错误率、延迟、资源使用
4. **自动化回滚** - 基于指标自动判断回滚
5. **渐进式发布** - 从小流量开始，逐步增加

## 工具推荐

- **Argo Rollouts** - 高级部署策略（金丝雀、蓝绿、实验）
- **Flagger** - 自动金丝雀分析
- **Istio** - 服务网格流量管理
- **Prometheus** - 指标监控
