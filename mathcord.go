package main

import (
	"fmt"
	"mathcord/utils"
)

func main() {
	/*
	hash := sha512.NewSha512()

	hash.Update("Hello World!")
	hash.Calculate()
	fmt.Print(hash.GetHexDigest())
*/
	fmt.Print(utils.BinaryToDecimal(utils.IntegerToBinary(88, 64)))
}
