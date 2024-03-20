package main

import (
	"html/template"
	"net/http"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"gorm.io/gorm"
)

func gothSetup() *ProviderIndex {
	goth.UseProviders(
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), "http://localhost:8080/auth/callback?provider=discord", discord.ScopeIdentify, discord.ScopeEmail),
	)
	m := map[string]string{
		"discord": "Discord",
	}
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return &ProviderIndex{Providers: keys, ProvidersMap: m}
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth?provider={{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

var userTemplate = `
<p><a href="/logout?provider={{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

// Função de retorno de chamada para o provedor Discord
func providerCallback(c *gin.Context, db *gorm.DB) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	//t, _ := template.New("user").Parse(userTemplate)
	//t.Execute(c.Writer, user)

	// TENHO INFORMAÇÕES DO USUARIO AQUI

	testarTokenDiscordGerarJwt(c, user)

	salvarUsuarioOauth(c, user, db)

}

func salvarUsuarioOauth(c *gin.Context, gothUser goth.User, db *gorm.DB) {

	var usuarioEncontradoBanco User

	usuarioCriado := &User{
		OauthUserId:  gothUser.UserID,
		Nome:         gothUser.Name,
		Provedor:     gothUser.Provider,
		Email:        gothUser.Email,
		NomeNick:     gothUser.NickName,
		Lugar:        gothUser.Location,
		UrlAvatar:    gothUser.AvatarURL,
		Descricao:    gothUser.Description,
		AccessToken:  gothUser.AccessToken,
		ExpiresAt:    &gothUser.ExpiresAt,
		RefreshToken: gothUser.RefreshToken,
	}

	// Verifica se o usuário já existe no banco de dados
	if err := db.Where("oauth_user_id = ?", gothUser.UserID).First(&usuarioEncontradoBanco).Error; err != nil {
		db.Create(usuarioCriado)
	} else {
		usuarioCriado.ID = usuarioEncontradoBanco.ID
		usuarioCriado.CreatedAt = usuarioEncontradoBanco.CreatedAt
		db.Save(usuarioCriado)
	}
}

// Função para fazer logout do provedor OAuth
func oauthLogout(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// Função para autenticar o usuário com o provedor OAuth
func authProvider(c *gin.Context) {
	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		t, _ := template.New("user").Parse(userTemplate)
		t.Execute(c.Writer, gothUser)
	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

// Função para renderizar o template de login com os provedores disponíveis
func getTemplate(c *gin.Context, pindex *ProviderIndex) {
	t, _ := template.New("index").Parse(indexTemplate)
	t.Execute(c.Writer, pindex)
}

func testarTokenDiscordGerarJwt(c *gin.Context, user goth.User) {
	accessToken := user.AccessToken
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token de acesso não fornecido"})
		return
	}
	url := "https://discord.com/api/v9/users/@me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar requisição"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao fazer requisição"})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		// TOKEN VÁLIDO, GERAR JWT.
		gerarERetornarTokenJwt(c, user.UserID)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de acesso inválido"})
	}
}
