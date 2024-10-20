package models

import (
	"database/sql"
	"errors"
	"time"
)

// The Guest struct holds the data for a specific social worker.
type Guest struct {
	ID            int
	LastName      string
	FirstName     string
	BirthDate     time.Time
	BirthPlace    string
	IdNumber      string
	Nationality   string
	LastResidence string
	HouseBan      bool
	HbEndDate     time.Time
	HbStartDate   time.Time
}

// The GuestModelInterface describes the methods that the actual SocialWorkerModel struct has.
type GuestModelInterface interface {
	Insert(lastName string, firstName string, birthDate time.Time, birthPlace string,
		idNumber string, nationality string, lastResidence string, houseBan bool,
		hbEndDate time.Time, hbStartDate time.Time) (int, error)
	Update(id int, lastName string, firstName string, birthDate time.Time, birthPlace string,
		idNumber string, nationality string, lastResidence string, houseBan bool,
		hbEndDate time.Time, hbStartDate time.Time) (int, error)
	Get(id int) (Guest, error)
	All() ([]Guest, error)
	Search(lastName string) ([]Guest, error)
	Exists(id int) (bool, error)
	Delete(id int) error
}

// The GuestModel struct which wraps a database connection pool.
type GuestModel struct {
	DB *sql.DB
}

// The Insert method adds a new record to the "guests" table
func (m *GuestModel) Insert(lastName string, firstName string, birthDate time.Time, birthPlace string,
	idNumber string, nationality string, lastResidence string, houseBan bool, hbStartDate time.Time,
	hbEndDate time.Time) (int, error) {
	stmt := `INSERT INTO guests (last_name, first_name, birth_date, birth_place, id_number, nationality,
                    last_residence, house_ban, hb_start_date, hb_end_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	// Use the Exec() method to insert the guest details into the guests table.
	result, err := m.DB.Exec(stmt, lastName, firstName, birthDate, birthPlace, idNumber, nationality,
		lastResidence, houseBan, hbStartDate, hbEndDate)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// The Update method updates the data of a specific guest.
func (m *GuestModel) Update(id int, lastName string, firstName string, birthDate time.Time,
	birthPlace string, idNumber string, nationality string, lastResidence string, houseBan bool,
	hbStartDate time.Time, hbEndDate time.Time) (int, error) {
	stmt := `UPDATE guests SET last_name = ?, first_name = ?, birth_date = ?, birth_place = ?, id_number = ?,
                  nationality = ?, last_residence = ?, house_ban = ?, hb_start_date = ?,
                  hb_end_date = ?  WHERE id = ?`
	// Use the Exec() method to update the guest details.
	result, err := m.DB.Exec(stmt, lastName, firstName, birthDate, birthPlace, idNumber, nationality,
		lastResidence, houseBan, hbStartDate, hbEndDate, id)
	if err != nil {
		return 0, err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// The Get method returns the data of a specific guest.
func (m *GuestModel) Get(id int) (Guest, error) {
	var guest Guest

	stmt := `SELECT id, last_name, first_name, birth_date, birth_place, id_number, nationality, last_residence, 
       house_ban, hb_start_date, hb_end_date FROM guests WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&guest.ID, &guest.LastName, &guest.FirstName, &guest.BirthDate,
		&guest.BirthPlace, &guest.IdNumber, &guest.Nationality, &guest.LastResidence, &guest.HouseBan,
		&guest.HbStartDate, &guest.HbEndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Guest{}, ErrNoRecord
		} else {
			return Guest{}, err
		}
	}
	return guest, nil
}

// The All method returns all guests.
func (m *GuestModel) All() ([]Guest, error) {
	// Create SQL statement to select all guests
	stmt := `SELECT id, last_name, first_name, birth_date FROM guests order by last_name asc`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Defer rows.Close() to ensure the sql.Rows resultset is always properly closed before
	// the All() method returns.
	defer rows.Close()
	// Initialize an empty slice to hold the Guest structs.
	var guests []Guest
	// Use rows.Next() to iterate through the rows in the resultset. The resultset closes itself automatically
	// after completion and frees-up the underlying database connection.
	for rows.Next() {
		var guest Guest
		// Use row.Scan() to copy the values from each field in the row to the new
		// Guest object. The arguments to row.Scan() must be pointers to the place
		// we want to copy the data into.
		err = rows.Scan(&guest.ID, &guest.LastName, &guest.FirstName, &guest.BirthDate)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of guests.
		guests = append(guests, guest)
	}
	// Retrieve any error that was encountered during the iteration.
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Guests slice.
	return guests, nil
}

// The Search method returns all guests with specific letters in the name.
func (m *GuestModel) Search(searchName string) ([]Guest, error) {
	// Add wildcard for SQLite search.
	lastName := "%" + searchName + "%"
	// Create SQL statement to filter guest names
	stmt := `SELECT id, last_name, first_name, birth_date FROM guests 
                                             WHERE last_name LIKE ?`
	rows, err := m.DB.Query(stmt, lastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Guest{}, ErrNoRecord
		} else {
			return []Guest{}, err
		}
	}
	// Defer rows.Close() to ensure the sql.Rows result-set is always properly closed before the
	// Search() method returns.
	defer rows.Close()
	// Initialize an empty slice to hold the Guest structs.
	var guests []Guest
	// Use rows.Next() to iterate through the rows in the result-set. The result-set closes itself
	// automatically after completion and frees-up the underlying database connection.
	for rows.Next() {
		var guest Guest
		err = rows.Scan(&guest.ID, &guest.LastName, &guest.FirstName, &guest.BirthDate)
		if err != nil {
			return nil, err
		}
		guests = append(guests, guest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the filtered Guests slice.
	return guests, nil
}

// The Exists method checks if a social worker exists with a specific ID.
func (m *GuestModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM guests WHERE id = ?)`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

// The Delete method deletes a social worker with a specific ID.
func (m *GuestModel) Delete(id int) error {
	stmt := `DELETE FROM guests WHERE id = ?`
	_, err := m.DB.Exec(stmt, id)
	return err
}
