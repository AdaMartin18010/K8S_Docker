# CI/CD 示例

> 持续集成/持续交付流水线配置

---

## GitHub Actions

### 可复用工作流

- [Go 构建测试](./github-actions/reusable-workflows/go-build-test.yaml)
- [Docker 构建推送](./github-actions/reusable-workflows/docker-build-push.yaml)
- [K8s 部署](./github-actions/reusable-workflows/k8s-deploy.yaml)

### 使用示例

```yaml
# .github/workflows/main.yml
name: CI/CD

on:
  push:
    branches: [main]

jobs:
  build:
    uses: ./.github/workflows/go-build-test.yaml
    with:
      go-version: '1.23'

  docker:
    needs: build
    uses: ./.github/workflows/docker-build-push.yaml
    with:
      image-name: myapp
      platforms: linux/amd64,linux/arm64

  deploy:
    needs: docker
    uses: ./.github/workflows/k8s-deploy.yaml
    with:
      environment: production
      namespace: default
      image: ghcr.io/org/myapp
      tag: ${{ github.sha }}
```

---

## GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - build
  - test
  - deploy

variables:
  DOCKER_IMAGE: $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

build:
  stage: build
  script:
    - docker build -t $DOCKER_IMAGE .
    - docker push $DOCKER_IMAGE

test:
  stage: test
  script:
    - go test ./...

deploy:
  stage: deploy
  script:
    - kubectl set image deployment/myapp app=$DOCKER_IMAGE
  only:
    - main
```

---

## 相关文档

- [CI/CD 指南](../../docs/06-practices/cicd-guide.md)
