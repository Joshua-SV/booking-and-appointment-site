package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Joshua-SV/booking-and-appointment-site/db/generated"
	"github.com/Joshua-SV/booking-and-appointment-site/internal/utils"
	"github.com/google/uuid"
)

func (cfg *ApiModel) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment AppointmentRequest
	defer r.Body.Close()

	// parse request body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&appointment)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	// Get user ID from access token
	userID := cfg.CheckAccessToken(w, r)
	if userID == uuid.Nil {
		return // Error already handled in CheckAccessToken
	}

	// create appointment in database
	params := generated.CreateAppointmentParams{
		UserID:          userID,
		AppointmentTime: appointment.AppointmentTime,
		Status:          appointment.Status,
		Notes:           appointment.Notes,
	}
	responseAppointment, err := cfg.db.CreateAppointment(r.Context(), params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create appointment in database")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, responseAppointment)
}
