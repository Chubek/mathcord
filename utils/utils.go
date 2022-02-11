package utils

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
)

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func ParseToFloat(num string) float64 {
	numParsed, err := strconv.ParseFloat(num, 32)
	if err != nil {
		log.Fatal(err)
	}

	return numParsed
}

func Min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func Max(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func DecimalToHex(dec uint) string {
	hexVal := fmt.Sprintf("%x", dec)

	return hexVal
}

func ByteToHex(dec byte) string {
	hexVal := fmt.Sprintf("%x", dec)

	return hexVal
}

func StrToBinary(s string) string {
	res := ""
	for _, c := range s {
		res = fmt.Sprintf("%s%.8b", res, c)
	}
	return res
}

func BinaryToDecimal(bin string) uint {
	num, err := strconv.ParseInt(bin, 2, 64)
	if err != nil {
		log.Fatal(err)
	}
	return uint(num)
}

func RotateRightByNBits(b []byte, n int) []byte {
	for i := range b {
		b[i] = (b[i] >> n) | (b[i] << (64 - n))
	}

	return b
}

func RotateUintRightByNBits(b uint, n int) uint {
	return (b >> n) | (b << (64 - n))
}

func RotateStringRightByNBits(s string, n int) string {
	integerForm := BinaryToDecimal(s)

	rotated := (integerForm >> n) | (integerForm << (64 - n))

	return IntegerToBinary(uint(rotated), len(s))
}

func ShiftStringRightByNBits(s string, n int) string {
	integerForm := BinaryToDecimal(s)

	shifted := integerForm >> n

	return IntegerToBinary(uint(shifted), len(s))
}

func RotateByteRightByNBits(b byte, n int) byte {
	return (b >> n) | (b << (64 - n))
}

func ShiftRightByNBits(b []byte, n int) []byte {
	for i := range b {
		b[i] = (b[i] >> n)
	}

	return b
}

func ShiftUintRightByNBits(b uint, n int) uint {
	return b >> n
}

func IntegerToBinary(num uint, len int) string {
	switch len {
	case 8:
		return fmt.Sprintf("%08b", num)
	case 16:
		return fmt.Sprintf("%016b", num)
	case 32:
		return fmt.Sprintf("%032b", num)
	case 64:
		return fmt.Sprintf("%064b", num)
	}
	return fmt.Sprintf("%0128b", num)

}

func IntegerTo128Bytes(num int) []byte {
	bs := make([]byte, 128)
	binary.LittleEndian.PutUint64(bs, uint64(num))

	return bs
}


func XorThree(one, two, three string) string {
	oneInt := BinaryToDecimal(one)
	twoInt := BinaryToDecimal(one)
	threeInt := BinaryToDecimal(one)

	res := oneInt ^ twoInt ^ threeInt

	return IntegerToBinary(res, len(one))
}

func AddFour(A, B, C, D string) string {
	AInt := BinaryToDecimal(A)
	BInt := BinaryToDecimal(B)
	CInt := BinaryToDecimal(C)
	DInt := BinaryToDecimal(D)

	sum := AInt + BInt + CInt + DInt

	return IntegerToBinary(sum, len(A))
}