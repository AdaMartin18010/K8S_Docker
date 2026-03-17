#!/bin/bash
# Cosign 签名验证示例

IMAGE="myregistry.io/myapp:v1.0.0"

echo "=== 无密钥签名 ==="
cosign sign --yes $IMAGE

echo "=== 验证签名 ==="
cosign verify $IMAGE \
  --certificate-identity="user@example.com" \
  --certificate-oidc-issuer="https://accounts.google.com"

echo "=== 生成 SBOM ==="
syft $IMAGE -o cyclonedx-json > sbom.json
cosign attach sbom --sbom sbom.json $IMAGE
