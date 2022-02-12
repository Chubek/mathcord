package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
	numParsed, err := strconv.ParseFloat(strings.Trim(num, ""), 32)
	if err != nil {
		log.Fatal(err)
	}

	return numParsed
}

func Min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func DecimalToHex(dec uint64) string {
	hexVal := fmt.Sprintf("%016x", dec)
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

func BinaryToDecimal(bin string) uint64 {
	num, err := strconv.ParseUint(bin, 2, 64)
	if err != nil {
		log.Fatal(err)
	}
	return uint64(num)
}

func BinaryToDecimal32(bin string) uint32 {
	num, err := strconv.ParseUint(bin, 2, 32)
	if err != nil {
		log.Fatal(err)
	}
	return uint32(num)
}

func RotateRightByNBits(b []byte, n int) []byte {
	for i := range b {
		b[i] = (b[i] >> n) | (b[i] << (64 - n))
	}

	return b
}

func RotateUintRightByNBits(b uint64, n int) uint64 {
	return (b >> n) | (b << (64 - n))
}

func RotateStringRightByNBits(s string, n int) string {
	integerForm := BinaryToDecimal(s)

	rotated := (integerForm >> n) | (integerForm << (64 - n))

	return IntegerToBinary(rotated, 64)
}

func ShiftStringRightByNBits(s string, n int) string {
	integerForm := BinaryToDecimal(s)

	shifted := integerForm >> n

	return IntegerToBinary(shifted, 64)
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

func ShiftUintRightByNBits(b uint64, n int) uint64 {
	return b >> n
}

func IntegerToBinary(num uint64, len int) string {
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

func Integer32ToBinary(num uint32, len int) string {
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

func XorThree(one, two, three string) string {
	oneInt := BinaryToDecimal(one)
	twoInt := BinaryToDecimal(one)
	threeInt := BinaryToDecimal(one)

	res := oneInt ^ twoInt ^ threeInt

	return IntegerToBinary(res, 64)
}

func AddFour(A, B, C, D string) string {
	AInt := BinaryToDecimal(A)
	BInt := BinaryToDecimal(B)
	CInt := BinaryToDecimal(C)
	DInt := BinaryToDecimal(D)

	sum := AInt + BInt + CInt + DInt

	return IntegerToBinary(sum, 64)
}

func ConvertBinaryToIntegerArray(arr []string) []uint64 {
	var ret []uint64

	for _, str := range arr {
		ret = append(ret, BinaryToDecimal(str))
	}

	return ret
}
