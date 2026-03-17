# K8S_Docker - 云原生技术知识体系

> 基于主题的 Docker & Kubernetes 知识库 (2025 Edition)

---

## 📚 项目简介

本项目是一个全面、系统的 Docker 和 Kubernetes 知识库，涵盖从基础概念到生产实践的全方位内容。

- **版本对齐**: Docker 25.x+ / Kubernetes 1.30-1.32
- **技术栈**: 采用 BuildKit 1.0、Gateway API、原生 Sidecar 等最新特性
- **实例语言**: Go

---

## 📁 项目结构

```
K8S_Docker/
├── docs/                    # 📖 主题式文档 (3.02MB, 39+ 文件)
│   ├── 00-overview/         # 项目概览与学习指南
│   ├── 01-fundamentals/     # 容器基础理论
│   ├── 02-docker/           # Docker 容器技术
│   ├── 03-kubernetes/       # Kubernetes 核心
│   ├── 04-ecosystem/        # 云原生生态
│   ├── 05-tools/            # 工具链
│   ├── 06-practices/        # 工程实践
│   └── 99-appendix/         # 附录
├── examples/                # 💻 代码示例 (22KB, 69+ 文件)
│   ├── docker/              # Dockerfile 最佳实践
│   ├── kubernetes/          # K8s 资源清单
│   ├── go-client/           # client-go 示例
│   ├── helm/                # Helm Charts
│   ├── ci-cd/               # CI/CD 流水线
│   └── patterns/            # 设计模式
└── docs-backup/             # 📦 归档的原始文档
```

---

## 🗺️ 文档导航

### 快速入门
- [项目概览](docs/00-overview/README.md)
- [学习路线图](docs/00-overview/roadmap.md)
- [阅读指南](docs/00-overview/guide.md)

### 核心内容

| 主题 | 文档 | 示例 |
|------|------|------|
| **容器基础** | [docs/01-fundamentals/](docs/01-fundamentals/) | - |
| **Docker** | [docs/02-docker/](docs/02-docker/) | [examples/docker/](examples/docker/) |
| **Kubernetes** | [docs/03-kubernetes/](docs/03-kubernetes/) | [examples/kubernetes/](examples/kubernetes/) |
| **云原生生态** | [docs/04-ecosystem/](docs/04-ecosystem/) | - |
| **工程实践** | [docs/06-practices/](docs/06-practices/) | [examples/patterns/](examples/patterns/) |

---

## ✨ 2025 关键更新

### 移除 (已废弃)
- ❌ PodSecurityPolicy (K8s 1.25 移除)
- ❌ Dockershim (K8s 1.24 移除)
- ❌ CephFS/RBD 内置插件

### 新增 (重点特性)
- ✅ BuildKit 1.0 (Docker 默认构建器)
- ✅ Docker Scout (安全扫描)
- ✅ Gateway API v1.1+
- ✅ 原生 Sidecar 容器 (v1.29+ GA)
- ✅ 原地 Pod 垂直扩缩容

---

## 🚀 快速开始

### 查看示例

```bash
# Docker 多阶段构建
cat examples/docker/multi-stage/Dockerfile

# Kubernetes Deployment
kubectl apply -f examples/kubernetes/02-workloads/deployment.yaml

# 运行 Go 示例
cd examples/go-client/01-basic-ops
go run main.go
```

---

## 📊 内容统计

| 类别 | 数量 | 大小 |
|------|------|------|
| 文档文件 | 39+ | 3.02 MB |
| 代码示例 | 69+ | 22 KB |
| 主题分类 | 8 大主题 | 21 子目录 |

---

## 📖 推荐阅读顺序

```
1. 基础篇 (01-fundamentals)
   ├── 容器技术概述
   ├── OCI 开放容器标准
   └── 容器 vs 虚拟机

2. Docker 篇 (02-docker)
   ├── Docker 核心概念
   ├── Dockerfile 最佳实践
   ├── 镜像与仓库
   └── Docker Compose

3. Kubernetes 篇 (03-kubernetes)
   ├── 核心概念
   ├── 工作负载管理
   ├── 网络与服务
   ├── 存储
   └── 运维实践

4. 进阶篇 (04-ecosystem + 06-practices)
   ├── 服务网格
   ├── GitOps
   ├── 设计模式
   └── 案例研究
```

---

## 🔗 相关资源

- [命令速查表](docs/99-appendix/commands.md)
- [资源推荐](docs/99-appendix/resources.md)
- [K8s 官方文档](https://kubernetes.io/docs/)
- [Docker 官方文档](https://docs.docker.com/)

---

## 📝 许可证

[LICENSE](LICENSE)
