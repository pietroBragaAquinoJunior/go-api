package main

import (
	"gorm.io/gorm"
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
}
