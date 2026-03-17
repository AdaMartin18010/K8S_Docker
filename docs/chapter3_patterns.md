# 第三章：Docker 与 Kubernetes 设计模式深度解析

> 本章深入分析 Docker 和 Kubernetes 生态系统中应用的设计模式，从经典的 GoF 模式到分布式系统特有的模式，为构建云原生应用提供理论指导和实践参考。

---

## 目录

- [第三章：Docker 与 Kubernetes 设计模式深度解析](#第三章docker-与-kubernetes-设计模式深度解析)
  - [目录](#目录)
  - [1. GoF 设计模式在 K8s 中的应用](#1-gof-设计模式在-k8s-中的应用)
    - [1.1 单例模式（Singleton）](#11-单例模式singleton)
      - [意图](#意图)
      - [K8s 中的应用](#k8s-中的应用)
      - [优缺点](#优缺点)
      - [适用场景](#适用场景)
    - [1.2 工厂模式（Factory）](#12-工厂模式factory)
      - [意图](#意图-1)
      - [K8s 中的应用](#k8s-中的应用-1)
      - [优缺点](#优缺点-1)
      - [适用场景](#适用场景-1)
    - [1.3 策略模式（Strategy）](#13-策略模式strategy)
      - [意图](#意图-2)
      - [K8s 中的应用](#k8s-中的应用-2)
      - [优缺点](#优缺点-2)
      - [适用场景](#适用场景-2)
    - [1.4 观察者模式（Observer）](#14-观察者模式observer)
      - [意图](#意图-3)
      - [K8s 中的应用](#k8s-中的应用-3)
      - [优缺点](#优缺点-3)
      - [适用场景](#适用场景-3)
    - [1.5 模板方法模式（Template Method）](#15-模板方法模式template-method)
      - [意图](#意图-4)
      - [K8s 中的应用](#k8s-中的应用-4)
      - [优缺点](#优缺点-4)
      - [适用场景](#适用场景-4)
    - [1.6 适配器模式（Adapter）](#16-适配器模式adapter)
      - [意图](#意图-5)
      - [K8s 中的应用](#k8s-中的应用-5)
      - [优缺点](#优缺点-5)
      - [适用场景](#适用场景-5)
    - [1.7 装饰器模式（Decorator）](#17-装饰器模式decorator)
      - [意图](#意图-6)
      - [K8s 中的应用](#k8s-中的应用-6)
      - [优缺点](#优缺点-6)
      - [适用场景](#适用场景-6)
  - [2. 分布式系统设计模式](#2-分布式系统设计模式)
    - [2.1 Sidecar 模式](#21-sidecar-模式)
      - [意图](#意图-7)
      - [结构](#结构)
      - [实现](#实现)
      - [优缺点](#优缺点-7)
      - [适用场景](#适用场景-7)
    - [2.2 Ambassador 模式](#22-ambassador-模式)
      - [意图](#意图-8)
      - [结构](#结构-1)
      - [实现](#实现-1)
      - [优缺点](#优缺点-8)
      - [适用场景](#适用场景-8)
    - [2.3 Adapter 模式（分布式版本）](#23-adapter-模式分布式版本)
      - [意图](#意图-9)
      - [结构](#结构-2)
      - [实现](#实现-2)
      - [优缺点](#优缺点-9)
      - [适用场景](#适用场景-9)
    - [2.4 Scatter-Gather 模式](#24-scatter-gather-模式)
      - [意图](#意图-10)
      - [结构](#结构-3)
      - [实现](#实现-3)
      - [优缺点](#优缺点-10)
      - [适用场景](#适用场景-10)
    - [2.5 Saga 模式](#25-saga-模式)
      - [意图](#意图-11)
      - [结构](#结构-4)
      - [实现](#实现-4)
      - [优缺点](#优缺点-11)
      - [适用场景](#适用场景-11)
    - [2.6 Circuit Breaker 模式](#26-circuit-breaker-模式)
      - [意图](#意图-12)
      - [结构](#结构-5)
      - [实现](#实现-5)
      - [优缺点](#优缺点-12)
      - [适用场景](#适用场景-12)
    - [2.7 Bulkhead 模式](#27-bulkhead-模式)
      - [意图](#意图-13)
      - [结构](#结构-6)
      - [实现](#实现-6)
      - [优缺点](#优缺点-13)
      - [适用场景](#适用场景-13)
    - [2.8 Retry 模式](#28-retry-模式)
      - [意图](#意图-14)
      - [结构](#结构-7)
      - [实现](#实现-7)
      - [优缺点](#优缺点-14)
      - [适用场景](#适用场景-14)
  - [3. 并发与并行模式](#3-并发与并行模式)
    - [3.1 Worker Pool 模式](#31-worker-pool-模式)
      - [意图](#意图-15)
      - [结构](#结构-8)
      - [实现](#实现-8)
      - [优缺点](#优缺点-15)
      - [适用场景](#适用场景-15)
    - [3.2 Pub/Sub 模式](#32-pubsub-模式)
      - [意图](#意图-16)
      - [结构](#结构-9)
      - [实现](#实现-9)
      - [优缺点](#优缺点-16)
      - [适用场景](#适用场景-16)
    - [3.3 Leader Election 模式](#33-leader-election-模式)
      - [意图](#意图-17)
      - [结构](#结构-10)
      - [实现](#实现-10)
      - [优缺点](#优缺点-17)
      - [适用场景](#适用场景-17)
    - [3.4 Distributed Lock 模式](#34-distributed-lock-模式)
      - [意图](#意图-18)
      - [结构](#结构-11)
      - [实现](#实现-11)
      - [优缺点](#优缺点-18)
      - [适用场景](#适用场景-18)
    - [3.5 Barrier 模式](#35-barrier-模式)
      - [意图](#意图-19)
      - [结构](#结构-12)
      - [实现](#实现-12)
      - [优缺点](#优缺点-19)
      - [适用场景](#适用场景-19)
    - [3.6 Pipeline 模式](#36-pipeline-模式)
      - [意图](#意图-20)
      - [结构](#结构-13)
      - [实现](#实现-13)
      - [优缺点](#优缺点-20)
      - [适用场景](#适用场景-20)
  - [4. 同步与异步模式](#4-同步与异步模式)
    - [4.1 同步调用模式](#41-同步调用模式)
      - [意图](#意图-21)
      - [结构](#结构-14)
      - [实现](#实现-14)
      - [优缺点](#优缺点-21)
      - [适用场景](#适用场景-21)
    - [4.2 异步消息模式](#42-异步消息模式)
      - [意图](#意图-22)
      - [结构](#结构-15)
      - [实现](#实现-15)
      - [优缺点](#优缺点-22)
      - [适用场景](#适用场景-22)
    - [4.3 Watch 机制](#43-watch-机制)
      - [意图](#意图-23)
      - [结构](#结构-16)
      - [实现](#实现-16)
      - [优缺点](#优缺点-23)
      - [适用场景](#适用场景-23)
    - [4.4 Informer 模式](#44-informer-模式)
      - [意图](#意图-24)
      - [结构](#结构-17)
      - [实现](#实现-17)
      - [优缺点](#优缺点-24)
      - [适用场景](#适用场景-24)
    - [4.5 Work Queue 模式](#45-work-queue-模式)
      - [意图](#意图-25)
      - [结构](#结构-18)
      - [实现](#实现-18)
      - [优缺点](#优缺点-25)
      - [适用场景](#适用场景-25)
    - [4.6 Channel 模式](#46-channel-模式)
      - [意图](#意图-26)
      - [实现](#实现-19)
      - [优缺点](#优缺点-26)
      - [适用场景](#适用场景-26)
  - [5. 工作流设计模式](#5-工作流设计模式)
    - [5.1 CronJob 模式](#51-cronjob-模式)
      - [意图](#意图-27)
      - [结构](#结构-19)
      - [实现](#实现-20)
      - [优缺点](#优缺点-27)
      - [适用场景](#适用场景-27)
    - [5.2 Job 模式](#52-job-模式)
      - [意图](#意图-28)
      - [结构](#结构-20)
      - [实现](#实现-21)
      - [优缺点](#优缺点-28)
      - [适用场景](#适用场景-28)
    - [5.3 DaemonSet 模式](#53-daemonset-模式)
      - [意图](#意图-29)
      - [结构](#结构-21)
      - [实现](#实现-22)
      - [优缺点](#优缺点-29)
      - [适用场景](#适用场景-29)
    - [5.4 Pipeline 模式](#54-pipeline-模式)
      - [意图](#意图-30)
      - [结构](#结构-22)
      - [实现](#实现-23)
      - [优缺点](#优缺点-30)
      - [适用场景](#适用场景-30)
    - [5.5 DAG 模式](#55-dag-模式)
      - [意图](#意图-31)
      - [结构](#结构-23)
      - [实现](#实现-24)
      - [优缺点](#优缺点-31)
      - [适用场景](#适用场景-31)
    - [5.6 State Machine 模式](#56-state-machine-模式)
      - [意图](#意图-32)
      - [结构](#结构-24)
      - [实现](#实现-25)
      - [优缺点](#优缺点-32)
      - [适用场景](#适用场景-32)
  - [6. K8s 特有模式](#6-k8s-特有模式)
    - [6.1 Controller 模式](#61-controller-模式)
      - [意图](#意图-33)
      - [结构](#结构-25)
      - [实现](#实现-26)
      - [优缺点](#优缺点-33)
      - [适用场景](#适用场景-33)
    - [6.2 Operator 模式](#62-operator-模式)
      - [意图](#意图-34)
      - [结构](#结构-26)
      - [实现](#实现-27)
      - [优缺点](#优缺点-34)
      - [适用场景](#适用场景-34)
    - [6.3 Initializer 模式](#63-initializer-模式)
      - [意图](#意图-35)
      - [结构](#结构-27)
      - [实现](#实现-28)
      - [优缺点](#优缺点-35)
      - [适用场景](#适用场景-35)
    - [6.4 Finalizer 模式](#64-finalizer-模式)
      - [意图](#意图-36)
      - [结构](#结构-28)
      - [实现](#实现-29)
      - [优缺点](#优缺点-36)
      - [适用场景](#适用场景-36)
    - [6.5 Owner Reference 模式](#65-owner-reference-模式)
      - [意图](#意图-37)
      - [结构](#结构-29)
      - [实现](#实现-30)
      - [优缺点](#优缺点-37)
      - [适用场景](#适用场景-37)
  - [7. Go 代码实现示例](#7-go-代码实现示例)
    - [7.1 完整 Controller 协调循环实现](#71-完整-controller-协调循环实现)
    - [7.2 完整 Informer 模式实现](#72-完整-informer-模式实现)
    - [7.3 完整 Work Queue 模式实现](#73-完整-work-queue-模式实现)
    - [7.4 完整 Leader Election 实现](#74-完整-leader-election-实现)
    - [7.5 完整 Sidecar 模式框架实现](#75-完整-sidecar-模式框架实现)
    - [7.6 完整 Circuit Breaker 实现](#76-完整-circuit-breaker-实现)
  - [总结](#总结)

---

## 1. GoF 设计模式在 K8s 中的应用

Kubernetes 的架构设计大量借鉴了经典的 GoF（Gang of Four）设计模式，这些模式在分布式系统中得到了新的诠释和应用。

### 1.1 单例模式（Singleton）

#### 意图

确保一个类只有一个实例，并提供一个全局访问点。在 K8s 中，单例模式用于管理共享资源，如数据库连接、配置管理器等。

#### K8s 中的应用

**API Server 的 etcd 连接**

```go
// EtcdClient 单例实现
package storage

import (
    "sync"
    "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
    client *clientv3.Client
    once   sync.Once
}

var (
    instance *EtcdClient
    mu       sync.Mutex
)

// GetInstance 返回 EtcdClient 的单例实例
func GetInstance(endpoints []string) (*EtcdClient, error) {
    mu.Lock()
    defer mu.Unlock()

    if instance == nil {
        instance = &EtcdClient{}
        err := instance.init(endpoints)
        if err != nil {
            instance = nil
            return nil, err
        }
    }
    return instance, nil
}

func (e *EtcdClient) init(endpoints []string) error {
    var err error
    e.once.Do(func() {
        e.client, err = clientv3.New(clientv3.Config{
            Endpoints: endpoints,
        })
    })
    return err
}

// GetClient 获取 etcd 客户端
func (e *EtcdClient) GetClient() *clientv3.Client {
    return e.client
}
```

**Kubernetes 中的实际应用：**

- `kube-apiserver` 与 etcd 的连接管理
- `Scheduler` 的调度器实例
- `Controller Manager` 的控制器实例

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 确保全局唯一性 | 违反单一职责原则 |
| 节省系统资源 | 难以单元测试 |
| 控制并发访问 | 隐藏类之间的依赖 |

#### 适用场景

- 数据库连接池管理
- 全局配置管理
- 共享缓存实例

---

### 1.2 工厂模式（Factory）

#### 意图

定义一个创建对象的接口，让子类决定实例化哪一个类。工厂方法使类的实例化延迟到子类。

#### K8s 中的应用

**Pod 和 Container 的创建**

```go
// PodFactory 接口定义
package factory

import (
    corev1 "k8s.io/api/core/v1"
)

// PodCreator 定义 Pod 创建接口
type PodCreator interface {
    CreatePod(name string, spec corev1.PodSpec) *corev1.Pod
}

// DeploymentPodFactory 创建 Deployment 管理的 Pod
type DeploymentPodFactory struct {
    DeploymentName string
    ReplicaSetName string
}

func (d *DeploymentPodFactory) CreatePod(name string, spec corev1.PodSpec) *corev1.Pod {
    pod := &corev1.Pod{}
    pod.Name = name
    pod.Spec = spec

    // 添加 Deployment 特有的标签
    pod.Labels = map[string]string{
        "app":               d.DeploymentName,
        "pod-template-hash": d.ReplicaSetName,
    }

    // 设置 OwnerReference
    pod.OwnerReferences = []metav1.OwnerReference{
        {
            APIVersion: "apps/v1",
            Kind:       "ReplicaSet",
            Name:       d.ReplicaSetName,
            Controller: boolPtr(true),
        },
    }

    return pod
}

// StatefulSetPodFactory 创建 StatefulSet 管理的 Pod
type StatefulSetPodFactory struct {
    StatefulSetName string
    Ordinal         int
}

func (s *StatefulSetPodFactory) CreatePod(name string, spec corev1.PodSpec) *corev1.Pod {
    pod := &corev1.Pod{}
    pod.Name = fmt.Sprintf("%s-%d", s.StatefulSetName, s.Ordinal)
    pod.Spec = spec

    // StatefulSet Pod 有固定的网络标识
    pod.Labels = map[string]string{
        "app":                  s.StatefulSetName,
        "statefulset.kubernetes.io/pod-name": pod.Name,
    }

    // 添加 hostname 和 subdomain
    pod.Spec.Hostname = pod.Name

    return pod
}

// JobPodFactory 创建 Job 管理的 Pod
type JobPodFactory struct {
    JobName string
}

func (j *JobPodFactory) CreatePod(name string, spec corev1.PodSpec) *corev1.Pod {
    pod := &corev1.Pod{}
    pod.Name = name
    pod.Spec = spec

    // Job Pod 通常不需要重启
    pod.Spec.RestartPolicy = corev1.RestartPolicyNever

    pod.Labels = map[string]string{
        "job-name": j.JobName,
    }

    return pod
}

// PodFactory 工厂方法
type PodFactory struct{}

func (f *PodFactory) CreatePod(kind string, name string, spec corev1.PodSpec) PodCreator {
    switch kind {
    case "Deployment":
        return &DeploymentPodFactory{DeploymentName: name}
    case "StatefulSet":
        return &StatefulSetPodFactory{StatefulSetName: name}
    case "Job":
        return &JobPodFactory{JobName: name}
    default:
        return &DeploymentPodFactory{DeploymentName: name}
    }
}
```

**Kubernetes 中的实际应用：**

- `kubelet` 创建不同类型的容器
- `kube-proxy` 创建不同的代理规则
- `Scheduler` 创建调度上下文

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 解耦对象创建和使用 | 增加类的数量 |
| 易于扩展新的产品类型 | 代码复杂度增加 |
| 符合开闭原则 | 需要额外的抽象层 |

#### 适用场景

- 创建复杂对象
- 需要根据不同条件创建不同类型对象
- 对象创建逻辑需要集中管理

---

### 1.3 策略模式（Strategy）

#### 意图

定义一系列算法，把它们一个个封装起来，并且使它们可互相替换。策略模式让算法的变化独立于使用算法的客户。

#### K8s 中的应用

**CNI/CRI/CSI 插件架构**

```go
// CNI 策略模式实现
package cni

import (
    "context"
    "encoding/json"
)

// NetworkPlugin 定义网络插件接口
type NetworkPlugin interface {
    Name() string
    SetupPodNetwork(ctx context.Context, podName, namespace string, annotations map[string]string) error
    TeardownPodNetwork(ctx context.Context, podName, namespace string) error
    GetPodNetworkStatus(ctx context.Context, podName, namespace string) (*NetworkStatus, error)
}

// NetworkStatus 网络状态
type NetworkStatus struct {
    IP       string   `json:"ip"`
    MAC      string   `json:"mac,omitempty"`
    DNS      []string `json:"dns,omitempty"`
    Gateways []string `json:"gateways,omitempty"`
}

// CalicoPlugin Calico CNI 实现
type CalicoPlugin struct {
    etcdEndpoints []string
    calicoConfig  *CalicoConfig
}

type CalicoConfig struct {
    IPPool    string `json:"ip_pool"`
    IPIPMode  string `json:"ipip_mode"`
    VXLANMode string `json:"vxlan_mode"`
}

func NewCalicoPlugin(config []byte) (*CalicoPlugin, error) {
    var calicoConfig CalicoConfig
    if err := json.Unmarshal(config, &calicoConfig); err != nil {
        return nil, err
    }
    return &CalicoPlugin{calicoConfig: &calicoConfig}, nil
}

func (c *CalicoPlugin) Name() string {
    return "calico"
}

func (c *CalicoPlugin) SetupPodNetwork(ctx context.Context, podName, namespace string, annotations map[string]string) error {
    // Calico 特定的网络设置逻辑
    // 1. 从 IP Pool 分配 IP
    // 2. 创建 veth pair
    // 3. 配置路由
    return nil
}

func (c *CalicoPlugin) TeardownPodNetwork(ctx context.Context, podName, namespace string) error {
    // 清理 Calico 网络配置
    return nil
}

func (c *CalicoPlugin) GetPodNetworkStatus(ctx context.Context, podName, namespace string) (*NetworkStatus, error) {
    // 获取 Calico 网络状态
    return &NetworkStatus{}, nil
}

// FlannelPlugin Flannel CNI 实现
type FlannelPlugin struct {
    subnetFile string
    etcdPrefix string
}

func NewFlannelPlugin(config []byte) (*FlannelPlugin, error) {
    return &FlannelPlugin{subnetFile: "/run/flannel/subnet.env"}, nil
}

func (f *FlannelPlugin) Name() string {
    return "flannel"
}

func (f *FlannelPlugin) SetupPodNetwork(ctx context.Context, podName, namespace string, annotations map[string]string) error {
    // Flannel 特定的网络设置逻辑
    // 1. 读取 subnet.env
    // 2. 分配子网
    // 3. 配置 VXLAN
    return nil
}

func (f *FlannelPlugin) TeardownPodNetwork(ctx context.Context, podName, namespace string) error {
    // 清理 Flannel 网络配置
    return nil
}

func (f *FlannelPlugin) GetPodNetworkStatus(ctx context.Context, podName, namespace string) (*NetworkStatus, error) {
    return &NetworkStatus{}, nil
}

// CiliumPlugin Cilium CNI 实现
type CiliumPlugin struct {
    socketPath string
}

func NewCiliumPlugin(config []byte) (*CiliumPlugin, error) {
    return &CiliumPlugin{socketPath: "/var/run/cilium/cilium.sock"}, nil
}

func (c *CiliumPlugin) Name() string {
    return "cilium"
}

func (c *CiliumPlugin) SetupPodNetwork(ctx context.Context, podName, namespace string, annotations map[string]string) error {
    // Cilium 使用 eBPF 进行网络配置
    return nil
}

func (c *CiliumPlugin) TeardownPodNetwork(ctx context.Context, podName, namespace string) error {
    return nil
}

func (c *CiliumPlugin) GetPodNetworkStatus(ctx context.Context, podName, namespace string) (*NetworkStatus, error) {
    return &NetworkStatus{}, nil
}

// PluginManager 插件管理器
type PluginManager struct {
    plugins map[string]NetworkPlugin
}

func NewPluginManager() *PluginManager {
    return &PluginManager{
        plugins: make(map[string]NetworkPlugin),
    }
}

func (pm *PluginManager) RegisterPlugin(plugin NetworkPlugin) {
    pm.plugins[plugin.Name()] = plugin
}

func (pm *PluginManager) GetPlugin(name string) (NetworkPlugin, bool) {
    plugin, ok := pm.plugins[name]
    return plugin, ok
}

func (pm *PluginManager) SetupPodNetwork(ctx context.Context, pluginName, podName, namespace string, annotations map[string]string) error {
    plugin, ok := pm.GetPlugin(pluginName)
    if !ok {
        return fmt.Errorf("plugin %s not found", pluginName)
    }
    return plugin.SetupPodNetwork(ctx, podName, namespace, annotations)
}
```

**CRI 策略模式实现：**

```go
// CRI 运行时接口
package cri

import (
    runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// RuntimeService 定义容器运行时接口
type RuntimeService interface {
    // 容器生命周期管理
    CreateContainer(ctx context.Context, podSandboxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error)
    StartContainer(ctx context.Context, containerID string) error
    StopContainer(ctx context.Context, containerID string, timeout int64) error
    RemoveContainer(ctx context.Context, containerID string) error

    // 镜像管理
    PullImage(ctx context.Context, image *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error)
    ListImages(ctx context.Context, filter *runtimeapi.ImageFilter) ([]*runtimeapi.Image, error)
    RemoveImage(ctx context.Context, image *runtimeapi.ImageSpec) error

    // 运行时信息
    Version(ctx context.Context) (*runtimeapi.VersionResponse, error)
    Status(ctx context.Context, verbose bool) (*runtimeapi.StatusResponse, error)
}

// DockerRuntime Docker 运行时实现
type DockerRuntime struct {
    client *docker.Client
}

func NewDockerRuntime(endpoint string) (*DockerRuntime, error) {
    client, err := docker.NewClient(endpoint)
    if err != nil {
        return nil, err
    }
    return &DockerRuntime{client: client}, nil
}

func (d *DockerRuntime) CreateContainer(ctx context.Context, podSandboxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
    // Docker 特定的容器创建逻辑
    return "", nil
}

// ContainerdRuntime containerd 运行时实现
type ContainerdRuntime struct {
    client *containerd.Client
}

func NewContainerdRuntime(address string) (*ContainerdRuntime, error) {
    client, err := containerd.New(address)
    if err != nil {
        return nil, err
    }
    return &ContainerdRuntime{client: client}, nil
}

func (c *ContainerdRuntime) CreateContainer(ctx context.Context, podSandboxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
    // containerd 特定的容器创建逻辑
    return "", nil
}

// CRIORuntime CRI-O 运行时实现
type CRIORuntime struct {
    client *crio.Client
}

func NewCRIORuntime(endpoint string) (*CRIORuntime, error) {
    client, err := crio.New(endpoint)
    if err != nil {
        return nil, err
    }
    return &CRIORuntime{client: client}, nil
}

func (c *CRIORuntime) CreateContainer(ctx context.Context, podSandboxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
    // CRI-O 特定的容器创建逻辑
    return "", nil
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 算法可自由切换 | 客户端必须了解所有策略 |
| 避免多重条件语句 | 增加对象数量 |
| 易于扩展新算法 | 策略选择需要额外逻辑 |
| 符合开闭原则 | |

#### 适用场景

- 多种算法或行为需要动态选择
- 算法需要独立于使用者变化
- 需要消除大量的条件判断语句

---

### 1.4 观察者模式（Observer）

#### 意图

定义对象间的一对多依赖关系，当一个对象状态发生改变时，所有依赖于它的对象都得到通知并被自动更新。

#### K8s 中的应用

**Informer 和 Watch 机制**

```go
// Observer 模式实现
package observer

import (
    "context"
    "sync"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/watch"
)

// ResourceEvent 资源事件
type ResourceEvent struct {
    Type   watch.EventType
    Object metav1.Object
}

// ResourceEventHandler 资源事件处理器接口
type ResourceEventHandler interface {
    OnAdd(obj metav1.Object)
    OnUpdate(oldObj, newObj metav1.Object)
    OnDelete(obj metav1.Object)
}

// ResourceEventHandlerFuncs 函数式事件处理器
type ResourceEventHandlerFuncs struct {
    AddFunc    func(obj metav1.Object)
    UpdateFunc func(oldObj, newObj metav1.Object)
    DeleteFunc func(obj metav1.Object)
}

func (r *ResourceEventHandlerFuncs) OnAdd(obj metav1.Object) {
    if r.AddFunc != nil {
        r.AddFunc(obj)
    }
}

func (r *ResourceEventHandlerFuncs) OnUpdate(oldObj, newObj metav1.Object) {
    if r.UpdateFunc != nil {
        r.UpdateFunc(oldObj, newObj)
    }
}

func (r *ResourceEventHandlerFuncs) OnDelete(obj metav1.Object) {
    if r.DeleteFunc != nil {
        r.DeleteFunc(obj)
    }
}

// Observable 可观察对象接口
type Observable interface {
    AddObserver(handler ResourceEventHandler)
    RemoveObserver(handler ResourceEventHandler)
    NotifyObservers(event ResourceEvent)
}

// ResourceWatcher 资源观察者
type ResourceWatcher struct {
    observers []ResourceEventHandler
    mu        sync.RWMutex
    eventChan chan ResourceEvent
}

func NewResourceWatcher() *ResourceWatcher {
    return &ResourceWatcher{
        observers: make([]ResourceEventHandler, 0),
        eventChan: make(chan ResourceEvent, 100),
    }
}

func (rw *ResourceWatcher) AddObserver(handler ResourceEventHandler) {
    rw.mu.Lock()
    defer rw.mu.Unlock()
    rw.observers = append(rw.observers, handler)
}

func (rw *ResourceWatcher) RemoveObserver(handler ResourceEventHandler) {
    rw.mu.Lock()
    defer rw.mu.Unlock()

    for i, h := range rw.observers {
        if h == handler {
            rw.observers = append(rw.observers[:i], rw.observers[i+1:]...)
            break
        }
    }
}

func (rw *ResourceWatcher) NotifyObservers(event ResourceEvent) {
    rw.mu.RLock()
    observers := make([]ResourceEventHandler, len(rw.observers))
    copy(observers, rw.observers)
    rw.mu.RUnlock()

    for _, handler := range observers {
        switch event.Type {
        case watch.Added:
            handler.OnAdd(event.Object)
        case watch.Modified:
            // 对于更新事件，需要保存旧对象
            handler.OnUpdate(nil, event.Object)
        case watch.Deleted:
            handler.OnDelete(event.Object)
        }
    }
}

func (rw *ResourceWatcher) Start(ctx context.Context) {
    go func() {
        for {
            select {
            case event := <-rw.eventChan:
                rw.NotifyObservers(event)
            case <-ctx.Done():
                return
            }
        }
    }()
}

func (rw *ResourceWatcher) SendEvent(event ResourceEvent) {
    select {
    case rw.eventChan <- event:
    default:
        // 通道满，丢弃事件或记录日志
    }
}

// PodObserver Pod 观察者实现
type PodObserver struct {
    name string
}

func NewPodObserver(name string) *PodObserver {
    return &PodObserver{name: name}
}

func (p *PodObserver) OnAdd(obj metav1.Object) {
    log.Printf("[%s] Pod added: %s/%s", p.name, obj.GetNamespace(), obj.GetName())
}

func (p *PodObserver) OnUpdate(oldObj, newObj metav1.Object) {
    log.Printf("[%s] Pod updated: %s/%s", p.name, newObj.GetNamespace(), newObj.GetName())
}

func (p *PodObserver) OnDelete(obj metav1.Object) {
    log.Printf("[%s] Pod deleted: %s/%s", p.name, obj.GetNamespace(), obj.GetName())
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 松耦合的交互 | 可能导致循环依赖 |
| 支持广播通信 | 通知顺序不确定 |
| 易于扩展观察者 | 内存泄漏风险 |

#### 适用场景

- 事件驱动架构
- 需要实时响应状态变化
- 发布-订阅模式

---

### 1.5 模板方法模式（Template Method）

#### 意图

定义一个操作中的算法骨架，而将一些步骤延迟到子类中。模板方法使得子类可以不改变一个算法的结构即可重定义该算法的某些特定步骤。

#### K8s 中的应用

**控制器的协调循环**

```go
// 模板方法模式实现
package controller

import (
    "context"
    "fmt"
    "time"
)

// Controller 控制器接口
type Controller interface {
    Name() string
    Run(ctx context.Context, workers int) error
}

// BaseController 基础控制器（模板）
type BaseController struct {
    name           string
    queue          workqueue.RateLimitingInterface
    informer       cache.SharedIndexInformer
    syncHandler    func(key string) error
    processNextItem func() bool
}

// ControllerTemplate 控制器模板接口
type ControllerTemplate interface {
    // 必须由子类实现的方法
    GetObject(key string) (interface{}, error)
    Reconcile(obj interface{}) error
    ShouldProcess(obj interface{}) bool

    // 可选覆盖的方法
    PreProcess(key string) error
    PostProcess(key string, err error)
    OnError(key string, err error) error
}

// GenericController 通用控制器实现模板方法
type GenericController struct {
    *BaseController
    template ControllerTemplate
}

func NewGenericController(name string, template ControllerTemplate) *GenericController {
    c := &GenericController{
        BaseController: &BaseController{
            name:  name,
            queue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
        },
        template: template,
    }
    c.syncHandler = c.processItem
    return c
}

// processItem 模板方法 - 定义协调流程
func (c *GenericController) processItem(key string) error {
    // 1. 前置处理（钩子方法）
    if err := c.template.PreProcess(key); err != nil {
        return err
    }

    // 2. 获取对象
    obj, err := c.template.GetObject(key)
    if err != nil {
        return err
    }

    // 3. 判断是否处理
    if !c.template.ShouldProcess(obj) {
        return nil
    }

    // 4. 执行协调（抽象方法）
    err = c.template.Reconcile(obj)

    // 5. 错误处理（钩子方法）
    if err != nil {
        if handleErr := c.template.OnError(key, err); handleErr != nil {
            return handleErr
        }
    }

    // 6. 后置处理（钩子方法）
    c.template.PostProcess(key, err)

    return err
}

// DeploymentController Deployment 控制器实现
type DeploymentController struct {
    kubeClient kubernetes.Interface
}

func NewDeploymentController(client kubernetes.Interface) *DeploymentController {
    return &DeploymentController{kubeClient: client}
}

func (d *DeploymentController) GetObject(key string) (interface{}, error) {
    namespace, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        return nil, err
    }
    return d.kubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (d *DeploymentController) Reconcile(obj interface{}) error {
    deployment := obj.(*appsv1.Deployment)

    // 1. 获取关联的 ReplicaSet
    rsList, err := d.getReplicaSetsForDeployment(deployment)
    if err != nil {
        return err
    }

    // 2. 计算新的 ReplicaSet
    newRS, err := d.getNewReplicaSet(deployment, rsList)
    if err != nil {
        return err
    }

    // 3. 同步 ReplicaSet
    if err := d.syncReplicaSet(deployment, newRS); err != nil {
        return err
    }

    // 4. 更新 Deployment 状态
    return d.updateDeploymentStatus(deployment, rsList, newRS)
}

func (d *DeploymentController) ShouldProcess(obj interface{}) bool {
    deployment := obj.(*appsv1.Deployment)
    // 检查是否需要处理：暂停、删除等状态
    return deployment.DeletionTimestamp == nil && !deployment.Spec.Paused
}

func (d *DeploymentController) PreProcess(key string) error {
    // 可以记录开始处理的日志
    return nil
}

func (d *DeploymentController) PostProcess(key string, err error) {
    // 可以记录处理完成的日志
}

func (d *DeploymentController) OnError(key string, err error) error {
    // 实现错误处理逻辑，如重试、告警等
    return err
}

// ServiceController Service 控制器实现
type ServiceController struct {
    kubeClient kubernetes.Interface
}

func NewServiceController(client kubernetes.Interface) *ServiceController {
    return &ServiceController{kubeClient: client}
}

func (s *ServiceController) GetObject(key string) (interface{}, error) {
    namespace, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        return nil, err
    }
    return s.kubeClient.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (s *ServiceController) Reconcile(obj interface{}) error {
    service := obj.(*corev1.Service)

    // Service 特定的协调逻辑
    // 1. 获取 Endpoints
    // 2. 同步 EndpointSlice
    // 3. 更新负载均衡器状态

    return nil
}

func (s *ServiceController) ShouldProcess(obj interface{}) bool {
    service := obj.(*corev1.Service)
    return service.DeletionTimestamp == nil
}

func (s *ServiceController) PreProcess(key string) error  { return nil }
func (s *ServiceController) PostProcess(key string, err error) {}
func (s *ServiceController) OnError(key string, err error) error { return err }
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 代码复用 | 类数量增加 |
| 算法结构稳定 | 继承的局限性 |
| 易于扩展 | 钩子方法过多 |

#### 适用场景

- 多个类有共同的算法骨架
- 需要控制子类的扩展点
- 算法的某些步骤需要延迟实现

---

### 1.6 适配器模式（Adapter）

#### 意图

将一个类的接口转换成客户希望的另外一个接口。适配器模式使得原本由于接口不兼容而不能一起工作的那些类可以一起工作。

#### K8s 中的应用

**CRI 适配不同 Runtime**

```go
// 适配器模式实现
package adapter

import (
    "context"
    "encoding/json"

    runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// LegacyDockerClient 旧的 Docker 客户端接口
type LegacyDockerClient interface {
    CreateContainer(config *LegacyContainerConfig) (string, error)
    StartContainer(id string) error
    StopContainer(id string, timeout int) error
    InspectContainer(id string) (*LegacyContainerInfo, error)
}

type LegacyContainerConfig struct {
    Image        string
    Cmd          []string
    Env          map[string]string
    PortBindings map[string]string
    Volumes      map[string]string
}

type LegacyContainerInfo struct {
    ID      string
    Name    string
    State   string
    Pid     int
    IP      string
    Created string
}

// LegacyDockerAdapter 适配器实现
type LegacyDockerAdapter struct {
    legacyClient LegacyDockerClient
}

func NewLegacyDockerAdapter(client LegacyDockerClient) *LegacyDockerAdapter {
    return &LegacyDockerAdapter{legacyClient: client}
}

// 实现 RuntimeService 接口

func (a *LegacyDockerAdapter) Version(ctx context.Context) (*runtimeapi.VersionResponse, error) {
    return &runtimeapi.VersionResponse{
        Version:           "0.1.0",
        RuntimeName:       "docker-legacy",
        RuntimeVersion:    "1.0",
        RuntimeApiVersion: "v1",
    }, nil
}

func (a *LegacyDockerAdapter) RunPodSandbox(ctx context.Context, config *runtimeapi.PodSandboxConfig) (string, error) {
    // 适配 PodSandbox 到 Docker 容器
    legacyConfig := &LegacyContainerConfig{
        Image: "pause:3.6", // pause 容器
        Env:   config.GetLinux().GetSecurityContext().GetSeccompProfilePath(),
    }
    return a.legacyClient.CreateContainer(legacyConfig)
}

func (a *LegacyDockerAdapter) StopPodSandbox(ctx context.Context, podSandboxID string) error {
    return a.legacyClient.StopContainer(podSandboxID, 30)
}

func (a *LegacyDockerAdapter) RemovePodSandbox(ctx context.Context, podSandboxID string) error {
    // Docker 没有直接对应的 API，使用 Stop + 删除
    return a.legacyClient.StopContainer(podSandboxID, 0)
}

func (a *LegacyDockerAdapter) CreateContainer(ctx context.Context, podSandboxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
    // 转换 CRI 配置到 Docker 配置
    legacyConfig := a.convertToLegacyConfig(config)
    return a.legacyClient.CreateContainer(legacyConfig)
}

func (a *LegacyDockerAdapter) StartContainer(ctx context.Context, containerID string) error {
    return a.legacyClient.StartContainer(containerID)
}

func (a *LegacyDockerAdapter) StopContainer(ctx context.Context, containerID string, timeout int64) error {
    return a.legacyClient.StopContainer(containerID, int(timeout))
}

func (a *LegacyDockerAdapter) RemoveContainer(ctx context.Context, containerID string) error {
    return a.legacyClient.StopContainer(containerID, 0)
}

func (a *LegacyDockerAdapter) ListContainers(ctx context.Context, filter *runtimeapi.ContainerFilter) ([]*runtimeapi.Container, error) {
    // 适配过滤条件并查询
    return []*runtimeapi.Container{}, nil
}

func (a *LegacyDockerAdapter) ContainerStatus(ctx context.Context, containerID string, verbose bool) (*runtimeapi.ContainerStatusResponse, error) {
    info, err := a.legacyClient.InspectContainer(containerID)
    if err != nil {
        return nil, err
    }

    // 转换状态
    var state runtimeapi.ContainerState
    switch info.State {
    case "running":
        state = runtimeapi.ContainerState_CONTAINER_RUNNING
    case "exited":
        state = runtimeapi.ContainerState_CONTAINER_EXITED
    default:
        state = runtimeapi.ContainerState_CONTAINER_CREATED
    }

    return &runtimeapi.ContainerStatusResponse{
        Status: &runtimeapi.ContainerStatus{
            Id:       info.ID,
            State:    state,
            CreatedAt: info.Created,
        },
    }, nil
}

// 辅助转换方法
func (a *LegacyDockerAdapter) convertToLegacyConfig(config *runtimeapi.ContainerConfig) *LegacyContainerConfig {
    envMap := make(map[string]string)
    for _, env := range config.Envs {
        envMap[env.Key] = env.Value
    }

    return &LegacyContainerConfig{
        Image: config.Image.Image,
        Cmd:   config.Command,
        Env:   envMap,
    }
}

// ContainerdAdapter containerd 适配器
type ContainerdAdapter struct {
    client *containerd.Client
}

func NewContainerdAdapter(client *containerd.Client) *ContainerdAdapter {
    return &ContainerdAdapter{client: client}
}

func (c *ContainerdAdapter) Version(ctx context.Context) (*runtimeapi.VersionResponse, error) {
    version, err := c.client.Version(ctx)
    if err != nil {
        return nil, err
    }
    return &runtimeapi.VersionResponse{
        Version:           version.Version,
        RuntimeName:       "containerd",
        RuntimeVersion:    version.Version,
        RuntimeApiVersion: version.APIVersion,
    }, nil
}

// 实现其他 RuntimeService 接口方法...
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 复用现有类 | 增加系统复杂度 |
| 解耦目标类和适配者 | 需要额外的适配器类 |
| 灵活性和扩展性好 | 代码可读性降低 |

#### 适用场景

- 使用已有的类，但接口不符合需求
- 需要统一多个类的接口
- 需要与第三方库集成

---

### 1.7 装饰器模式（Decorator）

#### 意图

动态地给一个对象添加一些额外的职责。就增加功能来说，装饰器模式相比生成子类更为灵活。

#### K8s 中的应用

**Sidecar 模式**

```go
// 装饰器模式实现
package decorator

import (
    "context"
    corev1 "k8s.io/api/core/v1"
)

// PodDecorator Pod 装饰器接口
type PodDecorator interface {
    Decorate(pod *corev1.Pod) (*corev1.Pod, error)
    Name() string
}

// BasePodDecorator 基础装饰器
type BasePodDecorator struct {
    next PodDecorator
}

func (b *BasePodDecorator) SetNext(next PodDecorator) {
    b.next = next
}

func (b *BasePodDecorator) DecorateNext(pod *corev1.Pod) (*corev1.Pod, error) {
    if b.next != nil {
        return b.next.Decorate(pod)
    }
    return pod, nil
}

// LoggingSidecarDecorator 日志收集 Sidecar 装饰器
type LoggingSidecarDecorator struct {
    BasePodDecorator
    image        string
    logPath      string
    outputTarget string
}

func NewLoggingSidecarDecorator(image, logPath, outputTarget string) *LoggingSidecarDecorator {
    return &LoggingSidecarDecorator{
        image:        image,
        logPath:      logPath,
        outputTarget: outputTarget,
    }
}

func (l *LoggingSidecarDecorator) Name() string {
    return "logging-sidecar"
}

func (l *LoggingSidecarDecorator) Decorate(pod *corev1.Pod) (*corev1.Pod, error) {
    // 添加日志收集 Sidecar 容器
    sidecar := corev1.Container{
        Name:  "log-collector",
        Image: l.image,
        VolumeMounts: []corev1.VolumeMount{
            {
                Name:      "log-volume",
                MountPath: l.logPath,
            },
        },
        Env: []corev1.EnvVar{
            {
                Name:  "LOG_OUTPUT_TARGET",
                Value: l.outputTarget,
            },
            {
                Name: "POD_NAME",
                ValueFrom: &corev1.EnvVarSource{
                    FieldRef: &corev1.ObjectFieldSelector{
                        FieldPath: "metadata.name",
                    },
                },
            },
        },
    }

    pod.Spec.Containers = append(pod.Spec.Containers, sidecar)

    // 添加共享日志卷
    volumeExists := false
    for _, v := range pod.Spec.Volumes {
        if v.Name == "log-volume" {
            volumeExists = true
            break
        }
    }

    if !volumeExists {
        pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
            Name: "log-volume",
            VolumeSource: corev1.VolumeSource{
                EmptyDir: &corev1.EmptyDirVolumeSource{},
            },
        })
    }

    // 为主容器添加卷挂载
    for i := range pod.Spec.Containers {
        if pod.Spec.Containers[i].Name != "log-collector" {
            pod.Spec.Containers[i].VolumeMounts = append(
                pod.Spec.Containers[i].VolumeMounts,
                corev1.VolumeMount{
                    Name:      "log-volume",
                    MountPath: l.logPath,
                },
            )
        }
    }

    return l.DecorateNext(pod)
}

// MetricsSidecarDecorator 监控 Sidecar 装饰器
type MetricsSidecarDecorator struct {
    BasePodDecorator
    image    string
    port     int32
    endpoint string
}

func NewMetricsSidecarDecorator(image string, port int32, endpoint string) *MetricsSidecarDecorator {
    return &MetricsSidecarDecorator{
        image:    image,
        port:     port,
        endpoint: endpoint,
    }
}

func (m *MetricsSidecarDecorator) Name() string {
    return "metrics-sidecar"
}

func (m *MetricsSidecarDecorator) Decorate(pod *corev1.Pod) (*corev1.Pod, error) {
    sidecar := corev1.Container{
        Name:  "metrics-exporter",
        Image: m.image,
        Ports: []corev1.ContainerPort{
            {
                Name:          "metrics",
                ContainerPort: m.port,
                Protocol:      corev1.ProtocolTCP,
            },
        },
        Env: []corev1.EnvVar{
            {
                Name:  "METRICS_ENDPOINT",
                Value: m.endpoint,
            },
            {
                Name:  "METRICS_PORT",
                Value: string(m.port),
            },
        },
        Resources: corev1.ResourceRequirements{
            Limits: corev1.ResourceList{
                corev1.ResourceCPU:    resource.MustParse("100m"),
                corev1.ResourceMemory: resource.MustParse("128Mi"),
            },
        },
    }

    pod.Spec.Containers = append(pod.Spec.Containers, sidecar)

    // 添加 Prometheus 注解
    if pod.Annotations == nil {
        pod.Annotations = make(map[string]string)
    }
    pod.Annotations["prometheus.io/scrape"] = "true"
    pod.Annotations["prometheus.io/port"] = string(m.port)
    pod.Annotations["prometheus.io/path"] = "/metrics"

    return m.DecorateNext(pod)
}

// IstioSidecarDecorator Istio Sidecar 装饰器
type IstioSidecarDecorator struct {
    BasePodDecorator
    proxyImage    string
    proxyConfig   string
    includeIPRanges []string
}

func NewIstioSidecarDecorator(proxyImage, proxyConfig string) *IstioSidecarDecorator {
    return &IstioSidecarDecorator{
        proxyImage:      proxyImage,
        proxyConfig:     proxyConfig,
        includeIPRanges: []string{"*"},
    }
}

func (i *IstioSidecarDecorator) Name() string {
    return "istio-sidecar"
}

func (i *IstioSidecarDecorator) Decorate(pod *corev1.Pod) (*corev1.Pod, error) {
    // 添加 Init 容器用于配置 iptables
    initContainer := corev1.Container{
        Name:  "istio-init",
        Image: i.proxyImage,
        Args: []string{
            "istio-iptables",
            "-p", "15001",
            "-z", "15006",
            "-u", "1337",
            "-m", "REDIRECT",
            "-i", "*",
        },
        SecurityContext: &corev1.SecurityContext{
            Capabilities: &corev1.Capabilities{
                Add: []corev1.Capability{"NET_ADMIN"},
            },
            RunAsUser: int64Ptr(0),
        },
    }

    pod.Spec.InitContainers = append(pod.Spec.InitContainers, initContainer)

    // 添加 Envoy Sidecar
    envoySidecar := corev1.Container{
        Name:  "istio-proxy",
        Image: i.proxyImage,
        Args: []string{
            "proxy",
            "sidecar",
            "--configPath", "/etc/istio/proxy",
            "--binaryPath", "/usr/local/bin/envoy",
            "--serviceCluster", pod.Name,
            "--drainDuration", "45s",
            "--parentShutdownDuration", "1m0s",
        },
        Env: []corev1.EnvVar{
            {
                Name:  "ISTIO_META_INTERCEPTION_MODE",
                Value: "REDIRECT",
            },
            {
                Name: "POD_NAME",
                ValueFrom: &corev1.EnvVarSource{
                    FieldRef: &corev1.ObjectFieldSelector{
                        FieldPath: "metadata.name",
                    },
                },
            },
            {
                Name: "POD_NAMESPACE",
                ValueFrom: &corev1.EnvVarSource{
                    FieldRef: &corev1.ObjectFieldSelector{
                        FieldPath: "metadata.namespace",
                    },
                },
            },
        },
        VolumeMounts: []corev1.VolumeMount{
            {
                Name:      "istio-envoy",
                MountPath: "/etc/istio/proxy",
            },
            {
                Name:      "istio-certs",
                MountPath: "/etc/certs",
                ReadOnly:  true,
            },
        },
    }

    pod.Spec.Containers = append(pod.Spec.Containers, envoySidecar)

    return i.DecorateNext(pod)
}

// PodDecoratorChain 装饰器链
type PodDecoratorChain struct {
    decorators []PodDecorator
}

func NewPodDecoratorChain() *PodDecoratorChain {
    return &PodDecoratorChain{
        decorators: make([]PodDecorator, 0),
    }
}

func (c *PodDecoratorChain) AddDecorator(d PodDecorator) {
    c.decorators = append(c.decorators, d)
}

func (c *PodDecoratorChain) Decorate(pod *corev1.Pod) (*corev1.Pod, error) {
    var err error
    for _, decorator := range c.decorators {
        pod, err = decorator.Decorate(pod)
        if err != nil {
            return nil, err
        }
    }
    return pod, nil
}

// 辅助函数
func int64Ptr(i int64) *int64 {
    return &i
}

func boolPtr(b bool) *bool {
    return &b
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 比继承更灵活 | 产生很多小对象 |
| 可动态添加功能 | 调试困难 |
| 符合单一职责原则 | 装饰链过长影响性能 |

#### 适用场景

- 需要动态添加功能
- 不能或不想使用继承
- 需要透明地包装对象

---


## 2. 分布式系统设计模式

分布式系统设计模式是解决微服务架构中常见问题的最佳实践。Kubernetes 原生支持多种分布式模式。

---

### 2.1 Sidecar 模式

#### 意图

将应用程序的功能分解为独立的进程，作为 Sidecar 容器与主应用容器共同运行在同一 Pod 中。Sidecar 容器共享相同的网络和存储命名空间。

#### 结构

```
┌─────────────────────────────────────────────────────┐
│                      Pod                            │
│  ┌─────────────────┐  ┌─────────────────────────┐  │
│  │  Main Container │  │    Sidecar Container    │  │
│  │                 │  │                         │  │
│  │  ┌───────────┐  │  │  ┌─────────────────┐   │  │
│  │  │   App     │  │  │  │  Log Collector  │   │  │
│  │  │  Logic    │  │  │  │  / Monitoring   │   │  │
│  │  └───────────┘  │  │  │  / Proxy        │   │  │
│  │                 │  │  └─────────────────┘   │  │
│  └─────────────────┘  └─────────────────────────┘  │
│                                                     │
│  Shared: Network Namespace, Storage Volumes         │
└─────────────────────────────────────────────────────┘
```

#### 实现

```go
// Sidecar 模式实现
package sidecar

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "time"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// SidecarConfig Sidecar 配置
type SidecarConfig struct {
    Name          string
    Image         string
    Resources     corev1.ResourceRequirements
    VolumeMounts  []corev1.VolumeMount
    Env           []corev1.EnvVar
    Args          []string
    LivenessProbe *corev1.Probe
    ReadinessProbe *corev1.Probe
}

// LogCollectorSidecar 日志收集 Sidecar
type LogCollectorSidecar struct {
    config    *SidecarConfig
    logPath   string
    outputDir string
    client    kubernetes.Interface
}

func NewLogCollectorSidecar(client kubernetes.Interface, logPath, outputDir string) *LogCollectorSidecar {
    return &LogCollectorSidecar{
        client:    client,
        logPath:   logPath,
        outputDir: outputDir,
        config: &SidecarConfig{
            Name:  "log-collector",
            Image: "fluent/fluent-bit:1.9",
            Resources: corev1.ResourceRequirements{
                Limits: corev1.ResourceList{
                    corev1.ResourceCPU:    resource.MustParse("100m"),
                    corev1.ResourceMemory: resource.MustParse("128Mi"),
                },
                Requests: corev1.ResourceList{
                    corev1.ResourceCPU:    resource.MustParse("50m"),
                    corev1.ResourceMemory: resource.MustParse("64Mi"),
                },
            },
        },
    }
}

func (l *LogCollectorSidecar) CreateConfig() (*SidecarConfig, error) {
    configMapData := l.generateFluentBitConfig()

    // 创建 ConfigMap
    configMap := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: "fluent-bit-config",
        },
        Data: map[string]string{
            "fluent-bit.conf": configMapData,
        },
    }

    _, err := l.client.CoreV1().ConfigMaps("default").Create(context.TODO(), configMap, metav1.CreateOptions{})
    if err != nil {
        return nil, err
    }

    l.config.VolumeMounts = []corev1.VolumeMount{
        {
            Name:      "app-logs",
            MountPath: l.logPath,
            ReadOnly:  true,
        },
        {
            Name:      "fluent-bit-config",
            MountPath: "/fluent-bit/etc",
        },
    }

    l.config.Env = []corev1.EnvVar{
        {
            Name: "POD_NAME",
            ValueFrom: &corev1.EnvVarSource{
                FieldRef: &corev1.ObjectFieldSelector{
                    FieldPath: "metadata.name",
                },
            },
        },
        {
            Name: "NAMESPACE",
            ValueFrom: &corev1.EnvVarSource{
                FieldRef: &corev1.ObjectFieldSelector{
                    FieldPath: "metadata.namespace",
                },
            },
        },
    }

    return l.config, nil
}

func (l *LogCollectorSidecar) generateFluentBitConfig() string {
    return `[INPUT]
    Name              tail
    Tag               app.*
    Path              ` + l.logPath + `/*.log
    Parser            json
    DB                /tmp/flb_app.db
    Mem_Buf_Limit     5MB
    Skip_Long_Lines   On
    Refresh_Interval  10

[FILTER]
    Name              kubernetes
    Match             app.*
    Kube_URL          https://kubernetes.default.svc:443
    Kube_CA_File      /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    Kube_Token_File   /var/run/secrets/kubernetes.io/serviceaccount/token

[OUTPUT]
    Name              es
    Match             app.*
    Host              ${ELASTICSEARCH_HOST}
    Port              9200
    Index             app-logs
    Type              _doc`
}

// MetricsSidecar 监控指标 Sidecar
type MetricsSidecar struct {
    config   *SidecarConfig
    port     int32
    endpoint string
}

func NewMetricsSidecar(port int32, endpoint string) *MetricsSidecar {
    return &MetricsSidecar{
        port:     port,
        endpoint: endpoint,
        config: &SidecarConfig{
            Name:  "metrics-exporter",
            Image: "prom/node-exporter:v1.3.1",
            Resources: corev1.ResourceRequirements{
                Limits: corev1.ResourceList{
                    corev1.ResourceCPU:    resource.MustParse("100m"),
                    corev1.ResourceMemory: resource.MustParse("128Mi"),
                },
            },
            Args: []string{
                "--path.procfs=/host/proc",
                "--path.rootfs=/host/root",
                "--collector.filesystem.ignored-mount-points=^/(sys|proc|dev|host|etc)($$|/)",
            },
        },
    }
}

func (m *MetricsSidecar) CreateConfig() (*SidecarConfig, error) {
    m.config.Ports = []corev1.ContainerPort{
        {
            Name:          "metrics",
            ContainerPort: m.port,
            Protocol:      corev1.ProtocolTCP,
        },
    }

    m.config.VolumeMounts = []corev1.VolumeMount{
        {
            Name:      "proc",
            MountPath: "/host/proc",
            ReadOnly:  true,
        },
        {
            Name:      "sys",
            MountPath: "/host/sys",
            ReadOnly:  true,
        },
    }

    m.config.LivenessProbe = &corev1.Probe{
        ProbeHandler: corev1.ProbeHandler{
            HTTPGet: &corev1.HTTPGetAction{
                Path: "/metrics",
                Port: intstr.FromInt(int(m.port)),
            },
        },
        InitialDelaySeconds: 10,
        PeriodSeconds:       10,
    }

    return m.config, nil
}

// SidecarInjector Sidecar 注入器
type SidecarInjector struct {
    sidecars map[string]SidecarFactory
}

type SidecarFactory func() (Sidecar, error)

type Sidecar interface {
    CreateConfig() (*SidecarConfig, error)
    GetName() string
}

func NewSidecarInjector() *SidecarInjector {
    return &SidecarInjector{
        sidecars: make(map[string]SidecarFactory),
    }
}

func (s *SidecarInjector) Register(name string, factory SidecarFactory) {
    s.sidecars[name] = factory
}

func (s *SidecarInjector) Inject(pod *corev1.Pod, sidecarNames []string) error {
    for _, name := range sidecarNames {
        factory, ok := s.sidecars[name]
        if !ok {
            return fmt.Errorf("sidecar %s not found", name)
        }

        sidecar, err := factory()
        if err != nil {
            return err
        }

        config, err := sidecar.CreateConfig()
        if err != nil {
            return err
        }

        // 将 Sidecar 容器添加到 Pod
        container := corev1.Container{
            Name:           config.Name,
            Image:          config.Image,
            Resources:      config.Resources,
            VolumeMounts:   config.VolumeMounts,
            Env:            config.Env,
            Args:           config.Args,
            Ports:          config.Ports,
            LivenessProbe:  config.LivenessProbe,
            ReadinessProbe: config.ReadinessProbe,
        }

        pod.Spec.Containers = append(pod.Spec.Containers, container)
    }

    return nil
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 语言无关 | 资源开销增加 |
| 独立部署升级 | Pod 启动时间增加 |
| 复用性强 | 容器间通信复杂度 |
| 关注点分离 | 共享资源竞争 |

#### 适用场景

- 日志收集和聚合
- 监控和指标收集
- 配置同步
- 服务网格代理

---

### 2.2 Ambassador 模式

#### 意图

Ambassador 模式使用一个代理容器作为外部服务的本地代表，处理连接管理、重试、断路等复杂逻辑，让主应用容器专注于业务逻辑。

#### 结构

```
┌─────────────────────────────────────────────────────────┐
│                          Pod                            │
│  ┌──────────────────┐    ┌─────────────────────────┐   │
│  │  Main Container  │◄──►│   Ambassador Container  │   │
│  │                  │    │                         │   │
│  │  ┌────────────┐  │    │  ┌─────────────────┐   │   │
│  │  │   App      │  │    │  │  Connection     │   │   │
│  │  │  Logic     │  │    │  │  Pool Manager   │   │   │
│  │  │            │  │    │  │  Retry Logic    │   │   │
│  │  └────────────┘  │    │  │  Load Balancer  │   │   │
│  │                  │    │  └─────────────────┘   │   │
│  └──────────────────┘    └───────────┬─────────────┘   │
│                                      │                  │
└──────────────────────────────────────┼──────────────────┘
                                       │
                                       ▼
                            ┌──────────────────────┐
                            │   External Service   │
                            │  (Redis/MySQL/Kafka) │
                            └──────────────────────┘
```

#### 实现

```go
// Ambassador 模式实现
package ambassador

import (
    "context"
    "fmt"
    "net"
    "sync"
    "time"
)

// ConnectionPool 连接池接口
type ConnectionPool interface {
    Get() (net.Conn, error)
    Put(conn net.Conn) error
    Close() error
}

// Ambassador Ambassador 代理
type Ambassador struct {
    targetAddr    string
    localAddr     string
    pool          ConnectionPool
    retryPolicy   *RetryPolicy
    circuitBreaker *CircuitBreaker
    mu            sync.RWMutex
}

type RetryPolicy struct {
    MaxRetries  int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Multiplier  float64
}

func NewAmbassador(localAddr, targetAddr string, pool ConnectionPool) *Ambassador {
    return &Ambassador{
        targetAddr:     targetAddr,
        localAddr:      localAddr,
        pool:           pool,
        retryPolicy:    defaultRetryPolicy(),
        circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
    }
}

func defaultRetryPolicy() *RetryPolicy {
    return &RetryPolicy{
        MaxRetries: 3,
        BaseDelay:  100 * time.Millisecond,
        MaxDelay:   5 * time.Second,
        Multiplier: 2.0,
    }
}

// Start 启动 Ambassador 代理
func (a *Ambassador) Start(ctx context.Context) error {
    listener, err := net.Listen("tcp", a.localAddr)
    if err != nil {
        return fmt.Errorf("failed to listen on %s: %w", a.localAddr, err)
    }
    defer listener.Close()

    fmt.Printf("Ambassador listening on %s, proxying to %s\n", a.localAddr, a.targetAddr)

    go func() {
        <-ctx.Done()
        listener.Close()
    }()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        conn, err := listener.Accept()
        if err != nil {
            if ctx.Err() != nil {
                return ctx.Err()
            }
            continue
        }

        go a.handleConnection(ctx, conn)
    }
}

func (a *Ambassador) handleConnection(ctx context.Context, clientConn net.Conn) {
    defer clientConn.Close()

    // 检查断路器状态
    if a.circuitBreaker.IsOpen() {
        fmt.Println("Circuit breaker is open, rejecting connection")
        return
    }

    // 获取后端连接（带重试）
    backendConn, err := a.getBackendConnectionWithRetry(ctx)
    if err != nil {
        a.circuitBreaker.RecordFailure()
        fmt.Printf("Failed to get backend connection: %v\n", err)
        return
    }
    defer backendConn.Close()

    a.circuitBreaker.RecordSuccess()

    // 双向数据转发
    errChan := make(chan error, 2)

    go func() {
        _, err := copyData(backendConn, clientConn)
        errChan <- err
    }()

    go func() {
        _, err := copyData(clientConn, backendConn)
        errChan <- err
    }()

    // 等待任意一个方向完成
    <-errChan
}

func (a *Ambassador) getBackendConnectionWithRetry(ctx context.Context) (net.Conn, error) {
    var lastErr error
    delay := a.retryPolicy.BaseDelay

    for i := 0; i <= a.retryPolicy.MaxRetries; i++ {
        if i > 0 {
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(delay):
                delay = time.Duration(float64(delay) * a.retryPolicy.Multiplier)
                if delay > a.retryPolicy.MaxDelay {
                    delay = a.retryPolicy.MaxDelay
                }
            }
        }

        conn, err := a.pool.Get()
        if err == nil {
            return conn, nil
        }

        lastErr = err
        fmt.Printf("Retry %d: failed to get connection: %v\n", i, err)
    }

    return nil, fmt.Errorf("exhausted retries: %w", lastErr)
}

func copyData(dst, src net.Conn) (int64, error) {
    buf := make([]byte, 32*1024)
    var written int64

    for {
        nr, err := src.Read(buf)
        if nr > 0 {
            nw, ew := dst.Write(buf[0:nr])
            if nw > 0 {
                written += int64(nw)
            }
            if ew != nil {
                return written, ew
            }
            if nr != nw {
                return written, fmt.Errorf("short write")
            }
        }
        if err != nil {
            return written, err
        }
    }
}

// CircuitBreaker 断路器
type CircuitBreaker struct {
    failureThreshold int
    resetTimeout     time.Duration
    failures         int
    lastFailureTime  time.Time
    state            CircuitState
    mu               sync.RWMutex
}

type CircuitState int

const (
    StateClosed CircuitState = iota
    StateOpen
    StateHalfOpen
)

func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        failureThreshold: threshold,
        resetTimeout:     resetTimeout,
        state:            StateClosed,
    }
}

func (cb *CircuitBreaker) IsOpen() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    if cb.state == StateOpen {
        if time.Since(cb.lastFailureTime) > cb.resetTimeout {
            cb.mu.RUnlock()
            cb.mu.Lock()
            cb.state = StateHalfOpen
            cb.mu.Unlock()
            cb.mu.RLock()
            return false
        }
        return true
    }

    return false
}

func (cb *CircuitBreaker) RecordSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.failures = 0
    cb.state = StateClosed
}

func (cb *CircuitBreaker) RecordFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.failures++
    cb.lastFailureTime = time.Now()

    if cb.failures >= cb.failureThreshold {
        cb.state = StateOpen
    }
}

// RedisAmbassador Redis 代理
type RedisAmbassador struct {
    *Ambassador
    sentinelAddrs []string
    masterName    string
}

func NewRedisAmbassador(localAddr string, sentinelAddrs []string, masterName string) *RedisAmbassador {
    pool := NewRedisConnectionPool(sentinelAddrs, masterName)

    return &RedisAmbassador{
        Ambassador:    NewAmbassador(localAddr, "", pool),
        sentinelAddrs: sentinelAddrs,
        masterName:    masterName,
    }
}

// MySQLAmbassador MySQL 代理
type MySQLAmbassador struct {
    *Ambassador
    readReplicas []string
    primary      string
}

func NewMySQLAmbassador(localAddr, primary string, readReplicas []string) *MySQLAmbassador {
    pool := NewMySQLConnectionPool(primary, readReplicas)

    return &MySQLAmbassador{
        Ambassador:   NewAmbassador(localAddr, primary, pool),
        readReplicas: readReplicas,
        primary:      primary,
    }
}

// 连接池实现

// SimpleConnectionPool 简单连接池
type SimpleConnectionPool struct {
    targetAddr string
    maxConns   int
    conns      chan net.Conn
    mu         sync.Mutex
}

func NewSimpleConnectionPool(targetAddr string, maxConns int) *SimpleConnectionPool {
    return &SimpleConnectionPool{
        targetAddr: targetAddr,
        maxConns:   maxConns,
        conns:      make(chan net.Conn, maxConns),
    }
}

func (p *SimpleConnectionPool) Get() (net.Conn, error) {
    select {
    case conn := <-p.conns:
        return conn, nil
    default:
        return net.Dial("tcp", p.targetAddr)
    }
}

func (p *SimpleConnectionPool) Put(conn net.Conn) error {
    select {
    case p.conns <- conn:
        return nil
    default:
        return conn.Close()
    }
}

func (p *SimpleConnectionPool) Close() error {
    close(p.conns)
    for conn := range p.conns {
        conn.Close()
    }
    return nil
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 简化应用代码 | 增加网络跳数 |
| 统一连接管理 | 额外的资源消耗 |
| 支持高级功能 | 单点故障风险 |
| 透明代理 | 调试复杂度增加 |

#### 适用场景

- 数据库连接池管理
- 外部服务代理
- 多语言环境服务发现
- 遗留系统现代化

---

### 2.3 Adapter 模式（分布式版本）

#### 意图

在分布式系统中，Adapter 模式用于处理不同服务之间的协议转换、数据格式转换和接口适配。

#### 结构

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Service A  │────►│   Adapter    │────►│   Service B  │
│  (REST API)  │     │  (Protocol   │     │  (gRPC API)  │
└──────────────┘     │  Converter)  │     └──────────────┘
                     └──────────────┘
```

#### 实现

```go
// Adapter 模式实现
package adapter

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "google.golang.org/grpc"
    "google.golang.org/protobuf/types/known/emptypb"
)

// ProtocolAdapter 协议适配器接口
type ProtocolAdapter interface {
    ConvertRequest(src interface{}) (interface{}, error)
    ConvertResponse(src interface{}) (interface{}, error)
}

// RESTToGRPCAdapter REST 到 gRPC 适配器
type RESTToGRPCAdapter struct {
    grpcClient *grpc.ClientConn
}

func NewRESTToGRPCAdapter(grpcAddr string) (*RESTToGRPCAdapter, error) {
    conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    return &RESTToGRPCAdapter{grpcClient: conn}, nil
}

// HTTPHandler HTTP 处理器
func (a *RESTToGRPCAdapter) HTTPHandler(w http.ResponseWriter, r *http.Request) {
    // 解析 REST 请求
    var restReq RESTRequest
    if err := json.NewDecoder(r.Body).Decode(&restReq); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 转换为 gRPC 请求
    grpcReq, err := a.ConvertRequest(&restReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 调用 gRPC 服务
    grpcResp, err := a.callGRPCService(r.Context(), grpcReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 转换响应
    restResp, err := a.ConvertResponse(grpcResp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 返回 REST 响应
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(restResp)
}

type RESTRequest struct {
    UserID   string                 `json:"user_id"`
    Action   string                 `json:"action"`
    Payload  map[string]interface{} `json:"payload"`
    Metadata map[string]string      `json:"metadata"`
}

type RESTResponse struct {
    Success bool                   `json:"success"`
    Data    map[string]interface{} `json:"data,omitempty"`
    Error   string                 `json:"error,omitempty"`
}

func (a *RESTToGRPCAdapter) ConvertRequest(src interface{}) (interface{}, error) {
    restReq := src.(*RESTRequest)

    // 转换为 gRPC 请求
    grpcReq := &GRPCRequest{
        UserId:   restReq.UserID,
        Action:   restReq.Action,
        Payload:  mustMarshal(restReq.Payload),
        Metadata: restReq.Metadata,
    }

    return grpcReq, nil
}

func (a *RESTToGRPCAdapter) ConvertResponse(src interface{}) (interface{}, error) {
    grpcResp := src.(*GRPCResponse)

    var data map[string]interface{}
    if err := json.Unmarshal(grpcResp.Data, &data); err != nil {
        return nil, err
    }

    restResp := &RESTResponse{
        Success: grpcResp.Success,
        Data:    data,
        Error:   grpcResp.Error,
    }

    return restResp, nil
}

func (a *RESTToGRPCAdapter) callGRPCService(ctx context.Context, req interface{}) (*GRPCResponse, error) {
    // 实际调用 gRPC 服务
    grpcReq := req.(*GRPCRequest)

    // 这里使用示例 gRPC 客户端
    client := NewExampleGRPCClient(a.grpcClient)
    return client.Process(ctx, grpcReq)
}

// 数据格式转换器

// DataFormatConverter 数据格式转换器
type DataFormatConverter struct {
    converters map[string]Converter
}

type Converter interface {
    Convert(data []byte) ([]byte, error)
    ContentType() string
}

func NewDataFormatConverter() *DataFormatConverter {
    return &DataFormatConverter{
        converters: make(map[string]Converter),
    }
}

func (d *DataFormatConverter) Register(format string, converter Converter) {
    d.converters[format] = converter
}

func (d *DataFormatConverter) Convert(fromFormat, toFormat string, data []byte) ([]byte, error) {
    fromConverter, ok := d.converters[fromFormat]
    if !ok {
        return nil, fmt.Errorf("unsupported source format: %s", fromFormat)
    }

    toConverter, ok := d.converters[toFormat]
    if !ok {
        return nil, fmt.Errorf("unsupported target format: %s", toFormat)
    }

    // 先转换为中间格式（如 map[string]interface{}）
    intermediate, err := fromConverter.Convert(data)
    if err != nil {
        return nil, err
    }

    // 再转换为目标格式
    return toConverter.Convert(intermediate)
}

// JSONConverter JSON 转换器
type JSONConverter struct{}

func (j *JSONConverter) Convert(data []byte) ([]byte, error) {
    // JSON 已经是目标格式，直接返回
    return data, nil
}

func (j *JSONConverter) ContentType() string {
    return "application/json"
}

// XMLConverter XML 转换器
type XMLConverter struct{}

func (x *XMLConverter) Convert(data []byte) ([]byte, error) {
    var result map[string]interface{}
    if err := xml.Unmarshal(data, &result); err != nil {
        return nil, err
    }
    return json.Marshal(result)
}

func (x *XMLConverter) ContentType() string {
    return "application/xml"
}

// ProtocolBufferConverter Protobuf 转换器
type ProtocolBufferConverter struct {
    message proto.Message
}

func (p *ProtocolBufferConverter) Convert(data []byte) ([]byte, error) {
    if err := proto.Unmarshal(data, p.message); err != nil {
        return nil, err
    }

    // 使用 protojson 转换为 JSON
    return protojson.Marshal(p.message)
}

func (p *ProtocolBufferConverter) ContentType() string {
    return "application/x-protobuf"
}

// 辅助函数
func mustMarshal(v interface{}) []byte {
    data, _ := json.Marshal(v)
    return data
}

// GRPCRequest gRPC 请求
type GRPCRequest struct {
    UserId   string
    Action   string
    Payload  []byte
    Metadata map[string]string
}

// GRPCResponse gRPC 响应
type GRPCResponse struct {
    Success bool
    Data    []byte
    Error   string
}

// ExampleGRPCClient 示例 gRPC 客户端
type ExampleGRPCClient struct {
    conn *grpc.ClientConn
}

func NewExampleGRPCClient(conn *grpc.ClientConn) *ExampleGRPCClient {
    return &ExampleGRPCClient{conn: conn}
}

func (c *ExampleGRPCClient) Process(ctx context.Context, req *GRPCRequest) (*GRPCResponse, error) {
    // 实际 gRPC 调用实现
    return &GRPCResponse{
        Success: true,
        Data:    []byte(`{"result": "ok"}`),
    }, nil
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 解耦服务间依赖 | 增加延迟 |
| 支持异构系统集成 | 额外的维护成本 |
| 可独立演进 | 数据丢失风险 |

#### 适用场景

- 微服务间协议转换
- 遗留系统集成
- 多协议支持

---

### 2.4 Scatter-Gather 模式

#### 意图

将请求分散（Scatter）到多个服务实例并行处理，然后收集（Gather）所有响应并聚合结果。

#### 结构

```
                         ┌──────────────┐
                         │   Client     │
                         └──────┬───────┘
                                │
                    ┌───────────┼───────────┐
                    │           │           │
                    ▼           ▼           ▼
            ┌──────────┐ ┌──────────┐ ┌──────────┐
            │ Service  │ │ Service  │ │ Service  │
            │    A     │ │    B     │ │    C     │
            └────┬─────┘ └────┬─────┘ └────┬─────┘
                 │            │            │
                 └────────────┼────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Result Aggregator│
                    └──────────────────┘
```

#### 实现

```go
// Scatter-Gather 模式实现
package scattergather

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Request 请求接口
type Request interface {
    GetKey() string
}

// Response 响应接口
type Response interface {
    GetKey() string
    GetData() interface{}
    GetError() error
}

// Service 服务接口
type Service interface {
    Name() string
    Process(ctx context.Context, req Request) (Response, error)
}

// ScatterGather Scatter-Gather 处理器
type ScatterGather struct {
    services []Service
    timeout  time.Duration
    strategy AggregationStrategy
}

// AggregationStrategy 聚合策略
type AggregationStrategy interface {
    Aggregate(responses []Response) (Response, error)
}

// NewScatterGather 创建 ScatterGather 实例
func NewScatterGather(timeout time.Duration, strategy AggregationStrategy) *ScatterGather {
    return &ScatterGather{
        services: make([]Service, 0),
        timeout:  timeout,
        strategy: strategy,
    }
}

func (sg *ScatterGather) AddService(service Service) {
    sg.services = append(sg.services, service)
}

// Execute 执行 Scatter-Gather
func (sg *ScatterGather) Execute(ctx context.Context, req Request) (Response, error) {
    ctx, cancel := context.WithTimeout(ctx, sg.timeout)
    defer cancel()

    responseChan := make(chan Response, len(sg.services))
    var wg sync.WaitGroup

    // Scatter: 并行发送到所有服务
    for _, service := range sg.services {
        wg.Add(1)
        go func(s Service) {
            defer wg.Done()

            resp, err := s.Process(ctx, req)
            if err != nil {
                resp = &ErrorResponse{
                    key:   req.GetKey(),
                    error: err,
                }
            }

            select {
            case responseChan <- resp:
            case <-ctx.Done():
            }
        }(service)
    }

    // 等待所有服务完成
    go func() {
        wg.Wait()
        close(responseChan)
    }()

    // Gather: 收集响应
    var responses []Response
    for resp := range responseChan {
        responses = append(responses, resp)
    }

    // 聚合结果
    return sg.strategy.Aggregate(responses)
}

// FirstSuccessStrategy 第一个成功响应策略
type FirstSuccessStrategy struct{}

func (f *FirstSuccessStrategy) Aggregate(responses []Response) (Response, error) {
    for _, resp := range responses {
        if resp.GetError() == nil {
            return resp, nil
        }
    }
    return nil, fmt.Errorf("all services failed")
}

// MajorityVotingStrategy 多数投票策略
type MajorityVotingStrategy struct{}

func (m *MajorityVotingStrategy) Aggregate(responses []Response) (Response, error) {
    votes := make(map[string]int)
    var successfulResponses []Response

    for _, resp := range responses {
        if resp.GetError() != nil {
            continue
        }
        successfulResponses = append(successfulResponses, resp)

        key := fmt.Sprintf("%v", resp.GetData())
        votes[key]++
    }

    if len(successfulResponses) == 0 {
        return nil, fmt.Errorf("no successful responses")
    }

    // 找到多数票
    maxVotes := 0
    var majorityResponse Response
    for _, resp := range successfulResponses {
        key := fmt.Sprintf("%v", resp.GetData())
        if votes[key] > maxVotes {
            maxVotes = votes[key]
            majorityResponse = resp
        }
    }

    // 检查是否达到多数
    if maxVotes > len(successfulResponses)/2 {
        return majorityResponse, nil
    }

    return nil, fmt.Errorf("no majority consensus")
}

// WeightedAverageStrategy 加权平均策略
type WeightedAverageStrategy struct {
    weights map[string]float64
}

func NewWeightedAverageStrategy(weights map[string]float64) *WeightedAverageStrategy {
    return &WeightedAverageStrategy{weights: weights}
}

func (w *WeightedAverageStrategy) Aggregate(responses []Response) (Response, error) {
    var sum float64
    var totalWeight float64

    for _, resp := range responses {
        if resp.GetError() != nil {
            continue
        }

        weight := w.weights[resp.GetKey()]
        if weight == 0 {
            weight = 1.0
        }

        value, ok := resp.GetData().(float64)
        if !ok {
            continue
        }

        sum += value * weight
        totalWeight += weight
    }

    if totalWeight == 0 {
        return nil, fmt.Errorf("no valid responses for weighted average")
    }

    return &SimpleResponse{
        key:  "aggregated",
        data: sum / totalWeight,
    }, nil
}

// 具体实现

// SearchService 搜索服务
type SearchService struct {
    name   string
    source string
}

func NewSearchService(name, source string) *SearchService {
    return &SearchService{name: name, source: source}
}

func (s *SearchService) Name() string {
    return s.name
}

func (s *SearchService) Process(ctx context.Context, req Request) (Response, error) {
    searchReq := req.(*SearchRequest)

    // 模拟搜索
    select {
    case <-time.After(100 * time.Millisecond):
        return &SearchResponse{
            key:    s.name,
            source: s.source,
            results: []SearchResult{
                {Title: fmt.Sprintf("Result from %s", s.source), Score: 0.9},
            },
        }, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

type SearchRequest struct {
    query string
}

func (s *SearchRequest) GetKey() string {
    return s.query
}

type SearchResponse struct {
    key     string
    source  string
    results []SearchResult
}

func (s *SearchResponse) GetKey() string    { return s.key }
func (s *SearchResponse) GetData() interface{} { return s.results }
func (s *SearchResponse) GetError() error   { return nil }

type SearchResult struct {
    Title string
    Score float64
}

// SimpleResponse 简单响应
type SimpleResponse struct {
    key   string
    data  interface{}
    error error
}

func (s *SimpleResponse) GetKey() string      { return s.key }
func (s *SimpleResponse) GetData() interface{} { return s.data }
func (s *SimpleResponse) GetError() error     { return s.error }

// ErrorResponse 错误响应
type ErrorResponse struct {
    key   string
    error error
}

func (e *ErrorResponse) GetKey() string      { return e.key }
func (e *ErrorResponse) GetData() interface{} { return nil }
func (e *ErrorResponse) GetError() error     { return e.error }

// 使用示例
func ExampleScatterGather() {
    // 创建聚合策略
    strategy := &FirstSuccessStrategy{}

    // 创建 ScatterGather
    sg := NewScatterGather(5*time.Second, strategy)

    // 添加搜索服务
    sg.AddService(NewSearchService("elasticsearch", "es"))
    sg.AddService(NewSearchService("solr", "solr"))
    sg.AddService(NewSearchService("meilisearch", "meilisearch"))

    // 执行搜索
    ctx := context.Background()
    req := &SearchRequest{query: "kubernetes patterns"}

    resp, err := sg.Execute(ctx, req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Result: %v\n", resp.GetData())
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 提高响应速度 | 增加系统复杂度 |
| 提高可用性 | 资源消耗增加 |
| 支持多种聚合策略 | 结果一致性挑战 |

#### 适用场景

- 搜索引擎聚合
- 价格比较服务
- 分布式计算

---

### 2.5 Saga 模式

#### 意图

Saga 模式用于管理分布式事务，将长事务拆分为多个本地事务，每个本地事务有对应的补偿操作，以确保数据最终一致性。

#### 结构

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│  Step 1 │───►│  Step 2 │───►│  Step 3 │───►│  Step 4 │
│  (T1)   │    │  (T2)   │    │  (T3)   │    │  (T4)   │
│  [C1]   │    │  [C2]   │    │  [C3]   │    │  [C4]   │
└─────────┘    └────┬────┘    └─────────┘    └─────────┘
                    │
                    │ (失败)
                    ▼
              ┌─────────┐
              │ Compensate 2 │
              │ Compensate 1 │
              └─────────┘

T = Transaction, C = Compensation
```

#### 实现

```go
// Saga 模式实现
package saga

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// SagaState Saga 状态
type SagaState int

const (
    SagaPending SagaState = iota
    SagaRunning
    SagaCompleted
    SagaCompensating
    SagaCompensated
    SagaFailed
)

// SagaLog Saga 日志
type SagaLog struct {
    StepIndex   int
    Action      string // "execute" or "compensate"
    Success     bool
    Timestamp   time.Time
    Error       string
}

// SagaStep Saga 步骤
type SagaStep struct {
    Name        string
    Execute     func(ctx context.Context) error
    Compensate  func(ctx context.Context) error
    RetryCount  int
    RetryDelay  time.Duration
}

// Saga Saga 协调器
type Saga struct {
    name     string
    steps    []*SagaStep
    state    SagaState
    logs     []SagaLog
    mu       sync.RWMutex
    onError  func(stepIndex int, err error)
}

// NewSaga 创建 Saga
func NewSaga(name string) *Saga {
    return &Saga{
        name:  name,
        steps: make([]*SagaStep, 0),
        state: SagaPending,
        logs:  make([]SagaLog, 0),
    }
}

func (s *Saga) AddStep(step *SagaStep) {
    s.steps = append(s.steps, step)
}

func (s *Saga) SetOnError(handler func(stepIndex int, err error)) {
    s.onError = handler
}

// Execute 执行 Saga
func (s *Saga) Execute(ctx context.Context) error {
    s.mu.Lock()
    s.state = SagaRunning
    s.mu.Unlock()

    for i, step := range s.steps {
        // 执行步骤
        err := s.executeStep(ctx, i, step)
        if err != nil {
            // 执行失败，开始补偿
            s.log(SagaLog{
                StepIndex: i,
                Action:    "execute",
                Success:   false,
                Timestamp: time.Now(),
                Error:     err.Error(),
            })

            if s.onError != nil {
                s.onError(i, err)
            }

            // 执行补偿
            return s.compensate(ctx, i)
        }

        s.log(SagaLog{
            StepIndex: i,
            Action:    "execute",
            Success:   true,
            Timestamp: time.Now(),
        })
    }

    s.mu.Lock()
    s.state = SagaCompleted
    s.mu.Unlock()

    return nil
}

func (s *Saga) executeStep(ctx context.Context, index int, step *SagaStep) error {
    var err error

    for attempt := 0; attempt <= step.RetryCount; attempt++ {
        if attempt > 0 {
            select {
            case <-time.After(step.RetryDelay):
            case <-ctx.Done():
                return ctx.Err()
            }
        }

        err = step.Execute(ctx)
        if err == nil {
            return nil
        }
    }

    return err
}

func (s *Saga) compensate(ctx context.Context, failedIndex int) error {
    s.mu.Lock()
    s.state = SagaCompensating
    s.mu.Unlock()

    // 逆向补偿已执行的步骤
    for i := failedIndex - 1; i >= 0; i-- {
        step := s.steps[i]

        if step.Compensate != nil {
            err := step.Compensate(ctx)
            s.log(SagaLog{
                StepIndex: i,
                Action:    "compensate",
                Success:   err == nil,
                Timestamp: time.Now(),
                Error:     func() string { if err != nil { return err.Error() }; return "" }(),
            })

            if err != nil {
                // 补偿失败，需要人工介入
                s.mu.Lock()
                s.state = SagaFailed
                s.mu.Unlock()
                return fmt.Errorf("compensation failed at step %d: %w", i, err)
            }
        }
    }

    s.mu.Lock()
    s.state = SagaCompensated
    s.mu.Unlock()

    return fmt.Errorf("saga compensated due to failure at step %d", failedIndex)
}

func (s *Saga) log(entry SagaLog) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.logs = append(s.logs, entry)
}

func (s *Saga) GetState() SagaState {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.state
}

func (s *Saga) GetLogs() []SagaLog {
    s.mu.RLock()
    defer s.mu.RUnlock()
    logs := make([]SagaLog, len(s.logs))
    copy(logs, s.logs)
    return logs
}

// 订单 Saga 示例

// OrderSaga 订单处理 Saga
type OrderSaga struct {
    *Saga
    orderID   string
    userID    string
    amount    float64
    inventory map[string]int
}

func NewOrderSaga(orderID, userID string, amount float64, inventory map[string]int) *OrderSaga {
    saga := &OrderSaga{
        orderID:   orderID,
        userID:    userID,
        amount:    amount,
        inventory: inventory,
    }

    saga.Saga = NewSaga(fmt.Sprintf("order-%s", orderID))

    // 步骤1: 创建订单
    saga.AddStep(&SagaStep{
        Name: "CreateOrder",
        Execute: func(ctx context.Context) error {
            return saga.createOrder(ctx)
        },
        Compensate: func(ctx context.Context) error {
            return saga.cancelOrder(ctx)
        },
        RetryCount: 3,
        RetryDelay: 1 * time.Second,
    })

    // 步骤2: 扣减库存
    saga.AddStep(&SagaStep{
        Name: "DeductInventory",
        Execute: func(ctx context.Context) error {
            return saga.deductInventory(ctx)
        },
        Compensate: func(ctx context.Context) error {
            return saga.restoreInventory(ctx)
        },
        RetryCount: 3,
        RetryDelay: 1 * time.Second,
    })

    // 步骤3: 扣减余额
    saga.AddStep(&SagaStep{
        Name: "DeductBalance",
        Execute: func(ctx context.Context) error {
            return saga.deductBalance(ctx)
        },
        Compensate: func(ctx context.Context) error {
            return saga.refundBalance(ctx)
        },
        RetryCount: 3,
        RetryDelay: 1 * time.Second,
    })

    // 步骤4: 发送通知
    saga.AddStep(&SagaStep{
        Name: "SendNotification",
        Execute: func(ctx context.Context) error {
            return saga.sendNotification(ctx)
        },
        // 通知不需要补偿
        Compensate: nil,
        RetryCount: 3,
        RetryDelay: 500 * time.Millisecond,
    })

    return saga
}

func (o *OrderSaga) createOrder(ctx context.Context) error {
    fmt.Printf("Creating order %s for user %s\n", o.orderID, o.userID)
    // 调用订单服务创建订单
    return nil
}

func (o *OrderSaga) cancelOrder(ctx context.Context) error {
    fmt.Printf("Canceling order %s\n", o.orderID)
    // 调用订单服务取消订单
    return nil
}

func (o *OrderSaga) deductInventory(ctx context.Context) error {
    fmt.Printf("Deducting inventory for order %s\n", o.orderID)
    // 调用库存服务扣减库存
    return nil
}

func (o *OrderSaga) restoreInventory(ctx context.Context) error {
    fmt.Printf("Restoring inventory for order %s\n", o.orderID)
    // 调用库存服务恢复库存
    return nil
}

func (o *OrderSaga) deductBalance(ctx context.Context) error {
    fmt.Printf("Deducting balance for user %s: %.2f\n", o.userID, o.amount)
    // 调用支付服务扣减余额
    return nil
}

func (o *OrderSaga) refundBalance(ctx context.Context) error {
    fmt.Printf("Refunding balance for user %s: %.2f\n", o.userID, o.amount)
    // 调用支付服务退款
    return nil
}

func (o *OrderSaga) sendNotification(ctx context.Context) error {
    fmt.Printf("Sending notification for order %s\n", o.orderID)
    // 调用通知服务发送通知
    return nil
}

// SagaOrchestrator Saga 编排器
type SagaOrchestrator struct {
    sagas map[string]*Saga
    mu    sync.RWMutex
}

func NewSagaOrchestrator() *SagaOrchestrator {
    return &SagaOrchestrator{
        sagas: make(map[string]*Saga),
    }
}

func (o *SagaOrchestrator) RegisterSaga(saga *Saga) {
    o.mu.Lock()
    defer o.mu.Unlock()
    o.sagas[saga.name] = saga
}

func (o *SagaOrchestrator) ExecuteSaga(ctx context.Context, name string) error {
    o.mu.RLock()
    saga, ok := o.sagas[name]
    o.mu.RUnlock()

    if !ok {
        return fmt.Errorf("saga %s not found", name)
    }

    return saga.Execute(ctx)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 避免长事务锁定 | 实现复杂 |
| 支持异步执行 | 补偿逻辑难设计 |
| 提高系统可用性 | 最终一致性 |

#### 适用场景

- 电商订单处理
- 金融转账
- 分布式工作流

---



### 2.6 Circuit Breaker 模式

#### 意图

Circuit Breaker 模式用于防止故障扩散，当检测到服务故障时，快速失败而不是持续尝试可能导致失败的操作。

#### 结构

```
    ┌─────────────┐
    │   CLOSED    │ ──► 正常处理请求
    │  (正常状态)  │
    └──────┬──────┘
           │ 失败次数超过阈值
           ▼
    ┌─────────────┐
    │    OPEN     │ ──► 快速失败，不调用服务
    │  (熔断状态)  │
    └──────┬──────┘
           │ 超时时间到达
           ▼
    ┌─────────────┐
    │  HALF-OPEN  │ ──► 允许有限请求测试服务
    │ (半开状态)   │
    └─────────────┘
```

#### 实现

```go
// Circuit Breaker 模式实现
package circuitbreaker

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"
)

// State 断路器状态
type State int

const (
    StateClosed State = iota    // 关闭 - 正常处理
    StateOpen                   // 打开 - 快速失败
    StateHalfOpen               // 半开 - 测试服务恢复
)

func (s State) String() string {
    switch s {
    case StateClosed:
        return "CLOSED"
    case StateOpen:
        return "OPEN"
    case StateHalfOpen:
        return "HALF_OPEN"
    default:
        return "UNKNOWN"
    }
}

// Config 断路器配置
type Config struct {
    FailureThreshold    int           // 失败阈值
    SuccessThreshold    int           // 成功阈值（半开状态）
    Timeout             time.Duration // 熔断超时时间
    HalfOpenMaxRequests int           // 半开状态最大请求数
    OnStateChange       func(from, to State)
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
    return &Config{
        FailureThreshold:    5,
        SuccessThreshold:    3,
        Timeout:             30 * time.Second,
        HalfOpenMaxRequests: 3,
    }
}

// CircuitBreaker 断路器
type CircuitBreaker struct {
    name           string
    config         *Config
    state          State
    failures       int
    successes      int
    consecutiveSuccesses int
    lastFailureTime time.Time
    halfOpenRequests int
    mu             sync.RWMutex
}

// ErrCircuitOpen 断路器打开错误
var ErrCircuitOpen = errors.New("circuit breaker is open")

// NewCircuitBreaker 创建断路器
func NewCircuitBreaker(name string, config *Config) *CircuitBreaker {
    if config == nil {
        config = DefaultConfig()
    }

    return &CircuitBreaker{
        name:   name,
        config: config,
        state:  StateClosed,
    }
}

// Execute 执行受保护的操作
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
    if err := cb.beforeRequest(); err != nil {
        return err
    }

    err := fn()
    cb.afterRequest(err)

    return err
}

// ExecuteWithResult 执行带返回值的受保护操作
func (cb *CircuitBreaker) ExecuteWithResult(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
    if err := cb.beforeRequest(); err != nil {
        return nil, err
    }

    result, err := fn()
    cb.afterRequest(err)

    return result, err
}

func (cb *CircuitBreaker) beforeRequest() error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    switch cb.state {
    case StateClosed:
        return nil

    case StateOpen:
        if time.Since(cb.lastFailureTime) > cb.config.Timeout {
            cb.transitionTo(StateHalfOpen)
            cb.halfOpenRequests = 0
            return nil
        }
        return ErrCircuitOpen

    case StateHalfOpen:
        if cb.halfOpenRequests >= cb.config.HalfOpenMaxRequests {
            return ErrCircuitOpen
        }
        cb.halfOpenRequests++
        return nil
    }

    return nil
}

func (cb *CircuitBreaker) afterRequest(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err == nil {
        cb.onSuccess()
    } else {
        cb.onFailure()
    }
}

func (cb *CircuitBreaker) onSuccess() {
    switch cb.state {
    case StateClosed:
        cb.failures = 0

    case StateHalfOpen:
        cb.consecutiveSuccesses++
        if cb.consecutiveSuccesses >= cb.config.SuccessThreshold {
            cb.transitionTo(StateClosed)
            cb.failures = 0
            cb.consecutiveSuccesses = 0
        }
    }
}

func (cb *CircuitBreaker) onFailure() {
    cb.failures++
    cb.lastFailureTime = time.Now()

    switch cb.state {
    case StateClosed:
        if cb.failures >= cb.config.FailureThreshold {
            cb.transitionTo(StateOpen)
        }

    case StateHalfOpen:
        cb.transitionTo(StateOpen)
    }
}

func (cb *CircuitBreaker) transitionTo(newState State) {
    oldState := cb.state
    cb.state = newState

    if cb.config.OnStateChange != nil {
        cb.config.OnStateChange(oldState, newState)
    }

    fmt.Printf("CircuitBreaker %s: %s -> %s\n", cb.name, oldState, newState)
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() State {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    return cb.state
}

// GetMetrics 获取指标
func (cb *CircuitBreaker) GetMetrics() map[string]interface{} {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    return map[string]interface{}{
        "state":               cb.state.String(),
        "failures":            cb.failures,
        "successes":           cb.successes,
        "consecutiveSuccesses": cb.consecutiveSuccesses,
        "lastFailureTime":     cb.lastFailureTime,
    }
}

// CircuitBreakerManager 断路器管理器
type CircuitBreakerManager struct {
    breakers map[string]*CircuitBreaker
    mu       sync.RWMutex
    config   *Config
}

// NewCircuitBreakerManager 创建断路器管理器
func NewCircuitBreakerManager(config *Config) *CircuitBreakerManager {
    return &CircuitBreakerManager{
        breakers: make(map[string]*CircuitBreaker),
        config:   config,
    }
}

// GetCircuitBreaker 获取或创建断路器
func (m *CircuitBreakerManager) GetCircuitBreaker(name string) *CircuitBreaker {
    m.mu.RLock()
    cb, ok := m.breakers[name]
    m.mu.RUnlock()

    if ok {
        return cb
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    // 双重检查
    if cb, ok := m.breakers[name]; ok {
        return cb
    }

    cb = NewCircuitBreaker(name, m.config)
    m.breakers[name] = cb
    return cb
}

// GetAllMetrics 获取所有断路器指标
func (m *CircuitBreakerManager) GetAllMetrics() map[string]map[string]interface{} {
    m.mu.RLock()
    defer m.mu.RUnlock()

    metrics := make(map[string]map[string]interface{})
    for name, cb := range m.breakers {
        metrics[name] = cb.GetMetrics()
    }
    return metrics
}

// HTTPClient 带断路器的 HTTP 客户端
type HTTPClient struct {
    baseClient *http.Client
    cbManager  *CircuitBreakerManager
}

func NewHTTPClient(cbManager *CircuitBreakerManager) *HTTPClient {
    return &HTTPClient{
        baseClient: &http.Client{Timeout: 10 * time.Second},
        cbManager:  cbManager,
    }
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
    serviceName := req.URL.Host
    cb := c.cbManager.GetCircuitBreaker(serviceName)

    var resp *http.Response
    var err error

    execErr := cb.Execute(req.Context(), func() error {
        resp, err = c.baseClient.Do(req)
        if err != nil {
            return err
        }

        // 5xx 错误视为服务故障
        if resp.StatusCode >= 500 {
            return fmt.Errorf("server error: %d", resp.StatusCode)
        }

        return nil
    })

    if execErr != nil {
        return nil, execErr
    }

    return resp, err
}

// GRPCClient 带断路器的 gRPC 客户端
type GRPCClient struct {
    cbManager *CircuitBreakerManager
}

func NewGRPCClient(cbManager *CircuitBreakerManager) *GRPCClient {
    return &GRPCClient{cbManager: cbManager}
}

func (c *GRPCClient) Invoke(ctx context.Context, serviceName string, fn func() error) error {
    cb := c.cbManager.GetCircuitBreaker(serviceName)
    return cb.Execute(ctx, fn)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 防止级联故障 | 增加系统复杂度 |
| 快速失败 | 需要合理配置阈值 |
| 自动恢复 | 可能误熔断 |

#### 适用场景

- 微服务调用
- 外部 API 调用
- 数据库连接

---

### 2.7 Bulkhead 模式

#### 意图

Bulkhead（舱壁）模式通过隔离资源池来限制故障影响范围，防止一个组件的故障耗尽所有资源。

#### 结构

```
┌─────────────────────────────────────────────────────────┐
│                    Application                          │
│                                                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │  Service A  │  │  Service B  │  │  Service C  │     │
│  │   Pool      │  │   Pool      │  │   Pool      │     │
│  │ ┌───┬───┐   │  │ ┌───┬───┐   │  │ ┌───┬───┐   │     │
│  │ │ C │ C │   │  │ │ C │ C │   │  │ │ C │ C │   │     │
│  │ └───┴───┘   │  │ └───┴───┘   │  │ └───┴───┘   │     │
│  │  (max: 10)  │  │  (max: 20)  │  │  (max: 5)   │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
│                                                         │
│  隔离的资源池，一个服务耗尽不影响其他服务                  │
└─────────────────────────────────────────────────────────┘
```

#### 实现

```go
// Bulkhead 模式实现
package bulkhead

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"
)

// ErrBulkheadFull 资源池已满错误
var ErrBulkheadFull = errors.New("bulkhead is full")

// ErrBulkheadTimeout 等待超时错误
var ErrBulkheadTimeout = errors.New("bulkhead wait timeout")

// Bulkhead 舱壁
type Bulkhead struct {
    name         string
    maxConcurrent int
    semaphore    chan struct{}
    queueSize    int
    queue        chan func()
    wg           sync.WaitGroup
    mu           sync.RWMutex
    metrics      *BulkheadMetrics
}

// BulkheadMetrics 舱壁指标
type BulkheadMetrics struct {
    ActiveExecutions    int64
    QueuedExecutions    int64
    CompletedExecutions int64
    RejectedExecutions  int64
    AverageWaitTime     time.Duration
}

// Config 舱壁配置
type Config struct {
    MaxConcurrent int           // 最大并发数
    QueueSize     int           // 队列大小
    MaxWaitTime   time.Duration // 最大等待时间
}

// NewBulkhead 创建舱壁
func NewBulkhead(name string, config *Config) *Bulkhead {
    if config.MaxConcurrent <= 0 {
        config.MaxConcurrent = 10
    }
    if config.QueueSize < 0 {
        config.QueueSize = 0
    }

    b := &Bulkhead{
        name:          name,
        maxConcurrent: config.MaxConcurrent,
        semaphore:     make(chan struct{}, config.MaxConcurrent),
        queueSize:     config.QueueSize,
        metrics:       &BulkheadMetrics{},
    }

    if config.QueueSize > 0 {
        b.queue = make(chan func(), config.QueueSize)
        go b.processQueue()
    }

    return b
}

// Execute 执行受保护的操作
func (b *Bulkhead) Execute(ctx context.Context, fn func() error) error {
    select {
    case b.semaphore <- struct{}{}:
        // 获取到信号量，直接执行
        defer func() { <-b.semaphore }()
        return fn()

    default:
        // 信号量已满
        if b.queue == nil {
            return ErrBulkheadFull
        }

        // 尝试加入队列
        resultChan := make(chan error, 1)

        select {
        case b.queue <- func() {
            b.semaphore <- struct{}{}
            defer func() { <-b.semaphore }()
            resultChan <- fn()
        }:
            // 等待执行完成
            select {
            case err := <-resultChan:
                return err
            case <-ctx.Done():
                return ctx.Err()
            }

        default:
            // 队列已满
            return ErrBulkheadFull
        }
    }
}

// ExecuteWithTimeout 带超时的执行
func (b *Bulkhead) ExecuteWithTimeout(timeout time.Duration, fn func() error) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    return b.Execute(ctx, fn)
}

func (b *Bulkhead) processQueue() {
    for task := range b.queue {
        b.wg.Add(1)
        go func(t func()) {
            defer b.wg.Done()
            t()
        }(task)
    }
}

// GetMetrics 获取指标
func (b *Bulkhead) GetMetrics() BulkheadMetrics {
    b.mu.RLock()
    defer b.mu.RUnlock()
    return *b.metrics
}

// Close 关闭舱壁
func (b *Bulkhead) Close() error {
    if b.queue != nil {
        close(b.queue)
    }
    b.wg.Wait()
    return nil
}

// BulkheadManager 舱壁管理器
type BulkheadManager struct {
    bulkheads map[string]*Bulkhead
    mu        sync.RWMutex
    configs   map[string]*Config
}

// NewBulkheadManager 创建舱壁管理器
func NewBulkheadManager() *BulkheadManager {
    return &BulkheadManager{
        bulkheads: make(map[string]*Bulkhead),
        configs:   make(map[string]*Config),
    }
}

// Register 注册舱壁配置
func (m *BulkheadManager) Register(name string, config *Config) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.configs[name] = config
}

// GetBulkhead 获取舱壁
func (m *BulkheadManager) GetBulkhead(name string) *Bulkhead {
    m.mu.RLock()
    bulkhead, ok := m.bulkheads[name]
    m.mu.RUnlock()

    if ok {
        return bulkhead
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    // 双重检查
    if bulkhead, ok := m.bulkheads[name]; ok {
        return bulkhead
    }

    config := m.configs[name]
    if config == nil {
        config = &Config{MaxConcurrent: 10}
    }

    bulkhead = NewBulkhead(name, config)
    m.bulkheads[name] = bulkhead
    return bulkhead
}

// 数据库连接池舱壁示例

type DatabaseBulkhead struct {
    bulkhead *Bulkhead
    db       *sql.DB
}

func NewDatabaseBulkhead(db *sql.DB, maxConnections int) *DatabaseBulkhead {
    config := &Config{
        MaxConcurrent: maxConnections,
        QueueSize:     100,
        MaxWaitTime:   5 * time.Second,
    }

    return &DatabaseBulkhead{
        bulkhead: NewBulkhead("database", config),
        db:       db,
    }
}

func (d *DatabaseBulkhead) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    var rows *sql.Rows
    var err error

    execErr := d.bulkhead.Execute(ctx, func() error {
        rows, err = d.db.QueryContext(ctx, query, args...)
        return err
    })

    if execErr != nil {
        return nil, execErr
    }

    return rows, err
}

// HTTP 服务舱壁示例

type HTTPServiceBulkhead struct {
    bulkhead *Bulkhead
    client   *http.Client
    baseURL  string
}

func NewHTTPServiceBulkhead(name, baseURL string, maxConcurrent int) *HTTPServiceBulkhead {
    config := &Config{
        MaxConcurrent: maxConcurrent,
        QueueSize:     50,
        MaxWaitTime:   10 * time.Second,
    }

    return &HTTPServiceBulkhead{
        bulkhead: NewBulkhead(name, config),
        client:   &http.Client{Timeout: 30 * time.Second},
        baseURL:  baseURL,
    }
}

func (h *HTTPServiceBulkhead) Call(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
    var resp *http.Response
    var err error

    execErr := h.bulkhead.Execute(ctx, func() error {
        url := h.baseURL + path
        var bodyReader io.Reader
        if body != nil {
            bodyReader = bytes.NewReader(body)
        }

        req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
        if err != nil {
            return err
        }

        resp, err = h.client.Do(req)
        return err
    })

    if execErr != nil {
        return nil, execErr
    }

    return resp, err
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 资源隔离 | 配置复杂 |
| 故障隔离 | 资源利用率降低 |
| 可预测性能 | 需要监控 |

#### 适用场景

- 数据库连接池
- 外部服务调用
- 资源密集型操作

---

### 2.8 Retry 模式

#### 意图

Retry 模式通过自动重试失败的操作来处理临时性故障，提高系统的可靠性和可用性。

#### 结构

```
请求 ──► 失败 ──► 等待 ──► 重试 ──► 成功
         │         │
         │         └── 指数退避: 1s, 2s, 4s, 8s...
         │
         └── 可重试错误: 网络超时、5xx 错误
```

#### 实现

```go
// Retry 模式实现
package retry

import (
    "context"
    "errors"
    "fmt"
    "math"
    "math/rand"
    "time"
)

// RetryableError 可重试错误
type RetryableError struct {
    Err error
}

func (e *RetryableError) Error() string {
    return fmt.Sprintf("retryable: %v", e.Err)
}

func IsRetryable(err error) bool {
    var re *RetryableError
    return errors.As(err, &re)
}

// Policy 重试策略
type Policy struct {
    MaxRetries  int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Multiplier  float64
    Jitter      float64 // 抖动因子 0-1
    Retryable   func(error) bool
}

// DefaultPolicy 默认重试策略
func DefaultPolicy() *Policy {
    return &Policy{
        MaxRetries: 3,
        BaseDelay:  100 * time.Millisecond,
        MaxDelay:   30 * time.Second,
        Multiplier: 2.0,
        Jitter:     0.1,
        Retryable:  defaultRetryable,
    }
}

func defaultRetryable(err error) bool {
    return IsRetryable(err)
}

// ExponentialBackoff 指数退避
type ExponentialBackoff struct {
    policy *Policy
}

func NewExponentialBackoff(policy *Policy) *ExponentialBackoff {
    if policy == nil {
        policy = DefaultPolicy()
    }
    return &ExponentialBackoff{policy: policy}
}

// NextDelay 计算下一次延迟
func (e *ExponentialBackoff) NextDelay(attempt int) time.Duration {
    if attempt == 0 {
        return 0
    }

    // 指数计算
    delay := float64(e.policy.BaseDelay) * math.Pow(e.policy.Multiplier, float64(attempt-1))

    // 应用最大值限制
    if delay > float64(e.policy.MaxDelay) {
        delay = float64(e.policy.MaxDelay)
    }

    // 添加抖动
    if e.policy.Jitter > 0 {
        jitter := delay * e.policy.Jitter * (2*rand.Float64() - 1)
        delay += jitter
    }

    return time.Duration(delay)
}

// Retry 执行重试
func (e *ExponentialBackoff) Retry(ctx context.Context, fn func() error) error {
    var lastErr error

    for attempt := 0; attempt <= e.policy.MaxRetries; attempt++ {
        if attempt > 0 {
            delay := e.NextDelay(attempt)

            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return ctx.Err()
            }
        }

        err := fn()
        if err == nil {
            return nil
        }

        lastErr = err

        if !e.policy.Retryable(err) {
            return err
        }
    }

    return fmt.Errorf("exhausted %d retries: %w", e.policy.MaxRetries, lastErr)
}

// RetryWithResult 带返回值的重试
func (e *ExponentialBackoff) RetryWithResult(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
    var lastErr error

    for attempt := 0; attempt <= e.policy.MaxRetries; attempt++ {
        if attempt > 0 {
            delay := e.NextDelay(attempt)

            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return nil, ctx.Err()
            }
        }

        result, err := fn()
        if err == nil {
            return result, nil
        }

        lastErr = err

        if !e.policy.Retryable(err) {
            return nil, err
        }
    }

    return nil, fmt.Errorf("exhausted %d retries: %w", e.policy.MaxRetries, lastErr)
}

// LinearBackoff 线性退避
type LinearBackoff struct {
    policy *Policy
}

func NewLinearBackoff(policy *Policy) *LinearBackoff {
    if policy == nil {
        policy = DefaultPolicy()
    }
    return &LinearBackoff{policy: policy}
}

func (l *LinearBackoff) NextDelay(attempt int) time.Duration {
    if attempt == 0 {
        return 0
    }

    delay := time.Duration(attempt) * l.policy.BaseDelay
    if delay > l.policy.MaxDelay {
        delay = l.policy.MaxDelay
    }

    return delay
}

// FixedBackoff 固定间隔
type FixedBackoff struct {
    policy *Policy
}

func NewFixedBackoff(policy *Policy) *FixedBackoff {
    if policy == nil {
        policy = DefaultPolicy()
    }
    return &FixedBackoff{policy: policy}
}

func (f *FixedBackoff) NextDelay(attempt int) time.Duration {
    if attempt == 0 {
        return 0
    }
    return f.policy.BaseDelay
}

// DecorrelatedJitterBackoff 去相关抖动退避
type DecorrelatedJitterBackoff struct {
    policy *Policy
    lastDelay time.Duration
}

func NewDecorrelatedJitterBackoff(policy *Policy) *DecorrelatedJitterBackoff {
    if policy == nil {
        policy = DefaultPolicy()
    }
    return &DecorrelatedJitterBackoff{policy: policy}
}

func (d *DecorrelatedJitterBackoff) NextDelay(attempt int) time.Duration {
    if attempt == 0 {
        return 0
    }

    // 去相关抖动: sleep = min(cap, random_between(base, sleep * 3))
    maxDelay := d.lastDelay * 3
    if maxDelay == 0 {
        maxDelay = d.policy.BaseDelay
    }
    if maxDelay > d.policy.MaxDelay {
        maxDelay = d.policy.MaxDelay
    }

    delay := time.Duration(rand.Int63n(int64(maxDelay-d.policy.BaseDelay))) + d.policy.BaseDelay
    d.lastDelay = delay

    return delay
}

// RetryManager 重试管理器
type RetryManager struct {
    policies map[string]*Policy
}

func NewRetryManager() *RetryManager {
    return &RetryManager{
        policies: make(map[string]*Policy),
    }
}

func (r *RetryManager) Register(name string, policy *Policy) {
    r.policies[name] = policy
}

func (r *RetryManager) Execute(ctx context.Context, name string, fn func() error) error {
    policy, ok := r.policies[name]
    if !ok {
        policy = DefaultPolicy()
    }

    backoff := NewExponentialBackoff(policy)
    return backoff.Retry(ctx, fn)
}

// HTTPRetryableTransport 带重试的 HTTP Transport
type HTTPRetryableTransport struct {
    Base      http.RoundTripper
    RetryPolicy *Policy
}

func (t *HTTPRetryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    var resp *http.Response
    var err error

    backoff := NewExponentialBackoff(t.RetryPolicy)

    retryErr := backoff.Retry(req.Context(), func() error {
        resp, err = t.Base.RoundTrip(req.Clone(req.Context()))
        if err != nil {
            return &RetryableError{Err: err}
        }

        // 5xx 错误可重试
        if resp.StatusCode >= 500 {
            resp.Body.Close()
            return &RetryableError{Err: fmt.Errorf("server error: %d", resp.StatusCode)}
        }

        return nil
    })

    if retryErr != nil {
        return nil, retryErr
    }

    return resp, err
}

// 使用示例
func ExampleRetry() {
    policy := &Policy{
        MaxRetries: 5,
        BaseDelay:  100 * time.Millisecond,
        MaxDelay:   5 * time.Second,
        Multiplier: 2.0,
        Jitter:     0.1,
    }

    backoff := NewExponentialBackoff(policy)

    ctx := context.Background()
    err := backoff.Retry(ctx, func() error {
        // 执行可能失败的操作
        return callExternalService()
    })

    if err != nil {
        fmt.Printf("Failed after retries: %v\n", err)
    }
}

func callExternalService() error {
    // 模拟外部服务调用
    return nil
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 提高可靠性 | 增加延迟 |
| 自动恢复 | 可能加剧故障 |
| 简单易用 | 需要合理配置 |

#### 适用场景

- 网络请求
- 数据库操作
- 外部 API 调用

---



## 3. 并发与并行模式

并发与并行模式是构建高性能分布式系统的核心。Kubernetes 和 Go 语言天然支持这些模式。

---

### 3.1 Worker Pool 模式

#### 意图

Worker Pool 模式维护一组固定数量的工作协程，从任务队列中获取任务并执行，避免频繁创建和销毁协程的开销。

#### 结构

```
┌─────────────────────────────────────────────────────────┐
│                    Worker Pool                          │
│                                                         │
│   ┌─────────────┐                                       │
│   │ Task Queue  │◄────── Submit Task                   │
│   │  (Buffered) │                                       │
│   └──────┬──────┘                                       │
│          │                                              │
│    ┌─────┼─────┬─────────┬─────────┐                   │
│    │     │     │         │         │                   │
│    ▼     ▼     ▼         ▼         ▼                   │
│ ┌────┐ ┌────┐ ┌────┐  ┌────┐  ┌────┐                  │
│ │ W1 │ │ W2 │ │ W3 │  │ W4 │  │ W5 │  Workers          │
│ └────┘ └────┘ └────┘  └────┘  └────┘                  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### 实现

```go
// Worker Pool 模式实现
package workerpool

import (
    "context"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// Task 任务接口
type Task interface {
    Execute() error
    GetID() string
}

// TaskFunc 任务函数类型
type TaskFunc func() error

// SimpleTask 简单任务实现
type SimpleTask struct {
    id   string
    fn   TaskFunc
}

func NewSimpleTask(id string, fn TaskFunc) *SimpleTask {
    return &SimpleTask{id: id, fn: fn}
}

func (t *SimpleTask) Execute() error {
    return t.fn()
}

func (t *SimpleTask) GetID() string {
    return t.id
}

// Worker 工作协程
type Worker struct {
    id       int
    pool     *Pool
    taskChan chan Task
    quit     chan bool
}

func (w *Worker) start() {
    go func() {
        for {
            // 注册到 pool 的 readyWorkers
            w.pool.readyWorkers <- w.taskChan

            select {
            case task := <-w.taskChan:
                // 执行任务
                w.pool.wg.Add(1)
                w.execute(task)
                w.pool.wg.Done()

            case <-w.quit:
                return
            }
        }
    }()
}

func (w *Worker) execute(task Task) {
    defer func() {
        if r := recover(); r != nil {
            w.pool.recordError(fmt.Errorf("worker %d panic: %v", w.id, r))
        }
    }()

    start := time.Now()
    err := task.Execute()
    duration := time.Since(start)

    w.pool.recordMetrics(task.GetID(), duration, err)
}

func (w *Worker) stop() {
    go func() {
        w.quit <- true
    }()
}

// Pool 工作池
type Pool struct {
    workers      []*Worker
    workerCount  int
    taskQueue    chan Task
    readyWorkers chan chan Task
    wg           sync.WaitGroup
    quit         chan bool

    // 指标
    tasksSubmitted int64
    tasksCompleted int64
    tasksFailed    int64

    // 错误处理
    errorHandler func(error)

    // 配置
    maxQueueSize int
    timeout      time.Duration
}

// PoolConfig 工作池配置
type PoolConfig struct {
    WorkerCount  int
    MaxQueueSize int
    Timeout      time.Duration
    ErrorHandler func(error)
}

// NewPool 创建工作池
func NewPool(config *PoolConfig) *Pool {
    if config.WorkerCount <= 0 {
        config.WorkerCount = 10
    }
    if config.MaxQueueSize <= 0 {
        config.MaxQueueSize = 100
    }

    pool := &Pool{
        workerCount:  config.WorkerCount,
        taskQueue:    make(chan Task, config.MaxQueueSize),
        readyWorkers: make(chan chan Task, config.WorkerCount),
        quit:         make(chan bool),
        maxQueueSize: config.MaxQueueSize,
        timeout:      config.Timeout,
        errorHandler: config.ErrorHandler,
    }

    // 创建工作协程
    pool.workers = make([]*Worker, config.WorkerCount)
    for i := 0; i < config.WorkerCount; i++ {
        worker := &Worker{
            id:       i,
            pool:     pool,
            taskChan: make(chan Task),
            quit:     make(chan bool),
        }
        pool.workers[i] = worker
        worker.start()
    }

    // 启动分发器
    go pool.dispatch()

    return pool
}

func (p *Pool) dispatch() {
    for {
        select {
        case task := <-p.taskQueue:
            // 获取一个就绪的工作协程
            workerChan := <-p.readyWorkers
            workerChan <- task

        case <-p.quit:
            return
        }
    }
}

// Submit 提交任务
func (p *Pool) Submit(task Task) error {
    select {
    case p.taskQueue <- task:
        atomic.AddInt64(&p.tasksSubmitted, 1)
        return nil
    default:
        return fmt.Errorf("task queue is full")
    }
}

// SubmitWithTimeout 带超时的任务提交
func (p *Pool) SubmitWithTimeout(task Task, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    select {
    case p.taskQueue <- task:
        atomic.AddInt64(&p.tasksSubmitted, 1)
        return nil
    case <-ctx.Done():
        return fmt.Errorf("submit timeout")
    }
}

// Stop 停止工作池
func (p *Pool) Stop() {
    close(p.quit)

    // 停止所有工作协程
    for _, worker := range p.workers {
        worker.stop()
    }

    // 等待所有任务完成
    p.wg.Wait()
}

// StopWithTimeout 带超时的停止
func (p *Pool) StopWithTimeout(timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    done := make(chan struct{})
    go func() {
        p.Stop()
        close(done)
    }()

    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (p *Pool) recordMetrics(taskID string, duration time.Duration, err error) {
    if err != nil {
        atomic.AddInt64(&p.tasksFailed, 1)
    } else {
        atomic.AddInt64(&p.tasksCompleted, 1)
    }
}

func (p *Pool) recordError(err error) {
    if p.errorHandler != nil {
        p.errorHandler(err)
    }
}

// GetMetrics 获取指标
func (p *Pool) GetMetrics() map[string]int64 {
    return map[string]int64{
        "workers":         int64(p.workerCount),
        "queue_size":      int64(len(p.taskQueue)),
        "tasks_submitted": atomic.LoadInt64(&p.tasksSubmitted),
        "tasks_completed": atomic.LoadInt64(&p.tasksCompleted),
        "tasks_failed":    atomic.LoadInt64(&p.tasksFailed),
    }
}

// PriorityPool 优先级工作池
type PriorityPool struct {
    workers     []*Worker
    workerCount int

    // 多个优先级的队列
    highPriorityQueue   chan Task
    normalPriorityQueue chan Task
    lowPriorityQueue    chan Task

    wg   sync.WaitGroup
    quit chan bool
}

// NewPriorityPool 创建优先级工作池
func NewPriorityPool(workerCount int) *PriorityPool {
    if workerCount <= 0 {
        workerCount = 10
    }

    pool := &PriorityPool{
        workerCount:         workerCount,
        highPriorityQueue:   make(chan Task, 100),
        normalPriorityQueue: make(chan Task, 200),
        lowPriorityQueue:    make(chan Task, 300),
        quit:                make(chan bool),
    }

    // 启动优先级分发器
    go pool.priorityDispatch()

    return pool
}

func (p *PriorityPool) priorityDispatch() {
    for {
        select {
        // 高优先级优先处理
        case task := <-p.highPriorityQueue:
            p.executeTask(task)

        default:
            select {
            case task := <-p.highPriorityQueue:
                p.executeTask(task)
            case task := <-p.normalPriorityQueue:
                p.executeTask(task)
            case task := <-p.lowPriorityQueue:
                p.executeTask(task)
            case <-p.quit:
                return
            }
        }
    }
}

func (p *PriorityPool) executeTask(task Task) {
    p.wg.Add(1)
    go func() {
        defer p.wg.Done()
        task.Execute()
    }()
}

func (p *PriorityPool) SubmitHighPriority(task Task) {
    p.highPriorityQueue <- task
}

func (p *PriorityPool) SubmitNormalPriority(task Task) {
    p.normalPriorityQueue <- task
}

func (p *PriorityPool) SubmitLowPriority(task Task) {
    p.lowPriorityQueue <- task
}

// 使用示例
func ExampleWorkerPool() {
    config := &PoolConfig{
        WorkerCount:  5,
        MaxQueueSize: 100,
        Timeout:      30 * time.Second,
        ErrorHandler: func(err error) {
            fmt.Printf("Task error: %v\n", err)
        },
    }

    pool := NewPool(config)
    defer pool.Stop()

    // 提交任务
    for i := 0; i < 50; i++ {
        taskID := fmt.Sprintf("task-%d", i)
        task := NewSimpleTask(taskID, func() error {
            fmt.Printf("Executing %s\n", taskID)
            time.Sleep(100 * time.Millisecond)
            return nil
        })

        if err := pool.Submit(task); err != nil {
            fmt.Printf("Failed to submit %s: %v\n", taskID, err)
        }
    }

    // 等待所有任务完成
    time.Sleep(5 * time.Second)

    // 打印指标
    metrics := pool.GetMetrics()
    fmt.Printf("Metrics: %+v\n", metrics)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 复用协程 | 需要合理配置大小 |
| 控制并发数 | 任务队列可能堆积 |
| 提高性能 | 增加复杂度 |

#### 适用场景

- 批处理任务
- HTTP 请求处理
- 消息消费

---

### 3.2 Pub/Sub 模式

#### 意图

Pub/Sub（发布-订阅）模式实现松耦合的消息通信，发布者不需要知道订阅者的存在。

#### 结构

```
┌─────────┐     ┌─────────────┐     ┌─────────┐
│Publisher│────►│   Topic     │◄────│Subscriber│
└─────────┘     │  (Channel)  │     └─────────┘
                └──────┬──────┘
                       │
           ┌───────────┼───────────┐
           │           │           │
           ▼           ▼           ▼
      ┌─────────┐ ┌─────────┐ ┌─────────┐
      │  Sub 1  │ │  Sub 2  │ │  Sub 3  │
      └─────────┘ └─────────┘ └─────────┘
```

#### 实现

```go
// Pub/Sub 模式实现
package pubsub

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Message 消息
type Message struct {
    ID        string
    Topic     string
    Payload   interface{}
    Timestamp time.Time
    Metadata  map[string]string
}

// Handler 消息处理器
type Handler func(msg *Message) error

// Subscriber 订阅者
type Subscriber struct {
    id      string
    handler Handler
    topics  []string
}

// Broker 消息代理
type Broker struct {
    topics       map[string][]*Subscriber
    subscribers  map[string]*Subscriber
    messageChan  chan *Message
    mu           sync.RWMutex
    wg           sync.WaitGroup
    ctx          context.Context
    cancel       context.CancelFunc
}

// NewBroker 创建消息代理
func NewBroker(bufferSize int) *Broker {
    ctx, cancel := context.WithCancel(context.Background())

    broker := &Broker{
        topics:      make(map[string][]*Subscriber),
        subscribers: make(map[string]*Subscriber),
        messageChan: make(chan *Message, bufferSize),
        ctx:         ctx,
        cancel:      cancel,
    }

    // 启动消息分发
    go broker.dispatch()

    return broker
}

// Subscribe 订阅主题
func (b *Broker) Subscribe(subscriberID string, topics []string, handler Handler) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if _, exists := b.subscribers[subscriberID]; exists {
        return fmt.Errorf("subscriber %s already exists", subscriberID)
    }

    subscriber := &Subscriber{
        id:      subscriberID,
        handler: handler,
        topics:  topics,
    }

    b.subscribers[subscriberID] = subscriber

    // 注册到主题
    for _, topic := range topics {
        b.topics[topic] = append(b.topics[topic], subscriber)
    }

    return nil
}

// Unsubscribe 取消订阅
func (b *Broker) Unsubscribe(subscriberID string) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    subscriber, exists := b.subscribers[subscriberID]
    if !exists {
        return fmt.Errorf("subscriber %s not found", subscriberID)
    }

    // 从所有主题中移除
    for _, topic := range subscriber.topics {
        subscribers := b.topics[topic]
        for i, sub := range subscribers {
            if sub.id == subscriberID {
                b.topics[topic] = append(subscribers[:i], subscribers[i+1:]...)
                break
            }
        }
    }

    delete(b.subscribers, subscriberID)

    return nil
}

// Publish 发布消息
func (b *Broker) Publish(msg *Message) error {
    select {
    case b.messageChan <- msg:
        return nil
    default:
        return fmt.Errorf("message channel is full")
    }
}

// PublishWithTimeout 带超时的发布
func (b *Broker) PublishWithTimeout(msg *Message, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(b.ctx, timeout)
    defer cancel()

    select {
    case b.messageChan <- msg:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (b *Broker) dispatch() {
    for {
        select {
        case msg := <-b.messageChan:
            b.deliver(msg)
        case <-b.ctx.Done():
            return
        }
    }
}

func (b *Broker) deliver(msg *Message) {
    b.mu.RLock()
    subscribers := b.topics[msg.Topic]
    b.mu.RUnlock()

    for _, subscriber := range subscribers {
        b.wg.Add(1)
        go func(sub *Subscriber) {
            defer b.wg.Done()

            if err := sub.handler(msg); err != nil {
                fmt.Printf("Handler error for subscriber %s: %v\n", sub.id, err)
            }
        }(subscriber)
    }
}

// Close 关闭代理
func (b *Broker) Close() {
    b.cancel()
    b.wg.Wait()
    close(b.messageChan)
}

// GetStats 获取统计信息
func (b *Broker) GetStats() map[string]interface{} {
    b.mu.RLock()
    defer b.mu.RUnlock()

    topicStats := make(map[string]int)
    for topic, subs := range b.topics {
        topicStats[topic] = len(subs)
    }

    return map[string]interface{}{
        "subscribers": len(b.subscribers),
        "topics":      len(b.topics),
        "topic_stats": topicStats,
        "queue_size":  len(b.messageChan),
    }
}

// PersistentBroker 持久化消息代理
type PersistentBroker struct {
    *Broker
    storage MessageStorage
}

// MessageStorage 消息存储接口
type MessageStorage interface {
    Save(msg *Message) error
    Get(topic string, offset int64) ([]*Message, error)
    Ack(messageID string) error
}

func NewPersistentBroker(bufferSize int, storage MessageStorage) *PersistentBroker {
    return &PersistentBroker{
        Broker:  NewBroker(bufferSize),
        storage: storage,
    }
}

func (p *PersistentBroker) Publish(msg *Message) error {
    // 先持久化
    if err := p.storage.Save(msg); err != nil {
        return err
    }

    return p.Broker.Publish(msg)
}

// OrderedBroker 有序消息代理
type OrderedBroker struct {
    *Broker
    partitions int
    partitionsMap map[string]chan *Message
}

func NewOrderedBroker(bufferSize, partitions int) *OrderedBroker {
    broker := &OrderedBroker{
        Broker:        NewBroker(bufferSize),
        partitions:    partitions,
        partitionsMap: make(map[string]chan *Message),
    }

    // 为每个分区创建通道
    for i := 0; i < partitions; i++ {
        partitionChan := make(chan *Message, bufferSize/partitions)
        broker.partitionsMap[fmt.Sprintf("partition-%d", i)] = partitionChan
        go broker.processPartition(partitionChan)
    }

    return broker
}

func (o *OrderedBroker) processPartition(ch chan *Message) {
    for msg := range ch {
        o.Broker.deliver(msg)
    }
}

func (o *OrderedBroker) PublishOrdered(msg *Message, partitionKey string) error {
    // 根据 partitionKey 选择分区
    partition := o.getPartition(partitionKey)

    select {
    case o.partitionsMap[partition] <- msg:
        return nil
    default:
        return fmt.Errorf("partition %s is full", partition)
    }
}

func (o *OrderedBroker) getPartition(key string) string {
    // 简单的哈希分区
    hash := 0
    for _, c := range key {
        hash += int(c)
    }
    partition := hash % o.partitions
    return fmt.Sprintf("partition-%d", partition)
}

// 使用示例
func ExamplePubSub() {
    broker := NewBroker(1000)
    defer broker.Close()

    // 订阅者1：处理订单事件
    broker.Subscribe("order-processor", []string{"orders"}, func(msg *Message) error {
        fmt.Printf("Order processor received: %+v\n", msg.Payload)
        return nil
    })

    // 订阅者2：处理库存事件
    broker.Subscribe("inventory-processor", []string{"inventory"}, func(msg *Message) error {
        fmt.Printf("Inventory processor received: %+v\n", msg.Payload)
        return nil
    })

    // 订阅者3：处理所有事件（日志）
    broker.Subscribe("logger", []string{"orders", "inventory"}, func(msg *Message) error {
        fmt.Printf("Logger: [%s] %s\n", msg.Topic, msg.ID)
        return nil
    })

    // 发布消息
    for i := 0; i < 10; i++ {
        broker.Publish(&Message{
            ID:        fmt.Sprintf("order-%d", i),
            Topic:     "orders",
            Payload:   map[string]interface{}{"order_id": i, "amount": 100.0},
            Timestamp: time.Now(),
        })
    }

    time.Sleep(2 * time.Second)

    // 打印统计
    stats := broker.GetStats()
    fmt.Printf("Stats: %+v\n", stats)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 松耦合 | 消息可能丢失 |
| 可扩展 | 消息顺序难保证 |
| 异步通信 | 调试困难 |

#### 适用场景

- 事件驱动架构
- 日志聚合
- 通知系统

---

### 3.3 Leader Election 模式

#### 意图

Leader Election 模式确保在分布式系统中只有一个实例担任领导者角色，负责协调任务或处理独占资源。

#### 结构

```
┌─────────────────────────────────────────────────────────┐
│                  Kubernetes Cluster                     │
│                                                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                 │
│  │ Pod 1   │  │ Pod 2   │  │ Pod 3   │                 │
│  │         │  │         │  │         │                 │
│  │ LEADER  │  │Follower │  │Follower │                 │
│  │  [Active]│  │ [Standby]│  │ [Standby]│                │
│  └─────────┘  └─────────┘  └─────────┘                 │
│       │            │            │                      │
│       └────────────┴────────────┘                      │
│                    │                                    │
│              ┌─────┴─────┐                             │
│              │   Lease   │                             │
│              │  (etcd)   │                             │
│              └───────────┘                             │
└─────────────────────────────────────────────────────────┘
```

#### 实现

```go
// Leader Election 模式实现
package leaderelection

import (
    "context"
    "fmt"
    "sync"
    "time"

    coordinationv1 "k8s.io/api/coordination/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/leaderelection"
    "k8s.io/client-go/tools/leaderelection/resourcelock"
)

// LeaderCallbacks 领导者回调
type LeaderCallbacks struct {
    OnStartedLeading func(context.Context)
    OnStoppedLeading func()
    OnNewLeader      func(identity string)
}

// LeaderElector 领导者选举器
type LeaderElector struct {
    client        kubernetes.Interface
    lockName      string
    namespace     string
    identity      string
    leaseDuration time.Duration
    renewDeadline time.Duration
    retryPeriod   time.Duration
    callbacks     LeaderCallbacks

    isLeader bool
    mu       sync.RWMutex
}

// Config 选举配置
type Config struct {
    Client        kubernetes.Interface
    LockName      string
    Namespace     string
    Identity      string
    LeaseDuration time.Duration
    RenewDeadline time.Duration
    RetryPeriod   time.Duration
    Callbacks     LeaderCallbacks
}

// NewLeaderElector 创建领导者选举器
func NewLeaderElector(config *Config) (*LeaderElector, error) {
    if config.LeaseDuration == 0 {
        config.LeaseDuration = 15 * time.Second
    }
    if config.RenewDeadline == 0 {
        config.RenewDeadline = 10 * time.Second
    }
    if config.RetryPeriod == 0 {
        config.RetryPeriod = 2 * time.Second
    }

    return &LeaderElector{
        client:        config.Client,
        lockName:      config.LockName,
        namespace:     config.Namespace,
        identity:      config.Identity,
        leaseDuration: config.LeaseDuration,
        renewDeadline: config.RenewDeadline,
        retryPeriod:   config.RetryPeriod,
        callbacks:     config.Callbacks,
    }, nil
}

// Run 启动选举
func (le *LeaderElector) Run(ctx context.Context) error {
    // 创建资源锁
    lock := &resourcelock.LeaseLock{
        LeaseMeta: metav1.ObjectMeta{
            Name:      le.lockName,
            Namespace: le.namespace,
        },
        Client: le.client.CoordinationV1(),
        LockConfig: resourcelock.ResourceLockConfig{
            Identity: le.identity,
        },
    }

    // 配置选举
    lec := leaderelection.LeaderElectionConfig{
        Lock:            lock,
        LeaseDuration:   le.leaseDuration,
        RenewDeadline:   le.renewDeadline,
        RetryPeriod:     le.retryPeriod,
        ReleaseOnCancel: true,
        Callbacks: leaderelection.LeaderCallbacks{
            OnStartedLeading: func(ctx context.Context) {
                le.mu.Lock()
                le.isLeader = true
                le.mu.Unlock()

                fmt.Printf("[%s] Became leader\n", le.identity)

                if le.callbacks.OnStartedLeading != nil {
                    le.callbacks.OnStartedLeading(ctx)
                }
            },
            OnStoppedLeading: func() {
                le.mu.Lock()
                le.isLeader = false
                le.mu.Unlock()

                fmt.Printf("[%s] Stopped leading\n", le.identity)

                if le.callbacks.OnStoppedLeading != nil {
                    le.callbacks.OnStoppedLeading()
                }
            },
            OnNewLeader: func(identity string) {
                fmt.Printf("[%s] New leader elected: %s\n", le.identity, identity)

                if le.callbacks.OnNewLeader != nil {
                    le.callbacks.OnNewLeader(identity)
                }
            },
        },
    }

    elector, err := leaderelection.NewLeaderElector(lec)
    if err != nil {
        return err
    }

    elector.Run(ctx)
    return nil
}

// IsLeader 检查是否是领导者
func (le *LeaderElector) IsLeader() bool {
    le.mu.RLock()
    defer le.mu.RUnlock()
    return le.isLeader
}

// 基于 etcd 的 Leader Election（简化版）

type EtcdLeaderElector struct {
    client        *clientv3.Client
    electionName  string
    identity      string
    ttl           int
    session       *concurrency.Session
    election      *concurrency.Election

    isLeader bool
    mu       sync.RWMutex
}

func NewEtcdLeaderElector(client *clientv3.Client, electionName, identity string, ttl int) (*EtcdLeaderElector, error) {
    return &EtcdLeaderElector{
        client:       client,
        electionName: electionName,
        identity:     identity,
        ttl:          ttl,
    }, nil
}

func (e *EtcdLeaderElector) Run(ctx context.Context) error {
    // 创建会话
    session, err := concurrency.NewSession(e.client, concurrency.WithTTL(e.ttl))
    if err != nil {
        return err
    }
    defer session.Close()

    e.session = session
    e.election = concurrency.NewElection(session, e.electionName)

    // 参与选举
    campaignChan := make(chan error, 1)
    go func() {
        // 竞选领导者
        err := e.election.Campaign(ctx, e.identity)
        campaignChan <- err
    }()

    select {
    case err := <-campaignChan:
        if err != nil {
            return err
        }

        // 成为领导者
        e.mu.Lock()
        e.isLeader = true
        e.mu.Unlock()

        fmt.Printf("[%s] Became leader for election %s\n", e.identity, e.electionName)

        // 保持领导者身份
        <-ctx.Done()

    case <-ctx.Done():
        return ctx.Err()
    }

    return nil
}

func (e *EtcdLeaderElector) Resign(ctx context.Context) error {
    if e.election != nil {
        return e.election.Resign(ctx)
    }
    return nil
}

func (e *EtcdLeaderElector) IsLeader() bool {
    e.mu.RLock()
    defer e.mu.RUnlock()
    return e.isLeader
}

// LeaderAwareTask 领导者感知任务
type LeaderAwareTask struct {
    elector  *LeaderElector
    task     func(context.Context)
    interval time.Duration
}

func NewLeaderAwareTask(elector *LeaderElector, interval time.Duration, task func(context.Context)) *LeaderAwareTask {
    return &LeaderAwareTask{
        elector:  elector,
        task:     task,
        interval: interval,
    }
}

func (t *LeaderAwareTask) Run(ctx context.Context) {
    ticker := time.NewTicker(t.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if t.elector.IsLeader() {
                taskCtx, cancel := context.WithTimeout(ctx, t.interval)
                t.task(taskCtx)
                cancel()
            }
        case <-ctx.Done():
            return
        }
    }
}

// 使用示例
func ExampleLeaderElection() {
    // 创建 Kubernetes 客户端
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err)
    }

    client, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    // 配置选举
    leConfig := &Config{
        Client:        client,
        LockName:      "my-controller-leader",
        Namespace:     "default",
        Identity:      "pod-1",
        LeaseDuration: 15 * time.Second,
        RenewDeadline: 10 * time.Second,
        RetryPeriod:   2 * time.Second,
        Callbacks: LeaderCallbacks{
            OnStartedLeading: func(ctx context.Context) {
                fmt.Println("Started leading, starting controller...")
                // 启动控制器逻辑
            },
            OnStoppedLeading: func() {
                fmt.Println("Stopped leading, cleaning up...")
                // 清理资源
            },
            OnNewLeader: func(identity string) {
                fmt.Printf("New leader: %s\n", identity)
            },
        },
    }

    elector, err := NewLeaderElector(leConfig)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    elector.Run(ctx)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 高可用 | 选举需要时间 |
| 避免冲突 | 脑裂风险 |
| 自动故障转移 | 需要额外存储 |

#### 适用场景

- 控制器高可用
- 定时任务调度
- 独占资源访问

---

### 3.4 Distributed Lock 模式

#### 意图

Distributed Lock 模式用于在分布式环境中协调对共享资源的访问，防止竞争条件。

#### 结构

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│ Client 1│     │ Client 2│     │ Client 3│
└────┬────┘     └────┬────┘     └────┬────┘
     │               │               │
     └───────────────┼───────────────┘
                     │
              ┌──────┴──────┐
              │  Lock Store │
              │   (etcd)    │
              └─────────────┘
```

#### 实现

```go
// Distributed Lock 模式实现
package distributedlock

import (
    "context"
    "fmt"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
    "go.etcd.io/etcd/client/v3/concurrency"
)

// Lock 分布式锁接口
type Lock interface {
    Lock(ctx context.Context) error
    Unlock(ctx context.Context) error
    IsLocked() bool
}

// EtcdLock etcd 分布式锁
type EtcdLock struct {
    client    *clientv3.Client
    key       string
    ttl       int
    session   *concurrency.Session
    mutex     *concurrency.Mutex
    isLocked  bool
}

// NewEtcdLock 创建 etcd 锁
func NewEtcdLock(client *clientv3.Client, key string, ttl int) (*EtcdLock, error) {
    return &EtcdLock{
        client: client,
        key:    key,
        ttl:    ttl,
    }, nil
}

func (l *EtcdLock) Lock(ctx context.Context) error {
    // 创建会话
    session, err := concurrency.NewSession(l.client, concurrency.WithTTL(l.ttl))
    if err != nil {
        return err
    }

    l.session = session
    l.mutex = concurrency.NewMutex(session, l.key)

    // 获取锁
    if err := l.mutex.Lock(ctx); err != nil {
        session.Close()
        return err
    }

    l.isLocked = true
    return nil
}

func (l *EtcdLock) Unlock(ctx context.Context) error {
    if !l.isLocked {
        return nil
    }

    if l.mutex != nil {
        if err := l.mutex.Unlock(ctx); err != nil {
            return err
        }
    }

    if l.session != nil {
        l.session.Close()
    }

    l.isLocked = false
    return nil
}

func (l *EtcdLock) IsLocked() bool {
    return l.isLocked
}

// RedisLock Redis 分布式锁
type RedisLock struct {
    client *redis.Client
    key    string
    value  string
    ttl    time.Duration
}

func NewRedisLock(client *redis.Client, key string, ttl time.Duration) *RedisLock {
    return &RedisLock{
        client: client,
        key:    key,
        value:  generateUniqueID(),
        ttl:    ttl,
    }
}

func (r *RedisLock) Lock(ctx context.Context) error {
    for {
        ok, err := r.client.SetNX(ctx, r.key, r.value, r.ttl).Result()
        if err != nil {
            return err
        }

        if ok {
            return nil
        }

        // 等待后重试
        select {
        case <-time.After(100 * time.Millisecond):
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func (r *RedisLock) Unlock(ctx context.Context) error {
    // 使用 Lua 脚本确保原子性
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `

    result, err := r.client.Eval(ctx, script, []string{r.key}, r.value).Result()
    if err != nil {
        return err
    }

    if result.(int64) == 0 {
        return fmt.Errorf("lock not owned by this client")
    }

    return nil
}

func (r *RedisLock) IsLocked() bool {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    val, err := r.client.Get(ctx, r.key).Result()
    if err != nil {
        return false
    }

    return val == r.value
}

// LockManager 锁管理器
type LockManager struct {
    etcdClient  *clientv3.Client
    redisClient *redis.Client
    locks       map[string]Lock
}

func NewLockManager(etcdClient *clientv3.Client, redisClient *redis.Client) *LockManager {
    return &LockManager{
        etcdClient:  etcdClient,
        redisClient: redisClient,
        locks:       make(map[string]Lock),
    }
}

func (m *LockManager) AcquireEtcdLock(ctx context.Context, key string, ttl int) (Lock, error) {
    lock, err := NewEtcdLock(m.etcdClient, key, ttl)
    if err != nil {
        return nil, err
    }

    if err := lock.Lock(ctx); err != nil {
        return nil, err
    }

    m.locks[key] = lock
    return lock, nil
}

func (m *LockManager) AcquireRedisLock(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
    lock := NewRedisLock(m.redisClient, key, ttl)

    if err := lock.Lock(ctx); err != nil {
        return nil, err
    }

    m.locks[key] = lock
    return lock, nil
}

func (m *LockManager) ReleaseLock(ctx context.Context, key string) error {
    lock, ok := m.locks[key]
    if !ok {
        return fmt.Errorf("lock %s not found", key)
    }

    return lock.Unlock(ctx)
}

// WithLock 带锁执行函数
func WithLock(ctx context.Context, lock Lock, fn func() error) error {
    if err := lock.Lock(ctx); err != nil {
        return err
    }
    defer lock.Unlock(ctx)

    return fn()
}

// 辅助函数
func generateUniqueID() string {
    return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int())
}

// 使用示例
func ExampleDistributedLock() {
    // 创建 etcd 客户端
    etcdClient, err := clientv3.New(clientv3.Config{
        Endpoints: []string{"localhost:2379"},
    })
    if err != nil {
        panic(err)
    }
    defer etcdClient.Close()

    // 创建锁
    lock, err := NewEtcdLock(etcdClient, "/locks/resource-1", 10)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // 获取锁并执行操作
    err = WithLock(ctx, lock, func() error {
        fmt.Println("Lock acquired, performing operation...")
        time.Sleep(5 * time.Second)
        return nil
    })

    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 防止竞争条件 | 增加延迟 |
| 支持分布式 | 死锁风险 |
| 可超时释放 | 需要额外存储 |

#### 适用场景

- 资源分配
- 定时任务
- 配置更新

---

### 3.5 Barrier 模式

#### 意图

Barrier 模式用于协调多个并发任务，使它们在特定点同步，等待所有任务到达后再继续执行。

#### 结构

```
Task 1 ──► ┌─────────┐
Task 2 ──► │ Barrier │───► 所有任务到达后继续
Task 3 ──► │  (n=3)  │
Task 4 ──► └─────────┘
```

#### 实现

```go
// Barrier 模式实现
package barrier

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Barrier 屏障
type Barrier struct {
    count    int
    current  int
    mu       sync.Mutex
    cond     *sync.Cond
    broken   bool
}

// NewBarrier 创建屏障
func NewBarrier(count int) *Barrier {
    if count <= 0 {
        panic("count must be positive")
    }

    b := &Barrier{
        count: count,
    }
    b.cond = sync.NewCond(&b.mu)
    return b
}

// Wait 等待屏障
func (b *Barrier) Wait() error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.broken {
        return fmt.Errorf("barrier is broken")
    }

    b.current++

    if b.current >= b.count {
        // 最后一个到达，唤醒所有等待者
        b.current = 0
        b.cond.Broadcast()
        return nil
    }

    // 等待
    b.cond.Wait()

    if b.broken {
        return fmt.Errorf("barrier is broken")
    }

    return nil
}

// WaitWithTimeout 带超时的等待
func (b *Barrier) WaitWithTimeout(timeout time.Duration) error {
    done := make(chan struct{})
    var err error

    go func() {
        err = b.Wait()
        close(done)
    }()

    select {
    case <-done:
        return err
    case <-time.After(timeout):
        b.mu.Lock()
        b.broken = true
        b.mu.Unlock()
        b.cond.Broadcast()
        return fmt.Errorf("barrier wait timeout")
    }
}

// Reset 重置屏障
func (b *Barrier) Reset() {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.current = 0
    b.broken = false
}

// CyclicBarrier 循环屏障
type CyclicBarrier struct {
    count   int
    parties int
    mu      sync.Mutex
    cond    *sync.Cond
    action  func() // 到达屏障时执行的动作
}

// NewCyclicBarrier 创建循环屏障
func NewCyclicBarrier(parties int, action func()) *CyclicBarrier {
    cb := &CyclicBarrier{
        parties: parties,
        action:  action,
    }
    cb.cond = sync.NewCond(&cb.mu)
    return cb
}

func (cb *CyclicBarrier) Await() (int, error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    index := cb.parties - cb.count - 1
    cb.count++

    if cb.count >= cb.parties {
        // 最后一个到达
        if cb.action != nil {
            cb.action()
        }
        cb.count = 0
        cb.cond.Broadcast()
        return 0, nil
    }

    cb.cond.Wait()
    return index, nil
}

// CountDownLatch 倒计时门闩
type CountDownLatch struct {
    count int
    mu    sync.Mutex
    cond  *sync.Cond
}

// NewCountDownLatch 创建倒计时门闩
func NewCountDownLatch(count int) *CountDownLatch {
    cdl := &CountDownLatch{
        count: count,
    }
    cdl.cond = sync.NewCond(&cdl.mu)
    return cdl
}

// CountDown 计数减一
func (c *CountDownLatch) CountDown() {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.count > 0 {
        c.count--
        if c.count == 0 {
            c.cond.Broadcast()
        }
    }
}

// Await 等待计数归零
func (c *CountDownLatch) Await() {
    c.mu.Lock()
    defer c.mu.Unlock()

    for c.count > 0 {
        c.cond.Wait()
    }
}

// AwaitWithTimeout 带超时的等待
func (c *CountDownLatch) AwaitWithTimeout(timeout time.Duration) bool {
    done := make(chan struct{})

    go func() {
        c.Await()
        close(done)
    }()

    select {
    case <-done:
        return true
    case <-time.After(timeout):
        return false
    }
}

// GetCount 获取当前计数
func (c *CountDownLatch) GetCount() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

// Phaser 分阶段屏障
type Phaser struct {
    parties   int
    phase     int
    arrived   int
    mu        sync.Mutex
    cond      *sync.Cond
    onAdvance func(phase int, registeredParties int) bool
}

// NewPhaser 创建分阶段屏障
func NewPhaser(parties int, onAdvance func(phase int, registeredParties int) bool) *Phaser {
    p := &Phaser{
        parties:   parties,
        onAdvance: onAdvance,
    }
    p.cond = sync.NewCond(&p.mu)
    return p
}

// Register 注册参与者
func (p *Phaser) Register() int {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.parties++
    return p.phase
}

// ArriveAndAwaitAdvance 到达并等待
func (p *Phaser) ArriveAndAwaitAdvance() int {
    p.mu.Lock()
    defer p.mu.Unlock()

    currentPhase := p.phase
    p.arrived++

    if p.arrived >= p.parties {
        // 所有参与者到达
        if p.onAdvance != nil && !p.onAdvance(p.phase, p.parties) {
            // 终止
            return -1
        }

        p.phase++
        p.arrived = 0
        p.cond.Broadcast()
        return p.phase
    }

    // 等待
    for p.phase == currentPhase {
        p.cond.Wait()
    }

    return p.phase
}

// ArriveAndDeregister 到达并注销
func (p *Phaser) ArriveAndDeregister() int {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.arrived++
    p.parties--

    if p.arrived >= p.parties {
        p.phase++
        p.arrived = 0
        p.cond.Broadcast()
    }

    return p.phase
}

// 使用示例
func ExampleBarrier() {
    // 创建屏障（3个参与者）
    barrier := NewBarrier(3)

    var wg sync.WaitGroup

    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            fmt.Printf("Worker %d: Phase 1\n", id)
            time.Sleep(time.Duration(id*100) * time.Millisecond)

            // 等待所有工作到达屏障
            barrier.Wait()

            fmt.Printf("Worker %d: Phase 2\n", id)
        }(i)
    }

    wg.Wait()
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 同步协调 | 可能导致死锁 |
| 简单易用 | 性能开销 |
| 灵活多样 | 需要准确计数 |

#### 适用场景

- 并行计算
- 分阶段处理
- 测试协调

---

### 3.6 Pipeline 模式

#### 意图

Pipeline 模式将数据处理分解为多个阶段，每个阶段由独立的协程处理，数据在各个阶段之间流动。

#### 结构

```
Input ──► Stage 1 ──► Stage 2 ──► Stage 3 ──► Output
Source   (Parse)     (Transform)  (Save)      Sink
            │            │           │
            └────────────┴───────────┘
                    Channels
```

#### 实现

```go
// Pipeline 模式实现
package pipeline

import (
    "context"
    "fmt"
    "sync"
)

// Stage 管道阶段
type Stage func(in <-chan interface{}) <-chan interface{}

// Pipeline 管道
type Pipeline struct {
    stages []Stage
}

// NewPipeline 创建管道
func NewPipeline(stages ...Stage) *Pipeline {
    return &Pipeline{stages: stages}
}

// Run 运行管道
func (p *Pipeline) Run(ctx context.Context, source <-chan interface{}) <-chan interface{} {
    var current <-chan interface{} = source

    for _, stage := range p.stages {
        current = stage(current)
    }

    return current
}

// Generator 数据生成器
func Generator(ctx context.Context, values ...interface{}) <-chan interface{} {
    out := make(chan interface{})

    go func() {
        defer close(out)

        for _, v := range values {
            select {
            case out <- v:
            case <-ctx.Done():
                return
            }
        }
    }()

    return out
}

// FanOut 扇出
type FanOut struct {
    workers int
}

func NewFanOut(workers int) *FanOut {
    return &FanOut{workers: workers}
}

func (f *FanOut) Stage(fn func(interface{}) interface{}) Stage {
    return func(in <-chan interface{}) <-chan interface{} {
        var wg sync.WaitGroup

        // 创建多个输出通道
        outs := make([]chan interface{}, f.workers)
        for i := range outs {
            outs[i] = make(chan interface{})
        }

        // 分发任务
        wg.Add(1)
        go func() {
            defer wg.Done()

            i := 0
            for v := range in {
                select {
                case outs[i%f.workers] <- v:
                    i++
                }
            }

            for _, out := range outs {
                close(out)
            }
        }()

        // 每个工作协程处理自己的通道
        results := make([]chan interface{}, f.workers)
        for i := 0; i < f.workers; i++ {
            results[i] = make(chan interface{})
            wg.Add(1)

            go func(idx int) {
                defer wg.Done()
                defer close(results[idx])

                for v := range outs[idx] {
                    results[idx] <- fn(v)
                }
            }(i)
        }

        // 合并结果
        merged := make(chan interface{})
        go func() {
            defer close(merged)
            wg.Wait()

            for _, r := range results {
                for v := range r {
                    merged <- v
                }
            }
        }()

        return merged
    }
}

// FanIn 扇入
func FanIn(channels ...<-chan interface{}) <-chan interface{} {
    out := make(chan interface{})
    var wg sync.WaitGroup

    wg.Add(len(channels))

    for _, ch := range channels {
        go func(c <-chan interface{}) {
            defer wg.Done()

            for v := range c {
                out <- v
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

// BufferedStage 带缓冲的阶段
func BufferedStage(bufferSize int, fn func(interface{}) interface{}) Stage {
    return func(in <-chan interface{}) <-chan interface{} {
        out := make(chan interface{}, bufferSize)

        go func() {
            defer close(out)

            for v := range in {
                result := fn(v)
                out <- result
            }
        }()

        return out
    }
}

// FilterStage 过滤阶段
func FilterStage(predicate func(interface{}) bool) Stage {
    return func(in <-chan interface{}) <-chan interface{} {
        out := make(chan interface{})

        go func() {
            defer close(out)

            for v := range in {
                if predicate(v) {
                    out <- v
                }
            }
        }()

        return out
    }
}

// MapStage 映射阶段
func MapStage(fn func(interface{}) interface{}) Stage {
    return func(in <-chan interface{}) <-chan interface{} {
        out := make(chan interface{})

        go func() {
            defer close(out)

            for v := range in {
                out <- fn(v)
            }
        }()

        return out
    }
}

// ReduceStage 归约阶段
func ReduceStage(initial interface{}, fn func(acc, v interface{}) interface{}) Stage {
    return func(in <-chan interface{}) <-chan interface{} {
        out := make(chan interface{}, 1)

        go func() {
            defer close(out)

            acc := initial
            for v := range in {
                acc = fn(acc, v)
            }
            out <- acc
        }()

        return out
    }
}

// PipelineBuilder 管道构建器
type PipelineBuilder struct {
    stages []Stage
}

func NewPipelineBuilder() *PipelineBuilder {
    return &PipelineBuilder{
        stages: make([]Stage, 0),
    }
}

func (pb *PipelineBuilder) AddStage(stage Stage) *PipelineBuilder {
    pb.stages = append(pb.stages, stage)
    return pb
}

func (pb *PipelineBuilder) AddMapStage(fn func(interface{}) interface{}) *PipelineBuilder {
    return pb.AddStage(MapStage(fn))
}

func (pb *PipelineBuilder) AddFilterStage(predicate func(interface{}) bool) *PipelineBuilder {
    return pb.AddStage(FilterStage(predicate))
}

func (pb *PipelineBuilder) AddBufferedStage(bufferSize int, fn func(interface{}) interface{}) *PipelineBuilder {
    return pb.AddStage(BufferedStage(bufferSize, fn))
}

func (pb *PipelineBuilder) Build() *Pipeline {
    return NewPipeline(pb.stages...)
}

// 使用示例
func ExamplePipeline() {
    ctx := context.Background()

    // 构建管道
    pipeline := NewPipelineBuilder().
        AddFilterStage(func(v interface{}) bool {
            // 过滤偶数
            n, ok := v.(int)
            return ok && n%2 == 0
        }).
        AddMapStage(func(v interface{}) interface{} {
            // 平方
            n := v.(int)
            return n * n
        }).
        AddBufferedStage(10, func(v interface{}) interface{} {
            // 添加标签
            return map[string]interface{}{
                "value": v,
                "processed": true,
            }
        }).
        Build()

    // 生成数据
    source := make(chan interface{})
    go func() {
        defer close(source)
        for i := 1; i <= 20; i++ {
            source <- i
        }
    }()

    // 运行管道
    result := pipeline.Run(ctx, source)

    // 消费结果
    for v := range result {
        fmt.Printf("Result: %+v\n", v)
    }
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 解耦处理阶段 | 增加内存占用 |
| 提高吞吐量 | 调试困难 |
| 可组合 | 错误处理复杂 |

#### 适用场景

- 数据处理流水线
- ETL 作业
- 实时流处理

---



## 4. 同步与异步模式

同步与异步模式决定了系统组件之间的通信方式，直接影响系统的性能、可靠性和可扩展性。

---

### 4.1 同步调用模式

#### 意图

同步调用模式使调用方等待被调用方完成操作并返回结果后再继续执行。

#### 结构

```
Caller ──► Call ──► Wait ──► Response ──► Continue
                    │
                    └── 阻塞等待
```

#### 实现

```go
// 同步调用模式实现
package syncpattern

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

// SyncClient 同步客户端
type SyncClient struct {
    baseURL    string
    httpClient *http.Client
    timeout    time.Duration
}

// NewSyncClient 创建同步客户端
func NewSyncClient(baseURL string, timeout time.Duration) *SyncClient {
    return &SyncClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: timeout,
        },
        timeout: timeout,
    }
}

// Call 同步调用
func (c *SyncClient) Call(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
    url := c.baseURL + path

    req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")

    // 同步等待响应
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }

    return resp, nil
}

// CallWithRetry 带重试的同步调用
func (c *SyncClient) CallWithRetry(ctx context.Context, method, path string, body []byte, maxRetries int) (*http.Response, error) {
    var lastErr error

    for i := 0; i <= maxRetries; i++ {
        if i > 0 {
            // 指数退避
            backoff := time.Duration(i*i) * 100 * time.Millisecond
            select {
            case <-time.After(backoff):
            case <-ctx.Done():
                return nil, ctx.Err()
            }
        }

        resp, err := c.Call(ctx, method, path, body)
        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }

        if err != nil {
            lastErr = err
        } else {
            lastErr = fmt.Errorf("HTTP error: %d", resp.StatusCode)
            resp.Body.Close()
        }
    }

    return nil, fmt.Errorf("exhausted retries: %w", lastErr)
}

// ServiceClient 服务客户端
type ServiceClient struct {
    endpoints map[string]*SyncClient
    balancer  LoadBalancer
}

type LoadBalancer interface {
    Select(endpoints []string) string
}

func NewServiceClient(endpoints []string, timeout time.Duration) *ServiceClient {
    clients := make(map[string]*SyncClient)
    for _, ep := range endpoints {
        clients[ep] = NewSyncClient(ep, timeout)
    }

    return &ServiceClient{
        endpoints: clients,
        balancer:  &RoundRobinBalancer{},
    }
}

func (s *ServiceClient) Call(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
    // 选择端点
    endpoints := make([]string, 0, len(s.endpoints))
    for ep := range s.endpoints {
        endpoints = append(endpoints, ep)
    }

    selected := s.balancer.Select(endpoints)
    client := s.endpoints[selected]

    return client.Call(ctx, method, path, body)
}

// RoundRobinBalancer 轮询负载均衡器
type RoundRobinBalancer struct {
    current uint64
}

func (r *RoundRobinBalancer) Select(endpoints []string) string {
    if len(endpoints) == 0 {
        return ""
    }

    idx := atomic.AddUint64(&r.current, 1) % uint64(len(endpoints))
    return endpoints[idx]
}

// CircuitBreakerSyncClient 带熔断的同步客户端
type CircuitBreakerSyncClient struct {
    client       *SyncClient
    circuitBreaker *CircuitBreaker
}

func NewCircuitBreakerSyncClient(client *SyncClient, cb *CircuitBreaker) *CircuitBreakerSyncClient {
    return &CircuitBreakerSyncClient{
        client:         client,
        circuitBreaker: cb,
    }
}

func (c *CircuitBreakerSyncClient) Call(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
    var resp *http.Response
    var err error

    execErr := c.circuitBreaker.Execute(ctx, func() error {
        resp, err = c.client.Call(ctx, method, path, body)
        if err != nil {
            return err
        }

        if resp.StatusCode >= 500 {
            return fmt.Errorf("server error: %d", resp.StatusCode)
        }

        return nil
    })

    if execErr != nil {
        return nil, execErr
    }

    return resp, err
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 简单直观 | 阻塞等待 |
| 立即获得结果 | 级联延迟 |
| 易于调试 | 故障传播 |

#### 适用场景

- 需要立即结果
- 简单查询操作
- 事务处理

---

### 4.2 异步消息模式

#### 意图

异步消息模式使调用方无需等待被调用方完成，通过消息队列或事件机制进行通信。

#### 结构

```
Producer ──► Queue/Topic ──► Consumer
    │                           │
    └── 立即返回               └── 异步处理
```

#### 实现

```go
// 异步消息模式实现
package asyncpattern

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// Message 消息
type Message struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    Payload     map[string]interface{} `json:"payload"`
    Timestamp   time.Time              `json:"timestamp"`
    CorrelationID string               `json:"correlation_id,omitempty"`
    ReplyTo     string                 `json:"reply_to,omitempty"`
}

// AsyncProducer 异步生产者
type AsyncProducer struct {
    queue     chan *Message
    maxSize   int
    mu        sync.RWMutex
    handlers  map[string][]MessageHandler
}

type MessageHandler func(*Message) error

// NewAsyncProducer 创建异步生产者
func NewAsyncProducer(maxSize int) *AsyncProducer {
    return &AsyncProducer{
        queue:    make(chan *Message, maxSize),
        maxSize:  maxSize,
        handlers: make(map[string][]MessageHandler),
    }
}

// Send 发送消息（异步）
func (p *AsyncProducer) Send(msg *Message) error {
    select {
    case p.queue <- msg:
        return nil
    default:
        return fmt.Errorf("message queue is full")
    }
}

// SendWithTimeout 带超时的发送
func (p *AsyncProducer) SendWithTimeout(msg *Message, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    select {
    case p.queue <- msg:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// RegisterHandler 注册处理器
func (p *AsyncProducer) RegisterHandler(msgType string, handler MessageHandler) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.handlers[msgType] = append(p.handlers[msgType], handler)
}

// Start 启动消费者
func (p *AsyncProducer) Start(ctx context.Context, workers int) {
    for i := 0; i < workers; i++ {
        go p.consume(ctx)
    }
}

func (p *AsyncProducer) consume(ctx context.Context) {
    for {
        select {
        case msg := <-p.queue:
            p.processMessage(msg)
        case <-ctx.Done():
            return
        }
    }
}

func (p *AsyncProducer) processMessage(msg *Message) {
    p.mu.RLock()
    handlers := p.handlers[msg.Type]
    p.mu.RUnlock()

    for _, handler := range handlers {
        if err := handler(msg); err != nil {
            fmt.Printf("Handler error: %v\n", err)
        }
    }
}

// AsyncRequestReply 异步请求-响应
type AsyncRequestReply struct {
    producer    *AsyncProducer
    pending     map[string]chan *Message
    mu          sync.RWMutex
    timeout     time.Duration
}

// NewAsyncRequestReply 创建异步请求-响应
func NewAsyncRequestReply(producer *AsyncProducer, timeout time.Duration) *AsyncRequestReply {
    arr := &AsyncRequestReply{
        producer: producer,
        pending:  make(map[string]chan *Message),
        timeout:  timeout,
    }

    // 注册响应处理器
    producer.RegisterHandler("response", arr.handleResponse)

    return arr
}

func (a *AsyncRequestReply) Request(ctx context.Context, msg *Message) (*Message, error) {
    // 创建响应通道
    replyChan := make(chan *Message, 1)

    a.mu.Lock()
    a.pending[msg.ID] = replyChan
    a.mu.Unlock()

    defer func() {
        a.mu.Lock()
        delete(a.pending, msg.ID)
        a.mu.Unlock()
    }()

    // 设置回复地址
    msg.ReplyTo = "response"

    // 发送请求
    if err := a.producer.Send(msg); err != nil {
        return nil, err
    }

    // 等待响应
    select {
    case resp := <-replyChan:
        return resp, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    case <-time.After(a.timeout):
        return nil, fmt.Errorf("request timeout")
    }
}

func (a *AsyncRequestReply) handleResponse(msg *Message) error {
    a.mu.RLock()
    replyChan, ok := a.pending[msg.CorrelationID]
    a.mu.RUnlock()

    if !ok {
        return fmt.Errorf("no pending request for correlation ID: %s", msg.CorrelationID)
    }

    replyChan <- msg
    return nil
}

// EventBus 事件总线
type EventBus struct {
    subscribers map[string][]EventHandler
    mu          sync.RWMutex
    eventChan   chan *Event
}

type Event struct {
    Topic     string
    Data      interface{}
    Timestamp time.Time
}

type EventHandler func(*Event) error

// NewEventBus 创建事件总线
func NewEventBus(bufferSize int) *EventBus {
    return &EventBus{
        subscribers: make(map[string][]EventHandler),
        eventChan:   make(chan *Event, bufferSize),
    }
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(topic string, handler EventHandler) {
    eb.mu.Lock()
    defer eb.mu.Unlock()

    eb.subscribers[topic] = append(eb.subscribers[topic], handler)
}

// Publish 发布事件
func (eb *EventBus) Publish(event *Event) error {
    select {
    case eb.eventChan <- event:
        return nil
    default:
        return fmt.Errorf("event channel is full")
    }
}

// Start 启动事件分发
func (eb *EventBus) Start(ctx context.Context) {
    go func() {
        for {
            select {
            case event := <-eb.eventChan:
                eb.dispatch(event)
            case <-ctx.Done():
                return
            }
        }
    }()
}

func (eb *EventBus) dispatch(event *Event) {
    eb.mu.RLock()
    handlers := eb.subscribers[event.Topic]
    eb.mu.RUnlock()

    for _, handler := range handlers {
        go func(h EventHandler) {
            if err := h(event); err != nil {
                fmt.Printf("Event handler error: %v\n", err)
            }
        }(handler)
    }
}

// 使用示例
func ExampleAsyncPattern() {
    // 创建异步生产者
    producer := NewAsyncProducer(1000)

    // 注册处理器
    producer.RegisterHandler("order.created", func(msg *Message) error {
        fmt.Printf("Processing order: %+v\n", msg.Payload)
        return nil
    })

    // 启动消费者
    ctx := context.Background()
    producer.Start(ctx, 5)

    // 发送消息
    for i := 0; i < 10; i++ {
        msg := &Message{
            ID:        fmt.Sprintf("msg-%d", i),
            Type:      "order.created",
            Payload:   map[string]interface{}{"order_id": i, "amount": 100.0},
            Timestamp: time.Now(),
        }

        if err := producer.Send(msg); err != nil {
            fmt.Printf("Send error: %v\n", err)
        }
    }

    time.Sleep(2 * time.Second)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 解耦组件 | 增加复杂度 |
| 提高吞吐量 | 最终一致性 |
| 削峰填谷 | 需要消息系统 |

#### 适用场景

- 事件驱动架构
- 高并发处理
- 削峰填谷

---

### 4.3 Watch 机制

#### 意图

Watch 机制允许客户端监听资源的变化，当资源发生变化时，服务器主动推送通知。

#### 结构

```
Client ──► Watch Request ──► Server
    │                            │
    │◄────── Event Stream ──────┤
    │                            │
    └── 持续监听资源变化         └── 推送变更事件
```

#### 实现

```go
// Watch 机制实现
package watch

import (
    "context"
    "fmt"
    "sync"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/watch"
    "k8s.io/client-go/kubernetes"
)

// Watcher 观察者
type Watcher struct {
    client    kubernetes.Interface
    watchers  map[string]watch.Interface
    handlers  map[string][]EventHandler
    mu        sync.RWMutex
}

type EventHandler func(event watch.Event) error

// NewWatcher 创建观察者
func NewWatcher(client kubernetes.Interface) *Watcher {
    return &Watcher{
        client:   client,
        watchers: make(map[string]watch.Interface),
        handlers: make(map[string][]EventHandler),
    }
}

// WatchPods 监听 Pod 变化
func (w *Watcher) WatchPods(ctx context.Context, namespace string, selector string) error {
    listOptions := metav1.ListOptions{
        LabelSelector: selector,
    }

    watcher, err := w.client.CoreV1().Pods(namespace).Watch(ctx, listOptions)
    if err != nil {
        return err
    }

    w.mu.Lock()
    w.watchers["pods"] = watcher
    w.mu.Unlock()

    // 处理事件
    go w.handleEvents("pods", watcher.ResultChan())

    return nil
}

// WatchDeployments 监听 Deployment 变化
func (w *Watcher) WatchDeployments(ctx context.Context, namespace string) error {
    watcher, err := w.client.AppsV1().Deployments(namespace).Watch(ctx, metav1.ListOptions{})
    if err != nil {
        return err
    }

    w.mu.Lock()
    w.watchers["deployments"] = watcher
    w.mu.Unlock()

    go w.handleEvents("deployments", watcher.ResultChan())

    return nil
}

// WatchServices 监听 Service 变化
func (w *Watcher) WatchServices(ctx context.Context, namespace string) error {
    watcher, err := w.client.CoreV1().Services(namespace).Watch(ctx, metav1.ListOptions{})
    if err != nil {
        return err
    }

    w.mu.Lock()
    w.watchers["services"] = watcher
    w.mu.Unlock()

    go w.handleEvents("services", watcher.ResultChan())

    return nil
}

// RegisterHandler 注册事件处理器
func (w *Watcher) RegisterHandler(resourceType string, handler EventHandler) {
    w.mu.Lock()
    defer w.mu.Unlock()

    w.handlers[resourceType] = append(w.handlers[resourceType], handler)
}

func (w *Watcher) handleEvents(resourceType string, eventChan <-chan watch.Event) {
    for event := range eventChan {
        w.mu.RLock()
        handlers := w.handlers[resourceType]
        w.mu.RUnlock()

        for _, handler := range handlers {
            if err := handler(event); err != nil {
                fmt.Printf("Handler error: %v\n", err)
            }
        }
    }
}

// Stop 停止监听
func (w *Watcher) Stop(resourceType string) {
    w.mu.Lock()
    defer w.mu.Unlock()

    if watcher, ok := w.watchers[resourceType]; ok {
        watcher.Stop()
        delete(w.watchers, resourceType)
    }
}

// StopAll 停止所有监听
func (w *Watcher) StopAll() {
    w.mu.Lock()
    defer w.mu.Unlock()

    for _, watcher := range w.watchers {
        watcher.Stop()
    }

    w.watchers = make(map[string]watch.Interface)
}

// ResourceVersionWatcher 基于资源版本的观察者
type ResourceVersionWatcher struct {
    client          kubernetes.Interface
    resourceVersion string
    mu              sync.RWMutex
}

func NewResourceVersionWatcher(client kubernetes.Interface) *ResourceVersionWatcher {
    return &ResourceVersionWatcher{
        client: client,
    }
}

// WatchFromVersion 从指定版本开始监听
func (rw *ResourceVersionWatcher) WatchFromVersion(ctx context.Context, resourceType, namespace, resourceVersion string) (watch.Interface, error) {
    listOptions := metav1.ListOptions{
        ResourceVersion: resourceVersion,
        Watch:           true,
    }

    switch resourceType {
    case "pods":
        return rw.client.CoreV1().Pods(namespace).Watch(ctx, listOptions)
    case "deployments":
        return rw.client.AppsV1().Deployments(namespace).Watch(ctx, listOptions)
    case "services":
        return rw.client.CoreV1().Services(namespace).Watch(ctx, listOptions)
    default:
        return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
    }
}

// BookmarkWatcher 书签观察者（处理断连重连）
type BookmarkWatcher struct {
    *Watcher
    bookmarks map[string]string // resourceType -> resourceVersion
    mu        sync.RWMutex
}

func NewBookmarkWatcher(client kubernetes.Interface) *BookmarkWatcher {
    return &BookmarkWatcher{
        Watcher:   NewWatcher(client),
        bookmarks: make(map[string]string),
    }
}

func (bw *BookmarkWatcher) handleEventWithBookmark(resourceType string, event watch.Event) error {
    // 更新书签
    if obj, ok := event.Object.(metav1.Object); ok {
        bw.mu.Lock()
        // 从对象获取 resource version
        // 实际实现需要类型断言获取具体类型的 ResourceVersion
        bw.bookmarks[resourceType] = obj.GetResourceVersion()
        bw.mu.Unlock()
    }

    // 调用原始处理器
    return nil
}

// GetBookmark 获取书签
func (bw *BookmarkWatcher) GetBookmark(resourceType string) string {
    bw.mu.RLock()
    defer bw.mu.RUnlock()
    return bw.bookmarks[resourceType]
}

// 使用示例
func ExampleWatch() {
    // 创建 Kubernetes 客户端
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err)
    }

    client, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    // 创建观察者
    watcher := NewWatcher(client)

    // 注册处理器
    watcher.RegisterHandler("pods", func(event watch.Event) error {
        pod := event.Object.(*corev1.Pod)
        fmt.Printf("Pod event: %s - %s/%s\n", event.Type, pod.Namespace, pod.Name)
        return nil
    })

    // 开始监听
    ctx := context.Background()
    if err := watcher.WatchPods(ctx, "default", ""); err != nil {
        panic(err)
    }

    // 保持运行
    select {}
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 实时响应 | 连接管理复杂 |
| 减少轮询 | 需要处理断连 |
| 资源高效 | 状态同步挑战 |

#### 适用场景

- 配置中心
- 服务发现
- 实时监控

---

### 4.4 Informer 模式

#### 意图

Informer 模式结合了 Watch 机制和本地缓存，提供高效的资源变更通知和查询能力。

#### 结构

```
┌─────────────────────────────────────────────────────────┐
│                      Informer                           │
│                                                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │   Reflector │───►│ Local Cache │───►│  Indexer    │ │
│  │             │    │  (Store)    │    │             │ │
│  └──────┬──────┘    └─────────────┘    └─────────────┘ │
│         │                                               │
│         │ Watch API Server                              │
│         │                                               │
│  ┌──────┴──────┐                                        │
│  │ Work Queue  │                                        │
│  └──────┬──────┘                                        │
│         │                                               │
│         ▼                                               │
│  ┌─────────────┐                                        │
│  │  Processor  │───► Event Handlers                     │
│  └─────────────┘                                        │
└─────────────────────────────────────────────────────────┘
```

#### 实现

```go
// Informer 模式实现
package informer

import (
    "context"
    "fmt"
    "sync"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/watch"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
)

// Informer 资源 Informer
type Informer struct {
    client        kubernetes.Interface
    resourceType  string
    namespace     string
    listWatcher   *cache.ListWatch
    indexer       cache.Indexer
    controller    cache.Controller
    handlers      []ResourceEventHandler
    mu            sync.RWMutex
    resyncPeriod  time.Duration
}

// ResourceEventHandler 资源事件处理器
type ResourceEventHandler interface {
    OnAdd(obj interface{})
    OnUpdate(oldObj, newObj interface{})
    OnDelete(obj interface{})
}

// ResourceEventHandlerFuncs 函数式处理器
type ResourceEventHandlerFuncs struct {
    AddFunc    func(obj interface{})
    UpdateFunc func(oldObj, newObj interface{})
    DeleteFunc func(obj interface{})
}

func (r *ResourceEventHandlerFuncs) OnAdd(obj interface{}) {
    if r.AddFunc != nil {
        r.AddFunc(obj)
    }
}

func (r *ResourceEventHandlerFuncs) OnUpdate(oldObj, newObj interface{}) {
    if r.UpdateFunc != nil {
        r.UpdateFunc(oldObj, newObj)
    }
}

func (r *ResourceEventHandlerFuncs) OnDelete(obj interface{}) {
    if r.DeleteFunc != nil {
        r.DeleteFunc(obj)
    }
}

// NewPodInformer 创建 Pod Informer
func NewPodInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration) *Informer {
    listWatcher := &cache.ListWatch{
        ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
            return client.CoreV1().Pods(namespace).List(context.TODO(), options)
        },
        WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
            return client.CoreV1().Pods(namespace).Watch(context.TODO(), options)
        },
    }

    indexer, controller := cache.NewIndexerInformer(
        listWatcher,
        &corev1.Pod{},
        resyncPeriod,
        cache.ResourceEventHandlerFuncs{},
        cache.Indexers{
            cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
        },
    )

    return &Informer{
        client:       client,
        resourceType: "pods",
        namespace:    namespace,
        listWatcher:  listWatcher,
        indexer:      indexer,
        controller:   controller,
        resyncPeriod: resyncPeriod,
    }
}

// AddEventHandler 添加事件处理器
func (i *Informer) AddEventHandler(handler ResourceEventHandler) {
    i.mu.Lock()
    defer i.mu.Unlock()

    i.handlers = append(i.handlers, handler)
}

// Run 启动 Informer
func (i *Informer) Run(ctx context.Context) {
    // 包装处理器
    wrappedHandler := cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            i.mu.RLock()
            handlers := i.handlers
            i.mu.RUnlock()

            for _, h := range handlers {
                h.OnAdd(obj)
            }
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            i.mu.RLock()
            handlers := i.handlers
            i.mu.RUnlock()

            for _, h := range handlers {
                h.OnUpdate(oldObj, newObj)
            }
        },
        DeleteFunc: func(obj interface{}) {
            i.mu.RLock()
            handlers := i.handlers
            i.mu.RUnlock()

            for _, h := range handlers {
                h.OnDelete(obj)
            }
        },
    }

    // 重新创建 controller 以使用新的 handler
    indexer, controller := cache.NewIndexerInformer(
        i.listWatcher,
        &corev1.Pod{},
        i.resyncPeriod,
        wrappedHandler,
        cache.Indexers{
            cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
        },
    )

    i.indexer = indexer
    i.controller = controller

    // 启动 controller
    go i.controller.Run(ctx.Done())
}

// GetByKey 通过 key 获取对象
func (i *Informer) GetByKey(key string) (interface{}, bool, error) {
    return i.indexer.GetByKey(key)
}

// List 列出所有对象
func (i *Informer) List() []interface{} {
    return i.indexer.List()
}

// ListByNamespace 按命名空间列出
func (i *Informer) ListByNamespace(namespace string) ([]interface{}, error) {
    return i.indexer.ByIndex(cache.NamespaceIndex, namespace)
}

// HasSynced 检查是否已同步
func (i *Informer) HasSynced() bool {
    return i.controller.HasSynced()
}

// SharedInformer 共享 Informer
type SharedInformer struct {
    informers map[string]*Informer
    mu        sync.RWMutex
}

func NewSharedInformer() *SharedInformer {
    return &SharedInformer{
        informers: make(map[string]*Informer),
    }
}
}

func (s *SharedInformer) InformerFor(resourceType string, newFunc func() *Informer) *Informer {
    s.mu.Lock()
    defer s.mu.Unlock()

    informer, exists := s.informers[resourceType]
    if !exists {
        informer = newFunc()
        s.informers[resourceType] = informer
    }

    return informer
}

func (s *SharedInformer) Start(ctx context.Context) {
    s.mu.RLock()
    informers := make([]*Informer, 0, len(s.informers))
    for _, inf := range s.informers {
        informers = append(informers, inf)
    }
    s.mu.RUnlock()

    for _, inf := range informers {
        go inf.Run(ctx)
    }
}

// DeltaFIFO Delta FIFO 队列
type DeltaFIFO struct {
    queue []Delta
    mu    sync.Mutex
    cond  *sync.Cond
    knownObjects KeyListerGetter
}

type Delta struct {
    Type   DeltaType
    Object interface{}
}

type DeltaType string

const (
    Added   DeltaType = "Added"
    Updated DeltaType = "Updated"
    Deleted DeltaType = "Deleted"
    Sync    DeltaType = "Sync"
)

type KeyListerGetter interface {
    ListKeys() []string
    GetByKey(key string) (interface{}, bool, error)
}

func NewDeltaFIFO(knownObjects KeyListerGetter) *DeltaFIFO {
    f := &DeltaFIFO{
        knownObjects: knownObjects,
    }
    f.cond = sync.NewCond(&f.mu)
    return f
}

func (f *DeltaFIFO) Add(obj interface{}) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    f.queue = append(f.queue, Delta{Type: Added, Object: obj})
    f.cond.Broadcast()
    return nil
}

func (f *DeltaFIFO) Update(obj interface{}) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    f.queue = append(f.queue, Delta{Type: Updated, Object: obj})
    f.cond.Broadcast()
    return nil
}

func (f *DeltaFIFO) Delete(obj interface{}) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    f.queue = append(f.queue, Delta{Type: Deleted, Object: obj})
    f.cond.Broadcast()
    return nil
}

func (f *DeltaFIFO) Pop(process func(interface{}) error) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    for len(f.queue) == 0 {
        f.cond.Wait()
    }

    delta := f.queue[0]
    f.queue = f.queue[1:]

    f.mu.Unlock()
    err := process(delta)
    f.mu.Lock()

    return err
}

// Reflector Reflector 实现
type Reflector struct {
    name          string
    expectedType  interface{}
    store         cache.Store
    listerWatcher *cache.ListWatch
    period        time.Duration
}

func NewReflector(name string, expectedType interface{}, store cache.Store, listerWatcher *cache.ListWatch, period time.Duration) *Reflector {
    return &Reflector{
        name:          name,
        expectedType:  expectedType,
        store:         store,
        listerWatcher: listerWatcher,
        period:        period,
    }
}

func (r *Reflector) Run(ctx context.Context) {
    // 首次 List
    if err := r.listAndWatch(ctx); err != nil {
        fmt.Printf("Reflector error: %v\n", err)
    }

    // 定期重新同步
    ticker := time.NewTicker(r.period)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := r.listAndWatch(ctx); err != nil {
                fmt.Printf("Reflector error: %v\n", err)
            }
        case <-ctx.Done():
            return
        }
    }
}

func (r *Reflector) listAndWatch(ctx context.Context) error {
    // List 所有资源
    list, err := r.listerWatcher.List(metav1.ListOptions{})
    if err != nil {
        return err
    }

    // 同步到 store
    items, err := meta.ExtractList(list)
    if err != nil {
        return err
    }

    for _, item := range items {
        r.store.Update(item)
    }

    // 获取 resource version
    resourceVersion, err := meta.NewAccessor().ResourceVersion(list)
    if err != nil {
        return err
    }

    // Watch 变更
    watcher, err := r.listerWatcher.Watch(metav1.ListOptions{
        ResourceVersion: resourceVersion,
    })
    if err != nil {
        return err
    }
    defer watcher.Stop()

    // 处理事件
    for event := range watcher.ResultChan() {
        switch event.Type {
        case watch.Added:
            r.store.Add(event.Object)
        case watch.Modified:
            r.store.Update(event.Object)
        case watch.Deleted:
            r.store.Delete(event.Object)
        }
    }

    return nil
}

// 使用示例
func ExampleInformer() {
    // 创建 Kubernetes 客户端
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err)
    }

    client, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    // 创建 Informer
    informer := NewPodInformer(client, "default", 30*time.Second)

    // 添加处理器
    informer.AddEventHandler(&ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            pod := obj.(*corev1.Pod)
            fmt.Printf("Pod added: %s\n", pod.Name)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            pod := newObj.(*corev1.Pod)
            fmt.Printf("Pod updated: %s, phase: %s\n", pod.Name, pod.Status.Phase)
        },
        DeleteFunc: func(obj interface{}) {
            pod := obj.(*corev1.Pod)
            fmt.Printf("Pod deleted: %s\n", pod.Name)
        },
    })

    // 启动 Informer
    ctx := context.Background()
    informer.Run(ctx)

    // 等待同步完成
    for !informer.HasSynced() {
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Informer synced")

    // 查询缓存
    pods := informer.List()
    fmt.Printf("Cached pods: %d\n", len(pods))

    select {}
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 高性能查询 | 内存占用 |
| 实时更新 | 缓存一致性 |
| 减少 API 调用 | 实现复杂 |

#### 适用场景

- 控制器开发
- 资源缓存
- 事件处理

---

### 4.5 Work Queue 模式

#### 意图

Work Queue 模式提供可靠的任务队列，支持延迟、速率限制和重试机制。

#### 结构

```
Producer ──► ┌─────────────┐ ──► Consumer
             │  Work Queue │
             │             │
             │ - Delayed   │
             │ - RateLimit │
             │ - Retry     │
             └─────────────┘
```

#### 实现

```go
// Work Queue 模式实现
package workqueue

import (
    "container/heap"
    "context"
    "fmt"
    "sync"
    "time"
)

// Item 队列项
type Item struct {
    Key       string
    Timestamp time.Time
    Attempts  int
    Index     int // 用于堆
}

// Interface 工作队列接口
type Interface interface {
    Add(item interface{})
    Len() int
    Get() (item interface{}, shutdown bool)
    Done(item interface{})
    ShutDown()
    ShuttingDown() bool
}

// Type 基础队列实现
type Type struct {
    queue      []interface{}
    dirty      map[interface{}]struct{}
    processing map[interface{}]struct{}
    cond       *sync.Cond
    shuttingDown bool
}

// New 创建工作队列
func New() *Type {
    t := &Type{
        dirty:      make(map[interface{}]struct{}),
        processing: make(map[interface{}]struct{}),
    }
    t.cond = sync.NewCond(&sync.Mutex{})
    return t
}

func (t *Type) Add(item interface{}) {
    t.cond.L.Lock()
    defer t.cond.L.Unlock()

    if t.shuttingDown {
        return
    }

    if _, ok := t.dirty[item]; ok {
        return
    }

    t.dirty[item] = struct{}{}

    if _, ok := t.processing[item]; ok {
        return
    }

    t.queue = append(t.queue, item)
    t.cond.Signal()
}

func (t *Type) Len() int {
    t.cond.L.Lock()
    defer t.cond.L.Unlock()
    return len(t.queue)
}

func (t *Type) Get() (interface{}, bool) {
    t.cond.L.Lock()
    defer t.cond.L.Unlock()

    for len(t.queue) == 0 && !t.shuttingDown {
        t.cond.Wait()
    }

    if len(t.queue) == 0 {
        return nil, true
    }

    item := t.queue[0]
    t.queue = t.queue[1:]

    delete(t.dirty, item)
    t.processing[item] = struct{}{}

    return item, false
}

func (t *Type) Done(item interface{}) {
    t.cond.L.Lock()
    defer t.cond.L.Unlock()

    delete(t.processing, item)

    if _, ok := t.dirty[item]; ok {
        t.queue = append(t.queue, item)
        t.cond.Signal()
    }
}

func (t *Type) ShutDown() {
    t.cond.L.Lock()
    defer t.cond.L.Unlock()

    t.shuttingDown = true
    t.cond.Broadcast()
}

func (t *Type) ShuttingDown() bool {
    t.cond.L.Lock()
    defer t.cond.L.Unlock()
    return t.shuttingDown
}

// DelayingInterface 延迟队列接口
type DelayingInterface interface {
    Interface
    AddAfter(item interface{}, duration time.Duration)
}

// DelayingType 延迟队列实现
type DelayingType struct {
    *Type
    waitingForAddCh chan *waitFor
    stopCh          chan struct{}
}

type waitFor struct {
    data    interface{}
    readyAt time.Time
}

type waitForPriorityQueue []*waitFor

func (pq waitForPriorityQueue) Len() int { return len(pq) }

func (pq waitForPriorityQueue) Less(i, j int) bool {
    return pq[i].readyAt.Before(pq[j].readyAt)
}

func (pq waitForPriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
}

func (pq *waitForPriorityQueue) Push(x interface{}) {
    item := x.(*waitFor)
    *pq = append(*pq, item)
}

func (pq *waitForPriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    *pq = old[:n-1]
    return item
}

// NewDelayingQueue 创建延迟队列
func NewDelayingQueue() DelayingInterface {
    dq := &DelayingType{
        Type:            New(),
        waitingForAddCh: make(chan *waitFor, 1000),
        stopCh:          make(chan struct{}),
    }

    go dq.waitingLoop()

    return dq
}

func (dq *DelayingType) AddAfter(item interface{}, duration time.Duration) {
    if duration <= 0 {
        dq.Add(item)
        return
    }

    select {
    case dq.waitingForAddCh <- &waitFor{data: item, readyAt: time.Now().Add(duration)}:
    case <-dq.stopCh:
    }
}

func (dq *DelayingType) waitingLoop() {
    waitingForQueue := &waitForPriorityQueue{}
    heap.Init(waitingForQueue)

    nextReadyAtTimer := time.NewTimer(0)
    <-nextReadyAtTimer.C

    for {
        if dq.Type.ShamingDown() {
            return
        }

        now := time.Now()

        // 处理到期的项
        for waitingForQueue.Len() > 0 {
            item := (*waitingForQueue)[0]
            if item.readyAt.After(now) {
                break
            }

            heap.Pop(waitingForQueue)
            dq.Type.Add(item.data)
        }

        // 设置下一个定时器
        nextReadyAt := now.Add(1 * time.Second)
        if waitingForQueue.Len() > 0 {
            nextReadyAt = (*waitingForQueue)[0].readyAt
        }

        nextReadyAtTimer.Reset(nextReadyAt.Sub(now))

        select {
        case waitEntry := <-dq.waitingForAddCh:
            if waitEntry.readyAt.After(now) {
                heap.Push(waitingForQueue, waitEntry)
            } else {
                dq.Type.Add(waitEntry.data)
            }

        case <-nextReadyAtTimer.C:

        case <-dq.stopCh:
            return
        }
    }
}

// RateLimitingInterface 速率限制队列接口
type RateLimitingInterface interface {
    DelayingInterface
    AddRateLimited(item interface{})
    Forget(item interface{})
    NumRequeues(item interface{}) int
}

// RateLimiter 速率限制器接口
type RateLimiter interface {
    When(item interface{}) time.Duration
    Forget(item interface{})
    NumRequeues(item interface{}) int
}

// ItemExponentialFailureRateLimiter 指数退避速率限制器
type ItemExponentialFailureRateLimiter struct {
    failures     map[interface{}]int
    baseDelay    time.Duration
    maxDelay     time.Duration
    mu           sync.Mutex
}

func NewItemExponentialFailureRateLimiter(baseDelay, maxDelay time.Duration) *ItemExponentialFailureRateLimiter {
    return &ItemExponentialFailureRateLimiter{
        failures:  make(map[interface{}]int),
        baseDelay: baseDelay,
        maxDelay:  maxDelay,
    }
}

func (r *ItemExponentialFailureRateLimiter) When(item interface{}) time.Duration {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.failures[item]++

    exp := r.failures[item] - 1
    if exp < 0 {
        exp = 0
    }

    delay := r.baseDelay
    for i := 0; i < exp; i++ {
        delay *= 2
        if delay > r.maxDelay {
            delay = r.maxDelay
            break
        }
    }

    return delay
}

func (r *ItemExponentialFailureRateLimiter) Forget(item interface{}) {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.failures, item)
}

func (r *ItemExponentialFailureRateLimiter) NumRequeues(item interface{}) int {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.failures[item]
}

// RateLimitingType 速率限制队列实现
type RateLimitingType struct {
    *DelayingType
    rateLimiter RateLimiter
}

// NewRateLimitingQueue 创建速率限制队列
func NewRateLimitingQueue(rateLimiter RateLimiter) RateLimitingInterface {
    return &RateLimitingType{
        DelayingType: NewDelayingQueue().(*DelayingType),
        rateLimiter:  rateLimiter,
    }
}

func (r *RateLimitingType) AddRateLimited(item interface{}) {
    r.DelayingType.AddAfter(item, r.rateLimiter.When(item))
}

func (r *RateLimitingType) Forget(item interface{}) {
    r.rateLimiter.Forget(item)
}

func (r *RateLimitingType) NumRequeues(item interface{}) int {
    return r.rateLimiter.NumRequeues(item)
}

// 使用示例
func ExampleWorkQueue() {
    // 创建速率限制队列
    rateLimiter := NewItemExponentialFailureRateLimiter(1*time.Second, 60*time.Second)
    queue := NewRateLimitingQueue(rateLimiter)

    // 生产者
    go func() {
        for i := 0; i < 100; i++ {
            queue.Add(fmt.Sprintf("item-%d", i))
        }
    }()

    // 消费者
    go func() {
        for {
            item, shutdown := queue.Get()
            if shutdown {
                return
            }

            // 处理任务
            err := processItem(item.(string))

            if err != nil {
                // 失败，重新入队
                fmt.Printf("Processing failed: %v, requeuing\n", err)
                queue.AddRateLimited(item)
            } else {
                // 成功，忘记
                queue.Forget(item)
            }

            queue.Done(item)
        }
    }()

    time.Sleep(30 * time.Second)
    queue.ShutDown()
}

func processItem(item string) error {
    // 模拟处理
    if item == "item-5" {
        return fmt.Errorf("simulated error")
    }
    fmt.Printf("Processing: %s\n", item)
    return nil
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 可靠的任务处理 | 内存占用 |
| 支持重试 | 无持久化 |
| 速率限制 | 单点问题 |

#### 适用场景

- 控制器任务队列
- 批处理作业
- 异步任务

---

### 4.6 Channel 模式

#### 意图

Channel 模式利用 Go 语言的 channel 特性实现协程间的通信和同步。

#### 实现

```go
// Channel 模式实现
package channel

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// BufferedChannel 缓冲通道
type BufferedChannel struct {
    ch chan interface{}
}

func NewBufferedChannel(size int) *BufferedChannel {
    return &BufferedChannel{
        ch: make(chan interface{}, size),
    }
}

func (b *BufferedChannel) Send(item interface{}, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    select {
    case b.ch <- item:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (b *BufferedChannel) Receive(timeout time.Duration) (interface{}, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    select {
    case item := <-b.ch:
        return item, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// SelectPattern select 模式
type SelectPattern struct {
    ch1 chan interface{}
    ch2 chan interface{}
    quit chan bool
}

func (s *SelectPattern) Run() {
    for {
        select {
        case item := <-s.ch1:
            fmt.Printf("Received from ch1: %v\n", item)

        case item := <-s.ch2:
            fmt.Printf("Received from ch2: %v\n", item)

        case <-s.quit:
            return

        default:
            // 非阻塞
        }
    }
}

// FanOutPattern 扇出模式
func FanOutPattern(input <-chan interface{}, outputs []chan<- interface{}) {
    var wg sync.WaitGroup

    for item := range input {
        wg.Add(len(outputs))

        for _, out := range outputs {
            go func(o chan<- interface{}, val interface{}) {
                defer wg.Done()
                o <- val
            }(out, item)
        }

        wg.Wait()
    }

    for _, out := range outputs {
        close(out)
    }
}

// FanInPattern 扇入模式
func FanInPattern(inputs []<-chan interface{}) <-chan interface{} {
    output := make(chan interface{})
    var wg sync.WaitGroup

    wg.Add(len(inputs))

    for _, in := range inputs {
        go func(ch <-chan interface{}) {
            defer wg.Done()

            for item := range ch {
                output <- item
            }
        }(in)
    }

    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}

// WorkerPoolChannel 工作池模式
type WorkerPoolChannel struct {
    workers int
    jobs    chan func()
    wg      sync.WaitGroup
}

func NewWorkerPoolChannel(workers int) *WorkerPoolChannel {
    wp := &WorkerPoolChannel{
        workers: workers,
        jobs:    make(chan func()),
    }

    for i := 0; i < workers; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }

    return wp
}

func (wp *WorkerPoolChannel) worker() {
    defer wp.wg.Done()

    for job := range wp.jobs {
        job()
    }
}

func (wp *WorkerPoolChannel) Submit(job func()) {
    wp.jobs <- job
}

func (wp *WorkerPoolChannel) Close() {
    close(wp.jobs)
    wp.wg.Wait()
}

// PubSubChannel 发布订阅模式
type PubSubChannel struct {
    subscribers map[string][]chan interface{}
    mu          sync.RWMutex
}

func NewPubSubChannel() *PubSubChannel {
    return &PubSubChannel{
        subscribers: make(map[string][]chan interface{}),
    }
}

func (ps *PubSubChannel) Subscribe(topic string) <-chan interface{} {
    ps.mu.Lock()
    defer ps.mu.Unlock()

    ch := make(chan interface{}, 10)
    ps.subscribers[topic] = append(ps.subscribers[topic], ch)

    return ch
}

func (ps *PubSubChannel) Publish(topic string, msg interface{}) {
    ps.mu.RLock()
    subs := ps.subscribers[topic]
    ps.mu.RUnlock()

    for _, ch := range subs {
        select {
        case ch <- msg:
        default:
            // 通道满，跳过
        }
    }
}

// DoneChannel 完成信号模式
type DoneChannel struct {
    done chan struct{}
}

func NewDoneChannel() *DoneChannel {
    return &DoneChannel{
        done: make(chan struct{}),
    }
}

func (d *DoneChannel) Close() {
    close(d.done)
}

func (d *DoneChannel) Done() <-chan struct{} {
    return d.done
}

// ContextPattern Context 模式
func ContextPattern(ctx context.Context) {
    // 带取消的 context
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // 带超时的 context
    ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // 带截止时间的 context
    ctx, cancel = context.WithDeadline(ctx, time.Now().Add(5*time.Second))
    defer cancel()

    // 使用 context
    select {
    case <-ctx.Done():
        fmt.Println("Context done:", ctx.Err())
    case <-time.After(1 * time.Second):
        fmt.Println("Operation completed")
    }
}

// 使用示例
func ExampleChannel() {
    // 缓冲通道
    ch := NewBufferedChannel(10)

    // 发送
    go func() {
        for i := 0; i < 5; i++ {
            ch.Send(i, 1*time.Second)
        }
    }()

    // 接收
    for i := 0; i < 5; i++ {
        item, err := ch.Receive(1 * time.Second)
        if err != nil {
            fmt.Printf("Receive error: %v\n", err)
            continue
        }
        fmt.Printf("Received: %v\n", item)
    }
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 类型安全 | 需要小心死锁 |
| 简洁优雅 | 调试困难 |
| 高效 | 需要合理设计 |

#### 适用场景

- 协程通信
- 信号通知
- 数据流处理

---



## 5. 工作流设计模式

工作流设计模式用于定义、执行和监控业务流程。Kubernetes 提供了多种工作负载类型来支持不同的工作流场景。

---

### 5.1 CronJob 模式

#### 意图

CronJob 模式基于时间调度执行任务，支持类 Unix cron 表达式的定时任务。

#### 结构

```
Cron Schedule: "*/5 * * * *" (每5分钟)
                    │
                    ▼
            ┌───────────────┐
            │   CronJob     │
            │   Controller  │
            └───────┬───────┘
                    │
                    ▼
            ┌───────────────┐
            │     Job       │
            │   (Pod)       │
            └───────────────┘
```

#### 实现

```go
// CronJob 模式实现
package cronjob

import (
    "context"
    "fmt"
    "sync"
    "time"

    batchv1 "k8s.io/api/batch/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// CronSchedule Cron 调度
type CronSchedule struct {
    Minute     string
    Hour       string
    DayOfMonth string
    Month      string
    DayOfWeek  string
}

// CronJob CronJob 定义
type CronJob struct {
    Name              string
    Namespace         string
    Schedule          string
    JobTemplate       batchv1.JobTemplateSpec
    ConcurrencyPolicy ConcurrencyPolicy
    Suspend           bool
    SuccessfulJobsHistoryLimit int32
    FailedJobsHistoryLimit     int32
}

// ConcurrencyPolicy 并发策略
type ConcurrencyPolicy string

const (
    AllowConcurrent   ConcurrencyPolicy = "Allow"
    ForbidConcurrent  ConcurrencyPolicy = "Forbid"
    ReplaceConcurrent ConcurrencyPolicy = "Replace"
)

// CronJobController CronJob 控制器
type CronJobController struct {
    client     kubernetes.Interface
    cronJobs   map[string]*CronJob
    schedulers map[string]*Scheduler
    mu         sync.RWMutex
}

// Scheduler 调度器
type Scheduler struct {
    cronJob   *CronJob
    ticker    *time.Ticker
    stopCh    chan struct{}
    nextRun   time.Time
}

// NewCronJobController 创建 CronJob 控制器
func NewCronJobController(client kubernetes.Interface) *CronJobController {
    return &CronJobController{
        client:     client,
        cronJobs:   make(map[string]*CronJob),
        schedulers: make(map[string]*Scheduler),
    }
}

// CreateCronJob 创建 CronJob
func (c *CronJobController) CreateCronJob(cronJob *CronJob) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    key := fmt.Sprintf("%s/%s", cronJob.Namespace, cronJob.Name)
    c.cronJobs[key] = cronJob

    // 创建调度器
    scheduler := &Scheduler{
        cronJob: cronJob,
        stopCh:  make(chan struct{}),
    }
    c.schedulers[key] = scheduler

    // 启动调度
    go c.runScheduler(scheduler)

    return nil
}

// DeleteCronJob 删除 CronJob
func (c *CronJobController) DeleteCronJob(namespace, name string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    key := fmt.Sprintf("%s/%s", namespace, name)

    if scheduler, ok := c.schedulers[key]; ok {
        close(scheduler.stopCh)
        delete(c.schedulers, key)
    }

    delete(c.cronJobs, key)

    return nil
}

func (c *CronJobController) runScheduler(scheduler *Scheduler) {
    // 解析 cron 表达式
    schedule, err := ParseCron(scheduler.cronJob.Schedule)
    if err != nil {
        fmt.Printf("Failed to parse cron: %v\n", err)
        return
    }

    // 计算下次执行时间
    scheduler.nextRun = schedule.Next(time.Now())

    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if time.Now().After(scheduler.nextRun) && !scheduler.cronJob.Suspend {
                // 检查并发策略
                if c.canRun(scheduler.cronJob) {
                    // 创建 Job
                    if err := c.createJob(scheduler.cronJob); err != nil {
                        fmt.Printf("Failed to create job: %v\n", err)
                    }
                }

                // 更新下次执行时间
                scheduler.nextRun = schedule.Next(time.Now())
            }

        case <-scheduler.stopCh:
            return
        }
    }
}

func (c *CronJobController) canRun(cronJob *CronJob) bool {
    switch cronJob.ConcurrencyPolicy {
    case ForbidConcurrent:
        // 检查是否有正在运行的 Job
        jobs, err := c.client.BatchV1().Jobs(cronJob.Namespace).List(context.TODO(), metav1.ListOptions{
            LabelSelector: fmt.Sprintf("cronjob-name=%s", cronJob.Name),
        })
        if err != nil {
            return false
        }

        for _, job := range jobs.Items {
            if job.Status.Active > 0 {
                return false
            }
        }
        return true

    case ReplaceConcurrent:
        // 删除正在运行的 Job
        // ...
        return true

    default: // AllowConcurrent
        return true
    }
}

func (c *CronJobController) createJob(cronJob *CronJob) error {
    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            GenerateName: fmt.Sprintf("%s-", cronJob.Name),
            Namespace:    cronJob.Namespace,
            Labels: map[string]string{
                "cronjob-name": cronJob.Name,
            },
        },
        Spec: cronJob.JobTemplate.Spec,
    }

    _, err := c.client.BatchV1().Jobs(cronJob.Namespace).Create(context.TODO(), job, metav1.CreateOptions{})
    return err
}

// CronParser Cron 解析器
type CronParser struct{}

// ParseCron 解析 cron 表达式
func ParseCron(schedule string) (*Schedule, error) {
    // 简化实现，实际使用 robfig/cron 库
    return &Schedule{}, nil
}

// Schedule 调度计划
type Schedule struct {
    expr string
}

func (s *Schedule) Next(t time.Time) time.Time {
    // 简化实现
    return t.Add(5 * time.Minute)
}

// 使用示例
func ExampleCronJob() {
    cronJob := &CronJob{
        Name:      "backup-job",
        Namespace: "default",
        Schedule:  "0 2 * * *", // 每天凌晨2点
        JobTemplate: batchv1.JobTemplateSpec{
            Spec: batchv1.JobSpec{
                Template: corev1.PodTemplateSpec{
                    Spec: corev1.PodSpec{
                        Containers: []corev1.Container{
                            {
                                Name:    "backup",
                                Image:   "backup-tool:latest",
                                Command: []string{"/bin/backup.sh"},
                            },
                        },
                        RestartPolicy: corev1.RestartPolicyOnFailure,
                    },
                },
            },
        },
        ConcurrencyPolicy:          ForbidConcurrent,
        SuccessfulJobsHistoryLimit: 3,
        FailedJobsHistoryLimit:     1,
    }

    fmt.Printf("CronJob created: %s\n", cronJob.Name)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 时间精确 | 单点故障 |
| 易于管理 | 时区问题 |
| 历史记录 | 任务堆积 |

#### 适用场景

- 定时备份
- 数据清理
- 报表生成

---

### 5.2 Job 模式

#### 意图

Job 模式用于运行一次性或批处理任务，确保任务成功完成指定次数。

#### 结构

```
Job
├── Parallelism: 3
├── Completions: 5
├── BackoffLimit: 4
│
├── Pod 1 (Running)
├── Pod 2 (Succeeded)
├── Pod 3 (Failed)
└── Pod 4 (Pending)
```

#### 实现

```go
// Job 模式实现
package job

import (
    "context"
    "fmt"
    "sync"

    batchv1 "k8s.io/api/batch/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// JobController Job 控制器
type JobController struct {
    client    kubernetes.Interface
    jobs      map[string]*batchv1.Job
    mu        sync.RWMutex
}

// NewJobController 创建 Job 控制器
func NewJobController(client kubernetes.Interface) *JobController {
    return &JobController{
        client: client,
        jobs:   make(map[string]*batchv1.Job),
    }
}

// CreateJob 创建 Job
func (jc *JobController) CreateJob(job *batchv1.Job) (*batchv1.Job, error) {
    createdJob, err := jc.client.BatchV1().Jobs(job.Namespace).Create(context.TODO(), job, metav1.CreateOptions{})
    if err != nil {
        return nil, err
    }

    jc.mu.Lock()
    jc.jobs[fmt.Sprintf("%s/%s", job.Namespace, job.Name)] = createdJob
    jc.mu.Unlock()

    return createdJob, nil
}

// GetJobStatus 获取 Job 状态
func (jc *JobController) GetJobStatus(namespace, name string) (*JobStatus, error) {
    job, err := jc.client.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }

    return &JobStatus{
        Active:      job.Status.Active,
        Succeeded:   job.Status.Succeeded,
        Failed:      job.Status.Failed,
        Conditions:  job.Status.Conditions,
        StartTime:   job.Status.StartTime,
        CompletionTime: job.Status.CompletionTime,
    }, nil
}

// JobStatus Job 状态
type JobStatus struct {
    Active         int32
    Succeeded      int32
    Failed         int32
    Conditions     []batchv1.JobCondition
    StartTime      *metav1.Time
    CompletionTime *metav1.Time
}

// IsComplete 检查 Job 是否完成
func (js *JobStatus) IsComplete() bool {
    for _, cond := range js.Conditions {
        if cond.Type == batchv1.JobComplete && cond.Status == corev1.ConditionTrue {
            return true
        }
    }
    return false
}

// IsFailed 检查 Job 是否失败
func (js *JobStatus) IsFailed() bool {
    for _, cond := range js.Conditions {
        if cond.Type == batchv1.JobFailed && cond.Status == corev1.ConditionTrue {
            return true
        }
    }
    return false
}

// WaitForCompletion 等待 Job 完成
func (jc *JobController) WaitForCompletion(ctx context.Context, namespace, name string, timeout int) error {
    for i := 0; i < timeout; i++ {
        status, err := jc.GetJobStatus(namespace, name)
        if err != nil {
            return err
        }

        if status.IsComplete() {
            return nil
        }

        if status.IsFailed() {
            return fmt.Errorf("job failed")
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(1 * time.Second):
        }
    }

    return fmt.Errorf("timeout waiting for job completion")
}

// DeleteJob 删除 Job
func (jc *JobController) DeleteJob(namespace, name string) error {
    err := jc.client.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
    if err != nil {
        return err
    }

    jc.mu.Lock()
    delete(jc.jobs, fmt.Sprintf("%s/%s", namespace, name))
    jc.mu.Unlock()

    return nil
}

// JobBuilder Job 构建器
type JobBuilder struct {
    job *batchv1.Job
}

// NewJobBuilder 创建 Job 构建器
func NewJobBuilder(name, namespace string) *JobBuilder {
    return &JobBuilder{
        job: &batchv1.Job{
            ObjectMeta: metav1.ObjectMeta{
                Name:      name,
                Namespace: namespace,
            },
            Spec: batchv1.JobSpec{
                Template: corev1.PodTemplateSpec{
                    Spec: corev1.PodSpec{
                        RestartPolicy: corev1.RestartPolicyNever,
                    },
                },
            },
        },
    }
}

func (jb *JobBuilder) WithParallelism(n int32) *JobBuilder {
    jb.job.Spec.Parallelism = &n
    return jb
}

func (jb *JobBuilder) WithCompletions(n int32) *JobBuilder {
    jb.job.Spec.Completions = &n
    return jb
}

func (jb *JobBuilder) WithBackoffLimit(n int32) *JobBuilder {
    jb.job.Spec.BackoffLimit = &n
    return jb
}

func (jb *JobBuilder) WithTTLSecondsAfterFinished(n int32) *JobBuilder {
    jb.job.Spec.TTLSecondsAfterFinished = &n
    return jb
}

func (jb *JobBuilder) WithContainer(container corev1.Container) *JobBuilder {
    jb.job.Spec.Template.Spec.Containers = append(
        jb.job.Spec.Template.Spec.Containers,
        container,
    )
    return jb
}

func (jb *JobBuilder) Build() *batchv1.Job {
    return jb.job
}

// IndexedJob 索引 Job（并行处理）
type IndexedJob struct {
    *JobController
}

func NewIndexedJob(client kubernetes.Interface) *IndexedJob {
    return &IndexedJob{
        JobController: NewJobController(client),
    }
}

func (ij *IndexedJob) CreateIndexedJob(name, namespace string, completions int32, workFunc func(int)) (*batchv1.Job, error) {
    completionMode := batchv1.IndexedCompletion

    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: namespace,
        },
        Spec: batchv1.JobSpec{
            Completions:    &completions,
            Parallelism:    &completions,
            CompletionMode: &completionMode,
            Template: corev1.PodTemplateSpec{
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "worker",
                            Image: "worker:latest",
                            Env: []corev1.EnvVar{
                                {
                                    Name: "JOB_COMPLETION_INDEX",
                                    ValueFrom: &corev1.EnvVarSource{
                                        FieldRef: &corev1.ObjectFieldSelector{
                                            FieldPath: "metadata.annotations['batch.kubernetes.io/job-completion-index']",
                                        },
                                    },
                                },
                            },
                        },
                    },
                    RestartPolicy: corev1.RestartPolicyNever,
                },
            },
        },
    }

    return ij.CreateJob(job)
}

// 使用示例
func ExampleJob() {
    // 创建简单 Job
    job := NewJobBuilder("data-processing", "default").
        WithCompletions(1).
        WithBackoffLimit(3).
        WithTTLSecondsAfterFinished(3600).
        WithContainer(corev1.Container{
            Name:    "processor",
            Image:   "data-processor:latest",
            Command: []string{"/bin/process-data.sh"},
        }).
        Build()

    fmt.Printf("Job created: %s\n", job.Name)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 确保完成 | 资源占用 |
| 自动重试 | 无持久化 |
| 并行处理 | 调试困难 |

#### 适用场景

- 数据处理
- 批量任务
- 一次性任务

---

### 5.3 DaemonSet 模式

#### 意图

DaemonSet 模式确保每个节点（或符合条件的节点）上运行一个 Pod 副本，常用于系统级服务。

#### 结构

```
┌─────────────────────────────────────────┐
│           DaemonSet Controller          │
│                                         │
│  Node 1 ──► Pod 1 (node-exporter)      │
│  Node 2 ──► Pod 2 (node-exporter)      │
│  Node 3 ──► Pod 3 (node-exporter)      │
│  Node 4 ──► Pod 4 (node-exporter)      │
└─────────────────────────────────────────┘
```

#### 实现

```go
// DaemonSet 模式实现
package daemonset

import (
    "context"
    "fmt"

    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// DaemonSetController DaemonSet 控制器
type DaemonSetController struct {
    client kubernetes.Interface
}

// NewDaemonSetController 创建 DaemonSet 控制器
func NewDaemonSetController(client kubernetes.Interface) *DaemonSetController {
    return &DaemonSetController{client: client}
}

// CreateDaemonSet 创建 DaemonSet
func (dc *DaemonSetController) CreateDaemonSet(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
    return dc.client.AppsV1().DaemonSets(ds.Namespace).Create(context.TODO(), ds, metav1.CreateOptions{})
}

// GetDaemonSetStatus 获取 DaemonSet 状态
func (dc *DaemonSetController) GetDaemonSetStatus(namespace, name string) (*DaemonSetStatus, error) {
    ds, err := dc.client.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }

    return &DaemonSetStatus{
        CurrentNumberScheduled: ds.Status.CurrentNumberScheduled,
        DesiredNumberScheduled: ds.Status.DesiredNumberScheduled,
        NumberMisscheduled:     ds.Status.NumberMisscheduled,
        NumberReady:            ds.Status.NumberReady,
        UpdatedNumberScheduled: ds.Status.UpdatedNumberScheduled,
        NumberAvailable:        ds.Status.NumberAvailable,
        NumberUnavailable:      ds.Status.NumberUnavailable,
    }, nil
}

// DaemonSetStatus DaemonSet 状态
type DaemonSetStatus struct {
    CurrentNumberScheduled int32
    DesiredNumberScheduled int32
    NumberMisscheduled     int32
    NumberReady            int32
    UpdatedNumberScheduled int32
    NumberAvailable        int32
    NumberUnavailable      int32
}

// IsReady 检查 DaemonSet 是否就绪
func (ds *DaemonSetStatus) IsReady() bool {
    return ds.NumberReady == ds.DesiredNumberScheduled
}

// DaemonSetBuilder DaemonSet 构建器
type DaemonSetBuilder struct {
    ds *appsv1.DaemonSet
}

// NewDaemonSetBuilder 创建 DaemonSet 构建器
func NewDaemonSetBuilder(name, namespace string) *DaemonSetBuilder {
    return &DaemonSetBuilder{
        ds: &appsv1.DaemonSet{
            ObjectMeta: metav1.ObjectMeta{
                Name:      name,
                Namespace: namespace,
            },
            Spec: appsv1.DaemonSetSpec{
                Selector: &metav1.LabelSelector{},
                Template: corev1.PodTemplateSpec{},
            },
        },
    }
}

func (db *DaemonSetBuilder) WithLabels(labels map[string]string) *DaemonSetBuilder {
    db.ds.Labels = labels
    return db
}

func (db *DaemonSetBuilder) WithSelector(matchLabels map[string]string) *DaemonSetBuilder {
    db.ds.Spec.Selector.MatchLabels = matchLabels
    db.ds.Spec.Template.Labels = matchLabels
    return db
}

func (db *DaemonSetBuilder) WithPodTemplate(template corev1.PodTemplateSpec) *DaemonSetBuilder {
    db.ds.Spec.Template = template
    return db
}

func (db *DaemonSetBuilder) WithUpdateStrategy(strategy appsv1.DaemonSetUpdateStrategy) *DaemonSetBuilder {
    db.ds.Spec.UpdateStrategy = strategy
    return db
}

func (db *DaemonSetBuilder) WithNodeSelector(nodeSelector map[string]string) *DaemonSetBuilder {
    db.ds.Spec.Template.Spec.NodeSelector = nodeSelector
    return db
}

func (db *DaemonSetBuilder) WithTolerations(tolerations []corev1.Toleration) *DaemonSetBuilder {
    db.ds.Spec.Template.Spec.Tolerations = tolerations
    return db
}

func (db *DaemonSetBuilder) Build() *appsv1.DaemonSet {
    return db.ds
}

// NodeExporterDaemonSet 节点监控 DaemonSet
func NewNodeExporterDaemonSet() *appsv1.DaemonSet {
    return NewDaemonSetBuilder("node-exporter", "monitoring").
        WithSelector(map[string]string{
            "app": "node-exporter",
        }).
        WithPodTemplate(corev1.PodTemplateSpec{
            ObjectMeta: metav1.ObjectMeta{
                Labels: map[string]string{
                    "app": "node-exporter",
                },
            },
            Spec: corev1.PodSpec{
                HostNetwork: true,
                HostPID:     true,
                Containers: []corev1.Container{
                    {
                        Name:  "node-exporter",
                        Image: "prom/node-exporter:v1.3.1",
                        Args: []string{
                            "--path.procfs=/host/proc",
                            "--path.sysfs=/host/sys",
                            "--path.rootfs=/host/root",
                        },
                        Ports: []corev1.ContainerPort{
                            {
                                Name:          "metrics",
                                ContainerPort: 9100,
                                Protocol:      corev1.ProtocolTCP,
                            },
                        },
                        VolumeMounts: []corev1.VolumeMount{
                            {
                                Name:      "proc",
                                MountPath: "/host/proc",
                                ReadOnly:  true,
                            },
                            {
                                Name:      "sys",
                                MountPath: "/host/sys",
                                ReadOnly:  true,
                            },
                            {
                                Name:      "root",
                                MountPath: "/host/root",
                                ReadOnly:  true,
                            },
                        },
                    },
                },
                Volumes: []corev1.Volume{
                    {
                        Name: "proc",
                        VolumeSource: corev1.VolumeSource{
                            HostPath: &corev1.HostPathVolumeSource{
                                Path: "/proc",
                            },
                        },
                    },
                    {
                        Name: "sys",
                        VolumeSource: corev1.VolumeSource{
                            HostPath: &corev1.HostPathVolumeSource{
                                Path: "/sys",
                            },
                        },
                    },
                    {
                        Name: "root",
                        VolumeSource: corev1.VolumeSource{
                            HostPath: &corev1.HostPathVolumeSource{
                                Path: "/",
                            },
                        },
                    },
                },
            },
        }).
        WithTolerations([]corev1.Toleration{
            {
                Operator: corev1.TolerationOpExists,
            },
        }).
        Build()
}

// FluentBitDaemonSet 日志收集 DaemonSet
func NewFluentBitDaemonSet() *appsv1.DaemonSet {
    return NewDaemonSetBuilder("fluent-bit", "logging").
        WithSelector(map[string]string{
            "app": "fluent-bit",
        }).
        WithPodTemplate(corev1.PodTemplateSpec{
            ObjectMeta: metav1.ObjectMeta{
                Labels: map[string]string{
                    "app": "fluent-bit",
                },
            },
            Spec: corev1.PodSpec{
                Containers: []corev1.Container{
                    {
                        Name:  "fluent-bit",
                        Image: "fluent/fluent-bit:1.9",
                        VolumeMounts: []corev1.VolumeMount{
                            {
                                Name:      "varlog",
                                MountPath: "/var/log",
                            },
                            {
                                Name:      "varlibdockercontainers",
                                MountPath: "/var/lib/docker/containers",
                                ReadOnly:  true,
                            },
                            {
                                Name:      "fluent-bit-config",
                                MountPath: "/fluent-bit/etc/",
                            },
                        },
                    },
                },
                Volumes: []corev1.Volume{
                    {
                        Name: "varlog",
                        VolumeSource: corev1.VolumeSource{
                            HostPath: &corev1.HostPathVolumeSource{
                                Path: "/var/log",
                            },
                        },
                    },
                    {
                        Name: "varlibdockercontainers",
                        VolumeSource: corev1.VolumeSource{
                            HostPath: &corev1.HostPathVolumeSource{
                                Path: "/var/lib/docker/containers",
                            },
                        },
                    },
                    {
                        Name: "fluent-bit-config",
                        VolumeSource: corev1.VolumeSource{
                            ConfigMap: &corev1.ConfigMapVolumeSource{
                                LocalObjectReference: corev1.LocalObjectReference{
                                    Name: "fluent-bit-config",
                                },
                            },
                        },
                    },
                },
            },
        }).
        WithTolerations([]corev1.Toleration{
            {
                Key:      "node-role.kubernetes.io/master",
                Operator: corev1.TolerationOpExists,
                Effect:   corev1.TaintEffectNoSchedule,
            },
        }).
        Build()
}

// 使用示例
func ExampleDaemonSet() {
    // 创建节点监控 DaemonSet
    nodeExporter := NewNodeExporterDaemonSet()
    fmt.Printf("DaemonSet created: %s\n", nodeExporter.Name)

    // 创建日志收集 DaemonSet
    fluentBit := NewFluentBitDaemonSet()
    fmt.Printf("DaemonSet created: %s\n", fluentBit.Name)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 节点全覆盖 | 资源消耗 |
| 自动部署 | 升级复杂 |
| 系统服务友好 | 调度限制 |

#### 适用场景

- 日志收集
- 监控代理
- 网络代理

---

### 5.4 Pipeline 模式

#### 意图

Pipeline 模式将工作流分解为多个阶段，每个阶段处理特定的任务，数据在阶段之间流动。

#### 结构

```
Source ──► Stage 1 ──► Stage 2 ──► Stage 3 ──► Sink
            (Build)     (Test)      (Deploy)
```

#### 实现

```go
// Pipeline 模式实现
package pipeline

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Stage 阶段接口
type Stage interface {
    Name() string
    Execute(ctx context.Context, input interface{}) (interface{}, error)
    Rollback(ctx context.Context, input interface{}) error
}

// Pipeline CI/CD 流水线
type Pipeline struct {
    name      string
    stages    []Stage
    status    PipelineStatus
    mu        sync.RWMutex
}

// PipelineStatus 流水线状态
type PipelineStatus struct {
    State       PipelineState
    CurrentStage int
    StartTime   time.Time
    EndTime     *time.Time
    Error       error
}

// PipelineState 流水线状态
type PipelineState string

const (
    PipelinePending   PipelineState = "Pending"
    PipelineRunning   PipelineState = "Running"
    PipelineSucceeded PipelineState = "Succeeded"
    PipelineFailed    PipelineState = "Failed"
    PipelineCancelled PipelineState = "Cancelled"
)

// NewPipeline 创建流水线
func NewPipeline(name string) *Pipeline {
    return &Pipeline{
        name:   name,
        stages: make([]Stage, 0),
        status: PipelineStatus{
            State: PipelinePending,
        },
    }
}

// AddStage 添加阶段
func (p *Pipeline) AddStage(stage Stage) {
    p.stages = append(p.stages, stage)
}

// Run 运行流水线
func (p *Pipeline) Run(ctx context.Context, input interface{}) error {
    p.mu.Lock()
    p.status.State = PipelineRunning
    p.status.StartTime = time.Now()
    p.mu.Unlock()

    var result interface{} = input

    for i, stage := range p.stages {
        p.mu.Lock()
        p.status.CurrentStage = i
        p.mu.Unlock()

        select {
        case <-ctx.Done():
            p.setStatus(PipelineCancelled, ctx.Err())
            return ctx.Err()
        default:
        }

        fmt.Printf("Executing stage: %s\n", stage.Name())

        var err error
        result, err = stage.Execute(ctx, result)
        if err != nil {
            p.setStatus(PipelineFailed, err)

            // 回滚已执行的步骤
            p.rollback(ctx, i, input)

            return err
        }
    }

    p.setStatus(PipelineSucceeded, nil)
    return nil
}

func (p *Pipeline) setStatus(state PipelineState, err error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.status.State = state
    p.status.Error = err
    now := time.Now()
    p.status.EndTime = &now
}

func (p *Pipeline) rollback(ctx context.Context, failedIndex int, input interface{}) {
    fmt.Println("Rolling back pipeline...")

    for i := failedIndex - 1; i >= 0; i-- {
        stage := p.stages[i]
        fmt.Printf("Rolling back stage: %s\n", stage.Name())

        if err := stage.Rollback(ctx, input); err != nil {
            fmt.Printf("Rollback failed for stage %s: %v\n", stage.Name(), err)
        }
    }
}

// GetStatus 获取状态
func (p *Pipeline) GetStatus() PipelineStatus {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.status
}

// BuildStage 构建阶段
type BuildStage struct {
    image    string
    commands []string
}

func NewBuildStage(image string, commands []string) *BuildStage {
    return &BuildStage{
        image:    image,
        commands: commands,
    }
}

func (b *BuildStage) Name() string {
    return "build"
}

func (b *BuildStage) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    fmt.Printf("Building with image: %s\n", b.image)

    for _, cmd := range b.commands {
        fmt.Printf("Executing: %s\n", cmd)
        // 实际执行构建命令
    }

    return map[string]string{
        "image": b.image,
        "tag":   "latest",
    }, nil
}

func (b *BuildStage) Rollback(ctx context.Context, input interface{}) error {
    fmt.Println("Cleaning up build artifacts")
    return nil
}

// TestStage 测试阶段
type TestStage struct {
    testCmd string
}

func NewTestStage(testCmd string) *TestStage {
    return &TestStage{testCmd: testCmd}
}

func (t *TestStage) Name() string {
    return "test"
}

func (t *TestStage) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    fmt.Printf("Running tests: %s\n", t.testCmd)

    // 执行测试
    // 模拟测试
    time.Sleep(1 * time.Second)

    return input, nil
}

func (t *TestStage) Rollback(ctx context.Context, input interface{}) error {
    return nil
}

// DeployStage 部署阶段
type DeployStage struct {
    namespace   string
    deployment  string
    image       string
}

func NewDeployStage(namespace, deployment, image string) *DeployStage {
    return &DeployStage{
        namespace:  namespace,
        deployment: deployment,
        image:      image,
    }
}

func (d *DeployStage) Name() string {
    return "deploy"
}

func (d *DeployStage) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    buildResult := input.(map[string]string)
    image := buildResult["image"]
    tag := buildResult["tag"]

    fmt.Printf("Deploying %s/%s:%s to namespace %s\n", d.image, d.deployment, tag, d.namespace)

    // 实际部署逻辑

    return map[string]string{
        "deployment": d.deployment,
        "namespace":  d.namespace,
        "image":      fmt.Sprintf("%s:%s", image, tag),
    }, nil
}

func (d *DeployStage) Rollback(ctx context.Context, input interface{}) error {
    fmt.Printf("Rolling back deployment %s in namespace %s\n", d.deployment, d.namespace)
    return nil
}

// ParallelPipeline 并行流水线
type ParallelPipeline struct {
    name   string
    groups [][]Stage
}

func NewParallelPipeline(name string) *ParallelPipeline {
    return &ParallelPipeline{
        name:   name,
        groups: make([][]Stage, 0),
    }
}

func (pp *ParallelPipeline) AddParallelGroup(stages ...Stage) {
    pp.groups = append(pp.groups, stages)
}

func (pp *ParallelPipeline) Run(ctx context.Context, input interface{}) error {
    var currentInput interface{} = input

    for _, group := range pp.groups {
        results := make([]interface{}, len(group))
        errChan := make(chan error, len(group))

        var wg sync.WaitGroup
        wg.Add(len(group))

        for i, stage := range group {
            go func(idx int, s Stage) {
                defer wg.Done()

                result, err := s.Execute(ctx, currentInput)
                if err != nil {
                    errChan <- err
                    return
                }

                results[idx] = result
            }(i, stage)
        }

        wg.Wait()
        close(errChan)

        for err := range errChan {
            return err
        }

        // 合并结果
        currentInput = results
    }

    return nil
}

// 使用示例
func ExamplePipeline() {
    // 创建流水线
    pipeline := NewPipeline("ci-cd")

    // 添加阶段
    pipeline.AddStage(NewBuildStage("golang:1.19", []string{
        "go build -o app",
        "docker build -t myapp:latest .",
    }))

    pipeline.AddStage(NewTestStage("go test ./..."))

    pipeline.AddStage(NewDeployStage("production", "myapp", "myapp"))

    // 运行流水线
    ctx := context.Background()
    if err := pipeline.Run(ctx, nil); err != nil {
        fmt.Printf("Pipeline failed: %v\n", err)
    }

    // 获取状态
    status := pipeline.GetStatus()
    fmt.Printf("Pipeline status: %s\n", status.State)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 模块化 | 复杂度高 |
| 可复用 | 调试困难 |
| 可回滚 | 状态管理 |

#### 适用场景

- CI/CD 流水线
- 数据处理
- 工作流编排

---

### 5.5 DAG 模式

#### 意图

DAG（有向无环图）模式用于定义复杂的依赖关系，任务按照依赖顺序执行。

#### 结构

```
     A
    / \
   B   C
   |   |
   D   E
    \ /
     F
```

#### 实现

```go
// DAG 模式实现
package dag

import (
    "context"
    "fmt"
    "sync"
)

// Task DAG 任务
type Task struct {
    ID           string
    Name         string
    Dependencies []string
    Execute      func(ctx context.Context) error
    Status       TaskStatus
    mu           sync.RWMutex
}

// TaskStatus 任务状态
type TaskStatus int

const (
    TaskPending TaskStatus = iota
    TaskRunning
    TaskSucceeded
    TaskFailed
    TaskSkipped
)

func (t *Task) GetStatus() TaskStatus {
    t.mu.RLock()
    defer t.mu.RUnlock()
    return t.Status
}

func (t *Task) SetStatus(status TaskStatus) {
    t.mu.Lock()
    defer t.mu.Unlock()
    t.Status = status
}

// DAG 有向无环图
type DAG struct {
    tasks map[string]*Task
    edges map[string][]string // task -> dependents
    mu    sync.RWMutex
}

// NewDAG 创建 DAG
func NewDAG() *DAG {
    return &DAG{
        tasks: make(map[string]*Task),
        edges: make(map[string][]string),
    }
}

// AddTask 添加任务
func (d *DAG) AddTask(task *Task) error {
    d.mu.Lock()
    defer d.mu.Unlock()

    // 检查循环依赖
    if err := d.detectCycle(task); err != nil {
        return err
    }

    d.tasks[task.ID] = task

    // 构建依赖图
    for _, dep := range task.Dependencies {
        d.edges[dep] = append(d.edges[dep], task.ID)
    }

    return nil
}

func (d *DAG) detectCycle(task *Task) error {
    visited := make(map[string]bool)
    recStack := make(map[string]bool)

    var dfs func(string) bool
    dfs = func(id string) bool {
        visited[id] = true
        recStack[id] = true

        task := d.tasks[id]
        if task != nil {
            for _, dep := range task.Dependencies {
                if !visited[dep] {
                    if dfs(dep) {
                        return true
                    }
                } else if recStack[dep] {
                    return true
                }
            }
        }

        recStack[id] = false
        return false
    }

    // 检查新任务的依赖
    for _, dep := range task.Dependencies {
        if dep == task.ID {
            return fmt.Errorf("self-dependency detected")
        }
        if dfs(dep) {
            return fmt.Errorf("cycle detected")
        }
    }

    return nil
}

// GetTasks 获取所有任务
func (d *DAG) GetTasks() []*Task {
    d.mu.RLock()
    defer d.mu.RUnlock()

    tasks := make([]*Task, 0, len(d.tasks))
    for _, task := range d.tasks {
        tasks = append(tasks, task)
    }
    return tasks
}

// GetReadyTasks 获取就绪任务（依赖已完成）
func (d *DAG) GetReadyTasks() []*Task {
    d.mu.RLock()
    defer d.mu.RUnlock()

    var ready []*Task

    for _, task := range d.tasks {
        if task.GetStatus() != TaskPending {
            continue
        }

        // 检查所有依赖是否完成
        allDepsCompleted := true
        for _, dep := range task.Dependencies {
            depTask := d.tasks[dep]
            if depTask == nil || depTask.GetStatus() != TaskSucceeded {
                allDepsCompleted = false
                break
            }
        }

        if allDepsCompleted {
            ready = append(ready, task)
        }
    }

    return ready
}

// IsComplete 检查 DAG 是否完成
func (d *DAG) IsComplete() bool {
    d.mu.RLock()
    defer d.mu.RUnlock()

    for _, task := range d.tasks {
        status := task.GetStatus()
        if status != TaskSucceeded && status != TaskSkipped {
            return false
        }
    }

    return true
}

// DAGExecutor DAG 执行器
type DAGExecutor struct {
    dag       *DAG
    workers   int
    mu        sync.Mutex
}

// NewDAGExecutor 创建 DAG 执行器
func NewDAGExecutor(dag *DAG, workers int) *DAGExecutor {
    if workers <= 0 {
        workers = 4
    }

    return &DAGExecutor{
        dag:     dag,
        workers: workers,
    }
}

// Execute 执行 DAG
func (de *DAGExecutor) Execute(ctx context.Context) error {
    taskChan := make(chan *Task, de.workers)
    errChan := make(chan error, 1)
    var wg sync.WaitGroup

    // 启动工作协程
    for i := 0; i < de.workers; i++ {
        wg.Add(1)
        go de.worker(ctx, taskChan, errChan, &wg)
    }

    // 调度任务
    go de.schedule(ctx, taskChan, errChan)

    // 等待完成
    wg.Wait()
    close(errChan)

    // 检查错误
    for err := range errChan {
        if err != nil {
            return err
        }
    }

    return nil
}

func (de *DAGExecutor) schedule(ctx context.Context, taskChan chan<- *Task, errChan chan<- error) {
    pendingTasks := make(map[string]bool)

    for {
        select {
        case <-ctx.Done():
            return
        default:
        }

        // 获取就绪任务
        readyTasks := de.dag.GetReadyTasks()

        for _, task := range readyTasks {
            if !pendingTasks[task.ID] {
                pendingTasks[task.ID] = true
                task.SetStatus(TaskRunning)
                taskChan <- task
            }
        }

        // 检查是否完成
        if de.dag.IsComplete() {
            close(taskChan)
            return
        }

        // 检查是否有失败的任务
        for _, task := range de.dag.GetTasks() {
            if task.GetStatus() == TaskFailed {
                errChan <- fmt.Errorf("task %s failed", task.ID)
                close(taskChan)
                return
            }
        }
    }
}

func (de *DAGExecutor) worker(ctx context.Context, taskChan <-chan *Task, errChan chan<- error, wg *sync.WaitGroup) {
    defer wg.Done()

    for task := range taskChan {
        fmt.Printf("Executing task: %s\n", task.Name)

        err := task.Execute(ctx)

        if err != nil {
            task.SetStatus(TaskFailed)
            errChan <- fmt.Errorf("task %s failed: %w", task.ID, err)
        } else {
            task.SetStatus(TaskSucceeded)
            fmt.Printf("Task completed: %s\n", task.Name)
        }
    }
}

// ArgoWorkflow Argo 工作流风格
type ArgoWorkflow struct {
    dag       *DAG
    templates map[string]*WorkflowTemplate
}

type WorkflowTemplate struct {
    Name    string
    Inputs  map[string]string
    Outputs map[string]string
    Script  string
}

func NewArgoWorkflow() *ArgoWorkflow {
    return &ArgoWorkflow{
        dag:       NewDAG(),
        templates: make(map[string]*WorkflowTemplate),
    }
}

func (aw *ArgoWorkflow) AddTemplate(template *WorkflowTemplate) {
    aw.templates[template.Name] = template
}

func (aw *ArgoWorkflow) AddTask(task *Task) error {
    return aw.dag.AddTask(task)
}

func (aw *ArgoWorkflow) Submit(ctx context.Context) error {
    executor := NewDAGExecutor(aw.dag, 4)
    return executor.Execute(ctx)
}

// 使用示例
func ExampleDAG() {
    // 创建 DAG
    dag := NewDAG()

    // 添加任务
    dag.AddTask(&Task{
        ID:           "A",
        Name:         "Prepare Data",
        Dependencies: []string{},
        Execute: func(ctx context.Context) error {
            fmt.Println("Task A: Preparing data")
            return nil
        },
    })

    dag.AddTask(&Task{
        ID:           "B",
        Name:         "Process Data 1",
        Dependencies: []string{"A"},
        Execute: func(ctx context.Context) error {
            fmt.Println("Task B: Processing data 1")
            return nil
        },
    })

    dag.AddTask(&Task{
        ID:           "C",
        Name:         "Process Data 2",
        Dependencies: []string{"A"},
        Execute: func(ctx context.Context) error {
            fmt.Println("Task C: Processing data 2")
            return nil
        },
    })

    dag.AddTask(&Task{
        ID:           "D",
        Name:         "Merge Results",
        Dependencies: []string{"B", "C"},
        Execute: func(ctx context.Context) error {
            fmt.Println("Task D: Merging results")
            return nil
        },
    })

    // 执行 DAG
    executor := NewDAGExecutor(dag, 4)
    ctx := context.Background()

    if err := executor.Execute(ctx); err != nil {
        fmt.Printf("DAG execution failed: %v\n", err)
    } else {
        fmt.Println("DAG execution completed successfully")
    }
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 表达复杂依赖 | 实现复杂 |
| 并行执行 | 调试困难 |
| 可视化友好 | 死锁风险 |

#### 适用场景

- 复杂工作流
- 数据处理
- 机器学习流水线

---

### 5.6 State Machine 模式

#### 意图

State Machine 模式通过定义状态和转换规则来管理工作流的生命周期。

#### 结构

```
     ┌─────────┐
     │  Start  │
     └────┬────┘
          │
          ▼
    ┌───────────┐
    │  Pending  │◄────────┐
    └─────┬─────┘         │
          │               │
     ┌────┴────┐          │
     ▼         ▼          │
┌────────┐ ┌────────┐     │
│Running │ │ Failed │─────┘
└───┬────┘ └────────┘
    │
    ▼
┌──────────┐
│Completed │
└──────────┘
```

#### 实现

```go
// State Machine 模式实现
package statemachine

import (
    "context"
    "fmt"
    "sync"
)

// State 状态
type State string

// Event 事件
type Event string

// Transition 转换
type Transition struct {
    From  State
    Event Event
    To    State
    Action func(ctx context.Context, data interface{}) error
}

// StateMachine 状态机
type StateMachine struct {
    current     State
    transitions map[State]map[Event]*Transition
    handlers    map[State]func(ctx context.Context, data interface{}) error
    mu          sync.RWMutex
}

// NewStateMachine 创建状态机
func NewStateMachine(initial State) *StateMachine {
    return &StateMachine{
        current:     initial,
        transitions: make(map[State]map[Event]*Transition),
        handlers:    make(map[State]func(ctx context.Context, data interface{}) error),
    }
}

// AddTransition 添加转换
func (sm *StateMachine) AddTransition(t *Transition) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if sm.transitions[t.From] == nil {
        sm.transitions[t.From] = make(map[Event]*Transition)
    }

    sm.transitions[t.From][t.Event] = t
}

// AddStateHandler 添加状态处理器
func (sm *StateMachine) AddStateHandler(state State, handler func(ctx context.Context, data interface{}) error) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    sm.handlers[state] = handler
}

// Trigger 触发事件
func (sm *StateMachine) Trigger(ctx context.Context, event Event, data interface{}) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    transitions := sm.transitions[sm.current]
    if transitions == nil {
        return fmt.Errorf("no transitions from state %s", sm.current)
    }

    transition := transitions[event]
    if transition == nil {
        return fmt.Errorf("no transition for event %s from state %s", event, sm.current)
    }

    // 执行转换动作
    if transition.Action != nil {
        if err := transition.Action(ctx, data); err != nil {
            return err
        }
    }

    // 转换状态
    sm.current = transition.To

    // 执行状态处理器
    if handler := sm.handlers[sm.current]; handler != nil {
        return handler(ctx, data)
    }

    return nil
}

// GetCurrentState 获取当前状态
func (sm *StateMachine) GetCurrentState() State {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return sm.current
}

// CanTrigger 检查是否可以触发事件
func (sm *StateMachine) CanTrigger(event Event) bool {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    transitions := sm.transitions[sm.current]
    if transitions == nil {
        return false
    }

    return transitions[event] != nil
}

// WorkflowStateMachine 工作流状态机
type WorkflowStateMachine struct {
    *StateMachine
    workflowID string
}

// 工作流状态
const (
    StateStart      State = "Start"
    StatePending    State = "Pending"
    StateRunning    State = "Running"
    StatePaused     State = "Paused"
    StateSucceeded  State = "Succeeded"
    StateFailed     State = "Failed"
    StateCancelled  State = "Cancelled"
)

// 工作流事件
const (
    EventSubmit    Event = "Submit"
    EventStart     Event = "Start"
    EventPause     Event = "Pause"
    EventResume    Event = "Resume"
    EventComplete  Event = "Complete"
    EventFail      Event = "Fail"
    EventCancel    Event = "Cancel"
    EventRetry     Event = "Retry"
)

// NewWorkflowStateMachine 创建工作流状态机
func NewWorkflowStateMachine(workflowID string) *WorkflowStateMachine {
    sm := NewStateMachine(StateStart)

    // 添加转换
    sm.AddTransition(&Transition{
        From:  StateStart,
        Event: EventSubmit,
        To:    StatePending,
    })

    sm.AddTransition(&Transition{
        From:  StatePending,
        Event: EventStart,
        To:    StateRunning,
    })

    sm.AddTransition(&Transition{
        From:  StateRunning,
        Event: EventPause,
        To:    StatePaused,
    })

    sm.AddTransition(&Transition{
        From:  StatePaused,
        Event: EventResume,
        To:    StateRunning,
    })

    sm.AddTransition(&Transition{
        From:  StateRunning,
        Event: EventComplete,
        To:    StateSucceeded,
    })

    sm.AddTransition(&Transition{
        From:  StateRunning,
        Event: EventFail,
        To:    StateFailed,
    })

    sm.AddTransition(&Transition{
        From:  StateFailed,
        Event: EventRetry,
        To:    StatePending,
    })

    sm.AddTransition(&Transition{
        From:  StatePending,
        Event: EventCancel,
        To:    StateCancelled,
    })

    sm.AddTransition(&Transition{
        From:  StateRunning,
        Event: EventCancel,
        To:    StateCancelled,
    })

    sm.AddTransition(&Transition{
        From:  StatePaused,
        Event: EventCancel,
        To:    StateCancelled,
    })

    // 添加状态处理器
    sm.AddStateHandler(StateRunning, func(ctx context.Context, data interface{}) error {
        fmt.Printf("Workflow %s is now running\n", workflowID)
        return nil
    })

    sm.AddStateHandler(StateSucceeded, func(ctx context.Context, data interface{}) error {
        fmt.Printf("Workflow %s completed successfully\n", workflowID)
        return nil
    })

    sm.AddStateHandler(StateFailed, func(ctx context.Context, data interface{}) error {
        fmt.Printf("Workflow %s failed\n", workflowID)
        return nil
    })

    return &WorkflowStateMachine{
        StateMachine: sm,
        workflowID:   workflowID,
    }
}

// StateMachineWorkflow 基于状态机的工作流
type StateMachineWorkflow struct {
    sm        *WorkflowStateMachine
    steps     map[State][]WorkflowStep
    currentStep int
}

type WorkflowStep struct {
    Name   string
    Action func(ctx context.Context) error
}

func NewStateMachineWorkflow(id string) *StateMachineWorkflow {
    return &StateMachineWorkflow{
        sm:    NewWorkflowStateMachine(id),
        steps: make(map[State][]WorkflowStep),
    }
}

func (w *StateMachineWorkflow) AddStep(state State, step WorkflowStep) {
    w.steps[state] = append(w.steps[state], step)
}

func (w *StateMachineWorkflow) Execute(ctx context.Context) error {
    // 提交工作流
    if err := w.sm.Trigger(ctx, EventSubmit, nil); err != nil {
        return err
    }

    // 启动工作流
    if err := w.sm.Trigger(ctx, EventStart, nil); err != nil {
        return err
    }

    // 执行步骤
    for w.sm.GetCurrentState() == StateRunning {
        steps := w.steps[StateRunning]

        for _, step := range steps {
            fmt.Printf("Executing step: %s\n", step.Name)

            if err := step.Action(ctx); err != nil {
                w.sm.Trigger(ctx, EventFail, err)
                return err
            }
        }

        // 完成
        w.sm.Trigger(ctx, EventComplete, nil)
    }

    return nil
}

// 使用示例
func ExampleStateMachine() {
    // 创建工作流状态机
    workflow := NewStateMachineWorkflow("workflow-1")

    // 添加步骤
    workflow.AddStep(StateRunning, WorkflowStep{
        Name: "Step 1: Initialize",
        Action: func(ctx context.Context) error {
            fmt.Println("Initializing...")
            return nil
        },
    })

    workflow.AddStep(StateRunning, WorkflowStep{
        Name: "Step 2: Process",
        Action: func(ctx context.Context) error {
            fmt.Println("Processing...")
            return nil
        },
    })

    workflow.AddStep(StateRunning, WorkflowStep{
        Name: "Step 3: Finalize",
        Action: func(ctx context.Context) error {
            fmt.Println("Finalizing...")
            return nil
        },
    })

    // 执行工作流
    ctx := context.Background()
    if err := workflow.Execute(ctx); err != nil {
        fmt.Printf("Workflow failed: %v\n", err)
    }

    fmt.Printf("Final state: %s\n", workflow.sm.GetCurrentState())
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 清晰的状态管理 | 状态爆炸 |
| 易于验证 | 复杂转换 |
| 可视化友好 | 实现复杂 |

#### 适用场景

- 订单处理
- 审批流程
- 任务生命周期管理

---



## 6. K8s 特有模式

Kubernetes 引入了一些特有的设计模式，这些模式是云原生应用开发的核心。

---

### 6.1 Controller 模式

#### 意图

Controller 模式通过持续监控资源状态，将实际状态调整为期望状态，实现声明式管理。

#### 结构

```
┌─────────────────────────────────────────────────────────┐
│                    Controller Loop                        │
│                                                           │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐           │
│  │  Watch   │───►│  Diff    │───►│  Reconcile│           │
│  │  Changes │    │  State   │    │  State    │           │
│  └──────────┘    └──────────┘    └────┬─────┘           │
│                                       │                   │
│                                       ▼                   │
│                                ┌──────────┐              │
│                                │  Update  │              │
│                                │  Status  │              │
│                                └──────────┘              │
└─────────────────────────────────────────────────────────┘
```

#### 实现

```go
// Controller 模式实现
package controller

import (
    "context"
    "fmt"
    "sync"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/util/workqueue"
)

// Controller 控制器接口
type Controller interface {
    Run(ctx context.Context, workers int) error
    Name() string
}

// Reconciler 协调器接口
type Reconciler interface {
    Reconcile(ctx context.Context, key string) error
}

// BaseController 基础控制器
type BaseController struct {
    name         string
    informer     cache.SharedIndexInformer
    queue        workqueue.RateLimitingInterface
    reconciler   Reconciler
    syncPeriod   time.Duration
    mu           sync.Mutex
}

// NewBaseController 创建基础控制器
func NewBaseController(name string, informer cache.SharedIndexInformer, reconciler Reconciler) *BaseController {
    return &BaseController{
        name:       name,
        informer:   informer,
        queue:      workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
        reconciler: reconciler,
        syncPeriod: 10 * time.Minute,
    }
}

// Run 启动控制器
func (c *BaseController) Run(ctx context.Context, workers int) error {
    defer c.queue.ShutDown()

    fmt.Printf("Starting controller: %s\n", c.name)

    // 添加事件处理器
    c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: c.enqueue,
        UpdateFunc: func(old, new interface{}) {
            c.enqueue(new)
        },
        DeleteFunc: c.enqueue,
    })

    // 启动 informer
    go c.informer.Run(ctx.Done())

    // 等待缓存同步
    if !cache.WaitForCacheSync(ctx.Done(), c.informer.HasSynced) {
        return fmt.Errorf("failed to sync cache")
    }

    fmt.Printf("Controller %s cache synced\n", c.name)

    // 启动工作协程
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.worker(ctx)
        }()
    }

    // 等待停止信号
    <-ctx.Done()

    // 等待工作协程完成
    wg.Wait()

    fmt.Printf("Controller %s stopped\n", c.name)
    return nil
}

func (c *BaseController) enqueue(obj interface{}) {
    key, err := cache.MetaNamespaceKeyFunc(obj)
    if err != nil {
        fmt.Printf("Failed to get key: %v\n", err)
        return
    }

    c.queue.Add(key)
}

func (c *BaseController) worker(ctx context.Context) {
    for c.processNextItem(ctx) {
    }
}

func (c *BaseController) processNextItem(ctx context.Context) bool {
    key, quit := c.queue.Get()
    if quit {
        return false
    }

    defer c.queue.Done(key)

    err := c.reconciler.Reconcile(ctx, key.(string))

    if err != nil {
        // 重试
        c.queue.AddRateLimited(key)
        fmt.Printf("Reconcile error for %s: %v\n", key, err)
    } else {
        c.queue.Forget(key)
    }

    return true
}

func (c *BaseController) Name() string {
    return c.name
}

// GenericReconciler 通用协调器
type GenericReconciler struct {
    getter    cache.Getter
    updater   cache.Updater
    syncFunc  func(ctx context.Context, obj interface{}) error
}

func NewGenericReconciler(getter cache.Getter, updater cache.Updater, syncFunc func(ctx context.Context, obj interface{}) error) *GenericReconciler {
    return &GenericReconciler{
        getter:   getter,
        updater:  updater,
        syncFunc: syncFunc,
    }
}

func (gr *GenericReconciler) Reconcile(ctx context.Context, key string) error {
    // 从缓存获取对象
    obj, exists, err := gr.getter.GetByKey(key)
    if err != nil {
        return err
    }

    if !exists {
        // 对象已删除
        return nil
    }

    // 执行同步
    return gr.syncFunc(ctx, obj)
}

// DeploymentReconciler Deployment 协调器
type DeploymentReconciler struct {
    client kubernetes.Interface
}

func NewDeploymentReconciler(client kubernetes.Interface) *DeploymentReconciler {
    return &DeploymentReconciler{client: client}
}

func (dr *DeploymentReconciler) Reconcile(ctx context.Context, key string) error {
    namespace, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        return err
    }

    // 获取 Deployment
    deployment, err := dr.client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
    if err != nil {
        return err
    }

    // 获取关联的 ReplicaSet
    rsList, err := dr.getReplicaSetsForDeployment(deployment)
    if err != nil {
        return err
    }

    // 检查是否需要滚动更新
    if dr.isRollingUpdateNeeded(deployment, rsList) {
        if err := dr.performRollingUpdate(ctx, deployment, rsList); err != nil {
            return err
        }
    }

    // 同步 ReplicaSet
    if err := dr.syncReplicaSets(ctx, deployment, rsList); err != nil {
        return err
    }

    // 更新状态
    return dr.updateDeploymentStatus(ctx, deployment, rsList)
}

func (dr *DeploymentReconciler) getReplicaSetsForDeployment(deployment *appsv1.Deployment) ([]*appsv1.ReplicaSet, error) {
    // 实现获取 ReplicaSet 的逻辑
    return nil, nil
}

func (dr *DeploymentReconciler) isRollingUpdateNeeded(deployment *appsv1.Deployment, rsList []*appsv1.ReplicaSet) bool {
    // 检查是否需要滚动更新
    return false
}

func (dr *DeploymentReconciler) performRollingUpdate(ctx context.Context, deployment *appsv1.Deployment, rsList []*appsv1.ReplicaSet) error {
    // 执行滚动更新
    return nil
}

func (dr *DeploymentReconciler) syncReplicaSets(ctx context.Context, deployment *appsv1.Deployment, rsList []*appsv1.ReplicaSet) error {
    // 同步 ReplicaSet
    return nil
}

func (dr *DeploymentReconciler) updateDeploymentStatus(ctx context.Context, deployment *appsv1.Deployment, rsList []*appsv1.ReplicaSet) error {
    // 更新 Deployment 状态
    return nil
}

// 使用示例
func ExampleController() {
    // 创建 Kubernetes 客户端
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err)
    }

    client, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    // 创建 Informer
    informer := informers.NewSharedInformerFactory(client, 0).Apps().V1().Deployments().Informer()

    // 创建协调器
    reconciler := NewDeploymentReconciler(client)

    // 创建控制器
    controller := NewBaseController("deployment-controller", informer, reconciler)

    // 启动控制器
    ctx := context.Background()
    if err := controller.Run(ctx, 2); err != nil {
        panic(err)
    }
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 声明式管理 | 学习曲线陡峭 |
| 自动恢复 | 调试复杂 |
| 幂等性 | 状态延迟 |

#### 适用场景

- 资源管理
- 自动扩缩容
- 配置同步

---

### 6.2 Operator 模式

#### 意图

Operator 模式将运维知识编码到软件中，自动化复杂的应用生命周期管理。

#### 结构

```
┌─────────────────────────────────────────────────────────┐
│                      Operator                           │
│                                                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │   Custom    │───►│  Controller │───►│   Manager   │ │
│  │  Resource   │    │   (Watch)   │    │  (Reconcile)│ │
│  └─────────────┘    └─────────────┘    └──────┬──────┘ │
│                                                │        │
│  ┌─────────────┐    ┌─────────────┐    ┌──────┴──────┐ │
│  │   Deploy    │◄───│   Upgrade   │◄───│   Backup    │ │
│  └─────────────┘    └─────────────┘    └─────────────┘ │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### 实现

```go
// Operator 模式实现
package operator

import (
    "context"
    "fmt"
    "os"
    "runtime"

    "k8s.io/apimachinery/pkg/runtime"
    utilruntime "k8s.io/apimachinery/pkg/util/runtime"
    clientgoscheme "k8s.io/client-go/kubernetes/scheme"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/healthz"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"
    "sigs.k8s.io/controller-runtime/pkg/manager"
)

// Operator 操作器
type Operator struct {
    manager    manager.Manager
    name       string
    version    string
    scheme     *runtime.Scheme
}

// OperatorConfig 操作器配置
type OperatorConfig struct {
    Name                 string
    Version              string
    MetricsAddr          string
    ProbeAddr            string
    LeaderElect          bool
    LeaderElectionID     string
    LeaderElectionNamespace string
    Namespace            string
}

// NewOperator 创建操作器
func NewOperator(config *OperatorConfig) (*Operator, error) {
    // 创建 Scheme
    scheme := runtime.NewScheme()
    utilruntime.Must(clientgoscheme.AddToScheme(scheme))

    // 创建 Manager
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:                  scheme,
        MetricsBindAddress:      config.MetricsAddr,
        Port:                    9443,
        HealthProbeBindAddress:  config.ProbeAddr,
        LeaderElection:          config.LeaderElect,
        LeaderElectionID:        config.LeaderElectionID,
        LeaderElectionNamespace: config.LeaderElectionNamespace,
        Namespace:               config.Namespace,
    })
    if err != nil {
        return nil, fmt.Errorf("unable to create manager: %w", err)
    }

    return &Operator{
        manager: mgr,
        name:    config.Name,
        version: config.Version,
        scheme:  scheme,
    }, nil
}

// AddController 添加控制器
func (o *Operator) AddController(name string, reconciler interface{}) error {
    // 使用 controller-runtime 添加控制器
    return nil
}

// Start 启动操作器
func (o *Operator) Start(ctx context.Context) error {
    // 添加健康检查
    if err := o.manager.AddHealthzCheck("healthz", healthz.Ping); err != nil {
        return fmt.Errorf("unable to set up health check: %w", err)
    }

    if err := o.manager.AddReadyzCheck("readyz", healthz.Ping); err != nil {
        return fmt.Errorf("unable to set up ready check: %w", err)
    }

    fmt.Printf("Starting operator: %s v%s\n", o.name, o.version)

    // 启动 Manager
    return o.manager.Start(ctx)
}

// GetClient 获取客户端
func (o *Operator) GetClient() client.Client {
    return o.manager.GetClient()
}

// GetScheme 获取 Scheme
func (o *Operator) GetScheme() *runtime.Scheme {
    return o.scheme
}

// CustomResource 自定义资源
type CustomResource struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec              CustomResourceSpec   `json:"spec,omitempty"`
    Status            CustomResourceStatus `json:"status,omitempty"`
}

// CustomResourceSpec 自定义资源规格
type CustomResourceSpec struct {
    Replicas int32  `json:"replicas"`
    Image    string `json:"image"`
    Config   map[string]string `json:"config,omitempty"`
}

// CustomResourceStatus 自定义资源状态
type CustomResourceStatus struct {
    Phase      string `json:"phase,omitempty"`
    Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// CustomResourceReconciler 自定义资源协调器
type CustomResourceReconciler struct {
    client.Client
    Scheme *runtime.Scheme
    Log    logr.Logger
}

// Reconcile 协调
func (r *CustomResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := r.Log.WithValues("customresource", req.NamespacedName)

    // 获取自定义资源
    cr := &CustomResource{}
    if err := r.Get(ctx, req.NamespacedName, cr); err != nil {
        if errors.IsNotFound(err) {
            return ctrl.Result{}, nil
        }
        return ctrl.Result{}, err
    }

    // 执行协调逻辑
    if err := r.reconcileDeployment(ctx, cr); err != nil {
        return ctrl.Result{}, err
    }

    if err := r.reconcileService(ctx, cr); err != nil {
        return ctrl.Result{}, err
    }

    if err := r.reconcileConfigMap(ctx, cr); err != nil {
        return ctrl.Result{}, err
    }

    // 更新状态
    cr.Status.Phase = "Running"
    if err := r.Status().Update(ctx, cr); err != nil {
        return ctrl.Result{}, err
    }

    log.Info("Reconciliation complete")

    return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *CustomResourceReconciler) reconcileDeployment(ctx context.Context, cr *CustomResource) error {
    // 创建或更新 Deployment
    deployment := &appsv1.Deployment{}
    deployment.Name = cr.Name
    deployment.Namespace = cr.Namespace

    _, err := ctrl.CreateOrUpdate(ctx, r.Client, deployment, func() error {
        deployment.Spec.Replicas = &cr.Spec.Replicas
        deployment.Spec.Selector = &metav1.LabelSelector{
            MatchLabels: map[string]string{
                "app": cr.Name,
            },
        }
        deployment.Spec.Template.Labels = map[string]string{
            "app": cr.Name,
        }
        deployment.Spec.Template.Spec.Containers = []corev1.Container{
            {
                Name:  "app",
                Image: cr.Spec.Image,
            },
        }
        return ctrl.SetControllerReference(cr, deployment, r.Scheme)
    })

    return err
}

func (r *CustomResourceReconciler) reconcileService(ctx context.Context, cr *CustomResource) error {
    // 创建或更新 Service
    return nil
}

func (r *CustomResourceReconciler) reconcileConfigMap(ctx context.Context, cr *CustomResource) error {
    // 创建或更新 ConfigMap
    return nil
}

// SetupWithManager 设置 Manager
func (r *CustomResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&CustomResource{}).
        Owns(&appsv1.Deployment{}).
        Owns(&corev1.Service{}).
        Owns(&corev1.ConfigMap{}).
        Complete(r)
}

// LifecycleManager 生命周期管理器
type LifecycleManager struct {
    client client.Client
}

func NewLifecycleManager(client client.Client) *LifecycleManager {
    return &LifecycleManager{client: client}
}

// Install 安装应用
func (lm *LifecycleManager) Install(ctx context.Context, cr *CustomResource) error {
    // 执行安装逻辑
    return nil
}

// Upgrade 升级应用
func (lm *LifecycleManager) Upgrade(ctx context.Context, cr *CustomResource) error {
    // 执行升级逻辑
    return nil
}

// Backup 备份应用
func (lm *LifecycleManager) Backup(ctx context.Context, cr *CustomResource) error {
    // 执行备份逻辑
    return nil
}

// Restore 恢复应用
func (lm *LifecycleManager) Restore(ctx context.Context, cr *CustomResource, backupPath string) error {
    // 执行恢复逻辑
    return nil
}

// Uninstall 卸载应用
func (lm *LifecycleManager) Uninstall(ctx context.Context, cr *CustomResource) error {
    // 执行卸载逻辑
    return nil
}

// 使用示例
func ExampleOperator() {
    config := &OperatorConfig{
        Name:                 "myapp-operator",
        Version:              "1.0.0",
        MetricsAddr:          ":8080",
        ProbeAddr:            ":8081",
        LeaderElect:          true,
        LeaderElectionID:     "myapp-operator-lock",
        LeaderElectionNamespace: "default",
    }

    operator, err := NewOperator(config)
    if err != nil {
        panic(err)
    }

    // 设置协调器
    reconciler := &CustomResourceReconciler{
        Client: operator.GetClient(),
        Scheme: operator.GetScheme(),
        Log:    ctrl.Log.WithName("controllers").WithName("CustomResource"),
    }

    if err := reconciler.SetupWithManager(operator.manager); err != nil {
        panic(err)
    }

    // 启动操作器
    ctx := ctrl.SetupSignalHandler()
    if err := operator.Start(ctx); err != nil {
        panic(err)
    }
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 自动化运维 | 开发复杂 |
| 领域知识编码 | 学习曲线陡峭 |
| 可复用 | 调试困难 |

#### 适用场景

- 数据库管理
- 中间件运维
- 复杂应用生命周期

---

### 6.3 Initializer 模式

#### 意图

Initializer 模式在资源创建时执行初始化逻辑，如设置默认值、注入配置等。

#### 结构

```
API Request ──► Admission Controller ──► Initializer ──► etcd
                     │                       │
                     │ Mutating              │ Async Init
                     │ Webhook               │
                     ▼                       ▼
               Set Defaults          Inject Sidecar
```

#### 实现

```go
// Initializer 模式实现
package initializer

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    admissionv1 "k8s.io/api/admission/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/utils/pointer"
)

// Initializer 初始化器接口
type Initializer interface {
    Initialize(obj runtime.Object) error
    Name() string
}

// PodInitializer Pod 初始化器
type PodInitializer struct {
    name string
}

func NewPodInitializer(name string) *PodInitializer {
    return &PodInitializer{name: name}
}

func (pi *PodInitializer) Name() string {
    return pi.name
}

func (pi *PodInitializer) Initialize(obj runtime.Object) error {
    pod, ok := obj.(*corev1.Pod)
    if !ok {
        return fmt.Errorf("not a pod object")
    }

    // 设置默认值
    if err := pi.setDefaults(pod); err != nil {
        return err
    }

    // 注入初始化容器
    if err := pi.injectInitContainers(pod); err != nil {
        return err
    }

    // 注入 Sidecar
    if err := pi.injectSidecars(pod); err != nil {
        return err
    }

    return nil
}

func (pi *PodInitializer) setDefaults(pod *corev1.Pod) error {
    // 设置默认资源限制
    for i := range pod.Spec.Containers {
        container := &pod.Spec.Containers[i]

        if container.Resources.Limits == nil {
            container.Resources.Limits = corev1.ResourceList{}
        }

        if _, ok := container.Resources.Limits[corev1.ResourceCPU]; !ok {
            container.Resources.Limits[corev1.ResourceCPU] = resource.MustParse("500m")
        }

        if _, ok := container.Resources.Limits[corev1.ResourceMemory]; !ok {
            container.Resources.Limits[corev1.ResourceMemory] = resource.MustParse("512Mi")
        }
    }

    // 设置默认安全上下文
    if pod.Spec.SecurityContext == nil {
        pod.Spec.SecurityContext = &corev1.PodSecurityContext{
            RunAsNonRoot: pointer.Bool(true),
            SeccompProfile: &corev1.SeccompProfile{
                Type: corev1.SeccompProfileTypeRuntimeDefault,
            },
        }
    }

    return nil
}

func (pi *PodInitializer) injectInitContainers(pod *corev1.Pod) error {
    // 检查是否需要初始化
    if pod.Annotations["initializer.example.com/skip"] == "true" {
        return nil
    }

    // 注入初始化容器
    initContainer := corev1.Container{
        Name:  "initializer",
        Image: "initializer:latest",
        Args: []string{
            "--namespace", pod.Namespace,
            "--pod-name", pod.Name,
        },
    }

    pod.Spec.InitContainers = append([]corev1.Container{initContainer}, pod.Spec.InitContainers...)

    return nil
}

func (pi *PodInitializer) injectSidecars(pod *corev1.Pod) error {
    // 根据注解注入 Sidecar
    if pod.Annotations["sidecar.example.com/inject"] != "true" {
        return nil
    }

    sidecar := corev1.Container{
        Name:  "proxy",
        Image: "proxy:latest",
    }

    pod.Spec.Containers = append(pod.Spec.Containers, sidecar)

    return nil
}

// MutatingWebhook 变更 Webhook
type MutatingWebhook struct {
    initializers []Initializer
}

func NewMutatingWebhook(initializers ...Initializer) *MutatingWebhook {
    return &MutatingWebhook{initializers: initializers}
}

func (mw *MutatingWebhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var admissionReview admissionv1.AdmissionReview

    if err := json.NewDecoder(r.Body).Decode(&admissionReview); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    admissionResponse := mw.mutate(r.Context(), admissionReview.Request)

    admissionReview.Response = admissionResponse
    admissionReview.Response.UID = admissionReview.Request.UID

    resp, err := json.Marshal(admissionReview)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(resp)
}

func (mw *MutatingWebhook) mutate(ctx context.Context, req *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
    // 解析对象
    obj, err := mw.parseObject(req)
    if err != nil {
        return &admissionv1.AdmissionResponse{
            Allowed: false,
            Result: &metav1.Status{
                Message: err.Error(),
            },
        }
    }

    // 执行初始化
    for _, initializer := range mw.initializers {
        if err := initializer.Initialize(obj); err != nil {
            return &admissionv1.AdmissionResponse{
                Allowed: false,
                Result: &metav1.Status{
                    Message: err.Error(),
                },
            }
        }
    }

    // 生成补丁
    patch, err := mw.generatePatch(req.Object.Raw, obj)
    if err != nil {
        return &admissionv1.AdmissionResponse{
            Allowed: false,
            Result: &metav1.Status{
                Message: err.Error(),
            },
        }
    }

    patchType := admissionv1.PatchTypeJSONPatch

    return &admissionv1.AdmissionResponse{
        Allowed:   true,
        Patch:     patch,
        PatchType: &patchType,
    }
}

func (mw *MutatingWebhook) parseObject(req *admissionv1.AdmissionRequest) (runtime.Object, error) {
    // 根据 Kind 解析对象
    switch req.Kind.Kind {
    case "Pod":
        pod := &corev1.Pod{}
        if err := json.Unmarshal(req.Object.Raw, pod); err != nil {
            return nil, err
        }
        return pod, nil
    default:
        return nil, fmt.Errorf("unsupported kind: %s", req.Kind.Kind)
    }
}

func (mw *MutatingWebhook) generatePatch(original []byte, modified runtime.Object) ([]byte, error) {
    modifiedBytes, err := json.Marshal(modified)
    if err != nil {
        return nil, err
    }

    // 使用 JSON Patch 生成差异
    // 实际实现需要使用 jsonpatch 库
    return modifiedBytes, nil
}

// InitializerManager 初始化器管理器
type InitializerManager struct {
    initializers map[string]Initializer
    mu           sync.RWMutex
}

func NewInitializerManager() *InitializerManager {
    return &InitializerManager{
        initializers: make(map[string]Initializer),
    }
}

func (im *InitializerManager) Register(initializer Initializer) {
    im.mu.Lock()
    defer im.mu.Unlock()

    im.initializers[initializer.Name()] = initializer
}

func (im *InitializerManager) Get(name string) (Initializer, bool) {
    im.mu.RLock()
    defer im.mu.RUnlock()

    init, ok := im.initializers[name]
    return init, ok
}

// 使用示例
func ExampleInitializer() {
    // 创建初始化器
    podInitializer := NewPodInitializer("pod-initializer")

    // 创建 Webhook
    webhook := NewMutatingWebhook(podInitializer)

    // 启动 Webhook 服务器
    http.Handle("/mutate", webhook)
    http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 统一初始化 | 增加延迟 |
| 可扩展 | 调试困难 |
| 解耦 | 依赖 Webhook |

#### 适用场景

- 默认值设置
- Sidecar 注入
- 配置验证

---

### 6.4 Finalizer 模式

#### 意图

Finalizer 模式在资源删除前执行清理逻辑，确保资源被正确清理。

#### 结构

```
Delete Request ──► Check Finalizers ──► Execute Cleanup ──► Remove Finalizer ──► Delete
                        │                      │
                        │ 有 Finalizer          │ 清理失败
                        ▼                      ▼
                  阻止删除              重新入队
```

#### 实现

```go
// Finalizer 模式实现
package finalizer

import (
    "context"
    "fmt"
    "slices"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Finalizer 清理器接口
type Finalizer interface {
    Name() string
    Cleanup(ctx context.Context, obj client.Object) error
}

// FinalizerManager 清理器管理器
type FinalizerManager struct {
    client     client.Client
    finalizers map[string]Finalizer
}

// NewFinalizerManager 创建清理器管理器
func NewFinalizerManager(client client.Client) *FinalizerManager {
    return &FinalizerManager{
        client:     client,
        finalizers: make(map[string]Finalizer),
    }
}

// Register 注册清理器
func (fm *FinalizerManager) Register(finalizer Finalizer) {
    fm.finalizers[finalizer.Name()] = finalizer
}

// AddFinalizer 添加 Finalizer
func (fm *FinalizerManager) AddFinalizer(ctx context.Context, obj client.Object, finalizerName string) error {
    // 检查 Finalizer 是否已存在
    if slices.Contains(obj.GetFinalizers(), finalizerName) {
        return nil
    }

    // 添加 Finalizer
    controllerutil.AddFinalizer(obj, finalizerName)

    // 更新对象
    return fm.client.Update(ctx, obj)
}

// RemoveFinalizer 移除 Finalizer
func (fm *FinalizerManager) RemoveFinalizer(ctx context.Context, obj client.Object, finalizerName string) error {
    // 检查 Finalizer 是否存在
    if !slices.Contains(obj.GetFinalizers(), finalizerName) {
        return nil
    }

    // 移除 Finalizer
    controllerutil.RemoveFinalizer(obj, finalizerName)

    // 更新对象
    return fm.client.Update(ctx, obj)
}

// HandleDeletion 处理删除
func (fm *FinalizerManager) HandleDeletion(ctx context.Context, obj client.Object) (bool, error) {
    // 检查是否正在删除
    if obj.GetDeletionTimestamp() == nil {
        // 添加 Finalizer
        for name := range fm.finalizers {
            if err := fm.AddFinalizer(ctx, obj, name); err != nil {
                return false, err
            }
        }
        return false, nil
    }

    // 执行清理
    for _, finalizerName := range obj.GetFinalizers() {
        finalizer, ok := fm.finalizers[finalizerName]
        if !ok {
            // 未知的 Finalizer，直接移除
            if err := fm.RemoveFinalizer(ctx, obj, finalizerName); err != nil {
                return false, err
            }
            continue
        }

        // 执行清理
        if err := finalizer.Cleanup(ctx, obj); err != nil {
            return false, fmt.Errorf("cleanup failed for finalizer %s: %w", finalizerName, err)
        }

        // 移除 Finalizer
        if err := fm.RemoveFinalizer(ctx, obj, finalizerName); err != nil {
            return false, err
        }
    }

    // 所有 Finalizer 已移除，允许删除
    return len(obj.GetFinalizers()) == 0, nil
}

// ResourceCleanupFinalizer 资源清理 Finalizer
type ResourceCleanupFinalizer struct {
    client client.Client
}

func NewResourceCleanupFinalizer(client client.Client) *ResourceCleanupFinalizer {
    return &ResourceCleanupFinalizer{client: client}
}

func (rcf *ResourceCleanupFinalizer) Name() string {
    return "cleanup.example.com/resources"
}

func (rcf *ResourceCleanupFinalizer) Cleanup(ctx context.Context, obj client.Object) error {
    // 清理关联资源
    // 例如：删除创建的 Deployment、Service 等

    fmt.Printf("Cleaning up resources for %s/%s\n", obj.GetNamespace(), obj.GetName())

    return nil
}

// ExternalResourceCleanupFinalizer 外部资源清理 Finalizer
type ExternalResourceCleanupFinalizer struct {
    // 外部服务客户端
}

func NewExternalResourceCleanupFinalizer() *ExternalResourceCleanupFinalizer {
    return &ExternalResourceCleanupFinalizer{}
}

func (ercf *ExternalResourceCleanupFinalizer) Name() string {
    return "cleanup.example.com/external"
}

func (ercf *ExternalResourceCleanupFinalizer) Cleanup(ctx context.Context, obj client.Object) error {
    // 清理外部资源
    // 例如：删除云资源、数据库记录等

    fmt.Printf("Cleaning up external resources for %s/%s\n", obj.GetNamespace(), obj.GetName())

    return nil
}

// BackupFinalizer 备份 Finalizer
type BackupFinalizer struct {
    backupService BackupService
}

type BackupService interface {
    Backup(ctx context.Context, obj client.Object) error
}

func NewBackupFinalizer(service BackupService) *BackupFinalizer {
    return &BackupFinalizer{backupService: service}
}

func (bf *BackupFinalizer) Name() string {
    return "cleanup.example.com/backup"
}

func (bf *BackupFinalizer) Cleanup(ctx context.Context, obj client.Object) error {
    // 在删除前备份数据
    fmt.Printf("Backing up data for %s/%s\n", obj.GetNamespace(), obj.GetName())

    return bf.backupService.Backup(ctx, obj)
}

// 使用示例
func ExampleFinalizer() {
    // 创建客户端
    // client := ...

    // 创建 Finalizer 管理器
    manager := NewFinalizerManager(client)

    // 注册清理器
    manager.Register(NewResourceCleanupFinalizer(client))
    manager.Register(NewExternalResourceCleanupFinalizer())

    // 在协调器中处理删除
    // func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    //     obj := &MyResource{}
    //     if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
    //         return ctrl.Result{}, err
    //     }
    //
    //     canDelete, err := r.finalizerManager.HandleDeletion(ctx, obj)
    //     if err != nil {
    //         return ctrl.Result{}, err
    //     }
    //
    //     if !canDelete {
    //         return ctrl.Result{Requeue: true}, nil
    //     }
    //
    //     // 继续正常协调
    // }
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 确保清理 | 删除延迟 |
| 可重试 | 死锁风险 |
| 可扩展 | 复杂性 |

#### 适用场景

- 资源清理
- 数据备份
- 外部服务解绑

---

### 6.5 Owner Reference 模式

#### 意图

Owner Reference 模式建立资源之间的父子关系，实现级联删除和垃圾回收。

#### 结构

```
Deployment (Owner)
    │
    ├── ReplicaSet (Owned)
    │       │
    │       └── Pod (Owned)
    │
    └── ReplicaSet (Owned - old)
            │
            └── Pod (Owned - to be deleted)
```

#### 实现

```go
// Owner Reference 模式实现
package ownerref

import (
    "context"
    "fmt"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/utils/pointer"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// OwnerReferenceManager Owner Reference 管理器
type OwnerReferenceManager struct {
    client client.Client
}

// NewOwnerReferenceManager 创建 Owner Reference 管理器
func NewOwnerReferenceManager(client client.Client) *OwnerReferenceManager {
    return &OwnerReferenceManager{client: client}
}

// SetOwnerReference 设置 Owner Reference
func (orm *OwnerReferenceManager) SetOwnerReference(owner, owned client.Object) error {
    return controllerutil.SetControllerReference(owner, owned, orm.client.Scheme())
}

// SetOwnerReferenceWithOptions 设置 Owner Reference（带选项）
func (orm *OwnerReferenceManager) SetOwnerReferenceWithOptions(owner, owned client.Object, blockOwnerDeletion, controller bool) error {
    // 获取 Owner 的 APIVersion 和 Kind
    gvk, err := orm.getGVK(owner)
    if err != nil {
        return err
    }

    ownerRef := metav1.OwnerReference{
        APIVersion:         gvk.GroupVersion().String(),
        Kind:               gvk.Kind,
        Name:               owner.GetName(),
        UID:                owner.GetUID(),
        Controller:         pointer.Bool(controller),
        BlockOwnerDeletion: pointer.Bool(blockOwnerDeletion),
    }

    // 设置 Owner Reference
    owned.SetOwnerReferences(append(owned.GetOwnerReferences(), ownerRef))

    return nil
}

func (orm *OwnerReferenceManager) getGVK(obj client.Object) (schema.GroupVersionKind, error) {
    return orm.client.Scheme().ObjectKinds(obj)
}

// GetOwnedResources 获取拥有的资源
func (orm *OwnerReferenceManager) GetOwnedResources(ctx context.Context, owner client.Object, ownedList client.ObjectList) error {
    // 列出所有资源
    if err := orm.client.List(ctx, ownedList, client.InNamespace(owner.GetNamespace())); err != nil {
        return err
    }

    // 过滤出属于 Owner 的资源
    // 实际实现需要遍历列表并检查 Owner Reference

    return nil
}

// DeleteOwnedResources 删除拥有的资源
func (orm *OwnerReferenceManager) DeleteOwnedResources(ctx context.Context, owner client.Object, ownedList client.ObjectList) error {
    // 获取所有拥有的资源
    if err := orm.GetOwnedResources(ctx, owner, ownedList); err != nil {
        return err
    }

    // 删除资源
    // 实际实现需要遍历列表并删除

    return nil
}

// GarbageCollector 垃圾回收器
type GarbageCollector struct {
    client client.Client
}

// NewGarbageCollector 创建垃圾回收器
func NewGarbageCollector(client client.Client) *GarbageCollector {
    return &GarbageCollector{client: client}
}

// CollectOrphans 收集孤儿资源
func (gc *GarbageCollector) CollectOrphans(ctx context.Context, ownedList client.ObjectList) error {
    // 列出所有资源
    if err := gc.client.List(ctx, ownedList); err != nil {
        return err
    }

    // 检查每个资源的 Owner 是否存在
    // 如果 Owner 不存在，则删除该资源

    return nil
}

// ResourceManager 资源管理器
type ResourceManager struct {
    client client.Client
    owner  client.Object
}

// NewResourceManager 创建资源管理器
func NewResourceManager(client client.Client, owner client.Object) *ResourceManager {
    return &ResourceManager{
        client: client,
        owner:  owner,
    }
}

// CreateOwnedResource 创建拥有的资源
func (rm *ResourceManager) CreateOwnedResource(ctx context.Context, obj client.Object) error {
    // 设置 Owner Reference
    if err := controllerutil.SetControllerReference(rm.owner, obj, rm.client.Scheme()); err != nil {
        return err
    }

    // 创建资源
    return rm.client.Create(ctx, obj)
}

// UpdateOwnedResource 更新拥有的资源
func (rm *ResourceManager) UpdateOwnedResource(ctx context.Context, obj client.Object) error {
    // 确保 Owner Reference 存在
    if err := controllerutil.SetControllerReference(rm.owner, obj, rm.client.Scheme()); err != nil {
        return err
    }

    // 更新资源
    return rm.client.Update(ctx, obj)
}

// DeleteOwnedResource 删除拥有的资源
func (rm *ResourceManager) DeleteOwnedResource(ctx context.Context, obj client.Object) error {
    return rm.client.Delete(ctx, obj)
}

// DeploymentResourceManager Deployment 资源管理器
type DeploymentResourceManager struct {
    *ResourceManager
    deployment *appsv1.Deployment
}

func NewDeploymentResourceManager(client client.Client, deployment *appsv1.Deployment) *DeploymentResourceManager {
    return &DeploymentResourceManager{
        ResourceManager: NewResourceManager(client, deployment),
        deployment:      deployment,
    }
}

func (drm *DeploymentResourceManager) CreateService(ctx context.Context, ports []corev1.ServicePort) (*corev1.Service, error) {
    service := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name:      drm.deployment.Name,
            Namespace: drm.deployment.Namespace,
            Labels:    drm.deployment.Labels,
        },
        Spec: corev1.ServiceSpec{
            Selector: drm.deployment.Spec.Selector.MatchLabels,
            Ports:    ports,
        },
    }

    if err := drm.CreateOwnedResource(ctx, service); err != nil {
        return nil, err
    }

    return service, nil
}

func (drm *DeploymentResourceManager) CreateConfigMap(ctx context.Context, data map[string]string) (*corev1.ConfigMap, error) {
    configMap := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      fmt.Sprintf("%s-config", drm.deployment.Name),
            Namespace: drm.deployment.Namespace,
        },
        Data: data,
    }

    if err := drm.CreateOwnedResource(ctx, configMap); err != nil {
        return nil, err
    }

    return configMap, nil
}

// 使用示例
func ExampleOwnerReference() {
    // 创建客户端
    // client := ...

    // 创建 Deployment
    deployment := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "myapp",
            Namespace: "default",
        },
    }

    // 创建资源管理器
    manager := NewDeploymentResourceManager(client, deployment)

    ctx := context.Background()

    // 创建 Service（自动设置 Owner Reference）
    service, err := manager.CreateService(ctx, []corev1.ServicePort{
        {
            Name: "http",
            Port: 80,
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Service created with Owner Reference: %s/%s\n", service.Namespace, service.Name)

    // 当 Deployment 被删除时，Service 会被自动删除
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| 自动清理 | 级联删除风险 |
| 关系清晰 | 循环依赖 |
| 垃圾回收 | 调试困难 |

#### 适用场景

- 资源生命周期管理
- 级联删除
- 依赖关系追踪

---



## 7. Go 代码实现示例

本章提供完整的 Go 代码示例，展示如何在实际项目中应用上述设计模式。

---

### 7.1 完整 Controller 协调循环实现

```go
// controller_example.go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/runtime"
    "k8s.io/apimachinery/pkg/util/wait"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/leaderelection"
    "k8s.io/client-go/tools/leaderelection/resourcelock"
    "k8s.io/client-go/util/workqueue"
    "k8s.io/klog/v2"
)

// PodController Pod 控制器
type PodController struct {
    clientset  kubernetes.Interface
    informer   cache.SharedIndexInformer
    queue      workqueue.RateLimitingInterface
    workerCount int
}

// NewPodController 创建 Pod 控制器
func NewPodController(clientset kubernetes.Interface) *PodController {
    // 创建 Informer 工厂
    factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

    // 获取 Pod Informer
    podInformer := factory.Core().V1().Pods().Informer()

    // 创建工作队列
    queue := workqueue.NewRateLimitingQueue(workqueue.NewMaxOfRateLimiter(
        workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
        &workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
    ))

    controller := &PodController{
        clientset:   clientset,
        informer:    podInformer,
        queue:       queue,
        workerCount: 3,
    }

    // 添加事件处理器
    podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: controller.enqueuePod,
        UpdateFunc: func(old, new interface{}) {
            controller.enqueuePod(new)
        },
        DeleteFunc: controller.enqueuePod,
    })

    return controller
}

func (c *PodController) enqueuePod(obj interface{}) {
    key, err := cache.MetaNamespaceKeyFunc(obj)
    if err != nil {
        runtime.HandleError(err)
        return
    }
    c.queue.Add(key)
}

// Run 启动控制器
func (c *PodController) Run(ctx context.Context) {
    defer runtime.HandleCrash()
    defer c.queue.ShutDown()

    klog.Info("Starting Pod Controller")

    // 启动 Informer
    go c.informer.Run(ctx.Done())

    // 等待缓存同步
    if !cache.WaitForCacheSync(ctx.Done(), c.informer.HasSynced) {
        runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
        return
    }

    klog.Info("Caches synced, starting workers")

    // 启动工作协程
    for i := 0; i < c.workerCount; i++ {
        go wait.UntilWithContext(ctx, c.runWorker, time.Second)
    }

    <-ctx.Done()
    klog.Info("Shutting down Pod Controller")
}

func (c *PodController) runWorker(ctx context.Context) {
    for c.processNextItem(ctx) {
    }
}

func (c *PodController) processNextItem(ctx context.Context) bool {
    key, quit := c.queue.Get()
    if quit {
        return false
    }
    defer c.queue.Done(key)

    err := c.syncHandler(ctx, key.(string))

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

func (c *PodController) syncHandler(ctx context.Context, key string) error {
    namespace, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
        return nil
    }

    // 从缓存获取 Pod
    obj, exists, err := c.informer.GetIndexer().GetByKey(key)
    if err != nil {
        return err
    }

    if !exists {
        klog.Infof("Pod %s/%s does not exist anymore\n", namespace, name)
        return nil
    }

    pod := obj.(*corev1.Pod)

    // 执行协调逻辑
    klog.Infof("Reconciling Pod: %s/%s (Phase: %s)\n", pod.Namespace, pod.Name, pod.Status.Phase)

    // 示例：检查 Pod 状态并执行操作
    if pod.Status.Phase == corev1.PodFailed {
        // 处理失败的 Pod
        return c.handleFailedPod(ctx, pod)
    }

    if pod.Status.Phase == corev1.PodPending {
        // 检查 Pending 原因
        return c.handlePendingPod(ctx, pod)
    }

    return nil
}

func (c *PodController) handleFailedPod(ctx context.Context, pod *corev1.Pod) error {
    klog.Infof("Handling failed Pod: %s/%s\n", pod.Namespace, pod.Name)

    // 实现失败处理逻辑
    // 例如：发送告警、记录事件等

    return nil
}

func (c *PodController) handlePendingPod(ctx context.Context, pod *corev1.Pod) error {
    klog.Infof("Handling pending Pod: %s/%s\n", pod.Namespace, pod.Name)

    // 检查 Pending 原因
    for _, condition := range pod.Status.Conditions {
        if condition.Type == corev1.PodScheduled && condition.Status == corev1.ConditionFalse {
            klog.Infof("Pod %s/%s not scheduled: %s\n", pod.Namespace, pod.Name, condition.Reason)
        }
    }

    return nil
}

// 带 Leader Election 的控制器
func runWithLeaderElection(ctx context.Context, clientset kubernetes.Interface) {
    // 创建资源锁
    lock := &resourcelock.LeaseLock{
        LeaseMeta: metav1.ObjectMeta{
            Name:      "pod-controller",
            Namespace: "kube-system",
        },
        Client: clientset.CoordinationV1(),
        LockConfig: resourcelock.ResourceLockConfig{
            Identity: os.Getenv("POD_NAME"),
        },
    }

    // 启动 Leader Election
    leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
        Lock:            lock,
        LeaseDuration:   15 * time.Second,
        RenewDeadline:   10 * time.Second,
        RetryPeriod:     2 * time.Second,
        ReleaseOnCancel: true,
        Callbacks: leaderelection.LeaderCallbacks{
            OnStartedLeading: func(ctx context.Context) {
                klog.Info("Became leader, starting controller")
                controller := NewPodController(clientset)
                controller.Run(ctx)
            },
            OnStoppedLeading: func() {
                klog.Info("Lost leadership, stopping controller")
            },
            OnNewLeader: func(identity string) {
                klog.Infof("New leader elected: %s", identity)
            },
        },
    })
}

func main() {
    // 创建 Kubernetes 客户端配置
    config, err := rest.InClusterConfig()
    if err != nil {
        klog.Fatalf("Error building kubeconfig: %s", err.Error())
    }

    // 创建 Kubernetes 客户端
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
    }

    // 创建上下文
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 处理信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigCh
        klog.Info("Received shutdown signal")
        cancel()
    }()

    // 启动控制器（带 Leader Election）
    runWithLeaderElection(ctx, clientset)
}
```

---

### 7.2 完整 Informer 模式实现

```go
// informer_example.go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/fields"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/cache"
)

// CustomIndexer 自定义索引器
type CustomIndexer struct {
    indexer cache.Indexer
}

// NewCustomIndexer 创建自定义索引器
func NewCustomIndexer() *CustomIndexer {
    // 创建自定义索引函数
    indexers := cache.Indexers{
        cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
        "nodeName": func(obj interface{}) ([]string, error) {
            pod, ok := obj.(*corev1.Pod)
            if !ok {
                return []string{}, nil
            }
            return []string{pod.Spec.NodeName}, nil
        },
        "phase": func(obj interface{}) ([]string, error) {
            pod, ok := obj.(*corev1.Pod)
            if !ok {
                return []string{}, nil
            }
            return []string{string(pod.Status.Phase)}, nil
        },
    }

    indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, indexers)

    return &CustomIndexer{indexer: indexer}
}

// InformerExample Informer 示例
type InformerExample struct {
    clientset      kubernetes.Interface
    sharedInformer informers.SharedInformerFactory
    podInformer    cache.SharedIndexInformer
    nodeInformer   cache.SharedIndexInformer
}

// NewInformerExample 创建 Informer 示例
func NewInformerExample(clientset kubernetes.Interface) *InformerExample {
    // 创建共享 Informer 工厂
    sharedInformer := informers.NewSharedInformerFactoryWithOptions(
        clientset,
        10*time.Minute,
        informers.WithNamespace(""),
        informers.WithTweakListOptions(func(options *metav1.ListOptions) {
            options.FieldSelector = fields.Everything().String()
        }),
    )

    return &InformerExample{
        clientset:      clientset,
        sharedInformer: sharedInformer,
        podInformer:    sharedInformer.Core().V1().Pods().Informer(),
        nodeInformer:   sharedInformer.Core().V1().Nodes().Informer(),
    }
}

// SetupHandlers 设置事件处理器
func (ie *InformerExample) SetupHandlers() {
    // Pod 事件处理器
    ie.podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            pod := obj.(*corev1.Pod)
            fmt.Printf("[ADD] Pod: %s/%s (Node: %s, Phase: %s)\n",
                pod.Namespace, pod.Name, pod.Spec.NodeName, pod.Status.Phase)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            oldPod := oldObj.(*corev1.Pod)
            newPod := newObj.(*corev1.Pod)

            // 只处理有意义的变更
            if oldPod.Status.Phase != newPod.Status.Phase {
                fmt.Printf("[UPDATE] Pod: %s/%s Phase changed: %s -> %s\n",
                    newPod.Namespace, newPod.Name, oldPod.Status.Phase, newPod.Status.Phase)
            }

            if oldPod.Spec.NodeName != newPod.Spec.NodeName {
                fmt.Printf("[UPDATE] Pod: %s/%s Node changed: %s -> %s\n",
                    newPod.Namespace, newPod.Name, oldPod.Spec.NodeName, newPod.Spec.NodeName)
            }
        },
        DeleteFunc: func(obj interface{}) {
            pod := obj.(*corev1.Pod)
            fmt.Printf("[DELETE] Pod: %s/%s\n", pod.Namespace, pod.Name)
        },
    })

    // Node 事件处理器
    ie.nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            node := obj.(*corev1.Node)
            fmt.Printf("[ADD] Node: %s (Ready: %v)\n", node.Name, isNodeReady(node))
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            oldNode := oldObj.(*corev1.Node)
            newNode := newObj.(*corev1.Node)

            oldReady := isNodeReady(oldNode)
            newReady := isNodeReady(newNode)

            if oldReady != newReady {
                fmt.Printf("[UPDATE] Node: %s Ready changed: %v -> %v\n",
                    newNode.Name, oldReady, newReady)
            }
        },
        DeleteFunc: func(obj interface{}) {
            node := obj.(*corev1.Node)
            fmt.Printf("[DELETE] Node: %s\n", node.Name)
        },
    })
}

// Run 启动 Informer
func (ie *InformerExample) Run(ctx context.Context) {
    // 启动所有 Informer
    ie.sharedInformer.Start(ctx.Done())

    // 等待缓存同步
    fmt.Println("Waiting for cache sync...")
    if !cache.WaitForCacheSync(ctx.Done(), ie.podInformer.HasSynced, ie.nodeInformer.HasSynced) {
        fmt.Println("Failed to sync cache")
        return
    }
    fmt.Println("Cache synced")

    // 演示查询功能
    ie.demoQueries()

    <-ctx.Done()
}

// demoQueries 演示查询功能
func (ie *InformerExample) demoQueries() {
    fmt.Println("\n=== Query Examples ===")

    // 1. 获取所有 Pod
    allPods := ie.podInformer.GetIndexer().List()
    fmt.Printf("Total Pods: %d\n", len(allPods))

    // 2. 按命名空间查询
    kubeSystemPods, err := ie.podInformer.GetIndexer().ByIndex(cache.NamespaceIndex, "kube-system")
    if err != nil {
        fmt.Printf("Error listing kube-system pods: %v\n", err)
    } else {
        fmt.Printf("kube-system Pods: %d\n", len(kubeSystemPods))
    }

    // 3. 按节点查询
    nodeName := "" // 替换为实际节点名
    podsOnNode, err := ie.podInformer.GetIndexer().ByIndex("nodeName", nodeName)
    if err != nil {
        fmt.Printf("Error listing pods on node %s: %v\n", nodeName, err)
    } else {
        fmt.Printf("Pods on node %s: %d\n", nodeName, len(podsOnNode))
    }

    // 4. 按状态查询
    runningPods, err := ie.podInformer.GetIndexer().ByIndex("phase", string(corev1.PodRunning))
    if err != nil {
        fmt.Printf("Error listing running pods: %v\n", err)
    } else {
        fmt.Printf("Running Pods: %d\n", len(runningPods))
    }

    // 5. 获取特定 Pod
    pod, exists, err := ie.podInformer.GetIndexer().GetByKey("default/nginx")
    if err != nil {
        fmt.Printf("Error getting pod: %v\n", err)
    } else if exists {
        p := pod.(*corev1.Pod)
        fmt.Printf("Found Pod: %s/%s\n", p.Namespace, p.Name)
    }
}

// DeltaFIFOExample DeltaFIFO 示例
type DeltaFIFOExample struct {
    fifo *cache.DeltaFIFO
}

// NewDeltaFIFOExample 创建 DeltaFIFO 示例
func NewDeltaFIFOExample() *DeltaFIFOExample {
    // 创建已知对象存储
    knownObjects := &knownObjectsStore{}

    // 创建 DeltaFIFO
    fifo := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
        KeyFunction:  cache.MetaNamespaceKeyFunc,
        KnownObjects: knownObjects,
    })

    return &DeltaFIFOExample{fifo: fifo}
}

// Run 运行 DeltaFIFO 示例
func (dfe *DeltaFIFOExample) Run(ctx context.Context) {
    // 生产者
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        count := 0
        for {
            select {
            case <-ticker.C:
                count++
                pod := &corev1.Pod{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      fmt.Sprintf("pod-%d", count),
                        Namespace: "default",
                    },
                    Status: corev1.PodStatus{
                        Phase: corev1.PodRunning,
                    },
                }

                // 添加 Delta
                dfe.fifo.Add(pod)

                // 模拟更新
                if count%3 == 0 {
                    pod.Status.Phase = corev1.PodSucceeded
                    dfe.fifo.Update(pod)
                }

            case <-ctx.Done():
                return
            }
        }
    }()

    // 消费者
    go func() {
        for {
            // 处理 Delta
            dfe.fifo.Pop(func(obj interface{}) error {
                deltas := obj.(cache.Deltas)

                for _, delta := range deltas {
                    pod := delta.Object.(*corev1.Pod)
                    fmt.Printf("Delta: %s, Pod: %s/%s\n", delta.Type, pod.Namespace, pod.Name)
                }

                return nil
            })
        }
    }()

    <-ctx.Done()
}

type knownObjectsStore struct {
    objects map[string]interface{}
}

func (k *knownObjectsStore) ListKeys() []string {
    return nil
}

func (k *knownObjectsStore) GetByKey(key string) (interface{}, bool, error) {
    return nil, false, nil
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
    // 创建 Kubernetes 客户端配置
    config, err := rest.InClusterConfig()
    if err != nil {
        // 尝试使用本地配置
        config, err = rest.InClusterConfig()
        if err != nil {
            panic(err)
        }
    }

    // 创建 Kubernetes 客户端
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    // 创建上下文
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 处理信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigCh
        fmt.Println("\nReceived shutdown signal")
        cancel()
    }()

    // 创建并运行 Informer 示例
    example := NewInformerExample(clientset)
    example.SetupHandlers()
    example.Run(ctx)
}
```

---

### 7.3 完整 Work Queue 模式实现

```go
// workqueue_example.go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "k8s.io/client-go/util/workqueue"
)

// Task 任务
type Task struct {
    ID       string
    Payload  interface{}
    Priority int
}

// WorkerPool 工作池
type WorkerPool struct {
    queue       workqueue.RateLimitingInterface
    workers     int
    processor   func(task *Task) error
    metrics     *Metrics
}

// Metrics 指标
type Metrics struct {
    Processed int64
    Failed    int64
    Retried   int64
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workers int, processor func(task *Task) error) *WorkerPool {
    // 创建速率限制队列
    queue := workqueue.NewRateLimitingQueue(
        workqueue.NewMaxOfRateLimiter(
            // 指数退避
            workqueue.NewItemExponentialFailureRateLimiter(
                5*time.Millisecond,  // 基础延迟
                1000*time.Second,    // 最大延迟
            ),
            // 令牌桶限流
            &workqueue.BucketRateLimiter{
                // 每秒 10 个令牌，桶容量 100
            },
        ),
    )

    return &WorkerPool{
        queue:     queue,
        workers:   workers,
        processor: processor,
        metrics:   &Metrics{},
    }
}

// Submit 提交任务
func (wp *WorkerPool) Submit(task *Task) {
    wp.queue.Add(task)
}

// SubmitWithDelay 延迟提交任务
func (wp *WorkerPool) SubmitWithDelay(task *Task, delay time.Duration) {
    // 使用延迟队列
    // 实际实现需要 DelayingInterface
    time.AfterFunc(delay, func() {
        wp.queue.Add(task)
    })
}

// Run 启动工作池
func (wp *WorkerPool) Run(ctx context.Context) {
    var wg sync.WaitGroup

    // 启动工作协程
    for i := 0; i < wp.workers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            wp.worker(ctx, id)
        }(i)
    }

    // 等待停止信号
    <-ctx.Done()

    // 停止队列
    wp.queue.ShutDown()

    // 等待工作协程完成
    wg.Wait()
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
    fmt.Printf("Worker %d started\n", id)

    for {
        item, shutdown := wp.queue.Get()
        if shutdown {
            fmt.Printf("Worker %d shutting down\n", id)
            return
        }

        task := item.(*Task)

        // 处理任务
        err := wp.processTask(task)

        wp.queue.Done(item)

        if err != nil {
            wp.metrics.Failed++

            // 检查重试次数
            if wp.queue.NumRequeues(item) < 5 {
                wp.metrics.Retried++
                wp.queue.AddRateLimited(item)
                fmt.Printf("Worker %d: Task %s failed, requeuing (retry %d)\n",
                    id, task.ID, wp.queue.NumRequeues(item))
            } else {
                wp.queue.Forget(item)
                fmt.Printf("Worker %d: Task %s failed after max retries, dropping\n", id, task.ID)
            }
        } else {
            wp.metrics.Processed++
            wp.queue.Forget(item)
            fmt.Printf("Worker %d: Task %s completed successfully\n", id, task.ID)
        }
    }
}

func (wp *WorkerPool) processTask(task *Task) error {
    fmt.Printf("Processing task: %s\n", task.ID)

    // 模拟处理
    time.Sleep(100 * time.Millisecond)

    // 模拟随机失败
    if time.Now().Unix()%5 == 0 {
        return fmt.Errorf("simulated error")
    }

    return wp.processor(task)
}

// GetMetrics 获取指标
func (wp *WorkerPool) GetMetrics() Metrics {
    return *wp.metrics
}

// PriorityQueue 优先级队列
type PriorityQueue struct {
    items []PriorityItem
    mu    sync.Mutex
}

type PriorityItem struct {
    Item     interface{}
    Priority int
}

// NewPriorityQueue 创建优先级队列
func NewPriorityQueue() *PriorityQueue {
    return &PriorityQueue{
        items: make([]PriorityItem, 0),
    }
}

func (pq *PriorityQueue) Add(item interface{}, priority int) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    pq.items = append(pq.items, PriorityItem{Item: item, Priority: priority})

    // 按优先级排序（高优先级在前）
    sort.Slice(pq.items, func(i, j int) bool {
        return pq.items[i].Priority > pq.items[j].Priority
    })
}

func (pq *PriorityQueue) Get() (interface{}, bool) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    if len(pq.items) == 0 {
        return nil, false
    }

    item := pq.items[0]
    pq.items = pq.items[1:]

    return item.Item, true
}

// DelayingQueueExample 延迟队列示例
type DelayingQueueExample struct {
    queue workqueue.DelayingInterface
}

// NewDelayingQueueExample 创建延迟队列示例
func NewDelayingQueueExample() *DelayingQueueExample {
    return &DelayingQueueExample{
        queue: workqueue.NewDelayingQueue(),
    }
}

// AddAfter 延迟添加
func (dqe *DelayingQueueExample) AddAfter(item interface{}, duration time.Duration) {
    dqe.queue.AddAfter(item, duration)
}

// Run 运行
func (dqe *DelayingQueueExample) Run(ctx context.Context) {
    go func() {
        for {
            item, shutdown := dqe.queue.Get()
            if shutdown {
                return
            }

            fmt.Printf("Processing delayed item: %v at %s\n", item, time.Now().Format(time.RFC3339))
            dqe.queue.Done(item)
        }
    }()

    <-ctx.Done()
    dqe.queue.ShutDown()
}

func main() {
    // 创建上下文
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 处理信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigCh
        fmt.Println("\nReceived shutdown signal")
        cancel()
    }()

    // 创建处理器
    processor := func(task *Task) error {
        fmt.Printf("Executing task: %s with payload: %v\n", task.ID, task.Payload)
        return nil
    }

    // 创建工作池
    pool := NewWorkerPool(3, processor)

    // 提交任务
    go func() {
        for i := 0; i < 20; i++ {
            task := &Task{
                ID:       fmt.Sprintf("task-%d", i),
                Payload:  fmt.Sprintf("data-%d", i),
                Priority: i % 5,
            }
            pool.Submit(task)
            time.Sleep(100 * time.Millisecond)
        }
    }()

    // 运行工作池
    pool.Run(ctx)

    // 打印指标
    metrics := pool.GetMetrics()
    fmt.Printf("\nMetrics:\n")
    fmt.Printf("  Processed: %d\n", metrics.Processed)
    fmt.Printf("  Failed: %d\n", metrics.Failed)
    fmt.Printf("  Retried: %d\n", metrics.Retried)
}
```

---

### 7.4 完整 Leader Election 实现

```go
// leader_election_example.go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/leaderelection"
    "k8s.io/client-go/tools/leaderelection/resourcelock"
    "k8s.io/klog/v2"
)

// LeaderElectionConfig Leader Election 配置
type LeaderElectionConfig struct {
    LeaseLockName      string
    LeaseLockNamespace string
    PodName            string
    PodNamespace       string
}

// LeaderElectedTask Leader 选举后执行的任务
type LeaderElectedTask struct {
    name string
    fn   func(ctx context.Context)
}

// LeaderElectionManager Leader Election 管理器
type LeaderElectionManager struct {
    config LeaderElectionConfig
    client kubernetes.Interface
    tasks  []LeaderElectedTask
}

// NewLeaderElectionManager 创建 Leader Election 管理器
func NewLeaderElectionManager(config LeaderElectionConfig, client kubernetes.Interface) *LeaderElectionManager {
    return &LeaderElectionManager{
        config: config,
        client: client,
        tasks:  make([]LeaderElectedTask, 0),
    }
}

// AddTask 添加任务
func (lem *LeaderElectionManager) AddTask(name string, fn func(ctx context.Context)) {
    lem.tasks = append(lem.tasks, LeaderElectedTask{name: name, fn: fn})
}

// Run 启动 Leader Election
func (lem *LeaderElectionManager) Run(ctx context.Context) {
    // 创建资源锁
    lock := &resourcelock.LeaseLock{
        LeaseMeta: metav1.ObjectMeta{
            Name:      lem.config.LeaseLockName,
            Namespace: lem.config.LeaseLockNamespace,
        },
        Client: lem.client.CoordinationV1(),
        LockConfig: resourcelock.ResourceLockConfig{
            Identity: lem.config.PodName,
        },
    }

    // 配置 Leader Election
    lec := leaderelection.LeaderElectionConfig{
        Lock:            lock,
        LeaseDuration:   15 * time.Second,
        RenewDeadline:   10 * time.Second,
        RetryPeriod:     2 * time.Second,
        ReleaseOnCancel: true,
        Callbacks: leaderelection.LeaderCallbacks{
            OnStartedLeading: func(ctx context.Context) {
                klog.Infof("[%s] Became leader, starting tasks", lem.config.PodName)

                // 启动所有任务
                for _, task := range lem.tasks {
                    go task.fn(ctx)
                }
            },
            OnStoppedLeading: func() {
                klog.Infof("[%s] Lost leadership, stopping tasks", lem.config.PodName)
            },
            OnNewLeader: func(identity string) {
                if identity != lem.config.PodName {
                    klog.Infof("[%s] New leader elected: %s", lem.config.PodName, identity)
                }
            },
        },
    }

    // 启动 Leader Election
    leaderelection.RunOrDie(ctx, lec)
}

// ScheduledTask 定时任务
type ScheduledTask struct {
    name     string
    interval time.Duration
    fn       func()
}

// NewScheduledTask 创建定时任务
func NewScheduledTask(name string, interval time.Duration, fn func()) *ScheduledTask {
    return &ScheduledTask{
        name:     name,
        interval: interval,
        fn:       fn,
    }
}

// Run 运行定时任务
func (st *ScheduledTask) Run(ctx context.Context) {
    ticker := time.NewTicker(st.interval)
    defer ticker.Stop()

    // 立即执行一次
    st.fn()

    for {
        select {
        case <-ticker.C:
            st.fn()
        case <-ctx.Done():
            klog.Infof("Task %s stopped", st.name)
            return
        }
    }
}

// LeaderAwareController Leader 感知控制器
type LeaderAwareController struct {
    name      string
    isLeader  bool
    mu        sync.RWMutex
}

// NewLeaderAwareController 创建 Leader 感知控制器
func NewLeaderAwareController(name string) *LeaderAwareController {
    return &LeaderAwareController{name: name}
}

// SetLeader 设置 Leader 状态
func (lac *LeaderAwareController) SetLeader(isLeader bool) {
    lac.mu.Lock()
    defer lac.mu.Unlock()
    lac.isLeader = isLeader
}

// IsLeader 检查是否是 Leader
func (lac *LeaderAwareController) IsLeader() bool {
    lac.mu.RLock()
    defer lac.mu.RUnlock()
    return lac.isLeader
}

// Run 运行控制器（仅在 Leader 时执行）
func (lac *LeaderAwareController) Run(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if lac.IsLeader() {
                lac.doWork()
            } else {
                klog.V(2).Infof("[%s] Not leader, skipping work", lac.name)
            }
        case <-ctx.Done():
            return
        }
    }
}

func (lac *LeaderAwareController) doWork() {
    klog.Infof("[%s] Doing work as leader", lac.name)
    // 执行实际工作
}

// MultiLockLeaderElection 多锁 Leader Election
type MultiLockLeaderElection struct {
    locks  []resourcelock.Interface
    config leaderelection.LeaderElectionConfig
}

// NewMultiLockLeaderElection 创建多锁 Leader Election
func NewMultiLockLeaderElection(config leaderelection.LeaderElectionConfig, locks ...resourcelock.Interface) *MultiLockLeaderElection {
    return &MultiLockLeaderElection{
        locks:  locks,
        config: config,
    }
}

// 使用示例
func main() {
    // 获取 Pod 信息
    podName := os.Getenv("POD_NAME")
    if podName == "" {
        podName = "leader-election-example"
    }

    podNamespace := os.Getenv("POD_NAMESPACE")
    if podNamespace == "" {
        podNamespace = "default"
    }

    // 创建 Kubernetes 客户端
    config, err := rest.InClusterConfig()
    if err != nil {
        klog.Fatalf("Error building kubeconfig: %s", err.Error())
    }

    client, err := kubernetes.NewForConfig(config)
    if err != nil {
        klog.Fatalf("Error building kubernetes client: %s", err.Error())
    }

    // 创建 Leader Election 配置
    leConfig := LeaderElectionConfig{
        LeaseLockName:      "example-leader-election",
        LeaseLockNamespace: podNamespace,
        PodName:            podName,
        PodNamespace:       podNamespace,
    }

    // 创建 Leader Election 管理器
    manager := NewLeaderElectionManager(leConfig, client)

    // 添加定时任务
    manager.AddTask("cleanup-task", func(ctx context.Context) {
        task := NewScheduledTask("cleanup", 30*time.Second, func() {
            klog.Info("Running cleanup task")
            // 执行清理逻辑
        })
        task.Run(ctx)
    })

    // 添加控制器任务
    controller := NewLeaderAwareController("example-controller")
    manager.AddTask("controller-task", func(ctx context.Context) {
        controller.SetLeader(true)
        controller.Run(ctx)
    })

    // 创建上下文
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 处理信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigCh
        klog.Info("Received shutdown signal")
        cancel()
    }()

    // 启动 Leader Election
    manager.Run(ctx)
}
```

---

### 7.5 完整 Sidecar 模式框架实现

```go
// sidecar_example.go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// Sidecar Sidecar 接口
type Sidecar interface {
    Name() string
    Start(ctx context.Context) error
    Stop() error
    Health() error
}

// SidecarManager Sidecar 管理器
type SidecarManager struct {
    sidecars []Sidecar
    mu       sync.RWMutex
}

// NewSidecarManager 创建 Sidecar 管理器
func NewSidecarManager() *SidecarManager {
    return &SidecarManager{
        sidecars: make([]Sidecar, 0),
    }
}

// Register 注册 Sidecar
func (sm *SidecarManager) Register(sidecar Sidecar) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.sidecars = append(sm.sidecars, sidecar)
}

// StartAll 启动所有 Sidecar
func (sm *SidecarManager) StartAll(ctx context.Context) error {
    sm.mu.RLock()
    sidecars := make([]Sidecar, len(sm.sidecars))
    copy(sidecars, sm.sidecars)
    sm.mu.RUnlock()

    for _, sidecar := range sidecars {
        go func(s Sidecar) {
            if err := s.Start(ctx); err != nil {
                fmt.Printf("Sidecar %s failed: %v\n", s.Name(), err)
            }
        }(sidecar)
    }

    return nil
}

// StopAll 停止所有 Sidecar
func (sm *SidecarManager) StopAll() error {
    sm.mu.RLock()
    sidecars := make([]Sidecar, len(sm.sidecars))
    copy(sidecars, sm.sidecars)
    sm.mu.RUnlock()

    for _, sidecar := range sidecars {
        if err := sidecar.Stop(); err != nil {
            fmt.Printf("Failed to stop sidecar %s: %v\n", sidecar.Name(), err)
        }
    }

    return nil
}

// HealthCheck 健康检查
func (sm *SidecarManager) HealthCheck() map[string]error {
    sm.mu.RLock()
    sidecars := make([]Sidecar, len(sm.sidecars))
    copy(sidecars, sm.sidecars)
    sm.mu.RUnlock()

    results := make(map[string]error)
    for _, sidecar := range sidecars {
        results[sidecar.Name()] = sidecar.Health()
    }

    return results
}

// LoggingSidecar 日志收集 Sidecar
type LoggingSidecar struct {
    name      string
    logPath   string
    outputDir string
    server    *http.Server
}

// NewLoggingSidecar 创建日志收集 Sidecar
func NewLoggingSidecar(logPath, outputDir string) *LoggingSidecar {
    return &LoggingSidecar{
        name:      "logging-sidecar",
        logPath:   logPath,
        outputDir: outputDir,
    }
}

func (ls *LoggingSidecar) Name() string {
    return ls.name
}

func (ls *LoggingSidecar) Start(ctx context.Context) error {
    fmt.Printf("Starting %s\n", ls.name)

    // 启动日志收集
    go ls.collectLogs(ctx)

    // 启动 HTTP 服务
    mux := http.NewServeMux()
    mux.HandleFunc("/health", ls.healthHandler)
    mux.HandleFunc("/logs", ls.logsHandler)

    ls.server = &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    go func() {
        <-ctx.Done()
        ls.server.Shutdown(context.Background())
    }()

    return ls.server.ListenAndServe()
}

func (ls *LoggingSidecar) Stop() error {
    fmt.Printf("Stopping %s\n", ls.name)
    if ls.server != nil {
        return ls.server.Shutdown(context.Background())
    }
    return nil
}

func (ls *LoggingSidecar) Health() error {
    // 检查日志收集是否正常
    return nil
}

func (ls *LoggingSidecar) collectLogs(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // 收集日志
            ls.processLogs()
        case <-ctx.Done():
            return
        }
    }
}

func (ls *LoggingSidecar) processLogs() {
    // 实现日志收集逻辑
    fmt.Printf("Processing logs from %s\n", ls.logPath)
}

func (ls *LoggingSidecar) healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func (ls *LoggingSidecar) logsHandler(w http.ResponseWriter, r *http.Request) {
    // 返回日志
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Logs..."))
}

// MetricsSidecar 监控 Sidecar
type MetricsSidecar struct {
    name   string
    port   int
    server *http.Server
}

// NewMetricsSidecar 创建监控 Sidecar
func NewMetricsSidecar(port int) *MetricsSidecar {
    return &MetricsSidecar{
        name: "metrics-sidecar",
        port: port,
    }
}

func (ms *MetricsSidecar) Name() string {
    return ms.name
}

func (ms *MetricsSidecar) Start(ctx context.Context) error {
    fmt.Printf("Starting %s on port %d\n", ms.name, ms.port)

    mux := http.NewServeMux()
    mux.HandleFunc("/metrics", ms.metricsHandler)
    mux.HandleFunc("/health", ms.healthHandler)

    ms.server = &http.Server{
        Addr:    fmt.Sprintf(":%d", ms.port),
        Handler: mux,
    }

    go func() {
        <-ctx.Done()
        ms.server.Shutdown(context.Background())
    }()

    return ms.server.ListenAndServe()
}

func (ms *MetricsSidecar) Stop() error {
    fmt.Printf("Stopping %s\n", ms.name)
    if ms.server != nil {
        return ms.server.Shutdown(context.Background())
    }
    return nil
}

func (ms *MetricsSidecar) Health() error {
    return nil
}

func (ms *MetricsSidecar) metricsHandler(w http.ResponseWriter, r *http.Request) {
    // 返回 Prometheus 格式的指标
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("# HELP example_metric Example metric\n"))
    w.Write([]byte("# TYPE example_metric gauge\n"))
    w.Write([]byte("example_metric 42\n"))
}

func (ms *MetricsSidecar) healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

// ProxySidecar 代理 Sidecar
type ProxySidecar struct {
    name       string
    listenAddr string
    targetAddr string
    server     *http.Server
}

// NewProxySidecar 创建代理 Sidecar
func NewProxySidecar(listenAddr, targetAddr string) *ProxySidecar {
    return &ProxySidecar{
        name:       "proxy-sidecar",
        listenAddr: listenAddr,
        targetAddr: targetAddr,
    }
}

func (ps *ProxySidecar) Name() string {
    return ps.name
}

func (ps *ProxySidecar) Start(ctx context.Context) error {
    fmt.Printf("Starting %s, proxying %s to %s\n", ps.name, ps.listenAddr, ps.targetAddr)

    mux := http.NewServeMux()
    mux.HandleFunc("/", ps.proxyHandler)

    ps.server = &http.Server{
        Addr:    ps.listenAddr,
        Handler: mux,
    }

    go func() {
        <-ctx.Done()
        ps.server.Shutdown(context.Background())
    }()

    return ps.server.ListenAndServe()
}

func (ps *ProxySidecar) Stop() error {
    fmt.Printf("Stopping %s\n", ps.name)
    if ps.server != nil {
        return ps.server.Shutdown(context.Background())
    }
    return nil
}

func (ps *ProxySidecar) Health() error {
    // 检查代理是否正常
    resp, err := http.Get("http://" + ps.listenAddr + "/health")
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unhealthy")
    }

    return nil
}

func (ps *ProxySidecar) proxyHandler(w http.ResponseWriter, r *http.Request) {
    // 转发请求到目标
    targetURL := "http://" + ps.targetAddr + r.URL.Path

    resp, err := http.Get(targetURL)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // 复制响应
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

// 使用示例
func main() {
    // 创建 Sidecar 管理器
    manager := NewSidecarManager()

    // 注册 Sidecar
    manager.Register(NewLoggingSidecar("/var/log/app", "/var/log/collected"))
    manager.Register(NewMetricsSidecar(9090))
    manager.Register(NewProxySidecar(":8080", "localhost:8081"))

    // 创建上下文
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 处理信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigCh
        fmt.Println("\nReceived shutdown signal")
        manager.StopAll()
        cancel()
    }()

    // 启动所有 Sidecar
    if err := manager.StartAll(ctx); err != nil {
        fmt.Printf("Failed to start sidecars: %v\n", err)
        os.Exit(1)
    }

    // 定期健康检查
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            results := manager.HealthCheck()
            for name, err := range results {
                if err != nil {
                    fmt.Printf("Sidecar %s is unhealthy: %v\n", name, err)
                } else {
                    fmt.Printf("Sidecar %s is healthy\n", name)
                }
            }
        case <-ctx.Done():
            return
        }
    }
}
```

---

### 7.6 完整 Circuit Breaker 实现

```go
// circuit_breaker_example.go
package main

import (
    "context"
    "errors"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

// State 断路器状态
type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

func (s State) String() string {
    switch s {
    case StateClosed:
        return "CLOSED"
    case StateOpen:
        return "OPEN"
    case StateHalfOpen:
        return "HALF_OPEN"
    default:
        return "UNKNOWN"
    }
}

// CircuitBreaker 断路器
type CircuitBreaker struct {
    name              string
    failureThreshold  int
    successThreshold  int
    timeout           time.Duration
    halfOpenMaxCalls  int

    state             State
    failures          int
    successes         int
    consecutiveSuccesses int
    lastFailureTime   time.Time
    halfOpenCalls     int

    mu                sync.RWMutex
    onStateChange     func(from, to State)
}

// CircuitBreakerConfig 断路器配置
type CircuitBreakerConfig struct {
    Name             string
    FailureThreshold int
    SuccessThreshold int
    Timeout          time.Duration
    HalfOpenMaxCalls int
    OnStateChange    func(from, to State)
}

// NewCircuitBreaker 创建断路器
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
    if config.FailureThreshold == 0 {
        config.FailureThreshold = 5
    }
    if config.SuccessThreshold == 0 {
        config.SuccessThreshold = 3
    }
    if config.Timeout == 0 {
        config.Timeout = 30 * time.Second
    }
    if config.HalfOpenMaxCalls == 0 {
        config.HalfOpenMaxCalls = 3
    }

    return &CircuitBreaker{
        name:             config.Name,
        failureThreshold: config.FailureThreshold,
        successThreshold: config.SuccessThreshold,
        timeout:          config.Timeout,
        halfOpenMaxCalls: config.HalfOpenMaxCalls,
        state:            StateClosed,
        onStateChange:    config.OnStateChange,
    }
}

// Execute 执行受保护的操作
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
    if err := cb.beforeRequest(); err != nil {
        return err
    }

    err := fn()
    cb.afterRequest(err)

    return err
}

// ExecuteWithFallback 执行带回退的操作
func (cb *CircuitBreaker) ExecuteWithFallback(ctx context.Context, fn func() error, fallback func() error) error {
    if err := cb.beforeRequest(); err != nil {
        // 断路器打开，执行回退
        if fallback != nil {
            return fallback()
        }
        return err
    }

    err := fn()
    cb.afterRequest(err)

    return err
}

func (cb *CircuitBreaker) beforeRequest() error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    switch cb.state {
    case StateClosed:
        return nil

    case StateOpen:
        // 检查是否超时
        if time.Since(cb.lastFailureTime) > cb.timeout {
            cb.transitionTo(StateHalfOpen)
            cb.halfOpenCalls = 0
            return nil
        }
        return ErrCircuitOpen

    case StateHalfOpen:
        if cb.halfOpenCalls >= cb.halfOpenMaxCalls {
            return ErrCircuitOpen
        }
        cb.halfOpenCalls++
        return nil
    }

    return nil
}

func (cb *CircuitBreaker) afterRequest(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err == nil {
        cb.onSuccess()
    } else {
        cb.onFailure()
    }
}

func (cb *CircuitBreaker) onSuccess() {
    switch cb.state {
    case StateClosed:
        cb.failures = 0

    case StateHalfOpen:
        cb.consecutiveSuccesses++
        if cb.consecutiveSuccesses >= cb.successThreshold {
            cb.transitionTo(StateClosed)
            cb.failures = 0
            cb.consecutiveSuccesses = 0
        }
    }
}

func (cb *CircuitBreaker) onFailure() {
    cb.failures++
    cb.lastFailureTime = time.Now()

    switch cb.state {
    case StateClosed:
        if cb.failures >= cb.failureThreshold {
            cb.transitionTo(StateOpen)
        }

    case StateHalfOpen:
        cb.transitionTo(StateOpen)
    }
}

func (cb *CircuitBreaker) transitionTo(newState State) {
    oldState := cb.state
    cb.state = newState

    if cb.onStateChange != nil {
        cb.onStateChange(oldState, newState)
    }

    fmt.Printf("CircuitBreaker %s: %s -> %s\n", cb.name, oldState, newState)
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() State {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    return cb.state
}

// GetMetrics 获取指标
func (cb *CircuitBreaker) GetMetrics() map[string]interface{} {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    return map[string]interface{}{
        "name":                cb.name,
        "state":               cb.state.String(),
        "failures":            cb.failures,
        "successes":           cb.successes,
        "consecutiveSuccesses": cb.consecutiveSuccesses,
        "lastFailureTime":     cb.lastFailureTime,
    }
}

// ErrCircuitOpen 断路器打开错误
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreakerGroup 断路器组
type CircuitBreakerGroup struct {
    breakers map[string]*CircuitBreaker
    mu       sync.RWMutex
    config   CircuitBreakerConfig
}

// NewCircuitBreakerGroup 创建断路器组
func NewCircuitBreakerGroup(config CircuitBreakerConfig) *CircuitBreakerGroup {
    return &CircuitBreakerGroup{
        breakers: make(map[string]*CircuitBreaker),
        config:   config,
    }
}

// Get 获取或创建断路器
func (cbg *CircuitBreakerGroup) Get(name string) *CircuitBreaker {
    cbg.mu.RLock()
    cb, ok := cbg.breakers[name]
    cbg.mu.RUnlock()

    if ok {
        return cb
    }

    cbg.mu.Lock()
    defer cbg.mu.Unlock()

    // 双重检查
    if cb, ok := cbg.breakers[name]; ok {
        return cb
    }

    config := cbg.config
    config.Name = name
    cb = NewCircuitBreaker(config)
    cbg.breakers[name] = cb

    return cb
}

// GetAllMetrics 获取所有断路器指标
func (cbg *CircuitBreakerGroup) GetAllMetrics() map[string]map[string]interface{} {
    cbg.mu.RLock()
    defer cbg.mu.RUnlock()

    metrics := make(map[string]map[string]interface{})
    for name, cb := range cbg.breakers {
        metrics[name] = cb.GetMetrics()
    }

    return metrics
}

// HTTPClient 带断路器的 HTTP 客户端
type HTTPClient struct {
    baseClient *http.Client
    breakers   *CircuitBreakerGroup
}

// NewHTTPClient 创建 HTTP 客户端
func NewHTTPClient(timeout time.Duration, breakerConfig CircuitBreakerConfig) *HTTPClient {
    return &HTTPClient{
        baseClient: &http.Client{
            Timeout: timeout,
        },
        breakers: NewCircuitBreakerGroup(breakerConfig),
    }
}

// Do 执行 HTTP 请求
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
    host := req.URL.Host
    cb := c.breakers.Get(host)

    var resp *http.Response
    var err error

    execErr := cb.Execute(req.Context(), func() error {
        resp, err = c.baseClient.Do(req)
        if err != nil {
            return err
        }

        // 5xx 错误视为服务故障
        if resp.StatusCode >= 500 {
            body, _ := io.ReadAll(resp.Body)
            resp.Body.Close()
            return fmt.Errorf("server error: %d, body: %s", resp.StatusCode, string(body))
        }

        return nil
    })

    if execErr != nil {
        return nil, execErr
    }

    return resp, err
}

// Get 执行 GET 请求
func (c *HTTPClient) Get(url string) (*http.Response, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    return c.Do(req)
}

// Post 执行 POST 请求
func (c *HTTPClient) Post(url string, body io.Reader) (*http.Response, error) {
    req, err := http.NewRequest("POST", url, body)
    if err != nil {
        return nil, err
    }

    return c.Do(req)
}

// GRPCClient 带断路器的 gRPC 客户端
type GRPCClient struct {
    breakers *CircuitBreakerGroup
}

// NewGRPCClient 创建 gRPC 客户端
func NewGRPCClient(breakerConfig CircuitBreakerConfig) *GRPCClient {
    return &GRPCClient{
        breakers: NewCircuitBreakerGroup(breakerConfig),
    }
}

// Invoke 调用 gRPC 方法
func (c *GRPCClient) Invoke(ctx context.Context, serviceName string, fn func() error) error {
    cb := c.breakers.Get(serviceName)
    return cb.Execute(ctx, fn)
}

// MetricsServer 指标服务器
type MetricsServer struct {
    breakers *CircuitBreakerGroup
    port     int
}

// NewMetricsServer 创建指标服务器
func NewMetricsServer(breakers *CircuitBreakerGroup, port int) *MetricsServer {
    return &MetricsServer{
        breakers: breakers,
        port:     port,
    }
}

// Start 启动服务器
func (ms *MetricsServer) Start(ctx context.Context) error {
    mux := http.NewServeMux()
    mux.HandleFunc("/metrics", ms.metricsHandler)
    mux.HandleFunc("/health", ms.healthHandler)

    server := &http.Server{
        Addr:    fmt.Sprintf(":%d", ms.port),
        Handler: mux,
    }

    go func() {
        <-ctx.Done()
        server.Shutdown(context.Background())
    }()

    return server.ListenAndServe()
}

func (ms *MetricsServer) metricsHandler(w http.ResponseWriter, r *http.Request) {
    metrics := ms.breakers.GetAllMetrics()

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, "{\n")

    first := true
    for name, m := range metrics {
        if !first {
            fmt.Fprintf(w, ",\n")
        }
        first = false

        fmt.Fprintf(w, "  %q: {\n", name)
        fmt.Fprintf(w, "    \"state\": %q,\n", m["state"])
        fmt.Fprintf(w, "    \"failures\": %v\n", m["failures"])
        fmt.Fprintf(w, "  }")
    }

    fmt.Fprintf(w, "\n}\n")
}

func (ms *MetricsServer) healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

// 使用示例
func main() {
    // 创建断路器配置
    breakerConfig := CircuitBreakerConfig{
        Name:             "default",
        FailureThreshold: 5,
        SuccessThreshold: 3,
        Timeout:          30 * time.Second,
        HalfOpenMaxCalls: 3,
        OnStateChange: func(from, to State) {
            fmt.Printf("Circuit breaker state changed: %s -> %s\n", from, to)
        },
    }

    // 创建 HTTP 客户端
    httpClient := NewHTTPClient(10*time.Second, breakerConfig)

    // 创建指标服务器
    metricsServer := NewMetricsServer(httpClient.breakers, 9090)

    // 创建上下文
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 处理信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigCh
        fmt.Println("\nReceived shutdown signal")
        cancel()
    }()

    // 启动指标服务器
    go metricsServer.Start(ctx)

    // 模拟请求
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    requestCount := 0
    for {
        select {
        case <-ticker.C:
            requestCount++

            // 模拟失败
            url := "http://localhost:8080/api/test"
            if requestCount%10 == 0 {
                url = "http://invalid-host:8080/api/test" // 模拟失败
            }

            resp, err := httpClient.Get(url)
            if err != nil {
                fmt.Printf("Request failed: %v\n", err)
            } else {
                resp.Body.Close()
                fmt.Printf("Request succeeded: %d\n", resp.StatusCode)
            }

        case <-ctx.Done():
            return
        }
    }
}
```

---

## 总结

本章全面分析了 Docker 和 Kubernetes 中的设计模式，包括：

1. **GoF 设计模式**：单例、工厂、策略、观察者、模板方法、适配器、装饰器
2. **分布式系统设计模式**：Sidecar、Ambassador、Adapter、Scatter-Gather、Saga、Circuit Breaker、Bulkhead、Retry
3. **并发与并行模式**：Worker Pool、Pub/Sub、Leader Election、Distributed Lock、Barrier、Pipeline
4. **同步与异步模式**：同步调用、异步消息、Watch、Informer、Work Queue、Channel
5. **工作流设计模式**：CronJob、Job、DaemonSet、Pipeline、DAG、State Machine
6. **K8s 特有模式**：Controller、Operator、Initializer、Finalizer、Owner Reference

每种模式都配有详细的 Go 代码实现，可直接用于实际项目开发。

---

*文档版本: 1.0*
*最后更新: 2024*
