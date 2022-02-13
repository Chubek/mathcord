package ed25519

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"mathcord/sha512"
	"mathcord/utils"
	"strconv"
)

var (
	b  int64
	L  *big.Int
	d  *big.Int
	I  *big.Int
	BY *big.Int
	BX *big.Int
	B  []*big.Int
	Q *big.Int
)

func init() {
	fltStr := strconv.FormatFloat(math.Pow(2.0, 252.0), 'f', 0, 64)
	pow := new(big.Int)
	flt, _, err := big.ParseFloat(fltStr, 10, 0, big.ToNearestEven)

	if err != nil {
		log.Fatal(err)
	}

	powTBA, _ := flt.Int(pow)

	TBA := new(big.Int)

	_, success := TBA.SetString("27742317777372353535851937790883648493", 10)

	if !success {
		log.Fatal("Set failed")
	}

	b = 256
	QS := int64(math.Pow(2.0, 255.0)) - 19
	L = new(big.Int)

	L.Add(TBA, powTBA)

	Q = big.NewInt(QS)

	d.SetInt64(-1211665)
	d.Mul(d, Invert(big.NewInt(121666)))
	
	QI := new(big.Int)

	QI.Sub(Q, big.NewInt(1))
	QI.Div(QI, big.NewInt(4))
	
	I = ExpMod(big.NewInt(2), QI, Q)

	BY.Mul(big.NewInt(4), Invert(big.NewInt(5)))
	BX = XRecover(BY)

	BXmodQ := new(big.Int)
	BYmodQ := new(big.Int)

	BXmodQ.Mod(BX, Q)
	BYmodQ.Mod(BY, Q)

	B = []*big.Int{BXmodQ, BYmodQ}
}

func ExpMod(b, e, m *big.Int) *big.Int {
	oneBigInt  := big.NewInt(1)
	
	if e == big.NewInt(0) {
		return oneBigInt
	}

	eDivTwo := new(big.Int)

	eDivTwo.Div(e, big.NewInt(2))
	
	t := ExpMod(b, eDivTwo, m)

	tPTwo := new(big.Int)

	tPTwo.Exp(t, big.NewInt(2), nil)

	tPTMod := new(big.Int)

	tPTMod.Mod(tPTwo, m)

	mBitwiseAndOne := new(big.Int)

	mBitwiseAndOne.And(m, oneBigInt)

	if mBitwiseAndOne != big.NewInt(0) {
		tPTMod.Mul(t, b)
		tPTMod.Mod(tPTMod, m)
	}

	return tPTMod

}

func Invert(x *big.Int) *big.Int {
	qMinTwo := new(big.Int)

	qMinTwo.Sub(Q, big.NewInt(2))

	return ExpMod(x, qMinTwo, Q)
}

func XRecover(y *big.Int) *big.Int {
	yMuly := new(big.Int)

	yMuly.Mul(y, y)

	yMuly.Sub(yMuly, big.NewInt(1))


	toInv := new(big.Int)
	added := new(big.Int)

	added.Mul(y, y)
	added.Add(added, big.NewInt(1))

	toInv.Mul(d, added)

	yMulyy := new(big.Int)

	yMulyy.Mul(y, y)
	yMulyy.Add(yMuly, big.NewInt(1))

	inverted := Invert(toInv)

	xX := new(big.Int)

	xX.Mul(yMuly, inverted)

	qMod := new(big.Int)

	qMod.Add(Q, big.NewInt(3))
	qMod.Div(qMod, big.NewInt(8))

	x := ExpMod(xX, qMod, Q)

	toCompare := new(big.Int)

	toCompare.Mul(x, x)
	toCompare.Sub(toCompare, xX)
	toCompare.Mod(toCompare, Q)


	if toCompare != big.NewInt(0)  {
		x.Mul(x, I)
		x.Mod(x, Q)		
	}

	x.Mod(x, big.NewInt(2))

	if x != big.NewInt(0) {
		x.Sub(Q, x)
	}

	return x

}

func Edwards(p, q []*big.Int) []*big.Int {
	X1 := p[0]
	Y1 := p[1]
	X2 := q[0]
	Y2 := q[1]

	opOne := new(big.Int)
	opTwo := new(big.Int)
	opThree := big.NewInt(1)
	opFour := d

	opOne.Mul(X1, Y1)
	opOne.Add(X2, Y2)

	opTwo.Mul(Y1, Y2)
	opTwo.Add(X1, X2)
	
	opThree.Add(opThree, d)
	opThree.Mul(X1, X2)
	opThree.Mul(Y1, Y2)

	opFour.Neg(opFour)
	opFour.Mul(X1, X2)
	opFour.Mul(Y1, Y2)

	opThreeInverted := Invert(opThree)
	opFourInverted := Invert(opFour)

	opThreeInverted.Mul(opThreeInverted, opOne)
	opFourInverted.Mul(opFourInverted, opTwo)

	X3 := new(big.Int)
	Y3 := new(big.Int)

	X3.Mod(opThreeInverted, Q)
	Y3.Mod(opFourInverted, Q)

	return []*big.Int{X3, Y3}
}

func ScalarMult(p []*big.Int, e *big.Int) []*big.Int {
	bigIntOne := big.NewInt(1)
	bigIntZero := big.NewInt(0)
	
	if e == big.NewInt(0) {
		return []*big.Int{bigIntZero, bigIntOne}
	}

	eDivTwo := new(big.Int)

	eDivTwo.Div(e, big.NewInt(2))

	qZ := ScalarMult(p, eDivTwo)
	qZ = Edwards(qZ, qZ)

	eBitwizeAndOne := new(big.Int)

	eBitwizeAndOne.And(e, bigIntOne)

	if eBitwizeAndOne != big.NewInt(0) {
		qZ = Edwards(qZ, p)
	}

	return qZ
}

func EncodeInt(y *big.Int) string {
	bits := make([]*big.Int, b)
	bigIntOne := big.NewInt(1)
	for i := range bits {
		res := new(big.Int)

		res.Rsh(y, uint(i))
		res.And(res, bigIntOne)

		bits[i] = res
	}

	finStr := ""
	for i := 0; i < int(b/8); i++ {
		toSum := big.NewInt(0)
		for j := 0; j < 8; j++ {
			toSum.Add(toSum, bits[i*8+j])
		}

		finStr += fmt.Sprintf("%c", toSum)
	}

	return finStr
}

func EncodePoint(P []*big.Int) string {
	x := P[0]
	y := P[1]

	bits := make([]*big.Int, b)
	bigIntOne := big.NewInt(1)
	for i := range bits {
		res := new(big.Int)

		res.Rsh(y, uint(i))
		res.And(res, bigIntOne)

		bits[i] = res
	}

	resX := new(big.Int)

	resX.And(x, big.NewInt(1))

	bits[len(bits)-1] = resX

	finStr := ""
	for i := 0; i < int(b/8); i++ {
		toSum := big.NewInt(0)
		for j := 0; j < 8; j++ {
			toSum.Add(toSum, bits[i*8+j])
		}

		finStr += fmt.Sprintf("%c", toSum)
	}

	return finStr

}

func Bit(h string, i int64) int64 {
	val := h[i/8]
	ordStr := fmt.Sprintf("%d", val)
	ordInt, err := strconv.ParseInt(ordStr, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	ordShifted := ordInt >> (i % 8)

	ordBitwiseAnd := ordShifted & 1

	return ordBitwiseAnd
}

func Hint(m string) *big.Int {
	h := sha512.HashWithSha512(m)

	sum := big.NewInt(0)

	for i := int64(0); i < 2*b; i++ {
		powTwoI := new(big.Int)
		powTwoI.Exp(big.NewInt(2), big.NewInt(i), nil)
		bitHI := Bit(h, i)

		sum.Add(powTwoI, big.NewInt(bitHI))
	}

	return sum
}

func PublicKey(sk string) string {
	h := sha512.HashWithSha512(sk)

	a1 := new(big.Int)
	a1.Sub(big.NewInt(b), big.NewInt(2))

	a2 := new(big.Int)

	a2.Exp(big.NewInt(2), a1, nil)

	sumA := big.NewInt(0)

	for i := int64(3); i < b-2; i++ {
		powTwoI := new(big.Int)
		powTwoI.Exp(big.NewInt(2), big.NewInt(i), nil)
		bitHI := Bit(h, i)

		sumA.Add(powTwoI, big.NewInt(bitHI))
	}

	a := new(big.Int)

	a.Add(sumA, a2)

	A := ScalarMult(B, a)

	return EncodePoint(A)
}

func Signature(m, sk, pk string) *big.Int {
	h := sha512.HashWithSha512(sk)

	a1 := new(big.Int)
	a1.Sub(big.NewInt(b), big.NewInt(2))

	a2 := new(big.Int)

	a2.Exp(big.NewInt(2), a1, nil)

	sumA := big.NewInt(0)

	for i := int64(3); i < b-2; i++ {
		powTwoI := new(big.Int)
		powTwoI.Exp(big.NewInt(2), big.NewInt(i), nil)
		bitHI := Bit(h, i)

		sumA.Add(powTwoI, big.NewInt(bitHI))
	}

	a := new(big.Int)

	a.Add(sumA, a2)


	hashSub := ""

	for i := b / 8; i < b/4; i++ {
		hashSub += string(h[i])
	}

	r := Hint(hashSub + m)

	R := ScalarMult(B, r)

	hintR := Hint(EncodePoint(R)+pk+m) 

	hintRBig := new(big.Int)
	hintRBig.Mul(hintR, a)

	hintRBig.Add(hintRBig, r)

	bigIntR := new(big.Int)

	bigIntR.SetString(fmt.Sprintf("%d", hintR), 10)

	bigIntS := new(big.Int)

	bigIntS.Mod(bigIntR, L)

	return bigIntS

}

func IsOnCurve(P []*big.Int) bool {
	x := P[0]
	y := P[1]

	res := new(big.Int)

	res.Neg(x)
	res.Mul(res, x)

	mult := new(big.Int)

	mult.Mul(y, y)

	res.Add(res, mult)

	res.Sub(res, big.NewInt(-1))

	multNew := new(big.Int)

	multNew.Mul(d, x)
	multNew.Mul(multNew, x)
	multNew.Mul(multNew, y)
	multNew.Mul(multNew, y)

	multNew.Neg(multNew)

	res.Add(res, multNew)

	mod := new(big.Int)

	mod.Mod(res, Q)

	return mod == big.NewInt(0)
}

func DecodeInt(s string) *big.Int {
	sum := big.NewInt(0)
	twoBig := big.NewInt(2)
	for i := int64(0); i < b; i++ {
		iBig := big.NewInt(i)
		
		powI := new(big.Int)

		powI.Exp(twoBig, iBig, nil)

		sum.Add(powI, big.NewInt(Bit(s, i)))
	}


	return sum
}

func DecodePoint(s string) []*big.Int {
	y := new(big.Int)
	twoBig := big.NewInt(2)
	for i := int64(0); i < b-1; i++ {
		powI := new(big.Int)
		iBig := big.NewInt(i)
		
		powI.Exp(twoBig, iBig, nil)
		bitI := big.NewInt(Bit(s, i))
		y.Mul(bitI, powI)
	}


	x := XRecover(y)

	xBitWiseAnd := new(big.Int)

	xBitWiseAnd.And(x, big.NewInt(1))

	if xBitWiseAnd != big.NewInt(Bit(s, b-1)) {
		x.Sub(Q, x)
	}

	P := []*big.Int{x, y}

	if !IsOnCurve(P) {
		log.Fatal("Not on curve")
	}

	return P

}

func CompreArray(a, b []*big.Int) bool {
	for i := range a {
		if a[i] != b[i] {
			fmt.Print(a[i], " was unequal to ", b[i], "\n")
			return false
		}
	}

	return true
}

func CheckValid(sEnc, m, pkEnc string) bool {
	s := utils.DecodeBase64(sEnc)
	pk := utils.DecodeBase64(pkEnc)

	if int64(len(s)) != b/4 {
		log.Fatal("Signature length wrong")
	}

	if int64(len(pk)) != b/8 {
		log.Fatal("Public Key length wrong")
	}

	R := DecodePoint(s[0 : b/8])
	A := DecodePoint(pk)
	S := DecodeInt(s[b/8 : b/4])
	h := Hint(EncodePoint(R) + pk + m)

	scMultBS := ScalarMult(B, S)
	scMultAH := ScalarMult(A, h)
	edWards := Edwards(R, scMultAH)

	return CompreArray(scMultBS, edWards)
}
