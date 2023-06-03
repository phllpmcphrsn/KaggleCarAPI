package main

import (
	"os"

	log "golang.org/x/exp/slog"
)

const inputFile string = "resources/Car_Models.csv"


func setLogger() {
	logger := log.New(log.NewJSONHandler(os.Stdout, nil))
	log.SetDefault(logger)
}

func main() {
	var (
		carRecords []*CarRecord
		err error
	)

	setLogger()

	// open file
	f, err := os.Open(inputFile)
	if err != nil {
		log.Error("Unable to read/open file", "filename", inputFile, "err", err)
		panic(err)
	}
	defer f.Close()
	
	carRecords, err = CsvReader(f)
	if err != nil {
		panic(err)
	}

	StartRouter(":9090", carRecords)
}
