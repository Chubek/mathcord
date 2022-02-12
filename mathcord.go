package main

import (
	"fmt"
	"mathcord/sha512"
	"strings"
)

func main() {

	hash := sha512.NewSha512()

	hash.Update(strings.Repeat(".", 512))
	hash.Calculate()
	fmt.Print(hash.GetHexDigest())

}
