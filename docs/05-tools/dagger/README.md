# Dagger - 可编程 CI/CD 引擎

## 概述

Dagger 是一个可编程的 CI/CD 引擎，允许开发者使用 Go、Python 或 TypeScript 等真实编程语言编写流水线代码，在本地和 CI 环境中以一致的方式运行。

## 核心特性

| 特性 | 描述 |
|------|------|
| 流水线即代码 | 使用真实编程语言而非 YAML |
| 容器化执行 | 基于 BuildKit，所有操作在容器中运行 |
| 即时本地测试 | 本地开发和调试 CI 流水线 |
| 自动缓存 | 内容寻址缓存，自动优化构建速度 |
| 跨平台 | 支持多架构构建（AMD64/ARM64） |

## 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                        Dagger 架构                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │  Pipeline    │    │   Dagger     │    │   Docker     │      │
│  │   Code       │───▶│    Engine    │───▶│  Daemon      │      │
│  │ (Go/Python)  │    │  (BuildKit)  │    │              │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
│                               │                                 │
│                               ▼                                 │
│                    ┌──────────────────────┐                    │
│                    │  DAG (有向无环图)      │                    │
│                    │  - 拉取镜像            │                    │
│                    │  - 执行命令            │                    │
│                    │  - 文件操作            │                    │
│                    │  - 缓存层              │                    │
│                    └──────────────────────┘                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 安装部署

### 安装 Dagger CLI

```bash
# macOS/Linux
curl -fsSL https://dl.dagger.io/dagger/install.sh | sh
sudo mv ./bin/dagger /usr/local/bin/

# 验证安装
dagger version

# 初始化项目
dagger init --sdk=go --name=ci
cd dagger
```

### Go SDK 示例

```go
// dagger/main.go
package main

import (
    "context"
    "fmt"
    "dagger/ci/internal/dagger"
)

type Ci struct{}

// Build 编译 Go 应用
func (m *Ci) Build(ctx context.Context, source *dagger.Directory) *dagger.Container {
    // 使用 Go 基础镜像构建
    builder := dag.Container().
        From("golang:1.22-alpine").
        WithWorkdir("/app").
        WithDirectory("/app", source).
        WithExec([]string{"go", "mod", "download"}).
        WithExec([]string{"go", "build", "-o", "/app/server", "./cmd/server"})

    // 创建最小运行时镜像
    runtime := dag.Container().
        From("alpine:3.19").
        WithFile("/usr/local/bin/server", builder.File("/app/server")).
        WithEntrypoint([]string{"/usr/local/bin/server"})

    return runtime
}

// Test 运行测试套件
func (m *Ci) Test(ctx context.Context, source *dagger.Directory) (string, error) {
    return dag.Container().
        From("golang:1.22-alpine").
        WithWorkdir("/app").
        WithDirectory("/app", source).
        WithExec([]string{"go", "mod", "download"}).
        WithExec([]string{"go", "test", "-v", "-race", "./..."}).
        Stdout(ctx)
}

// Publish 构建并推送镜像
func (m *Ci) Publish(ctx context.Context, source *dagger.Directory, registry string, tag string) (string, error) {
    container := m.Build(ctx, source)
    addr := registry + ":" + tag
    return container.Publish(ctx, addr)
}

// MultiPlatform 多平台构建
func (m *Ci) MultiPlatform(ctx context.Context, source *dagger.Directory) (string, error) {
    platforms := []dagger.Platform{
        "linux/amd64",
        "linux/arm64",
    }

    platformVariants := make([]*dagger.Container, len(platforms))
    for i, platform := range platforms {
        platformVariants[i] = dag.Container(dagger.ContainerOpts{Platform: platform}).
            From("golang:1.22-alpine").
            WithDirectory("/src", source).
            WithWorkdir("/src").
            WithEnvVariable("CGO_ENABLED", "0").
            WithExec([]string{"go", "build", "-o", "/app", "./cmd/server"})
    }

    return dag.Container().
        Publish(ctx, "ghcr.io/my-org/my-app:latest",
            dagger.ContainerPublishOpts{
                PlatformVariants: platformVariants,
            })
}
```

### Python SDK 示例

```python
# dagger/src/main/__init__.py
import dagger
from dagger import dag, function, object_type

@object_type
class MyPipeline:
    @function
    async def test(self, source: dagger.Directory) -> str:
        """运行 Python 测试"""
        return await (
            dag.container()
            .from_("python:3.12-slim")
            .with_directory("/src", source)
            .with_workdir("/src")
            .with_exec(["pip", "install", "-r", "requirements.txt"])
            .with_exec(["pip", "install", "pytest", "pytest-cov"])
            .with_exec(["pytest", "-v", "--cov=app", "tests/"])
            .stdout()
        )

    @function
    async def build(self, source: dagger.Directory) -> dagger.Container:
        """构建 Docker 镜像"""
        pip_cache = dag.cache_volume("pip-cache")

        return (
            dag.container()
            .from_("python:3.12-slim")
            .with_mounted_cache("/root/.cache/pip", pip_cache)
            .with_directory("/app", source)
            .with_workdir("/app")
            .with_exec(["pip", "install", "--no-cache-dir", "-r", "requirements.txt"])
            .with_entrypoint(["python", "main.py"])
        )

    @function
    async def publish(self, source: dagger.Directory, registry: str = "ghcr.io",
                     image_name: str = "my-app", tag: str = "latest",
                     registry_user: dagger.Secret | None = None,
                     registry_pass: dagger.Secret | None = None) -> str:
        """推送镜像到仓库"""
        container = await self.build(source)

        if registry_user and registry_pass:
            container = container.with_registry_auth(
                registry,
                await registry_user.plaintext(),
                registry_pass,
            )

        ref = f"{registry}/{image_name}:{tag}"
        digest = await container.publish(ref)
        return digest
```

## CI/CD 集成

### GitHub Actions

```yaml
# .github/workflows/dagger.yml
name: Dagger CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Run Dagger
      uses: dagger/dagger-for-github@v6
      with:
        version: "latest"
        cmds: |
          call test --source=.
          call build --source=.
```

### GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - test
  - build

variables:
  DAGGER_VERSION: "0.14.0"

test:
  stage: test
  image: docker:latest
  services:
    - docker:dind
  script:
    - apk add curl
    - curl -fsSL https://dl.dagger.io/dagger/install.sh | sh
    - ./bin/dagger call test --source=.

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - apk add curl
    - curl -fsSL https://dl.dagger.io/dagger/install.sh | sh
    - ./bin/dagger call publish --source=. --registry=$CI_REGISTRY --tag=$CI_COMMIT_SHA
```

## 缓存策略

```go
// 使用缓存卷
func (m *Ci) BuildWithCache(source *dagger.Directory) *dagger.Container {
    goCache := dag.CacheVolume("go-mod-cache")

    return dag.Container().
        From("golang:1.22-alpine").
        WithMountedCache("/go/pkg/mod", goCache).
        WithWorkdir("/app").
        WithDirectory("/app", source).
        WithExec([]string{"go", "mod", "download"}).
        WithExec([]string{"go", "build", "-o", "server"})
}
```

## Secret 管理

```go
// 安全使用 Secret
func (m *Ci) Deploy(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
    return dag.Container().
        From("alpine:3.19").
        WithWorkdir("/app").
        WithDirectory("/app", source).
        WithSecretVariable("DEPLOY_TOKEN", token).
        WithExec([]string{"sh", "-c", "deploy.sh"}).
        Stdout(ctx)
}

// 调用时传递
dagger call deploy --source=. --token=env:DEPLOY_TOKEN
```

## 2025 新特性

- **Daggerverse**: 模块生态系统，共享和重用流水线组件
- **AI 集成**: 与 LLM 集成，智能修复失败的流水线
- **自修复流水线**: 自动检测和修复构建问题
- **增强调试**: 交互式调试模式，逐步执行流水线

## 相关资源

- [Dagger 官网](https://dagger.io/)
- [Dagger 文档](https://docs.dagger.io/)
- [Daggerverse](https://daggerverse.dev/)
