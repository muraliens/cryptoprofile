package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/muraliens/cryptoprofile"
)

func main() {

	var inFileName string
	var seqLength int
	var numSamples int
	var outFileName string

	flag.StringVar(&inFileName, "inFileName", "rand.txt", "Input File Name")
	flag.IntVar(&seqLength, "seqLength", 256, "Sequence Length")
	flag.IntVar(&numSamples, "numSamples", 10, "Number of Samples")
	flag.StringVar(&outFileName, "outFileName", "nlcpprofile.xlsx", "Output Filename")

	flag.Parse()

	f, err := os.Open(inFileName)
	if err != nil {
		fmt.Printf("Invalid File")
		return
	}
	data := make([]byte, numSamples*seqLength)
	n, err := f.Read(data)
	if err != nil || n != len(data) {
		fmt.Printf("Error in reading file")
		return
	}
	bs := cryptoprofile.BitStream{
		Value:  string(data),
		Length: len(data),
	}
	startTime := time.Now()
	nlcp := bs.CreateNLCProfile(seqLength, numSamples)
	mv := cryptoprofile.CalculateMeanVarience(nlcp)
	emv := cryptoprofile.ExpectedMeanVarience(seqLength)
	cryptoprofile.SaveNLCPProfile(outFileName, nlcp, mv, emv)
	elasped := time.Since(startTime)
	fmt.Printf("Runtime : %dh: %dm: %ds\n", int(elasped.Hours()), int(elasped.Minutes()), int(elasped.Seconds()))
}
