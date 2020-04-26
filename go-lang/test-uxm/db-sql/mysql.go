
// GO Lang :: Sample MySQL / MariaDB

package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var schema = `
CREATE TABLE IF NOT EXISTS person (
	first_name text,
	last_name text,
	email text
);
`


type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

func main() {

	db, err := sqlx.Connect("mysql", "root:root@(127.0.0.1:3306)/smart_framework")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)
	db.MustExec("TRUNCATE TABLE person")

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "unixman@test.golang.loc")
	tx.Commit()

	people := []Person{}
	db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")

	log.Println("persons...")
	log.Println(people)

}

// END
