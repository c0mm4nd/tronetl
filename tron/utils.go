package tron

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/btcsuite/btcd/btcutil/base58"
)

// EnsureTAddr converts a (unknwon) Hex to TAddr
func EnsureTAddr(hexStr string) string {
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

// EnsureHexAddr converts a T-string to Hex addr
func EnsureHexAddr(theTstr string) string {
	if theTstr[0] != 'T' {
		log.Println(theTstr + " is not a TAddr")
		return theTstr
	}

	bs58decoded := base58.Decode(theTstr)
	withoutSum := bs58decoded[:21]
	return hex.EncodeToString(withoutSum)
}
