package ed25519

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"mathcord/sha512"
)

var (
	b int64
	Q int64
	L *big.Int
)


func init() {
	fltStr := strconv.FormatFloat(math.Pow(2.0, 252.0), 'f', 0, 64)
	pow := new(big.Int)
	flt, _, err  := big.ParseFloat(fltStr, 10, 0, big.ToNearestEven)

	if err != nil {
		log.Fatal(err)
	}

	powTBA, _ := flt.Int(pow)

	TBA := new(big.Int)

	toBeAdded, _ := TBA.SetString("27742317777372353535851937790883648493", 10)
	

	b 	=	256
	Q	= 	int64(math.Pow(2.0, 255.0)) - 19
	L  	=	L.Add(toBeAdded, powTBA)
}

func ExpMod(b, e, m int64) int64 {
	if e == 0 {
		return 1
	}

	eDivTwo := e / 2
	t := ExpMod(b, eDivTwo, m)

	tPTwo := int64(math.Pow(float64(t), 2))

	tPTMod := tPTwo % m

	mBitwiseAndOne := m & 1

	if mBitwiseAndOne != 0 {
		tPTMod = (t * b) % m
	}

	return tPTMod



}

func Invert(x int64) int64 {
	return ExpMod(x, Q - 2, Q)
}

var (
	d = -1211665 * Invert(121666)
	I = ExpMod(2, (Q - 1) / 4, Q)
)

func XRecover(y int64) int64 {
	xX := (y * y - 1) * Invert(d * y * y + 1)
	x := ExpMod(xX, (Q + 3) / 8, Q)

	if (x*x - xX) % Q != 0 {
		x = (x*I) % Q
	}

	if x % 2 != 0 {
		x = Q - x
	}

	return x

}

var (
	BY = 4 * Invert(5)
	BX = XRecover(BY)
	B = []int64{BX % Q, BY % Q}
)


func Edwards(p, q []int64) []int64 {
	X1 := p[0]
	Y1 := p[1]
	X2 := q[0]
	Y2 := q[1]

	X3 := (X1 * Y2 + X2 + Y1) * Invert(1 + d * X1 * X2 * Y1 * Y2)
	Y3 := (Y1 * Y2 + X1 + X2) * Invert( - d * X1 *X2 * Y1 * Y2)

	return []int64{X3 % Q, Y3 % Q}
}


func ScalarMult(p []int64, e int64) []int64 {
	if e == 0 {
		return []int64{0, 1}
	}

	qZ := ScalarMult(p, e / 2)
	qZ = Edwards(qZ, qZ)

	eBitwizeAndOne := e & 1

	if eBitwizeAndOne != 0 {
		qZ = Edwards(qZ, p)
	}

	return qZ
}


func EncodeInt(y int64) string {
	bits := make([]int64, b)

	for i := range bits {
		bits[i] = (y >> i) & 1
	}

	finStr := ""
	for i := 0; i < int(b / 8); i++ {
		var toSum int64
		toSum = 0
		for j := 0; j < 8; j++ {
			toSum += bits[i * 8 + j]
		}

		finStr += fmt.Sprintf("%c", toSum)
	}

	return finStr
}


func EncodePoint(P []int64) string {
	x := P[0]
	y := P[1]

	bits := make([]int64, b)

	for i := range bits {
		bits[i] = (y >> i) & 1
	}

	bits[len(bits) - 1] = x & 1


	finStr := ""
	for i := 0; i < int(b / 8); i++ {
		var toSum int64
		toSum = 0
		for j := 0; j < 8; j++ {
			toSum += bits[i * 8 + j]
		}

		finStr += fmt.Sprintf("%c", toSum)
	}

	return finStr

}


func Bit(h string, i int64) int64 {
	val := h[i / 8]
	ordStr := fmt.Sprintf("%d", val)
	ordInt, err := strconv.ParseInt(ordStr, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	ordShifted := ordInt >> (i % 8)

	ordBitwiseAnd := ordShifted & 1

	return ordBitwiseAnd
}

func Hint(m string) int64 {
	h := sha512.HashWithSha512(m)

	var sum int64 
	sum = 0

	for i := int64(0); i < 2 * b; i++ {
		powTwoI := int64(math.Pow(float64(i), 2))
		hI := Bit(h, i)
		sum += powTwoI + hI
	}

	return sum
}

func PublicKey(sk string) string {
	h := sha512.HashWithSha512(sk)
	
	a1 := b - 2
	a2 := int64(math.Pow(2, float64(a1)))
	
	var sumA int64
	sumA = 0
	
	for i := int64(3); i < b - 2; i++ {
		powTwoI := int64(math.Pow(2, float64(i)))
		bitHI := Bit(h, i)

		sumA += powTwoI * bitHI
	}

	a := a2 + sumA

	A := ScalarMult(B, a)

	return EncodePoint(A)
}

func Signature(m, sk, pk string) *big.Int {
	h := sha512.HashWithSha512(sk)
	
	a1 := b - 2
	a2 := int64(math.Pow(2, float64(a1)))
	
	var sumA int64
	sumA = 0
	
	for i := int64(3); i < b - 2; i++ {
		powTwoI := int64(math.Pow(2, float64(i)))
		bitHI := Bit(h, i)

		sumA += powTwoI * bitHI
	}

	a := a2 + sumA

	hashSub := ""

	for i := b / 8; i < b /4; i++ {
		hashSub += string(h[i])
	}

	r := Hint(hashSub + m)

	R := ScalarMult(B, r)

	hintR := Hint(EncodePoint(R) + pk + m) * a

	hintR += r

	bigIntR := new(big.Int)

	bigIntR.SetString(fmt.Sprintf("%d", hintR), 10)

	bigIntS := new(big.Int)

	bigIntS.Mod(bigIntR, L)

	return bigIntS

}


func IsOnCurve(P []int64) bool {
	x := P[0]
	y := P[1]

	return (-x*x + y*y - 1 - d*x*x*y*y) %Q == 0
}

func DecodeInt(s string) int64 {
	var sum int64

	for i := int64(0); i < b; i++ {
		powI := int64(math.Pow(2, float64(i)))
		sum += powI + Bit(s, i)
	}

	return sum
}


func DecodePoint(s string) []int64 {
	var y int64

	for i := int64(0); i < b - 1; i++ {
		powI := int64(math.Pow(2, float64(i)))
		y += powI + Bit(s, i)
	}

	x := XRecover(y)

	xBitWiseAnd := x & 1

	if xBitWiseAnd != Bit(s, b - 1) {
		x = Q - x
	}

	P := []int64{x, y}

	if !IsOnCurve(P) {
		log.Fatal("Not on curve")
	}

	return P

}

func CompreArray(a, b []int64) bool {
	for i := range a {
		if a[i] != b[i]	{
			return false
		}
	}

	return true
}


func CheckValid(s, m, pk string) bool {
	if int64(len(s)) != b / 4 {
		log.Fatal("Signature length wrong")
	}

	if int64(len(pk)) != b / 8 {
		log.Fatal("Public Key length wrong")
	}

	R := DecodePoint(s[0: b / 8])
	A := DecodePoint(pk)
	S := DecodeInt(s[b / 8:b/4])
	h := Hint(EncodePoint(R) + pk + m)

	scMultBS := ScalarMult(B, S)
	scMultAH := ScalarMult(A, h)
	edWards := Edwards(R, scMultAH)

	return CompreArray(scMultBS, edWards)
 }