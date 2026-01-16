package handlers

import (
	"net/http"

	"github.com/Joshua-SV/booking-and-appointment-site/internal/auth"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/utils"
)

func (cfg *ApiModel) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token string
	reToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	// update token to be revoked
	cfg.db.RevokeRefreshToken(r.Context(), reToken)

	// assume error means that token does not exist which is the same outcome as revoked
	w.WriteHeader(http.StatusNoContent)
}
