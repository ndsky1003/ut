# ut

Go 工具库 `github.com/ndsky1003/ut`（package `ut`），集合若干日常开发用的小工具：均匀随机 ID、任意进制转换、并发 map、轮询队列、崩溃保护等。

## 安装

```sh
go get github.com/ndsky1003/ut
```

## GenID —— 均匀随机 ID

按十进制位数生成均匀分布的随机 ID（1~19 位），适合玩家 ID 等场景。

```go
import "github.com/ndsky1003/ut"

// 生成一个 8 位 ID；canUse 用于去重（返回 false 表示已被占用，需重新生成）
id := ut.GenID(8, func(v uint64) bool { return !db.Exists(v) })
```

- 签名：`func GenID(num_width uint8, canUse func(id uint64) bool) uint64`
- `num_width`：1~19 的十进制位数。
- `canUse`：判断 ID 是否可用（如未被占用）；传 `nil` 会 panic。
- 会一直随机生成直到 `canUse` 返回 `true`；若 `canUse` 永远返回 `false` 会死循环，调用方需保证最终有可用 ID。

## Pick / Peek / Shuffle —— 随机抽取与洗牌

```go
// Pick：不放回随机抽取（结果不重复）
picked := ut.Pick([]int{1, 2, 3, 4, 5, 6}, 3)

// Peek：不放回随机抽取（结果不重复）
peeked := ut.Peek([]int{1, 2, 3, 4, 5, 6}, 3)

// Shuffle：洗牌（可指定洗牌次数，默认 1 次）
ut.Shuffle(arr, ut.ShuffleOption().SetCount(2))
```

## Base / Parse —— 任意进制转换（2~62）

```go
s := ut.Base(255, 16)        // "ff"
n, err := ut.Parse("ff", 16) // 255
```

## RCUMap —— 并发安全的 map

读操作无锁，写操作 copy-on-write + CAS，适合读多写少的场景。

```go
m := ut.NewRCUMap[string, int]()
m.Set("a", 1)
v, ok := m.Get("a")
m.Range(func(k string, v int) bool {
	// 返回 false 可提前停止遍历
	return true
})
m.Remove("a")
```

## Line —— 轮询队列（自动洗牌）

按顺序取出元素，一轮结束后自动洗牌进入下一轮。

```go
l := ut.NewLine[int]()
l.SetArr([]int{1, 2, 3})
v, err := l.Step() // 依次返回 1、2、3，一轮后再随机
```

## ProtectRun —— panic 后自动重试

```go
ut.ProtectRun(func() {
	// 可能 panic 的逻辑
}, ut.ProtectRunOption().SetInterval(time.Second))
```

## 选项模式

所有可配置项遵循 `XxxOption()` + `SetXxx(...)` + `Merge(...)` 的函数式选项模式，可用 `Merge` 组合多个选项后传入。
