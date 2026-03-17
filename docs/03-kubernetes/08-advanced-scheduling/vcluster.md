# vCluster - 虚拟 Kubernetes 集群

> 轻量级多租户解决方案

---

## 什么是 vCluster？

vCluster 在共享宿主集群上创建功能完整的虚拟 Kubernetes 集群，每个虚拟集群拥有独立的控制平面。

```
┌─────────────────────────────────────────────────────────────┐
│                    vCluster 架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  宿主集群                                                     │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Kube API Server                                      │  │
│  │  ├─ Namespace: vcluster-team-a                        │  │
│  │  │  ├─ Pod: vcluster-0 (虚拟控制平面)                 │  │
│  │  │  ├─ Pod: coredns                                 │  │
│  │  │  ├─ Pod: app-1, app-2 (工作负载)                 │  │
│  │  │                                                  │  │
│  │  ├─ Namespace: vcluster-team-b                        │  │
│  │  │  ├─ Pod: vcluster-0 (独立控制平面)                │  │
│  │  │  ├─ Pod: coredns                                 │  │
│  │  │  └─ Pod: service-1 (工作负载)                    │  │
│  │                                                    │  │
│  └────────────────────────────────────────────────────┘  │
│                                                              │
│  虚拟集群 (Team A 视角)                                       │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  vCluster API Server                                  │  │
│  │  ├─ Namespace: production                             │  │
│  │  │  ├─ Deployment: frontend                          │  │
│  │  │  ├─ Deployment: backend                           │  │
│  │  │  └─ Service: database                             │  │
│  │  ├─ Namespace: development                            │  │
│  │  │  └─ Deployment: test-app                          │  │
│  │                                                    │  │
│  │  CRDs, RBAC, NetworkPolicy 完全隔离                  │  │
│  └────────────────────────────────────────────────────┘  │
│                                                              │
│  优势:                                                       │
│  • 强隔离 (控制平面级别)                                     │
│  • 低成本 (共享工作节点)                                     │
│  • 快速启动 (秒级)                                          │
│  • 完整权限 (管理员权限)                                     │
└─────────────────────────────────────────────────────────────┘
```

---

## vCluster 对比其他方案

| 特性 | Namespace | vCluster | 独立集群 |
|------|-----------|----------|----------|
| **隔离级别** | 弱 | 强 | 最强 |
| **资源开销** | 低 | 中 | 高 |
| **启动时间** | 秒 | 秒 | 分钟 |
| **权限范围** | 受限 | 完整 | 完整 |
| **成本** | 低 | 低 | 高 |
| **运维复杂度** | 低 | 中 | 高 |

---

## 安装和使用

### 安装 CLI

```bash
# macOS/Linux
curl -s -L https://github.com/loft-sh/vcluster/releases/latest | \
  sed -nE 's!.*"([^"]*vcluster-linux-amd64)".*!https://github.com\1!p' | \
  xargs -n 1 curl -L -o vcluster && chmod +x vcluster && sudo mv vcluster /usr/local/bin

# 或使用 Homebrew
brew install vcluster
```

### 创建虚拟集群

```bash
# 创建虚拟集群
vcluster create team-a-cluster --namespace team-a

# 自动切换到虚拟集群上下文
kubectl get namespaces
kubectl create namespace production
kubectl apply -f deployment.yaml

# 切回宿主集群
vcluster disconnect
```

### Helm 安装

```bash
# 添加仓库
helm repo add loft-sh https://charts.loft.sh
helm repo update

# 安装 vCluster
helm install team-b-cluster loft-sh/vcluster \
  --namespace team-b \
  --create-namespace \
  --set storage.persistence=false  # 开发测试环境
```

---

## 配置选项

```yaml
# vcluster.yaml - 高级配置
apiVersion: config.vcluster.loft.sh/v1alpha1
kind: VirtualClusterConfig

# 同步配置
sync:
  toHost:
    endpoints:
      enabled: true
    pods:
      enabled: true
      enforceTolerations:
        - key: vcluster-team-a
          operator: Exists
  fromHost:
    nodes:
      enabled: true
      selector:
        labels:
          workload-type: shared

# 控制平面配置
controlPlane:
  backingStore:
    etcd:
      enabled: true
      persistence:
        enabled: true
        size: 5Gi

  # 启用高级功能
  advanced:
    virtualScheduler:
      enabled: true

    defaultImageRegistry: my-registry.io/vcluster

# 网络配置
networking:
  advanced:
    proxyKubelets:
      byHostname: false
      byIP: false

  # 重用宿主集群的 Coredns
  reuseNamespaceCIDR: true

# 集成
integrations:
  kubeVirt:
    enabled: true
    sync:
      persistentVolumeClaims:
        enabled: true
```

---

## 多租户场景

### 租户隔离

```bash
# 为每个团队创建虚拟集群
for team in team-a team-b team-c; do
  vcluster create ${team}-cluster --namespace ${team} \
    --set "isolation.enabled=true" \
    --set "isolation.networkPolicy.enabled=true" \
    --set "isolation.resourceQuota.enabled=true"
done
```

### CI/CD 隔离

```yaml
# CI Pipeline 中使用
stages:
  - test

variables:
  VCLUSTER_NAME: "ci-$CI_PIPELINE_ID"

test:
  script:
    # 创建临时集群
    - vcluster create $VCLUSTER_NAME --namespace ci-tests

    # 部署应用
    - kubectl apply -f k8s/
    - kubectl rollout status deployment/app

    # 运行测试
    - kubectl run test-runner --image=test-image --rm -i -- ./run-tests.sh

    # 清理
    - vcluster delete $VCLUSTER_NAME --namespace ci-tests
```

---

## 使用场景

| 场景 | 说明 |
|------|------|
| **多租户 SaaS** | 为每个客户提供隔离的 K8s 环境 |
| **开发环境** | 每个开发者独立的集群 |
| **CI/CD 测试** | 并行的临时测试环境 |
| **培训演示** | 安全、可重置的实验环境 |
| **边缘场景** | 轻量级边缘 K8s 环境 |

---

## 最佳实践

1. **资源配额**: 为每个 vCluster 设置资源限制
2. **网络策略**: 启用网络隔离
3. **持久化**: 生产环境启用 etcd 持久化
4. **监控**: 分别监控虚拟集群和宿主集群
5. **备份**: 定期备份虚拟集群数据
