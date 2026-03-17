# 深度故障排查指南

> Kubernetes 生产环境问题诊断

---

## 排查方法论

```
┌─────────────────────────────────────────────────────────────┐
│              系统化故障排查流程                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. 现象观察                                                 │
│     • 错误消息、日志                                         │
│     • 监控指标变化                                           │
│     • 用户反馈                                               │
│                         ↓                                    │
│  2. 信息收集                                                 │
│     • kubectl describe / logs                               │
│     • 事件查看                                               │
│     • 网络抓包                                               │
│                         ↓                                    │
│  3. 假设验证                                                 │
│     • 缩小范围                                               │
│     • 对比正常/异常状态                                      │
│     • 重现问题                                               │
│                         ↓                                    │
│  4. 根因定位                                                 │
│     • 代码审查                                               │
│     • 配置检查                                               │
│     • 基础设施排查                                           │
│                         ↓                                    │
│  5. 修复验证                                                 │
│     • 应用修复                                               │
│     • 验证恢复                                               │
│     • 预防措施                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## Pod 启动问题

### Pending 状态排查

```bash
# 查看 Pod 事件
kubectl describe pod <pod-name> | grep -A 10 Events

# 常见原因:
# 1. 资源不足
kubectl describe node <node-name> | grep -A 5 "Allocated resources"

# 2. 节点亲和性
kubectl get pod <pod-name> -o yaml | grep -A 10 nodeAffinity

# 3. PVC 未绑定
kubectl get pvc
kubectl describe pvc <pvc-name>

# 4. 污点和容忍
kubectl describe node <node-name> | grep Taints
```

### CrashLoopBackOff 排查

```bash
# 查看上次启动日志
kubectl logs <pod-name> --previous

# 查看退出码
docker inspect <container-id> | grep ExitCode

# 常见退出码:
# 0   - 正常退出
# 1   - 应用错误
# 137 - SIGKILL (OOM)
# 143 - SIGTERM

# 检查 OOM
kubectl describe pod <pod-name> | grep -i oom
kubectl get pod <pod-name> -o yaml | grep -A 5 "lastState"
```

---

## 网络问题诊断

### DNS 问题

```bash
# 测试 DNS 解析
kubectl run -it --rm debug --image=busybox:1.36 -- nslookup kubernetes.default

# 检查 CoreDNS
kubectl get pods -n kube-system -l k8s-app=kube-dns
kubectl logs -n kube-system -l k8s-app=kube-dns

# 检查 DNS 配置
cat /etc/resolv.conf

# 网络抓包
kubectl debug node/<node-name> -it --image=nicolaka/netshoot -- tcpdump -i any port 53
```

### 服务连通性

```bash
# 检查 Endpoints
kubectl get endpoints <service-name>
kubectl get pods -l app=<label-selector>

# 测试连接
kubectl run -it --rm debug --image=nicolaka/netshoot -- curl -v http://<service>:<port>

# 检查网络策略
kubectl get networkpolicies
kubectl describe networkpolicy <policy-name>

# Cilium 连通性测试
cilium connectivity test
```

### 跨节点通信

```bash
# 检查 CNI 配置
cat /etc/cni/net.d/*.conf

# 检查路由表
ip route

# VXLAN 检查
ip -d link show vxlan

# 检查 iptables/eBPF 规则
iptables -L -n -v
bpftool prog list
```

---

## 存储问题排查

### PVC 问题

```bash
# 检查 StorageClass
kubectl get storageclass

# 检查 PVC 事件
kubectl describe pvc <pvc-name>

# 检查 PV
kubectl get pv
kubectl describe pv <pv-name>

# 检查 CSI 驱动
kubectl get pods -n kube-system | grep csi
kubectl logs -n kube-system <csi-driver>
```

### 文件系统满

```bash
# 进入容器检查
df -h

# 查找大文件
du -sh /* 2>/dev/null | sort -rh | head -20

# 检查 inode 使用
df -i

# 清理日志
find /var/log -name "*.log" -mtime +7 -delete
```

---

## 性能问题分析

### CPU 问题

```bash
# 查看 CPU 使用
kubectl top pods -A
kubectl top nodes

# 进入容器分析
kubectl exec -it <pod> -- sh

# 查找高 CPU 进程
top -H -p <pid>

# Go 应用 pprof
curl http://localhost:6060/debug/pprof/profile > cpu.prof
go tool pprof cpu.prof
```

### 内存问题

```bash
# 查看内存使用
kubectl top pod <pod-name> --containers

# OOM 分析
kubectl describe pod <pod-name> | grep -A 3 "Last State"

# Go 内存分析
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# 内存泄漏检测
# 对比多次 heap dump
```

### 网络延迟

```bash
# 使用 netshoot
time nc -zv <host> <port>

# iperf3 测试
iperf3 -c <server>

# tcpdump 抓包
tcpdump -i any -w capture.pcap

# 火焰图分析
perf record -F 99 -a -g -- sleep 30
perf script | stackcollapse-perf.pl | flamegraph.pl > flame.svg
```

---

## 控制平面问题

### API Server

```bash
# 检查 API Server 状态
kubectl get --raw='/healthz'

# 查看 API Server 日志
kubectl logs -n kube-system kube-apiserver-<node>

# 检查证书过期
openssl x509 -in /etc/kubernetes/pki/apiserver.crt -noout -text | grep "Not After"

# ETCD 健康检查
kubectl exec -n kube-system <etcd-pod> -- etcdctl endpoint health
```

### ETCD 问题

```bash
# ETCD 空间检查
kubectl exec -n kube-system <etcd-pod> -- etcdctl endpoint status

# 压缩历史版本
etcdctl compaction $(etcdctl endpoint status --write-out=json | jq .[0].Status.header.revision)
etcdctl defrag

# 备份
etcdctl snapshot save backup.db
```

---

## 实用调试工具

```yaml
# 调试 Pod 模板
apiVersion: v1
kind: Pod
metadata:
  name: debug-pod
spec:
  containers:
    - name: debug
      image: nicolaka/netshoot
      command: ["sleep", "3600"]
      securityContext:
        privileged: true  # 需要时启用
      volumeMounts:
        - name: host-root
          mountPath: /host
  volumes:
    - name: host-root
      hostPath:
        path: /
```

### 常用命令速查

```bash
# 节点调试
kubectl debug node/<node> -it --image=mcr.microsoft.com/dotnet/runtime-deps:6.0

# 复制容器调试
kubectl debug <pod> -it --copy-to=debug-pod --container=debug --image=busybox

# 启动临时容器 (ephemeral containers)
kubectl debug -it <pod> --image=busybox --target=<container>

# 集群信息转储
kubectl cluster-info dump --all-namespaces --output-directory=./cluster-state
```
