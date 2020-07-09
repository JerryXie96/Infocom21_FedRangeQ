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

// the structure of one index item
type IndexStru struct {
	label      uint32           // the label of one index
	pubKey     [16]*bn256.G1    // the public key
	tag        [16][3][3]uint8  // the tag of each ciphertext [the number of blocks][the number of tag types][the max conflicts in one tag type]
	ciphertext [16][3]*bn256.G1 // the ciphertext [the number of blocks][the number of variables]
}

// the structure of one query
type QueryStru struct {
	pubKey     *bn256.G2     // the public key of one query
	tag        [16]uint8     // the tag of each block ciphertext
	ciphertext [16]*bn256.G2 // the ciphertext of each block
}

var (
	index [5]IndexStru // the index list
	g1s   *bn256.G1    // the key of index
	g2s   *bn256.G2    // the key of query
	query QueryStru    // the query
)

// getHashedValue: compute the hash value in tag and ciphertext (iStr: the value to be hashed, prefix: the prefix of this block, blockId: the current block number)
func getHashedValue(iStr string, prefix int64, blockId int) *big.Int {
	// the first block, no prefix
	if blockId == 0 {
		byteToHash := []byte(iStr) // only hash the block string
		hashed := sha256.Sum256(byteToHash[:])
		hashedValue := new(big.Int).SetBytes(hashed[:])
		return hashedValue
	} else { // includes the prefix
		byteToHash := []byte(strconv.FormatInt(prefix, 10)) // hash the prefix first
		hashed := sha256.Sum256(byteToHash[:])
		blockByte := []byte(iStr)

		var buffer bytes.Buffer // combine the prefix hash value and the block string
		buffer.Write(hashed[:])
		buffer.Write(blockByte)
		finalByte := buffer.Bytes()

		finalHashed := sha256.Sum256(finalByte[:]) // hash the combined value
		hashedValue := new(big.Int).SetBytes(finalHashed[:])
		return hashedValue
	}
}

// indexBlockEnc: encrypt one block index (block: the block value, prefix: the prefix value, id: the current index number, blockId: the current block number)
func indexBlockEnc(block int64, prefix int64, id int, blockId int) {
	var i int64
	var tagPos = [3]int{0, 0, 0}                                     // the current available space of each tag's list
	var cipherPos = 0                                                // the current available space of ciphertext list
	gamma, _ := rand.Int(rand.Reader, bn256.Order)                   // the nonce
	index[id].pubKey[blockId] = new(bn256.G1).ScalarMult(g1s, gamma) // calculate the public key of one block in one index item
	for i = 0; i < 3; i++ {                                          // initialize the tag list with 100 (one value out of range)
		for j := 0; j < 3; j++ {
			index[id].tag[blockId][i][j] = 100
		}
	}
	// iStr includes the operator > or <
	for i = 0; i < 4; i++ {
		if i == block { // do not encrypt the equal block
			continue
		} else if i < block { // the current variable is smaller than the current block
			iStr := strconv.FormatInt(i, 10) + "<"       // add the inequality operator into the string to be hashed
			exp := getHashedValue(iStr, prefix, blockId) // calculate the hash value in tag and ciphertext
			// calculate the tag
			tag, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(3)).String()) // the tag
			index[id].tag[blockId][tag][tagPos[tag]] = uint8(cipherPos)           // store the list number of the ciphertext in the current available space of corresponding tag
			tagPos[tag]++

			// generate the ciphertext
			index[id].ciphertext[blockId][cipherPos] = new(bn256.G1).ScalarBaseMult(exp)
			index[id].ciphertext[blockId][cipherPos].ScalarMult(index[id].ciphertext[blockId][cipherPos], gamma)
			cipherPos++
		} else if i > block { // the current variable is larger than the current block (the process procedure is similar)
			iStr := strconv.FormatInt(i, 10) + ">"
			exp := getHashedValue(iStr, prefix, blockId)
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

// indexEnc: encrypt one value as an index item (v: the value to be encrypted, id: the current index number)
func indexEnc(v uint32, id int) {
	vStr := strconv.FormatInt(int64(v), 2) // calculate the binary value
	vStr = fmt.Sprintf("%032s", vStr)      // pad to 32 bits
	var prefix int64
	for i := 0; i < 32/2; i++ {
		block, _ := strconv.ParseInt(vStr[i*2:i*2+2], 2, 0) // the block contains 2 bits
		if i == 0 {                                         // the first block (no prefix)
			prefix = -1
		} else { // other (has prefix)
			prefix, _ = strconv.ParseInt(vStr[0:i*2], 2, 0)
		}
		indexBlockEnc(block, prefix, id, i) // encrypt the block
	}
	index[id].label = v // store the label
}

// queryBlockEnc: encrypt one block in query, we only consider "larger than" (block: the block value, prefix: the prefix of this block, blockId: the current block number, gamma: gamma of this query)
func queryBlockEnc(block int64, prefix int64, blockId int, gamma *big.Int) {
	blockStr := strconv.FormatInt(block, 10) + ">"                             // add the inequality operator
	exp := getHashedValue(blockStr, prefix, blockId)                           // calculate the hash value
	tagValue, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(3)).String()) // calculate the tag
	query.tag[blockId] = uint8(tagValue)                                       // store the tag in the list

	// calculate the ciphertext of one block
	query.ciphertext[blockId] = new(bn256.G2).ScalarBaseMult(exp)
	query.ciphertext[blockId].ScalarMult(query.ciphertext[blockId], gamma)
}

// queryEnc: encrypt one value as a query
func queryEnc(v uint32) {
	gamma, _ := rand.Int(rand.Reader, bn256.Order)      // generate gamma
	query.pubKey = new(bn256.G2).ScalarMult(g2s, gamma) // calculate the public key of one query
	// query.pubKey.Neg(query.pubKey)						// in on-chain test, one variable should be negative because the pre-compiled contract 0x08 test the sum of both exponents
	vStr := strconv.FormatInt(int64(v), 2) // calculate the binary value
	vStr = fmt.Sprintf("%032s", vStr)      //pad into 32 bits
	var prefix int64
	for i := 0; i < 32/2; i++ {
		block, _ := strconv.ParseInt(vStr[i*2:i*2+2], 2, 0) // the block contains 2 bits
		if i == 0 {                                         // the first block (no prefix)
			prefix = -1
		} else { // other (has prefix)
			prefix, _ = strconv.ParseInt(vStr[0:i*2], 2, 0)
		}
		queryBlockEnc(block, prefix, i, gamma) // encrypt the block
	}
}

// match: find the index which the condition holds
func match() {
	for i := 0; i < 5; i++ { // scan each index item
		var isMatched bool = false
		for j := 0; j < 16; j++ { // scan each block
			for k := 0; k < 3; k++ { // scan all the blocks which their tags are the same as the query's
				if index[i].tag[j][query.tag[j]][k] == 100 { // if all the aforementioned block are checked
					break
				}
				id := index[i].tag[j][query.tag[j]][k]

				// use the bilinear map
				k1 := bn256.Pair(index[i].ciphertext[j][id], query.pubKey)
				k2 := bn256.Pair(index[i].pubKey[j], query.ciphertext[j])
				k1Byte := k1.Marshal()
				k2Byte := k2.Marshal()

				if bytes.Equal(k1Byte, k2Byte) { // if the bilinear map equation holds, the index is matched
					isMatched = true
					break
				}
			}
			if isMatched == true {
				break
			}
		}
		if isMatched == true {
			fmt.Println(index[i].label) // show the matched index label
		}
	}
}

func main() {
	// generate the parameters
	s, _ := rand.Int(rand.Reader, bn256.Order)
	g1s = new(bn256.G1).ScalarBaseMult(s)
	g2s = new(bn256.G2).ScalarBaseMult(s)

	// add the index items
	indexEnc(500, 0)
	indexEnc(85, 1)
	indexEnc(3, 2)
	indexEnc(150, 3)
	indexEnc(329, 4)

	// generate the query
	queryEnc(500)

	// find the matched index items
	match()
}
