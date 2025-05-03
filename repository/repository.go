package repository

import (
	"github.com/c0rrupt3dlance/web-service-gin/models"
)

type AlbumRepository interface {
	Create(album models.Album) (models.Album, error)
	GetAll() (*[]models.Album, error)
	GetByID(id int64) (*models.Album, error)
	Update(id int64, album *models.Album) error
	Remove(id int64) error
}
