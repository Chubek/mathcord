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
		log.Println(err)
	}

	powTBA, _ := flt.Int(pow)

	TBA := new(big.Int)

	_, success := TBA.SetString("27742317777372353535851937790883648493", 10)

	if !success {
		log.Println("Set failed")
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

	if eBitwiseAndOne.Cmp(bigZero) != 0 {
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

	if toCompare.Cmp(bigZero) != 0 {
		xI := new(big.Int)

		xI.Mul(x, I)
		xI.Mod(xI, Q)

		x = xI

	}

	xMod := new(big.Int)

	xMod.Mod(x, big.NewInt(2))

	if xMod.Cmp(bigZero) != 0 {
		qSubX := new(big.Int)

		qSubX.Sub(Q, x)

		x = qSubX

	}

	return x

}

func Edwards(p, q []*big.Int) []*big.Int {
	X1 := p[0]
	Y1 := p[1]
	X2 := q[0]
	Y2 := q[1]

	X3 := new(big.Int).Mul(new(big.Int).Add(new(big.Int).Mul(X1, Y2), new(big.Int).Mul(X2, Y1)), Invert(new(big.Int).Add(big.NewInt(1), new(big.Int).Mul(d, new(big.Int).Mul(new(big.Int).Mul(X1, X2), new(big.Int).Mul(Y1, Y2))))))
	Y3 := new(big.Int).Mul(new(big.Int).Add(new(big.Int).Mul(Y1, Y2), new(big.Int).Mul(X1, X2)), Invert(new(big.Int).Sub(big.NewInt(1), new(big.Int).Mul(d, new(big.Int).Mul(new(big.Int).Mul(X1, X2), new(big.Int).Mul(Y1, Y2))))))

	X3.Mod(X3, Q)
	Y3.Mod(Y3, Q)

	return []*big.Int{X3, Y3}
}

func ScalarMult(p []*big.Int, e *big.Int) []*big.Int {
	bigIntZero := big.NewInt(0)
	bigIntOne := big.NewInt(1)

	if e.Cmp(bigIntZero) == 0 {
		ret := make([]*big.Int, 2)

		ret[0] = bigIntZero
		ret[1] = bigIntOne

		return ret
	}

	eDivTwo := new(big.Int)

	eDivTwo.Div(e, big.NewInt(2))

	qZ := ScalarMult(p, eDivTwo)
	qZ = Edwards(qZ, qZ)

	eBitwiseAndOne := new(big.Int)

	eBitwiseAndOne.And(e, bigIntOne)

	if eBitwiseAndOne.Cmp(bigIntZero) != 0 {
		qZ = Edwards(qZ, p)
	}

	return qZ
}

func EncodeInt(y *big.Int) []byte {
	bits := make([]*big.Int, b)
	bigIntOne := big.NewInt(1)
	for i := range bits {
		res := new(big.Int)

		res.Rsh(y, uint(i))
		res.And(res, bigIntOne)

		bits[i] = res
	}

	var finStr []byte
	for i := 0; i < int(b/8); i++ {
		toSum := big.NewInt(0)
		for j := 0; j < 8; j++ {
			lShift := new(big.Int).Lsh(bits[i*8+j], uint(j))

			toSum.Add(toSum, lShift)
		}
		finStr = append(finStr, toSum.Bytes()...)
	}

	return finStr
}

func EncodePoint(P []*big.Int) []byte {
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

	var finStr []byte
	for i := 0; i < int(b/8); i++ {
		toSum := big.NewInt(0)
		for j := 0; j < 8; j++ {
			shiftLeft := new(big.Int).Lsh(bits[i*8+j], uint(j))
			toSum.Add(toSum, shiftLeft)
		}

		finStr = append(finStr, toSum.Bytes()...)

	}

	return finStr

}

func Bit(h []byte, i int64) int64 {
	ordInt := int64(h[i/8])

	ordShifted := ordInt >> (i % 8)

	ordBitwiseAnd := ordShifted & 1

	return ordBitwiseAnd
}

func Hint(m []byte) *big.Int {
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

func PublicKey(sk []byte) []byte {
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

func Signature(m, sk, pk []byte) *big.Int {
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

	hashSub := make([]byte, b/4)

	for i := b / 8; i < b/4; i++ {
		hashSub[i] = h[i]
	}

	var hSM []byte

	hSM = append(hSM, hashSub...)
	hSM = append(hSM, m...)

	r := Hint(hSM)

	R := ScalarMult(B, r)

	encP := EncodePoint(R)

	var hRR []byte

	hRR = append(hRR, encP...)
	hRR = append(hRR, pk...)
	hRR = append(hRR, m...)

	hintR := Hint(hRR)

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

func DecodeInt(s []byte) *big.Int {
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

func DecodePoint(s []byte) []*big.Int {
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

	if xBitWiseAnd.Cmp(big.NewInt(Bit(s, b-1))) != 0 {
		nn := new(big.Int).Sub(Q, x)
		x = nn
	}

	P := []*big.Int{x, y}

	if !IsOnCurve(P) {
		log.Println("Not on curve")
	}

	return P

}

func CompreArray(a, b []*big.Int) bool {
	a1 := a[0].String()
	a2 := a[1].String()
	b1 := b[0].String()
	b2 := b[1].String()

	if a1 != b1 || a2 != b2 {
		return false
	}

	return true
}

func CheckValid(mStr, sEnc, pkEnc string) bool {
	s := []byte(utils.DecodeBase64(sEnc))
	pk := []byte(utils.DecodeBase64(pkEnc))
	m := []byte(mStr)

	if int64(len(s)) != b/4 {
		log.Println("Signature length wrong")
	}

	if int64(len(pk)) != b/8 {
		log.Println("Public Key length wrong")
	}

	R := DecodePoint(s[0 : b/8])

	A := DecodePoint(pk)

	S := DecodeInt(s[b/8 : b/4])

	encR := EncodePoint(R)

	var hRR []byte

	hRR = append(hRR, encR...)
	hRR = append(hRR, pk...)
	hRR = append(hRR, m...)

	h := Hint(hRR)

	scMultBS := ScalarMult(B, S)
	scMultAH := ScalarMult(A, h)
	edWards := Edwards(R, scMultAH)

	return CompreArray(scMultBS, edWards)
}

func CheckValidBytes(m, s, pk []byte) bool {

	if int64(len(s)) != b/4 {
		log.Println("Signature length wrong")
	}

	if int64(len(pk)) != b/8 {
		log.Println("Public Key length wrong")
	}

	R := DecodePoint(s[0 : b/8])

	A := DecodePoint(pk)

	S := DecodeInt(s[b/8 : b/4])

	encR := EncodePoint(R)

	var hRR []byte

	hRR = append(hRR, encR...)
	hRR = append(hRR, pk...)
	hRR = append(hRR, m...)

	h := Hint(hRR)

	scMultBS := ScalarMult(B, S)
	scMultAH := ScalarMult(A, h)
	edWards := Edwards(R, scMultAH)

	return CompreArray(scMultBS, edWards)
}
