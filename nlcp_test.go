package cryptoprofile

import "testing"

func TestNLCP(t *testing.T) {
	v := "0011011100101110110"
	bs := BitStream{
		Value:  v,
		Length: len(v),
	}
	bs.CreateNLCProfile(19, 1)
}
