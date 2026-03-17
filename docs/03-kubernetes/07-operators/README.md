# Kubernetes Operator 开发指南

> 扩展 Kubernetes API 以自动化复杂应用管理

---

## 什么是 Operator？

Operator 是一种扩展 Kubernetes 的方法，用于自动化复杂有状态应用的生命周期管理。它通过自定义资源定义 (CRD) 和控制器来实现。

```
┌─────────────────────────────────────────────────────────────┐
│              Operator 架构                                   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  用户                                                        │
│    │                                                         │
│    │ kubectl apply -f mydb.yaml                              │
│    ↓                                                         │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           Kubernetes API Server                       │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │  Custom Resource (CR) - 期望状态                │  │  │
│  │  │  apiVersion: myapp.io/v1                        │  │  │
│  │  │  kind: Database                                 │  │  │
│  │  │  spec:                                          │  │  │
│  │  │    replicas: 3                                  │  │  │
│  │  │    version: "15"                                │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └───────────────────────┬──────────────────────────────┘  │
│                          │ Watch                           │
│  ┌───────────────────────▼───────────────────────────────┐  │
│  │                    Controller (Operator)               │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │  Reconciliation Loop - 调谐循环                 │  │  │
│  │  │  1. 获取当前状态                                 │  │  │
│  │  │  2. 对比期望状态                                 │  │  │
│  │  │  3. 执行差异操作                                 │  │  │
│  │  │  4. 更新状态                                     │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └───────────────────────┬───────────────────────────────┘  │
│                          │ Create/Update/Delete            │
│  ┌───────────────────────▼───────────────────────────────┐  │
│  │              Kubernetes Native Resources               │  │
│  │  • Deployment  • StatefulSet  • ConfigMap             │  │
│  │  • Service     • Secret       • PVC                   │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Operator vs Controller

| 特性 | Controller | Operator |
|------|-----------|----------|
| **管理对象** | 内置资源 (Deployment, Pod) | 自定义资源 (CRD) |
| **复杂度** | 通用逻辑 | 应用特定逻辑 |
| **状态管理** | 无状态为主 | 有状态应用 |
| **领域知识** | 通用 | 应用特定 (数据库、消息队列等) |

---

## Operator 开发框架

| 框架 | 语言 | 特点 | 适用场景 |
|------|------|------|----------|
| **Operator SDK** | Go/Ansible/Helm | 功能完整，CNCF项目 | 生产级Operator |
| **Kubebuilder** | Go | 基于controller-runtime，代码生成 | Go开发者首选 |
| **KUDO** | YAML | 声明式，无需编程 | 简单Operator |
| **Metacontroller** | 任意 | 通过Webhook扩展 | 快速原型 |

---

## Kubebuilder 实战

### 安装

```bash
# 安装 kubebuilder
curl -L -o kubebuilder https://sigs.k8s.io/kubebuilder/releases/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder
sudo mv kubebuilder /usr/local/bin/
```

### 创建项目

```bash
# 初始化项目
kubebuilder init --domain myapp.io --repo github.com/myorg/database-operator

# 创建 API (CRD + Controller)
kubebuilder create api --group database --version v1 --kind Database

# 安装依赖
make manifests
make install
```

### CRD 定义

```go
// api/v1/database_types.go
package v1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatabaseSpec 定义期望状态
type DatabaseSpec struct {
    // +kubebuilder:validation:Minimum=1
    // +kubebuilder:validation:Maximum=10
    Replicas int32 `json:"replicas,omitempty"`

    // +kubebuilder:validation:Enum=13;14;15
    Version string `json:"version,omitempty"`

    Storage string `json:"storage,omitempty"`

    // +kubebuilder:default=false
    BackupEnabled bool `json:"backupEnabled,omitempty"`

    Resources ResourceRequirements `json:"resources,omitempty"`
}

// DatabaseStatus 定义观察状态
type DatabaseStatus struct {
    Phase string `json:"phase,omitempty"`
    ReadyReplicas int32 `json:"readyReplicas,omitempty"`
    CurrentVersion string `json:"currentVersion,omitempty"`
    Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
type Database struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec   DatabaseSpec   `json:"spec,omitempty"`
    Status DatabaseStatus `json:"status,omitempty"`
}
```

### Controller 实现

```go
// internal/controller/database_controller.go
package controller

import (
    "context"
    "time"

    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    databasev1 "myapp.io/database/api/v1"
)

// DatabaseReconciler 调谐器
type DatabaseReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=database.myapp.io,resources=databases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.myapp.io,resources=databases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)

    // 1. 获取 CR
    var db databasev1.Database
    if err := r.Get(ctx, req.NamespacedName, &db); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // 2. 调谐 StatefulSet
    var sts appsv1.StatefulSet
    sts.Name = db.Name
    sts.Namespace = db.Namespace

    _, err := ctrl.CreateOrUpdate(ctx, r.Client, &sts, func() error {
        // 设置 OwnerReference
        if err := ctrl.SetControllerReference(&db, &sts, r.Scheme); err != nil {
            return err
        }

        // 配置 StatefulSet
        sts.Spec.Replicas = &db.Spec.Replicas
        sts.Spec.ServiceName = db.Name
        sts.Spec.Selector = &metav1.LabelSelector{
            MatchLabels: map[string]string{"app": db.Name},
        }
        sts.Spec.Template.Labels = map[string]string{"app": db.Name}
        sts.Spec.Template.Spec.Containers = []corev1.Container{{
            Name:  "postgres",
            Image: "postgres:" + db.Spec.Version,
            Ports: []corev1.ContainerPort{{
                ContainerPort: 5432,
            }},
        }}

        return nil
    })

    if err != nil {
        return ctrl.Result{}, err
    }

    // 3. 更新状态
    db.Status.ReadyReplicas = sts.Status.ReadyReplicas
    db.Status.CurrentVersion = db.Spec.Version

    if sts.Status.ReadyReplicas == db.Spec.Replicas {
        db.Status.Phase = "Ready"
        meta.SetStatusCondition(&db.Status.Conditions, metav1.Condition{
            Type:   "Ready",
            Status: metav1.ConditionTrue,
            Reason: "AllReplicasReady",
        })
    } else {
        db.Status.Phase = "Provisioning"
        // 重新调谐
        return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
    }

    if err := r.Status().Update(ctx, &db); err != nil {
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}

func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&databasev1.Database{}).
        Owns(&appsv1.StatefulSet{}).
        Owns(&corev1.Service{}).
        Complete(r)
}
```

---

## Operator 最佳实践

### 1. 单一职责原则

```
✅ 推荐: 一个 Operator 管理一种应用
   - Postgres Operator 管理 PostgreSQL
   - Redis Operator 管理 Redis

❌ 避免: 一个 Operator 管理多种应用
   - Database Operator 同时管理 MySQL + PostgreSQL + MongoDB
```

### 2. 状态管理

```go
// 使用 finalizer 处理资源清理
const databaseFinalizer = "database.myapp.io/finalizer"

func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    var db databasev1.Database
    if err := r.Get(ctx, req.NamespacedName, &db); err != nil {
        return ctrl.Result{}, err
    }

    // 处理删除
    if !db.DeletionTimestamp.IsZero() {
        if controllerutil.ContainsFinalizer(&db, databaseFinalizer) {
            // 执行清理操作
            if err := r.cleanupDatabase(ctx, &db); err != nil {
                return ctrl.Result{}, err
            }
            controllerutil.RemoveFinalizer(&db, databaseFinalizer)
            return ctrl.Result{}, r.Update(ctx, &db)
        }
        return ctrl.Result{}, nil
    }

    // 添加 finalizer
    if !controllerutil.ContainsFinalizer(&db, databaseFinalizer) {
        controllerutil.AddFinalizer(&db, databaseFinalizer)
        return ctrl.Result{}, r.Update(ctx, &db)
    }

    // 正常调谐
    return r.reconcileNormal(ctx, &db)
}
```

### 3. 高可用性

```yaml
# 多副本部署 Operator
apiVersion: apps/v1
kind: Deployment
metadata:
  name: database-operator
spec:
  replicas: 2
  template:
    spec:
      containers:
        - name: manager
          image: myorg/database-operator:v1.0.0
          args:
            - --leader-elect
          resources:
            limits:
              cpu: 500m
              memory: 256Mi
```

---

## 常用 Operator 推荐

| Operator | 用途 | 成熟度 |
|----------|------|--------|
| **Prometheus Operator** | 监控 | 生产级 |
| **Postgres Operator (Zalando)** | PostgreSQL | 生产级 |
| **Strimzi** | Kafka | 生产级 |
| **Argo CD Operator** | GitOps | 生产级 |
| **cert-manager** | 证书管理 | 生产级 |
| **Velero** | 备份 | 生产级 |

---

## 2025 趋势

- **Operator Lifecycle Manager (OLM)**: 标准化 Operator 分发
- **Cluster API**: 基础设施即 Kubernetes 对象
- **AI/ML Operators**: 自动化 ML 工作流
