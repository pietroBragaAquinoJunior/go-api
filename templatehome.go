package main

var indexTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        * {
            margin:0;
			padding:0
			box-sizing: border-box;
			font-family: "Courier New", Courier, monospace;
        }
		body{
			font-family: "Courier New", Courier, monospace;
			background-color: #EAEAEA;
		}
		.container{
			display:flex;
			justify-content:center;
			align-items:center;
			flex-direction:column;
		}
		.anchor{
			color:blue;
			text-decoration:none;
		}
		.provider{
			color: black;
			text-decoration:none;
			border:1px solid black;
			padding-top:5px;
			padding-bottom:5px;
			font-weight:700;
			font-size:15px;
			background:#F6D2A2;
			margin-top:10px;
			display:flex;
			justify-content:center;
			align-items:center;
		  }

		  .provider a {
			text-decoration:none;
			color:black;
		  }
		  .provider input {
			border:none;
			background:none;
		  }
		  .colorido{
			border-radius:10px;
			padding-top:7px;
			padding-bottom:7px;
			padding-right:10px;
			padding-left:10px;
			background:#6AD7E5;
			border:3px solid black;
		  }
		  .bom-espaco{
			padding-top:30px;
			padding-bottom:30px;
			padding-left:50px;
			padding-right:50px;
		  }
		  .fonte-padrao{
			font-size:20px;
			font-weight:bold;
		  }
    </style>
</head>
<body>
    <div class="container">
			<img width="800px" src=/public/golang.png alt="Descrição da Imagem" />
		<div style="margin-bottom:20px;">
			<a class="anchor" style="font-size:20px;" href="https://github.com/pietroBragaAquinoJunior" target="_blank">Criado por Piêtro Braga Aquino Júnior</a>
		</div>

		{{if not .Autenticado}}
		<div class="container colorido bom-espaco" >
			<h2 style="margin-bottom:30px">Entrar com login e senha</h2>
			<form style="width:100%" action="/login" method="post">
				<label class="fonte-padrao" for="usuario">Usuário:</label><br>
				<input style="margin-bottom:10px; width:100%" type="text" id="usuario" name="usuario" required><br>
		
				<label class="fonte-padrao" for="senha">Senha:</label><br>
				<input style="margin-bottom:10px; width:100%" type="password" id="senha" name="senha" required><br><br>
		
				<div class="provider" style="width:100%">
					<input class="fonte-padrao" type="submit" value="Login">
				</div>
			</form>
			<div class="provider" style="width:100%">
			{{range $key, $value := .Providers}}
				<a class="fonte-padrao" href="/auth?provider={{$value}}">Entrar com {{index $.ProvidersMap $value}}</a>
			{{end}}
			</div>
		</div>
		{{end}}
		{{if .Autenticado}}
		<div>
		  <p>Você foi autenticado com sucesso!</p>
		  <a href="/api/albums">Consultar Api Get Albums</a>
		</div>
		{{end}}
    </div>
</body>
</html>
`

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
