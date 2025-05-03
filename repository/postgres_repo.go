package repository

import (
	"database/sql"
	"errors"
	"github.com/c0rrupt3dlance/web-service-gin/errs"
	"github.com/c0rrupt3dlance/web-service-gin/models"
)

type PostgresAlbumRepo struct {
	DB *sql.DB
}

func (r *PostgresAlbumRepo) GetAll() (*[]models.Album, error) {
	rows, err := r.DB.Query("SELECT * FROM albums")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var album models.Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	return &albums, nil
}

func (r *PostgresAlbumRepo) GetByID(id int64) (*models.Album, error) {
	var album models.Album
	row := r.DB.QueryRow("SELECT * FROM albums WHERE id = $1", id)
	err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		} else {
			return &album, err
		}
	}
	return &album, nil
}

func (r *PostgresAlbumRepo) Create(album models.Album) (models.Album, error) {
	err := r.DB.QueryRow("INSERT INTO albums (title, artist, price) VALUES ($1, $2, $3)",
		album.Title, album.Artist, album.Price).Scan(&album.ID)
	if err != nil {
		return album, err
	}
	return album, nil
}
func (r *PostgresAlbumRepo) Remove(id int64) error {
	_, err := r.DB.Exec("DELETE FROM albums WHERE id = $1", id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errs.ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (r *PostgresAlbumRepo) Update(id int64, album *models.Album) error {
	result, err := r.DB.Exec("update albums set title = $1, artist = $2, price = $3 where id = $4", album.Title, album.Artist, album.Price, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
