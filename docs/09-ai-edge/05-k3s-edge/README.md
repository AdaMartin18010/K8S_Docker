# K3s - 轻量级 Kubernetes 发行版

## 概述

K3s 是 Rancher Labs（现 SUSE）开发的 CNCF 认证轻量级 Kubernetes 发行版，专为边缘计算、IoT、CI/CD 和资源受限环境设计。

## 核心特性

### 轻量级架构
- **单二进制文件**: 所有组件打包在 <100MB 的二进制文件中
- **低资源消耗**: 最低 512MB RAM 和 1 CPU 核心即可运行
- **简化存储**: 默认使用 SQLite 替代 etcd（支持 etcd、MySQL、PostgreSQL）
- **内置组件**: Containerd、Flannel、CoreDNS、Traefik、ServiceLB、Metrics Server

### 适用场景
- 边缘计算 (Edge Computing)
- IoT 设备与传感器
- ARM 架构设备（Raspberry Pi）
- CI/CD 流水线
- 开发测试环境
- 离线/隔离环境

## 架构对比

| 特性 | K8s (标准) | K3s |
|------|-----------|-----|
| 部署方式 | 多组件独立运行 | 单二进制单进程 |
| 存储后端 | 仅 etcd | SQLite/PostgreSQL/MySQL/etcd |
| 安装时间 | 数小时 | 数分钟 |
| 资源需求 | 2GB+ RAM | 512MB+ RAM |
| 外部依赖 | 较多 | 极少 |

## 安装部署

### 快速安装

```bash
# 标准安装
curl -sfL https://get.k3s.io | sh -

# 验证安装
sudo k3s kubectl get nodes

# 配置 kubeconfig
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
```

### 高可用模式

```bash
# Server 节点 1（初始化）
curl -sfL https://get.k3s.io | K3S_TOKEN=<secret> sh -s - server \
  --cluster-init \
  --tls-san=<load-balancer-ip>

# Server 节点 2/3（加入）
curl -sfL https://get.k3s.io | K3S_TOKEN=<secret> sh -s - server \
  --server https://<server1-ip>:6443 \
  --tls-san=<load-balancer-ip>

# Agent 节点
curl -sfL https://get.k3s.io | K3S_TOKEN=<secret> K3S_URL=https://<server-ip>:6443 sh -
```

### 外部数据库模式

```bash
# 使用 PostgreSQL
curl -sfL https://get.k3s.io | sh -s - server \
  --datastore-endpoint="postgres://user:pass@postgres:5432/k3s"

# 使用 MySQL
curl -sfL https://get.k3s.io | sh -s - server \
  --datastore-endpoint="mysql://user:pass@tcp(mysql:3306)/k3s"
```

## 边缘计算最佳实践

### 离线部署

```bash
# 1. 准备离线镜像
k3s air-gap 模式需要：
# - k3s 二进制文件
# - airgap-images.tar（镜像包）

# 2. 导入镜像
sudo mkdir -p /var/lib/rancher/k3s/agent/images/
sudo cp airgap-images.tar /var/lib/rancher/k3s/agent/images/

# 3. 离线安装
sudo INSTALL_K3S_SKIP_DOWNLOAD=true ./install.sh
```

### 自动升级

```yaml
# system-upgrade-controller 配置
apiVersion: upgrade.cattle.io/v1
kind: Plan
metadata:
  name: server-plan
  namespace: system-upgrade
spec:
  concurrency: 1
  cordon: true
  nodeSelector:
    matchExpressions:
    - key: node-role.kubernetes.io/control-plane
      operator: Exists
  serviceAccountName: system-upgrade
  upgrade:
    image: rancher/k3s-upgrade
  version: v1.32.3+k3s1
```

### 边缘节点管理

```yaml
# 边缘节点标签和污点
apiVersion: v1
kind: Node
metadata:
  labels:
    node-type: edge
    location: factory-floor-1
spec:
  taints:
  - key: dedicated
    value: edge
    effect: NoSchedule
```

```yaml
# 边缘应用部署
apiVersion: apps/v1
kind: Deployment
metadata:
  name: edge-analytics
spec:
  replicas: 1
  selector:
    matchLabels:
      app: edge-analytics
  template:
    spec:
      nodeSelector:
        node-type: edge
      tolerations:
      - key: dedicated
        operator: Equal
        value: edge
        effect: NoSchedule
      containers:
      - name: analyzer
        image: edge-analyzer:v1.0
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
```

## 安全加固

### Secrets 加密

```bash
# 启用 secrets 加密
k3s secrets-encrypt status
k3s secrets-encrypt rotate
k3s secrets-encrypt reencrypt
```

### 网络策略

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

### Pod 安全标准

```bash
# 启用 Pod Security Admission
# /etc/rancher/k3s/registries.yaml
kube-apiserver-arg:
  - admission-control-config-file=/etc/rancher/k3s/psa.yaml
```

## 监控与可观测性

```yaml
# 边缘友好的监控配置
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
data:
  prometheus.yml: |
    global:
      scrape_interval: 30s  # 降低频率节省带宽
      evaluation_interval: 30s
    remote_write:
    - url: "https://central-prometheus/api/v1/write"
      queue_config:
        max_samples_per_send: 100
        max_shards: 2
```

## 2025 新特性

- **Spegel 镜像分发**: 分布式镜像 Registry 加速边缘部署
- **嵌入式 etcd**: 单节点默认使用嵌入式 etcd 替代 SQLite
- **Containerd 2.0**: 支持 OCI 镜像规范 v1.1
- **IPv6 双栈**: 完整的 IPv4/IPv6 双栈支持
- **Traefik v3**: 默认 Ingress 控制器升级

## 相关资源

- [K3s 官方文档](https://docs.k3s.io/)
- [K3s GitHub](https://github.com/k3s-io/k3s)
- [System Upgrade Controller](https://github.com/rancher/system-upgrade-controller)
