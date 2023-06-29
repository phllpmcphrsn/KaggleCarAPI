package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestPing tests the ping handler
func TestPing(t *testing.T) {
	// create a mock gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	
	// Initialize the db field with a mock database implementation for testing
	a := NewAPIServer(&MockDB{}, "")

	// call the ping handler
	a.ping(c)

	// assert the response status and body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "\"PONG\"", w.Body.String())
}

// TestGetCars tests the getCars handler
func TestGetCars(t *testing.T) {
	cars := []Car{
		{ID: 1, Company: "Toyota", Model: "Corolla"},
		{ID: 2, Company: "Ford", Model: "F150"},
		{ID: 3, Company: "Chevrolet", Model: "Cobalt"},
	}

	// create a mock gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	a := NewAPIServer(&MockDB{}, "")

	// call the getCars handler
	a.getCars(c)

	// assert the response status and body
	assert.Equal(t, http.StatusOK, w.Code)
	
	var actualBody []Car
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)
	if assert.NoError(t, err){
		assert.Equal(t, cars, actualBody)
	}
}

func TestGetCarById(t *testing.T) {
	// Define the test cases as a table
	testCases := []struct {
		name           string
		carID          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Car ID",
			carID:          "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id": 1, "company": "Toyota", "model": "Corolla", "horsepower": "", "torque": "", "transmissionType": "", "drivetrain": "", "fuelEconomy": "", "numberOfDoors": "", "price": "", "startYear": 0, "endYear": 0, "bodyType": "", "engineType": "", "numberOfCylinders": "", "createdAt": "0001-01-01T00:00:00Z"}`,
		},
		{
			name:           "Invalid Car ID",
			carID:          "456",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message": "Car not found."}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new Gin router and set up the route for testing
			router := gin.Default()
			api := NewAPIServer(&MockDB{}, "")
			router.GET("/cars/:id", api.getCarById)

			// Create a new HTTP request with the test car ID
			req, _ := http.NewRequest("GET", "/cars/"+tc.carID, nil)
			rec := httptest.NewRecorder()

			// Serve the HTTP request
			router.ServeHTTP(rec, req)

			// Check the response status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert that the response body matches the expected value
			// First we'll treat both as JSON and unmarshal to a map
			var actualBody map[string]interface{}
			var expectedBody map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &actualBody)
			tcErr := json.Unmarshal([]byte(tc.expectedBody), &expectedBody)

			// If no errors occurred during unmarshalling then make assertions about bodies
			if assert.NoError(t, err) && assert.NoError(t, tcErr){
				assert.Equal(t, expectedBody, actualBody)
			}
		})
	}
}

func TestCreateCar(t *testing.T) {
	// Define the test cases as a table
	testCases := []struct {
		name           string
		expectedStatus int
		expectedBody   string
		requestBody	string
	}{
		{
			name:           "Valid Car",
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id": 1, "company": "Toyota", "model": "Corolla", "horsepower": "", "torque": "", "transmissionType": "", "drivetrain": "", "fuelEconomy": "", "numberOfDoors": "", "price": "", "startYear": 0, "endYear": 0, "bodyType": "", "engineType": "", "numberOfCylinders": "", "createdAt": "0001-01-01T00:00:00Z"}`,
			requestBody: `{"company": "Toyota", "model": "Corolla", "horsepower": "", "torque": "", "transmissionType": "", "drivetrain": "", "fuelEconomy": "", "numberOfDoors": "", "price": "", "startYear": 0, "endYear": 0, "bodyType": "", "engineType": "", "numberOfCylinders": ""}`,
		},
		{
			name:           "Invalid Car",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message": "Received bad request."}`,
			requestBody: `{"id":}`,
		},
		{
			name:           "Storage Issue",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message": "Could not insert Car into DB."}`,
			requestBody: `{"company": "BadCompany", "model": "", "horsepower": "", "torque": "", "transmissionType": "", "drivetrain": "", "fuelEconomy": "", "numberOfDoors": "", "price": "", "startYear": 0, "endYear": 0, "bodyType": "", "engineType": "", "numberOfCylinders": ""}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create a mock gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body := bytes.NewReader([]byte(tc.requestBody))
			c.Request = httptest.NewRequest("POST", "/api/v1/cars", body)

			a := NewAPIServer(&MockDB{}, "")
			a.createCar(c)

			// Check the response status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Assert that the response body matches the expected value
			// First we'll treat both as JSON and unmarshal to a map
			var actualBody map[string]interface{}
			var expectedBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &actualBody)
			tcErr := json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			
			// Removing the timestamps from comparison
			delete(actualBody, "createdAt")
			delete(expectedBody, "createdAt")

			// If no errors occurred during unmarshalling then make assertions about bodies
			if assert.NoError(t, err) && assert.NoError(t, tcErr){
				assert.Equal(t, expectedBody, actualBody)
			}
		})
	}
}