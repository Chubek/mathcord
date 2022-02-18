package ed25519_test

import (
	"mathcord/ed25519"
	"testing"
)

var (
	message   = "Test Message"
	pk        = "TA4x7SlBg7+eN3g27WC6PbNF96a6y1ss+EIuZutAfFU="
	signature = "Q84vlyRosdSfK8fh13UZoh4fstD4waGaAZVkDiFSMPlwAkePf+B9rMAdcTNjYQh0rto6/Lqw89wb+UIA562xAQ=="
)

func TestCheckIsValid(t *testing.T) {
	testName := "Test Validator"

	t.Run(testName, func(t *testing.T) {
		boolValid := ed25519.CheckValid(signature, message, pk)
		if !boolValid {
			t.Error("Test got it")
		}
	})
}
