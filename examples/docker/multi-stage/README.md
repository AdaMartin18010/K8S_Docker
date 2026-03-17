# 多阶段构建示例

演示 Docker 多阶段构建的高级用法，包含 8 个阶段。

## 阶段说明

| 阶段 | 目标 | 大小 |
|------|------|------|
| deps | 依赖下载 | ~350MB |
| codegen | 代码生成 | ~400MB |
| dev | 开发环境 | ~450MB |
| tester | 运行测试 | ~400MB |
| builder | 编译二进制 | ~400MB |
| security | 安全扫描 | ~200MB |
| production | 生产镜像 | ~15MB |
| debug | 调试镜像 | ~30MB |

## 使用 Makefile

```bash
# 完整流程
make all

# 仅生产构建
make production

# 开发模式
make dev
make run-debug

# 运行测试
make test

# 安全扫描
make security

# 查看各阶段大小
make sizes
```

## 关键特性

1. **BuildKit 缓存挂载**：`--mount=type=cache` 加速构建
2. **并行阶段**：deps 和 codegen 可并行执行
3. **输出产物**：使用 `--output` 导出测试结果
4. **条件构建**：通过 `--target` 选择特定阶段
