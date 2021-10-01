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
}

func main() {
	h := &Handle{}
	flag.IntVar(&h.bitLength, "bitLength", 4, "Bitlength for the profile")
	flag.IntVar(&h.streamLength, "streamLength", 64, "Stream for the profile")
	flag.StringVar(&h.crypto, "crypto", "", "Crypto algorithm")
	flag.StringVar(&h.bitStream, "bitStream", "", "Bit stream")

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
	var rs cryptoprofile.BitStream
	switch h.crypto {
	case "":
		rs = cryptoprofile.GenRandBitStream(h.streamLength)
	case "zuc":
		rs = h.ZucStream(h.streamLength)
	default:
		fmt.Printf("Crypto algorithm not supported")
		return
	}

	fmt.Printf("Stream length : %d\n", h.streamLength)

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
	evpsr := cryptoprofile.GetEigenProfiles(h.bitLength, rs)

	if newprofile {
		cryptoprofile.StoreEigenProfilesToDB(db, h.bitLength, evps)
	}

	cryptoprofile.PrintEigenProfiles(fmt.Sprintf("ProfileRand_%d.txt", h.bitLength), h.crypto, h.key, h.iv, evps, evpsr)
	// evps.PrintFile(fmt.Sprintf("Profile_%d.txt", bitLength))
	// evpsr.PrintFile(fmt.Sprintf("ProfileRand_%d.txt", bitLength))
}
