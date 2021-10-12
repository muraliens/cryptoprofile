package cryptoprofile

import (
	"fmt"
	"math"
	"testing"
)

func TestVocabulary(t *testing.T) {

	evps := &EigenProfiles{
		Profiles: make([]EigenProfileType, 0),
	}

	bitLength := 19
	totalValue := int(math.Pow(2, float64(bitLength)))
	fmt.Printf("")
	for i := 0; i < totalValue; i++ {
		bs := GenBitStream(bitLength, i)
		evp := bs.EigenProfile()

		evps.AddEigenProfile(evp)
		// fmt.Println(evp)
		// fmt.Printf("EVP : %d\n", len(evp))
		fmt.Printf("\r%d", int((i*100)/totalValue))
	}

	fmt.Printf("Total Profiles : %d\n", len(evps.Profiles))

}

func TestConvertion(t *testing.T) {
	res, err := BitStreamToBytes("011010101001101010100110101010")
	if err != nil {
		t.Fatalf("Failed to convert %v", err)
	}
	fmt.Printf("%x", res)
}
