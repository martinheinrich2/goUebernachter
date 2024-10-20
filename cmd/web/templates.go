package main

import (
	"github.com/martinheinrich2/goUebernachter/internal/models"
	"github.com/martinheinrich2/goUebernachter/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

// humanDate returns a human-readable string representation of a time.Time object
func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 um 15:04")
}

// humanBirthDay returns a human-readable string representation of a time.Time object
func humanBirthDay(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02.01.2006")
}

// dateReversed returns a human readable string representation of a time.Time object
func dateReversed(t time.Time) string {
	if t.IsZero() {
		t := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return t.UTC().Format("2006-01-02")
	}
	return t.UTC().Format("2006-01-02")
}

// dateTimeForm returns the datetime for a form datetime picker
func dateTimeForm(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02T15:04")
}

// humanDateTime returns a human readable string representation of a time.Time object
func humanDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02.01.2006 15:04:05")

} // humanDateTime2 returns a human readable string representation of a time.Time object
func humanDateTime2(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	u := time.Date(2010, 0, 0, 0, 0, 0, 0, time.UTC)
	if t.Before(u) {
		return "kein Termin vorhanden"
	}
	return t.UTC().Format("02.01.2006 um 15:04:05")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts al a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate":      humanDate,
	"humanBirthDay":  humanBirthDay,
	"dateReversed":   dateReversed,
	"humanDateTime":  humanDateTime,
	"humanDateTime2": humanDateTime2,
	"dateTimeForm":   dateTimeForm,
}

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to the HTML templates.
type templateData struct {
	CurrentYear     int
	Stay            models.Stay
	Stays           []models.Stay
	StayJoinUser    []models.StayJoinUser
	StayJoinGuest   []models.StayJoinGuest
	Guest           models.Guest
	Guests          []models.Guest
	SocialWorker    models.User
	User            models.User
	Users           []models.User
	AuthUser        models.User
	Date            string
	Form            any
	Flash           string
	IsAuthenticated bool
	IsAdmin         bool
	IsActive        bool
	StayCount       []models.StayCount
	StayCount2      []models.StayCount2
	Metadata        models.Metadata
	CSRFToken       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act al the cache.
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded
	// filesystem which match the pattern 'html/pages/*.tmpl'. This essentially
	// gives us a slice of all the 'page' templates for the application, just
	// like before.
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
