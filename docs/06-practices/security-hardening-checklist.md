# Kubernetes 安全加固检查清单

> 生产环境安全基线配置

---

## 🔐 安全加固层次

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          安全加固层次                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Layer 1: 镜像安全                                                            │
│  ════════════════════════════════════════════════════════════════════════   │
│  │ • 使用最小化基础镜像                                                       │
│  │ • 镜像签名和验证                                                           │
│  │ • 漏洞扫描                                                                 │
│  │ • SBOM 生成                                                                │
│                                                                              │
│  Layer 2: 运行时安全                                                          │
│  ════════════════════════════════════════════════════════════════════════   │
│  │ • Pod 安全标准                                                             │
│  │ • 安全上下文配置                                                           │
│  │ • 运行时监控 (Falco)                                                       │
│  │ • 网络策略                                                                 │
│                                                                              │
│  Layer 3: 准入控制                                                            │
│  ════════════════════════════════════════════════════════════════════════   │
│  │ • 镜像验证                                                                 │
│  │ • 策略执行 (OPA/Kyverno)                                                   │
│  │ • 资源限制验证                                                             │
│                                                                              │
│  Layer 4: 网络安全                                                            │
│  ════════════════════════════════════════════════════════════════════════   │
│  │ • mTLS (服务网格)                                                          │
│  │ • 网络分段                                                                 │
│  │ • API 安全                                                                 │
│                                                                              │
│  Layer 5: 密钥管理                                                            │
│  ════════════════════════════════════════════════════════════════════════   │
│  │ • 外部密钥管理 (Vault)                                                     │
│  │ • 密钥轮换                                                                 │
│  │ • 最小权限原则                                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. 镜像安全

### 基础镜像选择

```dockerfile
# ❌ 不推荐: 完整操作系统镜像
FROM ubuntu:latest

# ✅ 推荐: 最小化镜像
FROM distroless/static-debian11
# 或
FROM alpine:3.18
# 或
FROM scratch  # 对于静态编译的应用

# ✅ 使用特定版本标签
FROM nginx:1.25.3-alpine
```

### Dockerfile 安全最佳实践

```dockerfile
# ✅ 使用非 root 用户
FROM node:18-alpine
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001
USER nextjs

# ✅ 多阶段构建，减少攻击面
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
COPY --from=builder /app/main /main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/main"]

# ✅ 不安装不必要的包
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*
```

### 镜像签名与验证

```bash
# 签名镜像
cosign sign --yes myregistry.io/myapp:v1.0.0

# 验证签名
cosign verify myregistry.io/myapp:v1.0.0 \
  --certificate-identity=user@example.com \
  --certificate-oidc-issuer=https://accounts.google.com
```

---

## 2. Pod 安全标准

### Pod Security Standards (PSS)

```yaml
# 在命名空间上启用 Pod 安全标准
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

| 级别 | 说明 | 适用场景 |
|------|------|----------|
| **privileged** | 无限制 | 系统级 Pod |
| **baseline** | 最小限制 | 一般应用 |
| **restricted** | 最严格 | 安全敏感应用 |

### 安全上下文配置

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: secure-pod
spec:
  securityContext:
    runAsNonRoot: true          # 必须以非 root 运行
    runAsUser: 1000             # 指定用户 ID
    runAsGroup: 1000            # 指定组 ID
    fsGroup: 1000               # 卷挂载的组
    seccompProfile:
      type: RuntimeDefault      # 使用默认 seccomp 配置
  containers:
  - name: app
    image: myapp:v1.0
    securityContext:
      allowPrivilegeEscalation: false   # 禁止特权提升
      readOnlyRootFilesystem: true      # 只读根文件系统
      capabilities:
        drop:
        - ALL                        # 删除所有 capabilities
      resources:
        limits:
          memory: "256Mi"
          cpu: "500m"
```

### 检查清单

```markdown
## Pod 安全配置检查清单

- [ ] 使用非 root 用户运行
- [ ] 启用只读根文件系统
- [ ] 删除所有 capabilities
- [ ] 禁用特权提升
- [ ] 配置资源限制
- [ ] 使用 seccomp 配置
- [ ] 启用 AppArmor/SELinux (如适用)
```

---

## 3. 网络策略

### 默认拒绝策略

```yaml
# 默认拒绝所有入站流量
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
  namespace: production
spec:
  podSelector: {}
  policyTypes:
  - Ingress
---
# 默认拒绝所有出站流量
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-egress
  namespace: production
spec:
  podSelector: {}
  policyTypes:
  - Egress
```

### 应用间访问控制

```yaml
# 只允许前端访问后端
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: backend-policy
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
```

---

## 4. RBAC 配置

### 最小权限原则

```yaml
# 创建最小权限 ServiceAccount
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app-sa
  namespace: production
---
# 只允许获取特定 ConfigMap
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: app-role
  namespace: production
rules:
- apiGroups: [""]
  resources: ["configmaps"]
  resourceNames: ["app-config"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: app-rolebinding
  namespace: production
subjects:
- kind: ServiceAccount
  name: app-sa
  namespace: production
roleRef:
  kind: Role
  name: app-role
  apiGroup: rbac.authorization.k8s.io
```

### 避免使用 cluster-admin

```yaml
# ❌ 不要这样做
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: unsafe-binding
subjects:
- kind: ServiceAccount
  name: default
  namespace: production
roleRef:
  kind: ClusterRole
  name: cluster-admin  # 危险!
  apiGroup: rbac.authorization.k8s.io
```

---

## 5. 准入控制策略

### Kyverno 策略示例

```yaml
# 禁止最新标签
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: disallow-latest-tag
spec:
  validationFailureAction: enforce
  rules:
  - name: validate-image-tag
    match:
      any:
      - resources:
          kinds:
          - Pod
    validate:
      message: "Using 'latest' tag is not allowed"
      pattern:
        spec:
          containers:
          - image: "!*:latest"
---
# 要求资源限制
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-resources-limits
spec:
  validationFailureAction: enforce
  rules:
  - name: validate-resources
    match:
      any:
      - resources:
          kinds:
          - Pod
    validate:
      message: "CPU and memory limits are required"
      pattern:
        spec:
          containers:
          - resources:
              limits:
                memory: "?*"
                cpu: "?*"
```

### OPA/Gatekeeper 策略

```yaml
# 禁止特权容器
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: k8spspprivilegedcontainer
spec:
  crd:
    spec:
      names:
        kind: K8sPSPPrivilegedContainer
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package k8spspprivilegedcontainer
        violation[{"msg": msg}] {
          input.review.object.spec.containers[_].securityContext.privileged
          msg := "Privileged containers are not allowed"
        }
```

---

## 6. 密钥管理

### External Secrets 配置

```yaml
# 从 Vault 获取密钥
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: db-credentials
  namespace: production
spec:
  refreshInterval: 1h
  secretStoreRef:
    kind: ClusterSecretStore
    name: vault-backend
  target:
    name: db-credentials
    creationPolicy: Owner
  data:
  - secretKey: username
    remoteRef:
      key: secret/data/db
      property: username
  - secretKey: password
    remoteRef:
      key: secret/data/db
      property: password
```

### 密钥使用最佳实践

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: secure-app
spec:
  containers:
  - name: app
    image: myapp:v1.0
    env:
    - name: DB_USERNAME
      valueFrom:
        secretKeyRef:
          name: db-credentials
          key: username
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: db-credentials
          key: password
    volumeMounts:
    - name: secrets
      mountPath: /etc/secrets
      readOnly: true
  volumes:
  - name: secrets
    csi:
      driver: secrets-store.csi.k8s.io
      readOnly: true
      volumeAttributes:
        secretProviderClass: vault-db
```

---

## 7. 运行时安全

### Falco 规则示例

```yaml
# 检测容器内运行 shell
- rule: Terminal Shell in Container
  desc: Detect shell started inside a container
  condition: spawned_process and shell_procs and container
  output: "Shell started in container (user=%user.name container=%container.name)"
  priority: WARNING

# 检测敏感文件访问
- rule: Read Sensitive File
  desc: Detect access to sensitive files
  condition: open_read and (etc_passwd or etc_shadow)
  output: "Sensitive file accessed (user=%user.name file=%fd.name)"
  priority: WARNING
```

### Tetragon 配置

```yaml
# 监控特权升级
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: privilege-escalation
spec:
  kprobes:
  - call: "__x64_sys_setuid"
    syscall: true
    args:
    - index: 0
      type: int
```

---

## 8. 审计日志

### 启用审计日志

```yaml
# kube-apiserver 配置
apiVersion: apiserver.config.k8s.io/v1
kind: Policy
rules:
# 记录所有请求
- level: RequestResponse
  resources:
  - group: ""
    resources: ["pods", "secrets", "configmaps"]

# 记录元数据
- level: Metadata
  omitStages:
  - RequestReceived

# 不记录健康检查
- level: None
  nonResourceURLs:
  - /healthz*
  - /version
```

---

## 9. 安全加固检查清单

### 集群级别

```markdown
- [ ] 启用 Pod 安全标准 (PSS)
- [ ] 配置 NetworkPolicy 默认拒绝
- [ ] 启用审计日志
- [ ] 配置准入控制器
- [ ] 禁用匿名访问
- [ ] 启用 mTLS (服务网格)
- [ ] 定期轮换证书
- [ ] 启用 Falco/Tetragon 运行时监控
```

### 命名空间级别

```markdown
- [ ] 设置 Pod 安全标准标签
- [ ] 创建 NetworkPolicy
- [ ] 配置 ResourceQuota
- [ ] 配置 LimitRange
- [ ] 创建专用 ServiceAccount
```

### Pod 级别

```markdown
- [ ] 使用非 root 用户
- [ ] 配置 SecurityContext
- [ ] 设置资源限制
- [ ] 使用镜像签名
- [ ] 配置健康检查
- [ ] 使用外部密钥
```

---

## 10. 安全扫描工具

```bash
# Trivy 镜像扫描
trivy image myapp:latest

# Trivy 文件系统扫描
trivy fs .

# Trivy K8s 配置扫描
trivy config k8s/

# Kube-bench CIS 基准测试
kube-bench

# Kubesec 配置安全评分
kubesec scan pod.yaml

# Prowler 云安全扫描
prowler kubernetes
```

---

## 参考资源

- [CIS Kubernetes Benchmark](https://www.cisecurity.org/benchmark/kubernetes)
- [NSA/CISA Kubernetes Hardening Guide](https://media.defense.gov/2022/Aug/29/2003066362/-1/-1/0/CTR_KUBERNETES_HARDENING_GUIDANCE_1.2_20220829.PDF)
- [OWASP Kubernetes Security](https://cheatsheetseries.owasp.org/cheatsheets/Kubernetes_Security_Cheat_Sheet.html)
