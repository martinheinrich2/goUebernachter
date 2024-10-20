package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

func main() {
	var email string
	flag.StringVar(&email, "email", "", "email address")
	flag.Parse()
	if len(email) == 0 {
		fmt.Println("Please provide a valid email address")
		fmt.Println("Usage: createAdmin -email=\"email@domain.com\"")
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println("Querying database for: ", email)
	//  Create a connection pool from the openDB() function.
	//fmt.Println("Setting user 'Admin' to admin!")
	db, err := openDB()
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Check if email exists in database
	var exists bool
	exists, err = Exists(email)
	if err != nil {
		fmt.Printf(err.Error())
	}
	if !exists {
		fmt.Println("Email does not exist: ", email)
		os.Exit(1)
	}
	fmt.Println("Granting admin privileges to ", email)
	stmt := "UPDATE users SET admin = 1 WHERE email = ?"

	_, err = db.Exec(stmt, email)
	if err != nil {
		log.Printf(err.Error())
	}
	return
}

// The Exists method checks if a user with email exists
func Exists(email string) (bool, error) {
	var exists bool
	//  Create a connection pool from the openDB() function.
	//fmt.Println("Setting user 'Admin' to admin!")
	db, err := openDB()
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	stmt := "SELECT EXISTS (SELECT true FROM users WHERE email = ?)"
	err = db.QueryRow(stmt, email).Scan(&exists)
	return exists, err
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./reception.db")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		defer db.Close()
		return nil, err
	}
	return db, nil
}
