# errors.AsType

Go 1.26 在 `errors` 包新增泛型函数 **`AsType[T]`**, 是 **`As` 的类型安全版本**: 成功时直接得到 `T`, 无需类型断言.

## 对比

```go
var ce codedError
if errors.As(err, &ce) { ... }
```

```go
if ce, ok := errors.AsType[codedError](err); ok { ... }
```

本目录演示包装后的 `error` 能否解出 `codedError`, 以及无关错误返回 `ok == false`.

## 运行

```bash
GOWORK=off go test ./errorsastype/...
```

## 参考

- [Go 1.26 - errors](https://go.dev/doc/go1.26#errors)
