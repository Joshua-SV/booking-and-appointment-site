package handlers

import (
	"time"

	"github.com/Joshua-SV/booking-and-appointment-site/db/generated"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ApiModel struct {
	db          *generated.Queries
	devAccess   string
	serverKey   string
	apiKey      string
	rabbitmqURL string
}

type AppointmentRequest struct {
	AppointmentTime time.Time `json:"appointment_time"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
}

// declare struct to hold user info
type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password"`
}

// setter methods for api model
func (a *ApiModel) SetDB(db *generated.Queries) {
	a.db = db
}

func (a *ApiModel) SetDevAccess(devAccess string) {
	a.devAccess = devAccess
}

func (a *ApiModel) SetServerKey(serverKey string) {
	a.serverKey = serverKey
}

func (a *ApiModel) SetAPIKey(apiKey string) {
	a.apiKey = apiKey
}

func (a *ApiModel) SetRabbitmqURL(rabbitmqURL string) {
	a.rabbitmqURL = rabbitmqURL
}

func (a *ApiModel) GetRabbitmqURL() string {
	return a.rabbitmqURL
}
