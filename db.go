package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	log "golang.org/x/exp/slog"
)

// TODO create custom db errors
type CarDB interface {
	CreateCar(context.Context, *Car) error
	GetCar(*Car) error
	GetCars() error
	GetCarById(int) (*Car, error)
}

type PostGresStore struct {
	// will handle our DB instance
	db *sql.DB
}

func NewPostgresStore() (*PostGresStore, error) {
	connStr := "host=localhost dbname=postgres user=postgres password=password port=5432 sslmode=disable"

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

// TODO figure out Lakh is easily taken by the (Postgres) money data type
func (p *PostGresStore) createTable() error {
	stmt := `create table if not exists cars (
		id serial primary key,
		company varchar(50),
		model varchar(50), 
		horsepower varchar(50), 
		torque varchar(50), 
		transmission_type varchar(50), 
		drivetrain varchar(50), 
		fuel_economy varchar(50), 
		number_of_doors varchar(50), 
		price varchar(50), 
		start_year integer, 
		emd_year integer,
		body_type varchar(50), 
		engine_type varchar(50), 
		number_of_cylinders varchar(50),
		create_at timestamp 
	)`

	_, err := p.db.Exec(stmt)
	return err
}

func (p *PostGresStore) CreateCar(ctx context.Context, c *Car) error { 
	log.Debug("Inserting a car into DB", "car", c.String())
	insertStmt := `
	INSERT INTO cars (
		id, 
		company, 
		model, 
		horsepower,
		torque, 
		transmission_type, 
		drivetrain, 
		fuel_economy, 
		number_of_doors, 
		price, 
		start_year, 
		end_year, 
		body_type, 
		engine_type, 
		number_of_cyclinders, 
		created_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := p.db.PrepareContext(ctx, insertStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		c.ID, 
		c.Company, 
		c.Model, 
		c.Horsepower, 
		c.Torque, 
		c.TransmissionType, 
		c.Drivetrain, 
		c.FuelEconomy, 
		c.NumberOfDoors,
		c.Price,
		c.StartYear,
		c.EndYear,
		c.BodyType,
		c.EngineType,
		c.NumberofCylinders,
		c.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil 
}

func (p *PostGresStore) GetCar(c *Car) error { return nil }
func (p *PostGresStore) GetCars() error { return nil }
func (p *PostGresStore) GetCarById(int) (*Car, error) {return &Car{}, nil}