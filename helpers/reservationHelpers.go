package helpers

import (
	"movie/models"

	"github.com/jmoiron/sqlx"
)

func GetShowtimeData(db *sqlx.DB, showtimeId string) (*models.Showtime, error) {
	var showtime models.Showtime
	query := `
	SELECT * FROM showtimes WHERE showtimeid =  $1
	`
	err := db.Get(&showtime, query, showtimeId)
	if err != nil {
		return nil, err
	}

	return &showtime, nil
}
