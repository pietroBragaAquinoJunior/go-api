package main

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

func ginSetup(db * gorm.DB)  * gin.Engine {
	router := gin.Default()
	router.GET("/albums", func(c *gin.Context) { getAlbums(c, db) })
	router.GET("/albums/:id", func(c *gin.Context) { getAlbumByID(c, db) })
	router.POST("/albums", func(c *gin.Context) { postAlbum(c, db) })
	return router
}