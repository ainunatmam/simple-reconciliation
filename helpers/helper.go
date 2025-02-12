package helpers

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadFile(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to read input file "+path, err)
	}
	record,err := ReadCsv(*f)
	if err != nil {
		log.Fatal("Unable to read input file "+path, err)
	}
	defer f.Close()
	return record
}

func ReadCsv(file os.File) ([][]string, error)  {
	reader := csv.NewReader(&file)
	record, err := reader.ReadAll()
	if err != nil {
        log.Fatal("Unable to parse file as CSV", err)
    }
	return record, err
}