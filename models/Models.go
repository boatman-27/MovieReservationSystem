package models

import (
	"time"

	"github.com/google/uuid"
)

// === === === === ===
//
//	=== User Data ===
//
// === === === === ===
type User struct {
	UserId    uuid.UUID `json:"userId" db:"userid"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedat"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// === === === === ===
//
// === Movie Data ===
//
// === === === === ===
type Movie struct {
	MovieId     string    `json:"MovieId" db:"movieid"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Genre       string    `json:"genre" db:"genre"`
	Duration    int       `json:"duration" db:"duration"`
	Director    string    `json:"director" db:"director"`
	PosterImage string    `json:"posterImage" db:"posterimage"`
	ReleaseDate time.Time `json:"releaseDate" db:"releasedate"`
	CreatedAt   time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updatedat"`
}

// === === === === ===
//
// === Showtime Data ===
//
// === === === === ===
type Showtime struct {
	ShowtimeId     string    `json:"showtimeId" db:"showtimeid"`
	MovieId        string    `json:"movieId" db:"movieid"`
	StartTime      time.Time `json:"startTime" db:"starttime"`
	EndTime        time.Time `json:"endTime" db:"endtime"`
	Venue          string    `json:"venue" db:"venue"`
	PricePerSeat   float64   `json:"pricePerSeat" db:"priceperseat"`
	AvailableSeats int       `json:"availableSeats" db:"availableseats"`
	CreatedAt      time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updatedat"`
}

type ShowtimeAndMovie struct {
	ShowtimeId     string    `db:"showtimeid" json:"showtimeId"`
	StartTime      time.Time `db:"starttime" json:"startTime"`
	EndTime        time.Time `db:"endtime" json:"endTime"`
	Venue          string    `db:"venue" json:"venue"`
	PricePerSeat   float64   `db:"priceperseat" json:"pricePerSeat"`
	AvailableSeats int       `db:"availableseats" json:"availableSeats"`

	MovieId     string `db:"movieid" json:"movieId"`
	Title       string `db:"title" json:"title"`
	Genre       string `db:"genre" json:"genre"`
	Director    string `db:"director" json:"director"`
	PosterImage string `db:"posterimage" json:"posterImage"`
}

type SeatsAndPrice struct {
	PricePerSeat   float64 `db:"priceperseat" json:"pricePerSeat"`
	AvailableSeats int     `db:"availableseats" json:"availableSeats"`
}

// === === === === ===
//
// === Reservation Data ===
//
// === === === === ===
type Reservation struct {
	ReservationId   string    `json:"reservationId" db:"reservationid"`
	UserId          uuid.UUID `json:"userId" db:"userid"`
	ShowtimeId      string    `json:"showtimeId" db:"showtimeid"`
	NumberOfSeats   int       `json:"numberOfSeats" db:"numberofseats"`
	TotalPrice      float64   `json:"totalPrice" db:"totalprice"`
	ReservationDate time.Time `json:"reservationDate" db:"reservationdate"`
}

type BookingData struct {
	ShowtimeId string    `json:"showtimeId"`
	UserId     uuid.UUID `json:"userId"`
	Seats      int       `json:"seats"`
}
