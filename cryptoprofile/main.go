package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // Blank import needed to import sqlite3
	"github.com/muraliens/cryptoprofile"
)

type Handle struct {
	bitLength    int
	streamLength int
	crypto       string
	bitStream    string
	key          []byte
	iv           []byte
	numSamples   int
}

func main() {
	var keyBits string
	var ivBits string
	var keyBytes string
	var ivBytes string
	var key []byte
	var iv []byte
	h := &Handle{}
	flag.IntVar(&h.bitLength, "bitLength", 4, "Bitlength for the profile")
	flag.IntVar(&h.streamLength, "streamLength", 64, "Stream for the profile")
	flag.StringVar(&h.crypto, "crypto", "", "Crypto algorithm")
	flag.StringVar(&h.bitStream, "bitStream", "", "Bit stream")
	flag.IntVar(&h.numSamples, "numSamples", 1, "Number of samples")
	flag.StringVar(&keyBits, "keyBits", "", "Key bit stream")
	flag.StringVar(&ivBits, "ivBits", "", "IV bit stream")
	flag.StringVar(&keyBytes, "keyBytes", "", "Key byte stream")
	flag.StringVar(&ivBytes, "ivBytes", "", "IV byte stream")

	flag.Parse()

	db, err := gorm.Open("sqlite3", "profile.db")

	if err != nil {
		fmt.Printf("Error : %v\n", err)
		db.LogMode(false)
		db.DB().SetMaxOpenConns(1)
		return
	}
	db.AutoMigrate(&cryptoprofile.DBModel{})

	if h.bitLength < 4 || h.bitLength > 19 {
		fmt.Printf("Invalid bit length, supported between 4 to 19")
		return
	}

	if h.streamLength < int(math.Pow(2, float64(h.bitLength+2))) {
		fmt.Printf("Invalid stream length, minimum required %d", int(math.Pow(2, float64(h.bitLength+2))))
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

	var evps *cryptoprofile.EigenProfiles
	newprofile := false
	evps = cryptoprofile.GetEigenProfilesFromDB(db, h.bitLength)
	if evps == nil {
		evps = cryptoprofile.GetAllEigenProfiles(h.bitLength)
		newprofile = true
	}
	if newprofile {
		cryptoprofile.StoreEigenProfilesToDB(db, h.bitLength, evps)
	}
	for i := 0; i < h.numSamples; i++ {
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
			rs = h.EspressoStream()
		default:
			fmt.Printf("Crypto algorithm not supported")
			return
		}

		evpsr := cryptoprofile.GetEigenProfiles(h.bitLength, rs)

		cryptoprofile.PrintEigenProfiles(fmt.Sprintf("Profile%s_%d_%d.txt", h.crypto, h.bitLength, i+1), h.crypto, h.key, h.iv, rs, evps, evpsr)
		ChiSquareTest(fmt.Sprintf("ChiSquare%s_%d_%d.txt", h.crypto, h.bitLength, i+1), h.crypto, h.key, h.iv, rs, evps, evpsr)
	}
	// evps.PrintFile(fmt.Sprintf("Profile_%d.txt", bitLength))
	// evpsr.PrintFile(fmt.Sprintf("ProfileRand_%d.txt", bitLength))
}
