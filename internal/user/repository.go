package UserRepository

import (
	"database/sql"
	"errors"
	"nonamich/image-processing-service/internal/database"
)

type User struct {
	ID       int64
	Username string
	Password string
}

func FindUserByUsername(username string) (User, error) {
	user := User{}
	err := database.DB.QueryRow(
		`SELECT ID, username, password FROM users WHERE username = ? LIMIT 1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, nil
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func InsertUser(username string, password string) (User, error) {
	result, err := database.DB.Exec(
		`
			INSERT INTO users (username, password)
			VALUES (?, ?)
		`,
		username,
		password,
	)

	if err != nil {
		return User{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return User{}, err
	}

	return User{
		ID:       id,
		Username: username,
		Password: password,
	}, nil
}
