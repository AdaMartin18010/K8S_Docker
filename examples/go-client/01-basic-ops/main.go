package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

		corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// K8sClient 封装 Kubernetes 客户端操作
type K8sClient struct {
	clientset *kubernetes.Clientset
}

// NewK8sClient 创建新的 K8s 客户端
func NewK8sClient(kubeconfig string) (*K8sClient, error) {
	// 如果未提供 kubeconfig，尝试从默认位置加载
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	// 构建配置
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		// 如果无法加载 kubeconfig，尝试使用 in-cluster 配置
		config, err = clientcmd.NewInClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create config: %w", err)
		}
	}

	// 创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &K8sClient{clientset: clientset}, nil
}

// ListPods 列出指定命名空间的所有 Pod
func (c *K8sClient) ListPods(namespace string) (*corev1.PodList, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	return pods, nil
}

// GetPod 获取单个 Pod 详情
func (c *K8sClient) GetPod(namespace, name string) (*corev1.Pod, error) {
	pod, err := c.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s: %w", name, err)
	}
	return pod, nil
}

// CreatePod 创建一个新的 Pod
func (c *K8sClient) CreatePod(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	created, err := c.clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %w", err)
	}
	return created, nil
}

// DeletePod 删除 Pod
func (c *K8sClient) DeletePod(namespace, name string) error {
	err := c.clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod: %w", err)
	}
	return nil
}

// WatchPods 监视 Pod 变化
func (c *K8sClient) WatchPods(namespace string) {
	watch, err := c.clientset.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: fields.Everything().String(),
	})
	if err != nil {
		fmt.Printf("Failed to watch pods: %v\n", err)
		return
	}
	defer watch.Stop()

	fmt.Println("Watching pods... (Press Ctrl+C to stop)")
	for event := range watch.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			continue
		}
		fmt.Printf("Event: %s, Pod: %s/%s, Phase: %s\n",
			event.Type, pod.Namespace, pod.Name, pod.Status.Phase)
	}
}

// CreateSimplePod 创建一个简单的 Nginx Pod
func (c *K8sClient) CreateSimplePod(namespace, name string) (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "nginx",
				"managed-by": "client-go-example",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:alpine",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 80,
							Protocol:      corev1.ProtocolTCP,
						},
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("512Mi"),
						},
					},
				},
			},
		},
	}

	return c.CreatePod(namespace, pod)
}

// PrintPodInfo 打印 Pod 信息
func PrintPodInfo(pod *corev1.Pod) {
	fmt.Printf("Pod: %s/%s\n", pod.Namespace, pod.Name)
	fmt.Printf("  Status: %s\n", pod.Status.Phase)
	fmt.Printf("  Node: %s\n", pod.Spec.NodeName)
	fmt.Printf("  IP: %s\n", pod.Status.PodIP)
	fmt.Printf("  Start Time: %v\n", pod.Status.StartTime)
	fmt.Printf("  Containers:\n")
	for _, container := range pod.Spec.Containers {
		fmt.Printf("    - %s: %s\n", container.Name, container.Image)
	}
}

func main() {
	var (
		kubeconfig string
		namespace  string
		operation  string
		podName    string
	)

	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	flag.StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
	flag.StringVar(&operation, "op", "list", "Operation: list, get, create, delete, watch")
	flag.StringVar(&podName, "name", "", "Pod name for get/create/delete operations")
	flag.Parse()

	// 创建客户端
	client, err := NewK8sClient(kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	switch operation {
	case "list":
		pods, err := client.ListPods(namespace)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing pods: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Pods in namespace %s:\n", namespace)
		for _, pod := range pods.Items {
			PrintPodInfo(&pod)
			fmt.Println()
		}

	case "get":
		if podName == "" {
			fmt.Fprintln(os.Stderr, "Pod name is required for get operation")
			os.Exit(1)
		}
		pod, err := client.GetPod(namespace, podName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting pod: %v\n", err)
			os.Exit(1)
		}
		PrintPodInfo(pod)

	case "create":
		if podName == "" {
			podName = "nginx-test"
		}
		pod, err := client.CreateSimplePod(namespace, podName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating pod: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created pod: %s/%s\n", pod.Namespace, pod.Name)

	case "delete":
		if podName == "" {
			fmt.Fprintln(os.Stderr, "Pod name is required for delete operation")
			os.Exit(1)
		}
		if err := client.DeletePod(namespace, podName); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting pod: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted pod: %s/%s\n", namespace, podName)

	case "watch":
		client.WatchPods(namespace)

	default:
		fmt.Fprintf(os.Stderr, "Unknown operation: %s\n", operation)
		fmt.Fprintf(os.Stderr, "Supported operations: list, get, create, delete, watch\n")
		os.Exit(1)
	}

	_ = ctx
}
