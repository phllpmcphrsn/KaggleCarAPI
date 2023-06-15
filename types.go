package main

import (
	"fmt"
	"time"
)

// CreateCarRequest
type CreateCarRequest struct {
	Company           string `json:"Company"`
	Model             string `json:"Model"`
	Horsepower        string `json:"Horsepower"`
	Torque            string `json:"Torque"`
	TransmissionType  string `json:"Transmission Type"`
	Drivetrain        string `json:"Drivetrain"`
	FuelEconomy       string `json:"Fuel Economy"`
	NumberOfDoors     string `json:"Number of Doors"`
	Price             string `json:"Price"`
	StartYear         int    `json:"startYear"`
	EndYear           int    `json:"endYear"`
	BodyType          string `json:"Body Type"`
	EngineType        string `json:"Engine Type"`
	NumberofCylinders string `json:"Number of Cylinders"`
}

// TODO split out model_year_range to be start and end years. also, should those years be int or date?
// TODO transform this into a Car instead of a CarRecord. Essentially, we'll do any and all trasformations
// within this struct instead of making two separate structs (i.e. separating year range, money conversions)
// CarRecord will hold each row from the dataset
type CarRecord struct {
	Car *Car
	ModelYearRange    string    `csv:"Model Year Range"`
}

type Car struct {
	ID                int       `json:"id"`
	Company           string    `csv:"Company" json:"company" binding:"required"`
	Model             string    `csv:"Model" json:"model" binding:"required"`
	Horsepower        string    `csv:"Horsepower" json:"horsepower"`
	Torque            string    `csv:"Torque" json:"torque"`
	TransmissionType  string    `csv:"Transmission Type" json:"transmissionType"`
	Drivetrain        string    `csv:"Drivetrain" json:"drivetrain"`
	FuelEconomy       string    `csv:"Fuel Economy" json:"fuelEconomy"`
	NumberOfDoors     string    `csv:"Number of Doors" json:"numberOfDoors" binding:"required"`
	Price             string    `csv:"Price" json:"price" binding:"required"`
	StartYear         int       `json:"startYear" binding:"required"`
	EndYear           int       `json:"endYear"`
	BodyType          string    `csv:"Body Type" json:"bodyType"`
	EngineType        string    `csv:"Engine Type" json:"engineType"`
	NumberofCylinders string    `csv:"Number of Cylinders" json:"numberOfCylinders"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewCar(company string, model string, price string, year int, numOfDoors string) *Car {
	return &Car{
		Company: company,
		Model: model,
		Price: price,
		StartYear: year,
		NumberOfDoors: numOfDoors,
	}
}

func (c *Car) String() string {
	return fmt.Sprintf("%d %s %s", c.ID, c.Company, c.Model)
}