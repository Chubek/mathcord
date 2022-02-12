package main

import (
	"fmt"
	"mathcord/sha512"
)

func main() {

	hash := sha512.NewSha512()

	hash.Update("GeeksForGeeks")
	hash.Calculate()
	fmt.Print(hash.GetHexDigest())

}
