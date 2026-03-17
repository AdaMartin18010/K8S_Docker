# 学习路径指南

本文档提供不同角色的个性化学习路径。

---

## 🎓 角色一：后端开发工程师

**目标**: 能够使用 Docker 和 K8s 部署自己的应用

### 第一阶段：Docker 基础 (2 周)

**学习内容**:
1. [01-fundamentals/container-basics.md](../01-fundamentals/container-basics.md)
2. [02-docker/01-core-concepts/](../02-docker/01-core-concepts/)
3. [02-docker/02-dockerfile/best-practices.md](../02-docker/02-dockerfile/best-practices.md)

**实践项目**:
- 将自己的应用 Docker 化
- 编写多阶段 Dockerfile
- 使用 Docker Compose 本地开发

### 第二阶段：K8s 基础 (3 周)

**学习内容**:
1. [03-kubernetes/01-architecture/](../03-kubernetes/01-architecture/)
2. [03-kubernetes/02-workloads/pods.md](../03-kubernetes/02-workloads/pods.md)
3. [03-kubernetes/02-workloads/deployments.md](../03-kubernetes/02-workloads/deployments.md)
4. [03-kubernetes/03-networking/services.md](../03-kubernetes/03-networking/services.md)

**实践项目**:
- 在本地集群 (minikube/kind) 部署应用
- 配置健康检查和资源限制
- 实现配置分离 (ConfigMap/Secret)

### 第三阶段：生产实践 (2 周)

**学习内容**:
1. [03-kubernetes/05-security/pod-security.md](../03-kubernetes/05-security/pod-security.md)
2. [06-practices/cicd-guide.md](../06-practices/cicd-guide.md)

---

## 🛠️ 角色二：DevOps 工程师

**目标**: 构建和维护企业级容器平台

### 第一阶段：深度 Docker (2 周)

**学习内容**:
1. [02-docker/04-security/](../02-docker/04-security/)
2. [02-docker/02-dockerfile/multi-stage.md](../02-docker/02-dockerfile/multi-stage.md)

### 第二阶段：K8s 深度 (4 周)

**学习内容**:
1. [03-kubernetes/01-architecture/](../03-kubernetes/01-architecture/) - 深度理解
2. [03-kubernetes/03-networking/](../03-kubernetes/03-networking/)
3. [03-kubernetes/04-storage/](../03-kubernetes/04-storage/)
4. [03-kubernetes/05-security/](../03-kubernetes/05-security/)
5. [03-kubernetes/06-operations/](../03-kubernetes/06-operations/)

### 第三阶段：云原生生态 (3 周)

**学习内容**:
1. [04-cloud-native/02-gitops/](../04-cloud-native/02-gitops/)
2. [04-cloud-native/03-observability/](../04-cloud-native/03-observability/)
3. [05-patterns/](../05-patterns/)

---

## 🏗️ 角色三：云原生架构师

**目标**: 设计高可用、可扩展的云原生架构

### 第一阶段：理论基础 (2 周)

**学习内容**:
1. [01-fundamentals/](../01-fundamentals/)
2. [03-kubernetes/01-architecture/](../03-kubernetes/01-architecture/)

### 第二阶段：设计模式 (3 周)

**学习内容**:
1. [05-patterns/microservices.md](../05-patterns/microservices.md)
2. [05-patterns/deployment-strategies.md](../05-patterns/deployment-strategies.md)
3. [05-patterns/resilience.md](../05-patterns/resilience.md)

### 第三阶段：完整生态 (4 周)

**学习内容**:
1. [04-cloud-native/01-service-mesh/](../04-cloud-native/01-service-mesh/)
2. [04-cloud-native/02-gitops/](../04-cloud-native/02-gitops/)
3. [04-cloud-native/03-observability/](../04-cloud-native/03-observability/)
4. [06-practices/case-studies/](../06-practices/case-studies/)

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
