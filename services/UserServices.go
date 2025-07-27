package services

import (
	"fmt"
	"movie/helpers"
	"movie/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserService struct {
	DB *sqlx.DB
}

func NewuserService(db *sqlx.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (us *UserService) Login(creds *models.Credentials) (*models.User, string, string, error) {
	var user models.User
	query := `
	SELECT userid, name, email, password, role FROM users WHERE email = $1
	`
	// Get User
	err := us.DB.Get(&user, query, creds.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to fetch user: %w", err)
	}

	// Compare passwords
	ok := !helpers.CheckPasswords(user.Password, creds.Password)
	if !ok {
		return nil, "", "", fmt.Errorf("passwords don't match")
	}

	accessToken, err := helpers.GenerateAccessToken(user.UserId.String(), user.Email, user.Role)
	if err != nil {
		return nil, "", "", fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.UserId.String(), user.Email, user.Role)
	if err != nil {
		return nil, "", "", fmt.Errorf("error generating refresh Token: %w", err)
	}

	return &user, accessToken, refreshToken, nil
}

func (us *UserService) Signup(user *models.User) (*models.User, string, string, error) {
	// checks if entered email is used
	emailAvailable, err := helpers.IsEmailAvailable(us.DB, user.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("could not check email availability: %w", err)
	}
	if !emailAvailable {
		return nil, "", "", fmt.Errorf("email already taken")
	}

	// generate new userId format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	user.UserId = uuid.New()

	user.Role = "user" // default is user, only admin can promote

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return nil, "", "", fmt.Errorf("could not hash password: %w", err)
	}

	user.Password = hashedPassword

	// Insert into DB and return inserted user
	query := `
		INSERT INTO users (userid, name, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING userid, name, email, password, role, createdat, updatedat
	`
	err = us.DB.QueryRow(query, user.UserId, user.Name, user.Email, user.Password, user.Role).
		Scan(&user.UserId, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, "", "", fmt.Errorf("error inserting user: %w", err)
	}

	// Generate tokens
	accessToken, err := helpers.GenerateAccessToken(user.UserId.String(), user.Email, user.Role)
	if err != nil {
		return nil, "", "", fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.UserId.String(), user.Email, user.Role)
	if err != nil {
		return nil, "", "", fmt.Errorf("error generating refresh token: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

func (us *UserService) PromoteToAdmin(userId string) error {
	var currentRole string
	checkQuery := `SELECT role FROM users WHERE userid = $1`
	err := us.DB.Get(&currentRole, checkQuery, userId)
	if err != nil {
		return fmt.Errorf("user not found or error checking user: %w", err)
	}

	if currentRole == "admin" {
		return fmt.Errorf("user is already an admin")
	}

	updateQuery := `
		UPDATE users
		SET role = $1
		WHERE userid = $2
	`
	_, err = us.DB.Exec(updateQuery, "admin", userId)
	if err != nil {
		return fmt.Errorf("error promoting to admin: %w", err)
	}

	return nil
}
