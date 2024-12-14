package models

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"modernc.org/sqlite"
	"strings"
	"time"
)

// The User struct holds the data for a specific user.
type User struct {
	ID             int
	LastName       string
	FirstName      string
	Email          string
	JobTitle       string
	Room           string
	HashedPassword []byte
	Created        time.Time
	Admin          bool
	Active         bool
}

// The UserModelInterface describes the methods that the actual UserModel struct has.
type UserModelInterface interface {
	Insert(lastname, firstname, email, job, room, password string) error
	Update(id int, lastname, firstname, email, job, room string, admin bool, active bool) (int, error)
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	MailExists(email string) (bool, error)
	Get(id int) (User, error)
	GetUserId(email string) (int, error)
	SelectUserByJob(job string) ([]User, error)
	All() ([]User, error)
	AllNames() ([]User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
	ResetPassword(id int, newPassword string) error
}

// The UserModel struct which wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// The Insert method adds a new record to the "users" table.
func (m *UserModel) Insert(lastname, firstname, email, job, room, password string) error {
	fmt.Println("models/user.go")
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (last_name, first_name, email, job_title, room, hashed_password, created) VALUES(?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	// Use the Exec() method to insert the user details and hashed password
	// into the users table.
	fmt.Println("executing db stmt")
	_, err = m.DB.Exec(stmt, lastname, firstname, email, job, room, string(hashedPassword))
	fmt.Println("checking db error")
	if err != nil {
		fmt.Println("err != nil: ", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicateEmail
		}
		// probably no longer need this part

		// If this returns an error, we use the errors.As() function to check
		// whether the error has the type *sqlite.Error. If it does, the
		// error will be assigned to the sqliteError variable. We can then check
		// whether or not the error relates to our users_uc_email key by
		// checking if the error code equals 144 and the contents of the error
		// message string. If it does, we return an ErrDuplicateEmail error.
		var sqliteError *sqlite.Error
		fmt.Println("sqliteError is: ", sqliteError.Code())
		if errors.As(err, &sqliteError) {
			fmt.Println("the errors.AS is: ", sqliteError.Code())
			if strings.Contains(sqliteError.Error(), "UNIQUE constraint failed:") {
				return ErrDuplicateEmail
			}
			fmt.Println("users.go insert should return ErrDuplicateEmail ", sqliteError.Error())
		}
		return err
	}

	return nil
}

// The Update method updates the data of a specific user.
func (m *UserModel) Update(id int, lastName string, firstname string, email string, job string, room string,
	admin bool, active bool) (int, error) {
	stmt := `UPDATE users SET last_name = ?, first_name = ?, email = ?, job_title = ?, room = ?, admin = ?, active = ? WHERE id = ?`
	// Use the Exec() method to update user details.
	fmt.Println(stmt, lastName, firstname, email, job, room, admin, active, id)
	result, err := m.DB.Exec(stmt, lastName, firstname, email, job, room, admin, active, id)
	if err != nil {
		return 0, err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// The Authenticate method verifies whether a user exists with the provided
// email address and password. This will return the relevant user ID if
// they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Retrieve the id and hashed password associated with the given email. If
	// no matching email exists we return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

// The Exists method checks if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

// The MailExists method checks if a mail address exists in the database.
func (m *UserModel) MailExists(email string) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE email = ?)`
	err := m.DB.QueryRow(stmt, email).Scan(&exists)
	return exists, err
}

// The GetUserId method gets the user id for a mail address.
func (m *UserModel) GetUserId(email string) (int, error) {
	var id int
	stmt := `SELECT id FROM users WHERE email = ?`
	err := m.DB.QueryRow(stmt, email).Scan(&id)
	return id, err
}

// The All method returns all staff from the database.
func (m *UserModel) All() ([]User, error) {
	// Create SQL statement to select all articles
	stmt := `SELECT id, last_name, first_name, job_title, room FROM users order by last_name`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Defer rows.Close() to ensure the sql.Rows resultset is always properly closed before
	// the All() method returns.
	defer rows.Close()
	// Initialize an empty slice to hold the Guest structs.
	var users []User
	// Use rows.Next() to iterate through the rows in the resultset. The resultset closes itself automatically
	// after completion and frees-up the underlying database connection.
	for rows.Next() {
		var user User
		// Use row.Scan() to copy the values from each field in the row to the new
		// Guest object. The arguments to row.Scan() must be pointers to the place
		// we want to copy the data into.
		err = rows.Scan(&user.ID, &user.LastName, &user.FirstName, &user.JobTitle, &user.Room)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of guests.
		users = append(users, user)
	}
	// Retrieve any error that was encountered during the iteration.
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Guests slice.
	return users, nil
}

// The AllNames method returns all users ordered by lastname, firstname
func (m *UserModel) AllNames() ([]User, error) {
	stmt := `SELECT id, last_name, first_name FROM users ORDER BY last_name, first_name`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.LastName, &user.FirstName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// The Get method returns the data of a specific user.
func (m *UserModel) Get(id int) (User, error) {
	var user User
	stmt := `SELECT id, last_name, first_name, email, job_title, room, created, admin, active FROM users WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.LastName, &user.FirstName, &user.Email, &user.JobTitle,
		&user.Room, &user.Created, &user.Admin, &user.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		} else {
			return User{}, err
		}
	}
	return user, nil
}

// The SelectUserByJob method returns only users with a specific job title.
func (m *UserModel) SelectUserByJob(jobTitle string) ([]User, error) {
	stmt := `SELECT id, last_name, first_name, room FROM users WHERE job_title = ? ORDER BY last_name`
	rows, err := m.DB.Query(stmt, jobTitle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.LastName, &user.FirstName, &user.Room)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// The PasswordUpdate method updates the password in the database.
func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	var currentHashedPassword []byte

	stmt := "SELECT hashed_password FROM users WHERE id = ?"

	err := m.DB.QueryRow(stmt, id).Scan(&currentHashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(currentHashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	stmt = "UPDATE users SET hashed_password = ? WHERE id = ?"

	_, err = m.DB.Exec(stmt, string(newHashedPassword), id)
	return err
}

// The ResetPassword method resets the password in the database.
func (m *UserModel) ResetPassword(id int, newPassword string) error {

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	stmt := "UPDATE users SET hashed_password = ? WHERE id = ?"
	_, err = m.DB.Exec(stmt, string(newHashedPassword), id)
	return err
}
