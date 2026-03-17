# RKE2 - Rancher Kubernetes Engine 2

## 概述

RKE2（又名 RKE Government）是 Rancher 的下一代 Kubernetes 发行版，专为高安全性和合规性场景设计，特别适合美国政府机构和金融企业等对安全要求极高的环境。

## 核心特性

### 安全合规
- **CIS Benchmark**: 默认配置通过 CIS Kubernetes Benchmark v1.7/v1.8
- **FIPS 140-2**: 支持 FIPS 140-2 合规认证
- **SELinux**: 完整的 SELinux 策略和 MCS 标签强制执行
- **CVE 扫描**: 构建流程中使用 Trivy 定期扫描组件漏洞

### 与 K3s/RKE1 对比

| 特性 | RKE1 | K3s | RKE2 |
|------|------|-----|------|
| 架构 | Docker 容器 | 单二进制 | 静态 Pod |
| 容器运行时 | Docker | Containerd | Containerd |
| CIS 合规 | 需配置 | 需配置 | 默认通过 |
| FIPS 支持 | 否 | 否 | 是 |
| 目标场景 | 通用 | 边缘/IoT | 企业/政府 |

## 架构设计

### 控制平面
- **Static Pods**: 控制平面组件以静态 Pod 运行（由 kubelet 管理）
- **etcd**: 嵌入式 etcd 集群用于数据存储
- **Containerd**: 内置 containerd 作为 CRI

### 安全加固组件

```yaml
# /etc/rancher/rke2/config.yaml
# CIS 加固配置示例
kube-apiserver-arg:
  - authorization-mode=Node,RBAC
  - audit-log-path=/var/log/rke2/audit.log
  - audit-log-maxage=30
  - audit-log-maxbackup=10
  - audit-log-maxsize=100
  - request-timeout=300s
  - service-account-lookup=true
  - enable-admission-plugins=NodeRestriction,PodSecurityPolicy

kubelet-arg:
  - protect-kernel-defaults=true
  - read-only-port=0
  - streaming-connection-idle-timeout=5m
```

## 安装部署

### Tarball 安装（推荐用于加固环境）

```bash
# 1. 下载离线包
mkdir -p /opt/rke2-artifacts/
curl -LO https://github.com/rancher/rke2/releases/download/v1.32.3+rke2r1/rke2.linux-amd64.tar.gz
curl -LO https://github.com/rancher/rke2/releases/download/v1.32.3+rke2r1/sha256sum-amd64.txt

# 2. 解压安装
tar xzf rke2.linux-amd64.tar.gz -C /opt/rke2-artifacts/
cp /opt/rke2-artifacts/bin/rke2 /usr/local/bin/

# 3. 启用 Systemd 服务
install -m 755 /opt/rke2-artifacts/lib/systemd/system/rke2-server.service /etc/systemd/system/
systemctl enable --now rke2-server
```

### RPM 安装（RHEL/CentOS/Rocky）

```bash
# 配置 YUM 仓库
cat <<EOF > /etc/yum.repos.d/rke2.repo
[rke2-server]
name=RKE2 Server Repository
baseurl=https://rpm.rancher.io/rke2/stable/1.32/centos/9/x86_64
enabled=1
gpgcheck=1
gpgkey=https://rpm.rancher.io/public.key
EOF

# 安装 Server
dnf install rke2-server
systemctl enable --now rke2-server

# 安装 Agent（工作节点）
dnf install rke2-agent
```

### 高可用安装

```yaml
# /etc/rancher/rke2/config.yaml (所有 Server 节点)
server: https://<load-balancer-ip>:9345
token: <cluster-secret>
tls-san:
  - <load-balancer-ip>
  - <load-balancer-hostname>
cni:
  - canal
disable:
  - rke2-ingress-nginx
```

## CIS 加固指南

### 主机级加固

```bash
#!/bin/bash
# CIS 主机加固脚本（RHEL 9）

# 1. 禁用 USB 存储（防止数据泄露）
echo "install usb-storage /bin/true" > /etc/modprobe.d/usb-storage.conf

# 2. 配置审计规则
cat <<EOF > /etc/audit/rules.d/rke2.rules
-w /etc/rancher/rke2/ -p wa -k rke2-config
-w /var/lib/rancher/rke2/ -p wa -k rke2-data
-a always,exit -F arch=b64 -S setuid -S setgid -k privilege_escalation
EOF

# 3. 配置文件权限
chmod 600 /etc/rancher/rke2/rke2.yaml
chmod 700 /var/lib/rancher/rke2/server/tls

# 4. 启用 SELinux
setenforce 1
sed -i 's/SELINUX=permissive/SELINUX=enforcing/' /etc/selinux/config
```

### Kubernetes 级加固

```yaml
# Pod Security Standards (PSS)
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

```yaml
# Network Policy（零信任）
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: production
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-specific
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: secure-app
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8443
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: database
    ports:
    - protocol: TCP
      port: 5432
```

## 合规扫描

### 使用 kube-bench

```bash
# 安装 kube-bench
kubectl apply -f https://raw.githubusercontent.com/aquasecurity/kube-bench/main/job-rke2.yaml

# 查看结果
kubectl logs job/kube-bench
```

### Rancher 合规扫描

```yaml
# CIS Scan CRD
apiVersion: cis.cattle.io/v1
kind: ClusterScan
metadata:
  name: rke2-hardened-scan
spec:
  scanProfileName: rke2-hardened-1.23
  scoreWarning: pass
```

## 密钥管理

### HashiCorp Vault 集成

```yaml
# Vault CSI Provider
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: vault-backend
spec:
  provider: vault
  parameters:
    vaultAddress: https://vault.internal:8200
    roleName: rke2-apps
    objects: |
      - objectName: "db-password"
        secretPath: "secret/data/rke2/app1"
        secretKey: "password"
```

## 审计与日志

```yaml
# 审计策略
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
# 记录所有身份验证失败
- level: Metadata
  verbs: ["create"]
  resources:
  - group: "authentication.k8s.io"
    resources: ["tokenreviews"]
  omitStages:
  - RequestReceived

# 记录敏感资源的所有操作
- level: RequestResponse
  resources:
  - group: ""
    resources: ["secrets", "configmaps"]
  - group: "rbac.authorization.k8s.io"
    resources: ["roles", "rolebindings"]
```

## 2025 更新

- **CIS v1.9**: 支持最新的 CIS Kubernetes Benchmark v1.9
- **Kubernetes 1.33**: 跟进上游 K8s 最新版本
- **etcd v3.6**: 支持 etcd 3.6 性能改进
- **SELinux 优化**: 改进的 SELinux 策略减少误报

## 相关资源

- [RKE2 官方文档](https://docs.rke2.io/)
- [CIS Kubernetes Benchmark](https://www.cisecurity.org/benchmark/kubernetes)
- [NIST 800-53 Controls](https://csrc.nist.gov/publications/detail/sp/800-53/final)
