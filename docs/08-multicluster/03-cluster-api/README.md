# Cluster API (CAPI) - Kubernetes 集群生命周期管理

## 概述

Cluster API (CAPI) 是 Kubernetes 的子项目，提供声明式 API 和工具来管理 Kubernetes 集群的生命周期（创建、配置、升级、删除）。2025 年，CAPI v1beta2 版本带来更稳定的 API 和更强大的多集群管理能力。

> **关键数据**: CAPI 是 VKS (vSphere Kubernetes Service)、EKS-A、GKE On-Prem 等产品的底层技术。

## 架构原理

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Management Cluster (管理集群)                            │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Cluster API Controllers                          │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│  │  │   Cluster    │  │   Machine    │  │  Kubeconfig  │              │   │
│  │  │  Controller  │  │  Controller  │  │  Controller  │              │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Infrastructure Provider                          │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │   │
│  │  │  CAPA    │  │  CAPG    │  │  CAPZ    │  │  CAPV    │            │   │
│  │  │ (AWS)    │  │ (GCP)    │  │ (Azure)  │  │(vSphere) │            │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
┌─────────────▼────────┐ ┌────────▼────────┐ ┌────────▼────────┐
│   Workload Cluster 1 │ │ Workload Cluster│ │ Workload Cluster│
│     (Production)     │ │   (Staging)     │ │   (Dev/Test)    │
└──────────────────────┘ └─────────────────┘ └─────────────────┘
```

## 安装与使用

### 安装 clusterctl

```bash
curl -L https://github.com/kubernetes-sigs/cluster-api/releases/latest/download/clusterctl-linux-amd64 -o clusterctl
chmod +x clusterctl
sudo mv clusterctl /usr/local/bin/
```

### 初始化管理集群

```bash
# AWS Provider
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=<your-access-key>
export AWS_SECRET_ACCESS_KEY=<your-secret-key>
clusterctl init --infrastructure aws
```

### 创建工作负载集群

```bash
# 生成集群配置
clusterctl generate cluster production-cluster \
  --infrastructure aws \
  --kubernetes-version v1.30.0 \
  --control-plane-machine-count 3 \
  --worker-machine-count 5 \
  > production-cluster.yaml

# 应用配置
kubectl apply -f production-cluster.yaml

# 获取 kubeconfig
clusterctl get kubeconfig production-cluster > production.kubeconfig
```

## 核心资源

| 资源 | 作用 |
|------|------|
| Cluster | 集群整体配置 |
| Machine | 单个节点生命周期 |
| MachineDeployment | 节点组管理 |
| KubeadmControlPlane | 控制平面配置 |

## 升级集群

```bash
# 升级 Kubernetes 版本
kubectl patch kubeadmcontrolplane production-cluster-control-plane \
  --type merge \
  --patch '{"spec":{"version":"v1.30.1"}}'
```

## ClusterClass 模板

```yaml
apiVersion: cluster.x-k8s.io/v1beta2
kind: ClusterClass
metadata:
  name: production-class
spec:
  controlPlane:
    ref:
      apiVersion: controlplane.cluster.x-k8s.io/v1beta2
      kind: KubeadmControlPlaneTemplate
      name: production-control-plane
  workers:
    machineDeployments:
    - class: default-worker
      template:
        infrastructure:
          ref:
            apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
            kind: AWSMachineTemplate
            name: production-worker
```

## 总结

Cluster API 实现 Kubernetes 集群的声明式管理，是平台工程和多云策略的核心组件。
