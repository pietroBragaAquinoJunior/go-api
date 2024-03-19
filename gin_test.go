package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	//GORM
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	db := gormSetup()
	router = ginSetup(db)
	code := m.Run()
	os.Exit(code)
}

func TestGetAlbums(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/albums", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestPostAlbum(t *testing.T) {
	w := httptest.NewRecorder()
	jsonData := `{"Title": "TEST3", "Artist": "Piêtro Braga", "Price": 39.99}`
	req, _ := http.NewRequest("POST", "/albums", strings.NewReader(jsonData))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
