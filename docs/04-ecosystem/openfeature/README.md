# OpenFeature - 标准化特性标志管理

## 概述

OpenFeature 是一个 CNCF Incubating 项目，提供供应商无关的、社区驱动的特性标志 API 规范。它支持标准化特性标志管理，避免供应商锁定，可与任何特性标志管理工具或内部解决方案配合使用。

## 核心特性

| 特性 | 描述 |
|------|------|
| 供应商无关 | 统一 API，支持多种后端 |
| 多语言 SDK | 支持 10+ 编程语言 |
| 上下文感知 | 基于用户属性动态评估 |
| 可观测性 | 与 OpenTelemetry 集成 |
| Provider 插件 | 灵活切换后端实现 |

## 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                    OpenFeature 架构                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Application Code                      │   │
│  │                                                          │   │
│  │   client.BooleanValue("new-feature", false, context)    │   │
│  │   client.StringValue("theme", "dark", context)          │   │
│  │   client.NumberValue("threshold", 100, context)         │   │
│  │                                                          │   │
│  └──────────────────────┬──────────────────────────────────┘   │
│                         │                                       │
│                         ▼                                       │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  OpenFeature SDK                         │   │
│  │                                                          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │   │
│  │  │   Hooks     │  │  Evaluation │  │   Events    │      │   │
│  │  │  (钩子)      │  │   (评估)     │  │   (事件)     │      │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘      │   │
│  │                                                          │   │
│  └──────────────────────┬──────────────────────────────────┘   │
│                         │                                       │
│                         ▼                                       │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  Provider Interface                      │   │
│  └──────────────────────┬──────────────────────────────────┘   │
│                         │                                       │
│           ┌─────────────┼─────────────┐                        │
│           ▼             ▼             ▼                        │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐           │
│  │  LaunchDarkly│ │  Flagd       │ │  Custom      │           │
│  │  Provider    │ │  Provider    │ │  Provider    │           │
│  └──────────────┘ └──────────────┘ └──────────────┘           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 多语言 SDK

### Go SDK

```go
package main

import (
    "context"
    "fmt"

    "github.com/open-feature/go-sdk/openfeature"
)

func main() {
    // 配置 Provider
    provider := openfeature.NewNoopProvider()
    openfeature.SetProvider(provider)

    // 创建客户端
    client := openfeature.NewClient("my-app")

    // 创建评估上下文
    ctx := openfeature.NewEvaluationContext(
        "user-123",
        map[string]interface{}{
            "tenant": "premium",
            "region": "us-west",
        },
    )

    // 评估布尔标志
    enabled, err := client.BooleanValue(
        context.Background(),
        "new-feature",
        false,  // 默认值
        ctx,
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if enabled {
        fmt.Println("New feature is enabled!")
    }

    // 评估字符串标志
    theme, _ := client.StringValue(
        context.Background(),
        "ui-theme",
        "light",
        ctx,
    )
    fmt.Printf("Theme: %s\n", theme)

    // 评估数字标志
    threshold, _ := client.NumberValue(
        context.Background(),
        "cache-threshold",
        1000,
        ctx,
    )
    fmt.Printf("Threshold: %f\n", threshold)

    // 评估对象标志
    config, _ := client.ObjectValue(
        context.Background(),
        "feature-config",
        map[string]interface{}{},
        ctx,
    )
    fmt.Printf("Config: %v\n", config)
}
```

### Python SDK

```python
from openfeature import api
from openfeature.evaluation_context import EvaluationContext

def main():
    # 设置 Provider
    api.set_provider(openfeature.provider.NoOpProvider())

    # 获取客户端
    client = api.get_client()

    # 创建评估上下文
    context = EvaluationContext(
        targeting_key="user-123",
        attributes={
            "tenant": "premium",
            "region": "us-west",
        }
    )

    # 评估标志
    enabled = client.get_boolean_value(
        flag_key="new-feature",
        default_value=False,
        evaluation_context=context,
    )

    if enabled:
        print("New feature is enabled!")

    # 字符串标志
    theme = client.get_string_value(
        flag_key="ui-theme",
        default_value="light",
        evaluation_context=context,
    )
    print(f"Theme: {theme}")

    # 数字标志
    threshold = client.get_number_value(
        flag_key="cache-threshold",
        default_value=1000,
        evaluation_context=context,
    )
    print(f"Threshold: {threshold}")

    # 对象标志
    config = client.get_object_value(
        flag_key="feature-config",
        default_value={},
        evaluation_context=context,
    )
    print(f"Config: {config}")

if __name__ == "__main__":
    main()
```

## Provider 实现

### Flagd Provider

```go
import (
    flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
)

func main() {
    // 使用 Flagd (开源特性标志服务)
    provider, err := flagd.NewProvider(
        flagd.WithHost("localhost"),
        flagd.WithPort(8013),
    )
    if err != nil {
        log.Fatal(err)
    }

    openfeature.SetProvider(provider)
}
```

### 自定义 Provider

```go
type MyCustomProvider struct{}

func (p *MyCustomProvider) Metadata() openfeature.Metadata {
    return openfeature.Metadata{
        Name: "MyCustomProvider",
    }
}

func (p *MyCustomProvider) BooleanEvaluation(
    ctx context.Context,
    flag string,
    defaultValue bool,
    evalCtx openfeature.FlattenedContext,
) openfeature.BoolResolutionDetail {
    // 自定义评估逻辑
    value := evaluateFromMyBackend(flag, evalCtx)

    return openfeature.BoolResolutionDetail{
        Value: value,
        ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
            Reason: openfeature.TargetingMatchReason,
        },
    }
}

// 实现其他评估方法...
```

## Hooks 机制

```go
// 日志 Hook
type LoggingHook struct{}

func (h LoggingHook) Before(
    ctx context.Context,
    hookContext openfeature.HookContext,
    hints openfeature.HookHints,
) (*openfeature.EvaluationContext, error) {
    log.Printf("Evaluating flag: %s", hookContext.FlagKey())
    return nil, nil
}

func (h LoggingHook) After(
    ctx context.Context,
    hookContext openfeature.HookContext,
    flagDetails openfeature.InterfaceResolutionDetail,
    hints openfeature.HookHints,
) error {
    log.Printf("Flag evaluated: %s = %v", hookContext.FlagKey(), flagDetails.Value)
    return nil
}

// 使用 Hook
client := openfeature.NewClient("my-app")
client.AddHooks(LoggingHook{})
```

## OpenTelemetry 集成

```go
import (
    "github.com/open-feature/go-sdk-contrib/hooks/open-telemetry/pkg"
)

func main() {
    // 添加 OTel Hook
    client := openfeature.NewClient("my-app")
    client.AddHooks(otelhook.NewHook())

    // 现在每次标志评估都会生成追踪数据
}
```

## Kubernetes 集成

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
      - name: app
        image: my-app:latest
        env:
        - name: FLAGD_HOST
          value: "flagd.flagd.svc.cluster.local"
        - name: FLAGD_PORT
          value: "8013"
```

## 支持的 Provider

| Provider | 类型 | 描述 |
|----------|------|------|
| LaunchDarkly | 商业 | 企业级特性管理平台 |
| Flagsmith | 开源/商业 | 开源特性标志服务 |
| Flagd | 开源 | OpenFeature 原生服务 |
| Flipt | 开源 | 自托管特性标志 |
| GrowthBook | 开源 | A/B 测试和特性标志 |

## 相关资源

- [OpenFeature 官网](https://openfeature.dev/)
- [GitHub](https://github.com/open-feature)
- [规范文档](https://openfeature.dev/specification/)
