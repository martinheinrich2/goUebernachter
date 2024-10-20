package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// The Stay struct holds the data for a specific stay.
type Stay struct {
	ID              int
	StartDate       time.Time
	EndDate         time.Time
	TypeOfStay      string
	Room            string
	GuestId         int
	SocialWorkerId  int
	UserId          int
	Appointment     time.Time
	AppointmentDone bool
	StayProcessed   bool
}

// The StayJoinUser struct holds the data for the stays of a guest and the joined Names of the reception user
// and social worker
type StayJoinUser struct {
	ID              int
	StartDate       time.Time
	EndDate         time.Time
	TypeOfStay      string
	Room            string
	RLastName       string
	RFirstName      string
	SwLastName      string
	SwFirstName     string
	Appointment     time.Time
	AppointmentDone bool
	StayProcessed   bool
}

// The StayJoinGuest struct holds the data of a stay and the joined names of the guest and social worker
type StayJoinGuest struct {
	ID              int
	StartDate       time.Time
	EndDate         time.Time
	TypeOfStay      string
	Room            string
	GuestLastName   string
	GuestFirstName  string
	SwLastName      string
	SwFirstName     string
	Appointment     time.Time
	AppointmentDone bool
	StayProcessed   bool
}

type StayCount struct {
	Year  string
	Count int
}

type StayCount2 struct {
	Year       string
	Count      int
	TypeOfStay string
}

type StayModelInterface interface {
	Insert(startDate time.Time, endDate time.Time, typeOfStay string, room string, guestId int, socialWorkerId int,
		userId int, appointment time.Time) (int, error)
	Update(id int, startDate time.Time, endDate time.Time, typeOfStay string, room string, socialWorkerId int,
		appointment time.Time) (int, error)
	UpdateAppointmentDone(id int, appointmentDone bool) (int, error)
	UpdateStayProcessed(id int, appointmentDone bool) (int, error)
	Get(id int) (Stay, error)
	All(filters Filters) ([]StayJoinGuest, Metadata, error)
	AppointmentOpen(filters Filters) ([]StayJoinGuest, Metadata, error)
	StayNotProcessed(filters Filters) ([]StayJoinGuest, Metadata, error)
	GetGuestStays(guestId int) ([]StayJoinUser, error)
	Latest() ([]Stay, error)
	Statistics() ([]StayCount, error)
	Statistics2() ([]StayCount2, error)
}

// The StayModel struct which wraps a database connection pool.
type StayModel struct {
	DB *sql.DB
}

// The Insert method adds a new record to the "stay" table.
func (m *StayModel) Insert(startDate time.Time, endDate time.Time, typeOfStay string, room string, guestId int,
	socialWorkerId int, userId int, appointment time.Time) (int, error) {
	// Create SQL query
	stmt := `INSERT INTO stay (start_date, end_date, type_of_stay, room, guest_id, social_worker_id, user_id, appointment) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	// Use the Exec() method to insert the stay details into the stay table.
	result, err := m.DB.Exec(stmt, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"),
		typeOfStay, room, guestId, socialWorkerId, userId, appointment.Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// The Update method updates the data of a specific stay.
func (m *StayModel) Update(id int, startDate time.Time, endDate time.Time, typeOfStay string, room string,
	socialWorkerId int, appointment time.Time) (int, error) {
	stmt := `UPDATE stay SET start_date = ?, end_date = ?, type_of_stay = ?, room = ?,
                social_worker_id =?, appointment = ? WHERE id = ?`
	// Use the Exec() method to update the stay details.
	result, err := m.DB.Exec(stmt, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), typeOfStay, room, socialWorkerId,
		appointment.Format("2006-01-02 15:04:05"), id)
	if err != nil {
		return 0, err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateAppointmentDone sets the appointment_done bool value.
func (m *StayModel) UpdateAppointmentDone(id int, appointmentDone bool) (int, error) {
	stmt := `UPDATE stay SET appointment_done = ? WHERE id = ?`
	result, err := m.DB.Exec(stmt, appointmentDone, id)
	if err != nil {
		return 0, err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateStayProcessed sets the stay_processed bool value.
func (m *StayModel) UpdateStayProcessed(id int, stayProcessed bool) (int, error) {
	stmt := `UPDATE stay SET stay_processed = ? WHERE id = ?`
	result, err := m.DB.Exec(stmt, stayProcessed, id)
	if err != nil {
		return 0, err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Get returns a specific stay based on its id.
func (m *StayModel) Get(id int) (Stay, error) {
	stmt := `SELECT id, start_date, end_date, type_of_stay, room, guest_id, social_worker_id, user_id, appointment, 
       appointment_done, stay_processed FROM stay WHERE id = ?`
	var stay Stay
	// Use row.Scan() to copy the values from each field in sql.Row to the corresponding field
	// in the Stay struct. The arguments to row.Scan are *pointers* to the place we want
	// to copy the data into.
	err := m.DB.QueryRow(stmt, id).Scan(&stay.ID, &stay.StartDate, &stay.EndDate, &stay.TypeOfStay, &stay.Room,
		&stay.GuestId, &stay.SocialWorkerId, &stay.UserId, &stay.Appointment, &stay.AppointmentDone, &stay.StayProcessed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Stay{}, ErrNoRecord
		} else {
			return Stay{}, err
		}
	}
	return stay, nil
}

// GetGuestStays returns all stays and respective guest names.
func (m *StayModel) GetGuestStays(guestId int) ([]StayJoinUser, error) {
	//stmt := `SELECT id, start_date, end_date, type_of_stay, room, user_id, social_worker_id, appointment FROM stay WHERE guest_id = ? ORDER BY start_date desc `
	stmt := `SELECT stay.id, start_date, end_date, type_of_stay, stay.room, users.last_name, users.first_name, 
       users.last_name, users.first_name, appointment FROM stay INNER JOIN users on users.id = stay.social_worker_id WHERE guest_id = ? ORDER BY start_date desc `
	rows, err := m.DB.Query(stmt, guestId)
	if err != nil {
		return nil, err
	}
	// Defer rows.Close() to ensure the sql.Rows resultset is always closed properly before the method returns.
	defer rows.Close()
	// Initialize an empty slice go hold the stay structs.
	var stays []StayJoinUser
	// Iterate through the rows in the restultset.
	for rows.Next() {
		var stay StayJoinUser
		err = rows.Scan(&stay.ID, &stay.StartDate, &stay.EndDate, &stay.TypeOfStay, &stay.Room,
			&stay.RLastName, &stay.RFirstName, &stay.SwLastName, &stay.SwFirstName, &stay.Appointment)
		if err != nil {
			return nil, err
		}
		stays = append(stays, stay)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return stays, nil
}

// Latest returns the 10 most recent stays.
func (m *StayModel) Latest() ([]Stay, error) {
	// Create SQL statement to select the 10 latest stays.
	stmt := `SELECT id, start_date, end_date, type_of_stay, room, guest_id, social_worker_id, user_id, appointment FROM stay ORDER BY start_date DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the Stay structs.
	var stays []Stay
	for rows.Next() {
		// Create a pointer to a new zeroed Stay struct.
		var s Stay
		err = rows.Scan(&s.ID, &s.StartDate, &s.EndDate, &s.TypeOfStay, &s.Room, &s.GuestId, &s.SocialWorkerId, &s.UserId, &s.Appointment)
		if err != nil {
			return nil, err
		}
		stays = append(stays, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stays, nil
}

// All will return all stays.
func (m *StayModel) All(filters Filters) ([]StayJoinGuest, Metadata, error) {
	// Create SQL statement to select the 10 latest stays.
	stmt := fmt.Sprintf(
		`SELECT count(*) OVER(), stay.id, start_date, end_date, type_of_stay, appointment_done, stay_processed, stay.room,
	  guests.last_name, guests.first_name, appointment, users.last_name, users.first_name FROM stay
	      INNER JOIN guests ON guests.id = stay.guest_id INNER JOIN users ON users.id = stay.social_worker_id
	                    ORDER BY %s %s, start_date ASC LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	// Collect the values for the SQL statement placeholders in a slice.
	args := []any{filters.limit(), filters.offset()}

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	// Initialize an empty slice to hold the Stay structs.
	var stays []StayJoinGuest

	for rows.Next() {
		// Create a pointer to a new zeroed Stay struct.
		var stay StayJoinGuest

		err = rows.Scan(&totalRecords, &stay.ID, &stay.StartDate, &stay.EndDate, &stay.TypeOfStay, &stay.AppointmentDone,
			&stay.StayProcessed, &stay.Room, &stay.GuestLastName, &stay.GuestFirstName, &stay.Appointment,
			&stay.SwLastName, &stay.SwFirstName)
		if err != nil {
			return nil, Metadata{}, err
		}
		stays = append(stays, stay)
	}
	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	metadata.SortDirection = filters.sortDirection()
	return stays, metadata, nil
}

// AppointmentOpen will return a filtered list by appointment_done.
func (m *StayModel) AppointmentOpen(filters Filters) ([]StayJoinGuest, Metadata, error) {
	// Create SQL statement using SprintF to add variables.
	stmt := fmt.Sprintf(
		`SELECT count(*) OVER(), stay.id, start_date, end_date, type_of_stay, appointment_done, stay_processed, stay.room,
	  guests.last_name, guests.first_name, appointment, users.last_name, users.first_name FROM stay
	      INNER JOIN guests ON guests.id = stay.guest_id INNER JOIN users ON users.id = stay.social_worker_id
	                                                                                      WHERE appointment_done = false
	                    ORDER BY %s %s, start_date DESC LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())
	// Collect the values for the SQL statement placeholders in a slice.
	args := []any{filters.limit(), filters.offset()}

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice to hold the Stay structs.
	var stays []StayJoinGuest
	for rows.Next() {
		// Create a pointer to a new zeroed Stay struct.
		var stay StayJoinGuest
		err = rows.Scan(&totalRecords, &stay.ID, &stay.StartDate, &stay.EndDate, &stay.TypeOfStay, &stay.AppointmentDone,
			&stay.StayProcessed, &stay.Room, &stay.GuestLastName, &stay.GuestFirstName, &stay.Appointment,
			&stay.SwLastName, &stay.SwFirstName)
		if err != nil {
			return nil, Metadata{}, err
		}
		stays = append(stays, stay)
	}
	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	// Generate a Metadata struct, passing in the total record count and pagination parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return stays, metadata, nil

} // StayNotProcessed will return a filtered list by stay_processed
func (m *StayModel) StayNotProcessed(filters Filters) ([]StayJoinGuest, Metadata, error) {
	// Create SQL statement using SprintF to add variables.
	stmt := fmt.Sprintf(
		`SELECT count(*) OVER(), stay.id, start_date, end_date, type_of_stay, appointment_done, stay_processed, stay.room,
	  guests.last_name, guests.first_name, appointment, users.last_name, users.first_name FROM stay
	      INNER JOIN guests ON guests.id = stay.guest_id INNER JOIN users ON users.id = stay.social_worker_id
	                                                                                      WHERE stay_processed = false
	                    ORDER BY %s %s, start_date DESC LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())
	// Collect the values for the SQL statement placeholders in a slice.
	args := []any{filters.limit(), filters.offset()}

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice to hold the Stay structs.
	var stays []StayJoinGuest
	for rows.Next() {
		// Create a pointer to a new zeroed Stay struct.
		var stay StayJoinGuest
		err = rows.Scan(&totalRecords, &stay.ID, &stay.StartDate, &stay.EndDate, &stay.TypeOfStay, &stay.AppointmentDone,
			&stay.StayProcessed, &stay.Room, &stay.GuestLastName, &stay.GuestFirstName, &stay.Appointment,
			&stay.SwLastName, &stay.SwFirstName)
		if err != nil {
			return nil, Metadata{}, err
		}
		stays = append(stays, stay)
	}
	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	// Generate a Metadata struct, passing in the total record count and pagination parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return stays, metadata, nil
}

// Statistics summarizes/groups stays for a period
func (m *StayModel) Statistics() ([]StayCount, error) {
	// Create SQL statement to select stays fo a distinct year.
	//stmt := `SELECT count(*) FROM stay WHERE strftime('%Y', start_date)='2024'`
	stmt := `SELECT strftime('%Y', start_date) AS start_year,  COUNT(*) as year_count FROM stay 
	                    GROUP BY strftime('%Y', start_date);`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the StayCount structs.
	var stayscount []StayCount

	for rows.Next() {
		// Creat a pointer to a new zeroed Staycount struct.
		var staycount StayCount
		err := rows.Scan(&staycount.Year, &staycount.Count)
		if err != nil {
			return nil, err
		}
		stayscount = append(stayscount, staycount)
	}
	return stayscount, nil
}

// Statistics2 summarizes/groups stays for a period
func (m *StayModel) Statistics2() ([]StayCount2, error) {
	// Create SQL statement to select stays fo a distinct year.
	//stmt := `SELECT count(*) FROM stay WHERE strftime('%Y', start_date)='2024'`
	stmt := `SELECT strftime('%Y', start_date) AS start_year, type_of_stay,  COUNT(*) as year_count FROM stay 
	                                                             GROUP BY strftime('%Y', start_date), type_of_stay;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the StayCount structs.
	var stayscount2 []StayCount2

	for rows.Next() {
		// Creat a pointer to a new zeroed Staycount struct.
		var staycount2 StayCount2
		err := rows.Scan(&staycount2.Year, &staycount2.TypeOfStay, &staycount2.Count)
		if err != nil {
			return nil, err
		}

		stayscount2 = append(stayscount2, staycount2)
	}
	return stayscount2, nil
}
