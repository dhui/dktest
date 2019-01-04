package dktest

import (
	"math/rand"
	"time"
)

const (
	chars               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	containerNamePrefix = "dktest_"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randString(n uint) string {
	if n == 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func genContainerName() string { return containerNamePrefix + randString(10) }
