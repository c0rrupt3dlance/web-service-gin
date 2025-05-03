package handlers

import (
	"database/sql"
	"errors"
	"github.com/c0rrupt3dlance/web-service-gin/errs"
	"github.com/c0rrupt3dlance/web-service-gin/models"
	"github.com/c0rrupt3dlance/web-service-gin/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type AlbumHandler struct {
	Repo repository.AlbumRepository
}

func (h *AlbumHandler) GetAlbums(c *gin.Context) {
	albums, err := h.Repo.GetAll()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errs.ErrInternalServerError.Error()})
		return
	}
	if albums == nil {
		albums = &[]models.Album{}
	}
	c.JSON(http.StatusOK, gin.H{"albums": albums})
}

func (h *AlbumHandler) GetAlbumByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}
	album, err := h.Repo.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
		case errors.Is(err, errs.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"album": album})
}

func (h *AlbumHandler) AddAlbum(c *gin.Context) {
	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.Repo.Create(album)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errs.ErrInvalidInput})
		return
	}
	err = h.Repo.Remove(id)
	if err != nil {
		switch err {
		case errs.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
			return

		case sql.ErrNoRows:
			c.JSON(http.StatusNotFound, gin.H{"message": "invalid input"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {
	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		log.Println(album)
		c.JSON(http.StatusBadRequest, gin.H{"message": errs.ErrInternalServerError.Error()})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	album.ID = id
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errs.ErrInvalidInput})
		return
	}
	err = h.Repo.Update(id, &album)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": errs.ErrNotFound.Error()})
		return
	} else if err != nil {
		log.Println(album)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errs.ErrInternalServerError.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Album was updated", "album": album})

}
