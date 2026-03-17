# Helm 速查表

> Kubernetes 包管理工具快速参考 (Helm 3.15+)

---

## 基础命令

```bash
# 版本信息
helm version

# 仓库管理
helm repo add stable https://charts.helm.sh/stable
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm repo list
helm repo remove bitnami

# 搜索 chart
helm search repo nginx
helm search repo nginx --versions
helm search hub nginx
```

---

## 安装与卸载

```bash
# 安装
helm install my-release bitnami/nginx
helm install my-release ./mychart
helm install my-release ./mychart -f values.yaml
helm install my-release ./mychart --set image.tag=v2.0
helm install my-release ./mychart --dry-run --debug

# 升级
helm upgrade my-release bitnami/nginx
helm upgrade my-release ./mychart -f values.yaml
helm upgrade --install my-release ./mychart  # 如果不存在则安装

# 回滚
helm rollback my-release 1
helm history my-release

# 卸载
helm uninstall my-release
helm uninstall my-release --keep-history

# 查看状态
helm status my-release
helm get values my-release
helm get manifest my-release
```

---

## Chart 开发

```bash
# 创建 chart
helm create mychart

# 验证
helm lint mychart
helm template mychart
helm template mychart --debug
helm template mychart -f custom-values.yaml

# 打包
helm package mychart
helm package mychart --version 1.0.0
helm package mychart --sign --key mykey

# 依赖
helm dependency update mychart
helm dependency build mychart
helm dependency list mychart
```

---

## 模板语法

### 基础

```yaml
# 值引用
{{ .Values.replicaCount }}
{{ .Values.image.repository }}:{{ .Values.image.tag }}

# 默认值
{{ .Values.replicaCount | default 1 }}
{{ .Values.enabled | default true }}

# 变量
{{- $name := "myapp" -}}
{{ $name }}

# 条件
{{- if .Values.enabled }}
enabled: true
{{- end }}

{{- if eq .Values.environment "prod" }}
replicas: 10
{{- else }}
replicas: 2
{{- end }}

# 循环
{{- range .Values.services }}
- name: {{ .name }}
  port: {{ .port }}
{{- end }}
```

### 内置对象

```yaml
{{ .Release.Name }}
{{ .Release.Namespace }}
{{ .Release.IsUpgrade }}
{{ .Release.IsInstall }}
{{ .Release.Revision }}

{{ .Chart.Name }}
{{ .Chart.Version }}
{{ .Chart.AppVersion }}

{{ .Capabilities.KubeVersion.Version }}
{{ .Capabilities.HelmVersion.Version }}
```

### 函数

```yaml
# 字符串
{{ upper .Values.name }}
{{ lower .Values.name }}
{{ quote .Values.name }}
{{ trunc 63 .Values.name }}
{{ trimSuffix "-" .Values.name }}
{{ randAlphaNum 10 }}

# 数学
{{ add 1 2 }}
{{ mul 2 3 }}
{{ div 10 2 }}

# 集合
{{ first .Values.list }}
{{ last .Values.list }}
{{ len .Values.list }}
{{ has "value" .Values.list }}

# 编码
{{ b64enc .Values.secret }}
{{ b64dec .Values.encoded }}
{{ sha256sum .Values.password }}
```

---

## Chart 结构

```
mychart/
├── Chart.yaml          # Chart 元数据
├── values.yaml         # 默认值
├── values.schema.json  # 值验证模式
├── charts/             # 依赖 chart
│   └── dependency-1.0.0.tgz
├── templates/          # 模板文件
│   ├── _helpers.tpl    # 辅助函数
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── hpa.yaml
│   ├── serviceaccount.yaml
│   ├── NOTES.txt       # 安装后说明
│   └── tests/          # 测试
│       └── test-connection.yaml
└── README.md
```

---

## Chart.yaml 示例

```yaml
apiVersion: v2
name: myapp
description: A Helm chart for Kubernetes
type: application
version: 1.2.3
appVersion: "2.0.0"
kubeVersion: ">=1.25.0-0"
home: https://example.com
sources:
  - https://github.com/example/myapp
maintainers:
  - name: John Doe
    email: john@example.com
keywords:
  - web
  - application
dependencies:
  - name: postgresql
    version: 12.x.x
    repository: https://charts.bitnami.com/bitnami
    condition: postgresql.enabled
    tags:
      - database
annotations:
  category: Application
```

---

## Hooks

```yaml
# 预安装
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Release.Name }}-preinstall
  annotations:
    "helm.sh/hook": pre-install,pre-upgrade
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: preinstall
        image: busybox
        command: ['sh', '-c', 'echo pre-install hook']
```

| Hook | 说明 |
|------|------|
| pre-install | 安装前执行 |
| post-install | 安装后执行 |
| pre-delete | 删除前执行 |
| post-delete | 删除后执行 |
| pre-upgrade | 升级前执行 |
| post-upgrade | 升级后执行 |
| pre-rollback | 回滚前执行 |
| post-rollback | 回滚后执行 |

---

## 测试

```bash
# 运行测试
helm test my-release

# 查看测试结果
kubectl logs my-release-test-connection
```

```yaml
# templates/tests/test-connection.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "mychart.fullname" . }}-test-connection"
  annotations:
    "helm.sh/hook": test
spec:
  containers:
    - name: wget
      image: busybox
      command: ['wget']
      args: ['{{ include "mychart.fullname" . }}:{{ .Values.service.port }}']
  restartPolicy: Never
```
