package util

import (
	"crypto/sha1"
	"fmt"
)

func Sha1(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))

	return fmt.Sprintf("%+x", hash.Sum(nil))
}
