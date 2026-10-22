# playground/126 (Go 1.26)

本目录演示 Go 1.26 **已正式发布**的语言与标准库能力 (不含 `GOEXPERIMENT` 实验项).

## 运行

```bash
make test
# 或
GOWORK=off go test ./...
```

保存测试产物到固定目录:

```bash
GOWORK=off go test -artifacts ./artifacts/...
```

## Demo 索引

| 目录 | 特性 |
|------|------|
| [newexpr](./newexpr/) | `new` 接受表达式操作数 |
| [generic_adder](./generic_adder/) | 泛型约束自引用 |
| [errorsastype](./errorsastype/) | `errors.AsType` |
| [bytespeek](./bytespeek/) | `bytes.Buffer.Peek` |
| [slogmulti](./slogmulti/) | `slog.NewMultiHandler` |
| [reflectiter](./reflectiter/) | `reflect` 字段迭代器 |
| [artifacts](./artifacts/) | `testing.T.ArtifactDir` |
| [hpke](./hpke/) | `crypto/hpke` |

参考: [Go 1.26 Release Notes](https://go.dev/doc/go1.26)
