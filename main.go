package main

import (
	"database/sql"
	"github.com/c0rrupt3dlance/web-service-gin/database"
	"github.com/c0rrupt3dlance/web-service-gin/handlers"
	"github.com/c0rrupt3dlance/web-service-gin/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	//"net/http"
)

var DB *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// f
func main() {

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	repo := &repository.PostgresAlbumRepo{DB: db}
	handler := &handlers.AlbumHandler{Repo: repo}
	router := gin.Default()
	router.GET("/albums", handler.GetAlbums)
	router.POST("/albums", handler.AddAlbum)
	router.GET("/albums/:id", handler.GetAlbumByID)
	router.DELETE("/albums/:id", handler.DeleteAlbum)
	router.PUT("/albums/:id", handler.UpdateAlbum)
	router.Run("localhost:8080")
}
