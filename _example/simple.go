package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kaboc/sqlp"
)

/*
CREATE TABLE user (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	name varchar(32) NOT NULL,
	age tinyint(3) unsigned NOT NULL,
	PRIMARY KEY (id)
) ENGINE=InnoDB;
*/

type tUser struct {
	Name string
	Age  int
}

func main() {
	// Change the settings here to match your environment
	db, err := sqlp.Open("mysql", "user:pw@tcp(host:3306)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Bulk insert
	res, err := db.Insert("user", []tUser{
		{Name: "User1", Age: 22},
		{Name: "User2", Age: 27},
		{Name: "User3", Age: 22},
	})
	if err != nil {
		log.Fatal(err)
	}

	cnt, _ := res.RowsAffected()
	id, _ := res.LastInsertId()
	fmt.Printf("Number of affected rows: %d\n", cnt)
	fmt.Printf("Last inserted ID: %d\n", id)

	// Update rows in transaction using prepared statement and placeholders
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("UPDATE user SET name = ? WHERE name = ?")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec("User1-2", "User1")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec("User2-2", "User2")
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// Select rows using named placeholders and put them into slice of structures
	var u []tUser
	q := "SELECT name, age FROM user WHERE age IN :age[2] AND name LIKE :name"
	err = db.SelectToStruct(&u, q, map[string]interface{}{
		"age":  []interface{}{22, 27},
		"name": "%-2",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range u {
		fmt.Printf("%s (%d yo)\n", v.Name, v.Age)
	}
}
