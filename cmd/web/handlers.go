package main

import (
	"errors"
	"fmt"
	"github.com/martinheinrich2/goUebernachter/internal/models"
	"github.com/martinheinrich2/goUebernachter/internal/validator"
	"github.com/wneessen/go-mail"
	"net/http"
	"os"
	"strconv"
	"time"
)

// The home handler is called by the / route and displays the home page.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data, and add the stays slice to it.
	data := app.newTemplateData(r)
	// Pass the data to the render() helper.
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

// The stats handler is called by the GET /stats route and displays statistical data.
func (app *application) stats(w http.ResponseWriter, r *http.Request) {
	statistics, err := app.stays.Statistics()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	statistics2, err := app.stays.Statistics2()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.StayCount = statistics
	data.StayCount2 = statistics2
	app.render(w, r, http.StatusOK, "statistics.tmpl", data)
}

// The accountInactive handler is called from the Get /inactive route and displays
// warning page if the account is deactivated.
func (app *application) accountInactive(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "account_inactive.tmpl", app.newTemplateData(r))
}

// The guestCreate handler is called from the GET /guest/create route and displays a
// guest create form..
func (app *application) guestCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// Initialize a new createGuestForm instance and pass it to the template
	data.Form = guestCreateForm{}
	app.render(w, r, http.StatusOK, "guest_create.tmpl", data)
}

// The guestCreateForm struct holds the struct tags which tell the decoder how to map HTML form value into the
// different struct fields. The struct tag `form:":"` tells the decoder to completely ignore a field during decoding.
type guestCreateForm struct {
	LastName            string    `form:"lastname"`
	FirstName           string    `form:"firstname"`
	BirthDate           time.Time `form:"birthdate"`
	BirthPlace          string    `form:"birthplace"`
	IdNumber            string    `form:"idnumber"`
	Nationality         string    `form:"nationality"`
	LastResidence       string    `form:"lastresidence"`
	HouseBan            bool      `form:"houseban"`
	HBStartDate         time.Time `form:"hbstartdate"`
	HBEndDate           time.Time `form:"hbenddate"`
	validator.Validator `form:"-"`
}

// The guestUpdateForm struct holds the struct tags which tell the decoder how to map HTML vorm values
// into the different struct fields. The struct tag `form:":"` tells the decoder to completely ignore
// a field during decoding.
type guestUpdateForm struct {
	Id                  int       `form:"id"`
	LastName            string    `form:"lastname"`
	FirstName           string    `form:"firstname"`
	BirthDate           time.Time `form:"birthdate"`
	BirthPlace          string    `form:"birthplace"`
	IdNumber            string    `form:"idnumber"`
	Nationality         string    `form:"nationality"`
	LastResidence       string    `form:"lastresidence"`
	HouseBan            bool      `form:"houseban"`
	HBStartDate         time.Time `form:"hbstartdate"`
	HBEndDate           time.Time `form:"hbenddate"`
	validator.Validator `form:"-"`
}

// The guestCreatePost handler is called by the POST /guest/create route. It checks the validity of
// the input and inserts a new guest into the database.
func (app *application) guestCreatePost(w http.ResponseWriter, r *http.Request) {
	// Declare a new empty instance of the guestCreateForm
	var form guestCreateForm
	// Use the decodePostForm helper to handle the form data.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// The Validator struct is embedded in the guestCreateForm struct.
	// Call CheckField() directly to execute the validation checks.
	form.CheckField(validator.NotBlank(form.LastName), "lastname", "Eingabe erforderlich")
	form.CheckField(validator.NotBlank(form.FirstName), "firstname", "Eingabe erforderlich")
	form.CheckField(validator.NotBlank(form.BirthPlace), "birthplace", "Eingabe erforderlich")
	form.CheckField(validator.NotBlank(form.IdNumber), "idNumber", "Eingabe erforderlich")
	form.CheckField(validator.NotBlank(form.Nationality), "nationality", "Eingabe erforderlich")
	form.CheckField(validator.NotBlank(form.LastResidence), "lastResidence", "Eingabe erforderlich")
	// Use the Valid() method to see if any of the checks failed. Re-render the template
	// passing in the form, if check failed.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "guest_create.tmpl", data)
		return
	}

	_, err = app.guests.Insert(form.LastName, form.FirstName, form.BirthDate, form.BirthPlace, form.IdNumber,
		form.Nationality, form.LastResidence, form.HouseBan, form.HBStartDate, form.HBEndDate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Use the Put() method to add a string value ("Guest successfully
	// created!") and the corresponding key ("flash") to the session data
	app.sessionManager.Put(r.Context(), "flash", "Klient erfolgreich angelegt!")

	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}

// The guestUpdate handler is called from the GET /guest/update route and displays the guest_update page.
func (app *application) guestUpdate(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the GuestModel's Get() method to retrieve the data for a specific guest based on its ID.
	// Return a 404 Not Found response if there is no matching record.
	guest, err := app.guests.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Form = guest
	// Use the render helper.
	app.render(w, r, http.StatusOK, "guest_update.tmpl", data)
}

// The guestUpdatePost handler is called by the Post /guest/update/{id} route.
func (app *application) guestUpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}

	var form guestUpdateForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	_, err = app.guests.Update(id, form.LastName, form.FirstName, form.BirthDate, form.BirthPlace, form.IdNumber,
		form.Nationality, form.LastResidence, form.HouseBan, form.HBStartDate, form.HBEndDate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Klient erfolgreich geändert!")

	http.Redirect(w, r, fmt.Sprintf("/guests"), http.StatusSeeOther)
}

// The guestSearch handler is called by the GET /guest/search route.
func (app *application) guestSearch(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "search.tmpl", data)
}

// Create a guestSearchForm struct
type guestSearchForm struct {
	Search              string `form:"search"`
	validator.Validator `form:"-"`
}

// The guestSearchPost handler is called by the POST /guest/search route.
func (app *application) guestSearchPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("guestSearchPost searching ...")
	var form guestSearchForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	lastName := form.Search

	guests, err := app.guests.Search(lastName)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Guests = guests
	app.render(w, r, http.StatusOK, "search_results.tmpl", data)
}

// The allGuestsView handler is called by the GET /guests route.
func (app *application) allGuestsView(w http.ResponseWriter, r *http.Request) {
	// Get all guests
	guests, err := app.guests.All()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Guests = guests
	app.render(w, r, http.StatusOK, "guests_view.tmpl", data)
}

// The guestView handler is called by the GET /guest/view route.
func (app *application) guestView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the GuestModel's Get() method to retrieve the data for a specific guest based
	// on its ID. Return a 404 Not Found response if there is no matching record.
	guest, err := app.guests.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	stays, err := app.stays.GetGuestStays(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Guest = guest
	data.StayJoinUser = stays

	// Use the render helper.
	app.render(w, r, http.StatusOK, "guest_view.tmpl", data)
}

// Create a new stayCreateForm struct.
type stayCreateForm struct {
	ID                  int
	LastName            string
	FirstName           string
	StartDate           time.Time `form:"start_date"`
	EndDate             time.Time `form:"end_date"`
	TypeOfStay          string    `form:"type_of_stay"`
	Room                string    `form:"room"`
	GuestId             int
	SocialWorkerId      int       `form:"social_worker_id"`
	Appointment         time.Time `form:"appointment"`
	UserId              int
	SocialWorkers       []models.User
	AppointmentDone     bool `form:"appointment_done"`
	StayProcessed       bool `form:"stay_processed"`
	validator.Validator `form:"-"`
}

// Create a new stayUpdateForm struct.
type stayUpdateForm struct {
	ID                  int
	LastName            string
	FirstName           string
	StartDate           time.Time `form:"start_date"`
	EndDate             time.Time `form:"end_date"`
	TypeOfStay          string    `form:"type_of_stay"`
	Room                string    `form:"room"`
	SocialWorkerId      int       `form:"social_worker_id"`
	Appointment         time.Time `form:"appointment"`
	AppointmentDone     bool      `form:"appointment_done"`
	StayProcessed       bool      `form:"stay_processed"`
	SocialWorkers       []models.User
	validator.Validator `form:"-"`
}

// The stayCreate handler is called by the GET /stay/create route.
func (app *application) stayCreate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	guest, err := app.guests.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	socialWorkers, err := app.users.SelectUserByJob("Sozialarbeit")
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	lastName := guest.LastName
	firstName := guest.FirstName

	data := app.newTemplateData(r)
	data.Form = stayCreateForm{ID: id, LastName: lastName, FirstName: firstName, UserId: userID,
		SocialWorkers: socialWorkers}
	app.render(w, r, http.StatusOK, "stay_create.tmpl", data)
}

// The stayCreatePost handler ist called by the POST /stay/create route.
func (app *application) stayCreatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	var form stayCreateForm
	// Parse the form data into the stayCreateForm struct.
	err = app.decodePostForm(r, &form)
	if err != nil {
		fmt.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using the helper function.
	form.CheckField(validator.NotBlank(form.Room), "room", "This field cannot be blank")

	// If there are any errors, redisplay the stay create form with a 422 status code.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "stay_create.tmpl", data)
		return
	}

	fmt.Println(form.StartDate, form.EndDate, form.TypeOfStay, form.Room, form.GuestId,
		form.SocialWorkerId, userID, form.Appointment)
	// Try to create a stay record in the database.
	_, err = app.stays.Insert(form.StartDate, form.EndDate, form.TypeOfStay, form.Room, id,
		form.SocialWorkerId, userID, form.Appointment)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Aufenthalt erfolgreich angelegt!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// The stayUpdate handler is called by the GET /stay/update/{id} route.
func (app *application) stayUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	stay, err := app.stays.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	guest, err := app.guests.Get(stay.GuestId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	socialWorkers, err := app.users.SelectUserByJob("Sozialarbeit")
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		fmt.Println(socialWorkers)
		return

	}
	socialWorker, err := app.users.Get(stay.SocialWorkerId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	user, err := app.users.Get(stay.UserId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	authUserID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	authUser, err := app.users.Get(authUserID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Form = stayUpdateForm{ID: stay.ID, StartDate: stay.StartDate, EndDate: stay.EndDate, TypeOfStay: stay.TypeOfStay, Room: stay.Room,
		SocialWorkerId: stay.SocialWorkerId, Appointment: stay.Appointment, AppointmentDone: stay.AppointmentDone,
		StayProcessed: stay.StayProcessed, SocialWorkers: socialWorkers}
	//data.Form = stay
	data.Guest = guest
	data.SocialWorker = socialWorker
	data.User = user
	data.AuthUser = authUser
	app.render(w, r, http.StatusOK, "stay_update.tmpl", data)
}

// The stayUpdatePost handler is called byt the POST /stay/update route.
func (app *application) stayUpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	var form stayCreateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	_, err = app.stays.Update(id, form.StartDate, form.EndDate, form.TypeOfStay, form.Room, form.SocialWorkerId,
		form.Appointment)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Aufenthalt erfolgreich geändert!")
	http.Redirect(w, r, "/stays", http.StatusSeeOther)
}

// The allStaysView handler is called by the GET /stays route.
func (app *application) allStaysView(w http.ResponseWriter, r *http.Request) {
	// Define an input struct to hold the expected values from the request.
	var input struct {
		AppointmentDone int
		models.Filters
	}
	// Initialize a new Validator instance
	v := validator.New()
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	// Use the helper to extract the page and page size query string values as integers.
	// Set default values to page one and page size to 10.
	input.AppointmentDone = app.readInt(qs, "appointment_done", 1, v)
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 12, v)
	input.Filters.Sort = app.readString(qs, "sort", "start_date")
	input.Filters.SortSafeList = []string{"guests.last_name", "guests.first_name", "start_date", "end_date",
		"type_of_stay", "appointment_done", "stay_processed", "-appointment_done", "-stay_processed",
		"-guests.last_name", "-guests.first_name", "-start_date", "-end_date", "-type_of_stay", "stay.room",
		"-stay.room", "users.last_name", "-users.last_name"}
	input.Filters.SortDirection = app.readString(qs, "sort", "start_date")
	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send a response if necessary.
	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		//app.clientError(w, http.StatusBadRequest)
		app.validationError(w)
		fmt.Println(models.FailedValidation)
		return
	}
	// Get all guests
	stays, metadata, err := app.stays.All(input.Filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.StayJoinGuest = stays
	data.Metadata = metadata
	data.Metadata.SortDirection = input.Filters.SortDirection
	app.render(w, r, http.StatusOK, "stays_view.tmpl", data)
}

// The appointmentOpenView handler ist called by the GET /stays/appointment route
func (app *application) appointmentOpenView(w http.ResponseWriter, r *http.Request) {
	// Define an input struct to hold the expected values from the request.
	var input struct {
		AppointmentDone int
		models.Filters
	}
	// Initialize a new Validator instance
	v := validator.New()
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	// Use the helper to extract the page and page size query string values as integers.
	// Set default values to page one and page size to 10.
	input.AppointmentDone = app.readInt(qs, "appointment_done", 1, v)
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 10, v)
	input.Filters.Sort = app.readString(qs, "sort", "start_date")
	input.Filters.SortDirection = app.readString(qs, "sort_direction", "desc")
	input.Filters.SortSafeList = []string{"guests.last_name", "guests.first_name", "start_date", "end_date",
		"type_of_stay", "appointment_done", "stay_processed", "-appointment_done", "-stay_processed",
		"-guests.last_name", "-guests.first_name", "-start_date", "-end_date", "-type_of_stay", "stay.room",
		"-stay.room", "users.last_name", "-users.last_name"}
	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send a response if necessary.
	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		//app.clientError(w, http.StatusBadRequest)
		app.validationError(w)
		fmt.Println(models.FailedValidation)
		return
	}

	// Get all guests
	stays, metadata, err := app.stays.AppointmentOpen(input.Filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.StayJoinGuest = stays
	data.Metadata = metadata
	app.render(w, r, http.StatusOK, "stays_view.tmpl", data)
}

// The stayNotProcessedView handler ist called by the GET /stays/appointment route
func (app *application) stayNotProcessedView(w http.ResponseWriter, r *http.Request) {
	// Define an input struct to hold the expected values from the request.
	var input struct {
		StayProcessed int
		models.Filters
	}
	// Initialize a new Validator instance
	v := validator.New()
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	// Use the helper to extract the page and page size query string values as integers.
	// Set default values to page one and page size to 10.
	input.StayProcessed = app.readInt(qs, "appointment_done", 0, v)
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 10, v)
	input.Filters.Sort = app.readString(qs, "sort", "start_date")
	input.Filters.SortDirection = app.readString(qs, "sort_direction", "desc")
	input.Filters.SortSafeList = []string{"guests.last_name", "guests.first_name", "start_date", "end_date",
		"type_of_stay", "appointment_done", "stay_processed", "-appointment_done", "-stay_processed",
		"-guests.last_name", "-guests.first_name", "-start_date", "-end_date", "-type_of_stay", "stay.room",
		"-stay.room", "users.last_name", "-users.last_name"}
	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send a response if necessary.
	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		//app.clientError(w, http.StatusBadRequest)
		app.validationError(w)
		fmt.Println(models.FailedValidation)
		return
	}
	// Get all guests
	stays, metadata, err := app.stays.StayNotProcessed(input.Filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.StayJoinGuest = stays
	data.Metadata = metadata
	app.render(w, r, http.StatusOK, "stays_view.tmpl", data)
}

// The stayDetail handler is called by the GET /stay/detail route.
func (app *application) stayDetail(w http.ResponseWriter, r *http.Request) {
	authUserID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the StayModel's Get() method to retrieve the data for a specific stay based on its ID.
	// Return a 404 Not Found response if there is no matching record.
	stay, err := app.stays.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	guest, err := app.guests.Get(stay.GuestId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	socialWorker, err := app.users.Get(stay.SocialWorkerId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	user, err := app.users.Get(stay.UserId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	authUser, err := app.users.Get(authUserID)
	now := time.Now().Format("02.01.2006")
	data := app.newTemplateData(r)
	data.Stay = stay
	data.Guest = guest
	data.SocialWorker = socialWorker
	data.User = user
	data.AuthUser = authUser
	data.Date = now
	// Use the render helper.
	app.render(w, r, http.StatusOK, "stay_detail.tmpl", data)
}

// The postAppointmentDone handler toggles the appointment_done field
func (app *application) postAppointmentDone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
	}
	stay, err := app.stays.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	if stay.AppointmentDone == true {
		stay.AppointmentDone = false
	} else {
		stay.AppointmentDone = true
	}
	_, err = app.stays.UpdateAppointmentDone(id, stay.AppointmentDone)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	url := "/stay/detail/" + strconv.Itoa(id)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// The postDataProcessed handler toggles the stay_processed field
func (app *application) postDataProcessed(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
	}
	stay, err := app.stays.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	if stay.StayProcessed == true {
		stay.StayProcessed = false
	} else {
		stay.StayProcessed = true
	}
	_, err = app.stays.UpdateStayProcessed(id, stay.StayProcessed)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	url := "/stay/detail/" + strconv.Itoa(id)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

type stayDetailPrintForm struct {
	StayForm         bool `form:"stayform"`
	DataForm         bool `form:"dataform"`
	AppointmentForm  bool `form:"appointmentform"`
	ConfirmationForm bool `form:"confirmationform"`
	ClearingForm     bool `form:"clearingform"`
}

// The stayDetailPrint handler is called by the GET /stay/detail_print route.
func (app *application) stayDetailPrint(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Declare a new empty instance of the stayDetailPrintForm
	var form stayDetailPrintForm
	// Use the decodePostForm helper to handle the form data.
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the StayModel's Get() method to retrieve the data for a specific stay based
	// on its ID. Return a 404 Not Found response if there is no matching record.
	stay, err := app.stays.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	guest, err := app.guests.Get(stay.GuestId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	socialWorker, err := app.users.Get(stay.SocialWorkerId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	user, err := app.users.Get(stay.UserId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	now := time.Now().Format("02.01.2006")
	data := app.newTemplateData(r)
	data.Stay = stay
	data.Guest = guest
	data.SocialWorker = socialWorker
	data.User = user
	data.Date = now
	data.Form = form

	// Use the render helper.
	app.render(w, r, http.StatusOK, "stay_detail_print.tmpl", data)
}

// Create a new userSignupForm struct.
type userSignupForm struct {
	LastName             string `form:"lastname"`
	FirstName            string `form:"firstname"`
	Email                string `form:"email"`
	JobTitle             string `form:"job_title"`
	Room                 string `form:"room"`
	Password             string `form:"password"`
	PasswordConfirmation string `form:"passwordConfirmation"`
	validator.Validator  `form:"-"`
}

// Create a new userUpdateForm struct.
type userUpdateForm struct {
	LastName            string `form:"lastname"`
	FirstName           string `form:"firstname"`
	Email               string `form:"email"`
	JobTitle            string `form:"job_title"`
	Room                string `form:"room"`
	Admin               bool   `form:"admin"`
	Active              bool   `form:"active"`
	validator.Validator `form:"-"`
}

// The userSignup handler displays the signup page.
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

// The userSignupPost handler processes the signup data.
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare an zero-valued instance of our userSignupForm struct.
	var form userSignupForm

	// Parse the form data into the userSignupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.LastName), "lastname", "Eingabe erforderlich")
	form.CheckField(validator.NotBlank(form.FirstName), "firstname", "Eingabe erforderlich")
	form.CheckField(validator.NotBlank(form.Email), "email", "Eingabe erforderlich")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Keine gültige Email Adresse")
	form.CheckField(validator.NotBlank(form.JobTitle), "job_title", "Eingabe erforderlich")
	form.CheckField(validator.NotBlank(form.Room), "room", "Eingabe erforderlich")
	form.CheckField(validator.NotBlank(form.Password), "password", "Eingabe erforderlich")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Mindestens 8 Zeichen eingeben")
	form.CheckField(validator.NotBlank(form.PasswordConfirmation), "passwordConfirmation", "Eingabe erforderlich")
	form.CheckField(form.Password == form.PasswordConfirmation, "passwordConfirmation", "Passwörter stimmen nicht überein")

	// If there are any errors, redisplay the signup form along with a 422 status code.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Insert(form.LastName, form.FirstName, form.Email, form.JobTitle, form.Room, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Adresse ist schon in Verwendung")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// The allUsersView handler is called by the GET /users route.
func (app *application) allUsersView(w http.ResponseWriter, r *http.Request) {
	// Get all guests
	users, err := app.users.All()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Users = users
	app.render(w, r, http.StatusOK, "users_view.tmpl", data)
}

// Create a new userLoginForm struct.
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// The userLogin handler is called by the GET /user/login route.
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

// The userLoginPost handler is called by the POST /user/login route.
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// Decode the form data into the userLoginForm struct.
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Do some validation checks on the form. We check that both email and
	// password are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Email), "email", "Email Adresse fehlt")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Gültige Email Adresse eingeben")
	form.CheckField(validator.NotBlank(form.Password), "password", "Passwort fehlt")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	if !app.isActive(id) {
		http.Redirect(w, r, "/inactive", http.StatusSeeOther)
	}

	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Insert the ID of the current user to the session, so that they are now
	// 'logged in'.
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Use the PopString method to retrieve and remove a value from the session
	// data in one step. If no matching key exists this will return the empty
	// string.
	path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}

	// Redirect the user to the create articles page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// The userLogoutPost handler is called by the POST /user/logout route.
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Use the RenewToken() method on the current session to change the session
	// ID again.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// Insert a flash message to the session to confirm to the user that they've been
	// logged out.
	app.sessionManager.Put(r.Context(), "flash", "Erfolgreich abgemeldet")

	// Redirect the user to the application home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// The ping handler displays the status of the server
func ping(w http.ResponseWriter, r *http.Request) {
	str := "Server is running OK!"
	w.Write([]byte(str))
}

// The accountView handler displays the user account information.
func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	user, err := app.users.Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.User = user

	app.render(w, r, http.StatusOK, "account.tmpl", data)
}

// The userUpdate handler is called by the GET /user/update/{id} route.
func (app *application) userUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the GuestModel's Get() method to retrieve the data for a specific guest
	// based on its ID. Return a 404 Not Found response if there is no matching
	// record.
	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Form = user
	// Use the render helper.
	app.render(w, r, http.StatusOK, "user_update.tmpl", data)
}

// The userUpdatePost handler is called by the Post /user/update/{id} route.
func (app *application) userUpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.NotFound(w, r)
		return
	}
	var form userUpdateForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	_, err = app.users.Update(id, form.LastName, form.FirstName, form.Email, form.JobTitle, form.Room, form.Admin,
		form.Active)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Nutzer erfolgreich geändert!")

	http.Redirect(w, r, fmt.Sprintf("/users"), http.StatusSeeOther)
}

// Create a new accountPasswordUpdateForm struct.
type accountPasswordUpdateForm struct {
	CurrentPassword         string `form:"currentPassword"`
	NewPassword             string `form:"newPassword"`
	NewPasswordConfirmation string `form:"newPasswordConfirmation"`
	validator.Validator     `form:"-"`
}

// The accountPasswordUpdate handler creates the Password Update form template.
func (app *application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = accountPasswordUpdateForm{}

	app.render(w, r, http.StatusOK, "password.tmpl", data)
}

// The accountPasswordUpdatePost handler decodes the form, checks the input data and updates the password.
func (app *application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form accountPasswordUpdateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "Eingabe fehlt!")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "Eingabe fehlt!")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "Mindestens 8 Zeichen eingeben")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "Eingabe fehlt!")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwörter stimmen nicht überein")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "password.tmpl", data)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.users.PasswordUpdate(userID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Current password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form

			app.render(w, r, http.StatusUnprocessableEntity, "password.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your password has been updated!")

	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

// Create a new userPasswordResetForm struct.
type userPasswordResetForm struct {
	Id                      int
	NewPassword             string `form:"newPassword"`
	NewPasswordConfirmation string `form:"newPasswordConfirmation"`
	validator.Validator     `form:"-"`
}

// The userPasswordReset handler is called by the GET /user/password/reset/{id} route
func (app *application) userPasswordReset(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	data := app.newTemplateData(r)
	data.Form = userPasswordResetForm{Id: id}
	app.render(w, r, http.StatusOK, "password_reset.tmpl", data)
}

// The userPasswordResetPost handler is called by the POST /user/password/reset/{id} route.
func (app *application) userPasswordResetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}
	var form userPasswordResetForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "password_reset.tmpl", data)
		return
	}
	// Store the new password hash in the database.
	err = app.users.ResetPassword(id, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Current password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "password.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "The password has been reset!")

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// Create a new forgotPasswordForm struct.
type forgotPasswordForm struct {
	Email               string `form:"email"`
	validator.Validator `form:"-"`
}

// The forgotPassword handler is called by the GET /account/forgot_password route.
// Displays the enter email form und is the entry point for the forgot password function.
func (app *application) forgotPassword(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = forgotPasswordForm{}
	// Pass the data to the render() helper.
	app.render(w, r, http.StatusOK, "forgot_password.tmpl", data)
}

// The forgotPasswordPost handler is called by the POST /account/forgot_password route.
// It checks if a valid email address has been provided and if the email is in the users table.
// Then it stores the userid, random code and timestamp in the password_reset_tokens table.
// Initiates the sendToken function and redirects to /enter_token/user@somemail.com
func (app *application) forgotPasswordPost(w http.ResponseWriter, r *http.Request) {
	// get email from form data
	var form forgotPasswordForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Do some validation checks on the form.
	form.CheckField(validator.NotBlank(form.Email), "email", "Email Adresse fehlt")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Gültige Email Adresse eingeben")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "forgot_password.tmpl", data)
		return
	}
	// check if email is in database, throw error if not
	var exists bool
	exists, err = app.users.MailExists(form.Email)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
		}
	}
	// If email exists, store token in database.
	if exists {
		var result string
		result = app.storeToken(form.Email)
		if result == "" {
			app.logger.Error("storing token failed" + err.Error())
			return
		}
		// Send token tu user email address in separate goroutine.
		go app.sendToken(form.Email, result)
		app.logger.Info("token sent to user" + form.Email)

		// Redirect to the enter token page
		url := "/enter_token/" + form.Email
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// Create a new enterTokenForm struct to hold the token data.
type enterTokenForm struct {
	Email               string `form:"email"`
	Token               string `form:"token"`
	validator.Validator `form:"-"`
}

// The enterToken handler is called by the GET /enter/token/{email} route.
// Displays the enter token form.
func (app *application) enterToken(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	data := app.newTemplateData(r)
	data.Form = enterTokenForm{Email: email}
	app.render(w, r, http.StatusOK, "enter_token.tmpl", data)
}

// The enterTokenPost handler is called by the POST /enter/token/{email} route.
// It checks if a token is provided, gets the userId from the email address and stores the userId in the
// sessionManager context. It checks if token for userId is in database. If yes, it redirects
// to /password_reset.
func (app *application) enterTokenPost(w http.ResponseWriter, r *http.Request) {
	// Get Email from url
	email := r.PathValue("email")
	// Decode the form data into the enterTokenForm struct.
	var form enterTokenForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	// Do some validation checks on the form.
	form.CheckField(validator.NotBlank(form.Token), "token", "Code fehlt!")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "enter_token.tmpl", data)
		return
	}
	// Get userID from database
	userID, err := app.users.GetUserId(email)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	// Store userid in sessionManager
	id := strconv.Itoa(userID)
	app.sessionManager.Put(r.Context(), "userID", id)
	// Get the token for userid from the database.
	token, err := app.tokens.GetToken(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	if form.Token != token.Token {
		app.logger.Error("invalid token entered")
		app.sessionManager.Put(r.Context(), "flash", "Code ungültig")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/password_reset", http.StatusSeeOther)
}

type resetPasswordForm struct {
	UserId                  string
	NewPassword             string `form:"newPassword"`
	NewPasswordConfirmation string `form:"newPasswordConfirmation"`
	validator.Validator     `form:"-"`
}

// The resetPassword handler is called by the GET /password_reset/ route.
// It gets the userId from the sessionManager and displays the reset_forgotten_password page
func (app *application) resetPassword(w http.ResponseWriter, r *http.Request) {
	// Get userid from sessionManager
	id, ok := app.sessionManager.Get(r.Context(), "userID").(string)
	if !ok {
		app.logger.Error("type assertion to string failed")
		return
	}
	data := app.newTemplateData(r)
	data.Form = resetPasswordForm{UserId: id}
	app.render(w, r, http.StatusOK, "reset_forgotten_password.tmpl", data)
}

// The resetPasswordPost handler is called by the POST /password_reset_
// gets the userId from the sessionManager, parses the resetPasswordForm and sets the new password in the database.
func (app *application) resetPasswordPost(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.Get(r.Context(), "userID").(string)
	userID, _ := strconv.Atoi(id)
	var form resetPasswordForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "Eingabe fehlt!")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "Mindestens 8 Zeichen eingeben")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "Eingabe fehlt!")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwörter stimmen nicht überein")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "reset_forgotten_password.tmpl", data)
		return
	}
	err = app.users.ResetPassword(userID, form.NewPassword)
	if err != nil {
		app.serverError(w, r, err)
	}
	// Delete existing token for userID
	err = app.tokens.DeleteToken(userID)
	if err != nil {
		app.logger.Error("token deletion failed" + err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// The sendToken method sends the token to the email address of the user.
func (app *application) sendToken(email string, token string) {
	// Send token to email address.
	emailHost := os.Getenv("EMAIL_HOST")
	emailUser := os.Getenv("EMAIL_USERNAME")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	m := mail.NewMsg()
	if err := m.From(emailUser); err != nil {
		app.logger.Info("failed to set From address: %s" + err.Error())
	}
	if err := m.To(email); err != nil {
		app.logger.Info("failed to set To address: %s" + err.Error())
	}
	// Create text and html body.
	msg := "<h1>Reset code: " + token + "</h1>" +
		"<p>Code in Formular eingeben</p>"
	msg2 := "Reset code: " + token
	m.Subject("Passwort reset code")
	m.SetBodyString(mail.TypeTextHTML, msg)
	m.AddAlternativeString(mail.TypeTextPlain, msg2)
	// Create the mail client
	c, err := mail.NewClient(emailHost,
		mail.WithPort(587), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(emailUser),
		mail.WithPassword(emailPassword),
		mail.WithTLSPolicy(mail.TLSOpportunistic))
	if err != nil {
		msg := "failed to create mail client: " + err.Error()
		app.logger.Error(msg)
	}
	// Finally let's send out the mail
	if err := c.DialAndSend(m); err != nil {
		msg := "failed to send email: " + err.Error()
		app.logger.Error(msg)
	}
	app.logger.Info("sent token to " + email)
}

// The storeToken method creates a token and stores the userid, token and timestamp it in the database
func (app *application) storeToken(email string) string {
	// Create token
	token := app.createCode(6)
	// Get userID from database
	var userID int
	userID, err := app.users.GetUserId(email)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
		}
		return ""
	}
	// Delete existing token for userID
	err = app.tokens.DeleteToken(userID)
	if err != nil {
		app.logger.Error("token deletion failed" + err.Error())
		return ""
	}
	// Insert UserID, token and timestamp
	err = app.tokens.InsertToken(userID, token)
	if err != nil {
		app.logger.Error("token insertion failed" + err.Error())
		return ""
	}
	return token
}
