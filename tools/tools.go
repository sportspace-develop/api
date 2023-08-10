package tools

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var letterRunesNum = []rune("1234567890")

func RandNumRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunesNum[rand.Intn(len(letterRunesNum))]
	}
	return string(b)
}
