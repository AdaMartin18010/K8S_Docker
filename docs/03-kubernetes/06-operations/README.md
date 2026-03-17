# Kubernetes 运维指南

> 生产环境运维与故障排查

---

## 本章内容

1. [运维概述](./README.md)
2. [性能调优](./performance-tuning.md)
3. [故障排查](./troubleshooting.md)

---

## 运维检查清单

### 每日检查

- [ ] 检查 Pod 状态
- [ ] 查看节点资源使用
- [ ] 检查告警

### 每周检查

- [ ] 审查资源使用趋势
- [ ] 检查证书过期时间
- [ ] 审查安全事件

### 每月检查

- [ ] 容量规划审查
- [ ] 成本分析
- [ ] 备份验证

---

## 常用运维命令

```bash
# 查看集群状态
kubectl get nodes
kubectl get pods -A

# 查看资源使用
kubectl top nodes
kubectl top pods -A

# 查看事件
kubectl get events --sort-by='.lastTimestamp'

# 查看日志
kubectl logs -f <pod-name>
```

---

## 关联代码

- [examples/kubernetes/07-observability/](../../examples/kubernetes/07-observability/)
