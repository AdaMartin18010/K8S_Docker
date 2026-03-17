# 边缘计算

> Kubernetes 在边缘场景的部署

---

## 边缘计算挑战

```
┌─────────────────────────────────────────────────────────────┐
│                  边缘计算 vs 云计算                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  特性            云计算              边缘计算                 │
│  ────────────────────────────────────────────────────────   │
│  网络         高带宽、低延迟         不稳定、高延迟           │
│  资源         充足                   受限 (CPU/内存/存储)     │
│  可靠性       高可用                 间歇性连接               │
│  安全边界     数据中心               物理暴露                 │
│  规模         1000s 节点             10000s 边缘节点          │
│  运维         专业团队               无人值守                 │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## KubeEdge

华为开源的边缘计算框架，CNCF 孵化项目。

```
┌─────────────────────────────────────────────────────────────┐
│                     KubeEdge 架构                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Cloud Side                  Edge Side                      │
│  ────────────────────────────────────────────────────────   │
│                                                              │
│  ┌─────────────────┐         ┌─────────────────┐           │
│  │ CloudCore       │◄────────►│ EdgeCore        │           │
│  │ • EdgeController│  Sync    │ • Edged         │           │
│  │ • DeviceController│  (WebSocket)│ • MetaManager   │           │
│  │ • SyncController│         │ • EventBus      │           │
│  └─────────────────┘         │ • ServiceBus    │           │
│         │                    └─────────────────┘           │
│         │                           │                      │
│    ┌────┴────┐                 ┌────┴────┐                │
│    ↓         ↓                 ↓         ↓                │
│  ┌──────┐  ┌──────┐         ┌──────┐  ┌──────┐           │
│  │K8s API│  │Device│         │Pod   │  │Device│           │
│  │Server │  │Twin  │         │(Edge)│  │Twin  │           │
│  └──────┘  └──────┘         └──────┘  └──────┘           │
│                                                              │
│  优势: 边缘离线自治、轻量、设备孪生                           │
└─────────────────────────────────────────────────────────────┘
```

### 部署 EdgeCore

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: edgecore-config
data:
  edgecore.yaml: |
    modules:
      edged:
        nodeIP: 192.168.1.100
        runtimeType: containerd
      edgeHub:
        websocket:
          server: cloudcore.kubeedge:10000
```

---

## SuperEdge

腾讯开源的边缘容器方案。

```yaml
apiVersion: superedge.io/v1
kind: NodeUnit
metadata:
  name: beijing-edge
spec:
  nodes:
    - edge-node-1
    - edge-node-2
  selector:
    matchLabels:
      region: beijing
---
apiVersion: superedge.io/v1
kind: ServiceGrid
metadata:
  name: edge-service
spec:
  gridUniqKey: region
  template:
    selector:
      app: nginx
    ports:
      - port: 80
```

---

## 边缘自治

当边缘与云端断开连接时，边缘节点需要继续运行。

```yaml
apiVersion: v1
kind: Node
metadata:
  annotations:
    node.beta.kubernetes.io/authority: "true"
spec:
  # 边缘节点配置
status:
  conditions:
    - type: EdgeConnection
      status: "False"  # 断开连接
      reason: EdgeDisconnected
      # 但 Pod 继续运行
```

---

## 设备管理

```yaml
apiVersion: devices.kubeedge.io/v1alpha2
kind: Device
metadata:
  name: temperature-sensor
  labels:
    model: sensor
spec:
  deviceModelRef:
    name: temperature-model
  nodeSelector:
    nodeSelectorTerms:
      - matchExpressions:
          - key: edge-node
            operator: In
            values: ["edge-001"]
  twins:
    - propertyName: temperature
      desired:
        value: "22.5"
      reported:
        value: "22.3"
```

---

## 边缘方案对比

| 方案 | 重量 | 特点 | 适用场景 |
|------|------|------|----------|
| **KubeEdge** | 中 | 云边协同、设备孪生 | 物联网、工业 |
| **SuperEdge** | 中 | 国内优化、运维友好 |  CDN、视频 |
| **OpenYurt** | 轻 | 阿里云、原生扩展 | 大规模边缘 |
| **k3s** | 极轻 | 单二进制、完整K8s | 资源受限场景 |
| **MicroK8s** | 轻 | Ubuntu、Snap | 开发测试 |

---

## k3s - 轻量级 K8s

```bash
# 单命令安装
curl -sfL https://get.k3s.io | sh -

# 单节点即可运行
k3s kubectl get nodes

# 资源占用: < 500MB 内存
```

### 适用边缘场景

- 工厂网关
- 零售门店
- 自动驾驶车辆
- 无人机
