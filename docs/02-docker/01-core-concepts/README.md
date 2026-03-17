# Docker 核心概念

> Docker 25.x 完整指南

---

## 本章内容

1. [Docker 架构](./docker-architecture.md)
2. [镜像与容器](./images-and-containers.md)
3. [Docker 网络](./docker-networking.md)
4. [Docker 存储](./docker-storage.md)
5. [Docker API](./docker-api.md)

---

## Docker 架构演进

### Docker 20.10 之前

```
Docker Client → Docker Daemon → containerd → runc
```

### Docker 25.x (当前)

```
Docker Client → Docker Daemon (+ BuildKit) → containerd → runc
                     │
                     └→ Docker Scout (安全扫描)
                     └→ Docker Compose (编排)
```

---

## 核心组件

| 组件 | 功能 | 2025年更新 |
|------|------|-----------|
| Docker CLI | 命令行工具 | 支持 Compose 插件 |
| Docker Daemon | 守护进程 | 默认启用 BuildKit |
| containerd | 容器运行时 | 行业标准实现 |
| runc | 低层运行时 | OCI 标准实现 |
| BuildKit | 构建系统 | 1.0 GA，性能提升 40% |
| Docker Scout | 安全扫描 | 新增漏洞分析 |

---

## 快速开始

```bash
# 检查版本
docker version
docker info

# 运行第一个容器
docker run hello-world

# 运行 Nginx
docker run -d -p 80:80 --name my-nginx nginx:alpine

# 查看容器
docker ps

# 停止和删除
docker stop my-nginx
docker rm my-nginx
```
