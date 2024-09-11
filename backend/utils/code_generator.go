package utils

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID generates a ULID as a unique code
func GenerateULID() string {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(rand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func IsUlid(code string) bool {
	_, err := ulid.ParseStrict(code)
	return err == nil
}
