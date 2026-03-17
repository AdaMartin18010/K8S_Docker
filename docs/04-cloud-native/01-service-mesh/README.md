# 服务网格

> Service Mesh 技术详解

---

## 本章内容

- [Istio 指南](./istio-guide.md)

---

## 服务网格对比

| 特性 | Istio | Linkerd | Cilium |
|------|-------|---------|--------|
| 架构 | Sidecar/Ambient | Sidecar | eBPF |
| 性能 | 中等 | 高 | 最高 |
| 资源占用 | 高 | 低 | 低 |
| 功能完整度 | 最全 | 核心功能 | 网络+L7 |

---

## 2025 推荐

- **新部署**: Istio Ambient Mesh (无 Sidecar)
- **资源敏感**: Linkerd
- **已有 Cilium**: Cilium Service Mesh

---

## 相关文档

- [Cilium eBPF](../../04-ecosystem/ebpf-cilium/)
