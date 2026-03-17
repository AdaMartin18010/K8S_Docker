# 生产实践指南

> 从理论到生产落地的实战经验

---

## 本章内容

1. [CI/CD 实践](./cicd-guide.md)
2. [性能调优](./performance-benchmarks.md)
3. [故障排查](./debugging-guide.md)
4. [混沌工程](./chaos-engineering/README.md)
5. [案例研究](./case-studies/)

---

## 生产检查清单

### 部署前

- [ ] Dockerfile 使用多阶段构建
- [ ] 镜像使用非 root 用户
- [ ] 资源限制已配置 (requests/limits)
- [ ] 健康检查已配置 (liveness/readiness)
- [ ] 镜像已扫描无高危漏洞

### 部署时

- [ ] 使用 GitOps 工作流
- [ ] 配置 PDB (PodDisruptionBudget)
- [ ] HPA 自动扩缩容已配置
- [ ] 网络策略已启用
- [ ] 监控告警已配置

### 运行中

- [ ] 定期更新基础镜像
- [ ] 日志收集正常
- [ ] 备份策略已验证
- [ ] 灾难恢复计划已测试

---

## CI/CD 最佳实践

```yaml
# .github/workflows/deploy.yml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build
        run: docker build -t myapp:${{ github.sha }} .

      - name: Scan
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: myapp:${{ github.sha }}

      - name: Push
        run: |
          docker push myapp:${{ github.sha }}

      - name: Deploy
        run: |
          kubectl set image deployment/myapp app=myapp:${{ github.sha }}
          kubectl rollout status deployment/myapp
```

---

## 关联代码

- [CI/CD 指南](./cicd-guide.md)
- [Helm 图表](../05-tools/)
