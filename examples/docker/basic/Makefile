# Docker Basic Example Makefile
.PHONY: build run test clean compare

IMAGE_GOOD := myapp:good
IMAGE_BAD := myapp:bad
CONTAINER_NAME := myapp-test

build: ## 构建镜像
	@echo "=== 构建最佳实践镜像 ==="
	DOCKER_BUILDKIT=1 docker build -f Dockerfile.good -t $(IMAGE_GOOD) .
	@echo "=== 构建错误示例镜像 ==="
	DOCKER_BUILDKIT=0 docker build -f Dockerfile.bad -t $(IMAGE_BAD) .

run: ## 运行最佳实践镜像
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p 8080:8080 \
		-e PORT=8080 \
		$(IMAGE_GOOD)

	@echo "等待服务启动..."
	@sleep 3
	@echo "测试健康检查:"
	@curl -s http://localhost:8080/health | jq . || echo "服务未就绪"

test: ## 测试 API
	@echo "=== 测试端点 ==="
	curl -s http://localhost:8080/ | jq .
	curl -s http://localhost:8080/health | jq .
	curl -s http://localhost:8080/env | jq .

compare: ## 对比镜像大小
	@echo "=== 镜像大小对比 ==="
	docker images $(IMAGE_GOOD) $(IMAGE_BAD) --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

history: ## 查看构建历史
	@echo "=== 最佳实践镜像构建历史 ==="
	docker history $(IMAGE_GOOD)
	@echo ""
	@echo "=== 错误示例镜像构建历史 ==="
	docker history $(IMAGE_BAD)

clean: ## 清理资源
	-docker stop $(CONTAINER_NAME) 2>/dev/null
	-docker rm $(CONTAINER_NAME) 2>/dev/null
	-docker rmi $(IMAGE_GOOD) $(IMAGE_BAD) 2>/dev/null

all: build run test ## 完整流程
