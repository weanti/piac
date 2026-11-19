package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func savePartHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB ) {
	var id = request.URL.Query().Get("id")
	partDataFromRequest := getPartDataFromRequest(request)
	if partDataFromRequest.Brand != "" {
		if id == "" {
			// saving new item
			err := insertPart(db, partDataFromRequest)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// update existing item
			err := updatePart(db, partDataFromRequest)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
