package main

import (
	"html/template"
	"net/http"
)

var partsTemplate = template.Must( template.ParseFiles("templates/parts.html") )
var edit_partTemplate = template.Must( template.ParseFiles("templates/edit_part.html"))
var configurationTemplate = template.Must( template.ParseFiles( "templates/configurations.html"))
var edit_configurationTemplate =  template.Must( template.ParseFiles("templates/edit_configuration.html"))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	http.HandleFunc("/templates/parts.html", partsHandler)
	http.HandleFunc("/templates/edit_part.html", editPartHandler)
	http.HandleFunc("/templates/configurations.html", configurationHandler)
	http.ListenAndServe(":8080", nil)
}
