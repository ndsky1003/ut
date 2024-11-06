package ut

import (
	"fmt"

	"github.com/ndsky1003/buffer"
)

func base_getValue(c byte) int {
	var d byte
	switch {
	case '0' <= c && c <= '9':
		d = c - '0'
	case 'a' <= c && c <= 'z':
		d = c - 'a' + 10
	case 'A' <= c && c <= 'Z':
		d = c - 'A' + 36
	default:
		return -1
	}
	return int(d)
}

func base_getChar(d int) byte {
	var c byte
	switch {
	case 0 <= d && d <= 9:
		c = byte(d) + '0'
	case 10 <= d && d <= 35:
		c = byte(d) + 'a' - 10
	case 36 <= d && d <= 61:
		c = byte(d) + 'A' - 36
	}
	return c
}

// 将一个整数转换成指定进制的字符串
func Base(num int, base int) string {
	if base > 62 {
		base = 62
	}
	if base < 2 {
		base = 2
	}
	if num == 0 {
		return "0"
	}

	buf := buffer.Get()
	defer buf.Release()
	for num > 0 {
		remainder := num % base
		num /= base
		b := base_getChar(remainder)
		if err := buf.WriteByte(b); err != nil {
			panic(err)
		}
	}
	data := buf.Bytes()
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return buf.String()
}

// 将一个字符串转换成指定进制的整数
func Parse(encoded string, base int) (int, error) {
	num := 0
	for _, ch := range encoded {
		index := base_getValue(byte(ch))
		if index == -1 {
			return 0, fmt.Errorf("invalid character %c in encoded string", ch)
		}
		num = num*base + index
	}
	return num, nil
}
