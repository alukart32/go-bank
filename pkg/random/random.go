package random

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Int64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func String(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func GetString(strings ...string) string {
	return strings[rand.Intn(len(strings))]
}
