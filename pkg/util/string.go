package util

import (
	"math/rand"
)

var randx = rand.NewSource(42)

/*
	Feature #A:
		The string generator generates only hex values now.
*/

// RandString returns a random string of length n.
func RandHexString(n int) string {
	const hexBytes = "0123456789abcdef" // now contains only hexadecimal characters
	const (
		letterIdxBits = 4                    // Decreasing to 4 bits in order to represent 16 possible values
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, randx.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randx.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(hexBytes) {
			b[i] = hexBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
