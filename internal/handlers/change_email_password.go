package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Joshua-SV/booking-and-appointment-site/db/generated"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/auth"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/utils"
	"github.com/google/uuid"
)

func (cfg *ApiModel) ChangeEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	userID := cfg.CheckAccessToken(w, r)
	if userID == uuid.Nil {
		return
	}

	defer r.Body.Close()

	// get json body from request
	var user User
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not read request body")
		return
	}

	// create new hashed password
	hashedPass, err := auth.HashedPassword(user.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create hash")
		return
	}

	// store new email and password in database
	params := generated.SetEmailAndPasswordParams{
		Email:        user.Email,
		PasswordHash: hashedPass,
		ID:           userID,
	}
	err = cfg.db.SetEmailAndPassword(r.Context(), params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not update email and password")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"email": user.Email})
}
