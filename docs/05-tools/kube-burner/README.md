# Kube-burner - Kubernetes 性能测试工具

## 概述

Kube-burner 是一个 CNCF Sandbox 项目，用于 Kubernetes 性能和规模测试编排。它可以大规模创建、删除、读取和修补 Kubernetes 资源，并收集 Prometheus 指标。

## 核心特性

| 特性 | 描述 |
|------|------|
| 资源操作 | 大规模创建、删除、读取、修补 K8s 资源 |
| 指标收集 | 从 Prometheus 收集性能指标 |
| 测量功能 | Pod 延迟、VMI 延迟、Pprof 分析 |
| 告警检测 | 集成 Prometheus 告警 |
| OCP 包装器 | 针对 OpenShift 的专用工作负载 |

## 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                     Kube-burner 架构                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │   Config     │    │   Job        │    │   Measure    │      │
│  │   (YAML)     │───▶│   Execution  │───▶│   & Index    │      │
│  └──────────────┘    └──────────────┘    └──────┬───────┘      │
│         │                                       │              │
│         │         ┌─────────────────────────────┘              │
│         │         ▼                                            │
│         │    ┌──────────────┐    ┌──────────────┐             │
│         │    │  Prometheus  │    │ Elasticsearch│             │
│         │    │  (Metrics)   │    │  (Results)   │             │
│         │    └──────────────┘    └──────────────┘             │
│         │                                                      │
│         ▼                                                      │
│  ┌────────────────────────────────────────────────────────┐   │
│  │                   Kubernetes API                        │   │
│  │  - 创建/删除/更新资源                                   │   │
│  │  - Pod/Deployment/Service/ConfigMap/etc                 │   │
│  │  - 测量 Pod 启动延迟                                    │   │
│  └────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 安装部署

### 二进制安装

```bash
# 下载最新版本
curl -L https://github.com/kube-burner/kube-burner/releases/latest/download/kube-burner-linux-amd64.tar.gz | tar xz
sudo mv kube-burner /usr/local/bin/

# 验证安装
kube-burner version
```

### 容器运行

```bash
# 使用容器镜像
docker run --rm -v $(pwd)/config:/config quay.io/kube-burner/kube-burner:latest init -c /config/config.yaml
```

## 配置示例

```yaml
# config.yaml
---
global:
  qps: 20
  burst: 20
  gc: true
  gcMetrics: true
  measurements:
  - name: podLatency
    esIndex: kube-burner-pod
  - name: pprof
    esIndex: kube-burner-pprof
    pprofInterval: 5m
  indexerConfig:
    enabled: true
    type: opensearch
    esServers: ["https://opensearch:9200"]
    defaultIndex: kube-burner

jobs:
# 第一阶段：创建 namespace
- name: create-namespaces
  jobType: create
  qps: 10
  burst: 10
  jobIterations: 10
  namespace: kube-burner
  objects:
  - kind: Namespace
    objectTemplate: templates/namespace.yaml
    replicas: 1

# 第二阶段：部署应用
- name: deploy-applications
  jobType: create
  qps: 20
  burst: 20
  jobIterations: 100
  namespace: kube-burner
  objects:
  - kind: Deployment
    objectTemplate: templates/deployment.yaml
    replicas: 1
    inputVars:
      containerImage: nginx:latest
      replicas: 3
  - kind: Service
    objectTemplate: templates/service.yaml
    replicas: 1
  - kind: ConfigMap
    objectTemplate: templates/configmap.yaml
    replicas: 5

# 第三阶段：清理
- name: cleanup
  jobType: delete
  qps: 10
  burst: 10
  objects:
  - kind: Namespace
    labelSelector: {kube-burner-job: create-namespaces}
```

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.JobName}}-{{.Iteration}}
spec:
  replicas: {{.replicas}}
  selector:
    matchLabels:
      app: {{.JobName}}-{{.Iteration}}
  template:
    metadata:
      labels:
        app: {{.JobName}}-{{.Iteration}}
    spec:
      containers:
      - name: app
        image: {{.containerImage}}
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
```

## OpenShift OCP 包装器

```bash
# 查看可用工作负载
kube-burner ocp help

# Node Density 测试 - 100 pods 每节点
kube-burner ocp node-density --pods-per-node=100

# Cluster Density 测试
kube-burner ocp cluster-density-v2 --iterations=1000

# 带索引和告警检测
kube-burner ocp node-density --pods-per-node=100 \
  --es-server=https://opensearch:9200 \
  --es-index=kube-burner \
  --alerting=true
```

## 测量类型

### Pod 延迟测量

```yaml
measurements:
- name: podLatency
  esIndex: kube-burner-pod
  thresholds:
  - conditionType: Ready
    metric: P99
    threshold: 10s
```

### Pprof 分析

```yaml
measurements:
- name: pprof
  esIndex: kube-burner-pprof
  pprofInterval: 5m
  pprofTargets:
  - name: kube-apiserver
    namespace: openshift-kube-apiserver
    labelSelector: {app: openshift-kube-apiserver}
    container: kube-apiserver
  - name: etcd
    namespace: openshift-etcd
    labelSelector: {app: etcd}
    container: etcd
```

## 指标收集

```yaml
metricsEndpoints:
- endpoint: http://localhost:9090
  token: <prometheus-token>
  profile: metrics.yaml
  alertProfile: alert-profile.yaml

# 自定义指标配置
- endpoint: http://thanos-querier:9090
  token: <token>
  profile: metrics-reporting.yaml
```

```yaml
# metrics.yaml - 自定义指标查询
- query: avg(avg_over_time(irate(container_cpu_usage_seconds_total{name!="",namespace=~"openshift-.+|kube-.+|default"}[2m])[{{.elapsed}}:]))
  metricName: containerCPU

- query: avg(avg_over_time(container_memory_rss{name!="",namespace=~"openshift-.+|kube-.+|default"}[{{.elapsed}}:]))
  metricName: containerMemory

- query: histogram_quantile(0.99, sum(rate(apiserver_request_duration_seconds_bucket[5m])) by (le))
  metricName: apiserverRequestLatency
```

## 持续性能测试 (CPT)

```yaml
# .github/workflows/performance.yml
name: Continuous Performance Testing

on:
  push:
    branches: [main]

jobs:
  performance-test:
    runs-on: self-hosted
    steps:
    - uses: actions/checkout@v4

    - name: Run kube-burner
      run: |
        kube-burner ocp cluster-density-v2 \
          --iterations=100 \
          --es-server=${{ secrets.ES_SERVER }} \
          --es-index=ci-performance \
          --uuid=${{ github.run_id }}

    - name: Check regressions with Orion
      run: |
        orion cli --config=orion.yaml \
          --uuid=${{ github.run_id }} \
          --baseline=main
```

## 性能基线

| 指标 | 可接受范围 | 告警阈值 |
|------|-----------|----------|
| Pod 启动 P99 | < 30s | > 60s |
| API Server CPU | < 50% | > 80% |
| etcd 延迟 P99 | < 10ms | > 50ms |
| 调度延迟 P99 | < 5s | > 15s |

## 相关资源

- [Kube-burner 文档](https://kube-burner.github.io/kube-burner/)
- [GitHub](https://github.com/kube-burner/kube-burner)
- [CNCF Sandbox](https://landscape.cncf.io/?selected=kube-burner)
