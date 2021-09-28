package cryptoprofile

import (
	"crypto/rand"
	"fmt"
	"math"
	"os"
	"strconv"
)

const maxString string = "0000000000000000000000000000000000000000000000000000000000000000"

type BitStream struct {
	Value  string
	Length int
}

type KeyStream struct {
	Vaule     string
	Length    int
	BitLength int
}

type EigenValue int
type EigenProfile []EigenValue

type EigenProfileType struct {
	Profile EigenProfile
	Count   int
}

type EigenProfiles struct {
	Profiles []EigenProfileType
}

func GenBitStream(bitlength int, value int) BitStream {
	str := strconv.FormatInt(int64(value), 2)
	if len(str) < bitlength {
		str = maxString[:(bitlength-len(str))] + str
	}
	return BitStream{
		Value:  str,
		Length: len(str),
	}
}

func GenRandBitStream(bitlength int) BitStream {

	numOfBytes := bitlength / 8

	if numOfBytes*8 != bitlength {
		numOfBytes++
	}
	randData := make([]byte, numOfBytes)
	rand.Read(randData)
	var str string
	for _, data := range randData {
		str = str + fmt.Sprintf("%08b", data)
	}

	return BitStream{
		Value:  str[:bitlength],
		Length: bitlength,
	}
}

func ParseBitStream(stream string) (BitStream, error) {
	for i := 0; i < len(stream); i++ {
		if stream[i] != 0x30 && stream[i] != 0x31 {
			return BitStream{}, fmt.Errorf("invalid stream")
		}
	}
	return BitStream{
		Value:  stream,
		Length: len(stream),
	}, nil
}

func isTupleExist(bss []BitStream, tup string) bool {
	for i := range bss {
		if bss[i].Value == tup {
			return true
		}
	}
	return false
}

func (bs BitStream) Vocabulary() []BitStream {
	voc := make([]BitStream, 0)
	for i := 0; i < bs.Length; i++ {
		j := 0
		for {
			sub := bs.Value[j : j+i+1]
			if len(voc) == 0 || !isTupleExist(voc, sub) {
				tup := BitStream{
					Value:  sub,
					Length: len(sub),
				}
				voc = append(voc, tup)
			}
			j++
			if j+i == bs.Length {
				break
			}
		}
	}
	return voc
}

func (bs BitStream) Prefix() BitStream {
	if bs.Length == 1 {
		return bs
	}
	return BitStream{
		Value:  bs.Value[:bs.Length-1],
		Length: bs.Length - 1,
	}
}

func (bs BitStream) AllPrefixes() []BitStream {
	bss := make([]BitStream, 0)

	bss = append(bss, bs)
	pbs := bs
	for i := 0; i < bs.Length-1; i++ {
		pbs = pbs.Prefix()
		bss = append([]BitStream{pbs}, bss...)
	}
	return bss
}

func (bs BitStream) EigenVaule() EigenValue {
	if bs.Length == 1 {
		return 1
	}
	voc := bs.Vocabulary()
	pbs := bs.Prefix()
	pvoc := pbs.Vocabulary()
	ev := EigenValue(len(voc) - len(pvoc))
	return ev
}

func (bs BitStream) EigenProfile() EigenProfile {
	evp := make([]EigenValue, 0)
	bss := bs.AllPrefixes()
	for _, pbs := range bss {
		ev := pbs.EigenVaule()
		evp = append(evp, ev)
	}
	return evp
}

func isEigenProfileMatch(ev1 EigenProfile, ev2 EigenProfile) bool {
	if len(ev1) != len(ev2) {
		return false
	}
	for i := range ev1 {
		if ev1[i] != ev2[i] {
			return false
		}
	}
	return true
}

func (evps *EigenProfiles) AddEigenProfile(evp EigenProfile) {
	if len(evps.Profiles) == 0 {
		evps.Profiles = append(evps.Profiles, EigenProfileType{Profile: evp, Count: 1})
		return
	}
	for i := range evps.Profiles {
		if isEigenProfileMatch(evps.Profiles[i].Profile, evp) {
			evps.Profiles[i].Count++
			return
		}
	}
	evps.Profiles = append(evps.Profiles, EigenProfileType{Profile: evp, Count: 1})
}

func (evps *EigenProfiles) PrintFile(filename string) {
	totalOccurence := 0
	for i := 0; i < len(evps.Profiles); i++ {
		totalOccurence = totalOccurence + evps.Profiles[i].Count
	}
	str := fmt.Sprintf("Total Profiles : %d\n", len(evps.Profiles))
	fmt.Println(str)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file")
		return
	}
	defer f.Close()

	f.WriteString(str)
	f.WriteString("\n\n")
	for i := 0; i < len(evps.Profiles); i++ {

		str := fmt.Sprintf("Profile - %d\n", i+1)
		f.WriteString("--------------------------\n")
		f.WriteString(str)
		f.WriteString("--------------------------\n")
		str = fmt.Sprintf("%v\n", evps.Profiles[i].Profile)
		f.WriteString(str)
		str = fmt.Sprintf("Number of Occurrence : %d\n", evps.Profiles[i].Count)
		f.WriteString(str)
		str = fmt.Sprintf("Ratio of Occurrence : %0.06f\n\n\n", float64(evps.Profiles[i].Count)/float64(totalOccurence))
		f.WriteString(str)
	}
}

func GetAllEigenProfiles(bitLength int) *EigenProfiles {
	evps := &EigenProfiles{
		Profiles: make([]EigenProfileType, 0),
	}

	totalLength := int(math.Pow(2, float64(bitLength)))

	fmt.Printf("")
	for i := 0; i < totalLength; i++ {
		bs := GenBitStream(bitLength, i)
		evp := bs.EigenProfile()

		evps.AddEigenProfile(evp)
		// fmt.Println(evp)
		// fmt.Printf("EVP : %ld\n", len(evp))
		fmt.Printf("\r%d%%", int((i*100)/totalLength))
	}
	fmt.Printf("\r")
	return evps
}

func GetEigenProfiles(bitLength int, rs BitStream) *EigenProfiles {
	evps := &EigenProfiles{
		Profiles: make([]EigenProfileType, 0),
	}

	fmt.Printf("")
	for i := 0; i < rs.Length-(bitLength-1); i++ {
		bs := BitStream{
			Value:  rs.Value[i : i+bitLength],
			Length: bitLength,
		}
		evp := bs.EigenProfile()

		evps.AddEigenProfile(evp)
		// fmt.Println(evp)
		// fmt.Printf("EVP : %ld\n", len(evp))
		fmt.Printf("\r%d%%", int((i*100)/(rs.Length-(bitLength-1))))
	}
	fmt.Printf("\r")
	return evps
}

// func (ks KeyStream) EigenProfile() EigenProfile {

// }
