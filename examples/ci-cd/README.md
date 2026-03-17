# CI/CD 示例

本目录包含持续集成和持续部署的配置示例。

## GitHub Actions

### 工作流文件

| 文件 | 描述 |
|------|------|
| `docker-build.yml` | Docker 构建、扫描、推送 |
| `k8s-deploy.yml` | 多环境 Kubernetes 部署 |

### 功能特性

#### Docker 构建工作流

- ✅ 多架构构建（AMD64/ARM64）
- ✅ BuildKit 缓存加速
- ✅ Trivy 安全扫描
- ✅ Docker Scout 分析
- ✅ SBOM 和 Provenance 生成
- ✅ 代码覆盖率报告

#### Kubernetes 部署工作流

- ✅ 多环境部署（开发/测试/生产）
- ✅ Helm 部署
- ✅ 金丝雀发布
- ✅ 自动回滚
- ✅ Slack 通知

## 使用指南

### 设置 Secrets

在 GitHub 仓库中设置以下 Secrets：

```
# 容器仓库
GITHUB_TOKEN (自动提供)

# AWS 凭证（如使用 EKS）
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY

# Kubernetes 配置
KUBECONFIG_STAGING (base64 编码)
KUBECONFIG_PRODUCTION (base64 编码)

# 通知
SLACK_WEBHOOK_URL
```

### 目录结构

```
.github/
└── workflows/
    ├── docker-build.yml      # 构建工作流
    ├── k8s-deploy.yml        # 部署工作流
    └── pr-check.yml          # PR 检查
```

## 最佳实践

### 分支策略

```
main        → 生产环境
develop     → 开发环境
feature/*   → PR 检查
```

### 镜像标签策略

- `latest` - 最新开发版本
- `v1.0.0` - 语义化版本
- `sha-xxxx` - Git commit SHA
- `pr-123` - PR 编号

### 安全实践

1. **漏洞扫描**: 每次构建都扫描
2. **最小权限**: 使用专用 Service Account
3. **密钥管理**: 使用 GitHub Secrets
4. **审计日志**: 记录所有部署

## 其他 CI/CD 工具

### GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - build
  - test
  - deploy

build:
  stage: build
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
```

### Jenkins

```groovy
// Jenkinsfile
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'docker build -t myapp .'
            }
        }
    }
}
```

### Argo CD

```yaml
# Application 定义
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: web-app
spec:
  project: default
  source:
    repoURL: https://github.com/example/app
    targetRevision: HEAD
    path: k8s
  destination:
    server: https://kubernetes.default.svc
    namespace: production
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```
