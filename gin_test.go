package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func fazerLogin() string {
	loginData := map[string]string{
		"usuario": "pipiboy",
		"senha":   "172983456",
	}

	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		fmt.Println("Erro ao codificar dados de login:", err)
		return ""
	}

	w := httptest.NewRecorder()
	reqAuth, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	reqAuth.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqAuth)
	if w.Code != http.StatusOK {
		fmt.Printf("Falha ao fazer login. Código de status: %d\n", w.Code)
		return ""
	}

	var response loginResponse

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return ""
	}
	return response.Token
}

func TestGetAlbums(t *testing.T) {
	token := fazerLogin()

	req, _ := http.NewRequest("GET", "/api/albums", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPostAlbum(t *testing.T) {
	token := fazerLogin()

	w := httptest.NewRecorder()
	jsonData := `{"Title": "TEST3", "Artist": "Piêtro Braga", "Price": 39.99}`
	req, _ := http.NewRequest("POST", "/api/albums", strings.NewReader(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
