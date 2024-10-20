package main

import (
	"github.com/justinas/alice"
	"github.com/martinheinrich2/goUebernachter/ui"
	"net/http"
)

// The routes() method returns a servemux containing the application routes.
func (app *application) routes() http.Handler {
	// Use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()

	// Use the http.FileServerFS() function to create a HTTP handler which
	// serves the embedded files in ui.Files. The static files are
	// contained in the "static" and "html" folders of the ui.Files
	// embedded filesystem.
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	// Unprotected application routes using the "dynamic" middleware chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	//mux.Handle("GET /article/view/{id}", dynamic.ThenFunc(app.articleView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("GET /account/forgot_password", dynamic.ThenFunc(app.forgotPassword))
	mux.Handle("POST /account/forgot_password", dynamic.ThenFunc(app.forgotPasswordPost))
	mux.Handle("GET /enter_token/{email}", dynamic.ThenFunc(app.enterToken))
	mux.Handle("POST /enter_token/{email}", dynamic.ThenFunc(app.enterTokenPost))
	mux.Handle("GET /password_reset", dynamic.ThenFunc(app.resetPassword))
	mux.Handle("POST /password_reset", dynamic.ThenFunc(app.resetPasswordPost))
	mux.Handle("GET /inactive", dynamic.ThenFunc(app.accountInactive))

	// Protected (authenticated only) application routes using the "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	admin := dynamic.Append(app.requireAdmin)

	mux.Handle("GET /guest/create", protected.ThenFunc(app.guestCreate))
	mux.Handle("POST /guest/create", protected.ThenFunc(app.guestCreatePost))
	mux.Handle("GET /guest/view/{id}", protected.ThenFunc(app.guestView))
	mux.Handle("GET /guests", protected.ThenFunc(app.allGuestsView))
	mux.Handle("GET /guest/search", protected.ThenFunc(app.guestSearch))
	mux.Handle("POST /guest/search", protected.ThenFunc(app.guestSearchPost))
	mux.Handle("GET /guest/update/{id}", protected.ThenFunc(app.guestUpdate))
	mux.Handle("POST /guest/update/{id}", protected.ThenFunc(app.guestUpdatePost))

	mux.Handle("GET /stay/create/{id}", protected.ThenFunc(app.stayCreate))
	mux.Handle("POST /stay/create/{id}", protected.ThenFunc(app.stayCreatePost))
	mux.Handle("GET /stays", protected.ThenFunc(app.allStaysView))
	mux.Handle("GET /stay/appointmentopen", protected.ThenFunc(app.appointmentOpenView))
	mux.Handle("GET /stay/staynotprocessed", protected.ThenFunc(app.stayNotProcessedView))
	mux.Handle("GET /stay/detail/{id}", protected.ThenFunc(app.stayDetail))
	mux.Handle("POST /stay/detail_print/{id}", protected.ThenFunc(app.stayDetailPrint))
	mux.Handle("GET /stay/update/{id}", protected.ThenFunc(app.stayUpdate))
	mux.Handle("POST /stay/update/{id}", protected.ThenFunc(app.stayUpdatePost))
	mux.Handle("GET /stay/appointment/{id}", protected.ThenFunc(app.postAppointmentDone))
	mux.Handle("GET /stay/processed/{id}", protected.ThenFunc(app.postDataProcessed))

	mux.Handle("GET /stats", protected.ThenFunc(app.stats))

	mux.Handle("GET /user/update/{id}", admin.ThenFunc(app.userUpdate))
	mux.Handle("POST /user/update/{id}", admin.ThenFunc(app.userUpdatePost))
	mux.Handle("GET /users", protected.ThenFunc(app.allUsersView))
	mux.Handle("GET /user/password/reset/{id}", admin.ThenFunc(app.userPasswordReset))
	mux.Handle("POST /user/password/reset/{id}", admin.ThenFunc(app.userPasswordResetPost))

	mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))
	mux.Handle("GET /account/password/update", protected.ThenFunc(app.accountPasswordUpdate))
	mux.Handle("POST /account/password/update", protected.ThenFunc(app.accountPasswordUpdatePost))

	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
