package sha512_test

import (
	"fmt"
	"mathcord/sha512"
	"mathcord/utils"
	"testing"
)

func TestSha512Hash_Update(t *testing.T) {
	hash := sha512.NewSha512()

	hash.Update([]byte{67, 206, 47, 151, 36, 104, 177, 212, 159, 43, 199, 225, 215, 117, 25, 162, 30, 31, 178, 208, 248, 193, 161, 154, 1, 149, 100, 14, 33, 82, 48, 249, 76, 14, 49, 237, 41, 65, 131, 191, 158, 55, 120, 54, 237, 96, 186, 61, 179, 69, 247, 166, 186, 203, 91, 44, 248, 66, 46, 102, 235, 64, 124, 85, 84, 101, 115, 116, 32, 77, 101, 115, 115, 97, 103, 101})
	hash.Calculate()

	fmt.Println(utils.ToAscii(hash.GetDigest()))

}
