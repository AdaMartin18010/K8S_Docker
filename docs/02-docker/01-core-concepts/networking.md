# Docker 网络详解

> 容器网络原理与实践

---

## 网络模式

| 模式 | 说明 |
|------|------|
| **bridge** | 默认，NAT 网络 |
| **host** | 共享宿主机网络 |
| **none** | 无网络 |
| **container** | 共享其他容器网络 |

---

## Bridge 网络

```bash
# 创建自定义网络
docker network create my-network

# 运行容器并加入网络
docker run -d --name web --network my-network nginx

# 容器间通信（通过容器名）
docker run --rm --network my-network curlimages/curl http://web
```

---

## 网络命令

```bash
# 列出网络
docker network ls

# 查看网络详情
docker network inspect bridge

# 连接容器到网络
docker network connect my-network container

# 断开网络
docker network disconnect my-network container
```

---

## 端口映射

```bash
# 简单映射
docker run -p 8080:80 nginx

# 指定 IP
docker run -p 127.0.0.1:8080:80 nginx

# 随机端口
docker run -p 80 nginx
```
