# GitOps 高级实践 - ArgoCD ApplicationSet

## 概述

ApplicationSet 是 ArgoCD 的多集群/多租户管理利器，通过生成器模式从单一模板自动创建多个 Application，支持 List、Cluster、Git、Matrix 等生成器类型。

> **2025 生产指标**: 使用 ApplicationSet 可将多集群部署时间从 30+ 分钟缩短至 5 分钟（减少 83%），消除配置漂移。

## 架构原理

```
┌─────────────────────────────────────────────────────────────────────┐
│                        ArgoCD Server                               │
│  ┌──────────────────────────────────────────────────────────┐      │
│  │                   ApplicationSet Controller               │      │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐          │      │
│  │  │  Generator │  │  Generator │  │  Generator │          │      │
│  │  │   (List)   │  │  (Cluster) │  │   (Git)    │          │      │
│  │  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘          │      │
│  │        └─────────────────┼─────────────────┘               │      │
│  │                          ▼                                 │      │
│  │  ┌──────────────────────────────────────────────────────┐  │      │
│  │  │                Parameter Sets                         │  │      │
│  │  │  [{cluster: east, url: https://...},                 │  │      │
│  │  │   {cluster: west, url: https://...}]                 │  │      │
│  │  └─────────────────────────┬────────────────────────────┘  │      │
│  │                            ▼                               │      │
│  │  ┌──────────────────────────────────────────────────────┐  │      │
│  │  │              Template Rendering                       │  │      │
│  │  │  Go Template + Sprig 函数                            │  │      │
│  │  └─────────────────────────┬────────────────────────────┘  │      │
│  │                            ▼                               │      │
│  │  ┌──────────────────────────────────────────────────────┐  │      │
│  │  │              Generated Applications                   │  │      │
│  │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐              │  │      │
│  │  │  │ myapp-east│ │myapp-west│ │myapp-prod│              │  │      │
│  │  │  └──────────┘ └──────────┘ └──────────┘              │  │      │
│  │  └──────────────────────────────────────────────────────┘  │      │
│  └────────────────────────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────────────────────┘
```

## 生成器类型详解

### 1. List Generator - 显式列表

适用于：已知固定集群列表，需自定义参数

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: my-app
  namespace: argocd
spec:
  generators:
  - list:
      elements:
      - cluster: prod-east
        url: https://prod-east.example.com
        region: eastus
        replicas: "5"
        environment: production
      - cluster: prod-west
        url: https://prod-west.example.com
        region: westus2
        replicas: "3"
        environment: production
      - cluster: staging
        url: https://staging.example.com
        region: eastus
        replicas: "1"
        environment: staging
  template:
    metadata:
      name: 'my-app-{{cluster}}'
    spec:
      project: default
      source:
        repoURL: https://github.com/org/gitops-repo.git
        targetRevision: main
        path: charts/my-app
        helm:
          valueFiles:
          - values.yaml
          parameters:
          - name: replicaCount
            value: '{{replicas}}'
          - name: region
            value: '{{region}}'
          - name: environment
            value: '{{environment}}'
      destination:
        server: '{{url}}'
        namespace: my-app
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
        - CreateNamespace=true
```

### 2. Cluster Generator - 自动发现

适用于：动态集群环境，自动为新注册集群创建应用

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: monitoring-stack
  namespace: argocd
spec:
  generators:
  - clusters:
      # 只选择带有指定标签的集群
      selector:
        matchLabels:
          environment: production
          managed-by: argocd
  template:
    metadata:
      name: 'monitoring-{{name}}'
      labels:
        component: monitoring
    spec:
      project: default
      source:
        repoURL: https://github.com/org/gitops-repo.git
        targetRevision: main
        path: base/monitoring
        kustomize:
          commonLabels:
            cluster: '{{name}}'
      destination:
        server: '{{server}}'
        namespace: monitoring
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
        - CreateNamespace=true
```

### 3. Git Generator - 基于目录/文件

适用于： monorepo 结构，每个目录对应一个应用

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: cluster-configs
  namespace: argocd
spec:
  generators:
  - git:
      repoURL: https://github.com/org/gitops-repo.git
      revision: main
      directories:
      # 为 clusters/ 下的每个子目录生成应用
      - path: clusters/*
        exclude: false
  template:
    metadata:
      name: 'config-{{path.basename}}'
    spec:
      project: default
      source:
        repoURL: https://github.com/org/gitops-repo.git
        targetRevision: main
        path: '{{path}}'
      destination:
        server: 'https://{{path.basename}}-api.example.com'
        namespace: default
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
```

文件生成器（不同环境不同配置）：

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: env-based-apps
  namespace: argocd
spec:
  generators:
  - git:
      repoURL: https://github.com/org/gitops-repo.git
      revision: main
      files:
      - path: "config/**/config.json"
  template:
    metadata:
      name: '{{app.name}}-{{env.name}}'
    spec:
      project: default
      source:
        repoURL: https://github.com/org/gitops-repo.git
        targetRevision: main
        path: 'apps/{{app.path}}'
        helm:
          valueFiles:
          - 'values-{{env.name}}.yaml'
      destination:
        server: '{{env.cluster}}'
        namespace: '{{app.namespace}}'
```

### 4. Matrix Generator - 矩阵组合

适用于：多维度部署（集群 × 服务 × 环境）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: all-services-all-clusters
  namespace: argocd
spec:
  generators:
  - matrix:
      generators:
      # 第一维度：集群
      - clusters:
          selector:
            matchLabels:
              environment: production
      # 第二维度：服务（从 Git 目录）
      - git:
          repoURL: https://github.com/org/gitops-repo.git
          revision: main
          directories:
          - path: services/*
  template:
    metadata:
      # 组合命名：服务名-集群名
      name: '{{path.basename}}-{{name}}'
      labels:
        service: '{{path.basename}}'
        cluster: '{{name}}'
    spec:
      project: default
      source:
        repoURL: https://github.com/org/gitops-repo.git
        targetRevision: main
        path: '{{path}}'
        helm:
          valueFiles:
          - values.yaml
          # 集群特定配置覆盖
          - '../../clusters/{{name}}/values.yaml'
      destination:
        server: '{{server}}'
        namespace: '{{path.basename}}'
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
        - CreateNamespace=true
```

### 5. SCM Provider Generator - PR 预览环境

适用于：为每个 Pull Request 创建预览环境

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: pr-preview
  namespace: argocd
spec:
  generators:
  - scmProvider:
      github:
        organization: myorg
        repositories:
        - myapp
      # 只处理带特定标签的 PR
      filters:
      - labelMatch: preview
  template:
    metadata:
      name: 'myapp-pr-{{number}}'
    spec:
      project: default
      source:
        repoURL: https://github.com/myorg/myapp.git
        targetRevision: '{{head_sha}}'
        path: helm/myapp
        helm:
          parameters:
          - name: image.tag
            value: 'pr-{{number}}'
          - name: ingress.host
            value: 'pr-{{number}}.preview.example.com'
      destination:
        server: https://kubernetes.default.svc
        namespace: 'pr-{{number}}'
      syncPolicy:
        automated:
          prune: true
        syncOptions:
        - CreateNamespace=true
        - PruneLast=true
```

## 高级模板功能

### Go Template 语法

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: advanced-templating
  namespace: argocd
spec:
  goTemplate: true
  goTemplateOptions: ["missingkey=error"]
  generators:
  - list:
      elements:
      - cluster: prod
        replicas: 5
        features: "feature-a,feature-b"
  template:
    metadata:
      # 字符串操作
      name: '{{.cluster | lower}}-app'
      annotations:
        # 条件判断
        'enabled-features': '{{if .features}}{{.features}}{{else}}none{{end}}'
    spec:
      source:
        helm:
          parameters:
          # 数学运算
          - name: maxReplicas
            value: '{{mul (int .replicas) 2}}'
          # 条件渲染
          - name: resources.enabled
            value: '{{if gt (int .replicas) 3}}true{{else}}false{{end}}'
```

### 进度发布策略 (Progressive Rollout)

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: progressive-deployment
  namespace: argocd
spec:
  generators:
  - list:
      elements:
      - cluster: canary
        wave: "1"
      - cluster: staging
        wave: "2"
      - cluster: prod-1
        wave: "3"
      - cluster: prod-2
        wave: "3"
      - cluster: prod-3
        wave: "4"
  strategy:
    type: RollingSync
    rollingSync:
      steps:
      - matchExpressions:
        - key: wave
          operator: In
          values: ["1"]
      - matchExpressions:
        - key: wave
          operator: In
          values: ["2"]
      - matchExpressions:
        - key: wave
          operator: In
          values: ["3"]
      - matchExpressions:
        - key: wave
          operator: In
          values: ["4"]
  template:
    metadata:
      name: 'app-{{cluster}}'
      labels:
        wave: '{{wave}}'
    spec:
      # ... 常规配置
      syncPolicy:
        automated: {}
```

## 多租户安全模式

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: tenant-apps
  namespace: argocd
spec:
  generators:
  - clusters:
      selector:
        matchLabels:
          tenant: "*"
  template:
    metadata:
      name: 'tenant-{{metadata.labels.tenant}}-app'
    spec:
      project: 'tenant-{{metadata.labels.tenant}}'
      source:
        repoURL: https://github.com/org/tenant-manifests.git
        targetRevision: main
        path: 'tenants/{{metadata.labels.tenant}}'
      destination:
        server: '{{server}}'
        namespace: 'tenant-{{metadata.labels.tenant}}'
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
      # 资源限制
      ignoreDifferences:
      - group: apps
        kind: Deployment
        jsonPointers:
        - /spec/replicas
```

## 同步策略详解

```yaml
syncPolicy:
  automated:
    prune: true           # 删除 Git 中移除的资源
    selfHeal: true        # 回滚手动修改
    allowEmpty: false     # 防止空源删除所有资源
  syncOptions:
    - CreateNamespace=true
    - PrunePropagationPolicy=foreground  # 前台级联删除
    - PruneLast=true                     # 最后执行删除
    - RespectIgnoreDifferences=true      # 尊重 ignoreDifferences
  retry:
    limit: 5
    backoff:
      duration: 5s
      factor: 2
      maxDuration: 3m
```

## 与 KubeVirt 集成

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: vm-workloads
  namespace: argocd
spec:
  generators:
  - clusters:
      selector:
        matchLabels:
          virtualization: enabled
  template:
    metadata:
      name: 'vms-{{name}}'
    spec:
      project: default
      source:
        repoURL: https://github.com/org/vm-manifests.git
        targetRevision: main
        path: vms/
      destination:
        server: '{{server}}'
        namespace: vms
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        # VM 可能需要更长的同步超时
        retry:
          limit: 10
          backoff:
            duration: 30s
            factor: 2
            maxDuration: 10m
```

## 监控与可观测性

```yaml
# 查看 ApplicationSet 状态
kubectl get applicationset -n argocd
kubectl describe applicationset my-appset -n argocd

# 查看生成的 Applications
argocd app list | grep my-appset

# JSON 格式检查
argocd app list -o json | jq '.[] | {
  name: .metadata.name,
  sync: .status.sync.status,
  health: .status.health.status
}'
```

## 总结

| 场景 | 推荐生成器 |
|------|-----------|
| 固定集群列表 | List |
| 动态集群发现 | Cluster |
| Monorepo 多应用 | Git Directory |
| 多维度组合 | Matrix |
| PR 预览环境 | SCM Provider |
| 渐进式发布 | RollingSync |

ApplicationSet 将多集群管理从人工脚本转变为声明式 GitOps 工作流，是 2025 年 Kubernetes 多集群管理的最佳实践。
