package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "golang.org/x/exp/slog"
)

var cars []*Car

// APIServer is a gin RESTful API that will handle incoming requests for Cars.
type APIServer struct {
	db CarDB
	listenAddr string
}

func NewAPIServer(db CarDB, listenAddr string) *APIServer {
	return &APIServer{
		db: db,
		listenAddr: listenAddr,
	}
}

// Ping test
func (a *APIServer) ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cars)
}

// GET endpoints/methods

// getCars returns all cars
func (a *APIServer) getCars(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cars)
}

// getCarById gets a car by the id supplied in the path
func (a *APIServer) getCarById(c *gin.Context) {
	id := c.Param("id")
	car, err := a.carById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Car not found."})
		log.Error("Car not found", "err", err)
		return
	}
	c.IndentedJSON(http.StatusOK, car)
}

// carById is a helper function for retrieving records by id
func(a *APIServer) carById(id string) (*Car, error) {
	return cars[1], nil
}

// POST endpoints/methods

// createCar adds a new car to the db
func (a *APIServer) createCar(c *gin.Context) {
	var (
		newCar Car
		err	error
	) 
	
	if err := c.BindJSON(&newCar); err != nil {
		return
	}

	// TODO check for duplicates. Call to DB with Car given (SELECT-statement)
	if err = a.db.CreateCar(c, &newCar); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not insert Car into DB."})
		return
	}
	c.IndentedJSON(http.StatusCreated, newCar)
}


func (a *APIServer) StartRouter() {
	r := gin.Default()
	
	r.GET("/ping", a.ping)
	r.GET("/cars", a.getCars)
	r.GET("/cars/:id", a.getCarById)
	r.POST("/cars", a.createCar)

	r.Run(a.listenAddr)
}