# KubeVirt - Kubernetes 虚拟化

## 概述

KubeVirt 是 CNCF Incubating 项目（正在申请 Graduation），允许在 Kubernetes 上运行虚拟机(VM)工作负载，实现容器和虚拟机的统一管理平台。

> **2025 最新状态**: KubeVirt 已有 41+ 生产采用者，进入 CNCF Top 20 活跃项目，26% 的 Kubernetes 用户用于 VMware 迁移场景。

## 核心组件

```
┌─────────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                          │
│  ┌─────────────┐  ┌─────────────┐  ┌───────────────────────┐   │
│  │ virt-api    │  │virt-controller│ │  virt-operator        │   │
│  │ (API扩展)   │  │(VM生命周期)   │ │  (安装/更新)          │   │
│  └─────────────┘  └─────────────┘  └───────────────────────┘   │
│                         │                                      │
│  ┌──────────────────────┴──────────────────────────────────┐   │
│  │                    Node Level                            │   │
│  │  ┌─────────────┐  ┌─────────────────────────────────┐   │   │
│  │  │virt-handler │  │     virt-launcher (Pod)         │   │   │
│  │  │(节点协调)   │  │  ┌───────────────────────────┐  │   │   │
│  │  └─────────────┘  │  │   QEMU/KVM Process        │  │   │   │
│  │                   │  │   ┌───────────────────┐   │  │   │   │
│  │                   │  │   │  VM Instance      │   │  │   │   │
│  │                   │  │   │  (Guest OS)       │   │  │   │   │
│  │                   │  │   └───────────────────┘   │  │   │   │
│  │                   │  └───────────────────────────┘  │   │   │
│  │                   └─────────────────────────────────┘   │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### 组件职责

| 组件 | 功能 | 部署方式 |
|------|------|----------|
| virt-operator | CRD 管理、生命周期、更新 | Deployment |
| virt-controller | VM 调度、创建 VMI、关联 Pod | Deployment |
| virt-api | API 扩展、验证、认证 | Deployment |
| virt-handler | 节点级操作、设备管理 | DaemonSet |
| virt-launcher | 每个 VM 一个 Pod，运行 QEMU | Pod |

## 关键概念

### VirtualMachine (VM)

持久化的虚拟机定义，类似 StatefulSet，支持启停和重启。

```yaml
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: my-vm
spec:
  running: true  # 控制 VM 运行状态
  template:
    spec:
      domain:
        cpu:
          cores: 4
        resources:
          requests:
            memory: 8Gi
        devices:
          disks:
          - name: rootdisk
            disk:
              bus: virtio
      volumes:
      - name: rootdisk
        persistentVolumeClaim:
          claimName: vm-disk-pvc
```

### VirtualMachineInstance (VMI)

正在运行的 VM 实例，类似 Pod。通常由 VM 控制器自动创建。

### DataVolume (CDI)

容器化数据导入器，用于从镜像创建 PVC 作为 VM 磁盘。

## 核心使用场景

### 1. 遗留应用现代化

无需重构代码，将传统 VM 应用迁移到 Kubernetes 平台。

```yaml
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: legacy-app-vm
  annotations:
    # 启用实时迁移
    kubevirt.io/allow-pod-bridge-network-live-migration: "true"
spec:
  running: true
  template:
    metadata:
      labels:
        app: legacy-app
    spec:
      domain:
        cpu:
          cores: 2
          sockets: 1
          threads: 1
        resources:
          requests:
            memory: 4Gi
        devices:
          disks:
          - name: rootdisk
            disk:
              bus: virtio
          - name: cloudinit
            cdrom:
              bus: sata
          interfaces:
          - name: default
            masquerade: {}  # 使用 NAT 模式，支持实时迁移
      networks:
      - name: default
        pod: {}
      volumes:
      - name: rootdisk
        dataVolume:
          name: legacy-app-rootdisk
      - name: cloudinit
        cloudInitNoCloud:
          userData: |
            #cloud-config
            password: changeme
            chpasswd: { expire: False }
            ssh_pwauth: True
```

### 2. 混合工作负载（容器 + VM）

同一集群运行微服务和传统应用。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hybrid-app-service
spec:
  selector:
    app: hybrid-app
  ports:
  - port: 8080
---
# VM 后端
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: vm-backend
  labels:
    app: hybrid-app
    tier: backend
spec:
  running: true
  template:
    spec:
      domain:
        resources:
          requests:
            memory: 2Gi
      networks:
      - name: default
        pod: {}
      volumes:
      - name: rootdisk
        dataVolume:
          name: backend-disk
---
# 容器前端
apiVersion: apps/v1
kind: Deployment
metadata:
  name: container-frontend
  labels:
    app: hybrid-app
    tier: frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hybrid-app
      tier: frontend
  template:
    metadata:
      labels:
        app: hybrid-app
        tier: frontend
    spec:
      containers:
      - name: frontend
        image: myapp/frontend:latest
        env:
        - name: BACKEND_URL
          value: "http://vm-backend:8080"
```

### 3. 多租户 VM 平台（KCP + KubeVirt）

使用 KCP 管理多集群 VM。

```
┌─────────────────────────────────────────────────────────┐
│  KCP Control Plane (逻辑控制平面)                        │
│  ┌─────────────┐  ┌─────────────────────────────────┐  │
│  │  Workspace  │  │  APIExport (kubevirt.io/v1)     │  │
│  │  /tenant-a  │  │  APIBinding (同步到物理集群)    │  │
│  │  /tenant-b  │  │                                 │  │
│  └─────────────┘  └─────────────────────────────────┘  │
└──────────────────────────┬──────────────────────────────┘
                           │ sync
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  Cluster 1   │   │  Cluster 2   │   │  Cluster 3   │
│  (Region A)  │   │  (Region B)  │   │  (Region C)  │
│ ┌──────────┐ │   │ ┌──────────┐ │   │ ┌──────────┐ │
│ │ KubeVirt │ │   │ │ KubeVirt │ │   │ │ KubeVirt │ │
│ │ VMs      │ │   │ │ VMs      │ │   │ │ VMs      │ │
│ └──────────┘ │   │ └──────────┘ │   │ └──────────┘ │
└──────────────┘   └──────────────┘   └──────────────┘
```

## 2025 高级特性

### 1. 微秒级时钟同步（硬件时间戳）

金融交易、实时分析等延迟敏感场景：

```yaml
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: low-latency-vm
spec:
  template:
    spec:
      domain:
        cpu:
          model: host-passthrough
        devices:
          hostDevices:
          # SR-IOV 网卡直通
          - deviceName: intel.com/sriov_net
            name: sriov-nic
          # PTP 硬件时间戳设备
          - deviceName: ptp.io/ptp-device
            name: ptp-hw
```

### 2. 实时迁移配置

```yaml
apiVersion: kubevirt.io/v1
kind: VirtualMachineInstanceMigration
metadata:
  name: migrate-vm-to-node-b
spec:
  vmiName: my-vm
  # 可选：指定目标节点
  targetNode: node-b
```

### 3. PSI 感知的负载均衡

结合 Kubernetes Descheduler 实现基于 Pressure Stall Information 的自动迁移。

## 生产挑战与解决方案

| 挑战 | 影响 | 解决方案 |
|------|------|----------|
| 持久化存储 | 45% 用户困难 | 使用支持快照的 CSI（如 Portworx、Rook-Ceph） |
| VM 格式转换 | 43% 手动工作 | 使用 virt-v2v、Ansible 自动化 |
| 网络 CNI 冲突 | 实时迁移失败 | 使用 `masquerade` 或 Multus 辅助网络 |
| 文化阻力 | 38% 内部阻力 | 渐进式迁移、培训、混合架构 |

## GitOps 集成（ArgoCD）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: vm-gitops
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/org/vm-manifests.git
    targetRevision: main
    path: vms/
  destination:
    server: https://kubernetes.default.svc
    namespace: vms
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

## 与存储集成

```yaml
apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  name: import-vm-image
spec:
  source:
    http:
      url: "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
  pvc:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 20Gi
    storageClassName: portworx-vm-storage
```

## 监控与告警

```yaml
# Prometheus 规则
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: kubevirt-rules
spec:
  groups:
  - name: kubevirt
    rules:
    - alert: VMCpuHigh
      expr: |
        kubevirt_vmi_vcpu_seconds > 0.8
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "VM CPU 使用率过高"
    - alert: VMMemoryPressure
      expr: |
        kubevirt_vmi_memory_available_bytes /
        kubevirt_vmi_memory_domain_bytes < 0.1
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "VM 内存压力过高"
```

## 总结

KubeVirt 让 Kubernetes 成为真正的统一平台，支持容器和 VM 工作负载。对于正在 VMware 迁移或需要运行不可容器化工作负载的组织，KubeVirt 提供了云原生的虚拟化解决方案。

> **2025 趋势**: KubeVirt 正申请 CNCF Graduation，预计将成为企业级虚拟化的标准选择，特别是在混合云和边缘计算场景。
