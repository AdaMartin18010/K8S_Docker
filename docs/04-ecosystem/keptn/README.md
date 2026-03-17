# Keptn - 云原生应用生命周期编排

## 概述

Keptn 是一个云原生应用生命周期管理工具，与 GitOps 解决方案集成，支持部署前评估、健康检查、标准化部署前后任务以及可观测性。

## 核心特性

| 特性 | 描述 |
|------|------|
| SLO 分析 | 自动化服务级别目标分析和验证 |
| 部署编排 | 支持多阶段应用交付 |
| 可观测性 | 基于 OpenTelemetry 的 DORA 指标 |
| 指标聚合 | 统一多数据源指标访问 |
| GitOps 集成 | 与 ArgoCD、Flux 等工具无缝集成 |

## 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                      Keptn 架构                                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   KeptnApp  │    │   Analysis  │    │   Metrics   │         │
│  │   (应用)     │◀──▶│   (分析)     │◀──▶│   (指标)     │         │
│  └──────┬──────┘    └─────────────┘    └─────────────┘         │
│         │                                                       │
│         ▼                                                       │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   Kubernetes Workloads                   │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐     │   │
│  │  │Deployment│  │StatefulSet│  │DaemonSet│  │  Job    │     │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘     │   │
│  │       └─────────────┴─────────────┴─────────────┘        │   │
│  │                         │                                 │   │
│  │              ┌──────────▼──────────┐                     │   │
│  │              │   Scheduling Gate   │                     │   │
│  │              │   (阻塞部署)         │                     │   │
│  │              └──────────┬──────────┘                     │   │
│  └─────────────────────────┼─────────────────────────────────┘   │
│                            │                                     │
│                   ┌────────▼────────┐                           │
│                   │  Pre/Post Tasks │                           │
│                   │  - 依赖检查      │                           │
│                   │  - 镜像扫描      │                           │
│                   │  - 测试执行      │                           │
│                   └─────────────────┘                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 安装部署

### Helm 安装

```bash
# 添加 Helm 仓库
helm repo add keptn https://charts.lifecycle.keptn.sh
helm repo update

# 安装 Keptn
helm upgrade --install keptn keptn/keptn \
  -n keptn-system \
  --create-namespace \
  --wait

# 启用应用自动发现
helm upgrade --install keptn keptn/keptn \
  -n keptn-system \
  --set features.automaticAppDiscovery.enabled=true
```

### 配置 OpenTelemetry

```yaml
# keptnconfig.yaml
apiVersion: options.keptn.sh/v1alpha1
kind: KeptnConfig
metadata:
  name: keptnconfig
  namespace: keptn-system
spec:
  OTelCollectorUrl: 'jaeger-collector.monitoring.svc.cluster.local:4317'
  keptnAppCreationRequestTimeoutSeconds: 30
  observabilityTimeout: 5m
```

```bash
kubectl apply -f keptnconfig.yaml
```

## 核心概念

### KeptnApp

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnApp
metadata:
  name: my-app
  namespace: production
spec:
  version: "1.0.0"
  workloads:
  - name: frontend
    version: "1.2.3"
  - name: backend
    version: "2.0.1"
  - name: database
    version: "14.5"
```

### 部署前任务

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnTaskDefinition
metadata:
  name: pre-deployment-checks
  namespace: production
spec:
  retries: 3
  timeout: "5m"
  deno:
    inline:
      code: |
        // 检查外部依赖
        const response = await fetch('https://api.example.com/health');
        if (!response.ok) {
          throw new Error('External dependency not ready');
        }
        console.log('All dependencies ready');
```

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnTaskDefinition
metadata:
  name: security-scan
  namespace: production
spec:
  container:
    name: trivy
    image: aquasec/trivy:latest
    command: ["trivy", "image", "--exit-code", "1", "my-app:latest"]
```

### 部署后分析

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: AnalysisDefinition
metadata:
  name: slo-analysis
  namespace: production
spec:
  objectives:
  - keptnMetricRef:
      name: request-latency
      namespace: production
    target:
      failure:
        greaterThan:
          fixedValue: 1000
      warning:
        greaterThan:
          fixedValue: 500
    weight: 1
  - keptnMetricRef:
      name: error-rate
      namespace: production
    target:
      failure:
        greaterThan:
          fixedValue: 5
      warning:
        greaterThan:
          fixedValue: 2
    weight: 2
  totalScore:
    passPercentage: 90
    warningPercentage: 75
```

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: Analysis
metadata:
  name: post-deployment-analysis
  namespace: production
spec:
  analysisDefinition:
    name: slo-analysis
    namespace: production
  timeframe:
    from: "2025-01-01T00:00:00Z"
    to: "2025-01-01T01:00:00Z"
```

### 指标提供者

```yaml
apiVersion: metrics.keptn.sh/v1
kind: KeptnMetricsProvider
metadata:
  name: prometheus
  namespace: production
spec:
  type: prometheus
  targetServer: "http://prometheus.monitoring.svc.cluster.local:9090"
---
apiVersion: metrics.keptn.sh/v1
kind: KeptnMetric
metadata:
  name: request-latency
  namespace: production
spec:
  provider:
    name: prometheus
  query: 'histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))'
  fetchIntervalSeconds: 30
```

## 多阶段交付

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnAppContext
metadata:
  name: my-app-promotion
  namespace: production
spec:
  preDeploymentTasks:
  - pre-deployment-checks
  - security-scan
  postDeploymentTasks:
  - integration-tests
  postDeploymentEvaluations:
  - slo-analysis
  promotionTasks:
  - name: promote-to-production
    promotion:
      targetCluster: prod-cluster
      namespace: production
```

## 与 ArgoCD 集成

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: my-app
  namespace: argocd
  annotations:
    keptn.sh/app: my-app
    keptn.sh/workload: my-workload
spec:
  project: default
  source:
    repoURL: https://github.com/example/my-app
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

```yaml
# 启用 Keptn 注解的 Namespace
apiVersion: v1
kind: Namespace
metadata:
  name: production
  annotations:
    keptn.sh/lifecycle-controller: "enabled"
```

## DORA 指标

Keptn 自动生成以下 DORA 指标：

| 指标 | 描述 |
|------|------|
| 部署频率 | 应用部署次数 |
| 变更前置时间 | 代码提交到部署的时间 |
| 变更失败率 | 导致失败的部署比例 |
| 恢复时间 | 从失败中恢复的时间 |

```yaml
# Grafana Dashboard 查询
# 部署频率
keptn_deployment_count_total

# 变更前置时间
keptn_lead_time_seconds

# 恢复时间
keptn_recovery_time_seconds
```

## 2025 更新

- **Keptn Lifecycle Toolkit**: 新的轻量级版本
- **Analysis CRD**: 增强的 SLO 分析
- **Promotion Tasks**: 自动化环境晋升
- **多集群支持**: 跨集群应用管理

## 相关资源

- [Keptn 官网](https://keptn.sh/)
- [GitHub](https://github.com/keptn/lifecycle-toolkit)
- [文档](https://keptn.sh/stable/docs/)
