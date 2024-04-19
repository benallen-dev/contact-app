package contacts

import (
	"errors"
	"strconv"
	"log"
	"encoding/csv"
	"os"
)

var filename string = "data/contacts.csv"

func readCsv(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	// Create a csv reader
	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	err = file.Close()
	if err != nil {
		return [][]string{}, err
	}

	return data, nil
}

func writeCsv(filename string, data [][]string) error {

	log.Println("Writing to file: ", filename)
	log.Println("Data: ", data)
	
	dataFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return err
	}

	// Clear the file
	err = dataFile.Truncate(0)

	// Write to file
	csvWriter := csv.NewWriter(dataFile)
	writeErr := csvWriter.WriteAll(data)
	if writeErr != nil {
		return writeErr
	}

	err = dataFile.Close()
	if err != nil {
		return err
	}

	// Verify the write
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	csvReader := csv.NewReader(file)
	checkData, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	if len(checkData) != len(data) {
		return errors.New("Row count mismatch")
	}

	for i, row := range checkData {
		if len(row) != len(data[i]) {
			return errors.New("Row length mismatch at row " + strconv.Itoa(i))
		}

		for j, cell := range row {
			if cell != data[i][j] {
				return errors.New("Cell mismatch at row " + strconv.Itoa(i) + " column " + strconv.Itoa(j))
	}}}			

	return nil

}
