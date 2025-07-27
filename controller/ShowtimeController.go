package controllers

import (
	"movie/models"
	"movie/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShowtimeController struct {
	ShowtimeService *services.ShowtimeService
}

func NewShowtimeController(showtimeService *services.ShowtimeService) *ShowtimeController {
	return &ShowtimeController{
		showtimeService,
	}
}

func (sc *ShowtimeController) AddShowtimes(c *gin.Context) {
	var newShowtime models.Showtime
	if err := c.ShouldBindJSON(&newShowtime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	showtime, err := sc.ShowtimeService.AddShowtimes(&newShowtime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"showtime": showtime})
}

func (sc *ShowtimeController) UpdateShowtime(c *gin.Context) {
	var updatedShowtime models.Showtime
	if err := c.ShouldBindJSON(&updatedShowtime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedShowtime.ShowtimeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "showtimeId is required"})
		return
	}

	if updatedShowtime.MovieId != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot change MovieId while updating showtime"})
		return
	}

	showtime, err := sc.ShowtimeService.UpdateShowtime(&updatedShowtime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"showtime": showtime})
}

func (sc *ShowtimeController) DeleteShowtime(c *gin.Context) {
	showtimeId := c.Query("showtimeId")
	if showtimeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "showtimeId is required"})
		return
	}

	err := sc.ShowtimeService.DeleteShowtime(showtimeId)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case strings.Contains(err.Error(), "not found"):
			status = http.StatusNotFound
		case strings.Contains(err.Error(), "delete showtime"):
			status = http.StatusInternalServerError
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Showtime deleted successfully"})
}

func (sc *ShowtimeController) GetShowtimeAndMovie(c *gin.Context) {
	movieId := c.Query("movieId")
	if movieId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MovieId is required and cannot be empty"})
		return
	}

	showtimeAndMovieData, err := sc.ShowtimeService.GetShowtimesAndMovieData(movieId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"showtimeAndMovieData": showtimeAndMovieData})
}

func (sc *ShowtimeController) CheckAvailableSeats(c *gin.Context) {
	movieId := c.Query("movieId")
	if movieId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MovieId is required and cannot be empty"})
		return
	}

	seatsAndPriceData, err := sc.ShowtimeService.CheckAvailableSeats(movieId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if seatsAndPriceData.AvailableSeats == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No seats available for that movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Free seats":     seatsAndPriceData.AvailableSeats,
		"Price per seat": seatsAndPriceData.PricePerSeat,
	})
}
