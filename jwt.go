package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// sha-256
var secretKey = []byte(os.Getenv("TOKEN_JWT_SECRET"))

func authMiddleware(c *gin.Context) {
	// Extrai o token JWT do cabeçalho "Authorization"
	tokenString := extractTokenFromHeader(c)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
		c.Abort()
		return
	}

	// Analisa o token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		c.Abort()
		return
	}

	c.Next()
}

func login(c *gin.Context, db *gorm.DB) {
	var requestBody struct {
		Usuario string `json:"usuario"`
		Senha   string `json:"senha"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao ler corpo da requisição"})
		return
	}
	user, err := authenticateUser(requestBody.Usuario, requestBody.Senha, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	gerarERetornarTokenJwt(c, strconv.Itoa(int(user.ID)))

}

func gerarERetornarTokenJwt(c *gin.Context, idUsuario string) {
	// Se o usuário for autenticado com sucesso, gere um token JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = idUsuario // Aqui você pode usar o ID do usuário, por exemplo

	// Define o tempo de expiração do token
	expirationTime := time.Now().Add(24 * time.Hour) // Expira em 24 horas
	claims["exp"] = expirationTime.Unix()

	// Assina o token com a chave secreta
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar token"})
		return
	}
	// Retorna o token JWT para o cliente
	c.JSON(http.StatusOK, gin.H{"tokenJWT": tokenString})
}

func authenticateUser(usuario string, senha string, db *gorm.DB) (User, error) {
	var user User
	if err := db.Where("usuario = ?", usuario).First(&user).Error; err != nil {
		return User{}, errors.New("usuário não encontrado")
	}
	if user.Senha != senha {
		return User{}, errors.New("senha incorreta")
	}
	return user, nil
}

func extractTokenFromHeader(c *gin.Context) string {
	// Obtém o cabeçalho "Authorization" da requisição
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		return ""
	}

	// Divide o cabeçalho "Authorization" para obter o token JWT
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}
