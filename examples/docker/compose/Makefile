# Docker Compose 管理 Makefile
.PHONY: all dev prod down logs ps clean

# 默认环境
ENV ?= dev

# 根据环境选择 compose 文件
ifeq ($(ENV),prod)
	COMPOSE_FILES := -f docker-compose.yml -f docker-compose.prod.yml
else
	COMPOSE_FILES := -f docker-compose.yml -f docker-compose.override.yml
endif

# Docker Compose 命令
DC := docker compose $(COMPOSE_FILES)

all: dev

# 启动开发环境
dev:
	$(DC) up -d --build

# 启动生产环境
prod:
	ENV=prod $(DC) up -d

# 停止所有服务
down:
	$(DC) down --volumes --remove-orphans

# 查看日志
logs:
	$(DC) logs -f --tail=100

# 查看状态
ps:
	$(DC) ps

# 重启服务
restart:
	$(DC) restart

# 构建镜像
build:
	$(DC) build --parallel

# 拉取最新镜像
pull:
	$(DC) pull

# 执行数据库迁移
migrate:
	$(DC) exec app /app migrate up

# 备份数据库
backup:
	mkdir -p backups
	docker exec postgres pg_dump -U app myapp > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql

# 进入容器 shell
shell-app:
	$(DC) exec app sh

shell-db:
	$(DC) exec postgres psql -U app -d myapp

shell-redis:
	$(DC) exec redis redis-cli

# 健康检查
health:
	$(DC) ps --format "table {{.Name}}\t{{.Status}}\t{{.Health}}"

# 清理
 clean:
	$(DC) down -v --rmi all --remove-orphans
	docker system prune -f

# 部署到生产（示例）
deploy-prod:
	@echo "Deploying to production..."
	docker context use production
	ENV=prod $(DC) up -d --build
