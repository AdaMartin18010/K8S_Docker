# cert-manager - Kubernetes 证书管理

## 概述

cert-manager 是 Kubernetes 上最广泛使用的证书管理工具，支持自动颁发和管理 TLS 证书，与 Let's Encrypt、HashiCorp Vault、私有 CA 等集成。

## 核心特性

- **自动证书颁发**: 支持 ACME (Let's Encrypt)、自签名、CA、Vault
- **自动续期**: 在证书过期前自动续期
- **多 Issuer 支持**: ClusterIssuer (集群范围) 和 Issuer (命名空间范围)
- **Gateway API 集成**: 支持 Kubernetes Gateway API
- **CSI 驱动**: 支持证书自动注入 Pod

## 安装部署

### Helm 安装

```bash
# 添加 Helm 仓库
helm repo add jetstack https://charts.jetstack.io
helm repo update

# 安装 cert-manager
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.17.0 \
  --set installCRDs=true \
  --set config.enableGatewayAPI=true

# 验证安装
kubectl get pods -n cert-manager
```

### Gateway API 支持

```bash
# 启用 Gateway API 支持
helm upgrade cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --set config.enableGatewayAPI=true

# 重启生效
kubectl rollout restart deployment cert-manager -n cert-manager
```

## Issuer 配置

### Let's Encrypt ClusterIssuer

```yaml
# letsencrypt-staging.yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    email: admin@example.com
    privateKeySecretRef:
      name: letsencrypt-staging
    solvers:
    - http01:
        ingress:
          class: nginx
    - http01:
        gatewayHTTPRoute:
          parentRefs:
          - name: external-gateway
            namespace: envoy-gateway-system
            kind: Gateway
---
# letsencrypt-production.yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    # 多域名使用不同挑战方式
    - selector:
        dnsNames:
        - "*.example.com"
        - "example.com"
      dns01:
        route53:
          region: us-west-2
          hostedZoneID: Z123456789
    - http01:
        ingress:
          class: nginx
```

### Vault Issuer

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: vault-issuer
spec:
  vault:
    server: https://vault.internal:8200
    path: pki/sign/example-dot-com
    auth:
      kubernetes:
        mountPath: /v1/auth/kubernetes
        role: cert-manager
        secretRef:
          name: cert-manager-vault-token
          key: token
```

## Gateway API 集成

### 自动证书颁发

```yaml
# Gateway 自动证书
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: external-gateway
  namespace: envoy-gateway-system
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    cert-manager.io/acme-challenge-type: http01
spec:
  gatewayClassName: envoy
  listeners:
  - name: https
    protocol: HTTPS
    port: 443
    hostname: "*.example.com"
    tls:
      mode: Terminate
      certificateRefs:
      - name: wildcard-example-com
        namespace: envoy-gateway-system
    allowedRoutes:
      namespaces:
        from: All
  - name: http
    protocol: HTTP
    port: 80
    allowedRoutes:
      namespaces:
        from: All
```

### HTTPRoute 证书

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: api-route
  namespace: production
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  parentRefs:
  - name: external-gateway
    namespace: envoy-gateway-system
    sectionName: https
  hostnames:
  - api.example.com
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /v1
    backendRefs:
    - name: api-service
      port: 8080
```

## Certificate 资源

### 基本 Certificate

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com
  namespace: production
spec:
  secretName: example-com-tls
  duration: 2160h  # 90天
  renewBefore: 720h  # 30天前续期
  subject:
    organizations:
    - Example Corp
  isCA: false
  privateKey:
    algorithm: RSA
    encoding: PKCS1
    size: 2048
  usages:
  - server auth
  - client auth
  dnsNames:
  - example.com
  - www.example.com
  - api.example.com
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
    group: cert-manager.io
```

### 带 JKS/PKCS12 的证书

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: java-app-cert
spec:
  secretName: java-app-tls
  dnsNames:
  - java-app.example.com
  issuerRef:
    name: internal-ca
    kind: ClusterIssuer
  keystores:
    jks:
      create: true
      passwordSecretRef:
        name: jks-password
        key: password
    pkcs12:
      create: true
      passwordSecretRef:
        name: pkcs12-password
        key: password
```

## 2025 新特性

- **Gateway API GA**: 完整支持 Gateway API v1
- **Name Constraints**: CA 证书支持名称约束 (Beta)
- **CA Injector Merging**: CA 证书合并而非替换
- **Password as Literal**: 支持明文密码配置
- **RSA 签名增强**: 4096位+密钥使用 SHA-512

## 监控指标

```yaml
# ServiceMonitor
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cert-manager
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: cert-manager
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
```

## 相关资源

- [cert-manager 官网](https://cert-manager.io/)
- [GitHub](https://github.com/cert-manager/cert-manager)
- [Gateway API 集成](https://cert-manager.io/docs/usage/gateway/)
