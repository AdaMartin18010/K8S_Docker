# Helm Charts 示例

本目录包含生产级 Helm Chart 示例。

## 目录结构

```
helm-charts/
└── web-app/
    ├── Chart.yaml          # Chart 元数据
    ├── values.yaml         # 默认配置值
    ├── values-dev.yaml     # 开发环境配置
    ├── values-prod.yaml    # 生产环境配置
    ├── README.md           # Chart 文档
    └── templates/          # K8s 模板文件
        ├── _helpers.tpl    # 辅助模板
        ├── deployment.yaml
        ├── service.yaml
        ├── ingress.yaml
        ├── hpa.yaml
        ├── pdb.yaml
        ├── configmap.yaml
        ├── secret.yaml
        └── serviceaccount.yaml
```

## 快速开始

```bash
# 进入 Chart 目录
cd helm-charts/web-app

# 依赖更新
helm dependency update

# 模板渲染（测试）
helm template web-app . --values values.yaml

# 安装
helm install my-app . -n production --create-namespace

# 升级
helm upgrade my-app . -n production

# 回滚
helm rollback my-app 1 -n production

# 卸载
helm uninstall my-app -n production
```

## 多环境部署

```bash
# 开发环境
helm install web-app-dev . \
  -n development \
  --values values.yaml \
  --values values-dev.yaml

# 生产环境
helm install web-app-prod . \
  -n production \
  --values values.yaml \
  --values values-prod.yaml
```

## Chart 开发最佳实践

1. **版本管理**:
   - `version`: Chart 版本（语义化版本）
   - `appVersion`: 应用版本

2. **配置分离**:
   - `values.yaml`: 默认配置
   - `values-<env>.yaml`: 环境特定配置

3. **模板函数**:
   - 使用 `_helpers.tpl` 定义可复用模板
   - 使用 `include` 函数引用

4. **安全检查**:
   - 使用 `required` 函数确保必需值
   - 使用 `default` 提供默认值

5. **文档**:
   - README.md 描述配置参数
   - values.yaml 中添加注释

## 常用命令

```bash
# 打包 Chart
helm package web-app/

# 验证 Chart
helm lint web-app/

# 查看所有版本
helm history my-app -n production

# 查看 Release 状态
helm status my-app -n production

# 获取 values
helm get values my-app -n production
```
