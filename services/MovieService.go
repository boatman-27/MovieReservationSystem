package services

import (
	"fmt"
	"movie/models"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MovieService struct {
	DB *sqlx.DB
}

func NewMovieService(db *sqlx.DB) *MovieService {
	return &MovieService{
		DB: db,
	}
}

func (ms *MovieService) AddMovie(movie *models.Movie) (*models.Movie, error) {
	movie.MovieId = uuid.New().String()[:10]

	query := `
	INSERT INTO movies (movieid, title, description, genre, duration, director, posterimage, releasedate)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING *
	`
	err := ms.DB.Get(movie, query,
		movie.MovieId,
		movie.Title,
		movie.Description,
		movie.Genre,
		movie.Duration,
		movie.Director,
		movie.PosterImage,
		movie.ReleaseDate,
	)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (ms *MovieService) DeleteMovie(movieId string) error {
	var title string
	checkQuery := "SELECT title FROM movies WHERE movieid = $1"
	err := ms.DB.Get(&title, checkQuery, movieId)
	if err != nil {
		return fmt.Errorf("movie not found: %w", err)
	}

	deleteQuery := "DELETE FROM movies WHERE movieid = $1"
	_, err = ms.DB.Exec(deleteQuery, movieId)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	return nil
}

func (ms *MovieService) GetMovies() ([]*models.Movie, error) {
	var movies []*models.Movie
	fetchQuery := "SELECT * FROM movies"
	err := ms.DB.Select(&movies, fetchQuery)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (ms *MovieService) UpdateMovies(movie *models.Movie) (*models.Movie, error) {
	setClauses := []string{}
	args := []any{}
	argIndex := 1

	if movie.Title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, movie.Title)
		argIndex++
	}
	if movie.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, movie.Description)
		argIndex++
	}
	if movie.Genre != "" {
		setClauses = append(setClauses, fmt.Sprintf("genre = $%d", argIndex))
		args = append(args, movie.Genre)
		argIndex++
	}
	if movie.Duration != 0 {
		setClauses = append(setClauses, fmt.Sprintf("duration = $%d", argIndex))
		args = append(args, movie.Duration)
		argIndex++
	}
	if movie.Director != "" {
		setClauses = append(setClauses, fmt.Sprintf("director = $%d", argIndex))
		args = append(args, movie.Director)
		argIndex++
	}
	if movie.PosterImage != "" {
		setClauses = append(setClauses, fmt.Sprintf("posterimage = $%d", argIndex))
		args = append(args, movie.PosterImage)
		argIndex++
	}
	if !movie.ReleaseDate.IsZero() {
		setClauses = append(setClauses, fmt.Sprintf("releasedate = $%d", argIndex))
		args = append(args, movie.ReleaseDate)
		argIndex++
	}

	// Always update updatedat
	setClauses = append(setClauses, "updatedat = CURRENT_TIMESTAMP")

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// WHERE clause
	query := fmt.Sprintf(`
		UPDATE movies
		SET %s
		WHERE movieid = $%d
		RETURNING *
	`, strings.Join(setClauses, ", "), argIndex)

	args = append(args, movie.MovieId)

	err := ms.DB.Get(movie, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update movie: %w", err)
	}

	return movie, nil
}

func (ms *MovieService) GetMovieById(movieId string) (*models.Movie, error) {
	var movie models.Movie
	query := `SELECT * FROM movies WHERE movieid =  $1`
	err := ms.DB.Get(&movie, query, movieId)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	return &movie, nil
}
