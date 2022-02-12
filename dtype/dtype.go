package dtype

//credits for queue: https://stackoverflow.com/a/55214816
//credits for stack: https://www.educative.io/edpresso/how-to-implement-a-stack-in-golang

import (
	"math"
	"mathcord/constants"
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

type Chunk [][]uint64

func InitChunk(initVal string, numChunks int) *Chunk {
	chunks := &Chunk{}

	for i := 0; i < numChunks; i++ {
		chunk := make([]uint64, 80)

		for j := i * 1024; j < 1024*(i+1); j += 64 {
			val := initVal[j : j+64]
			chunk[j/64] = utils.BinaryToDecimal(val)

		}
		for i := 16; i < 80; i++ {
			PadWithWords(i, &chunk)

		}

		(*chunks) = append((*chunks), chunk)
	}

	return chunks

}

func DoA(chunk *[]uint64, j int) uint64 {
	el := (*chunk)[j-2]

	elOne := utils.RotateUintRightByNBits(el, 19)
	elTwo := utils.RotateUintRightByNBits(el, 61)
	elThree := utils.ShiftUintRightByNBits(el, 6)

	retA := elOne ^ elTwo ^ elThree

	return retA

}

func DoC(chunk *[]uint64, j int) uint64 {
	el := (*chunk)[j-15]

	elOne := utils.RotateUintRightByNBits(el, 1)
	elTwo := utils.RotateUintRightByNBits(el, 8)
	elThree := utils.ShiftUintRightByNBits(el, 7)

	retC := elOne ^ elTwo ^ elThree

	return retC

}



func PadWithWords(g int, chunk *[]uint64) {
	A := DoA(chunk, g)
	C := DoC(chunk, g)
	B := (*chunk)[g-7]
	D := (*chunk)[g-16]

	added := A + B + C + D

	(*chunk)[g] = added

}

type Sha512Message struct {
	Original  string
	Message   string
	Chunks    *Chunk
	NumChunks int
}

func (message *Sha512Message) InitAndAppendBits() {
	message.Message = utils.StrToBinary(message.Original)
	lengthOriginal := uint64(len(message.Message))

	message.Message += "1"

	lengthDiv := int(math.Ceil(float64(len(message.Original)) / float64((512 - 64))))
	message.Message += strings.Repeat("0", ((1024-64)*lengthDiv - len(message.Message)))

	message.Message += utils.IntegerToBinary(lengthOriginal, 64)

	message.NumChunks = lengthDiv
}

func (message *Sha512Message) GetLength() int {
	return len(message.Message)
}

func NewMessage(str string) *Sha512Message {
	message := &Sha512Message{Original: str}
	message.InitAndAppendBits()
	message.Chunks = InitChunk(message.Message, message.NumChunks)

	return message
}

type Sha512Buffer struct {
	A     uint64
	B     uint64
	C     uint64
	D     uint64
	E     uint64
	F     uint64
	G     uint64
	H     uint64
	APrev uint64
	BPrev uint64
	CPrev uint64
	DPrev uint64
	EPrev uint64
	FPrev uint64
	GPrev uint64
	HPrev uint64
}

func MajVal(A, B, C uint64) uint64 {
	return (A & B) ^ (B & C) ^ (C & A)
}

func ChVal(E, F, G uint64) uint64 {
	return (E & F) ^ (^E & G)
}

func SigmaE(E uint64) uint64 {
	e14 := utils.RotateUintRightByNBits(E, 14)
	e18 := utils.RotateUintRightByNBits(E, 18)
	e41 := utils.RotateUintRightByNBits(E, 41)

	resE :=  e14 ^ e18 ^ e41

	return resE
}

func SigmaA(A uint64) uint64 {
	a28 := utils.RotateUintRightByNBits(A, 28)
	a34 := utils.RotateUintRightByNBits(A, 34)
	a39 := utils.RotateUintRightByNBits(A, 39)

	resA := a28 ^ a34 ^ a39

	return resA
}

func (buff *Sha512Buffer) AddAndSetPrev() {
	buff.A += buff.APrev
	buff.B += buff.BPrev
	buff.C += buff.CPrev
	buff.D += buff.DPrev
	buff.E += buff.EPrev
	buff.F += buff.FPrev
	buff.G += buff.GPrev
	buff.H += buff.HPrev

	buff.APrev = buff.A
	buff.BPrev = buff.B
	buff.CPrev = buff.C
	buff.DPrev = buff.D
	buff.EPrev = buff.E
	buff.FPrev = buff.F
	buff.GPrev = buff.G
	buff.HPrev = buff.H
}

func Rotate(A, B, C, D, E, F, G, H *uint64, K, MS *[]uint64, i *int) {
	messageK := (*MS)[*i]
	kK := (*K)[*i]

	T1 := *H + ChVal(*E, *F, *G) + SigmaE(*E) + messageK + kK
	T2 := SigmaA(*A) + MajVal(*A, *B, *C)

	(*D) = (*D) + T1
	(*H) = T1 + T2

	*i += 1

}

func (buff *Sha512Buffer) DoOneRound(KS, MSU *[]uint64) {
	i := 0
	for j := 0; j < 10; j++ {
		Rotate(&buff.A, &buff.B, &buff.C, &buff.D, &buff.E, &buff.F, &buff.G, &buff.H, KS, MSU, &i)
		Rotate(&buff.H, &buff.A, &buff.B, &buff.C, &buff.D, &buff.E, &buff.F, &buff.G, KS, MSU, &i)
		Rotate(&buff.G, &buff.H, &buff.A, &buff.B, &buff.C, &buff.D, &buff.E, &buff.F, KS, MSU, &i)
		Rotate(&buff.F, &buff.G, &buff.H, &buff.A, &buff.B, &buff.C, &buff.D, &buff.E, KS, MSU, &i)
		Rotate(&buff.E, &buff.F, &buff.G, &buff.H, &buff.A, &buff.B, &buff.C, &buff.D, KS, MSU, &i)
		Rotate(&buff.D, &buff.E, &buff.F, &buff.G, &buff.H, &buff.A, &buff.B, &buff.C, KS, MSU, &i)
		Rotate(&buff.C, &buff.D, &buff.E, &buff.F, &buff.G, &buff.H, &buff.A, &buff.B, KS, MSU, &i)
		Rotate(&buff.B, &buff.C, &buff.D, &buff.E, &buff.F, &buff.G, &buff.H, &buff.A, KS, MSU, &i)

	}

	buff.AddAndSetPrev()

}

func (buff *Sha512Buffer) ProcessBlock(MS *[]uint64) {
	buff.DoOneRound(&constants.Sha512K, MS)
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
	buffer := &Sha512Buffer{A: 0x6a09e667f3bcc908,
		B:     0xbb67ae8584caa73b,
		C:     0x3c6ef372fe94f82b,
		D:     0xa54ff53a5f1d36f1,
		E:     0x510e527fade682d1,
		F:     0x9b05688c2b3e6c1f,
		G:     0x1f83d9abfb41bd6b,
		H:     0x5be0cd19137e2179,
		APrev: 0x6a09e667f3bcc908,
		BPrev: 0xbb67ae8584caa73b,
		CPrev: 0x3c6ef372fe94f82b,
		DPrev: 0xa54ff53a5f1d36f1,
		EPrev: 0x510e527fade682d1,
		FPrev: 0x9b05688c2b3e6c1f,
		GPrev: 0x1f83d9abfb41bd6b,
		HPrev: 0x5be0cd19137e2179}

	return buffer
}
