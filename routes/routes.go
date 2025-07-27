package routes

import (
	controllers "movie/controller"
	"movie/middlewares"
	"movie/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(router *gin.Engine, db *sqlx.DB) {
	// Services
	userService := services.NewuserService(db)
	movieService := services.NewMovieService(db)
	showtimeService := services.NewShowtimeService(db)
	reservationService := services.NewReservationService(db)

	// Controllers
	userController := controllers.NewUserController(userService)
	movieController := controllers.NewMovieController(movieService)
	showtimeContoller := controllers.NewShowtimeController(showtimeService)
	reservationController := controllers.NewReservationServiceController(reservationService)

	// Authentication Routes
	accountRoutes := router.Group("/account")
	{
		accountRoutes.POST("/login", userController.Login)
		accountRoutes.POST("/signup", userController.Signup)
	}

	// Protected Routes
	protected := router.Group("/protected")
	protected.Use(middlewares.RequireAuth)
	{
		protected.GET("/movies", movieController.GetMovies)
		protected.POST("/get-movie-byid", movieController.GetMovieById)
		protected.POST("/get-showtime-and-movie", showtimeContoller.GetShowtimeAndMovie)
		protected.POST("/get-seatsinfo", showtimeContoller.CheckAvailableSeats)
		protected.POST("/upcoming-reservations", reservationController.GetUpcomingReservations)
		protected.POST("/cancel-reservation", reservationController.CancelReservation)
		protected.POST("/book-seats", reservationController.BookSeats)
	}

	// Admin Routes
	admin := router.Group("/admin")
	admin.Use(middlewares.AdminAuth)
	{
		admin.POST("/promote", userController.PromoteToAdmin)
		admin.POST("/add-movie", movieController.AddMovie)
		admin.POST("/delete-movie", movieController.DeleteMovie)
		admin.PATCH("/update-movie", movieController.UpdateMovies)
		admin.POST("/add-showtime", showtimeContoller.AddShowtimes)
		admin.POST("/delete-showtime", showtimeContoller.DeleteShowtime)
		admin.PATCH("/update-showtime", showtimeContoller.UpdateShowtime)
		admin.GET("/all-reservations", reservationController.GetAllReservations)
		admin.POST("/user-reservations", reservationController.GetUserReservations)
	}
}
