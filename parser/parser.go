package parser

import (
	c "mathcord/constants"
	"mathcord/dtype"
	"regexp"
	"strings"
)

func ShuntingYard(ex string) string {
	tokenPatt := regexp.MustCompile(`\b\d+\b|\(|\)|\/|\-|\*|\+`)

	numberPatt := regexp.MustCompile(`\d+(\.\d+)?`)
	lParaPatt := regexp.MustCompile(`\(`)
	rParaPatt := regexp.MustCompile(`\)`)
	opPatt := regexp.MustCompile(`\/|\-|\*|\+`)

	tokens := tokenPatt.FindAllString(ex, -1)

	queueOut := dtype.NewQueue()
	stackOp := dtype.NewStack()

	for _, token := range tokens {
		token = strings.Trim(token, " ")
		if numberPatt.MatchString(token) {
			queueOut.Push(token, false)
		} else if opPatt.MatchString(token) {
			for c.Precedences[stackOp.ReturnLast()] > c.Precedences[token] {
				queueOut.Push(stackOp.Pop(), false)
			}

			stackOp.Push(token)
		} else if rParaPatt.MatchString(token) {

			for !lParaPatt.MatchString(stackOp.ReturnLast()) {
				tok_ := stackOp.Pop()
				queueOut.Push(tok_, false)
			}

			stackOp.Pop()

		} else {
			stackOp.Push(token)
		}
	}

	for !stackOp.IsEmpty() {
		queueOut.Push(stackOp.Pop(), false)
	}

	rpnStack := dtype.NewStack()

	for !queueOut.IsEmpty() {
		token := queueOut.Pop(false)

		if numberPatt.MatchString(token) {
			rpnStack.Push(token)
		} else {
			opOne := rpnStack.Pop()
			opTwo := rpnStack.Pop()

			rpnStack.DoOp(opOne, opTwo, token)
		}

	}

	return rpnStack.Pop()

}
