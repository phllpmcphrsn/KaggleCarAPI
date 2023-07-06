package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	log "golang.org/x/exp/slog"
)

func CsvReader(file *os.File, db CarDB) error {
	carRecords := []*CarRecord{}
	if err := gocsv.Unmarshal(file, &carRecords); err != nil {
		log.Error("Unable to unmarshal file contents", "filename", file.Name(), "err", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	// set all created times to the same time
	createdAt := time.Now().UTC()

	// loop through each record that was unmarshalled from the CSV then clean
	// up the year range such that it is stored in as start- and end-year
	for _, carRecord := range carRecords {
		err := clean(carRecord); if err != nil {
			log.Error("There was an issue cleaning a CarRecord", "car", carRecord)
			return err
		}
		
		car := carRecord.Car
		car.CreatedAt = createdAt
		_, err = db.CreateCar(ctx, car); if err != nil {
			log.Error("Could not insert Car into database", "car", car.String(), "err", err)
			return err
		}
	}

	return nil
}

// clean takes a car and cleans up the data for the model year range
func clean(c *CarRecord) error {
	err := cleanYears(c)
	// err = cleanPrice(c)
	return err
}

func cleanYears(c *CarRecord) error {
	var err error
	
	trimmedYearRange := strings.ReplaceAll(c.ModelYearRange, " ", "")
	if trimmedYearRange == "" {
		return nil
	}

	// the year range may have spacing between the hyphen and years (i.e. 2008 - 2012)
	// we eliminate those then split about the hyphen to produce a two-element array
	// (i.e. [2008, 2012])
	yearRange := strings.Split(trimmedYearRange, "-")
	car := c.Car
	car.StartYear, err = strconv.Atoi(yearRange[0]); if err != nil {
		log.Error("There was an issue cleaning the Model Year Range for the starting year", "err", err)
		return err
	}

	// catch the case where there's no ending year
	if len(yearRange) > 1 {
		endYear := yearRange[1]	
		car.EndYear, err = checkEndYear(endYear); if err != nil {
			log.Error("There was an issue cleaning the Model Year Range for the ending year", "err", err)
			return err
		}
	}

	return nil
}

// checkEndYear checks if the ending year is "present"; if so replace "present" with the current year
func checkEndYear(endYear string) (int, error) {
	endYear = strings.ToLower(endYear)
	if endYear == "present" || endYear == "" {
		return time.Now().Year(), nil
	}
	return strconv.Atoi(endYear)
}

// TODO figure out if we even want to clean the price
// func cleanPrice(p string) {}