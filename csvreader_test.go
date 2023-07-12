package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/stretchr/testify/assert"
)

func TestCsvReader(t *testing.T) {
	// Create a temporary CSV file with the test case content
	file, err := os.CreateTemp("", "cars-*.csv")
	if err != nil {
		fmt.Println("Failed to create temporary file:", err)
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	
	if err = addCSVHeaders(file); err != nil {
		t.Fatalf("Could not add csv headers to file: %s", err.Error())
	}
	

	testCases := []struct {
		name           string
		fileContent		[]*CarRecord
		expectedError  bool
		expectedStartYears  []int
		expectedEndYear int
	}{
		{
			name: "Valid File",
			expectedError: false,
			// fileContent: []*CarRecord{
			// 	{
			// 		Car: &Car{
			// 			Company:           "Toyota",
			// 			Model:             "Corolla",
			// 			Horsepower:        "140 hp",
			// 			Torque:            "126 lb-ft",
			// 			TransmissionType:  "Automatic",
			// 			Drivetrain:        "Front-Wheel Drive",
			// 			FuelEconomy:       "31/40 mpg",
			// 			NumberOfDoors:     "4",
			// 			Price:             "$20,000",
			// 			BodyType:          "Sedan",
			// 			EngineType:        "Inline-4",
			// 			NumberofCylinders: "4",
			// 		},
			// 		ModelYearRange: "2015 - 2020",
			// 	},
			// 	{
			// 		Car: &Car{
			// 			Company:           "Honda",
			// 			Model:             "Civic",
			// 			Horsepower:        "158 hp",
			// 			Torque:            "138 lb-ft",
			// 			TransmissionType:  "CVT",
			// 			Drivetrain:        "Front-Wheel Drive",
			// 			FuelEconomy:       "30/38 mpg",
			// 			NumberOfDoors:     "4",
			// 			Price:             "$19,500",
			// 			BodyType:          "Sedan",
			// 			EngineType:        "Inline-4",
			// 			NumberofCylinders: "4",
			// 		},
			// 		ModelYearRange: "2017 - 2022",
			// 	},
			// },
			expectedStartYears:  []int{2015, 2017},
			expectedEndYear: time.Now().Year(),
		},
		{
			name:          "Empty File",
			fileContent:   nil,
			expectedError: true,
		},
		{
			name: "Invalid File",
			fileContent: []*CarRecord{},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock CarDB implementation for testing
			mockDB := &MockDB{}			

			// Use the gocsv package to write the array of CarRecords to the CSV file
			err = gocsv.MarshalFile(&tc.fileContent, file)
			assert.NoError(t, err)
			
			// Need to move cursor back to the top if we're expecting more than just headers
			if !tc.expectedError {
				file.Seek(1,0)
			}
			err = CsvReader(file, mockDB)

			// Assert the expected error
			if tc.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Assert the expected years
			ctx := context.TODO()
			cars, err = mockDB.GetCars(ctx, nil)
			assert.NoError(t, err)
			
			var actualYears []int
			for _, car := range cars{
				actualYears = append(actualYears, car.StartYear)
			}
			assert.ElementsMatch(t, tc.expectedStartYears, actualYears)

			// Assert the expected end year
			// var actualEndYear int
			// if len(mockDB.CreateCarCalls()) > 0 {
			// 	actualEndYear = mockDB.CreateCarCalls()[0].car.EndYear
			// }
			// assert.Equal(t, tc.expectedEndYear, actualEndYear)
		})
	}
}

// Helper function to create a temporary CSV file
func addCSVHeaders(file *os.File) error {
	carRecord := []*CarRecord{}

	// Use the gocsv package to write the array of CarRecords to the CSV file
	if err := gocsv.MarshalFile(carRecord, file); err != nil {
		fmt.Println("Failed to write headers to CSV file:", err)
		return err
	}
	return nil
}

