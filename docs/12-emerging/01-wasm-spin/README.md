# WebAssembly (WASM) on Kubernetes - 云原生新边界

## 概述

WebAssembly (WASM) 是轻量级、可移植的二进制指令格式，为云原生带来次秒级启动、沙箱安全和极致资源效率。2025 年，Wasm 已成为边缘计算和 Serverless 的主流选择。

> **关键数据**: Akamai 收购 Fermyon 后，Spin 在其边缘网络上实现 7500 万 RPS 吞吐，冷启动 <10ms。

## 容器 vs WebAssembly

```
对比维度:
┌─────────────────┬─────────────────────┬─────────────────────┐
│     特性        │      容器           │    WebAssembly      │
├─────────────────┼─────────────────────┼─────────────────────┤
│ 镜像大小        │ 100-500 MB          │ 1-10 MB             │
│ 启动时间        │ 100-1000 ms         │ 1-10 ms             │
│ 内存占用        │ 100+ MB             │ 5-35 MB             │
│ 安全模型        │ 共享内核            │ 沙箱 + 能力模型      │
│ 可移植性        │ OS/Arch 依赖        │ 一次编译，到处运行   │
│ 冷启动成本      │ 高                  │ 极低                │
│ 适用场景        │ 长时间运行服务      │ 事件驱动、边缘计算   │
└─────────────────┴─────────────────────┴─────────────────────┘
```

## WebAssembly 运行时

| 运行时 | 特点 | 适用场景 |
|--------|------|----------|
| **Wasmtime** | Bytecode Alliance 官方，安全优先 | 通用工作负载 |
| **WasmEdge** | 高性能，支持 AI/ML | 边缘 AI |
| **Spin** | Fermyon 开发，微服务优化 | Serverless |
| **Wasmer** | 多语言绑定 | 嵌入式 |

## Kubernetes 集成架构

```
┌─────────────────────────────────────────────────────────────────┐
│                    Kubernetes Control Plane                     │
│                     (API Server, Scheduler)                     │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────────────┐
│                      containerd Runtime                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │  runc (标准)    │  │  containerd-    │  │  containerd-    │ │
│  │  OCI 容器       │  │  shim-spin      │  │  shim-wasmedge  │ │
│  │                 │  │  (Spin/WASM)    │  │  (WasmEdge)     │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    RuntimeClass 配置                         ││
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       ││
│  │  │   runc       │  │ wasmtime-spin│  │  wasmedge    │       ││
│  │  │ (默认)       │  │ (WASM)       │  │ (WASM)       │       ││
│  │  └──────────────┘  └──────────────┘  └──────────────┘       ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## Spin on Kubernetes

### 安装 containerd-shim-spin

```bash
# 下载 shim
curl -fsSL https://github.com/containerd/runwasi/releases/download/v0.3.0/containerd-wasm-shim-v0.3.0-linux-x86_64.tar.gz \
  -o containerd-wasm-shim.tar.gz

# 安装到节点
sudo tar -C /usr/local/bin -xzf containerd-wasm-shim.tar.gz

# 验证
containerd-shim-spin-v1 --version
```

### 配置 containerd

```toml
# /etc/containerd/config.toml
version = 2

[plugins."io.containerd.grpc.v1.cri".containerd]
  default_runtime_name = "runc"

[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
  runtime_type = "io.containerd.runc.v2"

[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.spin]
  runtime_type = "io.containerd.spin.v1"
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.spin.options]
    BinaryName = "/usr/local/bin/containerd-shim-spin-v1"
```

```bash
# 重启 containerd
sudo systemctl restart containerd

# 验证运行时
crictl info | jq '.config.containerd.runtimes'
```

### 创建 RuntimeClass

```yaml
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: wasmtime-spin
handler: spin
scheduling:
  nodeSelector:
    wasm-enabled: "true"
```

```bash
kubectl apply -f runtimeclass.yaml

# 标记 WASM 节点
kubectl label nodes worker-1 wasm-enabled=true
kubectl label nodes worker-2 wasm-enabled=true
```

### 部署 Spin 应用

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: spin-hello
spec:
  replicas: 3
  selector:
    matchLabels:
      app: spin-hello
  template:
    metadata:
      labels:
        app: spin-hello
    spec:
      runtimeClassName: wasmtime-spin
      containers:
      - name: spin-app
        image: ghcr.io/fermyon/spin-hello:v1.0
        ports:
        - containerPort: 80
          name: http
        resources:
          requests:
            memory: "8Mi"
            cpu: "10m"
          limits:
            memory: "16Mi"
            cpu: "50m"
---
apiVersion: v1
kind: Service
metadata:
  name: spin-hello
spec:
  selector:
    app: spin-hello
  ports:
  - port: 80
    targetPort: 80
```

## 开发 Spin 应用

### 安装 Spin CLI

```bash
curl -fsSL https://developer.fermyon.com/downloads/install.sh | bash
sudo mv spin /usr/local/bin/
```

### 创建 Rust HTTP 服务

```bash
# 创建新项目
spin new http-rust my-api
cd my-api
```

```rust
// src/lib.rs
use spin_sdk::{
    http::{Request, Response},
    http_component,
};

#[http_component]
fn handle_request(req: Request) -> Result<Response, String> {
    match req.uri().path() {
        "/health" => Ok(http::Response::builder()
            .status(200)
            .header("content-type", "application/json")
            .body(Some(r#"{"status": "healthy"}"#.into()))?),

        "/api/users" => Ok(http::Response::builder()
            .status(200)
            .header("content-type", "application/json")
            .body(Some(r#"[{"id": 1, "name": "Alice"}]"#.into()))?),

        _ => Ok(http::Response::builder()
            .status(404)
            .body(Some("Not Found".into()))?),
    }
}
```

```toml
# spin.toml
spin_manifest_version = 2

[application]
name = "my-api"
version = "0.1.0"
authors = ["dev@example.com"]
description = "My WASM microservice"

[[trigger.http]]
route = "/..."
component = "my-api"

[component.my-api]
source = "target/wasm32-wasi/release/my_api.wasm"
allowed_outbound_hosts = []

[component.my-api.build]
command = "cargo build --target wasm32-wasi --release"
```

### 本地测试

```bash
# 构建
spin build

# 本地运行
spin up

# 测试
curl http://localhost:3000/health
curl http://localhost:3000/api/users
```

### 打包为 OCI 镜像

```bash
# 安装 registry 插件
spin plugins install registry

# 推送到仓库
spin registry push myregistry.io/my-api:v1.0

# 验证镜像大小
docker images myregistry.io/my-api:v1.0
# 预期: 4-5 MB（对比容器 200+ MB）
```

## 自动扩缩容

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: spin-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: spin-hello
  minReplicas: 3
  maxReplicas: 200  # WASM 快速启动支持激进扩缩容
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 50
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 10  # 快速扩容
    scaleDown:
      stabilizationWindowSeconds: 30
```

## 使用场景对比

```
场景1: API 网关认证
─────────────────────────────────────────
传统: Lambda@Edge (80-150ms 延迟)
WASM: Spin on Edge (3-7ms 延迟)
优势: 20x 性能提升，<10ms 冷启动

场景2: 图片处理
─────────────────────────────────────────
传统: 容器化服务 (200-350ms)
WASM: Spin 函数 (10-20ms)
优势: 内存占用降低 90%

场景3: 边缘 ML 推理 (7B 参数模型)
─────────────────────────────────────────
传统: 600+ ms
WASM: 40-80ms
优势: 通过 K8s 调度 GPU，聚合 7500 万 RPS

场景4: 事件驱动处理
─────────────────────────────────────────
传统: KEDA + 容器 (秒级启动)
WASM: Spin + KEDA (毫秒级启动)
优势: 真正的 scale-to-zero
```

## 与 Service Mesh 集成

```yaml
# Istio 支持 WASM 作为 Sidecar
apiVersion: apps/v1
kind: Deployment
metadata:
  name: wasm-with-istio
spec:
  template:
    metadata:
      annotations:
        sidecar.istio.io/inject: "true"
    spec:
      runtimeClassName: wasmtime-spin
      containers:
      - name: wasm-app
        image: myregistry.io/wasm-app:v1
        ports:
        - containerPort: 8080
```

## 监控和调试

```bash
# 查看 WASM Pod 状态
kubectl get pods -l app=spin-hello

# 查看日志（与普通容器相同）
kubectl logs -l app=spin-hello

# 进入调试（功能有限）
kubectl exec -it <pod> -- /bin/sh

# 性能分析
kubectl top pod -l app=spin-hello
```

## 生产注意事项

| 方面 | 建议 |
|------|------|
| 调试 | WASM 容器不支持标准 shell，需要远程调试 |
| 存储 | 使用 WASI 文件系统，持久化需外挂卷 |
| 网络 | 出站连接需显式声明 allowed_outbound_hosts |
| 生态 | 工具链仍在成熟，部分库兼容性待完善 |
| 团队 | Rust/Go/AssemblyScript 学习曲线 |

## 2025 发展趋势

1. **WASI 0.2** - 更完善的系统接口支持
2. **组件模型** - 可组合、可复用的 WASM 模块
3. **OCI 标准化** - WASM 镜像与容器镜像统一分发
4. **AI 推理优化** - WASI-NN 标准支持 GPU 加速

## 总结

| 场景 | 推荐技术 |
|------|----------|
| 长时间运行服务 | 传统容器 |
| Serverless/边缘 | WebAssembly |
| 事件驱动处理 | WebAssembly |
| 微服务 API | 两者皆可 |
| AI 推理 | WebAssembly + GPU |

WebAssembly 不是替代容器，而是补充。对于冷启动敏感、资源受限的场景，WASM 是 2025 年的最佳选择。
