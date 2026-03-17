# Multi-Stage Docker Build Makefile
.PHONY: all build dev test security production debug clean

# 变量定义
IMAGE_NAME := multistage-app
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# 默认目标
all: test build production

# 开发构建 - 包含热重载
dev:
	docker build \
		--target dev \
		-t $(IMAGE_NAME):dev \
		.

# 运行测试
test:
	docker build \
		--target tester \
		--output type=local,dest=./test-results \
		-t $(IMAGE_NAME):test \
		.

# 安全扫描
security:
	docker build \
		--target security \
		-t $(IMAGE_NAME):security \
		.

# 生产构建
production:
	docker build \
		--target production \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t $(IMAGE_NAME):$(VERSION) \
		-t $(IMAGE_NAME):latest \
		.

# 调试镜像
debug:
	docker build \
		--target debug \
		--build-arg VERSION=$(VERSION) \
		-t $(IMAGE_NAME):debug \
		.

# 查看各阶段大小
sizes:
	@echo "=== 各阶段镜像大小 ==="
	@docker images $(IMAGE_NAME) --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}"

# 对比构建时间
timing:
	@echo "=== 对比构建时间 ==="
	@time docker build --target production -t $(IMAGE_NAME):timed .

# 清理
clean:
	docker rmi -f $(IMAGE_NAME):dev $(IMAGE_NAME):test $(IMAGE_NAME):security \
		$(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest $(IMAGE_NAME):debug 2>/dev/null || true
	rm -rf test-results/

# 运行生产容器
run:
	docker run -d -p 8080:8080 --name $(IMAGE_NAME)-prod $(IMAGE_NAME):latest

# 运行调试容器
run-debug:
	docker run -it --rm -p 8080:8080 --name $(IMAGE_NAME)-debug $(IMAGE_NAME):debug

# 停止容器
stop:
	docker stop $(IMAGE_NAME)-prod $(IMAGE_NAME)-debug 2>/dev/null || true
	docker rm $(IMAGE_NAME)-prod $(IMAGE_NAME)-debug 2>/dev/null || true
