package contacts

import (
	"encoding/csv"
	"os"
)

func ReadCsv(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	// Create a csv reader
	csvReader := csv.NewReader(file)
	return csvReader.ReadAll()
}

func WriteCsv(filename string, data [][]string) error {
	dataFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	// Write to file
	csvWriter := csv.NewWriter(dataFile)
	csvWriter.WriteAll(data)

	return nil

}
