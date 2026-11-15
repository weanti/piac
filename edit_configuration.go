package main

import (
	"net/http"
	"database/sql"
)

// Main function for displaying edited configuration in edit_configuration.html.
func editConfigurationHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var id = request.URL.Query().Get("id")
	var configurationData ConfigurationData

	// id = "" -> unsaved item: blank page (just opened) OR new item that has no (database Id) yet
	if id == "" {
		// it's a completely new empty entry, nothing to do
	} else { // saved = existing item
		// by default the data from database will be displayed
		configurationData = getConfiguration(db, id)
	}
	var freeParts []PartData = getParts(db)
	for _, usedpart := range configurationData.Parts {
		for i, freepart := range freeParts {
			if usedpart.Id == freepart.Id {
				freeParts = remove(freeParts, i)
			}
		}
	}
	var data = map[string]any{
		"Configuration": configurationData,
		"FreeParts":     freeParts,
	}
	edit_configurationTemplate.Execute(writer, data)
}
