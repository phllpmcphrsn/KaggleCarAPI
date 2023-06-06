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
	// will handle our DB instance
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

func (p *PostGresStore) Init() error {
	return p.createTable()
}

// TODO split out model_year_range to be start and end years. also, should those years be int or date?
// TODO figure out Lakh is easily taken by the (Postgres) money data type
func (p *PostGresStore) createTable() error {
	query := `create table if not exists cars (
		id serial primary key,
		company varchar(50),
		model varchar(50), 
		horsepower varchar(50), 
		torque varchar(50), 
		transmission_type varchar(50), 
		drivetrain varchar(50), 
		fuel_economy varchar(50), 
		number_of_doors varchar(50), 
		price money, 
		start_year integer, 
		emd_year integer,
		body_type varchar(50), 
		engine_type varchar(50), 
		number_of_cylinders varchar(50),
		create_at timestamp 
	)`

	_, err := p.db.Exec(query)
	return err
}