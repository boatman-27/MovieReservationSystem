package controllers

import (
	"movie/models"
	"movie/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReservationController struct {
	ReservationService *services.ReservationService
}

func NewReservationServiceController(reservationService *services.ReservationService) *ReservationController {
	return &ReservationController{
		reservationService,
	}
}

func (rc *ReservationController) GetUpcomingReservations(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	reservations, err := rc.ReservationService.GetUpcomingEvents(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

func (rc *ReservationController) GetAllReservations(c *gin.Context) {
	reservations, err := rc.ReservationService.GetAllReservations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

func (rc *ReservationController) GetUserReservations(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	reservations, err := rc.ReservationService.GetAllUserReservations(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

func (rc *ReservationController) CancelReservation(c *gin.Context) {
	reservationId := c.Query("reservationId")
	if reservationId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reservationId is required"})
		return
	}

	err := rc.ReservationService.CancelReservation(reservationId)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case strings.Contains(err.Error(), "not found"):
			status = http.StatusNotFound
		case strings.Contains(err.Error(), "delete reservation"):
			status = http.StatusInternalServerError
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted successfully"})
}

func (rc *ReservationController) BookSeats(c *gin.Context) {
	var bookingData models.BookingData
	if err := c.ShouldBindJSON(&bookingData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if bookingData.ShowtimeId == "" || bookingData.UserId == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Need both showtimeId and userId"})
		return
	}

	if bookingData.Seats < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seats should be a positive integer"})
	}

	reservation, err := rc.ReservationService.BookSeats(&bookingData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservation": reservation})
}
