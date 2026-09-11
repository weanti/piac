package main

import (
	"fmt"
	"html/template"
	"net/http"
)

//var tmpl = template.Must( template.ParseFiles("templates/index.html","templates/part.html", "templates/edit_part.html", "templates/configuration.html", "templates/edit_configuration.html"))
var partTemplate = template.Must( template.ParseFiles("templates/part.html") )
var edit_partTemplate = template.Must( template.ParseFiles("templates/edit_part.html"))
var configurationTemplate = template.Must( template.ParseFiles( "templates/configuration.html"))
var edit_configurationTemplate =  template.Must( template.ParseFiles("templates/edit_configuration.html"))

type PartData struct {
    Brand string
    Model string
    Description string
    Count int16
    Price int32
    Image string
}

type ConfigurationData struct {
    Name string
    Description string
    Price int32
    AdditionalPrice int32
    AdditionalPriceDesc string
}

func partHandler( writer http.ResponseWriter, request *http.Request ) {
	examplePartData := PartData {
	    Brand : "Shimano",
	    Model : "FD-A070",
	    Description : "black, max tooth 50",
	    Count : 1,
	    Price : 4000,
	    Image : "",
	}
	e := partTemplate.Execute(writer, examplePartData)
	fmt.Println("partHandler called")
	if e != nil {
		fmt.Println( e )
	}
}

func main() {


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	http.HandleFunc("/parts", partHandler)
//	http.handleFunc("/configurations", configurationHandler)
	http.ListenAndServe(":8080", nil)
}
