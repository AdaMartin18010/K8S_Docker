# Docker 网络详解

> 容器网络原理与实践 (Docker 28.x)

---

## 网络模式

| 模式 | 说明 | 使用场景 |
|------|------|----------|
| **bridge** | 默认，NAT 网络 | 单机容器通信 |
| **host** | 共享宿主机网络 | 高性能网络需求 |
| **none** | 无网络 | 完全隔离 |
| **container** | 共享其他容器网络 | Sidecar 模式 |
| **overlay** | 跨主机网络 | Docker Swarm |

---

## Bridge 网络（默认）

```
┌─────────────────────────────────────────────────────────────┐
│                      宿主机                                   │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Docker0 (网桥) 172.17.0.1               │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐           │  │
│  │  │ Container│  │ Container│  │ Container│           │  │
│  │  │ 172.17.0.2│  │ 172.17.0.3│  │ 172.17.0.4│          │  │
│  │  └────┬─────┘  └────┬─────┘  └────┬─────┘           │  │
│  │       └─────────────┴─────────────┘                  │  │
│  │                         │                            │  │
│  │                    NAT 转发                          │  │
│  └─────────────────────────┼────────────────────────────┘  │
│                            │                                 │
│                      宿主机网络接口                            │
│                            │                                 │
└────────────────────────────┼─────────────────────────────────┘
                             ▼
                         外部网络
```

### 创建和使用自定义网络

```bash
# 创建自定义网络
docker network create --driver bridge my-network

# 运行容器并加入网络
docker run -d --name web --network my-network nginx

# 容器间通信（通过容器名解析）
docker run --rm --network my-network curlimages/curl http://web

# 查看网络详情
docker network inspect my-network
```

---

## Host 网络模式

```bash
# 使用 host 网络（性能最好，无 NAT）
docker run -d --network host nginx

# 容器直接使用宿主机网络接口
# 注意：端口直接暴露在宿主机上，无需 -p 映射
```

**适用场景**:

- 高性能网络应用
- 需要访问宿主机所有网络接口
- 端口冲突不敏感的场景

---

## 端口映射

```bash
# 简单映射（主机8080 -> 容器80）
docker run -p 8080:80 nginx

# 指定 IP 绑定
docker run -p 127.0.0.1:8080:80 nginx

# 随机主机端口
docker run -p 80 nginx
# 查看随机分配的端口
docker port <container>

# 多端口映射
docker run -p 8080:80 -p 443:443 nginx

# UDP 端口
docker run -p 53:53/udp dns-server
```

---

## 容器间通信方式

### 1. 通过容器名（推荐）

```bash
# 创建自定义网络
docker network create my-app

# 运行数据库
docker run -d --name db --network my-app postgres:15

# 运行应用，通过容器名访问数据库
docker run -d --name app --network my-app myapp:latest
# 应用内连接字符串: postgres://db:5432/mydb
```

### 2. 通过 IP 地址

```bash
# 查看容器 IP
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' db

# 直接访问（不推荐，IP 会变）
ping 172.18.0.2
```

### 3. 通过 Link（已废弃）

```bash
# 旧方式，不推荐使用
docker run --link db:database myapp
```

---

## DNS 和容器发现

```bash
# 自定义 DNS
docker run --dns 8.8.8.8 --dns 8.8.4.4 nginx

# 自定义主机名
docker run --hostname mycontainer --add-host db:192.168.1.100 nginx

# 使用外部 DNS 搜索域
docker run --dns-search example.com nginx
```

---

## 高级网络配置

### 自定义网桥网络

```bash
# 创建带子网的网络
docker network create \
  --driver bridge \
  --subnet 172.28.0.0/16 \
  --gateway 172.28.0.1 \
  --opt com.docker.network.bridge.name=my-bridge \
  custom-network

# 使用自定义网络运行容器
docker run -d --network custom-network --ip 172.28.0.10 nginx
```

### 禁用容器间通信

```bash
# 创建隔离网络
docker network create --opt com.docker.network.bridge.enable_icc=false isolated
```

---

## 网络故障排查

```bash
# 查看容器网络配置
docker inspect <container> | jq '.[0].NetworkSettings'

# 进入容器网络命名空间
nsenter -t <container-pid> -n ip addr

# 查看容器路由表
docker exec <container> ip route

# 测试连通性
docker run --rm --network my-network nicolaka/netshoot \
  ping -c 3 target-container

# 抓包分析
docker run --rm --net container:<container-name> nicolaka/netshoot \
  tcpdump -i eth0 -w /tmp/capture.pcap
```

---

## Docker 28.x 新特性

| 特性 | 说明 |
|------|------|
| **IPv6 默认启用** | 更好的 IPv6 支持 |
| **网络性能优化** | 减少 iptables 规则，提高吞吐量 |
| **自定义 MTU** | 更好的 overlay 网络支持 |
| **加密 overlay** | Swarm overlay 网络默认加密 |

---

## 最佳实践

1. **使用自定义网络**: 不要依赖默认 bridge 网络
2. **通过服务名通信**: 避免使用 IP 地址
3. **限制端口暴露**: 只暴露必要的端口
4. **使用 host 网络谨慎**: 可能带来安全风险
5. **监控网络流量**: 使用 cAdvisor 等工具监控

---

## 关联文档

- [Docker Compose 网络](docs/02-docker/03-compose/)
- [Kubernetes 网络](docs/03-kubernetes/03-networking/)
