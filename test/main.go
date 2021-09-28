package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/muraliens/cryptoprofile"
)

func main() {
	bitLength := 4
	var err error
	if len(os.Args) >= 2 {
		bitLength, err = strconv.Atoi(os.Args[1])
		if err != nil || bitLength > 19 || bitLength < 4 {
			fmt.Printf("Invalid bit length\n")
			return
		}
	}
	streamLength := int(math.Pow(2, float64(bitLength+2)))
	rs := cryptoprofile.GenRandBitStream(streamLength)

	fmt.Printf("Stream length : %d\n", streamLength)

	if len(os.Args) >= 3 {
		rs, err = cryptoprofile.ParseBitStream(os.Args[2])
		if err != nil || rs.Length < streamLength {
			fmt.Printf("Invalid bit stream\n")
			return
		}
	}

	evps := cryptoprofile.GetAllEigenProfiles(bitLength)
	evpsr := cryptoprofile.GetEigenProfiles(bitLength, rs)

	evps.PrintFile(fmt.Sprintf("Profile_%d.txt", bitLength))
	evpsr.PrintFile(fmt.Sprintf("ProfileRand_%d.txt", bitLength))
}
