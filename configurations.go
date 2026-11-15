package main

import (
	"database/sql"
	"net/http"
	"strconv"
)

type ConfigurationData struct {
	Name                string
	Description         string
	Price               int32
	AdditionalPrice     int32
	AdditionalPriceDesc string
	Parts               []PartData
	Id                  string
}

/*****************
*      HTML      *
******************/

// Helper function that convert GET method params to ConfigurationData
func getConfigurationDataFromRequest(request *http.Request, db *sql.DB) ConfigurationData {
	var values = request.URL.Query()
	var name = values.Get("name")
	var description = values.Get("description")
	var price = values.Get("price")
	var additionalPrice = values.Get("additionalPrice")
	var additionalPriceDescription = values.Get("additionalPriceDescription")
	var id = values.Get("id")
	var n int64
	var p int32 = 0
	var err error
	n, err = strconv.ParseInt(price, 10, 32)
	if err == nil {
		p = int32(n)
	}
	var a int32 = 0
	n, err = strconv.ParseInt(additionalPrice, 10, 32)
	if err == nil {
		a = int32(n)
	}

	var parts []PartData
	var partIds = values["part"]
	for partId := range partIds {
		parts = append(parts, getPart(db, partIds[partId]))
	}

	return ConfigurationData{name, description, p, a, additionalPriceDescription, parts, id}
}

// Main function for displaying configurations in configurations.html.
func configurationHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	configurations := getConfigurations(db)
	configurationTemplate.Execute(writer, configurations)
}

