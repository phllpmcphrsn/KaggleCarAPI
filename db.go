package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "golang.org/x/exp/slog"
)

// TODO create custom db errors
// TODO create functions without ctx
type CarDB interface {
	CreateCar(context.Context, *Car) error
	GetCar(*Car) (*Car, error)
	GetCars(context.Context) ([]*Car, error)
	GetCarById(context.Context, string) (*Car, error)
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
// TODO determine if we should turn engine type into an array since it can have multiple (possibly comma-separated) values
func (p *PostGresStore) createTable() error {
	stmt := `create table if not exists cars (
		id serial primary key,
		company varchar(50),
		model varchar(50), 
		horsepower varchar(50), 
		torque varchar(50), 
		transmission_type varchar(50), 
		drivetrain varchar(50), 
		fuel_economy varchar(250), 
		number_of_doors varchar(50), 
		price varchar(50), 
		start_year integer, 
		end_year integer,
		body_type varchar(50), 
		engine_type varchar(100),
		number_of_cylinders varchar(50),
		created_at timestamp 
	)`

	_, err := p.db.Exec(stmt)
	return err
}

func (p *PostGresStore) Count() (int, error) {
	var count int
	countStmt := "SELECT COUNT(company) from cars"
	row := p.db.QueryRow(countStmt)
	switch err := row.Scan(&count); err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return count, nil
	default:
		return 0, err
	}

}

func (p *PostGresStore) IndexOnCompany(ctx context.Context) error {
	indexStmt := "CREATE INDEX IF NOT EXISTS company_idx ON cars (company)"
	_, err := p.db.ExecContext(ctx, indexStmt)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostGresStore) TableExists() (bool, error) {
	var table string

	existsStmt := "SELECT table_name FROM information_schema.tables WHERE table_name='cars'"
	row := p.db.QueryRow(existsStmt)
	switch err := row.Scan(&table); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (p *PostGresStore) CreateCar(ctx context.Context, c *Car) error { 
	log.Debug("Inserting a car into DB", "car", c.String())
	insertStmt := `
	INSERT INTO cars (
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
		number_of_cylinders, 
		created_at
	)	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err := p.db.ExecContext(
		ctx, 
		insertStmt,
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

func (p *PostGresStore) GetCar(c *Car) (*Car, error) { return nil, nil }

func (p *PostGresStore) GetCarById(ctx context.Context, id string) (*Car, error) {
	var car *Car

    // Query for a value based on a single row.
    err := p.db.QueryRowContext(ctx, "SELECT * from album where id = $1", id).Scan(&car)
	if err != nil {
		// TODO use custom error
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unknown car: %d", id)
		}
		return nil, err
    }
    return car, nil
}

func (p *PostGresStore) GetCars(ctx context.Context) ([]*Car, error) { 
	selectAllStmt := "SELECT * FROM cars"
	stmt, err := p.db.PrepareContext(ctx, selectAllStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	
	cars := []*Car{}
	for rows.Next() {
		car := new(Car)
		if err := rows.Scan(&car); err != nil {
			return cars, nil
		}
		cars = append(cars, car)
	}
	if err = rows.Err(); err != nil {
		return cars, err
	}
	return cars, nil
}

