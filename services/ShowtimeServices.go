package services

import (
	"fmt"
	"movie/models"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShowtimeService struct {
	DB *sqlx.DB
}

func NewShowtimeService(db *sqlx.DB) *ShowtimeService {
	return &ShowtimeService{
		db,
	}
}

func (ss *ShowtimeService) AddShowtimes(showtime *models.Showtime) (*models.Showtime, error) {
	showtime.ShowtimeId = uuid.New().String()[:10]

	query := `
	INSERT INTO showtimes (showtimeid, movieid, starttime, endtime, venue, priceperseat, availableseats)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING *
	`
	err := ss.DB.Get(showtime, query,
		showtime.ShowtimeId,
		showtime.MovieId,
		showtime.StartTime,
		showtime.EndTime,
		showtime.Venue,
		showtime.PricePerSeat,
		showtime.AvailableSeats,
		showtime.CreatedAt,
		showtime.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return showtime, nil
}

func (ss *ShowtimeService) DeleteShowtime(ShowtimeId string) error {
	var movieId string
	checkQuery := "SELECT movieid FROM showtimes WHERE showtimeid = $1"

	err := ss.DB.Get(&movieId, checkQuery, ShowtimeId)
	if err != nil {
		return fmt.Errorf("showtime not found: %w", err)
	}

	deleteQuery := "DELETE FROM showtimes WHERE showtimeid = $1"
	_, err = ss.DB.Exec(deleteQuery, ShowtimeId)
	if err != nil {
		return fmt.Errorf("failed to delete showtime: %w", err)
	}

	return nil
}

func (ss *ShowtimeService) UpdateShowtime(Showtime *models.Showtime) (*models.Showtime, error) {
	setClauses := []string{}
	args := []any{}
	argIndex := 1

	if Showtime.MovieId != "" {
		return nil, fmt.Errorf("cannot update movieId")
	}

	if !Showtime.StartTime.IsZero() {
		setClauses = append(setClauses, fmt.Sprintf("starttime = $%d", argIndex))
		args = append(args, Showtime.StartTime)
		argIndex++
	}

	if !Showtime.EndTime.IsZero() {
		setClauses = append(setClauses, fmt.Sprintf("endtime = $%d", argIndex))
		args = append(args, Showtime.EndTime)
		argIndex++
	}

	if Showtime.Venue != "" {
		setClauses = append(setClauses, fmt.Sprintf("venue = $%d", argIndex))
		args = append(args, Showtime.Venue)
		argIndex++
	}

	if Showtime.PricePerSeat != 0 {
		setClauses = append(setClauses, fmt.Sprintf("priceperseat = $%d", argIndex))
		args = append(args, Showtime.PricePerSeat)
		argIndex++
	}

	if Showtime.AvailableSeats != 0 {
		setClauses = append(setClauses, fmt.Sprintf("availableseats = $%d", argIndex))
		args = append(args, Showtime.AvailableSeats)
		argIndex++
	}

	// Always update updatedat
	setClauses = append(setClauses, "updatedat = CURRENT_TIMESTAMP")

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE showtimes
		SET %s
		WHERE showtimeid = $%d
		RETURNING *
	`, strings.Join(setClauses, ", "), argIndex)

	args = append(args, Showtime.ShowtimeId)

	err := ss.DB.Get(Showtime, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update showtime: %w", err)
	}

	return Showtime, nil
}

func (ss *ShowtimeService) GetShowtimesAndMovieData(MovieId string) (*[]models.ShowtimeAndMovie, error) {
	if MovieId == "" {
		return nil, fmt.Errorf("need movieId to query showtimes")
	}

	var showtimes []models.ShowtimeAndMovie

	query := `
	SELECT
	  s.showtimeid,
	  s.starttime,
	  s.endtime,
	  s.venue,
	  s.priceperseat,
	  s.availableseats,
	  m.movieid,
	  m.title,
	  m.genre,
	  m.director,
	  m.posterimage
	FROM showtimes s
	JOIN movies m ON s.movieid = m.movieid
	WHERE s.movieid = $1
	`

	err := ss.DB.Select(&showtimes, query, MovieId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch showtimes: %w", err)
	}

	return &showtimes, nil
}

func (ss *ShowtimeService) CheckAvailableSeats(MovieId string) (*models.SeatsAndPrice, error) {
	if MovieId == "" {
		return nil, fmt.Errorf("need MovieId to check available seats in db")
	}
	var seatsAndPriceData models.SeatsAndPrice

	query := `
	SELECT availableseats, priceperseat
	FROM showtimes
	WHERE movieid = $1
	`
	err := ss.DB.Get(&seatsAndPriceData, query, MovieId)
	if err != nil {
		return nil, fmt.Errorf("error fetching seats for given movieid: %w", err)
	}

	return &seatsAndPriceData, nil
}
