// Package random provides tools for some random string
package random

import (
	"math/rand"
	"time"
)

const numberRunes = "0123456789"

//RandomNumString 返回一个指定长度的数字字符串
func RandomNumString(length int) string {
	if length < 0 {
		return ""
	}
	s := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range s {
		// 10 是阿拉伯数字个数
		//TODO: 这个生成方式太浪费了
		s[i] = numberRunes[rand.Intn(10)]
	}
	return string(s)
}
