package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

// Controller 演示一个简单的自定义控制器
type Controller struct {
	clientset    *kubernetes.Clientset
	indexer      cache.Indexer
	queue        workqueue.TypedRateLimitingInterface[string]
	informer     cache.Controller
	stopCh       chan struct{}
}

// NewController 创建新的控制器
func NewController(clientset *kubernetes.Clientset) *Controller {
	// 创建 ListWatcher
	listWatcher := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"pods",
		corev1.NamespaceAll,
		fields.Everything(),
	)

	// 创建工作队列
	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())

	// 创建 Indexer 和 Informer
	indexer, informer := cache.NewIndexerInformer(
		listWatcher,
		&corev1.Pod{},
		0, // 全量同步
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err == nil {
					queue.Add(key)
					klog.Infof("Pod added: %s", key)
				}
			},
			UpdateFunc: func(old interface{}, new interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(new)
				if err == nil {
					queue.Add(key)
					klog.Infof("Pod updated: %s", key)
				}
			},
			DeleteFunc: func(obj interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err == nil {
					queue.Add(key)
					klog.Infof("Pod deleted: %s", key)
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
		stopCh:    make(chan struct{}),
	}
}

// Run 启动控制器
func (c *Controller) Run(workers int) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Info("Starting controller...")

	// 启动 Informer
	go c.informer.Run(c.stopCh)

	// 等待缓存同步
	if !cache.WaitForCacheSync(c.stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	klog.Info("Caches synced, starting workers...")

	// 启动工作线程
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, c.stopCh)
	}

	<-c.stopCh
	klog.Info("Stopping controller...")
}

// Stop 停止控制器
func (c *Controller) Stop() {
	close(c.stopCh)
}

// runWorker 工作线程
func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

// processNextItem 处理队列中的下一个项目
func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncHandler(key)
	if err == nil {
		c.queue.Forget(key)
	} else if c.queue.NumRequeues(key) < 5 {
		klog.Errorf("Error syncing pod %v: %v", key, err)
		c.queue.AddRateLimited(key)
	} else {
		klog.Errorf("Dropping pod %q out of the queue: %v", key, err)
		c.queue.Forget(key)
		runtime.HandleError(err)
	}

	return true
}

// syncHandler 实际的业务逻辑
func (c *Controller) syncHandler(key string) error {
	// 从 key 解析 namespace 和 name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// 从 indexer 获取 Pod
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		klog.Infof("Pod %s/%s does not exist anymore\n", namespace, name)
		return nil
	}

	pod := obj.(*corev1.Pod)

	// 业务逻辑：检查 Pod 状态并执行操作
	klog.Infof("Processing pod: %s/%s, Phase: %s", 
		pod.Namespace, pod.Name, pod.Status.Phase)

	// 示例：如果 Pod 处于 Pending 状态超过 5 分钟，添加注解
	if pod.Status.Phase == corev1.PodPending {
		if pod.Status.StartTime != nil {
			pendingTime := time.Since(pod.Status.StartTime.Time)
			if pendingTime > 5*time.Minute {
				if err := c.addAnnotation(pod, "debug/pending-since", 
					pod.Status.StartTime.Format(time.RFC3339)); err != nil {
					return err
				}
			}
		}
	}

	// 示例：确保 Pod 有特定的标签
	if _, exists := pod.Labels["managed-by"]; !exists {
		if err := c.addLabel(pod, "managed-by", "custom-controller"); err != nil {
			return err
		}
	}

	return nil
}

// addAnnotation 添加注解到 Pod
func (c *Controller) addAnnotation(pod *corev1.Pod, key, value string) error {
	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}
	pod.Annotations[key] = value

	_, err := c.clientset.CoreV1().Pods(pod.Namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		if apierrors.IsConflict(err) {
			klog.Warningf("Conflict updating pod %s/%s, will retry", pod.Namespace, pod.Name)
			return nil
		}
		return err
	}
	klog.Infof("Added annotation %s=%s to pod %s/%s", key, value, pod.Namespace, pod.Name)
	return nil
}

// addLabel 添加标签到 Pod
func (c *Controller) addLabel(pod *corev1.Pod, key, value string) error {
	if pod.Labels == nil {
		pod.Labels = make(map[string]string)
	}
	pod.Labels[key] = value

	_, err := c.clientset.CoreV1().Pods(pod.Namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		if apierrors.IsConflict(err) {
			klog.Warningf("Conflict updating pod %s/%s, will retry", pod.Namespace, pod.Name)
			return nil
		}
		return err
	}
	klog.Infof("Added label %s=%s to pod %s/%s", key, value, pod.Namespace, pod.Name)
	return nil
}

func main() {
	var kubeconfig string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	flag.Parse()

	// 加载配置
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		config, err = clientcmd.NewInClusterConfig()
		if err != nil {
			klog.Fatalf("Error building kubeconfig: %s", err.Error())
		}
	}

	// 创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	// 创建并运行控制器
	controller := NewController(clientset)
	
	// 优雅退出
	stop := make(chan os.Signal, 1)
	go func() {
		<-stop
		controller.Stop()
	}()

	controller.Run(2) // 使用 2 个工作线程
}
