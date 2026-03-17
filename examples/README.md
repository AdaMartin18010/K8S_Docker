# Docker & Kubernetes Go 示例代码库

本目录包含 Docker 和 Kubernetes 的完整 Go 代码示例，涵盖最佳实践、常见反例和完整项目案例。

## 📁 目录结构

```
examples/
├── docker/                      # Docker 示例
│   ├── basic/                  # 基础 Dockerfile 示例
│   │   ├── Dockerfile.good     # ✅ 最佳实践
│   │   ├── Dockerfile.bad      # ❌ 常见错误
│   │   ├── main.go             # 示例 Go 应用
│   │   └── README.md
│   ├── multi-stage/            # 多阶段构建示例
│   │   ├── Dockerfile          # 8 阶段构建
│   │   └── Makefile
│   ├── compose/                # Docker Compose 配置
│   │   ├── docker-compose.yml
│   │   ├── docker-compose.override.yml
│   │   └── docker-compose.prod.yml
│   └── security/               # 安全加固示例
│       ├── Dockerfile.secure
│       └── seccomp-profile.json
│
├── kubernetes/                  # Kubernetes 配置示例
│   ├── 01-basic-resources/     # 基础资源
│   │   ├── pod-good.yaml
│   │   ├── pod-bad.yaml
│   │   ├── deployment-good.yaml
│   │   └── service-good.yaml
│   ├── 03-config-management/   # 配置管理
│   │   ├── configmap-good.yaml
│   │   └── secret-good.yaml
│   └── 06-security/            # 安全配置
│       ├── rbac-good.yaml
│       └── network-policy-good.yaml
│
├── go-client/                   # Go client-go 示例
│   ├── 01-basic-ops/           # 基础操作
│   │   ├── main.go             # CRUD 操作
│   │   └── go.mod
│   └── 02-controller/          # 自定义控制器
│       └── main.go
│
├── anti-patterns/               # 反例集合
│   ├── docker/                 # Docker 反例
│   │   ├── Dockerfile.security-risks
│   │   └── Dockerfile.performance-issues
│   └── kubernetes/             # Kubernetes 反例
│       ├── pod.security-risks.yaml
│       └── deployment.anti-patterns.yaml
│
└── microservices-demo/          # 微服务完整案例
    ├── README.md
    └── user-service/           # 用户服务示例
        ├── main.go
        ├── Dockerfile
        ├── go.mod
        └── k8s/
            └── deployment.yaml
```

## 🚀 快速开始

### 环境要求

- Go 1.22+
- Docker 24.0+ (启用 BuildKit)
- Kubernetes 1.28+ (可选)
- kubectl (可选)

### Docker 示例

```bash
# 基础示例
cd docker/basic
docker build -f Dockerfile.good -t myapp:good .
docker build -f Dockerfile.bad -t myapp:bad .

# 对比镜像大小
docker images myapp:good myapp:bad

# 多阶段构建
cd docker/multi-stage
make production
make sizes

# Docker Compose
cd docker/compose
cp .env.example .env
make dev
```

### Kubernetes 示例

```bash
# 应用配置
kubectl apply -f kubernetes/01-basic-resources/pod-good.yaml
kubectl apply -f kubernetes/01-basic-resources/deployment-good.yaml

# 查看状态
kubectl get pods
kubectl get deployments
kubectl describe pod <pod-name>
```

### Go 客户端示例

```bash
# 基础操作
cd go-client/01-basic-ops
go mod download

# 列出 Pod
go run main.go -op list -namespace default

# 创建 Pod
go run main.go -op create -name nginx-test

# 监视 Pod
go run main.go -op watch

# 自定义控制器
cd go-client/02-controller
go run main.go
```

### 微服务演示

```bash
# 构建镜像
cd microservices-demo/user-service
docker build -t user-service:v1.0.0 .

# 本地运行
docker run -p 8080:8080 user-service:v1.0.0

# 测试
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users

# Kubernetes 部署
kubectl create namespace microservices
kubectl apply -f k8s/
```

## 📚 学习路径

### 初学者
1. `docker/basic/` - 理解 Dockerfile 最佳实践
2. `kubernetes/01-basic-resources/` - 学习基础 K8s 资源
3. `microservices-demo/` - 运行完整示例

### 中级
1. `docker/multi-stage/` - 掌握多阶段构建
2. `docker/compose/` - 学习容器编排
3. `kubernetes/06-security/` - 安全配置
4. `go-client/01-basic-ops/` - 编程式操作 K8s

### 高级
1. `go-client/02-controller/` - 自定义控制器开发
2. `anti-patterns/` - 理解常见错误
3. `docker/security/` - 安全加固

## ✅ 最佳实践检查清单

### Docker
- [ ] 使用多阶段构建
- [ ] 最小基础镜像 (scratch/distroless/alpine)
- [ ] 非 root 用户运行
- [ ] 只读根文件系统
- [ ] 资源限制
- [ ] 健康检查
- [ ] .dockerignore 配置
- [ ] 安全扫描

### Kubernetes
- [ ] 资源请求和限制
- [ ] 健康检查（liveness/readiness）
- [ ] 安全上下文（runAsNonRoot, readOnlyRootFilesystem）
- [ ] 网络策略
- [ ] RBAC 最小权限
- [ ] Pod 中断预算
- [ ] Secret 管理
- [ ] 配置分离（ConfigMap）

## ❌ 常见错误

### Dockerfile
1. 使用 `latest` 标签
2. 以 root 用户运行
3. 镜像包含敏感信息
4. 缺少 .dockerignore
5. 单阶段构建导致镜像过大

### Kubernetes
1. 使用 privileged 容器
2. 没有资源限制
3. 硬编码敏感信息
4. 缺少健康检查
5. 使用默认 Service Account

## 🔧 工具推荐

- **构建**: Docker BuildKit, Buildx
- **扫描**: Trivy, Docker Scout, Grype
- **Lint**: hadolint (Dockerfile), kube-linter (K8s)
- **监控**: Prometheus, Grafana
- **开发**: Air (Go 热重载), Telepresence

## 📖 参考资源

- [Docker 最佳实践](https://docs.docker.com/develop/dev-best-practices/)
- [Kubernetes 文档](https://kubernetes.io/docs/)
- [Go client-go 示例](https://github.com/kubernetes/client-go/tree/master/examples)
- [CNCF 云原生全景](https://landscape.cncf.io/)

## 📝 贡献

欢迎提交 Issue 和 PR 来改进示例代码。
