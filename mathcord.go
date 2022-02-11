package main

import (
	"fmt"
	"mathcord/sha512"
)

func main() {

	hash := sha512.NewSha512()

	hash.Update("lll.")
	hash.Calculate()
	fmt.Print(len(hash.GetHexDigest()))

}
