package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/c0rrupt3dlance/web-service-gin/errs"
	"github.com/c0rrupt3dlance/web-service-gin/handlers"
	"github.com/c0rrupt3dlance/web-service-gin/mocks"
	"github.com/c0rrupt3dlance/web-service-gin/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func setupMockGetByID(repo *mocks.AlbumRepoMock, id int64, result *models.Album, err error) {
	repo.ExpectedCalls = nil
	repo.On("GetByID", id).Return(result, err)
}

func setupMockGetAll(repo *mocks.AlbumRepoMock, result *[]models.Album, err error) {
	repo.ExpectedCalls = nil
	repo.On("GetAll").Return(result, err)
}

func TestGetAlbumById(t *testing.T) {
	repo := new(mocks.AlbumRepoMock)
	handler := &handlers.AlbumHandler{Repo: repo}
	router := gin.Default()
	router.GET("/album/:id", handler.GetAlbumByID)
	tests := []struct {
		name         string
		id           string
		mockResult   *models.Album
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Found Album",
			id:           "1",
			mockResult:   &models.Album{ID: 1, Title: "Album 1", Artist: "Artist 1", Price: 1.0},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"album":{"id":1,"title":"Album 1","artist":"Artist 1","price":1}}`,
		},
		{
			name:         "Album not found",
			id:           "2",
			mockResult:   nil,
			mockError:    errs.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"message":"album not found"}`,
		},
		{
			name:         "Invalid ID",
			id:           "abc",
			mockResult:   nil,
			mockError:    errs.ErrInvalidInput,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"invalid input"}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.id != "abc" {
				id, _ := strconv.ParseInt(tc.id, 10, 64)
				setupMockGetByID(repo, id, tc.mockResult, tc.mockError)
			}
			req := httptest.NewRequest(http.MethodGet, "/album/"+tc.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())

			repo.AssertExpectations(t)
		})
	}
}

func TestGetAllAlbums(t *testing.T) {
	repo := new(mocks.AlbumRepoMock)
	handler := &handlers.AlbumHandler{Repo: repo}
	router := gin.Default()
	router.GET("/album", handler.GetAlbums)
	tests := []struct {
		name         string
		mockResult   *[]models.Album
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name: "Album list is not empty",
			mockResult: &[]models.Album{
				models.Album{
					ID:     1,
					Title:  "Album 1",
					Artist: "Artist 1",
					Price:  1.0,
				},
				models.Album{
					ID:     2,
					Title:  "Album 2",
					Artist: "Artist 2",
					Price:  1.0,
				},
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"albums":[{"id":1,"title":"Album 1","artist":"Artist 1","price":1},{"id":2,"title":"Album 2","artist":"Artist 2","price":1}]}`,
		},
		{
			name:         "Album list is empty",
			mockResult:   &[]models.Album{},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"albums":[]}`,
		},
		{
			name:         "Repository error",
			mockResult:   &[]models.Album{},
			mockError:    errs.ErrInternalServerError,
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"internal server error"}`,
		},
		{
			name:         "Unexpected nil",
			mockResult:   nil,
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"albums":[]}`,
		},
	}
	for i, tc := range tests {
		log.Printf("%v test %v", tc.name, i)
		t.Run(tc.name, func(t *testing.T) {
			setupMockGetAll(repo, tc.mockResult, tc.mockError)
			req := httptest.NewRequest(http.MethodGet, "/album", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())

			repo.AssertExpectations(t)
		})
	}
}

func setupMockUpdate(repo *mocks.AlbumRepoMock, id int64, result *models.Album, err error) {
	repo.ExpectedCalls = nil
	repo.On("Update", id, result).Return(result, err)
}

func TestUpdateAlbum(t *testing.T) {
	repo := new(mocks.AlbumRepoMock)
	handler := &handlers.AlbumHandler{Repo: repo}
	router := gin.Default()
	router.PUT("/album/:id", handler.UpdateAlbum)
	test := []struct {
		name         string
		id           string
		mockResult   *models.Album
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name: "Successful update",
			id:   "1",
			mockResult: &models.Album{
				ID:     1,
				Title:  "Album 1",
				Artist: "Artist 1",
				Price:  1.0,
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"album":{"id":1,"title":"Album 1","artist":"Artist 1","price":1},` +
				`"message":"Album was updated"}`,
		},
	}

	for _, tc := range test {
		log.Printf("~~~~~%v test %v~~~~~", tc.name, tc.id)
		t.Run(tc.name, func(t *testing.T) {
			setupMockUpdate(repo, tc.mockResult.ID, tc.mockResult, tc.mockError)
			body, _ := json.Marshal(tc.mockResult)
			req := httptest.NewRequest(http.MethodPut, "/album/"+tc.id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
			repo.AssertExpectations(t)
		})
	}
}

func setupMockDelete(repo *mocks.AlbumRepoMock, id int64, err error) {
	repo.ExpectedCalls = nil
	repo.On("Remove", id).Return(err)
}

func TestDeleteAlbum(t *testing.T) {
	repo := new(mocks.AlbumRepoMock)
	handler := &handlers.AlbumHandler{Repo: repo}
	router := gin.Default()
	router.DELETE("/album/:id", handler.DeleteAlbum)
	test := []struct {
		name         string
		id           string
		mockResult   []*models.Album
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name: "Successful delete",
			id:   "1",
			mockResult: []*models.Album{
				{
					ID:     1,
					Title:  "Album 1",
					Artist: "Artist 1",
					Price:  1.0,
				},
				{
					ID:     2,
					Title:  "Album 2",
					Artist: "Artist 2",
					Price:  1.0,
				},
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1}`,
		},
		{
			name:         "Album list is empty",
			id:           "2",
			mockResult:   []*models.Album{},
			mockError:    sql.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"message":"invalid input"}`,
		},
		{
			name: "Album not found",
			id:   "3",
			mockResult: []*models.Album{
				{
					ID:     1,
					Title:  "Album 1",
					Artist: "Artist 1",
					Price:  1.0,
				},
			},
			mockError:    errs.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"message":"album not found"}`,
		},
	}
	for _, tc := range test {
		log.Printf("~~~~~%v test %v~~~~~", tc.name, tc.id)
		t.Run(tc.name, func(t *testing.T) {
			setupMockDelete(repo, 1, tc.mockError)
			req := httptest.NewRequest(http.MethodDelete, "/album/"+"1", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
			repo.AssertExpectations(t)
		})
	}
}
