package main

import (
	"crypto/rand"

	"github.com/bmkessler/trivium"
	"github.com/muraliens/cryptoprofile"
)

func (h *Handle) TriviumStream() cryptoprofile.BitStream {
	if len(h.key) == 0 {
		h.key = make([]byte, 16)
		rand.Read(h.key)
	}
	if len(h.iv) == 0 {
		h.iv = make([]byte, 16)
		rand.Read(h.iv)
	}

	var key [10]byte
	var iv [10]byte

	copy(key[:], h.key)
	copy(iv[:], h.iv)

	t := trivium.NewTrivium(key, iv)

	numberStream := h.streamLength / 8

	if numberStream*8 != h.streamLength {
		numberStream++
	}

	stream := make([]byte, 0)

	for i := 0; i < numberStream; i++ {
		data := t.NextByte()
		stream = append(stream, data)
	}
	return cryptoprofile.ParseBytes(h.streamLength, stream)
}
