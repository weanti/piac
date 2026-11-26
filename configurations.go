package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"os"
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

// Helper function that convert GET method params to ConfigurationData
func getConfigurationDataFromRequest(request *http.Request, db *sql.DB) (ConfigurationData, error) {
	var values = request.URL.Query()
	var name = values.Get("name")
	var description = values.Get("description")
	var price = values.Get("price")
	var additionalPrice = values.Get("additional_price")
	var additionalPriceDescription = values.Get("additional_price_description")
	var id = values.Get("id")
	var n int64
	var p int32 = 0
	var err error
	n, err = strconv.ParseInt(price, 10, 32)
	if err != nil {
		return ConfigurationData{}, err
	}
	p = int32(n)
	var a int32 = 0
	n, err = strconv.ParseInt(additionalPrice, 10, 32)
	if err != nil {
		return ConfigurationData{}, err
	}
	a = int32(n)

	var parts []PartData
	var partIds = values["part"]
	for partId := range partIds {
		part, err := getPart(db, partIds[partId])
		if err != nil {
			return ConfigurationData{}, err
		}
		parts = append(parts, part)
	}

	return ConfigurationData{name, description, p, a, additionalPriceDescription, parts, id}, nil
}

// Main function for displaying configurations in configurations.html.
func configurationHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var action string = request.URL.Query().Get("action")
	switch action {
	case "show":
		showConfigurationsAction( writer, db )
	case "save":
		saveConfigurationAction( writer, request, db )
	case "delete":
		deleteConfigurationAction( request, db )
	default:
		fmt.Println("Unknown action. Exiting.")
		os.Exit(1)
	}
}

func showConfigurationsAction( writer http.ResponseWriter, db *sql.DB ) {
	configurations := getConfigurations(db)
	configurationTemplate.Execute(writer, configurations)
}

func saveConfigurationAction( writer http.ResponseWriter, request *http.Request, db *sql.DB ) {
	var id = request.URL.Query().Get("id")

	configurationDataFromRequest, err := getConfigurationDataFromRequest(request, db)
	if err != nil {
		failureMessage := "Failed to get data from form: " + err.Error()
		saveResultTemplate.Execute(writer, failureMessage)
		return
	}
	// reference data, name should be essential
	if configurationDataFromRequest.Name != "" {
		// id = "" -> unsaved item: blank page (just opened) OR new item that has no (database Id) yet
		if id == "" {
			err := insertConfiguration(db, configurationDataFromRequest)
			if err != nil {
				failureMessage := "Failed to save new item: " + err.Error()
				saveResultTemplate.Execute(writer, failureMessage)
				return
			}
		} else { // saved = existing item
			// if the request contains data (other then the id) then the edited data is sent for storing
			err := updateConfiguration(db, configurationDataFromRequest)
			if err != nil {
				failureMessage := "Failed to save item: " + err.Error()
				saveResultTemplate.Execute(writer, failureMessage)
				return
			}
		}
	}
	// if everything is OK then display the Parts list
	// TODO: here the URL need to be cleared, otherwise refresing may save the new item multiple times
	showConfigurationsAction( writer, db )
}

func deleteConfigurationAction( request *http.Request, db *sql.DB ) {
	var id string = request.URL.Query().Get("id")
	if id != "" {
		err := deleteConfiguration(db, id)
		if err != nil {
			fmt.Println(err)
		}
	}
}
