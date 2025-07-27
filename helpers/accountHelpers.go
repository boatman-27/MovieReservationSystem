package helpers

import "github.com/jmoiron/sqlx"

func IsEmailAvailable(db *sqlx.DB, email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil // true = available
}

func IsUserIdAvailable(db *sqlx.DB, userId string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil // true = available
}
