package uploader

import (
	"database/sql"
	"encoding/json"
	"errors"
	"nonamich/image-processing-service/internal/database"
	"nonamich/image-processing-service/internal/user"

	"github.com/google/uuid"
)

func FindFileByUuid(uuid string) (UploadedFile, error) {
	var metadataJSON []byte
	var file UploadedFile
	err := database.DB.QueryRow(
		`SELECT ID, uuid, mime, user_id, metadata FROM files WHERE uuid = ? LIMIT 1`,
		uuid,
	).Scan(&file.ID, &file.Uuid, &file.Mime, &file.UserId, &metadataJSON)

	if errors.Is(err, sql.ErrNoRows) {
		return UploadedFile{}, nil
	}

	if err != nil {
		return UploadedFile{}, err
	}

	err = json.Unmarshal(metadataJSON, &file.Metadata)

	if err != nil {
		return UploadedFile{}, err
	}

	return file, nil
}

func SaveFile(user user.User, mime string, metadata FileMetadata) (UploadedFile, error) {
	fileUuid := uuid.New()
	metadataJSON, err := json.Marshal(metadata)

	if err != nil {
		return UploadedFile{}, err
	}

	result, err := database.DB.Exec(
		`
			INSERT INTO files (uuid, mime, user_id, metadata)
			VALUES (?, ?, ?, ?)
		`,
		fileUuid,
		mime,
		user.ID,
		metadataJSON,
	)

	if err != nil {
		return UploadedFile{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return UploadedFile{}, err
	}

	return UploadedFile{
		ID:       id,
		UserId:   user.ID,
		Uuid:     fileUuid.String(),
		Mime:     mime,
		Metadata: metadata,
	}, nil
}

func FindFilesByUserId(userId int64, page int, limit int) ([]UploadedFile, int, error) {
	total := 0
	limit = max(limit, 1)
	limit = min(limit, 50)
	page = max(page, 1)
	files := []UploadedFile{}
	offset := limit*page - limit
	rows, err := database.DB.Query(
		`
			SELECT ID, uuid, mime, user_id, metadata
			FROM files
			WHERE user_id = ?
			ORDER BY ID DESC
			LIMIT ?
			OFFSET ?
		`,
		userId,
		limit,
		offset,
	)

	if err := rows.Err(); err != nil {
		return files, total, err
	}

	defer rows.Close()

	if err != nil {
		return files, total, err
	}

	for rows.Next() {
		var file UploadedFile
		var metadataJSON []byte

		err := rows.Scan(
			&file.ID,
			&file.Uuid,
			&file.Mime,
			&file.UserId,
			&metadataJSON,
		)

		if err != nil {
			return files, total, err
		}

		err = json.Unmarshal(metadataJSON, &file.Metadata)

		if err != nil {
			return files, total, err
		}

		files = append(files, file)
	}

	database.DB.QueryRow(
		`
			SELECT COUNT(ID)
			FROM files
			WHERE user_id = ?
		`,
		userId,
	).Scan(&total)

	return files, total, nil
}
