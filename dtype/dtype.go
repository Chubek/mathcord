package dtype

//credits for queue: https://stackoverflow.com/a/55214816
//credits for stack: https://www.educative.io/edpresso/how-to-implement-a-stack-in-golang

import (
	"mathcord/utils"
	"strconv"
)

type Stack []string

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// IsEmpty: check if stack is empty
func (s *Stack) ReturnFirst() string {
	if len(*s) == 0 {
		return ""
	}

	return (*s)[0]
}

// IsEmpty: check if stack is empty
func (s *Stack) ReturnLast() string {
	if len(*s) == 0 {
		return ""
	}

	return (*s)[len(*s)-1]
}

// Push a new value onto the stack
func (s *Stack) Push(str string) {

	*s = append(*s, str) // Simply append the new value to the end of the stack
}

func (s *Stack) DoOp(opOne, opTwo, operator string) {
	var res float64

	switch operator {
	case "+":
		res = utils.ParseToFloat(opOne) + utils.ParseToFloat(opTwo)
	case "-":
		res = utils.ParseToFloat(opOne) - utils.ParseToFloat(opTwo)
	case "*":
		res = utils.ParseToFloat(opOne) * utils.ParseToFloat(opTwo)
	case "/":
		res = utils.ParseToFloat(opOne) / utils.ParseToFloat(opTwo)

	}

	*s = append(*s, strconv.FormatFloat(res, 'f', 4, 64)) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() string {
	if s.IsEmpty() {
		return ""
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.

		return element
	}
}

func NewStack() *Stack {
	return &Stack{}
}

type Queue []string

func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

func (q *Queue) HasAtLeastOne() bool {
	return len(*q) == 1
}

func (q *Queue) Push(xx string, op bool) {
	if !op {
		*q = append(*q, xx)
	} else {
		*q = append([]string{xx}, *q...)
	}

}

func (q *Queue) Pop(op bool) string {
	h := *q
	var el string
	l := len(h)

	if l == 1 {
		el, *q = h[0], []string{}
	} else if !op {
		el, *q = h[0], h[1:l]
	} else {
		el, *q = h[l-1], h[0:l-1]
	}

	return el
}

func NewQueue() *Queue {
	return &Queue{}
}
