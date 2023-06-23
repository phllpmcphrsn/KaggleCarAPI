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

//	@title			Kaggle 2023 Car Models API
//	@version		1.0
//	@description	REST API for Kaggle 2023 Car Models Dataset
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	https://github.com/phllpmcphrsn/KaggleCarAPI/issues
//	@contact.email	phllpmcphrsn@yahoo.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:9090
//	@BasePath	/api/v1
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

	// want to check if table has any elements prior to read and populating from csv
	// if it does we'll assume that it's already been populated with data from csv
	count, err := store.Count()
	if err != nil {
		log.Error("An error occured while checking for table's count", "err", err)
		panic(err)
	} else if count == 0 {
		log.Info("Populating cars table...")
		go readCsv(store)
	}

	api := NewAPIServer(store, ":9090")
	api.StartRouter()

}

func readCsv(store *PostGresStore) {
	f, err := os.Open(inputFile)
	if err != nil {
		log.Error("Unable to read/open file", "filename", inputFile)
		panic(err)
	}
	defer f.Close()

	err = CsvReader(f, store)
	if err != nil {
		panic(err)
	}
}
