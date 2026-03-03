# 错误处理规范文档

## 核心原则

1. **使用 `github.com/pkg/errors` 包**：替代标准库的 `errors` 和 `fmt.Errorf`
2. **错误包装而非替换**：保留原始错误信息，添加上下文

## 错误处理模式

### 1. 错误创建和包装

#### ❌ 错误做法
```go
// 使用标准库 errors
return errors.New("timeout must be positive")

// 使用 fmt.Errorf 丢失错误堆栈
return fmt.Errorf("failed to create http client: %w", err)
```

#### ✅ 正确做法
```go
import "github.com/pkg/errors"

// 使用 errors.Errorf 创建新错误
return errors.Errorf("timeout must be positive")

// 使用 errors.Wrap/Wrapf 包装错误，保留堆栈
return errors.Wrap(err, "failed to create http client")
return errors.Wrapf(err, "failed to marshal request body. body = %v", body)

```

### 2. 参数验证错误

```go
func WithTimeout(timeout time.Duration) Option {
    return func(c *DSPClient) error {
        if timeout <= 0 {
            return errors.Errorf("timeout must be positive, got: %v", timeout)
        }
        c.timeout = timeout
        return nil
    }
}

func WithAKSK(ak, sk string) Option {
    return func(c *DSPClient) error {
        if ak == "" || sk == "" {
            return errors.Errorf("ak and sk cannot be empty, ak=%s, sk=%s", ak, sk)
        }
        c.ak = ak
        c.sk = sk
        return nil
    }
}
```

### 错误日志最佳实践

1. **包含上下文信息**：记录请求ID、参数等关键信息
2. **记录错误堆栈**：使用 `errors.WithStack` 保留堆栈
3. **避免重复日志**：在错误传播过程中避免重复记录

```go
// 创建错误时记录
if err := someOperation(); err != nil {
    logger.Error("req-%s Operation failed: %v", reqID, err)
    return errors.Wrap(err, "someOperation failed")
}

// 上层不需要重复记录，只需传播
if err := lowerLayer(); err != nil {
    // 只添加上下文，不重复记录日志
    return errors.Wrapf(err, "req-%s lowerLayer failed for request", reqID)
}
```

## 总结

- **始终使用 `github.com/pkg/errors`** 进行错误处理
- **用 Wrap/Wrapf 保留错误上下文**，而不是替换错误
- **合理记录日志**，避免重复记录

遵循这些规范将使代码更加健壮、可维护，并提供更好的调试体验。