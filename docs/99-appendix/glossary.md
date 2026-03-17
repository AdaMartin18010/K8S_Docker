# 术语表

> Docker & Kubernetes 核心术语

---

## 容器相关

| 术语 | 英文 | 说明 |
|------|------|------|
| 容器 | Container | 应用及其依赖的打包单元 |
| 镜像 | Image | 容器的只读模板 |
| 镜像层 | Layer | 镜像的层级结构 |
| Dockerfile | Dockerfile | 定义镜像构建步骤的文件 |
| 仓库 | Registry | 存储和分发镜像的服务 |
| 卷 | Volume | 容器持久化存储 |

---

## Kubernetes 相关

| 术语 | 英文 | 说明 |
|------|------|------|
| Pod | Pod | K8s 最小调度单元 |
| 节点 | Node | 工作机器 |
| 集群 | Cluster | 节点集合 |
| 命名空间 | Namespace | 资源隔离逻辑分组 |
| 标签 | Label | 用于标识资源的键值对 |
| 选择器 | Selector | 根据标签筛选资源 |
| 控制器 | Controller | 确保期望状态与实际状态一致 |
| 调度器 | Scheduler | 负责 Pod 到节点的调度 |
| API Server | API Server | K8s 控制平面的入口 |
| etcd | etcd | 分布式键值存储 |

---

## 工作负载相关

| 术语 | 英文 | 说明 |
|------|------|------|
| Deployment | Deployment | 无状态应用管理 |
| StatefulSet | StatefulSet | 有状态应用管理 |
| DaemonSet | DaemonSet | 每个节点运行一个 Pod |
| Job | Job | 一次性任务 |
| CronJob | CronJob | 定时任务 |
| ReplicaSet | ReplicaSet | Pod 副本管理 |

---

## 网络相关

| 术语 | 英文 | 说明 |
|------|------|------|
| Service | Service | Pod 的网络抽象 |
| Ingress | Ingress | 集群入口路由 |
| Endpoint | Endpoint | Service 后端 Pod 列表 |
| DNS | DNS | 服务发现机制 |
| CNI | CNI | 容器网络接口 |
| NetworkPolicy | NetworkPolicy | 网络隔离策略 |

---

## 存储相关

| 术语 | 英文 | 说明 |
|------|------|------|
| PV | PersistentVolume | 持久卷 |
| PVC | PersistentVolumeClaim | 持久卷声明 |
| StorageClass | StorageClass | 存储类 |
| ConfigMap | ConfigMap | 配置数据 |
| Secret | Secret | 敏感数据 |

---

## 安全相关

| 术语 | 英文 | 说明 |
|------|------|------|
| RBAC | RBAC | 基于角色的访问控制 |
| ServiceAccount | ServiceAccount | Pod 身份标识 |
| NetworkPolicy | NetworkPolicy | 网络策略 |
| PodSecurity | PodSecurity | Pod 安全标准 |
| TLS | TLS | 传输层安全 |
| mTLS | mTLS | 双向 TLS |
