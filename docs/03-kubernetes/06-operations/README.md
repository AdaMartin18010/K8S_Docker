# Kubernetes 运维指南

> 生产环境运维与故障排查

---

## 运维检查清单

### 每日检查

- [ ] Pod 状态检查
- [ ] 节点资源使用
- [ ] 告警处理
- [ ] 日志审查

### 每周检查

- [ ] 资源使用趋势
- [ ] 证书过期时间
- [ ] 安全事件审查
- [ ] 备份验证

### 每月检查

- [ ] 容量规划
- [ ] 成本分析
- [ ] 性能基准测试
- [ ] 灾难恢复演练

---

## 常用运维命令

```bash
# 集群状态
kubectl get nodes
kubectl get pods -A
kubectl top nodes
kubectl top pods -A

# 事件查看
kubectl get events --sort-by='.lastTimestamp'
kubectl get events --field-selector type=Warning

# 日志查看
kubectl logs -f <pod-name>
kubectl logs --previous <pod-name>
stern <pod-name> # 多 Pod 日志
```

---

## 本章内容

- [性能调优](./performance-tuning.md)
- [故障排查](./troubleshooting.md)

---

## 运维工具

| 工具 | 用途 |
|------|------|
| kubectl | 集群管理 |
| stern | 多 Pod 日志 |
| k9s | 终端 UI |
| Lens | GUI 管理 |

---

## 相关文档

- [可观测性](../../05-tools/observability/)
