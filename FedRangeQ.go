package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// the struct which stores the original data
type TestData struct {
	StockID string
	Price   int
}

var (
	fileName string         = "preprocessedList.csv" // the filename of csv file
	testData [2000]TestData                          // the original data for our experiment (the top 2000 items)
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
		if i >= 2000 { // read the top 2000 items in original data
			break
		}
	}
}

func main() {
	readCSV() // read the pre-processed csv file
}
