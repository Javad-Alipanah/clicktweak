package util

import (
	"errors"
	"math/rand"
	"strings"
)

const DefaultRegex = `^[0-9A-Za-z]{0,7}$`

type Encoder struct {
	radix      uint32
	encodedLen int
	encodeStr  string
}

func NewEncoder62() *Encoder {
	return &Encoder{
		radix:      62,
		encodedLen: 7,
		encodeStr:  "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	}
}

// Encode converts uint32 number to base-62 encoded string
func (e *Encoder) Encode(num uint32) string {
	var buff = make([]byte, e.encodedLen)

	// initialization
	for i := 0; i < e.encodedLen; i++ {
		buff[i] = e.encodeStr[0]
	}

	// encoding
	i := e.encodedLen - 1
	for {
		if num == 0 {
			break
		}
		buff[i] = e.encodeStr[num%e.radix]
		i--
		num /= e.radix
	}

	// remove redundant zeroes
	for i = 0; i < e.encodedLen; i++ {
		if buff[i] != e.encodeStr[0] {
			break
		}
	}
	if i != 0 {
		return string(buff[i:])
	}
	return string(buff[0])
}

// Decode returns base10 representation of base-radix encoded string
func (e *Encoder) Decode(str string) (uint32, error) {
	var num uint32 = 0
	for i := 0; i < len(str); i++ {
		num *= e.radix
		digit := strings.IndexByte(e.encodeStr, str[i])
		if digit == -1 {
			return 0, errors.New("invalid encoded string")
		}
		num += uint32(digit)
	}
	return num, nil
}

func (e *Encoder) SimilarSuggestion(str string) (string, error) {
	l := e.encodedLen - len(str)
	if l > 0 {
		var buff = make([]byte, l)
		for l -= 1; l >= 0; l-- {
			buff[l] = e.encodeStr[rand.Intn(int(e.radix))]
		}
		return str + string(buff), nil
	}

	// add random number to decoded number
	d, err := e.Decode(str)
	if err != nil {
		return str, err
	}
	d += rand.Uint32()
	return e.Encode(d), nil
}
