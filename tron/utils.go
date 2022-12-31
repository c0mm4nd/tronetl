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

	// if genesis block miner
	if len(hexStr) == 230 && hexStr == "206e65772073797374656d206d75737420616c6c6f77206578697374696e672073797374656d7320746f206265206c696e6b656420746f67657468657220776974686f757420726571756972696e6720616e792063656e7472616c20636f6e74726f6c206f7220636f6f7264696e6174696f6e" {
		return "GenesisMiner" //  as a placeholder
	}

	// if genesis tx
	if len(hexStr) == 46 && hexStr == "3078303030303030303030303030303030303030303030" {
		return "7YxAaK71utTpYJ8u4Zna7muWxd1pQwimpGxy8" // = 3078303030303030303030303030303030303030303030
	}

	if len(hexStr) == 20*2 {
		hexStr = "41" + hexStr // make sure the T-prefix
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
