# 开发工具链

> 提升容器开发效率的工具

---

## IDE 集成

### VS Code 扩展

| 扩展名 | 功能 |
|--------|------|
| Docker | 容器管理、Dockerfile 支持 |
| Kubernetes | K8s 资源管理、Helm 支持 |
| YAML | YAML 验证、格式化 |

### 推荐 CLI 工具

| 工具 | 用途 |
|------|------|
| **kubectx** | 快速切换集群/命名空间 |
| **k9s** | 终端 K8s UI |
| **Lens** | K8s IDE |
| **stern** | 多 Pod 日志查看 |

---

## 本地开发环境

### Kind

```bash
# 创建集群
kind create cluster --name dev
```

### Minikube

```bash
# 启动
minikube start
```

---

## 调试工具

```bash
# 临时调试容器
kubectl debug <pod> -it --image=nicolaka/netshoot
```
