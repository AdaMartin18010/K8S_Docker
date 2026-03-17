# 零信任安全架构 (Zero Trust Architecture)

## 概述

零信任安全模型遵循"永不信任，始终验证"原则，要求每个访问请求都必须经过认证、授权和加密，无论请求来自内部还是外部网络。

> **2025 关键数据**: GitHub State of Secrets Sprawl 报告显示 70% 的泄露密钥在 2 年后仍然有效。零信任通过短期凭证自动轮换解决此问题。

## 零信任架构组件

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Zero Trust Architecture                             │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Control Plane                                   │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │   │
│  │  │   IdP    │  │   PDP    │  │   PAP    │  │  SPIRE Server    │   │   │
│  │  │ (身份)   │  │(策略决策)│  │(策略管理)│  │  (身份颁发)      │   │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│  ┌─────────────────────────────────┴─────────────────────────────────┐     │
│  │                        Data Plane                                  │     │
│  │                                                                    │     │
│  │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐        │     │
│  │  │   Service A  │◄──►│   Service B  │◄──►│   Service C  │        │     │
│  │  │              │mTLS│              │mTLS│              │        │     │
│  │  │ ┌──────────┐ │    │ ┌──────────┐ │    │ ┌──────────┐ │        │     │
│  │  │ │SPIFFE ID │ │    │ │SPIFFE ID │ │    │ │SPIFFE ID │ │        │     │
│  │  │ │SVID(X.509)││    │ │SVID(X.509)││    │ │SVID(X.509)││        │     │
│  │  │ └──────────┘ │    │ └──────────┘ │    │ └──────────┘ │        │     │
│  │  └──────────────┘    └──────────────┘    └──────────────┘        │     │
│  │                                                                    │     │
│  │  ┌────────────────────────────────────────────────────────────┐   │     │
│  │  │              Service Mesh (Istio Ambient)                   │   │     │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                    │   │     │
│  │  │  │ ztunnel │  │ ztunnel │  │ ztunnel │ (每节点 L4)        │   │     │
│  │  │  │(mTLS)   │  │(mTLS)   │  │(mTLS)   │                    │   │     │
│  │  │  └─────────┘  └─────────┘  └─────────┘                    │   │     │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                    │   │     │
│  │  │  │ waypoints│ │ waypoints│ │ waypoints│ (L7，按需)        │   │     │
│  │  │  └─────────┘  └─────────┘  └─────────┘                    │   │     │
│  │  └────────────────────────────────────────────────────────────┘   │     │
│  └────────────────────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────────────────────┘
```

## SPIFFE/SPIRE - 工作负载身份

### SPIFFE 核心概念

**SPIFFE ID**: 标准化 URI 格式的身份标识

```
格式: spiffe://<trust-domain>/<path>

示例:
- spiffe://company.com/ns/payments/sa/backend
- spiffe://company.com/ns/orders/sa/api
- spiffe://company.com/cluster/prod/node/worker-1
```

**SVID (SPIFFE Verifiable Identity Document)**: 可验证身份文档

- **X.509-SVID**: 标准 TLS 证书，SPIFFE ID 编码在 SAN 字段
- **JWT-SVID**: 用于无法使用 mTLS 的场景
- 短期有效（分钟到小时级），自动轮换

**Workload API**: 本地 Unix Domain Socket 获取身份

### SPIRE 架构

```
┌────────────────────────────────────────────────────────────────┐
│                     SPIRE Server                               │
│  ┌─────────────┐  ┌─────────────┐  ┌───────────────────────┐  │
│  │   CA 插件   │  │  注册条目   │  │   Node 证明器         │  │
│  │(自签/Upstream)│ │   数据库    │  │  (AWS/GCP/Azure/K8s)  │  │
│  └─────────────┘  └─────────────┘  └───────────────────────┘  │
└──────────────────────────┬─────────────────────────────────────┘
                           │ gRPC/mTLS
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  SPIRE Agent │   │  SPIRE Agent │   │  SPIRE Agent │
│  (Node 1)    │   │  (Node 2)    │   │  (Node 3)    │
│ ┌──────────┐ │   │ ┌──────────┐ │   │ ┌──────────┐ │
│ │Workload  │ │   │ │Workload  │ │   │ │Workload  │ │
│ │API (UDS) │ │   │ │API (UDS) │ │   │ │API (UDS) │ │
│ └────┬─────┘ │   │ └────┬─────┘ │   │ └────┬─────┘ │
│      │       │   │      │       │   │      │       │
│ ┌────▼─────┐ │   │ ┌────▼─────┐ │   │ ┌────▼─────┐ │
│ │  Pod A   │ │   │ │  Pod C   │ │   │ │  Pod E   │ │
│ │ (SVID)   │ │   │ │ (SVID)   │ │   │ │ (SVID)   │ │
│ └──────────┘ │   │ └──────────┘ │   │ └──────────┘ │
│ ┌──────────┐ │   │ ┌──────────┐ │   │ ┌──────────┐ │
│ │  Pod B   │ │   │ │  Pod D   │ │   │ │  Pod F   │ │
│ │ (SVID)   │ │   │ │ (SVID)   │ │   │ │ (SVID)   │ │
│ └──────────┘ │   │ └──────────┘ │   │ └──────────┘ │
└──────────────┘   └──────────────┘   └──────────────┘
```

## Kubernetes 集成

### 1. SPIRE 部署

```yaml
# SPIRE Server
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: spire-server
  namespace: spire
spec:
  serviceName: spire-server
  replicas: 1
  selector:
    matchLabels:
      app: spire-server
  template:
    metadata:
      labels:
        app: spire-server
    spec:
      containers:
      - name: spire-server
        image: ghcr.io/spiffe/spire-server:1.9.0
        args:
        - -config
        - /run/spire/config/server.conf
        ports:
        - containerPort: 8081  # gRPC
        volumeMounts:
        - name: spire-config
          mountPath: /run/spire/config
        - name: spire-data
          mountPath: /run/spire/data
---
# SPIRE Agent (DaemonSet)
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: spire-agent
  namespace: spire
spec:
  selector:
    matchLabels:
      app: spire-agent
  template:
    spec:
      hostPID: true
      hostNetwork: true
      dnsPolicy: ClusterFirstWithHostNet
      containers:
      - name: spire-agent
        image: ghcr.io/spiffe/spire-agent:1.9.0
        args:
        - -config
        - /run/spire/config/agent.conf
        volumeMounts:
        - name: spire-config
          mountPath: /run/spire/config
        - name: spire-socket
          mountPath: /run/spire/sockets
        - name: spire-token
          mountPath: /var/run/secrets/tokens
      volumes:
      - name: spire-socket
        hostPath:
          path: /run/spire/sockets
          type: DirectoryOrCreate
```

### 2. 注册工作负载身份

```bash
# 为 payments namespace 的 backend service account 注册身份
spire-server entry create \
  -spiffeID spiffe://company.com/ns/payments/sa/backend \
  -parentID spiffe://company.com/spire/agent/k8s_psat/cluster/my-cluster \
  -selector k8s:ns:payments \
  -selector k8s:sa:backend \
  -ttl 3600  # 1小时有效期

# 为 orders namespace 的 api service account 注册
spire-server entry create \
  -spiffeID spiffe://company.com/ns/orders/sa/api \
  -parentID spiffe://company.com/spire/agent/k8s_psat/cluster/my-cluster \
  -selector k8s:ns:orders \
  -selector k8s:sa:api
```

### 3. 应用 Pod 配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  namespace: payments
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      serviceAccountName: backend
      containers:
      - name: backend
        image: company/backend:v1.0.0
        volumeMounts:
        # SPIRE 自动挂载 SVID
        - name: spiffe-socket
          mountPath: /spiffe-socket
        env:
        - name: SPIFFE_ENDPOINT_SOCKET
          value: unix:///spiffe-socket/socket
        # 使用 SPIFFE 身份连接下游
        - name: DOWNSTREAM_TLS_CERT
          value: /spiffe-socket/svid.pem
        - name: DOWNSTREAM_TLS_KEY
          value: /spiffe-socket/svid.key
        - name: DOWNSTREAM_TLS_CA
          value: /spiffe-socket/bundle.pem
      volumes:
      - name: spiffe-socket
        csi:
          driver: "csi.spiffe.io"
          readOnly: true
```

## 服务网格集成 (Istio Ambient)

### 为什么选择 Ambient Mesh?

传统 Sidecar 模式的问题：

- 每个 Pod 一个 Envoy，资源开销大
- AI/LLM 容器已占用大量内存，再加 Sidecar 压力更大

Ambient Mesh 优势：

- **ztunnel**: 每节点 L4 代理，处理 mTLS
- **waypoint**: 按需 L7 代理
- 内存占用降低 40%+

```yaml
# 启用 Ambient Mesh
apiVersion: v1
kind: Namespace
metadata:
  name: payments
  labels:
    istio.io/dataplane-mode: ambient
    istio.io/use-waypoint: payments-waypoint
---
# 创建 waypoint proxy（仅当需要 L7 功能）
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: payments-waypoint
  namespace: payments
spec:
  gatewayClassName: istio-waypoint
  listeners:
  - name: mesh
    port: 15008
    protocol: HBONE
---
# 授权策略 - 基于 SPIFFE ID
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: backend-policy
  namespace: payments
spec:
  selector:
    matchLabels:
      app: backend
  action: ALLOW
  rules:
  - from:
    - source:
        principals:
        # 只允许 orders namespace 的 api 访问
        - "cluster.local/ns/orders/sa/api"
    to:
    - operation:
        methods: ["GET", "POST"]
        paths: ["/api/v1/*"]
```

## Go 代码示例

### 使用 SPIFFE 身份进行 mTLS 通信

```go
package main

import (
    "context"
    "crypto/tls"
    "fmt"
    "log"
    "net/http"

    "github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
    "github.com/spiffe/go-spiffe/v2/workloadapi"
)

// Server - 使用 SPIFFE 身份的服务端
func runServer(ctx context.Context) error {
    // 创建 Workload API 客户端
    client, err := workloadapi.New(ctx, workloadapi.WithAddr("unix:///spiffe-socket/socket"))
    if err != nil {
        return fmt.Errorf("failed to create workload client: %w", err)
    }
    defer client.Close()

    // 获取 X.509 上下文（包含 SVID 和信任域）
    x509Ctx, err := client.FetchX509Context(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch X509 context: %w", err)
    }

    // 打印自己的 SPIFFE ID
    log.Printf("Server SPIFFE ID: %s", x509Ctx.DefaultSVID().ID)

    // 创建 TLS 服务器
    server := &http.Server{
        Addr: ":8443",
        TLSConfig: tlsconfig.MTLSServerConfig(
            x509Ctx.Source,           // SVID 源
            x509Ctx,                  // 信任域 bundle
            tlsconfig.AuthorizeAny(), // 可改为 AuthorizeID 限制特定身份
        ),
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // 获取客户端 SPIFFE ID
        if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
            // 从证书中提取 SPIFFE ID
            log.Printf("Request from: %s", r.TLS.PeerCertificates[0].URIs)
        }
        fmt.Fprintf(w, "Hello from SPIFFE server!")
    })

    log.Println("Server listening on :8443")
    return server.ListenAndServeTLS("", "")
}

// Client - 使用 SPIFFE 身份的客户端
func callServer(ctx context.Context) error {
    client, err := workloadapi.New(ctx, workloadapi.WithAddr("unix:///spiffe-socket/socket"))
    if err != nil {
        return fmt.Errorf("failed to create workload client: %w", err)
    }
    defer client.Close()

    x509Ctx, err := client.FetchX509Context(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch X509 context: %w", err)
    }

    // 创建 mTLS HTTP 客户端
    httpClient := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsconfig.MTLSClientConfig(
                x509Ctx.Source,
                x509Ctx,
                tlsconfig.AuthorizeID("spiffe://company.com/ns/payments/sa/backend"),
            ),
        },
    }

    resp, err := httpClient.Get("https://backend.payments.svc.cluster.local:8443/")
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    log.Printf("Response status: %s", resp.Status)
    return nil
}

func main() {
    ctx := context.Background()

    // 服务器端
    go func() {
        if err := runServer(ctx); err != nil {
            log.Fatal(err)
        }
    }()

    // 客户端调用
    if err := callServer(ctx); err != nil {
        log.Printf("Client error: %v", err)
    }
}
```

## 零信任实施路径

### 阶段一：身份基础 (1-2 月)

1. 部署 SPIRE Server 和 Agent
2. 为关键服务注册 SPIFFE 身份
3. 验证 SVID 签发和轮换

```bash
# 验证 SVID
kubectl exec -n payments deploy/backend -- \
  openssl x509 -in /spiffe-socket/svid.pem -noout -text | grep URI

# 输出: URI:spiffe://company.com/ns/payments/sa/backend
```

### 阶段二：mTLS 全覆盖 (2-3 月)

1. 启用 Istio Ambient Mesh
2. 配置全局 mTLS (PERMISSIVE -> STRICT)
3. 逐步迁移服务

```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: istio-system
spec:
  mtls:
    mode: STRICT  # 强制 mTLS
```

### 阶段三：细粒度授权 (3-4 月)

1. 基于 SPIFFE ID 的授权策略
2. 实现最小权限原则
3. 持续审计和调整

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: fine-grained
  namespace: production
spec:
  action: ALLOW
  rules:
  # 只允许特定服务访问特定接口
  - from:
    - source:
        principals: ["cluster.local/ns/web/sa/frontend"]
    to:
    - operation:
        methods: ["GET"]
        paths: ["/public/*"]
  - from:
    - source:
        principals: ["cluster.local/ns/api/sa/gateway"]
    to:
    - operation:
        methods: ["POST"]
        paths: ["/api/v1/orders"]
    when:
    - key: request.auth.claims[scope]
      values: ["orders:write"]
```

### 阶段四：持续验证 (持续)

1. 部署监控和审计
2. 异常行为检测
3. 自动响应和隔离

```yaml
# 异常检测策略
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: security-metrics
  namespace: istio-system
spec:
  metrics:
  - providers:
    - name: prometheus
    overrides:
    # 监控 mTLS 握手失败
    - match:
        metric: CONNECTIONS_CLOSED
      tagOverrides:
        mtls:
          value: "true"
    # 监控授权拒绝
    - match:
        metric: REQUEST_COUNT
      tagOverrides:
        response_code:
          value: "403"
```

## 总结

| 组件 | 功能 | 2025 趋势 |
|------|------|-----------|
| SPIFFE/SPIRE | 工作负载身份 | 成为 K8s 标准，集成到更多发行版 |
| Istio Ambient | 无 Sidecar 服务网格 | 降低 AI/大模型工作负载开销 |
| OPA/Gatekeeper | 策略即代码 | 与 SPIFFE 集成更紧密 |
| cert-manager | 证书自动轮换 | 支持 SPIFFE 颁发 |

零信任不是一次性项目，而是持续的安全文化转变。通过 SPIFFE 自动身份管理和 Istio Ambient 的无侵入 mTLS，2025 年的 Kubernetes 零信任实施比以往任何时候都更加简单。
