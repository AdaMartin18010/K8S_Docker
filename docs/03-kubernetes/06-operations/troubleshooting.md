# K8s 故障排查指南

> 生产环境常见问题的排查方法

---

## 问题排查流程

```
1. 识别症状 → 2. 收集信息 → 3. 定位问题 → 4. 解决问题 → 5. 验证修复
```

---

## 常见问题速查

### Pod 处于 Pending 状态

```bash
# 查看事件
kubectl describe pod <pod-name>

# 常见原因
- 资源不足: 检查节点资源
- 节点亲和性: 检查 nodeSelector/亲和性规则
- PVC 未绑定: 检查存储类
```

### Pod 处于 CrashLoopBackOff

```bash
# 查看日志
kubectl logs <pod-name> --previous

# 常见原因
- 应用启动失败
- 健康检查配置错误
- 资源限制过低
```

### Service 无法访问

```bash
# 检查 Endpoints
kubectl get endpoints <service-name>

# 检查 Pod 标签
kubectl get pods -l app=<label>

# 常见原因
- 标签选择器不匹配
- Pod 未就绪
- 网络策略限制
```

---

## 网络问题排查

### DNS 问题

```bash
# 测试 DNS
kubectl run -it --rm debug --image=busybox:1.28 -- nslookup kubernetes.default

# 检查 CoreDNS
kubectl get pods -n kube-system -l k8s-app=kube-dns
kubectl logs -n kube-system -l k8s-app=kube-dns
```

### 连通性问题

```bash
# 启动调试 Pod
kubectl run -it --rm debug --image=nicolaka/netshoot -- /bin/bash

# 测试连接
nc -zv <service-name> <port>
ping <pod-ip>
```

---

## 资源问题排查

### OOMKilled

```bash
# 查看 Pod 状态
kubectl describe pod <pod-name> | grep -A 5 "State:"

# 解决方案
- 增加 memory limit
- 优化应用内存使用
- 添加 VPA
```

### CPU 节流

```bash
# 查看 CPU 节流
kubectl top pod <pod-name>

# 解决方案
- 增加 CPU limit
- 优化应用性能
```

---

## 存储问题排查

### PVC 无法绑定

```bash
# 查看 PVC 状态
kubectl get pvc
kubectl describe pvc <pvc-name>

# 检查 StorageClass
kubectl get storageclass
```

### 文件系统满

```bash
# 进入容器查看
kubectl exec -it <pod-name> -- df -h

# 清理或扩容
```

---

## 控制平面问题

### API Server 不可用

```bash
# 检查 API Server Pod
kubectl get pods -n kube-system -l component=kube-apiserver

# 检查证书
openssl x509 -in /etc/kubernetes/pki/apiserver.crt -noout -text | grep "Not After"
```

### etcd 问题

```bash
# 检查 etcd 健康
kubectl exec -n kube-system <etcd-pod> -- etcdctl endpoint health

# 检查 etcd 空间
kubectl exec -n kube-system <etcd-pod> -- etcdctl endpoint status
```

---

## 实用工具

```bash
# 集群诊断
kubectl cluster-info dump

# 节点诊断
kubectl debug node/<node-name> -it --image=mcr.microsoft.com/dotnet/runtime-deps:6.0

# Pod 调试
kubectl debug <pod-name> -it --copy-to=debug-pod --container=debug --image=busybox
```
