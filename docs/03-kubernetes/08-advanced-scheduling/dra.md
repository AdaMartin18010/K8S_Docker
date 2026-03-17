# 动态资源分配 (DRA)

> Kubernetes 1.32+ 资源调度革命

---

## 什么是 DRA？

Dynamic Resource Allocation (DRA) 是 Kubernetes 1.32 中 GA 的新资源调度机制，允许 Pod 请求除 CPU/内存外的自定义资源（GPU、FPGA、RDMA 等）。

```
┌─────────────────────────────────────────────────────────────┐
│              传统资源调度 vs DRA                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  传统方式 (Device Plugin)          DRA 方式                  │
│  ┌──────────────────────┐        ┌──────────────────────┐  │
│  │ Pod 请求 GPU         │        │ Pod 申请 ResourceClaim│  │
│  │ resources:           │        │ ┌──────────────────┐ │  │
│  │   nvidia.com/gpu: 2  │   →    │ │ gpu-claim        │ │  │
│  └──────────┬───────────┘        │ │ - 2x A100        │ │  │
│             │                    │ │ - NVLink 连接    │ │  │
│             ↓                    │ │ - 共享/独占      │ │  │
│  ┌──────────────────────┐        │ └──────────────────┘ │  │
│  │ Device Plugin        │        └──────────┬───────────┘  │
│  │ 简单计数分配         │                   │              │
│  │ 无法表达复杂需求     │                   ↓              │
│  └──────────────────────┘        ┌──────────────────────┐  │
│                                  │ DRA Driver           │  │
│                                  │ 结构化参数匹配       │  │
│                                  │ 智能调度决策         │  │
│                                  └──────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心概念

| 资源 | 说明 |
|------|------|
| **ResourceClaim** | 资源声明，定义应用需要什么资源 |
| **ResourceClaimTemplate** | 模板，用于自动创建 ResourceClaim |
| **ResourceClass** | 资源类，定义资源的类型和参数 |
| **Device** | 设备，实际的硬件资源 |
| **AllocationResult** | 分配结果，记录资源如何被分配 |

---

## ResourceClaim 示例

```yaml
apiVersion: resource.k8s.io/v1beta1
kind: ResourceClaim
metadata:
  name: gpu-claim
spec:
  resourceClassName: nvidia-gpu
  parameters:
    apiVersion: gpu.resource.nvidia.com/v1alpha1
    kind: GpuClaimParameters
    count: 2
    selector:
      productName: "NVIDIA A100"
    sharing:
      strategy: TimeSlicing
```

---

## Pod 使用 DRA

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ml-training
spec:
  containers:
    - name: trainer
      image: pytorch:latest
      resources:
        claims:
          - name: gpu
  resourceClaims:
    - name: gpu
      source:
        resourceClaimTemplateName: gpu-claim-template
```

---

## 结构化参数 (Structured Parameters)

K8s 1.32 核心改进，kube-scheduler 可以直接模拟资源分配，无需第三方驱动。

```
优势:
1. 调度器可以预测资源可用性
2. 支持 Cluster Autoscaler 自动扩缩容
3. 无需调度插件即可进行资源分配
```

---

## DRA vs Device Plugin

| 特性 | Device Plugin | DRA |
|------|---------------|-----|
| **资源表达** | 简单计数 | **结构化参数** |
| **调度决策** | 节点本地 | **全局优化** |
| **共享策略** | 不支持 | **支持时分复用** |
| **自动扩缩容** | 不支持 | **支持** |
| **资源发现** | 静态 | **动态** |

---

## 2025 DRA 生态

| 厂商 | DRA 支持 |
|------|----------|
| **NVIDIA** | gpu-driver v0.9.0+ |
| **Intel** | intel-device-plugins |
| **AMD** | amdgpu-driver |
| **华为** | Ascend DRA driver |

---

## 迁移指南

```bash
# 1. 启用 DRA 特性门
--feature-gates=DynamicResourceAllocation=true

# 2. 安装 DRA 驱动
kubectl apply -f nvidia-dra-driver.yaml

# 3. 创建 ResourceClass
kubectl apply -f resource-class.yaml

# 4. 迁移 Pod 配置
# 从: resources.limits.nvidia.com/gpu: 2
# 到: resourceClaims 配置
```
