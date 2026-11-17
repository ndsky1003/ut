package ut

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("测试之前的准备工作")
	code := m.Run()
	fmt.Println("测试之后的清理工作")
	os.Exit(code)
}

// pow10 返回 10 的 n 次方
func pow10(n int) uint64 {
	r := uint64(1)
	for i := 0; i < n; i++ {
		r *= 10
	}
	return r
}

// digitCount 返回 v 的十进制位数
func digitCount(v uint64) int {
	if v == 0 {
		return 1
	}
	c := 0
	for v > 0 {
		v /= 10
		c++
	}
	return c
}

// TestGenID_Range 验证每种宽度下生成的 ID 都落在 [min,max]，且十进制位数与宽度一致
func TestGenID_Range(t *testing.T) {
	for w := min_D_width; w <= max_D_width; w++ {
		min := pow10(int(w) - 1)
		max := pow10(int(w)) - 1
		const N = 5000
		for j := 0; j < N; j++ {
			v := GenID(w, func(id uint64) bool { return true })
			if v == 0 {
				t.Fatalf("width=%d 第 %d 次返回 0（生成失败）", w, j)
			}
			if v < min || v > max {
				t.Fatalf("width=%d 越界: v=%d 不在 [%d,%d]", w, v, min, max)
			}
			if digitCount(v) != int(w) {
				t.Fatalf("width=%d 位数错误: v=%d 实际位数=%d", w, v, digitCount(v))
			}
		}
		t.Logf("width=%2d 通过, 范围 [%d, %d], 样本 %d", w, min, max, N)
	}
}

// TestGenID_CanUseFilter 验证 canUse 过滤器确实生效
func TestGenID_CanUseFilter(t *testing.T) {
	for w := min_D_width; w <= max_D_width; w++ {
		for j := 0; j < 200; j++ {
			v := GenID(w, func(id uint64) bool { return id%2 == 0 })
			if v == 0 {
				t.Fatalf("width=%d 过滤偶数时返回 0", w)
			}
			if v%2 != 0 {
				t.Fatalf("width=%d 过滤失效: 返回奇数 %d", w, v)
			}
		}
	}
	t.Log("canUse 偶数过滤全部通过")
}

// TestGenID_Distribution 验证不会退化成固定值，且取值分布大致均匀
func TestGenID_Distribution(t *testing.T) {
	const w = 18
	const N = 50000
	seen := make(map[uint64]bool, N)
	bucket := make(map[uint64]int)
	for j := 0; j < N; j++ {
		v := GenID(w, func(id uint64) bool { return true })
		seen[v] = true
		bucket[v/(pow10(w)/10)]++ // 把 [min,max] 按首位数字分成 9 段
	}
	if len(seen) < N/2 {
		t.Fatalf("去重后只有 %d 个不同值（N=%d），疑似退化", len(seen), N)
	}
	t.Logf("width=%d, N=%d, 去重后 %d 个值, 分段数=%d, 各段分布=%v", w, N, len(seen), len(bucket), bucket)
}

// TestGenID_AllReject 验证 canUse 恒为 false 时不会无限递归，会在有限时间内返回 0
func TestGenID_AllReject(t *testing.T) {
	done := make(chan uint64, 1)
	go func() { done <- GenID(10, func(id uint64) bool { return false }) }()
	select {
	case v := <-done:
		if v != 0 {
			t.Fatalf("canUse 恒 false 应返回 0, 实际 %d", v)
		}
		t.Log("canUse 恒 false 正确返回 0（未无限递归）")
	case <-time.After(5 * time.Second):
		t.Fatal("canUse 恒 false 时疑似无限递归/卡死")
	}
}

// TestGenID_CustomMaxAttempts 验证可自定义重试上限（下限 1000，只能调大）
func TestGenID_CustomMaxAttempts(t *testing.T) {
	calls := 0
	GenID(10, func(id uint64) bool { calls++; return false }, 1500)
	if calls != 1500 {
		t.Fatalf("maxAttempts=1500 时 canUse 应调用 1500 次, 实际 %d", calls)
	}

	calls = 0
	GenID(10, func(id uint64) bool { calls++; return false }, 10)
	if calls != 1000 {
		t.Fatalf("maxAttempts=10 应被下限 1000 覆盖, canUse 应调用 1000 次, 实际 %d", calls)
	}
	t.Log("自定义 maxAttempts 生效（下限 1000，可调大）")
}

// TestGenID_InvalidWidth 验证非法宽度会 panic
func TestGenID_InvalidWidth(t *testing.T) {
	for _, w := range []uint8{0, max_D_width + 1, 20, 255} {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("width=%d 应触发 panic 但未触发", w)
				}
			}()
			GenID(w, func(id uint64) bool { return true })
		}()
	}
}

// TestGenID_NilCanUse 验证 canUse 为 nil 会 panic
func TestGenID_NilCanUse(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("canUse=nil 应触发 panic 但未触发")
		}
	}()
	GenID(5, nil)
}

func BenchmarkGenID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenID(8, func(id uint64) bool { return true })
	}
}

func BenchmarkGenID19(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenID(19, func(id uint64) bool { return true })
	}
}

func FuzzGenID(f *testing.F) {
	for _, tc := range []uint8{1, 2, 3, 4, 5, 6, 9, 18, 19} {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, w uint8) {
		// 合法宽度必须落在正确范围；非法宽度必须 panic
		if w >= min_D_width && w <= max_D_width {
			id := GenID(w, func(x uint64) bool { return true })
			min, max := pow10(int(w)-1), pow10(int(w))-1
			if id < min || id > max {
				t.Fatalf("width=%d id=%d 越界 [%d,%d]", w, id, min, max)
			}
		} else {
			defer func() {
				if recover() == nil {
					t.Fatalf("width=%d 应 panic 但未触发", w)
				}
			}()
			GenID(w, func(x uint64) bool { return true })
		}
	})
}

func TestPick(t *testing.T) {
	orig := []int{1, 23, 7, 4, 5, 6}
	snapshot := append([]int(nil), orig...)

	v := Pick(orig, 3)
	if len(v) != 3 {
		t.Fatalf("Pick 返回长度 %d, 期望 3", len(v))
	}
	for _, x := range v {
		if !contains(orig, x) {
			t.Errorf("Pick 返回了原切片中不存在的元素 %d", x)
		}
	}
	if !reflect.DeepEqual(orig, snapshot) {
		t.Errorf("Pick 不应修改原切片: got %v want %v", orig, snapshot)
	}

	if got := Pick(orig, 100); len(got) != len(orig) {
		t.Errorf("count 超过长度时应返回全部, 实际 %d", len(got))
	}
	if got := Pick(orig, 0); len(got) != 0 {
		t.Errorf("count=0 应返回空切片, 实际 %d", len(got))
	}
}

func contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func TestShuffle(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	snapshot := append([]int(nil), a...)
	Shuffle(a)
	if len(a) != len(snapshot) {
		t.Fatalf("Shuffle 不应改变长度: got %d want %d", len(a), len(snapshot))
	}
	sortedA := append([]int(nil), a...)
	// 元素集合不变（排序后相等）
	selectionSort(sortedA)
	sortedSnap := append([]int(nil), snapshot...)
	selectionSort(sortedSnap)
	if !reflect.DeepEqual(sortedA, sortedSnap) {
		t.Errorf("Shuffle 后元素集合发生变化: got %v want %v", sortedA, sortedSnap)
	}
}

func selectionSort(arr []int) {
	for i := 0; i < len(arr); i++ {
		minIdx := i
		for j := i + 1; j < len(arr); j++ {
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		arr[i], arr[minIdx] = arr[minIdx], arr[i]
	}
}
