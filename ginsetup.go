package main

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

func ginSetup(db * gorm.DB)  * gin.Engine {
	router := gin.Default()
	
	// Rotas protegidas com JWT
	protected := router.Group("/api")
	protected.Use(authMiddleware)

	protected.GET("/albums", func(c *gin.Context) { getAlbums(c, db) })
	protected.GET("/albums/:id", func(c *gin.Context) { getAlbumByID(c, db) })
	protected.POST("/albums", func(c *gin.Context) { postAlbum(c, db) })
	router.POST("/login", func(c *gin.Context) { login(c, db) })
	return router
}