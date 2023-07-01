package main

import (
	"net/http"
	"os"

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
	env        string
}

func NewAPIServer(db CarDB, listenAddr, env string) *APIServer {
	return &APIServer{
		db:         db,
		listenAddr: listenAddr,
		env:        env,
	}
}

// Ping godoc
//
//	@Summary		Ping example
//	@Description	Endpoint to test for liveness. It simply returns "PONG"
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	PONG
//	@Router			/ping [get]
func (a *APIServer) ping(c *gin.Context) {
	c.JSON(http.StatusOK, "PONG")
}

// GET endpoints/methods

// GetCars godoc
//
//	@Summary		Get Cars array
//	@Description	Responds with the list of all cars as JSON
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	Car	"ok"
//	@Router			/cars [get]
func (a *APIServer) getCars(c *gin.Context) {
	cars, err := a.db.GetCars(c)
	if err != nil {
		log.Error("There was an issue retrieving rows of Cars", "err", err)
		return
	}
	c.IndentedJSON(http.StatusOK, cars)
}

// GetCarById godoc
//
//	@Summary		Get single car by id
//	@Description	Returns the car with the given id
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"search by id"
//	@Success		200	{object}	Car		"ok"
//	@Router			/cars/{id} [get]
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

// CreateCar godoc
//
//	@Summary		Store a new car
//	@Description	Takes a car JSON and stores in DB. Returned saved JSON
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			car	body		Car	true	"Car JSON"
//	@Success		200	{object}	Car	"ok"
//	@Failure		400	{object}	map[string]any
//	@Failure		500	{object}	map[string]any
//	@Router			/cars/ [post]
func (a *APIServer) createCar(c *gin.Context) {
	var (
		newCar *Car
		id     int
		err    error
	)

	if err := c.BindJSON(&newCar); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Received bad request."})
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
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not insert Car into DB."})
		log.Error("Could not insert Car into DB", "err", err)
		return
	}

	newCar.ID = id
	c.IndentedJSON(http.StatusCreated, newCar)
}

func (a *APIServer) StartRouter() {
	r := gin.Default()
	if os.Getenv(gin.EnvGinMode) == "" {
		mode := ginEnvMode(a.env)
		gin.SetMode(mode) // set this based on production or development env
	}
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

func ginEnvMode(env string) string {
	switch env {
	case "prod":
		return gin.ReleaseMode
	case "dev":
		return gin.DebugMode
	default:
		return gin.TestMode
	}
}
