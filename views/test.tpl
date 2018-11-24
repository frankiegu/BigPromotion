<html>
  <head>
      <title>{{.Title}}</title>
  </head>
  <body>
  {{.Title}}
  {{range .Users}}
    {{.Name}}{{$.len}}<br/>
  {{end}}
  </body>
</html>