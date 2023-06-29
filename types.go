package main

import (
	"fmt"
	"time"
)

type CarRecord struct {
	*Car
	ModelYearRange    string    `csv:"Model Year Range"`
}

type Car struct {
	ID                int       `json:"id"`
	Company           string    `csv:"Company" json:"company"`
	Model             string    `csv:"Model" json:"model"`
	Horsepower        string    `csv:"Horsepower" json:"horsepower"`
	Torque            string    `csv:"Torque" json:"torque"`
	TransmissionType  string    `csv:"Transmission Type" json:"transmissionType"`
	Drivetrain        string    `csv:"Drivetrain" json:"drivetrain"`
	FuelEconomy       string    `csv:"Fuel Economy" json:"fuelEconomy"`
	NumberOfDoors     string    `csv:"Number of Doors" json:"numberOfDoors"`
	Price             string    `csv:"Price" json:"price"`
	StartYear         int       `json:"startYear"`
	EndYear           int       `json:"endYear"`
	BodyType          string    `csv:"Body Type" json:"bodyType"`
	EngineType        string    `csv:"Engine Type" json:"engineType"`
	NumberofCylinders string    `csv:"Number of Cylinders" json:"numberOfCylinders"`
	CreatedAt         time.Time `json:"createdAt"`
}

// NewCar creates a new Car instance with the given parameters
func NewCar(company, model, horsepower, torque, transmissionType, drivetrain, fuelEconomy, numberOfDoors, price, bodyType, engineType, numberOfCylinders string, startYear, endYear int) *Car {
	return &Car{
		Company:           company,
		Model:             model,
		Horsepower:        horsepower,
		Torque:            torque,
		TransmissionType:  transmissionType,
		Drivetrain:        drivetrain,
		FuelEconomy:       fuelEconomy,
		NumberOfDoors:     numberOfDoors,
		Price:             price,
		StartYear:         startYear,
		EndYear:           endYear,
		BodyType:          bodyType,
		EngineType:        engineType,
		NumberofCylinders: numberOfCylinders,
		CreatedAt:         time.Now().UTC(),
	}
}

func (c *Car) String() string {
	return fmt.Sprintf("%s %s", c.Company, c.Model)
}