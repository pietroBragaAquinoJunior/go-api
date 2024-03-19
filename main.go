package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func main() {
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

	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

	db := gormSetup()

	r := ginSetup(db, providerIndex)

	r.Run(":8080")
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

func getAlbums(c *gin.Context, db *gorm.DB) {
	var albums []Album
	db.Find(&albums)
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbum(c *gin.Context, db *gorm.DB) {
	var album Album
	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&album)
	c.JSON(http.StatusOK, album)
}

func getAlbumByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var album Album
	if err := db.First(&album, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Album não encontrado"})
		return
	}
	c.JSON(http.StatusOK, album)
}

// Função de retorno de chamada para o provedor Discord
func providerCallback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	t, _ := template.New("user").Parse(userTemplate)



	
	t.Execute(c.Writer, user)
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
