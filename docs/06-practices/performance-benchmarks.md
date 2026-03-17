# 性能基准与测试

> 容器与 K8s 性能评估指南

---

## 容器运行时性能对比

### 启动时间

| 运行时 | 冷启动 | 镜像拉取 | 总时间 |
|--------|--------|----------|--------|
| **runc** | ~500ms | 30s | ~30.5s |
| **youki** | ~300ms | 30s | ~30.3s |
| **runwasi** | ~20ms | 1s | ~1.02s |
| **Firecracker** | ~125ms | N/A | ~125ms |

### 内存占用

```
┌─────────────────────────────────────────────────────────────┐
│              容器运行时内存占用 (空闲)                        │
├─────────────────────────────────────────────────────────────┤
│  Docker Engine        ████████████████████  ~100MB          │
│  containerd           ████████              ~30MB           │
│  CRI-O                ████████              ~25MB           │
│  youki                ██                    ~5MB            │
│  runwasi              █                     ~2MB            │
└─────────────────────────────────────────────────────────────┘
```

---

## CNI 性能基准

### 网络延迟 (Pod-to-Pod)

| CNI | P50 延迟 | P99 延迟 | CPU 占用 |
|-----|----------|----------|----------|
| **Flannel (VXLAN)** | 250μs | 500μs | 中 |
| **Calico (iptables)** | 200μs | 400μs | 高 |
| **Calico (eBPF)** | 150μs | 300μs | 中 |
| **Cilium (eBPF)** | **60μs** | **120μs** | **低** |

### 吞吐量测试

```bash
# 使用 iperf3 测试
kubectl run iperf-server --image=networkstatic/iperf3 -- iperf3 -s
kubectl run iperf-client --image=networkstatic/iperf3 -- iperf3 -c iperf-server

# 结果示例 (10Gbps 网络)
# Flannel:  8.5 Gbps
# Calico:   9.2 Gbps
# Cilium:   9.8 Gbps
```

---

## K8s 调度性能

### 大规模集群测试

```
┌─────────────────────────────────────────────────────────────┐
│              集群规模 vs 调度延迟                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  节点数    Pod 数      调度延迟      建议                    │
│  ────────────────────────────────────────────────────────   │
│  100       5,000      < 1s          小规模                  │
│  1,000     50,000     < 2s          中等规模                │
│  5,000     250,000    < 5s          大规模                  │
│  10,000    500,000    < 10s         超大规模 (需优化)        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 优化建议

```yaml
# 1. 启用优先级和抢占
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: high-priority
value: 1000000
globalDefault: false
preemptionPolicy: PreemptLowerPriority

# 2. 配置调度器性能参数
apiVersion: v1
kind: ConfigMap
metadata:
  name: scheduler-config
data:
  scheduler.yaml: |
    apiVersion: kubescheduler.config.k8s.io/v1
    kind: KubeSchedulerConfiguration
    profiles:
      - schedulerName: default-scheduler
    percentageOfNodesToScore: 50  # 只评估 50% 节点
```

---

## 存储性能测试

### CSI 驱动对比

| 驱动 | IOPS (4K随机读) | 延迟 (顺序写) | 适用场景 |
|------|-----------------|---------------|----------|
| **EBS CSI** | 16,000 | ~5ms | AWS 云盘 |
| **Azure Disk** | 20,000 | ~3ms | Azure 云盘 |
| **Ceph RBD** | 50,000 | ~2ms | 自建存储 |
| **Local SSD** | 500,000 | ~0.1ms | 高性能本地 |

### 测试命令

```bash
# 使用 fio 测试
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: fio-test
spec:
  containers:
    - name: fio
      image: nixery.dev/shell/fio
      command: ["fio"]
      args:
        - --name=random-write
        - --ioengine=libaio
        - --iodepth=32
        - --rw=randwrite
        - --bs=4k
        - --direct=1
        - --size=4G
        - --numjobs=4
        - --runtime=60
        - --group_reporting
      volumeMounts:
        - name: data
          mountPath: /data
  volumes:
    - name: data
      persistentVolumeClaim:
        claimName: test-pvc
EOF
```

---

## WebAssembly vs Docker 性能

### 冷启动对比

```
┌─────────────────────────────────────────────────────────────┐
│              冷启动时间对比 (越小越好)                        │
├─────────────────────────────────────────────────────────────┤
│  Docker (500MB镜像)   ████████████████████████████████████ │
│  Docker (100MB镜像)   ████████████                          │
│  Docker (10MB镜像)    ███                                   │
│  Wasm (10MB模块)      █                                     │
│                                                              │
│  0ms      100ms     500ms     1s       2s        5s         │
└─────────────────────────────────────────────────────────────┘
```

### 密度测试

| 运行时 | 单节点 Pod 数 | 内存/实例 | 启动时间 |
|--------|--------------|-----------|----------|
| **Docker** | ~500 | 100MB | 2s |
| **gVisor** | ~200 | 200MB | 3s |
| **Kata** | ~100 | 512MB | 5s |
| **Wasm** | **~5000** | **5MB** | **20ms** |

---

## 性能调优清单

### 容器层面

- [ ] 使用多阶段构建减小镜像
- [ ] 设置合适的资源请求/限制
- [ ] 启用 CPU Manager 静态策略
- [ ] 使用 HugePages 大页内存

### 网络层面

- [ ] 使用 Cilium 替代 kube-proxy
- [ ] 启用 eBPF 加速
- [ ] 配置合适的 MTU
- [ ] 使用 RDMA (如适用)

### 存储层面

- [ ] 使用 Local SSD 或 NVMe
- [ ] 启用存储缓存
- [ ] 配置合适的 StorageClass
- [ ] 使用 VolumeSnapshot 备份

### 调度层面

- [ ] 启用优先级和抢占
- [ ] 配置 Pod 亲和性/反亲和性
- [ ] 使用拓扑感知调度
- [ ] 启用 DRA (动态资源分配)
