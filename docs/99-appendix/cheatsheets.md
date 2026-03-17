# 速查表

---

## Docker 命令速查

### 镜像管理

```bash
# 构建镜像
docker build -t myapp:v1.0.0 .
docker build --target production -t myapp:v1.0.0 .
docker build --platform linux/amd64,linux/arm64 -t myapp:v1.0.0 .

# 镜像操作
docker images
docker rmi myapp:v1.0.0
docker tag myapp:v1.0.0 myregistry/myapp:v1.0.0
docker push myregistry/myapp:v1.0.0
```

### 容器管理

```bash
# 运行容器
docker run -d -p 8080:8080 --name myapp myapp:v1.0.0
docker run --rm -it myapp:v1.0.0 sh

# 容器操作
docker ps -a
docker start/stop/restart myapp
docker rm -f myapp
docker exec -it myapp sh

# 日志和监控
docker logs -f myapp
docker stats
docker top myapp
```

### 清理命令

```bash
docker system prune -f              # 清理未使用数据
docker volume prune -f              # 清理未使用卷
docker image prune -f               # 清理未使用镜像
```

---

## Kubernetes 命令速查

### 基础操作

```bash
# 查看资源
kubectl get pods -o wide
kubectl get nodes
kubectl get svc,deploy,ingress
kubectl get all -n namespace

# 详细信息
kubectl describe pod <pod-name>
kubectl logs <pod-name> -f
kubectl logs <pod-name> --previous
kubectl exec -it <pod-name> -- sh
kubectl top pods/nodes
```

### 创建与更新

```bash
kubectl apply -f config.yaml
kubectl apply -k kustomize/
kubectl create deployment nginx --image=nginx
kubectl set image deployment/nginx nginx=nginx:1.21
kubectl scale deployment nginx --replicas=5
kubectl rollout status deployment/nginx
kubectl rollout undo deployment/nginx
```

### 调试命令

```bash
kubectl port-forward pod/<pod-name> 8080:80
kubectl cp <pod-name>:/path/file ./file
kubectl debug pod/<pod-name> -it --image=busybox
kubectl get events --sort-by='.lastTimestamp'
```

### 资源模板

```bash
# 生成 YAML
kubectl run nginx --image=nginx --dry-run=client -o yaml
kubectl create deployment nginx --image=nginx --dry-run=client -o yaml
```

---

## kubectl 快捷别名

```bash
# 添加 ~/.bashrc 或 ~/.zshrc
alias k='kubectl'
alias kg='kubectl get'
alias kd='kubectl describe'
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgd='kubectl get deployment'
alias kl='kubectl logs'
alias klf='kubectl logs -f'
alias ke='kubectl exec -it'
alias ka='kubectl apply -f'
alias kdel='kubectl delete'
alias kctx='kubectl config use-context'
alias kns='kubectl config set-context --current --namespace'
```

---

## YAML 模板速查

### Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
  labels:
    app: myapp
spec:
  containers:
    - name: app
      image: nginx:alpine
      ports:
        - containerPort: 80
```

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
        - name: app
          image: myapp:v1.0.0
          ports:
            - containerPort: 8080
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  selector:
    app: myapp
  ports:
    - port: 80
      targetPort: 8080
  type: ClusterIP
```
