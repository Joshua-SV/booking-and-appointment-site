package handlers

import (
	"time"

	"github.com/Joshua-SV/booking-and-appointment-site/db/generated"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ApiModel struct {
	db        *generated.Queries
	devAccess string
	key       string
	polkaKey  string
}

// declare struct to hold user info
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

// setter methods for api model
func (a *ApiModel) SetDB(db *generated.Queries) {
	a.db = db
}

func (a *ApiModel) SetDevAccess(devAccess string) {
	a.devAccess = devAccess
}

func (a *ApiModel) SetKey(key string) {
	a.key = key
}

func (a *ApiModel) SetPolkaKey(polkaKey string) {
	a.polkaKey = polkaKey
}
