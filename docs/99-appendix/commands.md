# 命令速查表

> 常用 Docker & Kubernetes 命令

---

## Docker 命令

### 镜像管理

| 命令 | 说明 |
|------|------|
| `docker images` | 列出镜像 |
| `docker pull <image>` | 拉取镜像 |
| `docker build -t <tag> .` | 构建镜像 |
| `docker push <image>` | 推送镜像 |
| `docker rmi <image>` | 删除镜像 |
| `docker inspect <image>` | 查看镜像详情 |

### 容器管理

| 命令 | 说明 |
|------|------|
| `docker run -d <image>` | 后台运行容器 |
| `docker ps` | 列出运行中的容器 |
| `docker ps -a` | 列出所有容器 |
| `docker stop <container>` | 停止容器 |
| `docker rm <container>` | 删除容器 |
| `docker logs <container>` | 查看日志 |
| `docker exec -it <container> sh` | 进入容器 |

### 网络与存储

| 命令 | 说明 |
|------|------|
| `docker network ls` | 列出网络 |
| `docker volume ls` | 列出卷 |
| `docker system prune` | 清理未使用资源 |

---

## Kubernetes 命令

### 资源管理

| 命令 | 说明 |
|------|------|
| `kubectl get <resource>` | 列出资源 |
| `kubectl describe <resource> <name>` | 查看详情 |
| `kubectl apply -f <file>` | 应用配置 |
| `kubectl delete -f <file>` | 删除资源 |
| `kubectl edit <resource> <name>` | 编辑资源 |

### Pod 操作

| 命令 | 说明 |
|------|------|
| `kubectl get pods` | 列出 Pod |
| `kubectl logs <pod>` | 查看日志 |
| `kubectl logs -f <pod>` | 实时查看日志 |
| `kubectl exec -it <pod> -- sh` | 进入容器 |
| `kubectl port-forward <pod> <port>` | 端口转发 |
| `kubectl cp <pod>:<path> <local>` | 复制文件 |

### 调试命令

| 命令 | 说明 |
|------|------|
| `kubectl get events` | 查看事件 |
| `kubectl top nodes` | 节点资源使用 |
| `kubectl top pods` | Pod 资源使用 |
| `kubectl cluster-info` | 集群信息 |

### 快捷命令

```bash
# 设置别名
alias k=kubectl
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgd='kubectl get deploy'
alias kdf='kubectl delete -f'
alias kaf='kubectl apply -f'

# 常用组合
k get pods -o wide          # 详细信息
k get pods --show-labels    # 显示标签
k get pods -A               # 所有命名空间
k logs -f <pod> --tail=100  # 最后100行
```

---

## kubectl 输出格式

| 选项 | 说明 |
|------|------|
| `-o yaml` | YAML 格式 |
| `-o json` | JSON 格式 |
| `-o wide` | 详细信息 |
| `-o name` | 仅名称 |
| `-o custom-columns=...` | 自定义列 |
