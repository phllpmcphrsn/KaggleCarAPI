package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCsvReader(t *testing.T) {
	testCases := []struct {
		name           string
		fileContent    string
		expectedError  bool
		expectedCars   []*Car
		expectedYears  []int
		expectedEndYear int
	}{
		{
			name: "Valid File",
			fileContent: `make,model,model_year_range
Toyota,Corolla,2008 - 2012
Ford,Fusion,2015 - present
Honda,Civic,1999 - 2005
`,
			expectedError: false,
			expectedCars: []*Car{
				{Make: "Toyota", Model: "Corolla", StartYear: 2008, EndYear: 2012},
				{Make: "Ford", Model: "Fusion", StartYear: 2015, EndYear: time.Now().Year()},
				{Make: "Honda", Model: "Civic", StartYear: 1999, EndYear: 2005},
			},
			expectedYears:  []int{2008, 2015, 1999},
			expectedEndYear: time.Now().Year(),
		},
		{
			name:          "Empty File",
			fileContent:   "",
			expectedError: true,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary CSV file with the test case content
			tempFile, err := createTempCSVFile(tc.fileContent)
			assert.NoError(t, err)
			defer os.Remove(tempFile.Name())

			// Create a mock CarDB implementation for testing
			mockDB := &MockDB{}

			// Invoke the CsvReader function with the temporary file and mock DB
			err = CsvReader(tempFile, mockDB)

			// Assert the expected error
			if tc.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Assert the expected cars
			assert.Equal(t, len(tc.expectedCars), len(mockDB.CreateCar()))
			for i, call := range mockDB.CreateCarCalls() {
				assert.Equal(t, tc.expectedCars[i], call.ctx)
				assert.Equal(t, tc.expectedCars[i], call.car)
			}

			// Assert the expected years
			var actualYears []int
			for _, call := range mockDB.CreateCarCalls() {
				actualYears = append(actualYears, call.car.StartYear)
			}
			assert.ElementsMatch(t, tc.expectedYears, actualYears)

			// Assert the expected end year
			var actualEndYear int
			if len(mockDB.CreateCarCalls()) > 0 {
				actualEndYear = mockDB.CreateCarCalls()[0].car.EndYear
			}
			assert.Equal(t, tc.expectedEndYear, actualEndYear)
		})
	}
}

// Helper function to create a temporary CSV file with the given content
func createTempCSVFile(content string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "test.csv")
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString(content)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

const validCsv = `Company,Model,Horsepower,Torque,Transmission Type,Drivetrain,Fuel Economy,Number of Doors,Price,Model Year Range,Body Type,Engine Type,Number of Cylinders
Ferrari,812 Superfast,789 hp,530 lb-ft,7-speed automatic,RWD,13/20 mpg,2,"$366,712",2018 - Present,Coupe,6.5L V12,12
Ferrari,F8 Tributo,710 hp,568 lb-ft,7-speed automatic,RWD,15/19 mpg,2,"$276,550",2020 - Present,Coupe,3.9L V8,8
`
const invalidCsv = `Company,Model,Horsepower,Torque,Transmission Type,Drivetrain,Fuel Economy,Number of Doors,Price,Model Year Range,Body Type,Engine Type,Number of Cylinders`
// MockDB is a mock implementation
