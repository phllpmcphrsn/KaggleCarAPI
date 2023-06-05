package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type CarDB interface {
	CreateCar(*CarRecord) error
	GetCars() error
	GetCarById(int) (*CarRecord, error)
}

type PostGresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostGresStore, error) {
	connStr := "user=postgres dbname=postgres password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostGresStore{db: db}, nil
}

func (p *PostGresStore) init() error {
	return nil
}

func (p *PostGresStore) createTable error {
	query := `create table car if not exists (
		id serial primary key
		company varcha(50)
		model             string 
		horsepower        string 
		torque            string 
		transmissionType  string 
		drivetrain        string 
		fuelEconomy       string 
		numberOfDoors     string 
		price             string 
		modelYear_range    string 
		body_type          string 
		engine_type        string 
		number_of_cylinders string 
	)`
}