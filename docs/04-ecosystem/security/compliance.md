# 云原生安全合规

> Kubernetes 安全基线与合规框架

---

## CIS Kubernetes Benchmark

CIS (Center for Internet Security) 基线是 K8s 安全的事实标准。

```
┌─────────────────────────────────────────────────────────────┐
│              CIS Benchmark 覆盖范围                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Control Plane Components (控制平面组件)                  │
│     ├── 1.1 API Server (34 项检查)                          │
│     ├── 1.2 Scheduler (2 项检查)                            │
│     ├── 1.3 Controller Manager (8 项检查)                   │
│     └── 1.4 etcd (7 项检查)                                 │
│                                                              │
│  2. Control Plane Configuration (控制平面配置)               │
│     ├── 2.1 Authentication (4 项检查)                       │
│     └── 2.2 Logging (2 项检查)                              │
│                                                              │
│  3. Worker Nodes (工作节点)                                  │
│     ├── 3.1 Kubelet (13 项检查)                             │
│     └── 3.2 Configuration (5 项检查)                        │
│                                                              │
│  4. Policies (策略)                                          │
│     ├── 4.1 RBAC (7 项检查)                                 │
│     ├── 4.2 Pod Security (8 项检查)                         │
│     └── 4.3 Network (3 项检查)                              │
│                                                              │
│  总计: 100+ 项安全检查                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## kube-bench

CIS Benchmark 的自动化检查工具。

```bash
# 运行 CIS 检查
kubectl run kube-bench --image=aquasec/kube-bench:latest \
  --restart=Never --rm -i -- \
  kube-bench run --targets node

# 输出结果
[INFO] 4 Worker Node Security Configuration
[INFO] 4.1 Worker Node Configuration Files
[PASS] 4.1.1 Ensure that the kubelet service file permissions are set to 600
[PASS] 4.1.2 Ensure that the kubelet service file ownership is set to root:root
[FAIL] 4.1.5 Ensure that the --kubeconfig kubelet.conf file permissions are set to 600
```

---

## Pod Security Standards (PSS)

K8s 内置的 Pod 安全标准，替代已废弃的 PSP。

```yaml
# 启用 Baseline 标准
apiVersion: v1
kind: Namespace
metadata:
  name: restricted-ns
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/enforce-version: latest
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

### 三个级别

| 级别 | 说明 | 适用场景 |
|------|------|----------|
| **Privileged** | 无限制 | 系统组件 |
| **Baseline** | 最小限制 | 一般应用 |
| **Restricted** | 最严格 | 安全敏感 |

### Restricted 要求

```yaml
apiVersion: v1
kind: Pod
spec:
  securityContext:
    runAsNonRoot: true
    seccompProfile:
      type: RuntimeDefault
  containers:
    - name: app
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop: ["ALL"]
        readOnlyRootFilesystem: true
        runAsUser: 1000
        runAsGroup: 1000
        seLinuxOptions:
          type: container_t
```

---

## 合规框架

### 等保 2.0 (中国)

```
┌─────────────────────────────────────────────────────────────┐
│              等保 2.0 三级要求 (K8s 相关)                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  8.1.4.1 身份鉴别                                            │
│  • 应启用双因素认证                                          │
│  • 应配置账号锁定策略                                        │
│  • 应启用审计日志                                            │
│                                                              │
│  8.1.4.2 访问控制                                            │
│  • 应配置 RBAC 最小权限                                      │
│  • 应启用 NetworkPolicy                                      │
│  • 应配置 Pod Security Standards                             │
│                                                              │
│  8.1.4.3 安全审计                                            │
│  • 应启用 API Server 审计日志                                │
│  • 日志留存 >= 180 天                                        │
│  • 日志集中存储                                              │
│                                                              │
│  8.1.4.4 入侵防范                                            │
│  • 应部署 Falco/Tetragon 运行时安全                          │
│  • 应启用镜像安全扫描                                        │
│  • 应配置资源限制                                            │
└─────────────────────────────────────────────────────────────┘
```

### SOC 2 / ISO 27001

```yaml
# 审计日志配置
apiVersion: v1
kind: Policy
rules:
  - level: Metadata
    resources:
      - group: ""
        resources: ["pods"]
  - level: RequestResponse
    resources:
      - group: "rbac.authorization.k8s.io"
        resources: ["roles", "rolebindings"]
  - level: Request
    users: ["admin"]
    verbs: ["delete"]
```

---

## 镜像安全扫描

### Trivy 集成

```yaml
# CI/CD 中扫描
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: image-scan
spec:
  steps:
    - name: scan
      image: aquasec/trivy:latest
      command:
        - trivy
        - image
        - --severity=HIGH,CRITICAL
        - --exit-code=1
        - $(params.image)
```

### 准入控制扫描

```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: image-scanner
webhooks:
  - name: scanner.security.io
    rules:
      - apiGroups: [""]
        apiVersions: ["v1"]
        operations: ["CREATE", "UPDATE"]
        resources: ["pods"]
    clientConfig:
      service:
        name: scanner
        namespace: security
```

---

## 运行时安全

### Falco

```yaml
# 检测异常行为
- rule: Unauthorized K8s API Access
  desc: Detect unauthorized access to Kubernetes API
  condition: |
    spawned_process and
    container and
    (proc.name in (kubectl, helm) or
     proc.cmdline contains "kube-apiserver")
  output: |
    Unauthorized K8s API access
    user=%user.name command=%proc.cmdline
  priority: WARNING
```

### Tetragon (Cilium)

```yaml
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: detect-symlink
spec:
  kprobes:
    - call: "sys_symlinkat"
      syscall: true
      args:
        - index: 0
          type: "string"
        - index: 1
          type: "fd"
      selectors:
        - matchBinaries:
            - operator: "In"
              values:
                - "/bin/bash"
```

---

## 安全扫描清单

```bash
# 1. CIS 基线检查
kube-bench run

# 2. 镜像漏洞扫描
trivy image myapp:latest

# 3. 配置安全扫描
checkov -d ./k8s-manifests/

# 4. 运行时安全检测
falco -r /etc/falco/falco_rules.yaml

# 5. RBAC 审计
kubectl audit rbac

# 6. 网络策略验证
cilium connectivity test
```
