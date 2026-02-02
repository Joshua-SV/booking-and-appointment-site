package handlers

import (
	"net/http"

	"github.com/Joshua-SV/booking-and-appointment-site/internal/utils"
)

func (cfg *ApiModel) Reset(w http.ResponseWriter, r *http.Request) {
	// Only allow reset in Dev mode
	if cfg.devAccess != "Dev" {
		utils.RespondWithError(w, http.StatusForbidden, "reset endpoint is only available in Dev mode")
		return
	}

	// perform reset operations here (e.g., clear database tables)
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to reset users in database")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "system reset successfully",
	})
}
