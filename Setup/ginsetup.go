package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ginSetup(db *gorm.DB, pindex *ProviderIndex) *gin.Engine {
	router := gin.Default()

	// Rotas protegidas com JWT
	protected := router.Group("/api")
	protected.Use(authMiddleware)

	protected.GET("/albums", func(c *gin.Context) { getAlbums(c, db) })
	protected.GET("/albums/:id", func(c *gin.Context) { getAlbumByID(c, db) })
	protected.POST("/albums", func(c *gin.Context) { postAlbum(c, db) })

	// JWT
	router.POST("/login", func(c *gin.Context) { login(c, db) })

	//OAUTH
	router.GET("/auth/callback", func(c *gin.Context) { providerCallback(c, db) })
	router.GET("/logout", func(c *gin.Context) { oauthLogout(c) })
	router.GET("/auth", func(c *gin.Context) { authProvider(c) })
	router.GET("/", func(c *gin.Context) { getTemplate(c, pindex) })

	return router
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
