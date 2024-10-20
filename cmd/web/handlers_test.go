package main

import (
	"github.com/martinheinrich2/goUebernachter/internal/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	// Create a new instance of the application struct which uses the mocked dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running ent to end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "Server is running OK!")
}

func TestUserSignup(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set
	// up the test server for running an end-to-end test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET /user/signup request and then extract the CSRF token from the
	// response body.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	// Log the CSRF token value in our test output using the t.Logf() function.
	// The t.Logf() function works in the same way as fmt.Printf(), but writes
	// the provided message to the test output.
	t.Logf("CSRF token is: %q", csrfToken)
}

//func TestUserSignup(t *testing.T) {
//	// Create the application struct containing the mocked dependencies and set
//	// up the test server for running an end-to-end test.
//	app := newTestApplication(t)
//	ts := newTestServer(t, app.routes())
//	defer ts.Close()
//
//	// Make a GET /user/signup request and then extract the CSRF token from the
//	// response body.
//	_, _, body := ts.get(t, "/user/signup")
//	validCSRFToken := extractCSRFToken(t, body)
//
//	const (
//		validLastName  = "Ross"
//		validFirstName = "Bob"
//		validPassword  = "validPa$$word"
//		validEmail     = "bob@example.com"
//		validJobTitle  = "Socialworker"
//		validRoom      = "H123"
//		formTag        = "<form action='/user/signup' method='POST' novalidate>"
//	)
//
//	tests := []struct {
//		name         string
//		lastName     string
//		firstName    string
//		userEmail    string
//		userPassword string
//		userJobTitle string
//		userRoom     string
//		csrfToken    string
//		wantCode     int
//		wantFormTag  string
//	}{
//		{
//			name:         "Valid submission",
//			lastName:     validLastName,
//			firstName:    validFirstName,
//			userEmail:    validEmail,
//			userPassword: validPassword,
//			userJobTitle: validJobTitle,
//			userRoom:     validRoom,
//			csrfToken:    validCSRFToken,
//			wantCode:     http.StatusSeeOther,
//		},
//		{
//			name:         "Invalid CSRF Token",
//			lastName:     validLastName,
//			firstName:    validFirstName,
//			userEmail:    validEmail,
//			userPassword: validPassword,
//			userJobTitle: validJobTitle,
//			userRoom:     validRoom,
//			csrfToken:    "wrongToken",
//			wantCode:     http.StatusBadRequest,
//		},
//		{
//			name:         "Empty Last name",
//			lastName:     "",
//			firstName:    validFirstName,
//			userEmail:    validEmail,
//			userPassword: validPassword,
//			csrfToken:    validCSRFToken,
//			wantCode:     http.StatusUnprocessableEntity,
//			wantFormTag:  formTag,
//		},
//		{
//			name:         "Empty email",
//			lastName:     validLastName,
//			firstName:    validFirstName,
//			userEmail:    "",
//			userPassword: validPassword,
//			csrfToken:    validCSRFToken,
//			wantCode:     http.StatusUnprocessableEntity,
//			wantFormTag:  formTag,
//		},
//		{
//			name:         "Empty password",
//			lastName:     validLastName,
//			firstName:    validFirstName,
//			userEmail:    validEmail,
//			userPassword: "",
//			csrfToken:    validCSRFToken,
//			wantCode:     http.StatusUnprocessableEntity,
//			wantFormTag:  formTag,
//		},
//		{
//			name:         "Invalid email",
//			lastName:     validLastName,
//			firstName:    validFirstName,
//			userEmail:    "bob@example.",
//			userPassword: validPassword,
//			csrfToken:    validCSRFToken,
//			wantCode:     http.StatusUnprocessableEntity,
//			wantFormTag:  formTag,
//		},
//		{
//			name:         "Short password",
//			lastName:     validLastName,
//			firstName:    validFirstName,
//			userEmail:    validEmail,
//			userPassword: "pa$$",
//			csrfToken:    validCSRFToken,
//			wantCode:     http.StatusUnprocessableEntity,
//			wantFormTag:  formTag,
//		},
//		{
//			name:         "Duplicate email",
//			lastName:     validLastName,
//			firstName:    validFirstName,
//			userEmail:    "dupe@example.com",
//			userPassword: validPassword,
//			csrfToken:    validCSRFToken,
//			wantCode:     http.StatusUnprocessableEntity,
//			wantFormTag:  formTag,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			fmt.Println("handlers_test.go url.Values ", url.Values{})
//			form := url.Values{}
//			form.Add("lastname", tt.lastName)
//			form.Add("firstname", tt.firstName)
//			form.Add("email", tt.userEmail)
//			form.Add("job_title", tt.userJobTitle)
//			form.Add("room", tt.userRoom)
//			form.Add("password", tt.userPassword)
//			form.Add("csrf_token", tt.csrfToken)
//
//			code, _, body := ts.postForm(t, "/user/signup", form)
//			fmt.Println("test returns code: ", code)
//			assert.Equal(t, code, tt.wantCode)
//
//			if tt.wantFormTag != "" {
//				assert.StringContains(t, body, tt.wantFormTag)
//			}
//		})
//	}
//}
