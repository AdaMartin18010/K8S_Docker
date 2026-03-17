# WebAssembly (Wasm) 与容器

> 下一代云原生运行时的崛起

---

## 什么是 WebAssembly？

WebAssembly (Wasm) 是一种二进制指令格式，最初为浏览器设计，现已演变为通用运行时。它允许代码在任何环境中以接近原生的性能运行。

```
┌─────────────────────────────────────────────────────────────┐
│                    WebAssembly 架构                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   C/C++     │  │    Rust     │  │    Go       │          │
│  │   源码      │  │    源码      │  │   源码       │          │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘          │
│         │                │                │                 │
│         └────────────────┼────────────────┘                 │
│                          ↓                                  │
│              ┌─────────────────────┐                        │
│              │   Wasm 编译器        │                        │
│              │ (Clang, rustc, tinygo)                       │
│              └──────────┬──────────┘                        │
│                         ↓                                   │
│              ┌─────────────────────┐                        │
│              │    .wasm 二进制     │                        │
│              │  (平台无关、体积小)   │                        │
│              └──────────┬──────────┘                        │
│                         ↓                                   │
│         ┌───────────────┼───────────────┐                   │
│         ↓               ↓               ↓                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  WasmEdge   │  │  Wasmtime   │  │   Wasmer    │          │
│  │  (云原生)    │  │  (工业级)    │  │  (通用)      │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

---

## Wasm vs Docker 性能对比 (2025)

| 指标 | Docker | WebAssembly | 优势 |
|------|--------|-------------|------|
| **冷启动时间** | 500ms - 2s | **5ms - 20ms** | **100x 更快** |
| **镜像大小** | 100MB - 500MB | **1MB - 10MB** | **10-50x 更小** |
| **内存占用** | 100MB+ | **几 MB** | **显著降低** |
| **启动密度** | 数百/节点 | **数千/节点** | **10x+** |
| **沙箱开销** | 高 (OS 虚拟化) | **低 (轻量级)** | **更安全** |

### 冷启动对比

```
Docker Container:
┌──────────┬──────────┬──────────┬──────────┐
│ 拉取镜像  │ 解压层    │ 启动 OS   │ 启动应用  │
│  ~30s    │  ~10s    │  ~2s     │  ~1s     │
└──────────┴──────────┴──────────┴──────────┘
Total: ~43s

WebAssembly:
┌──────────┬──────────┐
│ 加载 wasm │ 实例化   │
│  ~5ms    │  ~15ms   │
└──────────┴──────────┘
Total: ~20ms
```

---

## Docker 对 Wasm 的支持 (2025)

Docker Desktop 现在原生支持 WebAssembly 工作负载，使用 containerd shim 执行 Wasm 应用。

### 支持的运行时

| 运行时 | 适用场景 | 特点 |
|--------|----------|------|
| **WasmEdge** | 云原生/边缘 | CNCF 项目，K8s 集成好 |
| **Wasmtime** | 企业/生产 | 标准合规，LTS 支持 |
| **Spin** | Serverless | 事件驱动，开发体验好 |
| **Wasmer** | 全栈开发 | 通用运行时 |

### 运行 Wasm 容器

```bash
# 使用 Docker 运行 Wasm 模块
docker run --runtime=io.containerd.wasmedge.v1 \
  --platform=wasm32/wasi \
  -p 8080:8080 \
  myapp:wasm-latest

# 构建 Wasm 镜像
docker buildx build --platform wasi/wasm \
  -t myapp:wasm .
```

---

## Kubernetes 中的 WebAssembly

### runwasi 项目

runwasi 允许 Wasm 工作负载与容器并行运行在 K8s 中。

```yaml
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: wasmtime-spin
handler: spin
scheduling:
  nodeSelector:
    wasm: "true"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: wasm-app
spec:
  template:
    spec:
      runtimeClassName: wasmtime-spin
      containers:
        - name: app
          image: ghcr.io/myorg/wasm-app:latest
          resources:
            limits:
              cpu: "0.1"
              memory: "32Mi"
```

---

## 适用场景

### ✅ 适合 Wasm

| 场景 | 原因 |
|------|------|
| **Serverless/函数计算** | 冷启动快，按需计费更精准 |
| **边缘计算** | 体积小，启动快，资源受限 |
| **微服务网关** | 低延迟，高吞吐 |
| **插件/扩展系统** | 安全沙箱，隔离性好 |
| **IoT 设备** | 极小的资源占用 |

### ❌ 不适合 Wasm

| 场景 | 原因 |
|------|------|
| **数据库/有状态服务** | 需要文件系统，fork/exec |
| **GPU 工作负载** | WASI 生态不完善 |
| **遗留应用** | 深度依赖操作系统 |
| **复杂网络应用** | 部分网络 API 受限 |

---

## WASI 演进路线图

| 版本 | 状态 | 关键特性 |
|------|------|----------|
| **Preview 1** | 已弃用 | 基础文件系统、环境变量 |
| **Preview 2** | 稳定 | Component Model、Sockets、HTTP |
| **Preview 3** | 开发中 | Async I/O、Streams |

---

## 迁移策略

```
Phase 1: 评估候选者
    - 识别无状态、计算密集型工作负载
    - 避开需要 fork/exec 的应用

Phase 2: 工具链准备
    - 安装 Rust/Go 编译器
    - 配置 wasm32-wasi 目标

Phase 3: 选择运行时
    - 生产: Wasmtime
    - 边缘: WasmEdge
    - 开发: Wasmer

Phase 4: 渐进式发布
    - 从非关键服务开始
    - 监控性能指标
    - 逐步扩大范围
```

---

## 2027 预测

- **Wasm 主导**: 边缘计算、Serverless、插件系统
- **Docker 继续**: 复杂有状态应用
- **混合部署**: 成为标准架构模式
