package main

import (
	"crypto/rand"

	"github.com/frankurcrazy/zuc"
	"github.com/muraliens/cryptoprofile"
)

func (h *Handle) ZucStream(bitlength int) cryptoprofile.BitStream {
	if len(h.key) == 0 {
		h.key = make([]byte, 16)
		rand.Read(h.key)
	}
	if len(h.iv) == 0 {
		h.iv = make([]byte, 16)
		rand.Read(h.iv)
	}

	z := zuc.NewZUC(h.key, h.iv)

	numberStream := h.streamLength / 4

	if numberStream*4 != h.streamLength {
		numberStream++
	}

	streams := z.GenerateKeystream(uint32(numberStream))

	return cryptoprofile.ParseUInt32(bitlength, streams)
}
