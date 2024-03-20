package main

import (
	"gorm.io/gorm"
	"time"
)

type Album struct {
	gorm.Model
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type User struct {
	gorm.Model
	Nome    string `json:"nome"`
	Usuario string `json:"usuario"`
	Senha   string `json:"senha"`

	// oauth
	OauthUserId string `json:"oauth_user_id"`
	Provedor     string `json:"provedor"`
	Email        string `json:"email"`
	NomeNick     string `json:"nome_nick"`
	Lugar        string `json:"lugar"`
	UrlAvatar    string `json:"url_avatar"`
	Descricao    string `json:"descricao"`
	AccessToken  string `json:"acess_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}
