package tron

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/btcsuite/btcd/btcutil/base58"
)

// Hex2TAddr converts a (unknwon) Hex to TAddr
func Hex2TAddr(hexStr string) string {
	if hexStr[0] == 'T' {
		log.Printf("Taddr %s input as a hex?", hexStr)
		return hexStr
	}

	if len(hexStr) == 20*2 {
		hexStr = "41" + hexStr
	}

	if len(hexStr) == 21*2 {
		addrBytes, err := hex.DecodeString(hexStr)
		chk(err)
		sum0 := sha256.Sum256(addrBytes)
		sum1 := sha256.Sum256(sum0[:])
		chksum := sum1[0:4]
		addrBytes = append(addrBytes, chksum...)

		return base58.Encode(addrBytes)
	}

	panic(hexStr + "is not a no-prefix hex addr")
}

// Tstring2HexAddr converts a T-string to Hex addr
func Tstring2HexAddr(theTstr string) string {
	if theTstr[0] != 'T' {
		panic(theTstr + " is not a TAddr")
	}

	bs58decoded := base58.Decode(theTstr)
	withoutSum := bs58decoded[:21]
	return hex.EncodeToString(withoutSum)
}
