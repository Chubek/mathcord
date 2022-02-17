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
	Q  *big.Int
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
	L = new(big.Int)

	L.Add(TBA, powTBA)

	bigTwo := big.NewInt(2)
	bigTwoFiveFive := big.NewInt(255)
	bigNineteen := big.NewInt(19)

	bigTwo.Exp(bigTwo, bigTwoFiveFive, nil)
	bigTwo.Sub(bigTwo, bigNineteen)

	Q = bigTwo

	d = big.NewInt(-121665)
	dd := big.NewInt(121666)

	inv := Invert(dd)

	d.Mul(d, inv)

	QI := new(big.Int)

	QI.Sub(Q, big.NewInt(1))
	QI.Div(QI, big.NewInt(4))

	I = ExpMod(big.NewInt(2), QI, Q)

	BY = new(big.Int)

	BY.Mul(big.NewInt(4), Invert(big.NewInt(5)))
	BX = XRecover(BY)

	BXmodQ := new(big.Int)
	BYmodQ := new(big.Int)

	BXmodQ.Mod(BX, Q)
	BYmodQ.Mod(BY, Q)

	B = []*big.Int{BXmodQ, BYmodQ}
}

func ExpMod(b, e, m *big.Int) *big.Int {
	oneBigInt := big.NewInt(1)
	bigZero := big.NewInt(0)

	if e.Cmp(bigZero) == 0 {

		return oneBigInt
	}

	eDivTwo := new(big.Int)

	eDivTwo.Div(e, big.NewInt(2))

	t := ExpMod(b, eDivTwo, m)
	t.Exp(t, big.NewInt(2), nil)
	t.Mod(t, m)

	eBitwiseAndOne := new(big.Int)

	eBitwiseAndOne.And(e, oneBigInt)

	if eBitwiseAndOne.Cmp(bigZero) == 1 {
		t.Mul(t, b)
		t.Mod(t, m)
	}

	return t

}

func Invert(x *big.Int) *big.Int {
	qMinTwo := new(big.Int)

	qMinTwo.Sub(Q, big.NewInt(2))

	res := ExpMod(x, qMinTwo, Q)

	return res
}

func XRecover(y *big.Int) *big.Int {
	bigOne := big.NewInt(1)

	opOne := new(big.Int)

	opOne.Mul(y, y)

	opOne.Sub(opOne, bigOne)

	opTwo := new(big.Int)
	opTwo.Mul(y, y)

	opTwo.Mul(opTwo, d)

	opTwo.Add(opTwo, bigOne)

	invT := Invert(opTwo)

	xX := new(big.Int)

	xX.Mul(opOne, invT)

	qMod := new(big.Int)

	qMod.Add(Q, big.NewInt(3))
	qMod.Div(qMod, big.NewInt(8))

	x := ExpMod(xX, qMod, Q)

	toCompare := new(big.Int)

	toCompare.Mul(x, x)
	toCompare.Sub(toCompare, xX)
	toCompare.Mod(toCompare, Q)

	bigZero := big.NewInt(0)

	if toCompare.Cmp(bigZero) == 1 {
		xI := new(big.Int)

		xI.Mul(x, I)
		xI.Mod(xI, Q)

		x = xI

	}

	xMod := new(big.Int)

	xMod.Mod(x, big.NewInt(2))

	if xMod.Cmp(bigZero) == 1 {
		qSubX := new(big.Int)

		qSubX.Sub(Q, x)

		x = qSubX

	}

	return x

}

func Pass(a string) string {
	return a
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

func ScalarMult(p []*big.Int, e *big.Float) []*big.Int {
	bigIntOne := big.NewInt(1)
	bigIntZero := big.NewInt(0)

	if e.Cmp(new(big.Float).SetInt(bigIntZero)) == 0 {
		return []*big.Int{bigIntZero, bigIntOne}
	}

	eDivTwo := new(big.Float)

	eDivTwo.Quo(e, big.NewFloat(2))

	qZ := ScalarMult(p, eDivTwo)
	qZ = Edwards(qZ, qZ)

	eBitwiseAndOne := new(big.Int)

	eInt, _ := new(big.Int).SetString(e.String(), 10)

	eBitwiseAndOne.And(eInt, bigIntOne)

	if eBitwiseAndOne.Cmp(bigIntZero) == 1 {
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
			lShift := new(big.Int).Lsh(bits[i*8+j], uint(j))

			toSum.Add(toSum, lShift)
		}
		finStr += fmt.Sprintf("%c", toSum.Int64())
	}

	return strconv.QuoteToASCII(finStr)
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

	var finStr = ""

	for i := 0; i < int(b/8); i++ {
		toSum := big.NewInt(0)
		for j := 0; j < 8; j++ {
			shiftLeft := new(big.Int).Lsh(bits[i*8+j], uint(j))
			toSum.Add(toSum, shiftLeft)
		}

		strChar := string(rune(int(toSum.Int64())))

		finStr += strChar
	}

	fmt.Print(finStr, "\n")

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

		mulI := new(big.Int).Mul(powTwoI, big.NewInt(bitHI))

		sum.Add(sum, mulI)
	}

	return sum
}

func PublicKey(sk string) string {
	h := sha512.HashWithSha512(sk)

	fmt.Print(h, "\n")

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

	A := ScalarMult(B, new(big.Float).SetInt(a))

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

	R := ScalarMult(B, new(big.Float).SetInt(r))

	hintR := Hint(EncodePoint(R) + pk + m)

	hintRBig := new(big.Int)
	hintRBig.Mul(hintR, a)

	hintRBig.Add(hintRBig, r)

	bigIntR := new(big.Int)

	bigIntR.SetString(fmt.Sprintf("%d", hintR.Int64()), 10)

	bigIntS := new(big.Int)

	bigIntS.Mod(bigIntR, L)

	return bigIntS

}

func IsOnCurve(P []*big.Int) bool {
	x := P[0]
	y := P[1]

	xMulx := new(big.Int)
	yMuly := new(big.Int)

	xMulx.Mul(x, x)
	yMuly.Mul(y, y)

	xMulx.Neg(xMulx)

	xMulx.Add(xMulx, yMuly)

	xMulx_ := new(big.Int)
	yMuly_ := new(big.Int)
	xMulx_.Mul(x, x)
	yMuly_.Mul(y, y)

	xMulx_.Mul(yMuly_, xMulx_)
	xMulx_.Mul(d, xMulx_)

	xMulx_.Add(xMulx_, big.NewInt(1))

	xMulx_.Neg(xMulx_)

	xMulx.Add(xMulx, xMulx_)

	xMulx.Mod(xMulx, Q)

	return xMulx.Cmp(big.NewInt(0)) == 0
}

func DecodeInt(s string) *big.Int {
	sum := big.NewInt(0)
	twoBig := big.NewInt(2)
	for i := int64(0); i < b; i++ {
		iBig := big.NewInt(i)

		powI := new(big.Int)

		powI.Exp(twoBig, iBig, nil)

		sumBig := new(big.Int)
		bInt := big.NewInt(Bit(s, i))

		sumBig.Mul(powI, bInt)

		sum.Add(sum, sumBig)
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

		mul := new(big.Int)

		mul.Mul(bitI, powI)

		y.Add(y, mul)
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
	fmt.Print(a, "   ", b, "\n")

	for i := range a {
		if a[i].Cmp(b[i]) == 1 {
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

	scMultBS := ScalarMult(B, new(big.Float).SetInt(S))
	scMultAH := ScalarMult(A, new(big.Float).SetInt(h))
	edWards := Edwards(R, scMultAH)

	return CompreArray(scMultBS, edWards)
}
