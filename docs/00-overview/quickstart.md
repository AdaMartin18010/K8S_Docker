# 快速开始指南

> 30分钟上手云原生技术栈

---

## 🎯 路径选择

根据你的角色选择合适的学习路径：

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  我是开发者 👨‍💻                    我是运维工程师 👨‍🔧                      │
│  ─────────────────────             ─────────────────────                      │
│  1. Docker 基础 (10分钟)           1. K8s 基础概念 (10分钟)                  │
│  2. K8s 应用部署 (10分钟)          2. kubectl 实操 (10分钟)                  │
│  3. 编写 Dockerfile (10分钟)       3. 排查 Pod 故障 (10分钟)                 │
│                                                                              │
│  [开始开发者路径](#开发者路径)     [开始运维路径](#运维工程师路径)            │
│                                                                              │
│  我是架构师 🏗️                     我是学生/初学者 🎓                        │
│  ─────────────────────             ─────────────────────                      │
│  1. 云原生架构概览 (10分钟)        1. 容器概念理解 (10分钟)                  │
│  2. 技术选型决策 (10分钟)          2. Docker 初体验 (10分钟)                 │
│  3. 查看对比矩阵 (10分钟)          3. K8s 入门实战 (10分钟)                  │
│                                                                              │
│  [开始架构师路径](#架构师路径)     [开始初学者路径](#初学者路径)              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 开发者路径

### Step 1: Docker 基础 (10分钟)

```bash
# 安装 Docker (已安装请跳过)
# macOS: https://docs.docker.com/desktop/install/mac-install/
# Windows: https://docs.docker.com/desktop/install/windows-install/
# Linux: curl -fsSL https://get.docker.com | sh

# 验证安装
docker --version

# 运行第一个容器
docker run hello-world

# 运行 Nginx 并访问
docker run -d -p 8080:80 --name my-nginx nginx
# 浏览器访问 http://localhost:8080

# 查看运行中的容器
docker ps

# 停止并删除容器
docker stop my-nginx
docker rm my-nginx
```

**关键概念**:

- 镜像 (Image): 应用的只读模板
- 容器 (Container): 镜像的运行实例
- 端口映射: 将容器端口映射到主机端口

---

### Step 2: Kubernetes 应用部署 (10分钟)

```bash
# 安装 minikube (本地 K8s 集群)
# macOS: brew install minikube
# Linux: curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
# Windows: https://storage.googleapis.com/minikube/releases/latest/minikube-installer.exe

# 启动集群
minikube start

# 验证
kubectl get nodes

# 部署应用
kubectl create deployment hello-node --image=registry.k8s.io/e2e-test-images/agnhost:2.39 -- /agnhost netexec --http-port=8080

# 暴露服务
kubectl expose deployment hello-node --type=LoadBalancer --port=8080

# 查看服务
kubectl get services

# 在浏览器中访问 (minikube 会自动打开)
minikube service hello-node

# 清理
kubectl delete service hello-node
kubectl delete deployment hello-node
```

**关键概念**:

- Deployment: 管理应用副本
- Service: 暴露应用访问入口
- Pod: 包含一个或多个容器的最小部署单元

---

### Step 3: 编写 Dockerfile (10分钟)

创建一个简单的 Go 应用:

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

```go
// main.go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from Kubernetes!")
    })
    http.ListenAndServe(":8080", nil)
}
```

构建和部署:

```bash
# 构建镜像
docker build -t myapp:v1 .

# 加载到 minikube
minikube image load myapp:v1

# 部署到 K8s
kubectl create deployment myapp --image=myapp:v1
kubectl expose deployment myapp --type=LoadBalancer --port=8080

# 验证
minikube service myapp
```

**下一步**: 学习 [Helm](../05-tools/helm/README.md) 和 [GitOps](../04-ecosystem/gitops/README.md)

---

## 运维工程师路径

### Step 1: K8s 基础概念 (10分钟)

```
核心组件:
├── 控制平面 (Control Plane)
│   ├── API Server: 所有请求的入口
│   ├── etcd: 集群状态存储
│   ├── Scheduler: 调度 Pod 到节点
│   └── Controller Manager: 维护期望状态
│
└── 工作节点 (Worker Node)
    ├── kubelet: 管理 Pod 生命周期
    ├── kube-proxy: 网络代理
    └── Container Runtime: 运行容器
```

```bash
# 查看集群信息
kubectl cluster-info

# 查看节点
kubectl get nodes -o wide

# 查看系统组件
kubectl get pods -n kube-system
```

---

### Step 2: kubectl 实操 (10分钟)

```bash
# 常用命令速记
kubectl get pods                    # 查看 Pod
kubectl get pods -o wide           # 详细信息
kubectl describe pod <name>        # 查看详情
kubectl logs <pod-name>            # 查看日志
kubectl logs <pod-name> -f         # 实时日志
kubectl exec -it <pod-name> -- sh  # 进入容器
kubectl apply -f manifest.yaml     # 应用配置
kubectl delete -f manifest.yaml    # 删除资源

# 资源类型缩写
kubectl get po     # pods
kubectl get svc    # services
kubectl get deploy # deployments
kubectl get ns     # namespaces
kubectl get cm     # configmaps
kubectl get secret # secrets
kubectl get pv     # persistentvolumes
kubectl get pvc    # persistentvolumeclaims

# 上下文切换
kubectl config get-contexts        # 列出上下文
kubectl config use-context <name>  # 切换上下文
kubectl config current-context     # 当前上下文
```

---

### Step 3: 排查 Pod 故障 (10分钟)

```bash
# Pod 处于 Pending 状态
kubectl describe pod <pod-name>
# 常见原因: 资源不足、节点亲和性、PVC 未绑定

# Pod 处于 CrashLoopBackOff
kubectl logs <pod-name> --previous
# 常见原因: 应用启动失败、健康检查失败

# Pod 处于 ImagePullBackOff
kubectl describe pod <pod-name>
# 常见原因: 镜像不存在、镜像拉取权限不足

# Pod 运行但无法访问
kubectl get svc                    # 检查 Service
kubectl endpoints <svc-name>       # 检查 Endpoints
kubectl port-forward pod/<name> 8080:80  # 本地测试

# 网络问题排查
kubectl run -it --rm debug --image=nicolaka/netshoot --restart=Never -- sh
# 在 debug 容器中:
# - ping <service-ip>
# - nslookup <service-name>
# - curl <endpoint>
```

**故障排查流程图**: [故障排查指南](../06-practices/debugging-guide.md)

**下一步**: 学习 [监控告警](../05-tools/observability/README.md) 和 [备份恢复](../06-storage/velero/README.md)

---

## 架构师路径

### Step 1: 云原生架构概览 (10分钟)

阅读 [知识全景图](../99-appendix/knowledge-graphs/mindmap-cloudnative.md) 了解:

- 云原生技术栈全景
- 层次架构关系
- 核心技术演进

### Step 2: 技术选型决策 (10分钟)

使用决策树进行选型:

| 决策点 | 推荐方案 |
|--------|----------|
| CNI 网络 | [Cilium vs Calico 对比](../99-appendix/knowledge-graphs/decision-trees.md) |
| 服务网格 | [Istio Ambient vs Linkerd](../99-appendix/knowledge-graphs/decision-trees.md) |
| GitOps | [ArgoCD vs Flux](../99-appendix/knowledge-graphs/decision-trees.md) |
| 存储 | [Rook vs Longhorn](../99-appendix/knowledge-graphs/decision-trees.md) |

### Step 3: 查看对比矩阵 (10分钟)

详细对比: [技术对比矩阵](../99-appendix/knowledge-graphs/comparison-matrix.md)

**下一步**: 设计你的 [系统架构](../99-appendix/knowledge-graphs/architecture-systems.md)

---

## 初学者路径

### Step 1: 容器概念理解 (10分钟)

阅读 [容器概述](../01-fundamentals/container-overview.md) 理解:

- 什么是容器?
- 容器 vs 虚拟机
- 容器解决了什么问题?

### Step 2: Docker 初体验 (10分钟)

跟随 [Docker 官方教程](https://docs.docker.com/get-started/) 完成:

1. 运行你的第一个容器
2. 构建你的第一个镜像
3. 使用 Docker Compose

### Step 3: K8s 入门实战 (10分钟)

完成 [K8s 官方教程](https://kubernetes.io/docs/tutorials/kubernetes-basics/):

1. 创建集群
2. 部署应用
3. 暴露应用
4. 扩缩容
5. 更新应用

**下一步**: [K8s 管理员学习路径](../99-appendix/knowledge-graphs/learning-pathways.md)

---

## 📚 继续深入学习

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         进阶学习资源                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  📖 必读文档                                                                 │
│  ├── [Docker 最佳实践](../02-docker/02-dockerfile/best-practices.md)        │
│  ├── [K8s 核心概念](../03-kubernetes/01-architecture/README.md)             │
│  ├── [服务网格对比](../04-ecosystem/service-mesh/README.md)                 │
│  └── [可观测性指南](../05-tools/observability/opentelemetry.md)             │
│                                                                              │
│  💻 动手实验                                                                 │
│  ├── [部署微服务应用](../examples/microservices-demo/)                       │
│  ├── [配置 CI/CD 流水线](../examples/ci-cd/)                                 │
│  └── [混沌工程实验](../examples/chaos-mesh/)                                 │
│                                                                              │
│  🎓 认证准备                                                                 │
│  ├── CKA (K8s 管理员)                                                        │
│  ├── CKAD (K8s 应用开发者)                                                   │
│  └── CKS (K8s 安全专家)                                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## ❓ 常见问题

**Q: 我的电脑配置不够，如何学习 K8s?**

A: 使用以下轻量级方案:

- [K3s](https://k3s.io/): 单二进制，512MB 内存即可运行
- [Killercoda](https://killercoda.com/): 免费在线 K8s 环境
- [Play with Kubernetes](https://labs.play-with-k8s.com/): 浏览器中的 K8s

**Q: 学习云原生需要会编程吗?**

A: 基础操作不需要，但进阶需要:

- 运维: Shell 脚本 + YAML
- 开发: Go/Python + 客户端库
- 架构: 理解系统设计的编程思维

**Q: 从哪个技术开始学?**

A: 推荐顺序:

1. Linux 基础 (如果还不熟悉)
2. Docker (容器基础)
3. Kubernetes (编排核心)
4. 根据兴趣选择: 网络(Cilium)/可观测性(Prometheus)/GitOps(ArgoCD)

**Q: 有推荐的学习社区吗?**

A:

- 中文: K8s 技术社区、云原生社区
- 国际: CNCF Slack、Kubernetes 论坛、Reddit r/kubernetes
