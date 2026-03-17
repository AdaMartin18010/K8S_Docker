# kubectl 速查表

> Kubernetes 命令行工具快速参考 (v1.33/1.34)

---

## 基础命令

```bash
# 查看版本
kubectl version
kubectl version --client

# 集群信息
kubectl cluster-info
kubectl cluster-info dump

# 配置信息
kubectl config view
kubectl config get-contexts
kubectl config use-context <context>
kubectl config set-context --current --namespace=<ns>
```

---

## 资源操作

```bash
# 查看资源
kubectl get pods
kubectl get pods -o wide
kubectl get pods -n kube-system
kubectl get pods --all-namespaces
kubectl get pods -l app=nginx
kubectl get pods --show-labels
kubectl get pods -w  # 监视

# 查看详情
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl logs <pod-name> -c <container>
kubectl logs <pod-name> --tail=100 -f
kubectl logs <pod-name> --previous

# 进入容器
kubectl exec -it <pod-name> -- /bin/sh
kubectl exec <pod-name> -- ls /
kubectl cp <pod-name>:/file ./local-file
kubectl cp ./local-file <pod-name>:/file
```

---

## 创建与删除

```bash
# 创建资源
kubectl apply -f manifest.yaml
kubectl apply -f ./directory/
kubectl apply -k kustomization/
kubectl create -f manifest.yaml
kubectl create deployment nginx --image=nginx

# 删除资源
kubectl delete -f manifest.yaml
kubectl delete pod <pod-name>
kubectl delete pod <pod-name> --grace-period=0 --force
kubectl delete pods -l app=nginx
kubectl delete all --all  # 删除所有资源（谨慎）
```

---

## Deployment 操作

```bash
# 查看
kubectl get deployments
kubectl get deploy -n default
kubectl describe deploy <name>

# 扩缩容
kubectl scale deploy <name> --replicas=5
kubectl autoscale deploy <name> --min=2 --max=10 --cpu-percent=80

# 更新镜像
kubectl set image deploy/<name> container=new-image:v2
kubectl edit deploy <name>

# 回滚
kubectl rollout status deploy/<name>
kubectl rollout history deploy/<name>
kubectl rollout undo deploy/<name>
kubectl rollout undo deploy/<name> --to-revision=2
kubectl rollout pause deploy/<name>
kubectl rollout resume deploy/<name>
```

---

## 调试技巧

```bash
# 查看事件
kubectl get events --sort-by='.lastTimestamp'
kubectl get events --field-selector type=Warning

# 端口转发
kubectl port-forward pod/<name> 8080:80
kubectl port-forward svc/<name> 8080:80
kubectl port-forward deploy/<name> 8080:80

# 调试容器
kubectl debug pod/<name> -it --image=nicolaka/netshoot
kubectl debug pod/<name> -it --copy-to=debug-pod --image=busybox

# 运行临时 Pod
kubectl run debug --rm -it --image=busybox -- /bin/sh
kubectl run curl --rm -it --image=curlimages/curl -- <url>
```

---

## 高级查询

```bash
# JSONPath 查询
kubectl get pods -o jsonpath='{.items[*].metadata.name}'
kubectl get pods -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.phase}{"\n"}{end}'

# 自定义列
kubectl get pods -o custom-columns='NAME:.metadata.name,STATUS:.status.phase,IP:.status.podIP'

# 排序
kubectl get pods --sort-by=.metadata.name
kubectl get pods --sort-by=.status.startTime

# 导出 YAML
kubectl get pod <name> -o yaml > pod.yaml
kubectl get pod <name> -o yaml --export > pod.yaml
```

---

## 常用快捷方式

```bash
# 别名（推荐添加到 .bashrc/.zshrc）
alias k='kubectl'
alias kc='kubectl config'
alias kg='kubectl get'
alias kd='kubectl describe'
alias kdel='kubectl delete'
alias ka='kubectl apply'
alias ke='kubectl edit'
alias kx='kubectl exec'
alias kl='kubectl logs'
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgn='kubectl get nodes'
alias kgd='kubectl get deploy'

# 上下文和命名空间
alias kctx='kubectl config use-context'
alias kns='kubectl config set-context --current --namespace'
```
