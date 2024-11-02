package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/joho/godotenv"
	"github.com/martinheinrich2/goUebernachter/internal/models"
	"html/template"
	"log"
	"log/slog"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"time"
)

// Define application struct to hold the application-wide dependencies for the web application.
// This allows to make the model object available to the handlers.
// Define handler functions as methods against application.
type application struct {
	debug          bool
	logger         *slog.Logger
	guests         models.GuestModelInterface
	users          models.UserModelInterface
	stays          models.StayModelInterface
	tokens         models.TokenModelInterface
	sessionManager *scs.SessionManager
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
}

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Define command-line flags with default values.
	// $ go run ./cmd/web -addr=":8000"
	addr := flag.String("addr", ":4000", "HTTP network address")
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	// Initialize logger only logging errors
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	//logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	//  Create a connection pool from the openDB() function.
	db, err := openDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Initialize a new template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance
	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager.
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Initialize a new instance of the app struct, containing the connection pool.
	// This instance contains the application dependencies.
	app := &application{
		debug:          *debug,
		logger:         logger,
		guests:         &models.GuestModel{DB: db},
		users:          &models.UserModel{DB: db},
		stays:          &models.StayModel{DB: db},
		tokens:         &models.TokenModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	// is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	// Initialize a new http.Server struct to hold server behavior.
	srv := &http.Server{
		Addr:      *addr,
		Handler:   app.routes(),
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
		// Insert Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// Get local IPs
	localIPs, err := GetLocalIPs()
	if err != nil {
		logger.Info(err.Error())
	}
	// Create a time.Ticker loop to run db backup every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for t := range ticker.C {
			logger.Info("Backup at ", t.UTC().Format(time.RFC3339))
			// Create backup
			err = backup()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			fmt.Println("backup successful!")

		}
	}()

	fmt.Println("Starting Server at:", localIPs[0].String()+srv.Addr)
	logger.Info("starting server at:", slog.String("addr", srv.Addr))
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", os.Getenv("SRC_DB"))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
