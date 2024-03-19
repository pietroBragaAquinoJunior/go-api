package main

import (
	//GIN
	"net/http"

	"github.com/gin-gonic/gin"

	//GORM
	"gorm.io/gorm"
)

func main() {
	db := gormSetup()
	r := ginSetup(db)
	r.Run("localhost:8080")
}

func getAlbums(c *gin.Context, db *gorm.DB) {
	var albums []Album
	db.Find(&albums)
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbum(c *gin.Context, db *gorm.DB) {
	var album Album
	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&album)
	c.JSON(http.StatusOK, album)
}

func getAlbumByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var album Album
	if err := db.First(&album, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Album não encontrado"})
		return
	}
	c.JSON(http.StatusOK, album)
}
