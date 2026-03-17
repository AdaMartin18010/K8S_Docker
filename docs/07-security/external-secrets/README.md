# External Secrets Operator (ESO) - 外部密钥管理

## 概述

External Secrets Operator (ESO) 是一个 Kubernetes Operator，用于将外部密钥管理系统（如 AWS Secrets Manager、HashiCorp Vault、Azure Key Vault 等）的密钥同步到 Kubernetes Secret 中。它解决了 Kubernetes 原生 Secret 安全性不足的问题，实现了密钥的集中管理和自动轮换。

> **关键数据**: ESO 支持 50+ 种外部密钥提供者，是 GitOps 工作流中管理密钥的最佳实践。

## 架构原理

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         External Secrets Operator                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    External Secrets Controller                      │   │
│  │                                                                     │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│  │  │ ExternalSecret│  │ SecretStore  │  │ClusterSecretStore│         │   │
│  │  │   Controller │  │   Controller │  │    Controller    │         │   │
│  │  └──────┬───────┘  └──────────────┘  └──────────────────┘         │   │
│  │         │                                                          │   │
│  │         ▼                                                          │   │
│  │  ┌──────────────────────────────────────────────────────────────┐  │   │
│  │  │                    Provider Interface                         │  │   │
│  │  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐     │  │   │
│  │  │  │  AWS   │ │  GCP   │ │ Azure  │ │ Vault  │ │  ...   │     │  │   │
│  │  │  │Secrets │ │Secret  │ │KeyVault│ │        │ │        │     │  │   │
│  │  │  │Manager │ │Manager │ │        │ │        │ │        │     │  │   │
│  │  │  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘     │  │   │
│  │  └──────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Kubernetes Secret                                │   │
│  │            (自动创建/更新/同步)                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 核心概念

| 资源 | 作用 | 作用域 |
|------|------|--------|
| ExternalSecret | 定义从外部获取哪些密钥，如何映射到 K8s Secret | Namespace |
| SecretStore | 定义外部密钥存储的连接配置 | Namespace |
| ClusterSecretStore | 集群级 SecretStore，可被多个 Namespace 使用 | Cluster |
| PushSecret | 将 K8s Secret 推送到外部存储（双向同步） | Namespace |

## 安装 ESO

### Helm 安装

```bash
# 添加 Helm 仓库
helm repo add external-secrets https://charts.external-secrets.io
helm repo update

# 安装
helm install external-secrets external-secrets/external-secrets \
  --namespace external-secrets \
  --create-namespace \
  --set installCRDs=true

# 验证安装
kubectl get pods -n external-secrets
kubectl get crd | grep external-secrets
```

## 基础使用示例

### 1. AWS Secrets Manager 集成

```yaml
# 1. 创建 AWS 凭证 Secret
apiVersion: v1
kind: Secret
metadata:
  name: aws-credentials
  namespace: default
type: Opaque
stringData:
  access-key-id: AKIAIOSFODNN7EXAMPLE
  secret-access-key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
---
# 2. 创建 SecretStore
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: aws-secretstore
  namespace: default
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
      auth:
        secretRef:
          accessKeyIDSecretRef:
            name: aws-credentials
            key: access-key-id
          secretAccessKeySecretRef:
            name: aws-credentials
            key: secret-access-key
---
# 3. 创建 ExternalSecret
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: db-credentials
  namespace: default
spec:
  refreshInterval: 1h  # 同步间隔
  secretStoreRef:
    name: aws-secretstore
    kind: SecretStore
  target:
    name: database-credentials  # 创建的 K8s Secret 名称
    creationPolicy: Owner
    template:
      type: Opaque
      data:
        connection-string: "postgresql://{{ .username }}:{{ .password }}@{{ .host }}:5432/{{ .database }}"
  data:
  - secretKey: username
    remoteRef:
      key: prod/db-credentials  # AWS Secrets Manager 中的密钥名称
      property: username
  - secretKey: password
    remoteRef:
      key: prod/db-credentials
      property: password
  - secretKey: host
    remoteRef:
      key: prod/db-credentials
      property: host
  - secretKey: database
    remoteRef:
      key: prod/db-credentials
      property: database
```

### 2. HashiCorp Vault 集成

```yaml
# 使用 Kubernetes 认证
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: vault-backend
spec:
  provider:
    vault:
      server: "http://vault.vault-system:8200"
      path: "secret"
      version: "v2"
      auth:
        kubernetes:
          mountPath: "kubernetes"
          role: "external-secrets"
          serviceAccountRef:
            name: external-secrets-sa
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: vault-example
spec:
  refreshInterval: "15s"
  secretStoreRef:
    name: vault-backend
    kind: SecretStore
  target:
    name: vault-secrets
  data:
  - secretKey: api-key
    remoteRef:
      key: secret/data/myapp
      property: api-key
  - secretKey: api-secret
    remoteRef:
      key: secret/data/myapp
      property: api-secret
```

### 3. Azure Key Vault 集成（Workload Identity）

```yaml
# 使用 Azure Workload Identity
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: azure-store
spec:
  provider:
    azurekv:
      authType: WorkloadIdentity
      vaultUrl: "https://my-keyvault.vault.azure.net"
      serviceAccountRef:
        name: workload-identity-sa
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: azure-db-credentials
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: azure-store
    kind: SecretStore
  target:
    name: database-credentials
  data:
  - secretKey: username
    remoteRef:
      key: database-username
  - secretKey: password
    remoteRef:
      key: database-password
```

### 4. ClusterSecretStore（集群级共享）

```yaml
# 创建 ClusterSecretStore
apiVersion: external-secrets.io/v1beta1
kind: ClusterSecretStore
metadata:
  name: global-aws-store
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
      auth:
        jwt:
          serviceAccountRef:
            name: external-secrets
            namespace: external-secrets
  # 条件：只允许特定 namespace 使用
  conditions:
  - namespaces:
    - production
    - staging
    - "app-*"
---
# 在任意 namespace 使用
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: shared-credentials
  namespace: production
spec:
  secretStoreRef:
    name: global-aws-store
    kind: ClusterSecretStore
  target:
    name: shared-secrets
  data:
  - secretKey: api-key
    remoteRef:
      key: shared/api-key
```

## 高级功能

### 模板化（Template）

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: templated-secret
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: my-store
    kind: SecretStore
  target:
    name: app-config
    template:
      type: Opaque
      metadata:
        annotations:
          managed-by: external-secrets
          last-sync: "{{ .metadata.creationTimestamp }}"
      data:
        # 使用 Go template 语法
        config.json: |
          {
            "database": {
              "host": "{{ .db_host }}",
              "port": {{ .db_port }},
              "username": "{{ .db_user }}",
              "password": "{{ .db_pass }}"
            },
            "api": {
              "key": "{{ .api_key }}"
            }
          }
        # 构造连接字符串
        DATABASE_URL: "postgres://{{ .db_user }}:{{ .db_pass }}@{{ .db_host }}:{{ .db_port }}/{{ .db_name }}"
  data:
  - secretKey: db_host
    remoteRef:
      key: prod/db-config
      property: host
  - secretKey: db_port
    remoteRef:
      key: prod/db-config
      property: port
  - secretKey: db_user
    remoteRef:
      key: prod/db-config
      property: username
  - secretKey: db_pass
    remoteRef:
      key: prod/db-config
      property: password
  - secretKey: db_name
    remoteRef:
      key: prod/db-config
      property: database
  - secretKey: api_key
    remoteRef:
      key: prod/api-credentials
      property: key
```

### 批量获取（dataFrom）

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: bulk-secrets
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: my-store
    kind: SecretStore
  target:
    name: all-app-secrets
  dataFrom:
  # 获取所有匹配前缀的密钥
  - find:
      name:
        regexp: "^prod/app-"
  # 提取 JSON 中的所有字段
  - extract:
      key: prod/app-config
```

### PushSecret（双向同步）

```yaml
# 将 K8s Secret 推送到外部存储
apiVersion: external-secrets.io/v1alpha1
kind: PushSecret
metadata:
  name: push-to-aws
spec:
  refreshInterval: 1h
  secretStoreRefs:
  - name: aws-store
    kind: SecretStore
  selector:
    secret:
      name: my-secret  # 源 K8s Secret
  template:
    data:
      connection-string: "{{ .username }}:{{ .password }}@{{ .host }}"
  data:
  - match:
      secretKey: username
      remoteRef:
        remoteKey: prod/synced-credentials
        property: username
  - match:
      secretKey: password
      remoteRef:
        remoteKey: prod/synced-credentials
        property: password
```

## 与 Cert Manager 集成

```yaml
# 自动创建 TLS 证书并同步到外部存储
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: my-cert
spec:
  secretName: my-cert-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
  - app.example.com
---
apiVersion: external-secrets.io/v1alpha1
kind: PushSecret
metadata:
  name: push-cert
spec:
  refreshInterval: 1h
  secretStoreRefs:
  - name: aws-store
    kind: SecretStore
  selector:
    secret:
      name: my-cert-tls
  data:
  - match:
      secretKey: tls.crt
      remoteRef:
        remoteKey: prod/certs/app-tls
        property: certificate
  - match:
      secretKey: tls.key
      remoteRef:
        remoteKey: prod/certs/app-tls
        property: private-key
```

## 多租户安全实践

```yaml
# Namespace 级别的隔离
apiVersion: v1
kind: Namespace
metadata:
  name: team-a
  labels:
    external-secrets.io/managed: "true"
---
# Team A 的 SecretStore（只能访问 team-a 路径）
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: team-a-store
  namespace: team-a
spec:
  provider:
    vault:
      server: "http://vault:8200"
      path: "secret/team-a"  # 限制路径
      auth:
        kubernetes:
          role: "team-a-role"
---
# ClusterSecretStore 的命名空间限制
apiVersion: external-secrets.io/v1beta1
kind: ClusterSecretStore
metadata:
  name: shared-store
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
  conditions:
  - namespaces:
    - production
    - staging
```

## 监控和告警

```yaml
# PrometheusRule
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: external-secrets-alerts
spec:
  groups:
  - name: external-secrets
    rules:
    - alert: ExternalSecretSyncFailed
      expr: externalsecret_status_condition{condition="Ready",status="False"} == 1
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "ExternalSecret {{ $labels.name }} sync failed"
    
    - alert: ExternalSecretNotSynced
      expr: time() - externalsecret_status_sync_timestamp > 7200
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "ExternalSecret {{ $labels.name }} not synced for 2 hours"
```

## 故障排查

```bash
# 查看 ExternalSecret 状态
kubectl get externalsecrets
kubectl describe externalsecret <name>

# 查看生成的 Secret
kubectl get secret <target-name> -o yaml

# 查看控制器日志
kubectl logs -n external-secrets deployment/external-secrets

# 验证外部连接
kubectl exec -it deployment/external-secrets -- nslookup vault.vault-system
```

## 总结

| 场景 | 推荐方案 |
|------|----------|
| AWS 环境 | AWS Secrets Manager + IRSA |
| Azure 环境 | Azure Key Vault + Workload Identity |
| GCP 环境 | GCP Secret Manager + Workload Identity |
| 混合云/多云 | HashiCorp Vault |
| 简单场景 | 云厂商原生服务 |
| 复杂权限管理 | Vault + Kubernetes 认证 |

ESO 是现代 Kubernetes 集群管理密钥的标准方式，它将密钥管理从应用配置中解耦，实现了真正的 GitOps 工作流。
