// this source code is to test the matching algorithm locally
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"

	"github.com/clearmatics/bn256"
)

type IndexStru struct {
	label      uint32
	pubKey     [16]*bn256.G1
	tag        [16][3][3]uint8
	ciphertext [16][3]*bn256.G1
}

type QueryStru struct {
	pubKey     *bn256.G2
	tag        [16]uint8
	ciphertext [16]*bn256.G2
}

var (
	index [5]IndexStru
	g1s   *bn256.G1
	g2s   *bn256.G2
	query QueryStru
)

// blockStr includes the operator > or <
func getHashedValue(blockStr string, prefix int64, blockId int) *big.Int {
	// the first block, no prefix
	if blockId == 0 {
		byteToHash := []byte(blockStr) // only hash the block string
		hashed := sha256.Sum256(byteToHash[:])
		hashedValue := new(big.Int).SetBytes(hashed[:])
		return hashedValue
	} else { // includes the prefix
		byteToHash := []byte(strconv.FormatInt(prefix, 10)) // hash the prefix first
		hashed := sha256.Sum256(byteToHash[:])
		blockByte := []byte(blockStr)

		var buffer bytes.Buffer // combine the prefix hash value and the block string
		buffer.Write(hashed[:])
		buffer.Write(blockByte)
		finalByte := buffer.Bytes()

		finalHashed := sha256.Sum256(finalByte[:]) // hash the combined value
		hashedValue := new(big.Int).SetBytes(finalHashed[:])
		return hashedValue
	}
}

func indexBlockEnc(block int64, prefix int64, id int, blockId int) {
	var i int64
	var tagPos = [3]int{0, 0, 0}
	var cipherPos = 0
	gamma, _ := rand.Int(rand.Reader, bn256.Order)
	index[id].pubKey[blockId] = new(bn256.G1).ScalarMult(g1s, gamma)
	for i = 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			index[id].tag[blockId][i][j] = 100
		}
	}
	for i = 0; i < 4; i++ {
		if i == block {
			continue
		} else if i < block {
			blockStr := strconv.FormatInt(i, 10) + "<"
			exp := getHashedValue(blockStr, prefix, blockId)
			// calculate the tag
			tag, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(3)).String()) // the tag
			index[id].tag[blockId][tag][tagPos[tag]] = uint8(cipherPos)
			tagPos[tag]++

			// generate the ciphertext
			index[id].ciphertext[blockId][cipherPos] = new(bn256.G1).ScalarBaseMult(exp)
			index[id].ciphertext[blockId][cipherPos].ScalarMult(index[id].ciphertext[blockId][cipherPos], gamma)
			cipherPos++
		} else if i > block {
			blockStr := strconv.FormatInt(i, 10) + ">"
			exp := getHashedValue(blockStr, prefix, blockId)
			// calculate the tag
			tag, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(3)).String()) // the tag
			index[id].tag[blockId][tag][tagPos[tag]] = uint8(cipherPos)
			tagPos[tag]++

			// generate the ciphertext
			index[id].ciphertext[blockId][cipherPos] = new(bn256.G1).ScalarBaseMult(exp)
			index[id].ciphertext[blockId][cipherPos].ScalarMult(index[id].ciphertext[blockId][cipherPos], gamma)
			cipherPos++
		}

	}
}

func indexEnc(v uint32, id int) {
	vStr := strconv.FormatInt(int64(v), 2)
	vStr = fmt.Sprintf("%032s", vStr)
	var prefix int64
	for i := 0; i < 32/2; i++ {
		block, _ := strconv.ParseInt(vStr[i*2:i*2+2], 2, 0)
		if i == 0 {
			prefix = -1
		} else {
			prefix, _ = strconv.ParseInt(vStr[0:i*2], 2, 0)
		}
		indexBlockEnc(block, prefix, id, i)
	}
	index[id].label = v
}

func queryBlockEnc(block int64, prefix int64, blockId int, gamma *big.Int) {
	blockStr := strconv.FormatInt(block, 10) + ">"
	exp := getHashedValue(blockStr, prefix, blockId)
	tagValue, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(3)).String())
	query.tag[blockId] = uint8(tagValue)
	query.ciphertext[blockId] = new(bn256.G2).ScalarBaseMult(exp)
	query.ciphertext[blockId].ScalarMult(query.ciphertext[blockId], gamma)
}

func queryEnc(v uint32) {
	gamma, _ := rand.Int(rand.Reader, bn256.Order)
	query.pubKey = new(bn256.G2).ScalarMult(g2s, gamma)
	// query.pubKey.Neg(query.pubKey)						// in on-chain test, one variable should be negative because the pre-compiled contract 0x08 test the sum of both exponents
	vStr := strconv.FormatInt(int64(v), 2)
	vStr = fmt.Sprintf("%032s", vStr)
	var prefix int64
	for i := 0; i < 32/2; i++ {
		block, _ := strconv.ParseInt(vStr[i*2:i*2+2], 2, 0)
		if i == 0 {
			prefix = -1
		} else {
			prefix, _ = strconv.ParseInt(vStr[0:i*2], 2, 0)
		}
		queryBlockEnc(block, prefix, i, gamma)
	}
}

func match() {
	for i := 0; i < 5; i++ {
		var isMatched bool = false
		for j := 0; j < 16; j++ {
			for k := 0; k < 3; k++ {
				if index[i].tag[j][query.tag[j]][k] == 100 {
					break
				}
				id := index[i].tag[j][query.tag[j]][k]

				k1 := bn256.Pair(index[i].ciphertext[j][id], query.pubKey)
				k2 := bn256.Pair(index[i].pubKey[j], query.ciphertext[j])
				k1Byte := k1.Marshal()
				k2Byte := k2.Marshal()

				if bytes.Equal(k1Byte, k2Byte) {
					isMatched = true
					break
				}
			}
			if isMatched == true {
				break
			}
		}
		if isMatched == true {
			fmt.Println(index[i].label)
		}
	}
}

func main() {
	s, _ := rand.Int(rand.Reader, bn256.Order)
	g1s = new(bn256.G1).ScalarBaseMult(s)
	g2s = new(bn256.G2).ScalarBaseMult(s)

	indexEnc(500, 0)
	indexEnc(85, 1)
	indexEnc(3, 2)
	indexEnc(150, 3)
	indexEnc(329, 4)

	queryEnc(500)

	match()
}
