# Go 代码实战指南

> 使用 client-go 与 Kubernetes API 交互

---

## 目录

1. [环境准备](#1-环境准备)
2. [基础操作](#2-基础操作)
3. [Informer 模式](#3-informer-模式)
4. [自定义控制器](#4-自定义控制器)
5. [关联代码示例](#5-关联代码示例)

---

## 1. 环境准备

### 1.1 依赖安装

```bash
go get k8s.io/client-go@v0.29.0
```

### 1.2 配置连接

```go
import (
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    "path/filepath"
)

func NewK8sClient(kubeconfig string) (*kubernetes.Clientset, error) {
    if kubeconfig == "" {
        kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
    }

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        // 尝试 in-cluster 配置
        config, err = clientcmd.NewInClusterConfig()
        if err != nil {
            return nil, err
        }
    }

    return kubernetes.NewForConfig(config)
}
```

---

## 2. 基础操作

### 2.1 CRUD 操作

```go
import (
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "context"
)

// 创建 Pod
func (c *Client) CreatePod(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
    return c.clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
}

// 获取 Pod
func (c *Client) GetPod(namespace, name string) (*corev1.Pod, error) {
    return c.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// 列出 Pod
func (c *Client) ListPods(namespace string) (*corev1.PodList, error) {
    return c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

// 删除 Pod
func (c *Client) DeletePod(namespace, name string) error {
    return c.clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
```

### 2.2 Watch 监视

```go
func (c *Client) WatchPods(namespace string) {
    watch, err := c.clientset.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatal(err)
    }
    defer watch.Stop()

    for event := range watch.ResultChan() {
        pod := event.Object.(*corev1.Pod)
        log.Printf("Event: %s, Pod: %s/%s", event.Type, pod.Namespace, pod.Name)
    }
}
```

---

## 3. Informer 模式

### 3.1 基本用法

```go
import (
    "k8s.io/client-go/informers"
    "k8s.io/client-go/tools/cache"
    "time"
)

func RunInformer(clientset *kubernetes.Clientset) {
    // 创建 SharedInformerFactory
    factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

    // 获取 Pod Informer
    podInformer := factory.Core().V1().Pods().Informer()

    // 添加事件处理器
    podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            pod := obj.(*corev1.Pod)
            log.Printf("Pod Added: %s/%s", pod.Namespace, pod.Name)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            oldPod := oldObj.(*corev1.Pod)
            newPod := newObj.(*corev1.Pod)
            if oldPod.Status.Phase != newPod.Status.Phase {
                log.Printf("Pod %s status: %s -> %s",
                    newPod.Name, oldPod.Status.Phase, newPod.Status.Phase)
            }
        },
        DeleteFunc: func(obj interface{}) {
            pod := obj.(*corev1.Pod)
            log.Printf("Pod Deleted: %s/%s", pod.Namespace, pod.Name)
        },
    })

    // 启动 Informer
    stopCh := make(chan struct{})
    factory.Start(stopCh)

    // 等待缓存同步
    if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {
        log.Fatal("Failed to sync cache")
    }

    <-stopCh
}
```

### 3.2 从缓存读取

```go
// 使用 Lister 从本地缓存读取（不访问 API Server）
lister := factory.Core().V1().Pods().Lister()
pods, err := lister.List(labels.Everything())
pod, err := lister.Pods("default").Get("mypod")
```

---

## 4. 自定义控制器

### 4.1 控制器结构

```go
type Controller struct {
    clientset *kubernetes.Clientset
    indexer   cache.Indexer
    queue     workqueue.RateLimitingInterface
    informer  cache.Controller
}

func NewController(clientset *kubernetes.Clientset) *Controller {
    listWatcher := cache.NewListWatchFromClient(
        clientset.CoreV1().RESTClient(),
        "pods",
        corev1.NamespaceAll,
        fields.Everything(),
    )

    queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

    indexer, informer := cache.NewIndexerInformer(
        listWatcher,
        &corev1.Pod{},
        0,
        cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                key, err := cache.MetaNamespaceKeyFunc(obj)
                if err == nil {
                    queue.Add(key)
                }
            },
            UpdateFunc: func(oldObj, newObj interface{}) {
                key, err := cache.MetaNamespaceKeyFunc(newObj)
                if err == nil {
                    queue.Add(key)
                }
            },
            DeleteFunc: func(obj interface{}) {
                key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
                if err == nil {
                    queue.Add(key)
                }
            },
        },
        cache.Indexers{},
    )

    return &Controller{
        clientset: clientset,
        indexer:   indexer,
        queue:     queue,
        informer:  informer,
    }
}
```

### 4.2 控制器主循环

```go
func (c *Controller) Run(workers int, stopCh chan struct{}) {
    defer runtime.HandleCrash()
    defer c.queue.ShutDown()

    go c.informer.Run(stopCh)

    if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
        runtime.HandleError(fmt.Errorf("cache sync failed"))
        return
    }

    for i := 0; i < workers; i++ {
        go wait.Until(c.runWorker, time.Second, stopCh)
    }

    <-stopCh
}

func (c *Controller) runWorker() {
    for c.processNextItem() {
    }
}

func (c *Controller) processNextItem() bool {
    key, quit := c.queue.Get()
    if quit {
        return false
    }
    defer c.queue.Done(key)

    err := c.syncHandler(key.(string))
    if err == nil {
        c.queue.Forget(key)
    } else if c.queue.NumRekeys(key) < 5 {
        c.queue.AddRateLimited(key)
    } else {
        c.queue.Forget(key)
        runtime.HandleError(err)
    }

    return true
}

func (c *Controller) syncHandler(key string) error {
    obj, exists, err := c.indexer.GetByKey(key)
    if err != nil {
        return err
    }
    if !exists {
        return nil
    }

    pod := obj.(*corev1.Pod)
    // 执行业务逻辑
    log.Printf("Processing pod: %s/%s", pod.Namespace, pod.Name)
    return nil
}
```

---

## 5. 关联代码示例

| 主题 | 代码位置 |
|------|----------|
| 基础 CRUD | `examples/go-client/01-basic-ops/main.go` |
| Informer | `examples/go-client/03-informer/main.go` |
| 自定义控制器 | `examples/go-client/02-controller/main.go` |
| 微服务 | `examples/microservices-demo/user-service/main.go` |

---

## 参考

- [client-go 官方示例](https://github.com/kubernetes/client-go/tree/master/examples)
- [Kubernetes 控制器模式](https://kubernetes.io/docs/concepts/architecture/controller/)
