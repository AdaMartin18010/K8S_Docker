# 供应链安全

> 保护软件交付全流程 - 2025 最新实践

---

## 供应链攻击现状

```
┌─────────────────────────────────────────────────────────────┐
│              软件供应链攻击面                                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  源代码          构建           依赖           部署          │
│  ────────────────────────────────────────────────────────   │
│    │              │              │              │           │
│    ▼              ▼              ▼              ▼           │
│  恶意提交      篡改 CI/CD    恶意包         镜像篡改         │
│  泄露密钥      构建环境      漏洞依赖       未授权部署       │
│  代码注入      供应链投毒    混淆依赖       配置漂移         │
│                                                              │
│  典型案例:                                                   │
│  • SolarWinds (2020) - 构建环境被入侵                        │
│  • CodeCov (2021) - Bash Uploader 脚本被篡改                 │
│  • Log4j (2021) - 漏洞依赖                                   │
│  • xz Utils (2024) - 后门植入尝试                            │
│  • 恶意 npm/PyPI (2025) - 包仓库投毒                         │
│                                                              │
│  2025 年数据:                                                │
│  • 供应链攻击占企业安全事件的 45%                            │
│  • 仅 1/5 的组织对开源组件有完整可见性                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心概念

### SBOM (Software Bill of Materials)

软件物料清单，记录软件所有组件的机器可读清单。

```json
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.6",
  "components": [
    {
      "type": "library",
      "name": "lodash",
      "version": "4.17.21",
      "purl": "pkg:npm/lodash@4.17.21",
      "licenses": [{"license": {"id": "MIT"}}],
      "hashes": [
        {
          "alg": "SHA-256",
          "content": "..."
        }
      ]
    },
    {
      "type": "container",
      "name": "myapp",
      "version": "1.0.0",
      "purl": "pkg:docker/myapp@1.0.0"
    }
  ]
}
```

**2025 要求**: CISA 要求关键软件必须提供 SBOM

### SLSA (Supply-chain Levels for Software Artifacts)

软件供应链安全等级框架。

| 等级 | 描述 | 要求 |
|------|------|------|
| **L1** | 来源证明 | 自动化构建，来源元数据 |
| **L2** | 签名来源 | 使用签名，版本控制 |
| **L3** | 强化构建 | 隔离、hermetic 构建 |
| **L4** | 最高等级 | 可复现构建，双人审查 |

**SLSA 1.2 (2025)**: 新增 AI 生成内容的追踪支持

### Sigstore

开源软件签名和验证生态系统。

- **Fulcio**: OIDC 证书颁发机构
- **Rekor**: 透明日志
- **Cosign**: 容器镜像签名工具

**2025 集成**: Sigstore 已集成到 NPM、PyPI、Maven、GitHub、brew、Kubernetes

---

## Cosign 实战 (2025 最新)

### 安装

```bash
# macOS
brew install cosign

# Linux
curl -O -L https://github.com/sigstore/cosign/releases/latest/download/cosign-linux-amd64
chmod +x cosign-linux-amd64
sudo mv cosign-linux-amd64 /usr/local/bin/cosign
```

### 无密钥签名 (Keyless Signing) - 推荐

```bash
# 使用 OIDC 身份签名
COSIGN_EXPERIMENTAL=1 cosign sign --yes myregistry.io/myapp:v1.0.0

# 验证签名
cosign verify myregistry.io/myapp:v1.0.0 \
  --certificate-identity=user@example.com \
  --certificate-oidc-issuer=https://accounts.google.com

# 验证并输出证明
cosign verify myregistry.io/myapp:v1.0.0 \
  --certificate-identity=user@example.com \
  --certificate-oidc-issuer=https://accounts.google.com \
  --output text
```

### 使用密钥签名

```bash
# 生成密钥对
cosign generate-key-pair

# 签名镜像
cosign sign --key cosign.key myregistry.io/myapp:v1.0.0

# 验证签名
cosign verify --key cosign.pub myregistry.io/myapp:v1.0.0
```

### 附加 SBOM

```bash
# 生成 SBOM
syft myregistry.io/myapp:v1.0.0 -o cyclonedx-json > sbom.json

# 附加到镜像
cosign attach sbom --sbom sbom.json myregistry.io/myapp:v1.0.0

# 验证 SBOM
cosign verify-attestation \
  --type cyclonedx \
  --key cosign.pub \
  myregistry.io/myapp:v1.0.0
```

### SLSA 证明

```bash
# 生成 SLSA 证明
cosign attest \
  --predicate slsa-provenance.json \
  --type slsaprovenance \
  --key cosign.key \
  myregistry.io/myapp:v1.0.0

# 验证 SLSA 证明
cosign verify-attestation \
  --type slsaprovenance \
  --key cosign.pub \
  myregistry.io/myapp:v1.0.0
```

---

## CI/CD 集成

### GitHub Actions (2025 最佳实践)

```yaml
name: Build and Sign

on: [push]

permissions:
  contents: read
  id-token: write  # 用于 OIDC
  attestations: write  # 用于 SLSA 证明
  packages: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build image
        run: docker build -t ghcr.io/${{ github.repository }}:${{ github.sha }} .

      - name: Install Cosign
        uses: sigstore/cosign-installer@v3

      - name: Sign image (keyless)
        run: |
          cosign sign --yes \
            ghcr.io/${{ github.repository }}:${{ github.sha }}

      - name: Generate SBOM
        uses: anchore/sbom-action@v0
        with:
          image: ghcr.io/${{ github.repository }}:${{ github.sha }}
          format: cyclonedx-json
          output-file: sbom.json

      - name: Attach SBOM
        run: |
          cosign attach sbom --sbom sbom.json \
            ghcr.io/${{ github.repository }}:${{ github.sha }}

      - name: Generate SLSA provenance
        uses: slsa-framework/slsa-github-generator/.github/workflows/generator_container_slsa3.yml@v2.0.0
        with:
          image: ghcr.io/${{ github.repository }}:${{ github.sha }}
```

### 私有注册表配置

```yaml
# 使用 Harhor 作为签名代理
apiVersion: v1
kind: ConfigMap
metadata:
  name: cosign-config
data:
  COSIGN_REPOSITORY: harbor.example.com/signatures
```

---

## K8s 准入控制

### Kyverno 策略 (推荐)

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: verify-image-signature
spec:
  validationFailureAction: Enforce
  rules:
    - name: check-signature
      match:
        resources:
          kinds:
            - Pod
      verifyImages:
        - imageReferences:
            - "ghcr.io/myorg/*"
            - "myregistry.io/*"
          required: true
          mutateDigest: true
          verifyDigest: true
          attestors:
            - entries:
                # Keyless 验证 (GitHub Actions)
                - keyless:
                    issuer: https://token.actions.githubusercontent.com
                    subject: https://github.com/myorg/myrepo/.github/workflows/build.yml@refs/heads/main
                # 密钥验证
                - keys:
                    publicKeys: |
                      -----BEGIN PUBLIC KEY-----
                      MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE8n...
                      -----END PUBLIC KEY-----
```

### OPA/Gatekeeper 策略

```yaml
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sVerifiedImage
metadata:
  name: verify-image
spec:
  match:
    kinds:
      - apiGroups: [""]
        kinds: ["Pod"]
  parameters:
    repos:
      - ghcr.io/myorg/*
      - myregistry.io/*
```

---

## 扫描工具

| 工具 | 用途 | 特点 | 2025 状态 |
|------|------|------|----------|
| **Trivy** | 镜像/文件系统扫描 | 全面，支持 SBOM | v0.58+ |
| **Grype** | 漏洞扫描 | Anchore 出品 | v0.86+ |
| **Syft** | SBOM 生成 | 多格式输出 | v1.18+ |
| **Snyk** | 依赖扫描 | 商业工具 | - |
| **SLSA GitHub Generator** | SLSA 证明 | 自动生成 | v2.0+ |
| **GUAC** | 供应链分析 | 开源图谱分析 | v0.13+ |
| **VEX** | 漏洞利用交换 | 减少误报 | 标准化 |

### Trivy 扫描 (2025)

```bash
# 镜像扫描
trivy image myapp:latest

# 生成 SBOM
trivy image --format cyclonedx -o sbom.json myapp:latest

# 扫描 IaC
trivy config ./terraform

# CI 模式 (返回非零退出码)
trivy image --exit-code 1 --severity HIGH,CRITICAL myapp:latest

# 扫描 SBOM
trivy sbom sbom.json

# 扫描 Git 仓库
trivy repo https://github.com/myorg/myrepo
```

---

## 供应链安全清单

### 开发阶段

- [ ] 代码签名 (GPG)
- [ ] 依赖审查 (Dependabot, Snyk, Renovate)
- [ ] 密钥扫描 (TruffleHog, GitLeaks)
- [ ] SAST (SonarQube, CodeQL)
- [ ] 预提交钩子 (pre-commit)

### 构建阶段

- [ ] 隔离构建环境 (隔离 runner)
- [ ] Hermetic 构建 (无网络访问)
- [ ] 签名容器镜像 (Cosign)
- [ ] 生成 SBOM (CycloneDX/SPDX)
- [ ] SLSA 证明 (Level 3+)
- [ ] 依赖锁定 (lock files)

### 部署阶段

- [ ] 镜像签名验证 (Kyverno/OPA)
- [ ] SBOM 验证
- [ ] 漏洞扫描 (Trivy)
- [ ] 运行时安全监控 (Falco/Tetragon)
- [ ] 准入控制策略

---

## 2025 趋势

- ✅ **SLSA 1.2**: 新增 AI 生成内容追踪
- ✅ **VEX 标准化**: 漏洞利用交换格式统一
- ✅ **GUAC**: 开源供应链图谱分析工具成熟
- ✅ **Sigstore 普及**: 集成到主流包管理器
- ✅ **Keyless Signing**: OIDC 无密钥签名成为标准
- 🔄 **in-toto**: 供应链元数据框架推广
- 🔄 **二进制授权**: Google 风格的部署验证
- 🔄 **Reproducible Builds**: 可复现构建实践

---

## 参考

- [SLSA 官方文档](https://slsa.dev/)
- [Sigstore 文档](https://docs.sigstore.dev/)
- [Cosign 文档](https://docs.sigstore.dev/cosign/overview/)
- [GUAC 文档](https://guac.sh/)
- [CISA SBOM 指南](https://www.cisa.gov/sbom)
- [OpenSSF 安全供应链](https://openssf.org/)
