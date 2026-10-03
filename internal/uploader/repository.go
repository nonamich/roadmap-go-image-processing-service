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
		`SELECT ID, uuid, mime, metadata FROM files WHERE uuid = ? LIMIT 1`,
		uuid,
	).Scan(&file.ID, &file.Uuid, &file.Mime, &metadataJSON)

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
