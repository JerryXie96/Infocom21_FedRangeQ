// this source code is to test the matching algorithm locally. this is the 2-bit version.
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"

	"github.com/clearmatics/bn256"
)

// the struct which stores the original data
type TestData struct {
	StockID string
	Price   int
}

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
	indexSize      int            = 100                          // the number of indexes
	blockSize      int            = 2                            // the number of bits in one block
	fileName       string         = "../../preprocessedList.csv" // the filename of csv file
	testData       [2000]TestData                                // the original data for our experiment (the top 2000 items)
	index          []IndexStru    = make([]IndexStru, indexSize) // the index list
	g1s            *bn256.G1                                     // the key of index
	g2s            *bn256.G2                                     // the key of query
	tPLength       int64                                         // the length of array tagPos
	blockPossValue int64                                         // the number of possible values in one block (2**blockSize)
	query          QueryStru                                     // the query
)

// initialize: initialize some basic parameters
func initialize() {
	tPLength = int64(math.Exp2(float64(blockSize)) - 1) // the length of array tagPos
	blockPossValue = tPLength + 1
}

// keyGene: generate the keys of both index and query (g1s and g2s)
func keyGene() {
	s, _ := rand.Int(rand.Reader, bn256.Order)
	g1s = new(bn256.G1).ScalarBaseMult(s)
	g2s = new(bn256.G2).ScalarBaseMult(s)
}

// readCSV(): read the pre-processed csv file
func readCSV() {
	var i uint = 0
	file, _ := os.Open(fileName)
	reader := csv.NewReader(file)
	csvContent, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Read CSV file error!")
	}
	for _, row := range csvContent {
		if row[0] == "" { // skip the first row (the title)
			continue
		}
		testData[i].StockID = row[1]
		testData[i].Price, _ = strconv.Atoi(row[2]) // convert the string in csv file to int
		i++
		if i >= 2000 { // read the top 2000 items in original data
			break
		}
	}
}

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

	var tagPos []int = make([]int, tPLength) // the current available space of each tag's list
	for i = 0; i < tPLength; i++ {
		tagPos[i] = 0
	}

	var cipherPos = 0                                                // the current available space of ciphertext list
	gamma, _ := rand.Int(rand.Reader, bn256.Order)                   // the nonce
	index[id].pubKey[blockId] = new(bn256.G1).ScalarMult(g1s, gamma) // calculate the public key of one block in one index item
	for i = 0; i < tPLength; i++ {                                   // initialize the tag list with 100 (one value out of range)
		for j := 0; j < int(tPLength); j++ {
			index[id].tag[blockId][i][j] = 100
		}
	}
	// iStr includes the operator > or <
	for i = 0; i < blockPossValue; i++ {
		if i == block { // do not encrypt the equal block
			continue
		} else if i < block { // the current variable is smaller than the current block
			iStr := strconv.FormatInt(i, 10) + "<"       // add the inequality operator into the string to be hashed
			exp := getHashedValue(iStr, prefix, blockId) // calculate the hash value in tag and ciphertext
			// calculate the tag
			tag, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(tPLength)).String()) // the tag
			index[id].tag[blockId][tag][tagPos[tag]] = uint8(cipherPos)                  // store the list number of the ciphertext in the current available space of corresponding tag
			tagPos[tag]++

			// generate the ciphertext
			index[id].ciphertext[blockId][cipherPos] = new(bn256.G1).ScalarBaseMult(exp)
			index[id].ciphertext[blockId][cipherPos].ScalarMult(index[id].ciphertext[blockId][cipherPos], gamma)
			cipherPos++
		} else if i > block { // the current variable is larger than the current block (the process procedure is similar)
			iStr := strconv.FormatInt(i, 10) + ">"
			exp := getHashedValue(iStr, prefix, blockId)
			// calculate the tag
			tag, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(tPLength)).String()) // the tag
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
func indexItemEnc(v uint32, id int) {
	vStr := strconv.FormatInt(int64(v), 2) // calculate the binary value
	vStr = fmt.Sprintf("%032s", vStr)      // pad to 32 bits
	var prefix int64
	for i := 0; i < 32/blockSize; i++ {
		block, _ := strconv.ParseInt(vStr[i*blockSize:i*blockSize+blockSize], 2, 0) // the block contains 2 bits
		if i == 0 {                                                                 // the first block (no prefix)
			prefix = -1
		} else { // other (has prefix)
			prefix, _ = strconv.ParseInt(vStr[0:i*blockSize], 2, 0)
		}
		indexBlockEnc(block, prefix, id, i) // encrypt the block
	}
	index[id].label = v // store the label
}

// indexEnc: encrypt all the index items
func indexEnc() {
	for i := 0; i < indexSize; i++ {
		indexItemEnc(uint32(testData[i].Price), i)
	}
}

// queryBlockEnc: encrypt one block in query, we only consider "larger than" (block: the block value, prefix: the prefix of this block, blockId: the current block number, gamma: gamma of this query)
func queryBlockEnc(block int64, prefix int64, blockId int, gamma *big.Int) {
	blockStr := strconv.FormatInt(block, 10) + ">"                                    // add the inequality operator
	exp := getHashedValue(blockStr, prefix, blockId)                                  // calculate the hash value
	tagValue, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(tPLength)).String()) // calculate the tag
	query.tag[blockId] = uint8(tagValue)                                              // store the tag in the list

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
	for i := 0; i < 32/blockSize; i++ {
		block, _ := strconv.ParseInt(vStr[i*blockSize:i*blockSize+blockSize], 2, 0) // the block contains 2 bits
		if i == 0 {                                                                 // the first block (no prefix)
			prefix = -1
		} else { // other (has prefix)
			prefix, _ = strconv.ParseInt(vStr[0:i*blockSize], 2, 0)
		}
		queryBlockEnc(block, prefix, i, gamma) // encrypt the block
	}
}

// match: find the index which the condition holds
func match() {
	for i := 0; i < indexSize; i++ { // scan each index item
		var isMatched bool = false
		for j := 0; j < 32/blockSize; j++ { // scan each block
			for k := 0; k < int(tPLength); k++ { // scan all the blocks which their tags are the same as the query's
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
	initialize() // initialize the basic parameters
	keyGene()    // generate the key set
	readCSV()    // read the pre-processed csv file

	// the main procedure
	fmt.Println("Index Encryption begins")
	indexEnc() // encrypt the index
	fmt.Println("Index Encryption done")

	fmt.Println("Query Encryption begins")
	queryEnc(uint32(testData[1000].Price)) // generate the query
	fmt.Println("Query Encryption done")

	fmt.Println("Searching begins")
	match() // print the list of matched values
	fmt.Println("Searching Done")
}
