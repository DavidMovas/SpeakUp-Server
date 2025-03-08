package helpers

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

const (
	maxRand = 9999
	minRand = 1000
)

func GenerateRandomUsername(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	return fmt.Sprintf("%s_%d", name, rand.IntN(maxRand-minRand)+minRand)
}
