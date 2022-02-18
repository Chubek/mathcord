package utils_test

import (
	"mathcord/utils"
	"testing"
)

func TestByteFromHex(t *testing.T) {
	testName := "Test Hex"

	t.Run(testName, func(t *testing.T) {
		resHex := utils.HexToByte("44a3b745287713fb9a6db47ee5ef47d34c34eb7f3ff4f69e19db6ee407df5a73")

		t.Log(resHex)
	})
}
