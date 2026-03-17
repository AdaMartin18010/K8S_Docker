# Docker & Kubernetes 速查表

---

## Docker 命令

```bash
# 构建
docker build -t myapp:v1.0.0 .
docker build --target production -t myapp:v1.0.0 .

# 运行
docker run -d -p 8080:8080 --name myapp myapp:v1.0.0
docker run --rm -it myapp:v1.0.0 sh

# 管理
docker ps -a
docker logs -f myapp
docker exec -it myapp sh
docker stop myapp && docker rm myapp
docker rmi myapp:v1.0.0

# 清理
docker system prune -f
docker volume prune -f

# Compose
docker compose up -d
docker compose down
docker compose logs -f
docker compose build
```

---

## Kubernetes 命令

```bash
# 基础操作
kubectl get pods -o wide
kubectl get nodes
kubectl get all -n namespace

# 详情
describe pod <pod-name>
kubectl logs <pod-name> -f
kubectl logs <pod-name> --previous
kubectl exec -it <pod-name> -- sh

# 创建/更新
kubectl apply -f config.yaml
kubectl apply -k kustomize/
kubectl create deployment nginx --image=nginx
kubectl set image deployment/nginx nginx=nginx:1.21

# 删除
kubectl delete -f config.yaml
kubectl delete pod <pod-name> --force --grace-period=0

# 调试
kubectl port-forward pod/<pod-name> 8080:80
kubectl cp <pod-name>:/path/file ./file
kubectl top pods
kubectl top nodes
```

---

## kubectl 快捷技巧

```bash
# 设置别名
alias k=kubectl
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgd='kubectl get deployment'
alias kd='kubectl describe'
alias kl='kubectl logs'
alias ke='kubectl exec -it'

# 快速切换上下文
kubectl config use-context prod-cluster
kubectl config current-context

# 导出资源
kubectl get deployment myapp -o yaml > myapp.yaml

# 干运行
kubectl apply -f config.yaml --dry-run=client
kubectl apply -f config.yaml --dry-run=server

# 解释资源
kubectl explain pod.spec.containers
```

---

## 资源 YAML 模板

### Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
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
```

---

## 参考

- [kubectl 官方文档](https://kubernetes.io/docs/reference/kubectl/)
- [Docker CLI 文档](https://docs.docker.com/engine/reference/commandline/cli/)
