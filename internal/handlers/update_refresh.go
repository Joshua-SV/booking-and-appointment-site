package handlers

import (
	"net/http"
	"time"

	"github.com/Joshua-SV/booking-and-appointment-site/internal/auth"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/utils"
)

func (cfg *ApiModel) Refresh(w http.ResponseWriter, r *http.Request) {
	// parse refresh token string
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	// get user info by checking if the refresh token is still valid
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), tokenStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	// create new JWT token
	jwtToken, err := auth.CreateJWT(user.ID, cfg.serverKey, time.Hour)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create access token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"token": jwtToken,
	})
}
