package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

// Main function for displaying edited part in edit_part.html
func editPartHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var id = request.URL.Query().Get("id")
	var partData PartData

	if id == "" {
		// a completely new empty entry, nothing to do
	} else {
		// existing entry, load data for editing
		var err error
		partData, err = getPart(db, id)
		if err != nil {
			fmt.Println(err)
		}
	}
	edit_partTemplate.Execute(writer, partData)
}
