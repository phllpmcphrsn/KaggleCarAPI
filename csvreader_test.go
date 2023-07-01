package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCsvReader(t *testing.T) {
	testCases := []struct {
		name           string
		fileContent    string
		expectedError  bool
		expectedYears  []int
		expectedEndYear int
	}{
		{
			name: "Valid File",
			fileContent: validCsv,
			expectedError: false,
			expectedYears:  []int{2018, 2020},
			expectedEndYear: time.Now().Year(),
		},
		{
			name:          "Empty File",
			fileContent:   "",
			expectedError: true,
		},
		{
			name: "Invalid File",
			fileContent: invalidCsv,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary CSV file with the test case content
			tempFile, err := os.CreateTemp("", "test.csv")
			defer os.Remove(tempFile.Name())
			defer tempFile.Close()
			assert.NoError(t, err)

			err = createTempCSVFile(tempFile, tc.fileContent)
			assert.NoError(t, err)

			// move the pointer back to the top of tempfile for further reading
			// ret, _ := tempFile.Seek(0,0)
			// t.Log(ret)
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

			// Assert the expected years
			ctx := context.TODO()
			cars, err = mockDB.GetCars(ctx)
			assert.NoError(t, err)
			
			var actualYears []int
			for _, car := range cars{
				actualYears = append(actualYears, car.StartYear)
			}
			assert.ElementsMatch(t, tc.expectedYears, actualYears)

			// Assert the expected end year
			// var actualEndYear int
			// if len(mockDB.CreateCarCalls()) > 0 {
			// 	actualEndYear = mockDB.CreateCarCalls()[0].car.EndYear
			// }
			// assert.Equal(t, tc.expectedEndYear, actualEndYear)
		})
	}
}

// Helper function to create a temporary CSV file with the given content
func createTempCSVFile(tempFile *os.File, content string) error {
	_, err := tempFile.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

var validCsv = []*CarRecord{
	{Car:{"Ferrari","812 Superfast","789 hp","530 lb-ft","7-speed automatic","RWD","13/20 mpg",2,"$366,712","Coupe","6.5L V12",12,},ModelYearRange: "2018 - Present"},
	{Car:{"Ferrari","F8 Tributo","710 hp","568 lb-ft","7-speed automatic","RWD","15/19 mpg",2,"$276,550","Coupe","3.9L V8",8},ModelYearRange: "2020 - Present"},
}
const invalidCsv = `Company,Model,Horsepower,Torque,Transmission Type,Drivetrain,Fuel Economy,Number of Doors,Price,Model Year Range,Body Type,Engine Type,Number of Cylinders`