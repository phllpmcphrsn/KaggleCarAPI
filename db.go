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
	CreateCar(context.Context, *Car) (int, error)
	GetCars(context.Context, *Pagination) ([]*Car, error)
	GetCarById(context.Context, string) (*Car, error)
	Count() (int, error)
}

type PostGresStore struct {
	// will handle our DB instance
	db *sql.DB
}

func NewPostgresStore(c *Config, creds *Credentials) (*PostGresStore, error) {
	dbConfig := c.Database
	
	var ssl string
	if dbConfig.SSL.Enabled {
		ssl = "enabled"
	} else {
		ssl = "disabled"
	}

	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=%s", dbConfig.Host, dbConfig.Name, creds.Username, string(creds.Password), dbConfig.Port, ssl)

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
	log.Error("An error occured while creating the cars table", "err", err)
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
		log.Error("Index for column not created", "err", err)
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

func (p *PostGresStore) CreateCar(ctx context.Context, car *Car) (int, error) { 
	log.Debug("Inserting a car into DB", "car", car.String())
	var id int

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
	)	
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id`

	err := p.db.QueryRowContext(
		ctx, 
		insertStmt,
		&car.Company, 
		&car.Model, 
		&car.Horsepower, 
		&car.Torque, 
		&car.TransmissionType, 
		&car.Drivetrain, 
		&car.FuelEconomy, 
		&car.NumberOfDoors,
		&car.Price,
		&car.StartYear,
		&car.EndYear,
		&car.BodyType,
		&car.EngineType,
		&car.NumberofCylinders,
		&car.CreatedAt,
	).Scan(&id)
	
	if err != nil {
		log.Error("An error occurred while inserting to db", "err", err)
		return 0, err
	}

	log.Debug("Succesfully inserted row", "id", id)
	return id, nil 
}

func (p *PostGresStore) GetCarById(ctx context.Context, id string) (*Car, error) {
	var car Car

    // Query for a value based on a single row.
    err := p.db.QueryRowContext(ctx, "SELECT * FROM cars WHERE id = $1", id).Scan(
		&car.ID,
		&car.Company, 
		&car.Model, 
		&car.Horsepower, 
		&car.Torque, 
		&car.TransmissionType, 
		&car.Drivetrain, 
		&car.FuelEconomy, 
		&car.NumberOfDoors,
		&car.Price,
		&car.StartYear,
		&car.EndYear,
		&car.BodyType,
		&car.EngineType,
		&car.NumberofCylinders,
		&car.CreatedAt,
	)
	if err != nil {
		// TODO use custom error
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("car not found: %s", id)
		}
		return nil, err
    }
    return &car, nil
}

// TODO implement pagination
func (p *PostGresStore) GetCars(ctx context.Context, page *Pagination) ([]*Car, error) { 
	selectAllStmt := "SELECT * FROM cars"
	if page != nil {
		if page.Limit > 0 {
			selectAllStmt = selectAllStmt + " LIMIT $1 "
		}
		if page.Offset > 1 {
			selectAllStmt = selectAllStmt + " OFFSET $1 "
		}
	}
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
		err = rows.Scan(		
			&car.ID,
			&car.Company, 
			&car.Model, 
			&car.Horsepower, 
			&car.Torque, 
			&car.TransmissionType, 
			&car.Drivetrain, 
			&car.FuelEconomy, 
			&car.NumberOfDoors,
			&car.Price,
			&car.StartYear,
			&car.EndYear,
			&car.BodyType,
			&car.EngineType,
			&car.NumberofCylinders,
			&car.CreatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	// catch any errors that may have occurred during loop
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return cars, nil
}

