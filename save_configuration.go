package main 

import (
	"database/sql"
	"fmt"
	"net/http"
)

func saveConfigurationHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var id = request.URL.Query().Get("id")

	// id = "" -> unsaved item: blank page (just opened) OR new item that has no (database Id) yet
	if id == "" {
		configurationDataFromRequest := getConfigurationDataFromRequest(request, db)
		// reference data, name should be essential
		if configurationDataFromRequest.Name != "" {
			err := insertNewConfiguration(db, configurationDataFromRequest)
			if err != nil {
				fmt.Println( err )
			}
		}
	} else { // saved = existing item
		configurationDataFromRequest := getConfigurationDataFromRequest(request, db)
		// if the request contains data (other then the id) then the edited data is sent for storing
		// reference data, name should be essential
		if configurationDataFromRequest.Name != "" {
			err := updateConfiguration(db, configurationDataFromRequest)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	// TODO: here a Saved or Failed message can be displayed
	// now redirect to Configurations
}
