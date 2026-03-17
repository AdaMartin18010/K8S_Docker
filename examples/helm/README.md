# Helm Charts 示例

> Kubernetes 应用打包与部署

---

## Chart 结构

```
mychart/
├── Chart.yaml          # Chart 元数据
├── values.yaml         # 默认值
├── charts/             # 依赖
└── templates/          # 模板文件
    ├── _helpers.tpl    # 辅助函数
    ├── deployment.yaml
    ├── service.yaml
    └── ingress.yaml
```

---

## 常用命令

```bash
# 创建 Chart
helm create mychart

# 安装
helm install myapp ./mychart

# 升级
helm upgrade myapp ./mychart

# 打包
helm package mychart

# 添加仓库
helm repo add stable https://charts.helm.sh/stable
```

---

## 示例 Charts

- [Web 应用](../helm-charts/web-app/)

---

## 相关文档

- [Helm 速查表](../../docs/99-appendix/helm-cheatsheet.md)
