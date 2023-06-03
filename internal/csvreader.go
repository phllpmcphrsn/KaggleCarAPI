package internal

import (
	"io"

	"github.com/gocarina/gocsv"
	log "golang.org/x/exp/slog"
)

// CarRecord will hold each row from the dataset
type CarRecord struct {
	Company           string `csv:"Company"`
	Model             string `csv:"Model"`
	Horsepower        string `csv:"Horsepower"`
	Torque            string `csv:"Torque"`
	TransmissionType  string `csv:"Transmission Type"`
	Drivetrain        string `csv:"Drivetrain"`
	FuelEconomy       string `csv:"Fuel Economy"`
	NumberOfDoors     string `csv:"Number of Doors"`
	Price             string `csv:"Price"`
	ModelYearRange    string `csv:"Model Year Range"`
	BodyType          string `csv:"Body Type"`
	EngineType        string `csv:"Engine Type"`
	NumberofCylinders string `csv:"Number of Cylinders"`
}

func CsvReader(in io.Reader, inputFile string) ([]*CarRecord, error) {
	// store each row into a object
	cars := []*CarRecord{}
	if err := gocsv.Unmarshal(in, &cars); err != nil {
		log.Error("Unable to unmarshal file contents", "filename", inputFile, "err", err)
		return nil, err
	}

	return cars, nil
}
