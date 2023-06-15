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
	var err error

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

	api := NewAPIServer(store, ":9090")
	api.StartRouter()

	// open file
	f, err := os.Open(inputFile)
	if err != nil {
		log.Error("Unable to read/open file", "filename", inputFile)
		panic(err)
	}
	defer f.Close()

	cars, err = CsvReader(f, store)
	if err != nil {
		panic(err)
	}

}
