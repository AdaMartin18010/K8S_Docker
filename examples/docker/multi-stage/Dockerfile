# ============================================================================
# 多阶段构建高级示例
# 目标：最小化最终镜像大小，最大化构建缓存利用
# ============================================================================

# syntax=docker/dockerfile:1.6

# ----------------------------------------------------------------------------
# 阶段 1: 依赖下载 (Dependencies)
# 目的：利用缓存避免重复下载依赖
# ----------------------------------------------------------------------------
FROM golang:1.22-alpine3.19 AS deps

WORKDIR /src

# 仅复制依赖文件，最大化缓存命中率
COPY go.mod go.sum ./

# 使用 BuildKit 缓存挂载加速下载
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && \
    go mod verify

# ----------------------------------------------------------------------------
# 阶段 2: 代码生成 (Code Generation)
# 目的：生成 protobuf、mock 等代码
# ----------------------------------------------------------------------------
FROM golang:1.22-alpine3.19 AS codegen

WORKDIR /src

# 安装代码生成工具
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/golang/mock/mockgen@latest

COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .

# 生成代码
RUN go generate ./...

# ----------------------------------------------------------------------------
# 阶段 3: 开发构建 (Development)
# 目的：包含调试工具，用于开发环境
# ----------------------------------------------------------------------------
FROM golang:1.22-alpine3.19 AS dev

WORKDIR /app

# 安装开发工具
RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .

# 开发模式使用 air 热重载
CMD ["air", "-c", ".air.toml"]

# ----------------------------------------------------------------------------
# 阶段 4: 测试 (Testing)
# 目的：运行单元测试、集成测试、代码覆盖率分析
# ----------------------------------------------------------------------------
FROM golang:1.22-alpine3.19 AS tester

WORKDIR /src

# 安装测试工具
RUN go install github.com/axw/gocov/gocov@latest && \
    go install github.com/AlekSi/gocov-xml@latest

COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY --from=codegen /src /src

# 运行测试并生成覆盖率报告
RUN --mount=type=cache,target=/go/pkg/mod \
    go test -v -race -coverprofile=coverage.out ./... && \
    go tool cover -html=coverage.out -o coverage.html

# ----------------------------------------------------------------------------
# 阶段 5: 构建 (Building)
# 目的：编译生产就绪的二进制文件
# ----------------------------------------------------------------------------
FROM golang:1.22-alpine3.19 AS builder

WORKDIR /src

# 安装编译依赖
RUN apk add --no-cache git ca-certificates

COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY --from=codegen /src /src

# 编译优化选项：
# -w: 去除 DWARF 调试信息
# -s: 去除符号表
# -extldflags "-static": 静态链接
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s \
    -X main.Version=${VERSION} \
    -X main.BuildTime=${BUILD_TIME} \
    -X main.GitCommit=${GIT_COMMIT} \
    -extldflags '-static'" \
    -a -installsuffix cgo \
    -o /bin/server \
    ./cmd/server

# 验证二进制文件
RUN /bin/server --version || true

# ----------------------------------------------------------------------------
# 阶段 6: 安全扫描 (Security)
# 目的：扫描依赖和代码漏洞
# ----------------------------------------------------------------------------
FROM aquasec/trivy:latest AS security

COPY --from=builder /bin/server /app/server

# 扫描二进制文件
RUN trivy filesystem --exit-code 0 --severity HIGH,CRITICAL /app/server

# ----------------------------------------------------------------------------
# 阶段 7: 生产镜像 (Production)
# 目的：最小化、最安全的最终镜像
# ----------------------------------------------------------------------------
FROM scratch AS production

# 元数据
LABEL org.opencontainers.image.title="Multi-Stage App" \
      org.opencontainers.image.description="Optimized multi-stage build" \
      org.opencontainers.image.version="${VERSION}"

# 从 builder 复制必要文件
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/server /server

# 创建非 root 用户（在 scratch 中手动创建）
COPY --from=builder /etc/passwd /etc/passwd
USER 65534:65534

EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/server", "health"]

ENTRYPOINT ["/server"]

# ----------------------------------------------------------------------------
# 阶段 8: 调试镜像 (Debug)
# 目的：包含调试工具，用于生产环境排障
# ----------------------------------------------------------------------------
FROM alpine:3.19 AS debug

RUN apk add --no-cache ca-certificates curl net-tools bind-tools strace

COPY --from=builder /bin/server /server

# 使用 root 运行以便调试
ENTRYPOINT ["/server"]
