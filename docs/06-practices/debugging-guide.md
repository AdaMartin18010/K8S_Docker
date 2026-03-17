# Kubernetes 故障排查完全指南

> 从症状到解决方案的系统化排查方法

---

## 故障排查流程图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          故障排查通用流程                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 收集信息                                                                 │
│     ├── kubectl get pods -o wide          # 查看 Pod 状态                    │
│     ├── kubectl describe pod <name>       # 查看详情                         │
│     └── kubectl logs <pod>                # 查看日志                         │
│              │                                                               │
│              ▼                                                               │
│  2. 识别症状 ───────► 状态分类                                                │
│              │                                                               │
│              ├── Pending                                                    │
│              ├── CrashLoopBackOff                                           │
│              ├── ImagePullBackOff                                           │
│              ├── OOMKilled                                                  │
│              ├── Evicted                                                    │
│              └── Running (但无法访问)                                        │
│              │                                                               │
│              ▼                                                               │
│  3. 根因分析 ───────► 定位问题                                                │
│              │                                                               │
│              ▼                                                               │
│  4. 解决问题                                                                 │
│              │                                                               │
│              ▼                                                               │
│  5. 验证修复                                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Pod 状态故障排查

### 1. Pod 处于 Pending 状态

**症状**: Pod 创建后长时间处于 Pending，无法调度到节点

**排查命令**:
```bash
kubectl describe pod <pod-name>
# 查看 Events 部分
```

**常见原因及解决方案**:

| 原因 | 错误信息示例 | 解决方案 |
|------|-------------|----------|
| 资源不足 | `Insufficient cpu` / `Insufficient memory` | 增加节点资源或缩减 Pod 请求 |
| 无可用节点 | `0/3 nodes are available` | 检查节点状态，扩容节点 |
| PVC 未绑定 | `unbound immediate PersistentVolumeClaims` | 检查 StorageClass 和 PV |
| 节点选择器不匹配 | `node(s) didn't match Pod's node affinity` | 检查 nodeSelector/亲和性配置 |
| 污点容忍 | `node(s) had taint {key: value}` | 添加 tolerations 或移除污点 |
| 镜像拉取密钥缺失 | `kubernetes.io/dockerconfigjson not found` | 创建 imagePullSecret |

**快速诊断脚本**:
```bash
#!/bin/bash
POD_NAME=$1
NAMESPACE=${2:-default}

echo "=== Pod 状态 ==="
kubectl get pod $POD_NAME -n $NAMESPACE -o wide

echo -e "\n=== 事件详情 ==="
kubectl describe pod $POD_NAME -n $NAMESPACE | grep -A 20 Events

echo -e "\n=== 资源请求 ==="
kubectl get pod $POD_NAME -n $NAMESPACE -o jsonpath='{range .spec.containers[*]}{.name}{"\t"}{.resources.requests}{"\n"}{end}'

echo -e "\n=== 节点资源 ==="
kubectl top nodes

echo -e "\n=== 调度约束 ==="
kubectl get pod $POD_NAME -n $NAMESPACE -o jsonpath='{.spec.nodeSelector}'
kubectl get pod $POD_NAME -n $NAMESPACE -o jsonpath='{.spec.affinity}'
kubectl get pod $POD_NAME -n $NAMESPACE -o jsonpath='{.spec.tolerations}'
```

---

### 2. Pod 处于 CrashLoopBackOff

**症状**: Pod 反复启动失败，状态在 Running 和 CrashLoopBackOff 间切换

**排查命令**:
```bash
# 查看当前日志
kubectl logs <pod-name>

# 查看上次崩溃的日志
kubectl logs <pod-name> --previous

# 查看重启次数和原因
kubectl get pod <pod-name> -o jsonpath='{.status.containerStatuses[0].lastState}'
```

**常见原因及解决方案**:

| 原因 | 排查方法 | 解决方案 |
|------|----------|----------|
| 应用启动失败 | `kubectl logs` | 修复应用代码/配置 |
| 配置错误 | 检查 ConfigMap/Secret | 更正配置文件 |
| 依赖服务不可用 | `kubectl logs` 查看连接错误 | 启动依赖服务或配置重试 |
| 健康检查失败 | 检查 livenessProbe | 调整探针参数或修复应用 |
| 权限不足 | `kubectl logs` 查看 permission denied | 配置 SecurityContext |
| 资源限制 | 检查 exit code 137/143 | 增加资源限制 |

**Exit Code 参考**:
- 0: 正常退出
- 1: 应用错误
- 137 (128+9): SIGKILL (OOM 或强制终止)
- 143 (128+15): SIGTERM (优雅终止)

---

### 3. Pod 处于 ImagePullBackOff

**症状**: 无法拉取容器镜像，状态为 ImagePullBackOff 或 ErrImagePull

**排查命令**:
```bash
kubectl describe pod <pod-name>
# 查看 Events 中的拉取错误
```

**常见原因及解决方案**:

| 原因 | 错误信息 | 解决方案 |
|------|----------|----------|
| 镜像不存在 | `not found` | 检查镜像名称和标签 |
| 私有仓库未认证 | `pull access denied` | 创建 imagePullSecret |
| 网络问题 | `timeout` / `no such host` | 检查网络连接 |
| 镜像标签错误 | `manifest unknown` | 使用正确的镜像标签 |
| 平台不匹配 | `no matching manifest` | 使用多架构镜像或指定平台 |

**配置私有仓库凭证**:
```bash
# 创建 Secret
docker login registry.example.com
kubectl create secret generic regcred \
  --from-file=.dockerconfigjson=$HOME/.docker/config.json \
  --type=kubernetes.io/dockerconfigjson

# 在 Pod 中引用
# spec:
#   imagePullSecrets:
#   - name: regcred
```

---

### 4. Pod 被 OOMKilled

**症状**: Pod 因内存不足被系统终止，exit code 为 137

**排查命令**:
```bash
# 查看 Pod 状态
kubectl get pod <pod-name> -o jsonpath='{.status.containerStatuses[0].lastState}'

# 查看内存使用历史
kubectl top pod <pod-name>
```

**解决方案**:

```yaml
# 增加内存限制
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    resources:
      limits:
        memory: "512Mi"  # 增加限制
      requests:
        memory: "256Mi"  # 合理设置请求
```

**排查内存泄漏**:
```bash
# 进入容器查看内存使用
kubectl exec -it <pod-name> -- sh
# 然后运行: top, ps aux, 或应用内置指标

# 查看历史内存使用 (需要 Metrics Server)
kubectl top pod <pod-name> --containers
```

---

### 5. Pod 被 Evicted

**症状**: Pod 被节点驱逐，通常由于资源压力

**排查命令**:
```bash
kubectl get pod <pod-name> -o jsonpath='{.status.reason}'
kubectl describe pod <pod-name>
```

**驱逐原因**:
- `DiskPressure`: 磁盘空间不足
- `MemoryPressure`: 内存不足
- `PIDPressure`: 进程数过多
- `NodeNotReady`: 节点不健康

**解决方案**:
1. 清理节点资源
2. 增加 Pod 的 QoS 等级
3. 配置 PodDisruptionBudget

---

## 网络故障排查

### Service 无法访问

```bash
# 1. 检查 Service 是否存在
kubectl get svc <service-name>

# 2. 检查 Endpoints
kubectl get endpoints <service-name>
# 如果为空，说明 Selector 不匹配或 Pod 未就绪

# 3. 检查 Pod 标签
kubectl get pods --show-labels
# 确认 labels 与 Service selector 匹配

# 4. 测试 Service 连通性
kubectl run -it --rm debug --image=nicolaka/netshoot --restart=Never -- nslookup <service-name>

# 5. 检查 kube-proxy
kubectl get pods -n kube-system -l k8s-app=kube-proxy
kubectl logs -n kube-system -l k8s-app=kube-proxy
```

### Ingress/Gateway 无法访问

```bash
# 1. 检查 Ingress/Gateway 配置
kubectl get ingress
gatewayctl get gateways

# 2. 检查 Controller Pod
kubectl get pods -n ingress-nginx  # 或其他命名空间
kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx

# 3. 检查后端 Service
kubectl get endpoints <backend-service>

# 4. 本地测试
kubectl port-forward svc/<service-name> 8080:80
# 然后 curl localhost:8080
```

### DNS 解析问题

```bash
# 测试 DNS 解析
kubectl run -it --rm debug --image=nicolaka/netshoot --restart=Never -- nslookup kubernetes.default

# 检查 CoreDNS
kubectl get pods -n kube-system -l k8s-app=kube-dns
kubectl logs -n kube-system -l k8s-app=kube-dns

# 检查 DNS 配置
cat /etc/resolv.conf  # 在 Pod 内执行
```

---

## 存储故障排查

### PVC 无法绑定

```bash
# 查看 PVC 状态
kubectl get pvc
kubectl describe pvc <pvc-name>

# 检查 StorageClass
kubectl get storageclass

# 查看 PV
kubectl get pv
```

**常见问题**:
- 无可用 PV: 检查 StorageClass 是否支持动态供应
- 权限不足: 检查 CSI 驱动的 RBAC 配置

### Pod 无法挂载卷

```bash
# 查看 Pod 事件
kubectl describe pod <pod-name> | grep -A 10 Events

# 检查节点上的挂载
kubectl get pod <pod-name> -o jsonpath='{.spec.nodeName}'
# SSH 到节点检查: mount | grep <volume>
```

---

## 性能问题排查

### CPU 使用率高

```bash
# 查看 Pod CPU 使用
kubectl top pods
kubectl top pods --containers

# 查看节点 CPU
kubectl top nodes

# 进入容器分析
kubectl exec -it <pod-name> -- sh
# top, pidstat, perf 等工具
```

### 应用响应慢

```bash
# 检查资源限制是否成为瓶颈
kubectl describe pod <pod-name> | grep -A 5 "Limits"

# 检查 HPA 状态
kubectl get hpa

# 查看网络延迟
kubectl exec -it <pod-name> -- ping <target>

# 查看连接数
kubectl exec -it <pod-name> -- netstat -an | wc -l
```

---

## 常用诊断工具

### 调试容器

```bash
# 使用临时调试容器 (Ephemeral Containers)
kubectl debug -it <pod-name> --image=nicolaka/netshoot --target=<container-name>

# 常用调试镜像
# - nicolaka/netshoot: 网络调试工具集
# - busybox: 轻量级工具
# - alpine: 基础环境
```

### 网络抓包

```bash
# 在节点上抓包 (需要特权)
kubectl debug node/<node-name> -it --image=nicolaka/netshoot
# 然后: tcpdump -i any -w /tmp/capture.pcap

# 在 Pod 网络命名空间抓包
kubectl exec -it <pod-name> -- tcpdump -i eth0
```

### 日志分析

```bash
# 聚合日志 (需要 Stern 或 kubetail)
stern <pod-name-pattern>
kubetail <pod-name>

# 导出日志
kubectl logs <pod-name> --since=24h > pod.log

# 多容器 Pod
kubectl logs <pod-name> -c <container-name>
```

---

## 故障排查检查清单

```markdown
## Pod 故障检查清单

### 基本信息
- [ ] Pod 状态: kubectl get pod <name> -o wide
- [ ] Pod 详情: kubectl describe pod <name>
- [ ] Pod 日志: kubectl logs <name>
- [ ] 之前日志: kubectl logs <name> --previous

### 资源配置
- [ ] 资源请求和限制是否合理
- [ ] QoS 等级 (Guaranteed/Burstable/BestEffort)
- [ ] 节点资源是否充足

### 调度相关
- [ ] 节点选择器配置
- [ ] 亲和性/反亲和性配置
- [ ] 污点容忍配置

### 网络和存储
- [ ] Service 和 Endpoints
- [ ] PVC 绑定状态
- [ ] 网络策略限制

### 安全和权限
- [ ] ServiceAccount
- [ ] SecurityContext
- [ ] RBAC 权限
```

---

## 常见错误速查

| 错误信息 | 可能原因 | 解决方案 |
|----------|----------|----------|
| `CrashLoopBackOff` | 应用崩溃 | 查看日志修复 |
| `ImagePullBackOff` | 镜像拉取失败 | 检查镜像和凭证 |
| `Pending` | 无法调度 | 检查资源和约束 |
| `OOMKilled` | 内存不足 | 增加内存限制 |
| `Evicted` | 节点资源压力 | 清理资源或扩容 |
| `CreateContainerConfigError` | 配置错误 | 检查 ConfigMap/Secret |
| `InvalidImageName` | 镜像名格式错误 | 修正镜像名称 |
| `NodeAffinity` | 节点选择失败 | 检查亲和性配置 |
| `Taints` | 污点不匹配 | 添加 tolerations |
| `UnboundPVC` | PVC 未绑定 | 检查 StorageClass |
