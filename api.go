package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	docs "github.com/phllpmcphrsn/KaggleCarAPI/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	log "golang.org/x/exp/slog"
)

const basePath = "/api/v1"

// APIServer is a gin RESTful API that will handle incoming requests for Cars.
type APIServer struct {
	db         CarDB
	listenAddr string
}

func NewAPIServer(db CarDB, listenAddr string) *APIServer {
	return &APIServer{
		db:         db,
		listenAddr: listenAddr,
	}
}

// PingExample godoc
//
//	@Summary	ping example
//	@Schemes
//	@Description	do ping
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/example/helloworld [get]
//
// Ping test
func (a *APIServer) ping(c *gin.Context) {
	c.JSON(http.StatusOK, "PONG")
}

// GET endpoints/methods

// getCars returns all cars
func (a *APIServer) getCars(c *gin.Context) {
	cars, err := a.db.GetCars(c)
	if err != nil {
		log.Error("There was an issue retrieving rows of Cars", "err", err)
		return
	}
	c.IndentedJSON(http.StatusOK, cars)
}

// getCarById gets a car by the id supplied in the path
func (a *APIServer) getCarById(c *gin.Context) {
	id := c.Param("id")
	car, err := a.db.GetCarById(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Car not found."})
		log.Error("Car not found", "err", err)
		return
	}
	c.IndentedJSON(http.StatusOK, car)
}

// POST endpoints/methods

// createCar adds a new car to the db
func (a *APIServer) createCar(c *gin.Context) {
	var (
		newCar *Car
		id     int
		err    error
	)

	if err := c.BindJSON(&newCar); err != nil {
		return
	}

	newCar = NewCar(
		newCar.Company,
		newCar.Model,
		newCar.Horsepower,
		newCar.Torque,
		newCar.TransmissionType,
		newCar.Drivetrain,
		newCar.FuelEconomy,
		newCar.NumberOfDoors,
		newCar.Price,
		newCar.BodyType,
		newCar.EngineType,
		newCar.NumberofCylinders,
		newCar.StartYear,
		newCar.EndYear,
	)

	// TODO check for duplicates. Call to DB with Car given (SELECT-statement)
	if id, err = a.db.CreateCar(c, newCar); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not insert Car into DB.", "err": err.Error()})
		return
	}

	newCar.ID = id
	c.IndentedJSON(http.StatusCreated, newCar)
}

func (a *APIServer) StartRouter() {
	r := gin.Default()

	// setup Swagger
	docs.SwaggerInfo.BasePath = basePath
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// setup v1 routes
	v1 := r.Group(basePath)
	{
		v1.GET("/ping", a.ping)
		v1.GET("/cars", a.getCars)
		v1.GET("/cars/:id", a.getCarById)
		v1.POST("/cars", a.createCar)

	}

	r.Run(a.listenAddr)
}
