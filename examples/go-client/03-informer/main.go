package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

// InformerDemo 演示 Informer 的使用
type InformerDemo struct {
	clientset *kubernetes.Clientset
	factory   informers.SharedInformerFactory
	stopCh    chan struct{}
}

// NewInformerDemo 创建新的 InformerDemo
func NewInformerDemo(kubeconfig string) (*InformerDemo, error) {
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = homedir.HomeDir() + "/.kube/config"
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		config, err = clientcmd.NewInClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// 创建 SharedInformerFactory
	// 第二个参数是 resyncPeriod，定期重新同步
	factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	return &InformerDemo{
		clientset: clientset,
		factory:   factory,
		stopCh:    make(chan struct{}),
	}, nil
}

// SetupPodInformer 设置 Pod Informer
func (d *InformerDemo) SetupPodInformer() {
	// 获取 Pod Informer
	podInformer := d.factory.Core().V1().Pods().Informer()

	// 添加事件处理器
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			klog.Infof("Pod Added: %s/%s, Phase: %s, Node: %s",
				pod.Namespace, pod.Name, pod.Status.Phase, pod.Spec.NodeName)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*corev1.Pod)
			newPod := newObj.(*corev1.Pod)

			// 只处理状态变化的 Pod
			if oldPod.Status.Phase != newPod.Status.Phase {
				klog.Infof("Pod Status Changed: %s/%s, %s -> %s",
					newPod.Namespace, newPod.Name, oldPod.Status.Phase, newPod.Status.Phase)
			}

			// 监控重启次数变化
			oldRestarts := getContainerRestarts(oldPod)
			newRestarts := getContainerRestarts(newPod)
			if newRestarts > oldRestarts {
				klog.Warningf("Container Restarted: %s/%s, Restarts: %d",
					newPod.Namespace, newPod.Name, newRestarts)
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			klog.Infof("Pod Deleted: %s/%s", pod.Namespace, pod.Name)
		},
	})
}

// SetupNodeInformer 设置 Node Informer
func (d *InformerDemo) SetupNodeInformer() {
	nodeInformer := d.factory.Core().V1().Nodes().Informer()

	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node := obj.(*corev1.Node)
			klog.Infof("Node Added: %s, Ready: %v", node.Name, isNodeReady(node))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldNode := oldObj.(*corev1.Node)
			newNode := newObj.(*corev1.Node)

			oldReady := isNodeReady(oldNode)
			newReady := isNodeReady(newNode)

			if oldReady != newReady {
				if newReady {
					klog.Infof("Node Ready: %s", newNode.Name)
				} else {
					klog.Warningf("Node NotReady: %s", newNode.Name)
				}
			}
		},
	})
}

// SetupEventInformer 设置 Event Informer（集群事件）
func (d *InformerDemo) SetupEventInformer() {
	eventInformer := d.factory.Core().V1().Events().Informer()

	eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event := obj.(*corev1.Event)
			// 只关注 Warning 级别的事件
			if event.Type == corev1.EventTypeWarning {
				klog.Warningf("Event: %s/%s, Reason: %s, Message: %s",
					event.InvolvedObject.Namespace,
					event.InvolvedObject.Name,
					event.Reason,
					event.Message)
			}
		},
	})
}

// Run 启动 Informer
func (d *InformerDemo) Run() {
	// 启动所有 Informer
	d.factory.Start(d.stopCh)

	// 等待缓存同步
	klog.Info("Waiting for cache sync...")
	if !cache.WaitForCacheSync(d.stopCh, d.factory.Core().V1().Pods().Informer().HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	klog.Info("Cache synced, watching for events...")

	// 阻塞直到收到停止信号
	<-d.stopCh
}

// Stop 停止 Informer
func (d *InformerDemo) Stop() {
	close(d.stopCh)
}

// ListPodsFromCache 从缓存中列出 Pod（不访问 API Server）
func (d *InformerDemo) ListPodsFromCache(namespace string) {
	lister := d.factory.Core().V1().Pods().Lister()

	var pods []*corev1.Pod
	var err error

	if namespace == "" {
		pods, err = lister.List(cache.Everything)
	} else {
		pods, err = lister.Pods(namespace).List(cache.Everything)
	}

	if err != nil {
		klog.Errorf("Failed to list pods from cache: %v", err)
		return
	}

	klog.Infof("Pods in cache (%d total):", len(pods))
	for _, pod := range pods {
		klog.Infof("  - %s/%s (Node: %s)", pod.Namespace, pod.Name, pod.Spec.NodeName)
	}
}

// GetPodFromCache 从缓存中获取单个 Pod
func (d *InformerDemo) GetPodFromCache(namespace, name string) (*corev1.Pod, error) {
	return d.factory.Core().V1().Pods().Lister().Pods(namespace).Get(name)
}

// 辅助函数

func getContainerRestarts(pod *corev1.Pod) int32 {
	var total int32
	for _, status := range pod.Status.ContainerStatuses {
		total += status.RestartCount
	}
	return total
}

func isNodeReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func main() {
	var (
		kubeconfig string
		mode       string
	)

	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig")
	flag.StringVar(&mode, "mode", "watch", "Mode: watch or list")
	flag.Parse()

	demo, err := NewInformerDemo(kubeconfig)
	if err != nil {
		klog.Fatalf("Failed to create InformerDemo: %v", err)
	}

	// 设置各种 Informer
	demo.SetupPodInformer()
	demo.SetupNodeInformer()
	demo.SetupEventInformer()

	// 处理优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		klog.Info("Received shutdown signal...")
		demo.Stop()
	}()

	if mode == "list" {
		// 短暂运行以填充缓存
		go demo.Run()
		time.Sleep(2 * time.Second)
		demo.ListPodsFromCache("")
		demo.Stop()
	} else {
		// 持续监听
		demo.Run()
	}
}
