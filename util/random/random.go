// Package random provides tools for some random string
package random

import (
	cryptoRand "crypto/rand"
	"math/rand"
	"strings"
	"time"
)

const (
	numberRunes   = "0123456789"
	numberIdxBits = 4
	numberIdxMask = 1<<numberIdxBits - 1

	letterRunes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
)

//NumString generate random string with number letters
func NumString(length int) string {
	if length <= 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(length)
	var dice int64
	writeLen := 0
	rand.Seed(time.Now().UnixNano())
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

//String generate random string with letters and "_-"
func String(length int) string {
	if length <= 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(length)
	var dice int64
	writeLen := 0
	rand.Seed(time.Now().UnixNano())

	leftLen := 0
	for writeLen < length {
		if leftLen < letterIdxBits {
			dice = rand.Int63()
			leftLen = 64
		}
		idx := int(dice & letterIdxMask)
		sb.WriteByte(letterRunes[idx])
		writeLen++
		dice >>= letterIdxBits
		leftLen -= letterIdxBits
	}
	return sb.String()
}

//Bytes generate random bytes
func Bytes(length int) []byte {
	k := make([]byte, length)
	if _, err := cryptoRand.Read(k); err != nil {
		//TODO: does err should trigger panic ?
		return nil
	}
	return k
}
