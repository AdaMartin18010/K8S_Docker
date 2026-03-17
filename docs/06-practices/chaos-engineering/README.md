# 混沌工程 (Chaos Engineering)

> 通过主动注入故障提升系统韧性

---

## 什么是混沌工程？

混沌工程是在生产环境或类生产环境中主动注入故障，以验证系统韧性和发现潜在弱点的工程实践。

```
┌─────────────────────────────────────────────────────────────┐
│              混沌工程原则 (Chaos Engineering Principles)      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. 建立稳态假设                                              │
│     定义系统正常运行的可测量指标                               │
│                                                              │
│  2. 引入真实世界的故障                                        │
│     模拟网络延迟、Pod 崩溃、节点故障等                         │
│                                                              │
│  3. 在生产环境运行                                            │
│     只有在真实环境才能发现真实问题                             │
│                                                              │
│  4. 最小化爆炸半径                                            │
│     控制影响范围，避免用户感知                                 │
│                                                              │
│  5. 自动化持续运行                                            │
│     将混沌实验集成到 CI/CD 流水线                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 混沌工程收益

| 指标 | 改进 |
|------|------|
| **MTTR** | 平均修复时间降低 90% |
| **可用性** | 提升至 99.9%+ |
| **故障发现** | 提前发现潜在问题 |
| **团队信心** | 增强对系统的信心 |

---

## 主流工具对比

| 工具 | 特点 | 适用场景 |
|------|------|----------|
| **Chaos Mesh** | CNCF 孵化项目，K8s 原生，可视化界面 | K8s 环境首选 |
| **Litmus** | CNCF 项目，丰富的实验库，GitOps 友好 | CI/CD 集成 |
| **Gremlin** | 企业级，SaaS 服务，GUI 强大 | 企业用户 |
| **ChaosBlade** | 阿里巴巴开源，CLI 优先，轻量级 | 快速实验 |
| **ToxiProxy** | 网络故障注入，开发测试环境 | 网络测试 |

---

## Chaos Mesh 实战

### 安装

```bash
# 使用 Helm 安装
helm install chaos-mesh chaos-mesh/chaos-mesh \
  --namespace=chaos-testing \
  --create-namespace \
  --version 2.7.0

# 访问 Dashboard
kubectl port-forward -n chaos-testing svc/chaos-dashboard 2333:2333
```

### Pod 故障注入

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: pod-kill-example
  namespace: chaos-testing
spec:
  action: pod-kill          # pod-failure, container-kill
  mode: one                 # one, all, fixed, fixed-percent, random-max-percent
  selector:
    namespaces:
      - production
    labelSelectors:
      app: payment-service
  duration: 30s
  gracePeriod: 0            # 立即终止
```

### 网络延迟注入

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: network-delay
spec:
  action: delay
  mode: all
  selector:
    namespaces:
      - production
    labelSelectors:
      app: frontend
  delay:
    latency: 100ms
    correlation: 100
    jitter: 0ms
  duration: 5m
  direction: to             # to, from, both
  target:
    selector:
      namespaces:
        - production
      labelSelectors:
        app: database
    mode: all
```

### 压力测试

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: StressChaos
metadata:
  name: cpu-stress
spec:
  mode: one
  selector:
    namespaces:
      - production
    labelSelectors:
      app: worker
  stressors:
    cpu:
      workers: 4
      load: 80              # CPU 使用率 80%
    memory:
      workers: 2
      size: 1Gi             # 分配 1GB 内存
  duration: 10m
```

---

## Litmus 实战

### 安装

```bash
# 安装 Litmus
kubectl apply -f https://litmuschaos.github.io/litmus/litmus-operator-v3.0.0.yaml

# 安装 ChaosCenter
kubectl apply -f https://litmuschaos.github.io/litmus/3.0.0/litmus-3.0.0.yaml
```

### 混沌工作流

```yaml
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: nginx-chaos
  namespace: litmus
spec:
  appinfo:
    appns: 'default'
    applabel: 'app=nginx'
    appkind: 'deployment'
  # 定义稳态探针
  experiments:
    - name: pod-delete
      spec:
        components:
          env:
            - name: TOTAL_CHAOS_DURATION
              value: '30'
            - name: CHAOS_INTERVAL
              value: '10'
            - name: FORCE
              value: 'false'
          probe:
            - name: check-nginx-access
              type: httpProbe
              mode: Continuous
              runProperties:
                probeTimeout: 5s
                retry: 2
                interval: 5s
                probePollingInterval: 2s
                initialDelay: 2s
              httpProbe/inputs:
                url: http://nginx.default.svc.cluster.local:80
                insecureSkipVerify: false
                method:
                  get:
                    criteria: ==
                    responseCode: '200'
```

---

## CI/CD 集成

### GitHub Actions + Chaos Mesh

```yaml
name: Chaos Engineering

on: [push, pull_request]

jobs:
  chaos-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Kubernetes
        uses: helm/kind-action@v1

      - name: Install Chaos Mesh
        run: |
          helm install chaos-mesh chaos-mesh/chaos-mesh \
            --namespace chaos-testing --create-namespace

      - name: Run Chaos Experiment
        run: |
          kubectl apply -f chaos-experiments/pod-kill.yaml
          sleep 60
          kubectl get pods -n production

      - name: Verify Recovery
        run: |
          # 验证系统是否恢复正常
          curl -f http://myapp.production.svc/health || exit 1
```

---

## 混沌实验设计

### 实验模板

```
实验名称: [Pod/Network/Stress] 故障注入
目标: [具体服务名称]
稳态指标:
  - 错误率 < 1%
  - P99 延迟 < 500ms
  - 吞吐量 > 1000 RPS

故障类型:
  - Pod Kill: 随机终止 50% Pod
  - Network Delay: 注入 100ms 延迟
  - CPU Stress: 80% CPU 使用率

自动终止条件:
  - 错误率 > 5%
  - P99 延迟 > 2s
  - 手动干预

回滚策略:
  - 自动恢复实验
  - 检查服务健康状态
  - 通知值班人员
```

---

## 最佳实践

1. **从小规模开始**: 先在开发/测试环境验证
2. **控制爆炸半径**: 使用命名空间隔离、金丝雀发布
3. **定义清晰指标**: 明确稳态行为和终止条件
4. **自动化运行**: 将混沌实验纳入日常流程
5. **事后分析**: 记录发现的问题并修复
6. **持续改进**: 不断扩展实验范围和复杂度

---

## 与 AI 结合的混沌工程

2025 年趋势：AI 驱动的异常检测 + 混沌工程

- **智能基线**: ML 学习正常行为模式
- **预测性实验**: AI 建议最可能发现问题的实验
- **自动恢复**: 检测到异常后自动终止实验
- **根因分析**: AI 辅助分析故障传播路径
