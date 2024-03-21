package main

import (
	"errors"
	"fmt"
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
		session := getGothSession(c.Request)
		token, ok := session.Values["jwt_token"].(string)
		tokenString = token
		if !ok || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
			c.Abort()
			return
		}
	}

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

	session := getGothSession(c.Request)
	// Armazene a string na sessão
	session.Values["jwt_token"] = tokenString

	// Salve a sessão
	err2 := session.Save(c.Request, c.Writer)
	if err2 != nil {
		// Lidar com o erro, se houver
		fmt.Println("Erro ao salvar a sessão:", err)
		return
	}

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
	authorizationHeader := c.Request.Header.Get("Authorization")
	if authorizationHeader == "" {
		return ""
	}
	// O cabeçalho de autorização deve estar no formato "Bearer token"
	tokenString := strings.Split(authorizationHeader, " ")
	if len(tokenString) != 2 || tokenString[0] != "Bearer" {
		return ""
	}

	return tokenString[1]
}
