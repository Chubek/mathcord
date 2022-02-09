package main

import (
	"fmt"
	"mathcord/parser"
)

func main() {
	fmt.Print("Result is ", parser.ShuntingYard("4 / 4 + (4 + 4) + 10"), "\n")


}
