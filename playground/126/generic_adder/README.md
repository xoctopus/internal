# 泛型约束自引用

Go 1.26 放宽限制: **泛型类型的类型参数列表里可以引用该泛型类型自身**. 便于表达「参数类型必须能和我做同一种运算」一类约束.

## 示例

```go
type Adder[A Adder[A]] interface {
    Add(A) A
}

func sum[A Adder[A]](x, y A) A {
    return x.Add(y)
}
```

此前 `Adder[A Adder[A]]` 中的自引用不被允许.

本目录用 `IntAdder` 实现 `Add`, 并通过 `sum` 验证.

## 运行

```bash
GOWORK=off go test ./generic_adder/...
```

## 参考

- [Go 1.26 - Changes to the language](https://go.dev/doc/go1.26#language)
