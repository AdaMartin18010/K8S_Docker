# 供应链安全

> 保护软件交付全流程

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
│  • xz Utils (2024) - 后门植入尝试                            │
│  • 3CX (2023) - 软件更新被投毒                               │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心概念

### SBOM (Software Bill of Materials)

软件物料清单，记录软件所有组件的机器可读清单。

```json
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.5",
  "components": [
    {
      "type": "library",
      "name": "lodash",
      "version": "4.17.21",
      "purl": "pkg:npm/lodash@4.17.21",
      "licenses": [{"license": {"id": "MIT"}}]
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

### SLSA (Supply-chain Levels for Software Artifacts)

软件供应链安全等级框架。

| 等级 | 描述 | 要求 |
|------|------|------|
| **L1** | 来源证明 | 自动化构建，来源元数据 |
| **L2** | 签名来源 | 使用签名，版本控制 |
| **L3** | 强化构建 | 隔离、 hermetic 构建 |
| **L4** | 最高等级 | 可复现构建，双人审查 |

### Sigstore

开源软件签名和验证生态系统。

- **Fulcio**: OIDC 证书颁发机构
- **Rekor**: 透明日志
- **Cosign**: 容器镜像签名工具

---

## Cosign 实战

### 安装

```bash
# macOS
brew install cosign

# Linux
curl -O -L https://github.com/sigstore/cosign/releases/latest/download/cosign-linux-amd64
chmod +x cosign-linux-amd64
sudo mv cosign-linux-amd64 /usr/local/bin/cosign
```

### 无密钥签名 (Keyless Signing)

```bash
# 使用 OIDC 身份签名 (推荐)
cosign sign --yes myregistry.io/myapp:v1.0.0

# 验证签名
cosign verify myregistry.io/myapp:v1.0.0 \
  --certificate-identity=user@example.com \
  --certificate-oidc-issuer=https://accounts.google.com
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

---

## CI/CD 集成

### GitHub Actions

```yaml
name: Build and Sign

on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      id-token: write  # 用于 OIDC
    steps:
      - uses: actions/checkout@v4

      - name: Build image
        run: docker build -t myapp:${{ github.sha }} .

      - name: Install Cosign
        uses: sigstore/cosign-installer@v3

      - name: Sign image (keyless)
        run: |
          cosign sign --yes \
            --oidc-issuer https://token.actions.githubusercontent.com \
            myregistry.io/myapp:${{ github.sha }}

      - name: Generate SBOM
        uses: anchore/sbom-action@v0
        with:
          image: myregistry.io/myapp:${{ github.sha }}
          format: cyclonedx-json
          output-file: sbom.json

      - name: Attach SBOM
        run: |
          cosign attach sbom --sbom sbom.json \
            myregistry.io/myapp:${{ github.sha }}
```

---

## K8s 准入控制

### Kyverno 策略

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
            - "myregistry.io/*"
          attestors:
            - entries:
                - keyless:
                    issuer: https://token.actions.githubusercontent.com
                    subject: https://github.com/myorg/myrepo/.github/workflows/build.yml@refs/heads/main
                - keys:
                    publicKeys: |
                      -----BEGIN PUBLIC KEY-----
                      MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE8n...
                      -----END PUBLIC KEY-----
```

---

## 扫描工具

| 工具 | 用途 | 特点 |
|------|------|------|
| **Trivy** | 镜像/文件系统扫描 | 全面，支持 SBOM |
| **Grype** | 漏洞扫描 | Anchore 出品 |
| **Syft** | SBOM 生成 | 多格式输出 |
| **Snyk** | 依赖扫描 | 商业工具 |
| **SLSA GitHub Generator** | SLSA 证明 | 自动生成 |

### Trivy 扫描

```bash
# 镜像扫描
trivy image myapp:latest

# 生成 SBOM
trivy image --format cyclonedx -o sbom.json myapp:latest

# 扫描 IaC
trivy config ./terraform

# CI 模式 (返回非零退出码)
trivy image --exit-code 1 --severity HIGH,CRITICAL myapp:latest
```

---

## 供应链安全清单

### 开发阶段

- [ ] 代码签名 (GPG)
- [ ] 依赖审查 (Dependabot, Snyk)
- [ ] 密钥扫描 (TruffleHog, GitLeaks)
- [ ] SAST (SonarQube, CodeQL)

### 构建阶段

- [ ] 隔离构建环境
- [ ] 不可变构建
- [ ] 签名容器镜像
- [ ] 生成 SBOM
- [ ] SLSA 证明

### 部署阶段

- [ ] 镜像签名验证
- [ ] 准入控制策略
- [ ] 运行时安全监控
- [ ] 漏洞扫描

---

## 2025 趋势

- **SLSA 1.1**: 新增源代码追踪
- **VEX**: 漏洞利用交换，减少误报
- **GUAC**: 开源供应链分析
- **in-toto**: 供应链元数据框架
- **二进制授权**: Google 风格的部署验证
