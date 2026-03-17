# Kubernetes 安全体系

> 生产级安全加固指南

---

## 本章内容

1. [安全概述](./security-overview.md)
2. [RBAC 权限控制](./rbac.md)
3. [Pod 安全标准](./pod-security-standards.md)
4. [NetworkPolicy](./network-policy.md)
5. [Secret 管理](./secrets.md)
6. [镜像安全](./image-security.md)

---

## 安全层次

```
┌─────────────────────────────────────────┐
│  1. 镜像安全                             │
│     - 漏洞扫描                           │
│     - 镜像签名                           │
├─────────────────────────────────────────┤
│  2. 运行时安全                           │
│     - Pod Security Standards            │
│     - Seccomp / AppArmor                │
├─────────────────────────────────────────┤
│  3. 网络安全                             │
│     - NetworkPolicy                     │
│     - 服务网格 mTLS                     │
├─────────────────────────────────────────┤
│  4. 访问控制                             │
│     - RBAC                              │
│     - 准入控制器                         │
├─────────────────────────────────────────┤
│  5. 审计与合规                           │
│     - 审计日志                           │
│     - Falco 运行时检测                   │
└─────────────────────────────────────────┘
```

---

## Pod Security Standards (替代 PSP)

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/enforce-version: latest
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

---

## RBAC 最小权限示例

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: web-app
automountServiceAccountToken: false
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: web-app-role
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get", "list"]
    resourceNames: ["web-app-config"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: web-app-binding
subjects:
  - kind: ServiceAccount
    name: web-app
roleRef:
  kind: Role
  name: web-app-role
  apiGroup: rbac.authorization.k8s.io
```

---

## 安全检查清单

- [ ] 使用 Pod Security Standards
- [ ] 配置 NetworkPolicy
- [ ] 最小权限 RBAC
- [ ] Secret 加密存储
- [ ] 镜像漏洞扫描
- [ ] 运行时安全监控

---

## 关联代码

- [examples/kubernetes/06-security/](../../examples/kubernetes/06-security/)
