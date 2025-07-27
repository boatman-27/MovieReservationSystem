package controllers

import (
	"movie/models"
	"movie/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type MovieController struct {
	MovieService *services.MovieService
}

func NewMovieController(movieService *services.MovieService) *MovieController {
	return &MovieController{
		movieService,
	}
}

func (mc *MovieController) AddMovie(c *gin.Context) {
	var newMovie models.Movie
	if err := c.ShouldBindJSON(&newMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := mc.MovieService.AddMovie(&newMovie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

func (mc *MovieController) DeleteMovie(c *gin.Context) {
	movieId := c.Query("movieId")
	if movieId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "movieId is required"})
		return
	}

	err := mc.MovieService.DeleteMovie(movieId)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case strings.Contains(err.Error(), "not found"):
			status = http.StatusNotFound
		case strings.Contains(err.Error(), "delete movie"):
			status = http.StatusInternalServerError
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie Deleted successfully"})
}

func (mc *MovieController) GetMovies(c *gin.Context) {
	movies, err := mc.MovieService.GetMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}

func (mc *MovieController) UpdateMovies(c *gin.Context) {
	var updatedMovie models.Movie
	if err := c.ShouldBindJSON(&updatedMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedMovie.MovieId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "movieId is required"})
		return
	}

	movie, err := mc.MovieService.UpdateMovies(&updatedMovie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

func (mc *MovieController) GetMovieById(c *gin.Context) {
	movieId := c.Query("movieId")
	if movieId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "movieId is required"})
		return
	}

	movie, err := mc.MovieService.GetMovieById(movieId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}
