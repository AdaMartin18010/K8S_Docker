# 平台工程 (Platform Engineering)

> 构建内部开发者平台，提升开发效率

---

## 什么是平台工程？

平台工程是构建和运营内部开发者平台 (IDP) 的学科，通过提供自助服务能力和降低认知负载，加速软件交付。

```
┌─────────────────────────────────────────────────────────────┐
│              DevOps vs 平台工程                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  DevOps                          平台工程                     │
│  ────────────────────────────────────────────────────────   │
│                                                              │
│  "You build it,               "平台团队构建工具              │
│   you run it"                 开发者自助服务"               │
│                                                              │
│  每个团队管理自己的             标准化的平台服务              │
│  基础设施和工具                                               │
│                                                              │
│  高认知负载                    降低认知负载                   │
│  (需要掌握 K8s/Terraform等)    (通过抽象和自助服务)           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 内部开发者平台 (IDP) 组成

```
┌─────────────────────────────────────────────────────────────┐
│                    内部开发者平台架构                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           开发者门户 (Developer Portal)               │  │
│  │  • Backstage / Port / Cortex                          │  │
│  │  • 服务目录、文档、脚手架模板                          │  │
│  └─────────────────────┬────────────────────────────────┘  │
│                        │                                     │
│  ┌─────────────────────▼────────────────────────────────┐  │
│  │           平台编排层 (Platform Orchestration)         │  │
│  │  • Crossplane / Terraform                             │  │
│  │  • 基础设施即代码、策略即代码                          │  │
│  └─────────────────────┬────────────────────────────────┘  │
│                        │                                     │
│  ┌─────────────────────▼────────────────────────────────┐  │
│  │           GitOps 交付 (GitOps Delivery)               │  │
│  │  • ArgoCD / Flux                                      │  │
│  │  • 持续部署、渐进式发布                                │  │
│  └─────────────────────┬────────────────────────────────┘  │
│                        │                                     │
│  ┌─────────────────────▼────────────────────────────────┐  │
│  │           可观测性 (Observability)                    │  │
│  │  • Prometheus / Grafana / Jaeger                      │  │
│  │  • 统一监控、日志、追踪                                  │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Backstage - 开发者门户

Spotify 开源的开发者门户框架，已成为 CNCF 项目。

### 核心功能

| 功能 | 说明 |
|------|------|
| **Software Catalog** | 服务目录，管理所有服务、API、库 |
| **Scaffolder** | 脚手架，快速创建符合标准的新项目 |
| **TechDocs** | 技术文档，文档即代码 |
| **Plugins** | 插件生态，集成各种工具 |

### 安装

```bash
# 使用 npx 创建 Backstage 应用
npx @backstage/create-app@latest

# 启动开发服务器
cd my-backstage-app
yarn dev
```

### 实体定义

```yaml
# catalog-info.yaml
apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: payment-service
  description: Payment processing service
  annotations:
    github.com/project-slug: myorg/payment-service
    backstage.io/techdocs-ref: dir:.
    argocd/app-name: payment-service
    grafana/dashboard-selector: "title == 'Payment Service'"
spec:
  type: service
  lifecycle: production
  owner: payments-team
  system: ecommerce
  dependsOn:
    - component:database/payments-db
    - component:queue/payment-events
  providesApis:
    - payment-api
```

---

## Crossplane - 云原生控制平面

使用 Kubernetes API 管理云资源。

```yaml
apiVersion: apiextensions.crossplane.io/v1
kind: CompositeResourceDefinition
metadata:
  name: xpostgresqls.database.example.org
spec:
  group: database.example.org
  names:
    kind: XPostgreSQL
    plural: xpostgresqls
  claimNames:
    kind: PostgreSQL
    plural: postgresqls
  versions:
    - name: v1alpha1
      served: true
      referenceable: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                parameters:
                  type: object
                  properties:
                    storageGB:
                      type: integer
                    version:
                      type: string
                      enum: ["13", "14", "15"]
---
# 开发者自助服务
apiVersion: database.example.org/v1alpha1
kind: PostgreSQL
metadata:
  name: my-app-db
  namespace: production
spec:
  parameters:
    storageGB: 100
    version: "15"
  compositionRef:
    name: aws-postgresql
```

---

## 平台工程成熟度模型

| 级别 | 特征 | 工具 |
|------|------|------|
| **1. 初始** | 手动配置，脚本管理 | Shell, Ansible |
| **2. 可重复** | IaC，版本控制 | Terraform, Pulumi |
| **3. 定义** | GitOps，自动化流水线 | ArgoCD, Flux |
| **4. 管理** | 自助服务，策略即代码 | Backstage, OPA |
| **5. 优化** | AI 辅助，全面可观测 | 智能平台 |

---

## 2025 趋势

### AI 驱动的平台

- **AI 辅助开发**: 自动生成配置、检测配置错误
- **智能容量规划**: ML 预测资源需求
- **自动故障修复**: AI 驱动的自愈系统

### 平台即产品

- 平台团队采用产品经理思维
- 用户研究、路线图、反馈循环
- 度量开发者体验和平台采用率

### Pareto 原则 (80/20)

- 覆盖 80% 的用例即可
- 允许 20% 的例外通过自定义解决
- 平衡标准化和灵活性

---

## 实施路线图

```
Phase 1: 基础设施自动化 (3-6 个月)
  - 标准化 Terraform 模块
  - 实施 GitOps (ArgoCD)
  - 建立基础可观测性

Phase 2: 自助服务能力 (6-12 个月)
  - 部署 Backstage
  - 创建黄金路径模板
  - 实施策略即代码

Phase 3: 智能化平台 (12-18 个月)
  - AI 辅助配置
  - 自动容量管理
  - 高级成本优化

Phase 4: 持续优化 (持续)
  - 度量开发者体验
  - 持续改进平台
  - 扩展平台能力
```

---

## 工具生态

| 类别 | 工具 |
|------|------|
| **开发者门户** | Backstage, Port, Cortex, OpsLevel |
| **基础设施编排** | Crossplane, Terraform, Pulumi |
| **GitOps** | ArgoCD, Flux, Tekton |
| **策略** | OPA, Kyverno, ValidatingAdmissionPolicy |
| **可观测性** | Grafana, Datadog, Dynatrace |
| **成本管理** | Kubecost, Vantage, CloudHealth |
