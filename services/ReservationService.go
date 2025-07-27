package services

import (
	"fmt"
	"movie/helpers"
	"movie/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReservationService struct {
	DB *sqlx.DB
}

func NewReservationService(db *sqlx.DB) *ReservationService {
	return &ReservationService{
		DB: db,
	}
}

func (rs *ReservationService) GetUpcomingEvents(userId string) (*[]models.Reservation, error) {
	var events []models.Reservation
	query := `SELECT * FROM reservations WHERE userid = $1 AND reservationdate > CURRENT_TIMESTAMP`

	err := rs.DB.Select(&events, query, userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching upcoming events: %w", err)
	}

	return &events, nil
}

func (rs *ReservationService) BookSeats(bookingData *models.BookingData) (*models.Reservation, error) {
	showtime, err := helpers.GetShowtimeData(rs.DB, bookingData.ShowtimeId)
	if err != nil {
		return nil, fmt.Errorf("error fetching showtime data: %w", err)
	}

	if showtime.AvailableSeats < bookingData.Seats {
		return nil, fmt.Errorf("not enough seats available, only %d are left", showtime.AvailableSeats)
	}

	reservation := &models.Reservation{
		ReservationId:   uuid.New().String()[:10],
		UserId:          bookingData.UserId,
		ShowtimeId:      bookingData.ShowtimeId,
		NumberOfSeats:   bookingData.Seats,
		TotalPrice:      float64(bookingData.Seats) * showtime.PricePerSeat,
		ReservationDate: showtime.StartTime,
	}

	// Start a transaction
	tx, err := rs.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	insertQuery := `
		INSERT INTO reservations (reservationid, userid, showtimeid, numberofseats, totalprice, reservationdate)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`

	err = tx.Get(reservation, insertQuery,
		reservation.ReservationId,
		reservation.UserId,
		reservation.ShowtimeId,
		reservation.NumberOfSeats,
		reservation.TotalPrice,
		reservation.ReservationDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert reservation: %w", err)
	}

	updateSeatsQuery := `
		UPDATE showtimes
		SET availableseats = $2
		WHERE showtimeid = $1
	`

	newSeats := showtime.AvailableSeats - bookingData.Seats
	_, err = tx.Exec(updateSeatsQuery, bookingData.ShowtimeId, newSeats)
	if err != nil {
		return nil, fmt.Errorf("failed to update available seats: %w", err)
	}

	return reservation, nil
}

func (rs *ReservationService) CancelReservation(reservationId string) error {
	var reservationDate time.Time
	fetchQuery := `
	SELECT reservationdate FROM reservations WHERE reservationid = $1
	`
	err := rs.DB.Get(&reservationDate, fetchQuery, reservationId)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	deleteQuery := `
	DELETE FROM reservations WHERE reservationid = $1
	`
	_, err = rs.DB.Exec(deleteQuery, reservationId)
	if err != nil {
		return fmt.Errorf("error deleting reservation: %w", err)
	}

	return nil
}

func (rs *ReservationService) GetAllUserReservations(userId string) (*[]models.Reservation, error) {
	var events []models.Reservation
	query := `SELECT * FROM reservations WHERE userid = $1`

	err := rs.DB.Select(&events, query, userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching upcoming events: %w", err)
	}

	return &events, nil
}

func (rs *ReservationService) GetAllReservations() (*[]models.Reservation, error) {
	var reservations []models.Reservation

	query := `SELECT * FROM reservations`
	err := rs.DB.Select(&reservations, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all reservations: %w", err)
	}

	return &reservations, nil
}
