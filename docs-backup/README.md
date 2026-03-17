# Docker & Kubernetes 学习指南

> **版本**: 2.0
> **更新日期**: 2025年3月
> **适用版本**: Docker 25.x+, Kubernetes 1.30-1.32

本目录包含 Docker 和 Kubernetes 的完整学习资料，理论与实践相结合。

---

## 📚 文档结构

| 文档 | 内容 | 关联代码 |
|------|------|----------|
| [01-docker-guide.md](./01-docker-guide.md) | Docker 完整指南 | `examples/docker/` |
| [02-kubernetes-guide.md](./02-kubernetes-guide.md) | Kubernetes 完整指南 | `examples/kubernetes/` |
| [03-go-examples.md](./03-go-examples.md) | Go 代码实战指南 | `examples/go-client/` |
| [04-best-practices.md](./04-best-practices.md) | 生产最佳实践 | `examples/anti-patterns/` |
| [05-cheatsheets.md](./05-cheatsheets.md) | 速查表 | - |

---

## 🚀 快速开始

### 学习路径

```
阶段1: Docker 基础
  └─ 阅读 01-docker-guide.md
  └─ 实践 examples/docker/basic/

阶段2: Kubernetes 基础
  └─ 阅读 02-kubernetes-guide.md (1-3章)
  └─ 实践 examples/kubernetes/01-basic-resources/

阶段3: 部署模式
  └─ 阅读 02-kubernetes-guide.md (第4章)
  └─ 实践 examples/kubernetes/02-deployment-patterns/

阶段4: 存储与安全
  └─ 阅读 02-kubernetes-guide.md (5-6章)
  └─ 实践 examples/kubernetes/05-storage/ 和 06-security/

阶段5: Go 开发
  └─ 阅读 03-go-examples.md
  └─ 实践 examples/go-client/

阶段6: 生产实践
  └─ 阅读 04-best-practices.md
  └─ 研究 examples/anti-patterns/
```

---

## 📖 内容说明

### 理论文档 (docs/)

- 精简实用的核心概念
- 2025年最新技术更新
- 与代码示例直接关联

### 代码示例 (examples/)

- 可直接运行的完整代码
- 最佳实践与反例对比
- 生产级配置模板

---

## ⚠️ 技术更新说明

| 旧技术 | 状态 | 替代方案 |
|--------|------|----------|
| PodSecurityPolicy | ❌ 已移除 (1.25+) | Pod Security Standards |
| Docker Shim | ❌ 已移除 (1.24+) | 直接使用 containerd |
| CephFS/RBD 内嵌插件 | ❌ 已移除 (1.31+) | CSI 驱动 |
| `failure-domain.beta.kubernetes.io` | ❌ 已废弃 | `topology.kubernetes.io` |

---

## 🔗 外部资源

- [Docker 官方文档](https://docs.docker.com/)
- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [CNCF 云原生全景](https://landscape.cncf.io/)
- [examples/ 代码库](../examples/)
