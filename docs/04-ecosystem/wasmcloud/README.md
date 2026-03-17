# WasmCloud - CNCF WebAssembly 应用运行时

## 概述

WasmCloud 是一个 CNCF 孵化项目，基于 WebAssembly (Wasm) 组件模型构建，支持使用多种编程语言开发可重用的 Wasm 组件，并在任何云平台、Kubernetes 集群、数据中心或边缘设备上弹性运行。

## 核心概念

### WebAssembly 组件模型

- **组件 (Component)**: 可组合、语言无关的 Wasm 二进制文件
- **接口 (WIT)**: WebAssembly 接口类型，定义组件间契约
- **能力 (Capability)**: 通过严格定义的接口访问外部资源

### 关键优势

| 特性 | 描述 |
|------|------|
| 多语言 | Rust、Go、Python、JavaScript、C/C++ 等 |
| 冷启动 | 亚毫秒级启动时间 |
| 沙箱安全 | 默认无权限，能力显式注入 |
| 可移植性 | 一次编译，随处运行 |
| 可组合性 | 组件像乐高积木一样组合 |

## 架构组件

### Wash CLI

WasmCloud 开发工具链：

```bash
# 安装 wash
curl -s https://wasmcloud.dev/install.sh | bash

# 创建新项目
wash new component hello --template rust
wash new component hello-go --template go

# 构建组件
cd hello
wash build

# 本地运行
wash up
wash app deploy ./wadm.yaml
```

### wasmCloud Host

运行时主机，管理组件生命周期：

```yaml
# wasmcloud-host.yaml
apiVersion: k8s.wasmcloud.dev/v1alpha1
kind: WasmCloudHostConfig
metadata:
  name: wasmcloud-host
spec:
  natsAddress: nats://nats:4222
  natsLeafImage: nats:2.10
  lattice: default
  version: 1.5.0
  resources:
    nats:
      limits:
        cpu: 500m
        memory: 512Mi
```

### WADM (wasmCloud Application Deployment Manager)

声明式应用部署：

```yaml
# wadm.yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: http-hello
  annotations:
    version: v1.0.0
spec:
  components:
  - name: http-component
    type: component
    properties:
      image: ghcr.io/wasmcloud/http-hello:0.1.0
    traits:
    - type: spreadscaler
      properties:
        replicas: 3
    - type: link
      properties:
        target: httpserver
        namespace: wasi
        package: http
        interfaces: [incoming-handler]

  - name: httpserver
    type: capability
    properties:
      image: ghcr.io/wasmcloud/http-server:0.26.0
    traits:
    - type: spreadscaler
      properties:
        replicas: 1
```

## 在 Kubernetes 上运行

### 安装 wasmCloud Operator

```bash
# 添加 Helm 仓库
helm repo add wasmcloud https://wasmcloud.github.io/wasmcloud-helm
helm repo update

# 安装 NATS（wasmCloud 依赖）
helm install nats nats/nats --set config.jetstream.enabled=true

# 安装 wasmCloud Operator
helm install wasmcloud-operator wasmcloud/wasmcloud-operator \
  --set image.tag=0.4.0

# 验证安装
kubectl get pods -n wasmcloud
```

### 部署 WebAssembly 应用

```yaml
# http-hello-app.yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: http-hello-world
  namespace: default
spec:
  components:
  - name: http-hello
    type: component
    properties:
      image: ghcr.io/wasmcloud/http-hello-world:0.1.0
      id: http-hello
    traits:
    - type: spreadscaler
      properties:
        replicas: 2

  - name: httpserver
    type: capability
    properties:
      image: ghcr.io/wasmcloud/http-server:0.26.0
    traits:
    - type: link
      properties:
        target: http-hello
        namespace: wasi
        package: http
        interfaces: [incoming-handler]
```

### 与 Gateway API 集成

```yaml
# wasmcloud-gateway.yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: wasmcloud-route
spec:
  parentRefs:
  - name: wasmcloud-gateway
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /api
    backendRefs:
    - name: wasmcloud-http-server
      port: 8080
```

## 开发示例

### Rust 组件

```rust
// src/lib.rs
use wasmcloud_component::http;

struct Component;

impl http::Server for Component {
    fn handle(
        request: http::Request,
    ) -> http::Result<http::Response> {
        Ok(http::Response {
            status: 200,
            headers: vec![("content-type".to_string(), "text/plain".to_string())],
            body: b"Hello from WebAssembly!".to_vec(),
        })
    }
}
```

### Go 组件

```go
// main.go
package main

import (
    "github.com/bytecodealliance/wasm-tools-go/cm"
    http "github.com/wasmcloud/wasmcloud/examples/go/http-hello-world/gen"
)

func init() {
    http.Exports.Handle = func(request http.Request) http.Response {
        return http.Response{
            StatusCode: 200,
            Headers: []http.Header{
                {Name: "Content-Type", Value: "text/plain"},
            },
            Body: cm.NewList("Hello from Go + WebAssembly!"),
        }
    }
}

func main() {}
```

## 企业场景

### 美国运通 (American Express)

使用 wasmCloud 构建内部 FaaS 平台：

- 编译代码为 Wasm 组件
- 包装安全装饰器
- 支持多种调用接口
- 多租户函数服务

### Adobe

使用 wasmCloud 将 C/C++ 应用编译为 WebAssembly 组件，跨平台部署。

### Akamai

2025年12月收购 Fermyon 后，使用 Spin + wasmCloud 提供边缘计算服务，7500万 RPS，亚毫秒冷启动。

## SpinKube

SpinKube 是 CNCF Sandbox 项目，将 Spin 框架与 Kubernetes 集成：

```yaml
# spinapp.yaml
apiVersion: spin.core.spinkube.dev/v1alpha1
kind: SpinApp
metadata:
  name: simple-spin-app
spec:
  image: ghcr.io/spinkube/simple-go-app:v0.1.0
  replicas: 2
  executor: containerd-shim-spin
```

```bash
# 安装 SpinKube
helm repo add spinkube https://spinkube.github.io/spinkube
helm install spinkube spinkube/spinkube-operator

# 部署应用
kubectl apply -f spinapp.yaml
```

## 2025 新特性

- **WASI 0.2**: 标准化 HTTP、文件 I/O、Socket 接口
- **WebAssembly 3.0**: W3C 标准，支持 GC、64位地址空间
- **SPIFFE 集成**: WebAssembly 工作负载身份
- **组件依赖管理**: OCI-based WIT 依赖管理
- **wash 插件系统**: Wasm 驱动的 CLI 插件

## 性能对比

| 指标 | 传统容器 | WebAssembly |
|------|---------|-------------|
| 冷启动 | 100-500ms | <1ms |
| 内存占用 | 100MB+ | 5-20MB |
| 镜像大小 | 100MB+ | 1-10MB |
| 密度 | 10-100/节点 | 1000+/节点 |

## 相关资源

- [wasmCloud 官网](https://wasmcloud.com/)
- [Wasm Component Model](https://component-model.bytecodealliance.org/)
- [WASI Preview 2](https://wasi.dev/)
- [Spin 框架](https://developer.fermyon.com/spin)
