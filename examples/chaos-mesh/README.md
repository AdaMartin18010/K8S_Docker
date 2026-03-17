# Chaos Mesh 混沌工程示例

> 使用 Chaos Mesh 进行故障注入实验

---

## 安装 Chaos Mesh

```bash
# 安装 Chaos Mesh
kubectl apply -f https://mirrors.chaos-mesh.org/v2.6.3/install.yaml

# 验证安装
kubectl get pods -n chaos-mesh
```

---

## 实验类型

### 1. Pod 故障实验

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: pod-failure-example
spec:
  action: pod-failure
  mode: one
  duration: "30s"
  selector:
    namespaces:
      - default
    labelSelectors:
      app: my-app
```

### 2. 网络延迟实验

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: network-delay-example
spec:
  action: delay
  mode: all
  selector:
    namespaces:
      - default
    labelSelectors:
      app: my-app
  delay:
    latency: "100ms"
    correlation: "100"
    jitter: "0ms"
  duration: "5m"
```

### 3. CPU 压力实验

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: StressChaos
metadata:
  name: cpu-stress-example
spec:
  mode: one
  selector:
    namespaces:
      - default
    labelSelectors:
      app: my-app
  stressors:
    cpu:
      workers: 2
      load: 80
  duration: "5m"
```

---

## 运行实验

```bash
# 应用实验
kubectl apply -f pod-failure.yaml

# 查看实验状态
kubectl get podchaos

# 删除实验
kubectl delete -f pod-failure.yaml
```

---

## 相关文档

- [Chaos Mesh 官方文档](https://chaos-mesh.org/docs/)
- [混沌工程指南](../../docs/06-practices/chaos-engineering/)
