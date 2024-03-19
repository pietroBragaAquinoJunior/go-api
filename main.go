package main

import (
	//GIN
	"net/http"

	"github.com/gin-gonic/gin"

	//GORM
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func main() {
	//GORM
	dsn := "root:172983456@tcp(localhost:3306)/goapi?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&Album{})
	//	db.Create(&Album{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99})
	//	db.Create(&Album{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99})
	//	db.Create(&Album{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99})
	// GIN
	router := gin.Default()
	router.GET("/albums", func(c *gin.Context) { getAlbums(c, db) })
	router.GET("/albums/:id", func(c *gin.Context) { getAlbumByID(c, db) })
	router.POST("/albums", func(c *gin.Context) { postAlbum(c, db) })
	router.Run("localhost:8080")
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
