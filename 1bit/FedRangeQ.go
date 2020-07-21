/*
	FedRangeQ.go
	By Hongcheng Xie at Department of Computer Science, City University of Hong Kong, Hong Kong SAR, China

	This code is the blockchain client in our proposed system, including index encrypting module, index posting module, query encrypting module, query posting module and searching module. In this source code, the block size is set as 4.

	Due to the block size limit in go-ethereum, the code between Line 518 and Line 520 in core/tx_pool.go of go-ethereum should be commented and re-compiled ("make" in project root directory) before deploying our proposed system.
*/
package main

import (
	"bytes"
	"context"
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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// the struct which stores the original data
type TestData struct {
	StockID string
	Price   int
}

var (
	fileName       string         = "../preprocessedList.csv"  // the filename of csv file
	testData       [3500]TestData                              // the original data for our experiment (the top 2000 items)
	indexSize      int            = 400                        // the size of index
	batchSize      int            = 100                        // the batch size of one post
	blockSize      int            = 1                          // the number of bits in one block
	index          []Struct1      = make([]Struct1, indexSize) // the encrypted index corresponding to IndexStru in smart contract
	g1s            *bn256.G1                                   // the key of index
	g2s            *bn256.G2                                   // the key of query
	tPLength       int64                                       // the length of array tagPos
	blockPossValue int64                                       // the number of possible values in one block (2**blockSize)
	query          Struct0                                     // the encrypted query corresponding to QueryStru in smart contract

	url           string = "http://localhost:8545"                                            // the access URL of the test chain
	scAddress     string = "0x48b46e4768cB981ebE9bCaF4992d26CE37CCD9b2"                       // the address of smart contract
	privateKeyStr string = "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19" // the private key of blockchain account
)

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
		if i >= 3500 { // read the top 3500 items in original data
			break
		}
	}
}

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

	var cipherPos = 0                                                          // the current available space of ciphertext list
	gamma, _ := rand.Int(rand.Reader, bn256.Order)                             // the nonce
	index[id].PubKey[blockId] = new(bn256.G1).ScalarMult(g1s, gamma).Marshal() // calculate the public key of one block in one index item

	for i = 0; i < tPLength; i++ { // initialize the tag list with 100 (one value out of range)
		for j := 0; j < int(tPLength); j++ {
			index[id].Tag[blockId][i][j] = 100
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
			tag, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(int64(tPLength))).String()) // the tag
			index[id].Tag[blockId][tag][tagPos[tag]] = uint8(cipherPos)                         // store the list number of the ciphertext in the current available space of corresponding tag
			tagPos[tag]++

			// generate the ciphertext
			t := new(bn256.G1).ScalarBaseMult(exp)
			index[id].Ciphertext[blockId][cipherPos] = new(bn256.G1).ScalarMult(t, gamma).Marshal()
			cipherPos++
		} else if i > block { // the current variable is larger than the current block (the process procedure is similar)
			iStr := strconv.FormatInt(i, 10) + ">"
			exp := getHashedValue(iStr, prefix, blockId)
			// calculate the tag
			tag, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(int64(tPLength))).String()) // the tag
			index[id].Tag[blockId][tag][tagPos[tag]] = uint8(cipherPos)
			tagPos[tag]++

			// generate the ciphertext
			t := new(bn256.G1).ScalarBaseMult(exp)
			index[id].Ciphertext[blockId][cipherPos] = new(bn256.G1).ScalarMult(t, gamma).Marshal()
			cipherPos++
		}
	}
}

// indexItemEnc: encrypt one value as an index item (v: the value to be encrypted, id: the current index number)
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
	index[id].Label = big.NewInt(int64(v)) // store the label
}

// indexEnc: encrypt all the index items
func indexEnc() {
	for i := 0; i < indexSize; i++ {
		indexItemEnc(uint32(testData[i].Price), i)
	}
}

func indexPost(instance *FedRangeQABI, auth *bind.TransactOpts, conn *ethclient.Client) {
	for i := 0; i < len(index)/batchSize; i++ {
		tx, err := instance.Store(auth, index[i*batchSize:i*batchSize+batchSize])
		if err != nil {
			fmt.Println(err)
		}
		ctx := context.Background()
		receipt, err := bind.WaitMined(ctx, conn, tx)
		fmt.Println("Batch ", i, ": ", receipt.Status)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// queryBlockEnc: encrypt one block in query, we only consider "larger than" (block: the block value, prefix: the prefix of this block, blockId: the current block number, gamma: gamma of this query)
func queryBlockEnc(block int64, prefix int64, blockId int, gamma *big.Int) {
	blockStr := strconv.FormatInt(block, 10) + ">"                                    // add the inequality operator
	exp := getHashedValue(blockStr, prefix, blockId)                                  // calculate the hash value
	tagValue, _ := strconv.Atoi(new(big.Int).Mod(exp, big.NewInt(tPLength)).String()) // calculate the tag
	query.Tag[blockId] = uint8(tagValue)                                              // store the tag in the list

	// calculate the ciphertext of one block
	t := new(bn256.G2).ScalarBaseMult(exp)
	query.Ciphertext[blockId] = new(bn256.G2).ScalarMult(t, gamma).Marshal()
}

// queryEnc: encrypt one value as a query
func queryEnc(v uint32) {
	gamma, _ := rand.Int(rand.Reader, bn256.Order) // generate gamma
	t := new(bn256.G2).ScalarMult(g2s, gamma)      // calculate the public key of one query
	query.PubKey = t.Neg(t).Marshal()              // in on-chain test, one variable should be negative because the pre-compiled contract 0x08 test the sum of both exponents
	vStr := strconv.FormatInt(int64(v), 2)         // calculate the binary value
	vStr = fmt.Sprintf("%032s", vStr)              //pad into 32 bits
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

// search: post the query to the test chain and get the result
func search(instance *FedRangeQABI, auth *bind.TransactOpts, conn *ethclient.Client) []*big.Int {
	tx, err := instance.Search(auth, query)
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, conn, tx)
	if err != nil {
		fmt.Println("Status: ", receipt.Status, " ", "Error:", err)
	}
	res, _ := instance.GetResult(nil)
	return res
}

func clearResult(instance *FedRangeQABI, auth *bind.TransactOpts, conn *ethclient.Client) {
	tx, err := instance.ClearResult(auth)
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, conn, tx)
	fmt.Println("status: ", receipt.Status)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	initialize() // initialize the basic parameters
	keyGene()    // generate the key set
	readCSV()    // read the pre-processed csv file

	// initialize the blockchain parameters
	client, _ := ethclient.Dial(url) // connect to the network
	fmt.Println("CONNECT")
	scAddressHex := common.HexToAddress(scAddress)
	instance, _ := NewFedRangeQABI(scAddressHex, client) // instantiate the smart contract
	fmt.Println("Instance Generated")
	privateKey, _ := crypto.HexToECDSA(privateKeyStr)
	fmt.Println("Privatekey Settled")
	auth := bind.NewKeyedTransactor(privateKey) // bind the account
	auth.Nonce = nil
	auth.Value = big.NewInt(0)                  // in wei
	auth.GasLimit = uint64(1844674407370955300) // in units
	auth.GasPrice = big.NewInt(0)

	// the main procedure
	fmt.Println("Index Encryption begins")
	indexEnc() // encrypt the index
	fmt.Println("Index Encryption done")

	fmt.Println("Post begins")
	indexPost(instance, auth, client) // post the index to the test chain
	fmt.Println("Post done")

	fmt.Println("Query Encryption begins")
	queryEnc(uint32(1000000000)) // generate the query
	fmt.Println("Query Encryption done")

	fmt.Println("Searching begins")
	fmt.Println(search(instance, auth, client)) // print the list of matched values
	fmt.Println("Searching Done")

	clearResult(instance, auth, client) // clear the on-chain result list
}
