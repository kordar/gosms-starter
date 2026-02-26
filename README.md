# gosms-starter

用于 [gosms](https://github.com/kordar/gosms) 的轻量启动器，提供按 name 注册/获取 Provider 的全局句柄，并支持参考 dig-starter 风格的模块化加载回调（在加载时把 moduleName、id 和构造好的 Provider 回调给调用方）。

## 功能
- 全局按名称管理 `SMSProvider`：`Provide`、`Get`、`ProvideFromConfig`、`ProvideEFromConfig`
- 模块化加载：`NewSMSModule(name, load)` 与 `Load`，支持单实例与多实例配置
- 与具体短信厂商解耦：通过 `gosms` 的 Provider 插件机制工作（如 `gosms-aliyun`）

## 安装
```bash
go get github.com/kordar/gosms-starter
```

## 依赖
- gosms：核心短信接口与模型
- 各厂商 Provider 实现（至少引入一种，例如 Aliyun）：
  ```go
  import _ "github.com/kordar/gosms-aliyun"
  ```

## 快速开始（直接注册并获取）
```go
package main

import (
    smsstarter "github.com/kordar/gosms-starter"
    "github.com/kordar/gosms"
    _ "github.com/kordar/gosms-aliyun" // 通过 init 注册 aliyun provider
)

func main() {
    cfg := gosms.NewSMSConfig("aliyun", "AK", "SK").
        WithSign("签名").
        WithTemplate("SMS_123456").
        WithExtraParam("endpoint", "dysmsapi.aliyuncs.com")

    // 构造并注册到全局，出现错误直接中止（Fatal）
    p := smsstarter.ProvideEFromConfig("aliyun-main", cfg)

    // 也可以在任意位置通过名称获取
    // p := smsstarter.Get("aliyun-main")

    // 使用 gosms 的请求模型发送短信
    req := gosms.NewSMSRequest([]string{"13800138000"}, "").
        WithTemplateID("SMS_123456").
        WithTemplateVar("code", "1234")

    _, _ = p.SendTemplate(*req)
}
```

## 模块化加载（参考 dig-starter 风格）
`NewSMSModule(name, load)` 接收一个 `load` 回调，在每个实例构建完成后回调：
```go
func(moduleName string, itemId string, p gosms.SMSProvider, item map[string]interface{})
```
- `moduleName`：模块名（创建模块时传入）
- `itemId`：本次加载的实例 id
- `p`：根据配置构建完成的 `SMSProvider`
- `item`：本实例用于构建的原始配置（map）

### 单实例配置
```go
import (
    smsstarter "github.com/kordar/gosms-starter"
    "github.com/kordar/gosms"
    _ "github.com/kordar/gosms-aliyun"
)

func main() {
    m := smsstarter.NewSMSModule("sms", func(moduleName, itemId string, p gosms.SMSProvider, item map[string]interface{}) {
        // 在这里把 p 放入你自己的容器/路由/业务上下文中
        // 例如：myContainer.BindSMS(itemId, p)
    })

    m.Load(map[string]interface{}{
        "id":         "aliyun-main",
        "provider":   "aliyun",
        "access_key": "AK",
        "secret_key": "SK",
        "sign":       "签名",
        "template":   "SMS_123456",
        "extra": map[string]string{
            "endpoint": "dysmsapi.aliyuncs.com",
        },
    })
}
```

### 多实例配置
```go
m.Load(map[string]interface{}{
    "aliyun-main": map[string]interface{}{
        "provider":   "aliyun",
        "access_key": "AK1",
        "secret_key": "SK1",
        "sign":       "签名1",
        "template":   "SMS_111111",
    },
    "aliyun-bak": map[string]interface{}{
        "provider":   "aliyun",
        "access_key": "AK2",
        "secret_key": "SK2",
        "sign":       "签名2",
        "template":   "SMS_222222",
        "extra": map[string]string{
            "endpoint": "dysmsapi.aliyuncs.com",
        },
    },
})
```

## 配置字段
- `id`：实例名称（单实例时可放在顶层；多实例时为 map 的 key）
- `provider`：短信服务提供商标识（如 `aliyun`）
- `access_key` / `secret_key`：访问凭证
- `sign`：短信签名（可选）
- `template`：默认模板 ID（可选）
- `extra`：`map[string]string` 扩展参数（可选），例如 `endpoint`

## API 速览
- 全局句柄
  - `Provide(name string, p gosms.SMSProvider)`：注册现成实例
  - `Get(name string) gosms.SMSProvider`：按名称获取实例（未找到将 Fatal）
  - `ProvideFromConfig(name string, cfg *gosms.SMSConfig) (gosms.SMSProvider, error)`：由配置构建并注册，返回实例
  - `ProvideEFromConfig(name string, cfg *gosms.SMSConfig) gosms.SMSProvider`：出错 Fatal，返回实例
- 模块
  - `NewSMSModule(name string, load func(moduleName, itemId string, p gosms.SMSProvider, item map[string]interface{})) *SMSModule`
  - `Load(value interface{})`：
    - 传入单实例 `map[string]interface{}` 且包含 `id` 字段；或
    - 传入多实例 `map[string]interface{}`（key 为实例 id）

## 常见问题
1. 为什么需要 `blank import` 具体厂商包？  
   Provider 的注册通过各实现包的 `init()` 完成，必须引入对应实现（如 `gosms-aliyun`）。
2. `Get(name)` 找不到实例会怎样？  
   内部使用 `Fatalf` 直接终止，请在应用启动阶段保证实例注册成功。
3. 有没有“默认实例”的概念？  
   当前实现不内置默认实例选择器；如需默认实例，请在业务层自行约定并通过 `Get("your-default")` 使用。

## License
MIT

