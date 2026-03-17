# Crossplane - 云原生控制平面框架

## 概述

Crossplane 是一个 CNCF 毕业项目，用于构建云原生控制平面的框架。它扩展 Kubernetes 来管理任何资源，将声明式 API 和协调模式带给基础设施和应用管理。

## 核心特性

| 特性 | 描述 |
|------|------|
| 声明式配置 | Kubernetes 风格的 API 驱动管理 |
| 自修复 | 自动纠正配置漂移 |
| 统一平台 | 应用和基础设施配置统一管理 |
| GitOps 集成 | 与 CI/CD 和 GitOps 最佳实践集成 |
| 多租户 | API 级别策略和权限控制 |

## 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                    Crossplane 控制平面                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │ Composite   │    │  Managed    │    │  Provider   │         │
│  │ Resource    │───▶│  Resource   │───▶│  (AWS/Azure)│         │
│  │ (XR)        │    │  (MR)       │    │             │         │
│  └─────────────┘    └─────────────┘    └──────┬──────┘         │
│         │                                      │                │
│         ▼                                      ▼                │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Kubernetes API                        │   │
│  │  ┌──────────────────────────────────────────────────┐   │   │
│  │  │            Crossplane Composition                 │   │   │
│  │  │  - 定义 XRD (CompositeResourceDefinition)         │   │   │
│  │  │  - 创建 Composition                             │   │   │
│  │  │  - 部署 Claim                                    │   │   │
│  │  └──────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │  Helm       │    │  K8s        │    │  Cloud      │         │
│  │  Charts     │    │  Resources  │    │  Resources  │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 安装部署

### Helm 安装

```bash
# 添加仓库
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update

# 安装 Crossplane
helm install crossplane crossplane-stable/crossplane \
  --namespace crossplane-system \
  --create-namespace \
  --version 1.18.0

# 验证安装
kubectl get pods -n crossplane-system

# 安装 CLI
curl -sL https://raw.githubusercontent.com/crossplane/crossplane/master/install.sh | sh
sudo mv crossplane /usr/local/bin/
```

### 配置 Provider

```yaml
# provider-aws.yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-aws-ec2
spec:
  package: xpkg.upbound.io/upbound/provider-aws-ec2:v1.18.0
---
# provider-config.yaml
apiVersion: aws.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: aws-creds
      key: credentials
```

## 定义 Composite Resource

```yaml
# xrd.yaml - Composite Resource Definition
apiVersion: apiextensions.crossplane.io/v1
kind: CompositeResourceDefinition
metadata:
  name: xpostgresqlinstances.database.example.com
spec:
  group: database.example.com
  names:
    kind: XPostgreSQLInstance
    plural: xpostgresqlinstances
  claimNames:
    kind: PostgreSQLInstance
    plural: postgresqlinstances
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
                    default: 10
                  version:
                    type: string
                    default: "14"
                required:
                - storageGB
            required:
            - parameters
```

```yaml
# composition.yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: xpostgresqlinstances.aws.database.example.com
  labels:
    provider: aws
    region: us-west-2
spec:
  compositeTypeRef:
    apiVersion: database.example.com/v1alpha1
    kind: XPostgreSQLInstance
  resources:
  - name: subnetGroup
    base:
      apiVersion: rds.aws.upbound.io/v1beta1
      kind: SubnetGroup
      spec:
        forProvider:
          region: us-west-2
          subnetIds:
          - subnet-abc123
          - subnet-def456
  - name: rdsInstance
    base:
      apiVersion: rds.aws.upbound.io/v1beta1
      kind: Instance
      spec:
        forProvider:
          region: us-west-2
          engine: postgres
          instanceClass: db.t3.micro
          allocatedStorage: 10
          username: masteruser
          skipFinalSnapshot: true
          dbSubnetGroupNameSelector:
            matchControllerRef: true
        writeConnectionSecretToRef:
          namespace: crossplane-system
    patches:
    - fromFieldPath: spec.parameters.storageGB
      toFieldPath: spec.forProvider.allocatedStorage
    - fromFieldPath: spec.parameters.version
      toFieldPath: spec.forProvider.engineVersion
    - fromFieldPath: metadata.uid
      toFieldPath: spec.forProvider.identifier
      transforms:
      - type: string
        string:
          fmt: "postgres-%s"
```

## 使用 Claim

```yaml
# claim.yaml
apiVersion: database.example.com/v1alpha1
kind: PostgreSQLInstance
metadata:
  name: my-db
  namespace: production
spec:
  parameters:
    storageGB: 20
    version: "15"
  compositionSelector:
    matchLabels:
      provider: aws
      region: us-west-2
  writeConnectionSecretToRef:
    name: my-db-connection
```

## 多环境配置

```yaml
# dev-composition.yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: xpostgresqlinstances.aws.dev
  labels:
    environment: dev
spec:
  compositeTypeRef:
    apiVersion: database.example.com/v1alpha1
    kind: XPostgreSQLInstance
  resources:
  - name: rdsInstance
    base:
      apiVersion: rds.aws.upbound.io/v1beta1
      kind: Instance
      spec:
        forProvider:
          instanceClass: db.t3.micro
          allocatedStorage: 10
          skipFinalSnapshot: true
```

```yaml
# prod-composition.yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: xpostgresqlinstances.aws.prod
  labels:
    environment: prod
spec:
  compositeTypeRef:
    apiVersion: database.example.com/v1alpha1
    kind: XPostgreSQLInstance
  resources:
  - name: rdsInstance
    base:
      apiVersion: rds.aws.upbound.io/v1beta1
      kind: Instance
      spec:
        forProvider:
          instanceClass: db.r5.xlarge
          allocatedStorage: 100
          multiAz: true
          backupRetentionPeriod: 7
          deletionProtection: true
```

## 2025 新特性

- **Crossplane v2**: 应用和基础设施统一管理
- **CNCF 毕业项目**: 2025年11月 CNCF 毕业
- **Intelligent Control Plane**: AI 增强运维
- **Upbound Crossplane (UXP) 2.0**: 企业级功能

## 相关资源

- [Crossplane 官网](https://www.crossplane.io/)
- [Crossplane 文档](https://docs.crossplane.io/)
- [Upbound](https://www.upbound.io/)
