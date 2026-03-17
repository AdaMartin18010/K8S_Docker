# Examples Makefile - 统一管理所有示例
.PHONY: help docker k8s test clean lint

# 默认目标
.DEFAULT_GOAL := help

help: ## 显示帮助信息
	@echo "Docker & Kubernetes 示例代码库"
	@echo ""
	@echo "使用方法:"
	@echo "  make docker-build    - 构建所有 Docker 示例"
	@echo "  make docker-test     - 测试 Docker 示例"
	@echo "  make k8s-apply       - 应用所有 K8s 配置"
	@echo "  make k8s-delete      - 删除所有 K8s 资源"
	@echo "  make go-test         - 运行 Go 测试"
	@echo "  make lint            - 运行代码检查"
	@echo "  make clean           - 清理所有资源"

# =============================================================================
# Docker 相关
# =============================================================================

docker-build: ## 构建所有 Docker 示例
	@echo "=== 构建 Docker 基础示例 ==="
	cd docker/basic && docker build -f Dockerfile.good -t example/app:good .
	@echo "=== 构建多阶段示例 ==="
	cd docker/multi-stage && docker build --target production -t example/multi-stage:latest .
	@echo "=== 构建安全示例 ==="
	cd docker/security && docker build -f Dockerfile.secure -t example/secure:latest .
	@echo "=== 构建微服务示例 ==="
	cd microservices-demo/user-service && docker build -t example/user-service:v1.0.0 .

docker-test: ## 测试 Docker 示例
	@echo "=== 测试 Docker 基础示例 ==="
	docker run -d --name test-app -p 8080:8080 example/app:good
	@sleep 2
	curl -s http://localhost:8080/health | jq .
	docker stop test-app && docker rm test-app

docker-clean: ## 清理 Docker 资源
	@echo "=== 清理 Docker 容器和镜像 ==="
	-docker stop $$(docker ps -a -q --filter "name=example") 2>/dev/null
	-docker rm $$(docker ps -a -q --filter "name=example") 2>/dev/null
	-docker rmi $$(docker images -q example/*) 2>/dev/null
	-docker system prune -f

# =============================================================================
# Kubernetes 相关
# =============================================================================

k8s-apply: ## 应用所有 K8s 配置
	@echo "=== 创建命名空间 ==="
	kubectl create namespace examples --dry-run=client -o yaml | kubectl apply -f -
	@echo "=== 应用基础资源 ==="
	kubectl apply -f kubernetes/01-basic-resources/ -n examples
	@echo "=== 应用安全配置 ==="
	kubectl apply -f kubernetes/06-security/ -n examples

k8s-delete: ## 删除所有 K8s 资源
	@echo "=== 删除资源 ==="
	kubectl delete namespace examples --ignore-not-found=true

k8s-status: ## 查看 K8s 资源状态
	@echo "=== Pods ==="
	kubectl get pods -n examples
	@echo "=== Services ==="
	kubectl get svc -n examples
	@echo "=== Deployments ==="
	kubectl get deployments -n examples

# =============================================================================
# Go 相关
# =============================================================================

go-test: ## 运行 Go 测试
	@echo "=== 测试 client-go 示例 ==="
	cd go-client/01-basic-ops && go build -o /dev/null .
	cd go-client/02-controller && go build -o /dev/null .
	@echo "=== 测试微服务 ==="
	cd microservices-demo/user-service && go build -o /dev/null .

go-mod-download: ## 下载所有 Go 模块
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "=== Downloading $$dir ==="; \
		cd "$$dir" && go mod download; \
	done

# =============================================================================
# 代码检查
# =============================================================================

lint: lint-docker lint-k8s ## 运行所有代码检查

lint-docker: ## 检查 Dockerfile
	@echo "=== 检查 Dockerfile ==="
	@which hadolint >/dev/null 2>&1 || (echo "安装 hadolint..." && \
		curl -sL -o /tmp/hadolint \
		https://github.com/hadolint/hadolint/releases/download/v2.12.0/hadolint-Linux-x86_64 && \
		chmod +x /tmp/hadolint)
	@find . -name "Dockerfile*" -exec /tmp/hadolint {} \; 2>/dev/null || true

lint-k8s: ## 检查 Kubernetes 配置
	@echo "=== 检查 Kubernetes 配置 ==="
	@which kube-linter >/dev/null 2>&1 || (echo "安装 kube-linter..." && \
		curl -sL https://github.com/stackrox/kube-linter/releases/download/v0.6.4/kube-linter-linux.tar.gz | \
		tar -xz -C /tmp)
	@/tmp/kube-linter lint kubernetes/ 2>/dev/null || true

# =============================================================================
# 清理
# =============================================================================

clean: docker-clean k8s-delete ## 清理所有资源
	@echo "=== 清理完成 ==="

# =============================================================================
# 完整流程
# =============================================================================

all: docker-build go-test ## 构建和测试所有内容
	@echo "=== 全部完成 ==="
