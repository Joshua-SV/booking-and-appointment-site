package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Joshua-SV/booking-and-appointment-site/db/generated"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/auth"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/utils"
)

func (cfg *ApiModel) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	defer r.Body.Close()

	// parse request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	// hash the password
	hashedPassword, err := auth.HashedPassword(user.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	// create user in database
	params := generated.CreateUserParams{
		Email:        user.Email,
		PasswordHash: hashedPassword, // In a real application, the password is stored hashed
	}
	responseUser, err := cfg.db.CreateUser(r.Context(), params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create user in database")
		return
	}

	// respond with created user (excluding password)
	user.ID = responseUser.ID
	user.CreatedAt = responseUser.CreatedAt
	user.UpdatedAt = responseUser.UpdatedAt
	user.Password = "" // do not send password back

	// send response
	err = utils.RespondWithJSON(w, http.StatusCreated, user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to send response")
		return
	}
}
