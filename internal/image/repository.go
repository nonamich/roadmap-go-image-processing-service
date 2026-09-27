package image

import (
	"database/sql"
	"errors"
	"nonamich/image-processing-service/internal/database"
)

func FindImageByUuid(uuid string) (Image, error) {
	image := Image{}
	err := database.DB.QueryRow(
		`SELECT ID, username, password FROM users WHERE uuid = ? LIMIT 1`,
		uuid,
	).Scan(&image.ID, &image.Imagename, &image.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return Image{}, nil
	}

	if err != nil {
		return Image{}, err
	}

	return image, nil
}

func SaveImage(username string, password string) (Image, error) {
	result, err := database.DB.Exec(
		`
			INSERT INTO users (username, password)
			VALUES (?, ?)
		`,
		username,
		password,
	)

	if err != nil {
		return Image{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return Image{}, err
	}

	return Image{
		ID:        id,
		Imagename: username,
		Password:  password,
	}, nil
}
