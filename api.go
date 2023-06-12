package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "golang.org/x/exp/slog"
)

var cars []*Car

// Ping test
func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cars)
}

// GET endpoints/methods

// getCars returns all cars
func getCars(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cars)
}

// getCarById gets a car by the id supplied in the path
func getCarById(c *gin.Context) {
	id := c.Param("id")
	car, err := carById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Car not found."})
		log.Error("Car not found", "err", err)
		return
	}
	c.IndentedJSON(http.StatusOK, car)
}

// carById is a helper function for retrieving records by id
func carById(id string) (*Car, error) {
	return cars[1], nil
}

// POST endpoints/methods

// createCar adds a new car to the db
func createCar(c *gin.Context) {
	var newCar Car

	if err := c.BindJSON(&newCar); err != nil {
		return
	}

	cars = append(cars, &newCar)
	c.IndentedJSON(http.StatusCreated, newCar)
}


func StartRouter(port string, c []*Car) {
	cars = c
	r := gin.Default()
	
	r.GET("/ping", ping)
	r.GET("/cars", getCars)
	r.GET("/cars/:id", getCarById)
	r.POST("/cars", createCar)

	r.Run(port)
}