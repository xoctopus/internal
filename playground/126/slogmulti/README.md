# slog.NewMultiHandler

Go 1.26 在 **`log/slog`** 中新增 **`NewMultiHandler`**, 把日志同时交给多个 **`Handler`** (例如一份 text, 一份 JSON).

## 行为

- **`Enabled`**: 任一子 handler 启用则返回 true.
- **`Handle` / `WithAttrs` / `WithGroup`**: 对每个已启用的子 handler 调用对应方法.

本目录用 `TextHandler` 与 `JSONHandler` 各写一份 buffer, 打一条带属性的日志, 断言两边都有输出.

## 运行

```bash
GOWORK=off go test ./slogmulti/...
```

## 参考

- [Go 1.26 - log/slog](https://go.dev/doc/go1.26#log_slog)
