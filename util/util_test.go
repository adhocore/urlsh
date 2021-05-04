package util

import (
	"math/rand"
	"testing"
)

func TestRandomString(t *testing.T) {
	t.Run("random string length", func(t *testing.T) {
		expect := rand.Intn(9) + 1
		actual := len(RandomString(expect))

		if expect != actual {
			t.Errorf("output must contain %v chars, not %v", expect, actual)
		}
	})

	t.Run("random string 0 length", func(t *testing.T) {
		if 0 != len(RandomString(0)) {
			t.Errorf("output must be empty string for 0 len")
		}
	})
}
