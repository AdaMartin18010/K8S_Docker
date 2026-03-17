# 微服务演示项目

这是一个完整的微服务演示项目，展示 Docker 和 Kubernetes 的最佳实践。

## 架构

```
┌─────────────────┐     ┌─────────────────┐
│   API Gateway   │────▶│  User Service   │
│   (Nginx/Envoy) │     │   (Go/Gin)      │
└────────┬────────┘     └─────────────────┘
         │
         │             ┌─────────────────┐
         └────────────▶│  Order Service  │
                       │   (Go/Gin)      │
                       └─────────────────┘
```

## 项目结构

```
microservices-demo/
├── user-service/           # 用户服务
│   ├── main.go            # 服务入口
│   ├── Dockerfile         # 多阶段构建
│   ├── go.mod             # Go 模块
│   └── k8s/               # Kubernetes 配置
│       ├── deployment.yaml
│       └── service.yaml
├── order-service/          # 订单服务（类似结构）
└── gateway/                # API 网关
    └── nginx.conf
```

## 快速开始

### 本地开发

```bash
# 用户服务
cd user-service
go mod download
go run main.go

# 测试
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users
```

### Docker 构建

```bash
cd user-service
docker build -t user-service:v1.0.0 .
docker run -p 8080:8080 user-service:v1.0.0
```

### Kubernetes 部署

```bash
# 创建命名空间
kubectl create namespace microservices

# 部署服务
kubectl apply -f user-service/k8s/

# 查看状态
kubectl get pods -n microservices
kubectl get svc -n microservices
```

## 特性

- ✅ 多阶段 Docker 构建
- ✅ 非 root 用户运行
- ✅ 健康检查和就绪检查
- ✅ 优雅关闭
- ✅ 结构化日志
- ✅ Kubernetes 部署配置
