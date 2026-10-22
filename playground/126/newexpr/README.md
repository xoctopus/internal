# new 表达式操作数

Go 1.26 起, 内建函数 `new` 的操作数可以是**表达式**, 而不仅是类型. 返回值是指向该表达式求值结果的指针.

## 之前

需要先算值, 再取址:

```go
age := yearsSince(born)
person.Age = &age
```

或借助临时变量.

## 现在

```go
Age: new(yearsSince(born))
```

本目录 `newexpr_test.go` 还包含 `new(int64(300))` 与 `encoding/json` 可选指针字段的示例.

## 运行

```bash
GOWORK=off go test ./newexpr/...
```

## 参考

- [Go 1.26 - Changes to the language](https://go.dev/doc/go1.26#language)
