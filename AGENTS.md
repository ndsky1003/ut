# AGENTS.md

单 package 的 Go 工具库 `github.com/ndsky1003/ut`（package `ut`），无 CI / Makefile / lint 配置。代码注释与文档统一用中文。

## 常用命令

- 构建 / 静态检查 / 测试：`go build ./...` / `go vet ./...` / `go test ./...`
- 跑单个测试：`go test -run TestGenID_Range ./...`
- fuzz：`go test -fuzz=FuzzGenID -fuzztime=10s -run '^$' ./...`
- benchmark：`go test -bench BenchmarkGenID -benchtime=200000x -run '^$' ./...`

## 关键约定（容易踩坑）

- `GenID` 签名：`func GenID(num_width uint8, canUse func(id uint64) bool) uint64`
  - 返回 `uint64`、`canUse` 接收 `uint64`。
  - 无限重试直到 `canUse` 返回 `true`；`canUse` 永远返回 `false` 会死循环，调用方需保证最终有可用 ID。
- `GenID` 用 `rand.Uint64N(max-min+1)` 保证均匀分布，是刻意设计；不要改回「从随机数取位窗口」的旧算法（首位数字会严重偏斜）。
- 方法接收者命名不统一：`line.go` / `option.go` 用 `this`，`rcumap.go` 用 `m`；改文件时跟随该文件现有风格。

## 模块速览

- `rand.go`：`GenID`（均匀随机 ID）、`Pick` / `Peek`（不放回抽样）、`Shuffle`
- `base.go`：`Base` / `Parse` 任意进制（2~62）转换
- `line.go`：`Line[T]` 轮询队列（自动洗牌）
- `option.go`：函数式选项模式（`*Option().Merge(...)`）
- `ut.go`：`ProtectRun`（panic 后重试）
- `rcumap.go`：`RCUMap` 并发 map（读无锁、写 copy-on-write + CAS）
