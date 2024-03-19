package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	router.GET("/auth/callback", func(c *gin.Context) { providerCallback(c) })
	router.GET("/logout", func(c *gin.Context) { oauthLogout(c) })
	router.GET("/auth", func(c *gin.Context) { authProvider(c) })
	router.GET("/", func(c *gin.Context) { getTemplate(c, pindex) })

	return router
}
