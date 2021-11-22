package main

import (
	"fmt"
	"os"

	"github.com/muraliens/cryptoprofile"
)

func (h *Handle) StoreStream(fileName string, bs cryptoprofile.BitStream) {
	f, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("Failed to create file")
		return
	}
	defer f.Close()
	if h.binaryFormat {
		stream, err := cryptoprofile.BitStreamToBytes(bs.Value)
		if err != nil {
			fmt.Printf("Failed to convert bit stream")
			return
		}
		f.Write(stream)
	} else {
		f.WriteString(bs.Value)
	}
}
