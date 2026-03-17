# K8s 开发工具链

> 云原生应用开发与调试工具 (2025)

---

## 开发环境工具

### 1. kubectl + 插件

```bash
# 安装 krew (kubectl 插件管理器)
(
  set -x; cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  KREW="krew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
  tar zxvf "${KREW}.tar.gz" &&
  ./"${KREW}" install krew
)

# 常用插件
kubectl krew install ctx          # 快速切换 context
kubectl krew install ns           # 快速切换 namespace
kubectl krew install tail         # 多 Pod 日志 tail
kubectl krew install tree         # 资源树视图
kubectl krew install exec-as      # 以特定用户执行
kubectl krew install sniff        # 抓包工具
kubectl krew install view-secret  # 查看 secret 内容
kubectl krew install resource-capacity  # 资源容量视图
```

### 2. Lens / OpenLens

```bash
# 安装 OpenLens (开源版)
# macOS
brew install --cask openlens

# Windows
choco install openlens

# Linux
wget https://github.com/MuhammedKalkan/OpenLens/releases/download/v6.5.2-366/OpenLens-6.5.2-366.x86_64.AppImage
chmod +x OpenLens-6.5.2-366.x86_64.AppImage
./OpenLens-6.5.2-366.x86_64.AppImage
```

### 3. k9s - 终端 UI

```bash
# 安装
brew install k9s  # macOS
choco install k9s  # Windows

# 使用
k9s
k9s -n kube-system  # 指定 namespace
k9s --context prod  # 指定 context

# 快捷键
# :pod, :svc, :deploy - 切换资源类型
# /filter - 过滤
# l - 日志
# s - Shell
# d - 描述
# y - YAML
```

---

## 本地集群工具

### kind (Kubernetes in Docker)

```bash
# 安装
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.27.0/kind-linux-amd64
chmod +x ./kind && sudo mv ./kind /usr/local/bin/

# 创建集群
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
- role: worker
- role: worker
EOF

# 加载本地镜像到 kind
kind load docker-image myapp:latest
```

### minikube

```bash
# 安装
brew install minikube

# 启动（使用 Docker 驱动）
minikube start --driver=docker --cpus=4 --memory=8192

# 启用插件
minikube addons enable ingress
minikube addons enable metrics-server

# 使用本地镜像
eval $(minikube docker-env)
docker build -t myapp:latest .
```

### k3d (K3s in Docker)

```bash
# 安装
brew install k3d

# 创建集群
k3d cluster create mycluster \
  --servers 1 --agents 3 \
  --port "80:80@loadbalancer" \
  --port "443:443@loadbalancer"

# 导入镜像
k3d image import myapp:latest -c mycluster
```

---

## 调试工具

### Telepresence

```bash
# 安装
brew install telepresence

# 连接集群
telepresence connect

# 拦截服务（本地开发调试）
telepresence intercept my-service --port 8080:http

# 卸载拦截
telepresence leave my-service
```

### mirrord

```bash
# 安装
brew install metalbear-co/mirrord/mirrord

# 运行本地应用并镜像流量
mirrord exec --target deployment/my-app ./my-local-app

# 窃取流量（仅部分流量到本地）
mirrord exec --target deployment/my-app --steal ./my-local-app
```

### kubectl-debug

```bash
# 安装
kubectl krew install debug

# 在目标 Pod 中启动调试容器
kubectl debug my-pod -it --image=nicolaka/netshoot --target=my-container

# 复制 Pod 并调试
kubectl debug my-pod -it --copy-to=my-pod-debug --image=busybox
```

---

## CI/CD 工具

### 1. Argo Rollouts CLI

```bash
# 安装
brew install argoproj/tap/kubectl-argo-rollouts

# 查看 rollout 状态
kubectl argo rollouts list rollout
kubectl argo rollouts get rollout my-app
kubectl argo rollouts promote my-app
kubectl argo rollouts undo my-app
```

### 2. Flux CLI

```bash
# 安装
brew install fluxcd/tap/flux

# 引导集群
flux bootstrap github \
  --owner=$GITHUB_USER \
  --repository=fleet-infra \
  --branch=main \
  --path=./clusters/my-cluster \
  --personal

# 查看状态
flux get all
flux logs
```

---

## 代码生成工具

### kubebuilder

```bash
# 安装
brew install kubebuilder

# 初始化项目
kubebuilder init --domain example.com --repo github.com/example/my-operator

# 创建 API
kubebuilder create api --group cache --version v1 --kind Memcached

# 生成 manifests
make generate
make manifests

# 本地运行
make run

# 构建镜像
make docker-build IMG=myregistry/my-operator:v1
make deploy IMG=myregistry/my-operator:v1
```

### operator-sdk

```bash
# 安装
brew install operator-sdk

# 初始化
operator-sdk init --domain example.com --repo github.com/example/memcached-operator

# 创建 API
operator-sdk create api --group cache --version v1 --kind Memcached --resource --controller

# 添加组件
operator-sdk add component --helm-chart mariadb
```

---

## 最佳实践

1. **使用 ctx/ns 插件**: 快速切换集群和命名空间
2. **本地集群选 kind**: 轻量、快速、适合 CI
3. **生产调试用 mirrord**: 更现代的流量镜像方案
4. **代码生成用 kubebuilder**: 标准 Operator 框架
