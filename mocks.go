package main

import (
	"context"
	"fmt"
	"strconv"
)

type MockDB struct{}

func (m *MockDB) CreateCar(c context.Context, car *Car) (int, error) {
	if car.Company == "BadCompany" {
		return 0, fmt.Errorf("Error")
	}
	return 1, nil
}

func (m *MockDB) GetCars(context.Context) ([]*Car, error) {
	cars := []*Car{
		{ID: 1, Company: "Toyota", Model: "Corolla"},
		{ID: 2, Company: "Ford", Model: "F150"},
		{ID: 3, Company: "Chevrolet", Model: "Cobalt"},
	}
	return cars, nil
}
func (m *MockDB) GetCarById(c context.Context, id string) (*Car, error) {
	var car *Car
	if id == "1" {
		i, _ := strconv.Atoi(id)
		car = &Car{ID: i, Company: "Toyota", Model: "Corolla"}
		return car, nil
	}
	return nil, fmt.Errorf("car not found: %s", id)
}