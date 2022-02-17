package main

import (
	"fmt"
	"mathcord/ed25519"
)

var (
	message   = "Test Messa1e"
	pk        = "TA4x7SlBg7+eN3g27WC6PbNF96a6y1ss+EIuZutAfFU="
	signature = "Q84vlyRosdSfK8fh13UZoh4fstD4waGaAZVkDiFSMPlwAkePf+B9rMAdcTNjYQh0rto6/Lqw89wb+UIA562xAQ=="
)

func main() {

	fmt.Print(ed25519.CheckValid(signature, message, pk))
}
