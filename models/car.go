package models

import (
	"fmt"

	"github.com/karanjar/cargobackend_fibre_framework.git/config"
)

type Car struct {
	Id    int     `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string  `json:"name" bson:"name"`
	Model string  `json:"model" bson:"model"`
	Year  int64   `json:"year" bson:"year"`
	Price float64 `json:"price" bson:"price"`
}

type Error struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func (c *Car) Insert() error {
	query := `INSERT INTO cars (name,model,year,price) 
				VALUES (:name,:model,:year,:price)
				returning id;`

	row, err := config.Db.NamedQuery(query, c)
	if err != nil {
		fmt.Printf("Error inserting car: %v\n", err)
	}

	defer row.Close()

	if row.Next() {
		if err := row.Scan(&c.Id); err != nil {
			fmt.Printf("Error getting id while inserting car: %v\n", err)

		}
	}
	return nil
}
func (c *Car) Get() error {
	query := `SELECT id,name,model,price,year FROM cars WHERE id=$1`
	if err := config.Db.Get(c, query, c.Id); err != nil {
		fmt.Printf("Error getting car: %v\n", err)
		return err
	}

	return nil
}
func (c *Car) Update() error {
	query := `UPDATE cars SET name = :name,model = :model,year = :year WHERE id = :id`
	_, err := config.Db.NamedQuery(query, c)
	if err != nil {
		fmt.Printf("Error updating car: %v\n", err)
		return err
	}
	return nil
}
func (c *Car) Delete() error {
	query := `DELETE FROM cars WHERE id=$1`
	if _, err := config.Db.Exec(query, c.Id); err != nil {
		fmt.Printf("Error deleting car: %v\n", err)
		return err
	}

	return nil
}
