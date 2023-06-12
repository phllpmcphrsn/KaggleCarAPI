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
		cars []*Car
		err error
	)

	setLogger()
	store, err := NewPostgresStore()
	if err != nil {
		log.Error("There was an issue reaching the database", "err", err)
		panic(err)
	}
	log.Info("Connected to database...", "db", store.db.Stats())

	if err := store.Init(); err != nil {
		log.Error("There was an issue initializing the database", "err", err)
		panic(err)
	}
	
	// open file
	f, err := os.Open(inputFile)
	if err != nil {
		log.Error("Unable to read/open file", "filename", inputFile)
		panic(err)
	}
	defer f.Close()
	
	cars, err = CsvReader(f)
	if err != nil {
		panic(err)
	}

	StartRouter(":9090", cars)
}
