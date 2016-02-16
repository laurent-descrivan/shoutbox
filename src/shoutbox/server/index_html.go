package server

import (
	"fmt"
	"html/template"
)

var INDEX_TEMPLATE *template.Template

func init() {
	tmpl, err := template.New("index").Parse(_INDEX_TEMPLATE_STR)
	if err != nil {
		panic(fmt.Sprintf("Error parsing template: %s", err.Error()))
	}
	INDEX_TEMPLATE = tmpl
}

const _INDEX_TEMPLATE_STR = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Shoutbox customization page !!!1</title>
	<style type="text/css">
		#text {
			width: 600px;
			height:400px;
		}
	</style>
</head>
<body>
	<h1>Shoutbox customization page !!!1</h1>
	<form method="POST" action="post">
		<textarea id="text" name="text">{{.txt}}</textarea>
		<br/>
		<input type="submit"/>
	</form>
</body>
</html>`
