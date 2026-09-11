package main

import (
	"html/template"
	"net/http"
)

//var tmpl = template.Must( template.ParseFiles("templates/index.html","templates/part.html", "templates/edit_part.html", "templates/configuration.html", "templates/edit_configuration.html"))
var partsTemplate = template.Must( template.ParseFiles("templates/parts.html") )
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
    InternalId int32
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
	    InternalId : 1,
	}
	examplePartData2 := PartData {
	    Brand : "Sunrace",
	    Model : "RD90",
	    Description : "max tooth 28",
	    Count : 1,
	    Price : 4000,
	    Image : "",
	    InternalId : 2,
	}
	parts := []PartData{ examplePartData, examplePartData2 }
	partsTemplate.Execute(writer, parts)
}

func main() {


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	http.HandleFunc("/parts", partHandler)
//	http.handleFunc("/configurations", configurationHandler)
	http.ListenAndServe(":8080", nil)
}
