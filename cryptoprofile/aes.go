package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/muraliens/cryptoprofile"
)

func (h *Handle) AESStream() cryptoprofile.BitStream {
	if len(h.key) == 0 {
		h.key = make([]byte, 32)
		rand.Read(h.key)
	}
	if len(h.iv) == 0 {
		h.iv = make([]byte, 16)
		rand.Read(h.iv)
	}

	block, err := aes.NewCipher(h.key)
	if err != nil {
		return cryptoprofile.BitStream{}
	}
	numberStream := h.streamLength / 8

	if numberStream*8 != h.streamLength {
		numberStream++
	}

	plainData := make([]byte, numberStream)
	cipherData := make([]byte, numberStream)

	stream := cipher.NewCTR(block, h.iv)
	stream.XORKeyStream(cipherData, plainData)

	return cryptoprofile.ParseBytes(h.streamLength, cipherData)
}
