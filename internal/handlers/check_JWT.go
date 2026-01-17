package handlers

import (
	"net/http"

	"github.com/Joshua-SV/booking-and-appointment-site/internal/auth"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/utils"
	"github.com/google/uuid"
)

// helper function to check access token and return user uuid
func (cfg *ApiModel) CheckAccessToken(w http.ResponseWriter, r *http.Request) uuid.UUID {
	// get access token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid token: missing token")
		return uuid.Nil
	}

	// check access token
	userID, err := auth.ValidateJWT(accessToken, cfg.serverKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid token")
		return uuid.Nil
	}

	return userID
}
