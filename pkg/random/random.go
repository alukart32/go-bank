package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Int64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func String(n int) []byte {
	array := make([]byte, n)
	k := len(alphabet)

	for i := 0; i < n; i++ {
		array[i] = alphabet[rand.Intn(k)]
	}

	return array
}

func GetString(strings ...string) string {
	return strings[rand.Intn(len(strings))]
}
