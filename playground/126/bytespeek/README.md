# bytes.Buffer.Peek

Go 1.26 为 **`bytes.Buffer`** 增加 **`Peek(n int)`**: 查看接下来 `n` 个字节, **不推进读指针**.

## 用途

解析协议或流式数据时, 常需要「先看一眼再决定读多少」. 以前可能要 `Bytes()` 再手动记偏移; `Peek` 语义更清晰, 且与「读走数据」分离.

## 行为

- `Peek(3)` 返回前 3 字节, buffer 内容不变.
- 随后 `ReadByte` 等操作仍从 buffer 开头按序消费.

见 `bytespeek_test.go`.

## 运行

```bash
GOWORK=off go test ./bytespeek/...
```

## 参考

- [Go 1.26 - bytes](https://go.dev/doc/go1.26#bytes)
