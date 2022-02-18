package sha512

import (
	"mathcord/dtype"
	"mathcord/utils"
)

type Hash struct {
	Message   *dtype.Sha512Message
	Buffer    *dtype.Sha512Buffer
	HexDigest string
	Digest    []byte
}

func NewSha512() *Hash {
	return &Hash{}
}

func (sha512 *Hash) Update(data []byte) {
	sha512.Message = dtype.NewMessage(data)
	sha512.Buffer = dtype.NewBuffer()
}

func (sha512 *Hash) Calculate() {
	for _, chunk := range *sha512.Message.Chunks {
		sha512.Buffer.ProcessBlock(&chunk)
	}

	sha512.HexDigest = sha512.Buffer.ToHexaDecimal()
	sha512.Digest = utils.HexToAscii(sha512.HexDigest)
}

func (sha512 *Hash) GetHexDigest() string {
	return sha512.HexDigest
}
func (sha512 *Hash) GetDigest() []byte {
	return sha512.Digest
}

func HashWithSha512(data []byte) []byte {
	var byteNew = make([]byte, len(data), len(data))

	for i, b := range data {
		byteNew[i] = b
	}

	h := NewSha512()

	h.Update(byteNew)
	h.Calculate()
	return h.GetDigest()
}
