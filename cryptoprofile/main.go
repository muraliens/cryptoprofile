package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // Blank import needed to import sqlite3
	"github.com/muraliens/cryptoprofile"
)

type Handle struct {
	sbitLength   int
	ebitLength   int
	streamLength int
	crypto       string
	bitStream    string
	key          []byte
	iv           []byte
	numSamples   int
	numRounds    int
	onlyStream   bool
	binaryFormat bool
}

func main() {
	var keyBits string
	var ivBits string
	var keyBytes string
	var ivBytes string
	var key []byte
	var iv []byte
	var bitLengthStr string
	h := &Handle{}
	flag.StringVar(&bitLengthStr, "bitLength", "4", "Bitlength range for the profile")
	flag.IntVar(&h.streamLength, "streamLength", 64, "Stream for the profile")
	flag.StringVar(&h.crypto, "crypto", "", "Crypto algorithm(aes, trivium, zuc, espresso, grain, kasumi")
	flag.StringVar(&h.bitStream, "bitStream", "", "Bit stream")
	flag.IntVar(&h.numSamples, "numSamples", 1, "Number of samples")
	flag.StringVar(&keyBits, "keyBits", "", "Key bit stream")
	flag.StringVar(&ivBits, "ivBits", "", "IV bit stream")
	flag.StringVar(&keyBytes, "keyBytes", "", "Key byte stream")
	flag.StringVar(&ivBytes, "ivBytes", "", "IV byte stream")
	flag.IntVar(&h.numRounds, "numRounds", 0, "Number of Crypto Init Rounds")
	flag.BoolVar(&h.onlyStream, "onlyStream", false, "Generate only stream")
	flag.BoolVar(&h.binaryFormat, "binaryFormat", false, "Stream output format")

	flag.Parse()

	strArray := strings.Split(bitLengthStr, "-")
	if len(strArray) > 2 {
		fmt.Printf("Invalid bit length range\n")
		return
	}

	temp, err := strconv.ParseInt(strArray[0], 10, 32)
	if err != nil {
		fmt.Printf("Invalid bitlength range\n")
		return
	}
	h.sbitLength = int(temp)
	h.ebitLength = int(temp)

	if len(strArray) == 2 {
		temp, err = strconv.ParseInt(strArray[1], 10, 32)
		if err != nil {
			fmt.Printf("Invalid bitlength range\n")
			return
		}
		h.ebitLength = int(temp)
	}

	db, err := gorm.Open("sqlite3", "profile.db")

	if err != nil {
		fmt.Printf("Error : %v\n", err)
		db.LogMode(false)
		db.DB().SetMaxOpenConns(1)
		return
	}
	db.AutoMigrate(&cryptoprofile.DBModel{})

	if h.sbitLength < 4 || h.sbitLength > 19 || h.ebitLength < 4 || h.ebitLength > 19 {
		fmt.Printf("Invalid bit length, supported between 4 to 19")
		return
	}

	if h.streamLength < int(math.Pow(2, float64(h.ebitLength+2))) {
		fmt.Printf("Invalid stream length, minimum required %d", int(math.Pow(2, float64(h.ebitLength+2))))
		return
	}

	if keyBits != "" {
		key, err = cryptoprofile.BitStreamToBytes(keyBits)
		if err != nil {
			fmt.Printf("Invalid key bit stream : %v", err)
			return
		}
	}

	if ivBits != "" {
		iv, err = cryptoprofile.BitStreamToBytes(ivBits)
		if err != nil {
			fmt.Printf("Invalid IV bit stream : %v", err)
			return
		}
	}

	// if len(os.Args) >= 3 {
	// 	rs, err = cryptoprofile.ParseBitStream(os.Args[2])
	// 	if err != nil || rs.Length < h.streamLength {
	// 		fmt.Printf("Invalid bit stream\n")
	// 		return
	// 	}
	// }
	for bitLength := h.sbitLength; bitLength <= h.ebitLength; bitLength++ {
		cryptoStr := h.crypto
		if cryptoStr == "" {
			cryptoStr = "random"
		}
		fmt.Printf("---------------------------------------\n")
		fmt.Printf("BitStream : %s, BitLength : %d\n", cryptoStr, bitLength)
		fmt.Printf("---------------------------------------\n")
		var evps *cryptoprofile.EigenProfiles
		newprofile := false
		evps = cryptoprofile.GetEigenProfilesFromDB(db, bitLength)
		if evps == nil {
			evps = cryptoprofile.GetAllEigenProfiles(bitLength)
			newprofile = true
		}
		if newprofile {
			cryptoprofile.StoreEigenProfilesToDB(db, bitLength, evps)
		}
		for i := 0; i < h.numSamples; i++ {
			fmt.Printf("Sample : %d\n", i+1)
			var rs cryptoprofile.BitStream
			switch h.crypto {
			case "":
				rs = cryptoprofile.GenRandBitStream(h.streamLength)
			case "zuc":
				h.iv = iv
				h.key = key
				rs = h.ZucStream()
			case "espresso":
				h.iv = iv
				h.key = key
				rs = h.EspressoStream(h.numRounds)
			case "trivium":
				h.iv = iv
				h.key = key
				rs = h.TriviumStream()
			case "aes":
				h.iv = iv
				h.key = key
				rs = h.AESStream()
			case "kasumi":
				h.iv = iv
				h.key = key
				rs = h.KasumiStream()
			case "grain":
				h.iv = iv
				h.key = key
				rs = h.GrainStream()
			case "sm4":
				h.iv = iv
				h.key = key
				rs = h.SM4Stream()
			default:
				fmt.Printf("Crypto algorithm not supported")
				return
			}
			if h.onlyStream {
				fileName := fmt.Sprintf("Stream%s.txt", h.crypto)
				if i == 0 {
					os.Remove(fileName)
				}
				h.StoreStream(fileName, rs)
			} else {
				evpsr := cryptoprofile.GetEigenProfiles(bitLength, rs)
				cryptoprofile.PrintEigenProfiles(fmt.Sprintf("Profile%s_%d_%d.txt", h.crypto, bitLength, i+1), h.crypto, h.key, h.iv, rs, evps, evpsr)
				pvalue := ChiSquareTest(fmt.Sprintf("ChiSquare%s_%d_%d.txt", h.crypto, bitLength, i+1), h.crypto, h.key, h.iv, rs, evps, evpsr)
				missingProfile := cryptoprofile.UpdateEigenProfiles(h.crypto, bitLength, h.key, h.iv, rs, evps, evpsr)
				cryptoprofile.UpdatePValue(h.crypto, bitLength, pvalue, missingProfile)
			}
		}
		fmt.Printf("---------------------------------------\n\n")
	}
	// evps.PrintFile(fmt.Sprintf("Profile_%d.txt", bitLength))
	// evpsr.PrintFile(fmt.Sprintf("ProfileRand_%d.txt", bitLength))
}
