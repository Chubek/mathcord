package dtype

//credits for queue: https://stackoverflow.com/a/55214816
//credits for stack: https://www.educative.io/edpresso/how-to-implement-a-stack-in-golang

import (
	"math"
	"mathcord/utils"
	"strconv"
	"strings"
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

type Chunk [][]string

func InitChunk(initVal string) *Chunk {
	chunks := &Chunk{}

	numChunks := len(initVal) / 512
	for i := 0; i < numChunks; i++ {
		chunk := make([]string, 64)

		for j := 0; j < 16; j++ {
			chunk[j] = initVal[j:j + 32]

		}

		for i := 16; i < 64; i++ {
			PadWithWords(i, &chunk)

		}

		(*chunks) = append((*chunks), chunk)
	}

	return chunks

}

func DoA(chunk *[]string, j int) string {
	el := (*chunk)[j-2]

	elOne := utils.RotateStringRightByNBits(el, 19)
	elTwo := utils.RotateStringRightByNBits(el, 61)
	elThree := utils.RotateStringRightByNBits(el, 6)

	ret := utils.XorThree(elOne, elTwo, elThree)

	return ret

}

func DoC(chunk *[]string, j int) string {
	el := (*chunk)[j-2]

	elOne := utils.RotateStringRightByNBits(el, 1)
	elTwo := utils.RotateStringRightByNBits(el, 8)
	elThree := utils.RotateStringRightByNBits(el, 7)

	ret := utils.XorThree(elOne, elTwo, elThree)

	return ret

}




func PadWithWords(g int, chunk *[]string) {
	A := DoA(chunk, g)
	C := DoC(chunk, g)
	B := (*chunk)[g-7]
	D := (*chunk)[g-16]

	*chunk = append(*chunk, utils.AddFour(A, B, C, D))

}

type Sha512Message struct {
	Original 	string
	Message  	string
	Chunks   	*Chunk
	NumChunks 	int
}

func (message *Sha512Message) InitAndAppendBits() {
	message.Message = utils.StrToBinary(message.Original)
	message.Message += "1"

	lengthDiv := float64(len(message.Original)) / float64((512 - 64))
	message.Message += strings.Repeat("0", len(message.Original)-int(math.Ceil(lengthDiv)))

	message.Message += utils.IntegerToBinary(uint(len(message.Original)), 64)

	message.NumChunks = len(message.Message) / 512
}

func (message *Sha512Message) GetLength() int {
	return len(message.Message)
}

func NewMessage(str string) *Sha512Message {
	message := &Sha512Message{Original: str}
	message.InitAndAppendBits()
	message.Chunks = InitChunk(message.Message)

	return message
}

type Sha512Buffer struct {
	A uint
	B uint
	C uint
	D uint
	E uint
	F uint
	G uint
	H uint
}

func (buff *Sha512Buffer) MajorVal() uint {
	return (buff.A & buff.B) ^ (buff.B & buff.C) ^ (buff.C & buff.A)
}

func (buff *Sha512Buffer) ChVal() uint {
	return (buff.E & buff.F) ^ (-buff.E & buff.G)
}

func (buff *Sha512Buffer) SigmaE() uint {
	e14 := utils.RotateUintRightByNBits(buff.E, 14)
	e18 := utils.RotateUintRightByNBits(buff.E, 18)
	e41 := utils.RotateUintRightByNBits(buff.E, 41)

	return e14 ^ e18 ^ e41
}

func (buff *Sha512Buffer) SigmaA() uint {
	a28 := utils.RotateUintRightByNBits(buff.A, 28)
	a34 := utils.RotateUintRightByNBits(buff.A, 34)
	a39 := utils.RotateUintRightByNBits(buff.A, 39)

	return a28 ^ a34 ^ a39
}

func (buff *Sha512Buffer) ProcessBlock(K uint, MS string) {
	M := utils.BinaryToDecimal(MS)

	chVal := buff.ChVal()
	majVal := buff.MajorVal()
	hVal := chVal + buff.SigmaE() + K + M
	buff.H = buff.G
	buff.G = buff.F
	buff.F = buff.E
	dVal := hVal + buff.D
	buff.E = dVal
	buff.D = buff.C
	buff.C = buff.B
	buff.B = buff.A
	buff.A = buff.SigmaA() + majVal + hVal
}

func (buff *Sha512Buffer) ToHexaDecimal() string {
	A := utils.DecimalToHex(buff.A)
	B := utils.DecimalToHex(buff.B)
	C := utils.DecimalToHex(buff.C)
	D := utils.DecimalToHex(buff.D)
	E := utils.DecimalToHex(buff.E)
	F := utils.DecimalToHex(buff.F)
	G := utils.DecimalToHex(buff.E)
	H := utils.DecimalToHex(buff.H)

	return A + B + C + D + E + F + G + H

}

func NewBuffer() *Sha512Buffer {
	buffer := &Sha512Buffer{0x6a09e667f3bcc908,
		0xbb67ae8584caa73b,
		0x3c6ef372fe94f82b,
		0xa54ff53a5f1d36f1,
		0x510e527fade682d1,
		0x9b05688c2b3e6c1f,
		0x1f83d9abfb41bd6b,
		0x5be0cd19137e2179}

	return buffer
}
