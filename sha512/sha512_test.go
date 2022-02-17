package sha512_test

import (
	"fmt"
	"mathcord/sha512"
	"testing"
)

func TestSha512Hash_Update(t *testing.T) {
	hash := sha512.NewSha512()

	hash.Update(".")
	hash.Calculate()

	fmt.Println(hash.GetHexDigest())

}
