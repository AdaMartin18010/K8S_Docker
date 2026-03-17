# Docker Compose 生产级配置

## 文件说明

| 文件 | 用途 |
|------|------|
| `docker-compose.yml` | 基础服务定义 |
| `docker-compose.override.yml` | 开发环境覆盖 |
| `docker-compose.prod.yml` | 生产环境配置 |
| `.env.example` | 环境变量模板 |
| `Makefile` | 常用命令快捷方式 |

## 快速开始

```bash
# 1. 准备环境变量
cp .env.example .env
# 编辑 .env 填写实际值

# 2. 启动开发环境
make dev

# 3. 查看日志
make logs

# 4. 停止服务
make down
```

## 多环境配置

### 开发环境

```bash
make dev
# 或
ENV=dev docker compose up -d
```

### 生产环境

```bash
make prod
# 或
ENV=prod docker compose up -d
```

## 包含服务

- **app**: Go 应用程序
- **postgres**: PostgreSQL 数据库
- **redis**: Redis 缓存
- **nginx**: Nginx 反向代理
- **prometheus**: 指标收集
- **grafana**: 可视化监控
- **adminer**: 数据库管理（仅开发）
- **redis-commander**: Redis 管理（仅开发）

## 关键特性

1. **资源限制**: 所有服务都有 CPU/内存限制
2. **健康检查**: 自动检测服务健康
3. **只读文件系统**: 生产环境安全加固
4. **Secret 管理**: 敏感信息分离
5. **网络隔离**: frontend/backend 分离
6. **日志轮转**: 防止磁盘占满
