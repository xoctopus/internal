# testing.ArtifactDir

Go 1.26 为 **`testing.T` / `B` / `F`** 增加 **`ArtifactDir()`**, 返回可写入**测试产物**的目录. 它解决的是: 测试结果以**文件**形式存在 (二进制, 大图, 长报告), 不适合塞进 `t.Log`, 也不一定要进 git 里的 `testdata`.

## 应用场景

**可视化与回归**

- 测试里生成图表, 截图, PDF 预览帧等, 写入 `ArtifactDir`.
- 本地或 CI 使用 `go test -artifacts` 后, 人工或脚本对比目录中的文件 (例如图像 diff).

**大二进制或大体积快照**

- JSON/HTML 报告, 压缩包, 音频/字体子集等. 失败时把**实际输出**落盘, 再与期望文件 diff, 比把 MB 级内容打进日志更可读.

**失败时的调试包**

- 导出输入样例, 中间 AST/SQL/执行计划, `pprof` profile, 配置快照等.
- 平时普通 `go test` 用临时目录, 测完即删; CI nightly 或加 `-artifacts` 时保留, 配合流水线 `upload-artifact` 下载复现.

**代码生成 / 编译器类测试**

- 给定输入生成 Go, IR, SQL 等, 生成物写入 artifact, 供 review 或 `diff -ru`.

**与现有做法的对比**

| 方式 | 典型问题 |
|------|----------|
| `t.Log` | 只适合短文本, 二进制与大输出难处理 |
| 固定 `/tmp` 路径 | CI 难收集, 并行测试易冲突 |
| `ArtifactDir` | 统一 scratch; `-artifacts` 时路径可预期, 便于归档 |

本目录测试仅写入 `note.txt` 验证 API; 真实项目里常见的是 profile, 报告, 生成代码等.

## 与 `-artifacts` 标志

- **未**传 `-artifacts`: 使用临时目录, 测试结束后删除.
- **传入** `go test -artifacts`: 目录落在 `-outputdir` (默认当前目录) 下的 artifact 树中; 首次调用会在测试日志里打印 `=== ARTIFACTS ...` 行.

## 运行

```bash
GOWORK=off go test ./artifacts/...

# 保留产物到 output 目录
GOWORK=off go test -artifacts ./artifacts/...
```

## 参考

- [Go 1.26 - testing](https://go.dev/doc/go1.26#testing)
