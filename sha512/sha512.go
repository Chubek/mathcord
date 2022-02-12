package sha512

import (
	"mathcord/dtype"
)

type Sha512Hash struct {
	Message   *dtype.Sha512Message
	Buffer    *dtype.Sha512Buffer
	HexDigest string
}

func NewSha512() *Sha512Hash {
	return &Sha512Hash{}
}

func (sha512 *Sha512Hash) Update(str string) {
	sha512.Message = dtype.NewMessage(str)
	sha512.Buffer = dtype.NewBuffer()
}

func (sha512 *Sha512Hash) Calculate() {
	for _, chunk := range *sha512.Message.Chunks {
			sha512.Buffer.ProcessBlock(&chunk)
	}

	sha512.HexDigest = sha512.Buffer.ToHexaDecimal()
}

func (sha512 *Sha512Hash) GetHexDigest() string {
	return sha512.HexDigest
}
