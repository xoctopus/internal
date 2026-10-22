# reflect 字段迭代器

Go 1.26 为 **`reflect.Type`** 与 **`reflect.Value`** 增加基于 **`range` 的迭代方法**, 遍历 struct 字段时不必再手写下标循环.

## API (本 demo 用到)

| 方法 | 说明 |
|------|------|
| `Type.Fields()` | 迭代 struct 的 `StructField` |
| `Value.Fields()` | 每次 yield 字段元信息与字段 `Value` |

同版本还有 `Type.Methods`, `Type.Ins`, `Type.Outs`, `Value.Methods` 等, 本目录仅演示 struct 字段.

## 运行

```bash
GOWORK=off go test ./reflectiter/...
```

## 参考

- [Go 1.26 - reflect](https://go.dev/doc/go1.26#reflect)
