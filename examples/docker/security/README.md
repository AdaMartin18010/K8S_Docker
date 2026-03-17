# Docker 安全加固示例

本目录包含容器安全最佳实践，遵循 CIS Docker Benchmark。

## 安全策略

### 1. 镜像安全

- ✅ 使用最小基础镜像 (scratch/distroless)
- ✅ 多阶段构建分离构建依赖
- ✅ 非 root 用户运行
- ✅ 移除 shell 和包管理器

### 2. 运行时安全

- ✅ 只读文件系统
- ✅ 临时文件系统限制
- ✅ 能力 (capabilities) 限制
- ✅ seccomp 配置文件
- ✅ AppArmor 策略

### 3. 网络安全

- ✅ 内部网络隔离
- ✅ 不暴露不必要端口
- ✅ 网络策略限制

### 4. 资源限制

- ✅ CPU/内存限制
- ✅ PID 限制（防止 fork 炸弹）
- ✅ 磁盘 I/O 限制

## 使用

```bash
# 构建安全镜像
docker build -f Dockerfile.secure -t secure-app .

# 运行（使用 docker-compose）
docker compose -f docker-compose.security.yml up -d

# 检查安全配置
docker exec secure-app ps aux  # 查看用户
docker inspect secure-app | grep -A 20 "HostConfig"
```

## 安全检查清单

- [ ] 镜像使用非 root 用户
- [ ] 文件系统设置为只读
- [ ] 所有能力已删除
- [ ] seccomp 配置文件已加载
- [ ] 资源限制已设置
- [ ] 健康检查已配置
- [ ] 没有敏感信息硬编码
