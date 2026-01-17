package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Joshua-SV/booking-and-appointment-site/db/generated"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/auth"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/utils"
)

func (cfg *ApiModel) LoginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User

	// parse request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	// retrieve user from database
	dbUser, err := cfg.db.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid user data")
		return
	}

	// verify password
	if err := auth.CheckPasswordvsHash(user.Password, dbUser.PasswordHash); err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid user data")
		return
	}

	// create the JWT token string
	accessToken, err := auth.CreateJWT(dbUser.ID, cfg.serverKey, time.Hour) // 1 hour expiration
	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, "could not create token")
		return
	}

	// create refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create refresh token")
		return
	}
	// prepare params for storing refresh token in database
	refreshExpires := time.Now().UTC().Add(60 * 24 * time.Hour) // expires in 60 days
	params := generated.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: refreshExpires,
	}

	// store refresh token in database
	_, err = cfg.db.CreateRefreshToken(r.Context(), params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not store refresh token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":            dbUser.ID.String(),
		"email":         user.Email,
		"created_at":    dbUser.CreatedAt.String(),
		"updated_at":    dbUser.UpdatedAt.String(),
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}
