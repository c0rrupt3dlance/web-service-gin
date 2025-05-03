package mocks

import (
	"github.com/c0rrupt3dlance/web-service-gin/models"
	"github.com/stretchr/testify/mock"
)

type AlbumRepoMock struct {
	mock.Mock
}

func (m *AlbumRepoMock) GetByID(id int64) (*models.Album, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Album), args.Error(1)
}

func (m *AlbumRepoMock) Create(album models.Album) (models.Album, error) {
	args := m.Called(album)
	return args.Get(0).(models.Album), args.Error(1)
}

func (m *AlbumRepoMock) GetAll() (*[]models.Album, error) {
	args := m.Called()
	return args.Get(0).(*[]models.Album), args.Error(1)
}

func (m *AlbumRepoMock) Remove(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *AlbumRepoMock) Update(id int64, album *models.Album) error {
	args := m.Called(id, album)
	return args.Error(1)
}
