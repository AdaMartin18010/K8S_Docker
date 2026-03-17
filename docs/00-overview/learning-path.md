# 学习路径指南

> 基于实际文档结构的个性化学习路径

---

## 🎓 角色一：后端开发工程师

**目标**: 能够使用 Docker 和 K8s 部署自己的应用

### 第一阶段：Docker 基础 (2 周)

**学习内容**:

1. [容器基础概念](../01-fundamentals/container-overview.md)
2. [容器 vs 虚拟机](../01-fundamentals/container-vs-vm.md)
3. [Linux Namespace](../01-fundamentals/linux-namespace.md)
4. [Docker 架构](../02-docker/01-core-concepts/docker-architecture.md)
5. [Dockerfile 最佳实践](../02-docker/02-dockerfile/best-practices.md)

**实践项目**:

- 将自己的应用 Docker 化
- 编写多阶段 Dockerfile
- 使用 Docker Compose 本地开发

### 第二阶段：K8s 基础 (3 周)

**学习内容**:

1. [K8s 架构](../03-kubernetes/01-architecture/)
2. [原生 Sidecar](../03-kubernetes/01-pod/sidecar-native.md)
3. [工作负载](../03-kubernetes/02-workloads/)
4. [网络](../03-kubernetes/03-networking/)
5. [存储](../03-kubernetes/04-storage/)

**实践项目**:

- 在本地集群 (minikube/kind) 部署应用
- 配置健康检查和资源限制
- 实现配置分离 (ConfigMap/Secret)

### 第三阶段：生产实践 (2 周)

**学习内容**:

1. [K8s 安全](../03-kubernetes/05-security/)
2. [运维实践](../03-kubernetes/06-operations/)
3. [CI/CD 指南](../06-practices/cicd-guide.md)

---

## 🛠️ 角色二：DevOps 工程师

**目标**: 构建和维护企业级容器平台

### 第一阶段：深度 Docker (2 周)

**学习内容**:

1. [Docker 网络](../02-docker/01-core-concepts/networking.md)
2. [Docker 存储](../02-docker/01-core-concepts/storage.md)
3. [Dockerfile 安全](../02-docker/02-dockerfile/security.md)
4. [容器安全](../02-docker/04-security/)

### 第二阶段：K8s 深度 (4 周)

**学习内容**:

1. [K8s 架构](../03-kubernetes/01-architecture/) - 深度理解
2. [网络](../03-kubernetes/03-networking/) - CNI、Service、Ingress
3. [存储](../03-kubernetes/04-storage/) - PV、PVC、StorageClass
4. [可观测性](../03-kubernetes/05-observability/)
5. [安全](../03-kubernetes/05-security/)
6. [运维](../03-kubernetes/06-operations/)

### 第三阶段：云原生生态 (3 周)

**学习内容**:

1. [GitOps](../04-ecosystem/gitops/)
2. [可观测性](../05-tools/observability/)
3. [监控工具](../05-tools/prometheus-3/)
4. [设计模式](../05-patterns/)

---

## 🏗️ 角色三：云原生架构师

**目标**: 设计高可用、可扩展的云原生架构

### 第一阶段：理论基础 (2 周)

**学习内容**:

1. [OCI 标准](../01-fundamentals/oci-standard.md)
2. [容器运行时](../01-fundamentals/containerd-runtimes.md)
3. [K8s 架构](../03-kubernetes/01-architecture/)
4. [K8s 1.33 新特性](../03-kubernetes/whats-new-1.33.md)

### 第二阶段：设计模式 (3 周)

**学习内容**:

1. [Sidecar 模式](../05-patterns/sidecar-pattern.md)
2. [微服务模式](../05-patterns/microservices.md)
3. [部署策略](../05-patterns/deployment-strategies.md)
4. [混沌工程](../06-practices/chaos-engineering/)

### 第三阶段：完整生态 (4 周)

**学习内容**:

1. [服务网格](../04-ecosystem/service-mesh/)
2. [eBPF/Cilium](../04-ecosystem/ebpf-cilium/)
3. [GitOps](../04-ecosystem/gitops/)
4. [Dapr](../04-ecosystem/dapr/)
5. [Knative](../04-ecosystem/knative/)
6. [Crossplane](../04-ecosystem/crossplane/)

---

## ⏱️ 时间投入建议

| 角色 | 总时间 | 每日投入 | 完成周期 |
|------|--------|----------|----------|
| 后端开发 | 70h | 2h/天 | 5-6 周 |
| DevOps | 120h | 3h/天 | 8-9 周 |
| 架构师 | 150h | 3h/天 | 10-12 周 |

---

## 📊 技能评估

### 初级 (入门)

- [ ] 理解容器概念
- [ ] 能够编写 Dockerfile
- [ ] 能够部署简单 Pod
- [ ] 理解 Service 概念

### 中级 (熟练)

- [ ] 掌握多阶段构建
- [ ] 能够设计 K8s 应用架构
- [ ] 理解网络策略
- [ ] 能够排查常见问题

### 高级 (专家)

- [ ] 能够设计高可用架构
- [ ] 掌握服务网格
- [ ] 能够进行性能调优
- [ ] 理解底层原理
