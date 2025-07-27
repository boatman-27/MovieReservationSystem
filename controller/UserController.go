package controllers

import (
	"movie/models"
	"movie/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService,
	}
}

func (uc *UserController) Login(c *gin.Context) {
	var credentials models.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, accessToken, refreshToken, err := uc.UserService.Login(&credentials)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie(
		"refreshToken", // cookie name
		refreshToken,   // value
		7*24*60*60,     // maxAge in seconds (7 days)
		"/",            // path
		"",             // domain (empty = current domain)
		false,          // secure (set true in production with HTTPS)
		true,           // httpOnly (can't be accessed by JS)
	)

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": accessToken,
	})
}

func (uc *UserController) Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdUser, accessToken, refreshToken, err := uc.UserService.Signup(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.SetCookie(
		"refreshToken", // cookie name
		refreshToken,   // value
		7*24*60*60,     // maxAge in seconds (7 days)
		"/",            // path
		"",             // domain (empty = current domain)
		false,          // secure (set true in production with HTTPS)
		true,           // httpOnly (can't be accessed by JS)
	)

	c.JSON(http.StatusOK, gin.H{
		"user":  createdUser,
		"token": accessToken,
	})
}

func (uc *UserController) PromoteToAdmin(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	err := uc.UserService.PromoteToAdmin(userId)
	if err != nil {
		status := http.StatusBadRequest

		switch {
		case strings.Contains(err.Error(), "not found"):
			status = http.StatusNotFound
		case strings.Contains(err.Error(), "already an admin"):
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}
