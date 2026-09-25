package handlers

import (
	"errors"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleGetShowTimeDetails(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("showTimeId")
	showTimeId, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid show time id", err)
		return
	}

	details, err := cfg.DB.GetShowTimeDetails(r.Context(), showTimeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Show time not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching show time details", err)
		return
	}

	showTime := details.ShowTime
	cinema := details.Cinema

	//get current bookings
	bookings, err := cfg.DB.GetShowTimeReservations(r.Context(), database.GetShowTimeReservationsParams{
		ShowTimeID: showTime.ID,
		Status:     "booked",
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, "Error fetching show time bookings", err)
			return
		}
	}

	type response struct {
		ShowTime        database.ShowTime `json:"show_time"`
		CinemaCapacity  int               `json:"cinema_capacity"`
		CurrentBookings int               `json:"current_bookings"`
		CurrentRevenue  int               `json:"current_revenue"`
	}

	respondWithJSON(w, http.StatusOK, response{
		ShowTime:        showTime,
		CinemaCapacity:  cinema.SeatMap.TotalSeats,
		CurrentBookings: len(bookings),
		CurrentRevenue:  len(bookings) * int(showTime.Price),
	})
}
