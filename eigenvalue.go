package cryptoprofile

import (
	"crypto/rand"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
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

type EigenValue byte
type EigenProfile []byte

type EigenProfileType struct {
	Profile EigenProfile
	Count   int
}

type ProfileModel struct {
	gorm.Model
	EigenProfile
}

type DBModel struct {
	gorm.Model
	BitLength int    `gorm:"column:BitLength"`
	Profile   []byte `gorm:"column:Profile;size:500"`
	Count     int    `gorm:"column:Count"`
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

func ParseUInt32(bitlength int, streams []uint32) BitStream {

	if len(streams)*32 < bitlength {
		return BitStream{}
	}

	var str string
	for _, stream := range streams {
		str = str + fmt.Sprintf("%032b", stream)
	}

	return BitStream{
		Value:  str[:bitlength],
		Length: bitlength,
	}
}

func ParseBytes(bitlength int, streams []byte) BitStream {

	if len(streams)*8 < bitlength {
		return BitStream{}
	}

	var str string
	for _, stream := range streams {
		str = str + fmt.Sprintf("%08b", stream)
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

func BitStreamToBytes(stream string) ([]byte, error) {
	result := make([]byte, 0)
	str := stream
	for {
		l := 0
		if len(str) > 8 {
			l = len(str) - 8
		}
		temp, err := strconv.ParseInt(str[l:], 2, 64)
		if err != nil {
			return nil, err
		}
		result = append([]byte{byte(temp)}, result...)
		if l == 0 {
			break
		} else {
			str = str[:l]
		}
	}
	return result, nil
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

	var evp EigenProfile
	//	evp := make([]EigenValue, 0)
	bss := bs.AllPrefixes()
	for _, pbs := range bss {
		ev := pbs.EigenVaule()
		evp = append(evp, byte(ev))
	}
	return EigenProfile(evp)
}

func IsEigenProfileMatch(ev1 EigenProfile, ev2 EigenProfile) bool {
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
		if IsEigenProfileMatch(evps.Profiles[i].Profile, evp) {
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

func GetEigenProfilesFromDB(db *gorm.DB, bitLength int) *EigenProfiles {

	var model []DBModel
	err := db.Where("BitLength=?", bitLength).Find(&model).Error
	if err != nil {
		return nil
	}
	if len(model) == 0 {
		return nil
	}
	evps := &EigenProfiles{
		Profiles: make([]EigenProfileType, 0),
	}
	for i := range model {
		// var profile EigenProfile
		// profile = model[i].Profile
		// // for _, ev := range model[i].Profile {
		// // 	profile = append(profile, ev)
		// }
		evps.Profiles = append(evps.Profiles, EigenProfileType{Profile: model[i].Profile, Count: model[i].Count})
	}
	return evps
}

func StoreEigenProfilesToDB(db *gorm.DB, bitLength int, evps *EigenProfiles) error {

	if evps == nil {
		return fmt.Errorf("invalid profiles")
	}

	for _, profile := range evps.Profiles {
		// var temp []int
		// for _, ev := range profile.Profile {
		// 	temp = append(temp, int(ev))
		// }
		model := DBModel{
			BitLength: bitLength,
			Profile:   profile.Profile,
			Count:     profile.Count,
		}
		err := db.Create(&model).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func PrintEigenProfiles(filename string, crypto string, key []byte, iv []byte, rs BitStream, evps *EigenProfiles, evpsr *EigenProfiles) {
	total := 0
	for i := 0; i < len(evps.Profiles); i++ {
		total = total + evps.Profiles[i].Count
	}
	totalr := 0
	for i := 0; i < len(evpsr.Profiles); i++ {
		totalr = totalr + evpsr.Profiles[i].Count
	}

	fmt.Printf("Orginal Total Profiles : %d\n", len(evps.Profiles))
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file")
		return
	}
	defer f.Close()
	str := fmt.Sprintf("Orginal Total Profiles     : %d\n", len(evps.Profiles))
	f.WriteString(str)
	if crypto == "" {
		str = fmt.Sprintf("Rand Total Profiles         : %d\n", len(evpsr.Profiles))
		f.WriteString(str)
	} else {
		str = fmt.Sprintf("%s Total Profiles           : %d\n", crypto, len(evpsr.Profiles))
		f.WriteString(str)
		if len(key) > 0 {
			temp := ParseBytes(len(key)*8, key)
			str = fmt.Sprintf("Key : %s\n", temp.Value)
			f.WriteString(str)
		}
		if len(iv) > 0 {
			temp := ParseBytes(len(iv)*8, iv)
			str = fmt.Sprintf("IV  : %s\n", temp.Value)
			f.WriteString(str)
		}
		// fmt.Printf("RS : %v\n", rs)
		// if rs.Length > 0 {
		// 	str = fmt.Sprintf("Keystream  : %s\n", rs.Value)
		// 	f.WriteString(str)
		// }
	}

	f.WriteString("\n\n")
	for i := 0; i < len(evps.Profiles); i++ {

		str := fmt.Sprintf("Profile - %d\n", i+1)
		f.WriteString("--------------------------\n")
		f.WriteString(str)
		f.WriteString("--------------------------\n")
		str = fmt.Sprintf("%v\n", evps.Profiles[i].Profile)
		f.WriteString(str)
		str = fmt.Sprintf("Number of Occurrence (Original) : %d\n", evps.Profiles[i].Count)
		f.WriteString(str)
		str = fmt.Sprintf("Ratio of Occurrence (Original)  : %0.06f\n", float64(evps.Profiles[i].Count)/float64(total))
		f.WriteString(str)
		found := false
		index := 0
		for j := range evpsr.Profiles {
			if IsEigenProfileMatch(evpsr.Profiles[j].Profile, evps.Profiles[i].Profile) {
				found = true
				index = j
				break
			}
		}

		if found {
			str = fmt.Sprintf("Number of Occurrence (Rand)     : %d\n", evpsr.Profiles[index].Count)
			f.WriteString(str)
			str = fmt.Sprintf("Ratio of Occurrence (Rand)      : %0.06f\n\n\n", float64(evpsr.Profiles[index].Count)/float64(totalr))
			f.WriteString(str)
		} else {
			str = fmt.Sprintf("Number of Occurrence (Rand)     : Missing\n")
			f.WriteString(str)
			str = fmt.Sprintf("Ratio of Occurrence (Rand)      : Missing\n\n")
			f.WriteString(str)
		}
	}
}

// func (ks KeyStream) EigenProfile() EigenProfile {

// }
