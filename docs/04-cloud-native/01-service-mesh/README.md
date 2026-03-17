# 服务网格

> Service Mesh 技术详解与选型

---

## 什么是服务网格？

服务网格是处理服务间通信的基础设施层：

- **流量管理**：路由、负载均衡、熔断
- **安全通信**：mTLS 自动加密
- **可观测性**：指标、日志、追踪

---

## 架构对比

| 模式 | Sidecar | Ambient | eBPF |
|------|---------|---------|------|
| 资源占用 | 高 | 低 | 最低 |
| 延迟 | 高 | 中 | 低 |
| 功能完整 | ✅ | ✅ | 核心功能 |
| 代表 | Istio | Istio Ambient | Cilium |

---

## 2025 推荐

- **新部署**: Istio Ambient Mesh
- **资源敏感**: Linkerd
- **已有 Cilium**: Cilium Service Mesh

---

## 本章内容

- [Istio 指南](./istio-guide.md)

---

## 服务网格 vs API 网关

| 功能 | 服务网格 | API 网关 |
|------|----------|----------|
| 服务发现 | ✅ | ✅ |
| 负载均衡 | ✅ | ✅ |
| 认证授权 | mTLS | OAuth/JWT |
| 限流 | ✅ | ✅ |
| 协议转换 | ❌ | ✅ |

---

## 相关文档

- [Cilium eBPF](../../04-ecosystem/ebpf-cilium/)
- [Gateway API](../../03-kubernetes/03-networking/gateway-api.md)
