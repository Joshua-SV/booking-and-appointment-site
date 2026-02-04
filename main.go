package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Joshua-SV/booking-and-appointment-site/db/generated"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/handlers"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/pubsub"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// optional: ignore in prod; helpful log in dev
		fmt.Println("warning: .env not loaded:", err)
	}

	// Connect to Postgres database
	dburl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		fmt.Printf("could not open connection to database: %v", err)
		os.Exit(1)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		fmt.Printf("could not connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// declare api model with database connection
	cfg := handlers.ApiModel{}
	cfg.SetDB(generated.New(db))
	// set server key for JWT signing
	cfg.SetServerKey(os.Getenv("SERVER_SECRET_KEY"))
	// set dev access key
	cfg.SetDevAccess(os.Getenv("PLATFORM"))
	// set rabbitmq url
	cfg.SetRabbitmqURL(os.Getenv("RABBITMQ_URL"))

	// connect to rabbitmq server
	rabbitConn, err := amqp.Dial(cfg.GetRabbitmqURL())
	if err != nil {
		fmt.Printf("could not connect to rabbitmq server: %v", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// subscribe to appointment creation queue
	err = pubsub.SubscribeJSON(rabbitConn, "direct_appointment_exch", "appointment_creation_queue", "appointment.created", pubsub.DurableQueue)

	// declare a serve mux to handle enpoint routing
	servMux := http.NewServeMux()

	// register handlers for api namespace
	servMux.HandleFunc("POST /api/users", cfg.CreateUser)
	servMux.HandleFunc("POST /api/login", cfg.LoginUser)
	servMux.HandleFunc("POST /api/refresh", cfg.Refresh)
	servMux.HandleFunc("POST /api/revoke", cfg.RevokeRefreshToken)
	servMux.HandleFunc("PUT /api/users", cfg.ChangeEmailAndPassword)

	// admin endpoints
	servMux.HandleFunc("POST /admin/reset", cfg.Reset)

	// configure server
	server := http.Server{
		Addr:              ":8080",
		Handler:           servMux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	// start server
	fmt.Println("Starting server on :8080...")
	err = server.ListenAndServe() // pauses here until server stops
	if err != nil {
		fmt.Printf("could not start server: %v", err)
		os.Exit(1)
	}
}
