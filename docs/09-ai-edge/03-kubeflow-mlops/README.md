# Kubeflow MLOps - 生产级 AI/ML 平台

## 概述

Kubeflow 是 Kubernetes 上的开源 MLOps 平台，覆盖 ML 生命周期全流程：数据准备、模型训练、超参调优、模型服务和监控。2025年，Kubeflow 1.10/1.11 支持 K8s 1.33+，服务企业级 AI 部署。

> **关键数据**: 2025年，78% 的企业 ML 部署运行在 Kubernetes 上，Kubeflow 是首选平台。

## 架构组件

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Kubeflow Platform                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │   Notebooks  │  │  Pipelines   │  │    Katib     │  │   Training   │    │
│  │  (Jupyter)   │  │  (KFP v2)    │  │ (AutoML/HPO) │  │  Operators   │    │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘    │
│         │                 │                 │                 │            │
│         └─────────────────┴─────────────────┴─────────────────┘            │
│                                     │                                       │
│  ┌──────────────────────────────────┼──────────────────────────────────┐   │
│  │                        Model Registry                               │   │
│  │                     (MLMD/MR v2.0)                                  │   │
│  └──────────────────────────────────┼──────────────────────────────────┘   │
│                                     │                                       │
│  ┌──────────────────────────────────┼──────────────────────────────────┐   │
│  │                        KServe                                       │   │
│  │              (Model Serving Platform)                               │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐              │   │
│  │  │Inference │ │  Auto    │ │  Model   │ │  Multi   │              │   │
│  │  │ Services │ │ Scaling  │ │  Canary  │ │ Framework│              │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘              │   │
│  └────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌────────────────────────────────────────────────────────────────────┐   │
│  │                    Supporting Components                            │   │
│  │  Spark Operator | Feast (Feature Store) | MLflow | Prometheus       │   │
│  └────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Kubeflow 1.11 新特性 (2025)

| 特性 | 说明 |
|------|------|
| **Trainer 2.0** | 支持 JAX 分布式训练 |
| **Spark Operator** | 成为核心组件，支持 YuniKorn 群体调度 |
| **模型注册表 UI** | 全新界面，简化模型管理 |
| **S3 多租户** | SeaweedFS 替代 MinIO，硬多租户支持 |
| **LLMOps** | Katib 新增 LLM 超参优化 API |
| **KServe 0.16** | Python SDK、OCI 模型存储、模型缓存 |
| **安全增强** | Pod Security Standards Restricted 默认启用 |

## 安装 Kubeflow

### 方式1：Manifest（生产推荐）

```bash
# K8s 1.33+ 环境
kubectl apply -k "github.com/kubeflow/manifests.git?ref=v1.11.0"

# 等待就绪
kubectl wait --for=condition=available deployment/ml-pipeline -n kubeflow --timeout=300s
```

### 方式2：Terraform（IaC）

```hcl
module "kubeflow" {
  source = "terraform-google-modules/kubernetes-engine/google//modules/kubeflow"

  project_id = var.project_id
  cluster_name = "ml-cluster"
  kubernetes_version = "1.33.2"

  # 组件选择
  enable_katib = true
  enable_kserve = true
  enable_pipelines = true
  enable_notebooks = true
}
```

### 方式3：Helm（实验性）

```bash
# Kubeflow 1.11+ 开始支持 Helm
helm repo add kubeflow https://kubeflow.github.io/manifests
helm install kubeflow kubeflow/kubeflow \
  --namespace kubeflow \
  --create-namespace \
  --set pipelines.enabled=true \
  --set kserve.enabled=true
```

## Kubeflow Pipelines (KFP)

### Pipeline SDK v2

```python
# pipeline.py
from kfp import dsl
from kfp.dsl import Input, Output, Dataset, Model

@dsl.component(
    base_image="python:3.10",
    packages_to_install=["pandas", "scikit-learn"]
)
def preprocess(
    input_data: Input[Dataset],
    output_data: Output[Dataset]
):
    import pandas as pd
    from sklearn.preprocessing import StandardScaler

    df = pd.read_csv(input_data.path)
    scaler = StandardScaler()
    df_scaled = scaler.fit_transform(df)

    pd.DataFrame(df_scaled).to_csv(output_data.path, index=False)

@dsl.component(
    base_image="python:3.10",
    packages_to_install=["scikit-learn", "joblib"]
)
def train(
    training_data: Input[Dataset],
    model: Output[Model],
    accuracy: float
):
    import joblib
    from sklearn.ensemble import RandomForestClassifier
    from sklearn.model_selection import train_test_split
    import pandas as pd

    df = pd.read_csv(training_data.path)
    X, y = df.drop('target', axis=1), df['target']
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2)

    clf = RandomForestClassifier(n_estimators=100)
    clf.fit(X_train, y_train)

    accuracy = clf.score(X_test, y_test)
    joblib.dump(clf, model.path)

@dsl.pipeline(
    name="ML Training Pipeline",
    description="End-to-end ML training"
)
def ml_pipeline(
    data_url: str = "gs://bucket/dataset.csv"
):
    # 数据预处理
    preprocess_task = preprocess(
        input_data=dsl.importer(
            artifact_uri=data_url,
            artifact_class=Dataset
        ).output
    )

    # 模型训练
    train_task = train(
        training_data=preprocess_task.outputs['output_data']
    )

    # 条件判断
    with dsl.Condition(train_task.outputs['accuracy'] > 0.8):
        deploy_task = deploy_model(
            model=train_task.outputs['model']
        )
```

### 编译和运行

```bash
# 编译 pipeline
kfp dsl compile --py pipeline.py --output pipeline.yaml

# 上传到 Kubeflow
kfp pipeline create --pipeline-name ml-pipeline pipeline.yaml

# 或通过 UI 创建运行
```

## Katib 超参优化

### 自动超参调优

```yaml
apiVersion: kubeflow.org/v1beta1
kind: Experiment
metadata:
  name: hyperparam-tuning
  namespace: kubeflow
spec:
  parallelTrialCount: 3
  maxTrialCount: 12
  maxFailedTrialCount: 3
  objective:
    type: maximize
    goal: 0.99
    objectiveMetricName: accuracy
  algorithm:
    algorithmName: bayesianoptimization
  parameters:
  - name: learning_rate
    parameterType: double
    feasibleSpace:
      min: "0.01"
      max: "0.1"
  - name: batch_size
    parameterType: categorical
    feasibleSpace:
      list:
      - "32"
      - "64"
      - "128"
  - name: optimizer
    parameterType: categorical
    feasibleSpace:
      list:
      - adam
      - sgd
  trialTemplate:
    primaryContainerName: training-container
    trialParameters:
    - name: learningRate
      description: Learning rate
      reference: learning_rate
    - name: batchSize
      description: Batch size
      reference: batch_size
    - name: optimizer
      description: Optimizer
      reference: optimizer
    trialSpec:
      apiVersion: batch/v1
      kind: Job
      spec:
        template:
          spec:
            containers:
            - name: training-container
              image: myregistry/ml-training:v1
              command:
              - python
              - train.py
              - --lr=${trialParameters.learningRate}
              - --batch-size=${trialParameters.batchSize}
              - --optimizer=${trialParameters.optimizer}
            restartPolicy: Never
```

### LLMOps 超参优化 (Kubeflow 1.11+)

```yaml
apiVersion: kubeflow.org/v1beta1
kind: Experiment
metadata:
  name: llm-fine-tuning
spec:
  objective:
    type: minimize
    objectiveMetricName: perplexity
  algorithm:
    algorithmName: tpe  # 针对大模型优化的 TPE
  parameters:
  - name: lora_r
    parameterType: int
    feasibleSpace:
      min: "8"
      max: "64"
  - name: lora_alpha
    parameterType: int
    feasibleSpace:
      min: "16"
      max: "128"
  - name: quantization_bits
    parameterType: categorical
    feasibleSpace:
      list: ["4", "8"]
```

## KServe 模型服务

### InferenceService

```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: sklearn-iris
  namespace: production
spec:
  predictor:
    # 自动扩缩容
    minReplicas: 1
    maxReplicas: 10

    # 模型存储
    storageUri: "gs://kfserving-examples/models/sklearn/1.0/model"

    # 运行时
    sklearn:
      protocolVersion: v1
      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 1000m
          memory: 1Gi

  # 变换器（预处理/后处理）
  transformer:
    containers:
    - name: transformer
      image: myregistry/iris-transformer:v1
      resources:
        requests:
          cpu: 100m
          memory: 256Mi

  # 解释器
  explainer:
    alibi:
      type: KernelShap
      storageUri: "gs://kfserving-examples/models/sklearn/1.0/explainer"
```

### GPU 推理

```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: llm-gpt
spec:
  predictor:
    model:
      modelFormat:
        name: huggingface
      storageUri: "gs://models/llama-3-8b"
      resources:
        limits:
          nvidia.com/gpu: 2
        requests:
          nvidia.com/gpu: 2
    # GPU 特定配置
    nodeSelector:
      node-type: gpu-a100
    tolerations:
    - key: nvidia.com/gpu
      operator: Exists
      effect: NoSchedule
```

### 金丝雀发布

```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: sentiment-analysis
spec:
  predictor:
    canaryTrafficPercent: 20  # 20% 流量到新版本
    model:
      modelFormat:
        name: sklearn
      storageUri: "gs://models/sentiment/v2"
      name: sentiment-v2
```

## Training Operator 分布式训练

### PyTorchJob

```yaml
apiVersion: kubeflow.org/v1
kind: PyTorchJob
metadata:
  name: pytorch-distributed
spec:
  pytorchReplicaSpecs:
    Master:
      replicas: 1
      restartPolicy: OnFailure
      template:
        spec:
          containers:
          - name: pytorch
            image: myregistry/pytorch-training:v1
            command:
            - python
            - train.py
            - --backend=nccl
            resources:
              limits:
                nvidia.com/gpu: 4
    Worker:
      replicas: 3
      restartPolicy: OnFailure
      template:
        spec:
          containers:
          - name: pytorch
            image: myregistry/pytorch-training:v1
            command:
            - python
            - train.py
            - --backend=nccl
            resources:
              limits:
                nvidia.com/gpu: 4
```

### JAX 训练 (Kubeflow 1.11+)

```yaml
apiVersion: kubeflow.org/v1
kind: TrainingJob
metadata:
  name: jax-flax-training
spec:
  runtimeRef:
    apiVersion: trainer.kubeflow.org/v1
    kind: TrainingRuntime
    name: jax-distributed
  trainer:
    numNodes: 4
    resources:
      limits:
        nvidia.com/gpu: 8
    image: myregistry/jax-training:v1
```

## 完整 MLOps 流水线

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    End-to-End MLOps Pipeline                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. 数据准备                                                                 │
│     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐                │
│     │ Data Ingest │ ──► │   Spark     │ ──► │  Feature    │                │
│     │  (Raw Data) │     │ Processing  │     │   Store     │                │
│     └─────────────┘     └─────────────┘     └─────────────┘                │
│                                                        │                     │
│  2. 模型训练                                            │                     │
│     ┌─────────────┐     ┌─────────────┐     ┌─────────▼─────┐              │
│     │   Katib     │ ──► │ Distributed │ ──► │   Feast       │              │
│     │    (HPO)    │     │  Training   │     │ (Get Features)│              │
│     └─────────────┘     └─────────────┘     └───────────────┘              │
│                                                        │                     │
│  3. 模型评估                                            ▼                     │
│     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐                │
│     │   Model     │ ──► │ Validation  │ ──► │   Register  │                │
│     │   Evaluate  │     │   Tests     │     │   (MLflow)  │                │
│     └─────────────┘     └─────────────┘     └──────┬──────┘                │
│                                                    │                        │
│  4. 模型部署                                        ▼                        │
│     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐                │
│     │   KServe    │ ◄── │   Model     │ ◄── │  Drift      │                │
│     │  (Serving)  │     │   Canary    │     │  Detection  │                │
│     └─────────────┘     └─────────────┘     └─────────────┘                │
│           │                                                                  │
│  5. 监控   │                                                                  │
│     ┌─────▼─────┐     ┌─────────────┐     ┌─────────────┐                  │
│     │ Prometheus│ ──► │   Grafana   │ ──► │   Alert     │                  │
│     │  Metrics  │     │ Dashboards  │     │   Retrain   │ ──► 回到步骤 2   │
│     └───────────┘     └─────────────┘     └─────────────┘                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 生产最佳实践

### 资源隔离

```yaml
# 使用 ResourceQuota 限制 ML 工作负载
apiVersion: v1
kind: ResourceQuota
metadata:
  name: ml-quota
  namespace: data-science
spec:
  hard:
    requests.nvidia.com/gpu: 16
    limits.nvidia.com/gpu: 16
    requests.memory: 512Gi
    requests.cpu: "64"
```

### 多租户安全

```yaml
# Pod Security Standards
apiVersion: v1
kind: Namespace
metadata:
  name: data-science
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

### 成本优化

```yaml
# 使用抢占式实例训练
apiVersion: kubeflow.org/v1
kind: PyTorchJob
metadata:
  name: spot-training
spec:
  pytorchReplicaSpecs:
    Worker:
      template:
        spec:
          nodeSelector:
            cloud.google.com/gke-spot: "true"
          tolerations:
          - key: cloud.google.com/gke-spot
            operator: Equal
            value: "true"
            effect: NoSchedule
          # 检查点配置
          containers:
          - name: pytorch
            volumeMounts:
            - name: checkpoint
              mountPath: /checkpoints
          volumes:
          - name: checkpoint
            persistentVolumeClaim:
              claimName: training-checkpoint-pvc
```

## 监控指标

| 指标 | 工具 | 告警条件 |
|------|------|----------|
| GPU 利用率 | Prometheus | < 30% 持续 10 分钟 |
| 训练损失 | MLflow | 3 epoch 不下降 |
| 推理延迟 | KServe | P99 > 100ms |
| 模型漂移 | Evidently | KS 统计量 > 阈值 |
| Pipeline 成功率 | KFP | < 95% |

## 总结

| 场景 | Kubeflow 组件 |
|------|--------------|
| 交互式开发 | Notebooks |
| 自动化训练 | Pipelines + Training Operator |
| 超参优化 | Katib |
| 模型服务 | KServe |
| 特征管理 | Feast |
| 实验跟踪 | MLflow |

Kubeflow 提供完整的 MLOps 解决方案，是 2025 年企业级 AI 平台的标配。
