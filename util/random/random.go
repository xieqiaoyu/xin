// Package random provides tools for some random string
package random

import (
	"math/rand"
	"strings"
	"time"
)

const (
	numberRunes   = "0123456789"
	numberIdxBits = 4
	numberIdxMask = 1<<numberIdxBits - 1
)

//RandomNumString 返回一个指定长度的数字字符串
func RandomNumString(length int) string {
	if length <= 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(length)
	var dice int64
	writeLen := 0
	for writeLen < length {
		if dice == 0 {
			dice = rand.Int63()
		}
		idx := int(dice & numberIdxMask)
		if idx < 9 {
			sb.WriteByte(numberRunes[idx])
			writeLen++
		}
		dice >>= numberIdxBits
	}
	return sb.String()
}
