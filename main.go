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
	store, err := NewPostgresStore()
	if err != nil {
		log.Error("There was an issue reaching the database")
		panic(err)
	}
	log.Info("Connected to database...", "db", store.db.Stats())

	// open file
	f, err := os.Open(inputFile)
	if err != nil {
		log.Error("Unable to read/open file", "filename", inputFile)
		panic(err)
	}
	defer f.Close()
	
	carRecords, err = CsvReader(f)
	if err != nil {
		panic(err)
	}

	StartRouter(":9090", carRecords)
}
