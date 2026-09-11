package main

import (
	"html/template"
	"net/http"
)

//var tmpl = template.Must( template.ParseFiles("templates/index.html","templates/part.html", "templates/edit_part.html", "templates/configuration.html", "templates/edit_configuration.html"))
var partsTemplate = template.Must( template.ParseFiles("templates/parts.html") )
var edit_partTemplate = template.Must( template.ParseFiles("templates/edit_part.html"))
var configurationTemplate = template.Must( template.ParseFiles( "templates/configurations.html"))
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
    Parts []PartData
    InternalId int32
}

func partsHandler( writer http.ResponseWriter, request *http.Request ) {
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
	// TODO: perform an SQL query and map the result to PartData array
	parts := []PartData{ examplePartData, examplePartData2 }
	partsTemplate.Execute(writer, parts)
}

func editPartHandler( writer http.ResponseWriter, request *http.Request ) {
	examplePartData := PartData {
	    Brand : "Shimano",
	    Model : "FD-A070",
	    Description : "black, max tooth 50",
	    Count : 1,
	    Price : 4000,
	    Image : "",
	    InternalId : 1,
	}
	// TODO: get ID from URL
	// if empty then all fields shall be empty and a STORE operation is to be performed
	// if not empty then query all fields (is this redundant? maybe recevie everything from query?)
	//	in this case an update need to be performed
	edit_partTemplate.Execute(writer, examplePartData)
}

func configurationHandler( writer http.ResponseWriter, request *http.Request ) {

	examplePartData := PartData {
	    Brand : "Shimano",
	    Model : "FD-A070",
	    Description : "black, max tooth 50",
	    Count : 1,
	    Price : 4000,
	    Image : "",
	    InternalId : 1,
	}
	exampleConfigurationData := ConfigurationData {
		Name : "Csepel Rapid",
		Description : "flatbar, 2x8, integrated crankset, 50/34, 35-622 Schwalbe",
		Price : 162000,
		AdditionalPrice : 0,
		AdditionalPriceDesc : "",
		Parts : []PartData{ examplePartData },
		InternalId : 1,
	}
	configurationTemplate.Execute( writer, []ConfigurationData { exampleConfigurationData } )
	// TODO: get ID from URL
	// if empty then all fields shall be empty and a STORE operation is to be performed
	// if not empty then query all fields (is this redundant? maybe recevie everything from query?)
	//	in this case an update need to be performed
//	configurationTemplate.Execute(writer, examplePartData)
}

func main() {


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	http.HandleFunc("/templates/parts.html", partsHandler)
	http.HandleFunc("/templates/edit_part.html", editPartHandler)
	http.HandleFunc("/templates/configurations.html", configurationHandler)
	http.ListenAndServe(":8080", nil)
}
