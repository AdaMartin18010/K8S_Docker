# Docker 存储

> 容器数据持久化

---

## 存储类型

| 类型 | 用途 | 生命周期 |
|------|------|----------|
| **Volume** | 持久化数据 | 独立于容器 |
| **Bind Mount** | 开发调试 | 宿主机目录 |
| **tmpfs** | 敏感数据 | 内存存储 |

---

## Volume

```bash
# 创建卷
docker volume create my-vol

# 使用卷
docker run -v my-vol:/data nginx

# 查看卷
docker volume ls
docker volume inspect my-vol
```

---

## Bind Mount

```bash
# 开发时使用
docker run -v $(pwd):/app -w /app node npm start

# 只读挂载
docker run -v $(pwd)/config:/etc/nginx:ro nginx
```

---

## 多容器共享

```bash
# 数据卷容器模式
docker run -v /data --name data-busybox busybox true

docker run --volumes-from data-busybox -v $(pwd):/backup \
  busybox tar cvf /backup/backup.tar /data
```

---

## 最佳实践

1. 生产环境使用 Named Volume
2. 开发环境使用 Bind Mount
3. 敏感数据使用 tmpfs
4. 定期备份重要卷
