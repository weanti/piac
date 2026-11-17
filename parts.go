package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type PartData struct {
	Brand       string
	Model       string
	Description string
	Count       int16
	Price       int32
	Image       string
	Id  		string
}

// Main function for displaying parts list in parts.html
func partsHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	parts := getParts(db)
	partsTemplate.Execute(writer, parts)
}

// Main function for displaying edited part in edit_part.html
func editPartHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var id = request.URL.Query().Get("id")
	var partData PartData

	if id == "" {
		var brand = request.URL.Query().Get("brand")
		// let brand be the reference, if it's filled then the item shall be saved
		if brand != "" {
			partData := getPartDataFromRequest( request )
			insertNewPart( partData, db )
		} // else it's a completely empty entry, nothing to do
	} else {
		// load existing data
		partData = getPart( db, id )
		// if request contains data then it is modified
		partDataFromRequest := getPartDataFromRequest( request )
		if partDataFromRequest.Model != "" {
			err := updatePart( partDataFromRequest, db )
			if err == nil {
				partData = partDataFromRequest
			} else {
				fmt.Println(err)
			}
		}
	}
	edit_partTemplate.Execute(writer, partData)
}

// Helper function that convert GET method params to PartData
func getPartDataFromRequest( request *http.Request ) PartData {
	var brand = request.URL.Query().Get("brand")
	var model = request.URL.Query().Get("model")
	var description = request.URL.Query().Get("description")
	var count = request.URL.Query().Get("count")
	var price = request.URL.Query().Get("price")
	var image = request.URL.Query().Get("image")
	var id = request.URL.Query().Get("id")
	var n int64
	var c int16 = 0
	var err error
	n, err = strconv.ParseInt(count, 10, 16)
	if err == nil {
		c = int16(n)
	}
	var p int32 = 0
	n, err = strconv.ParseInt(price, 10, 32)
	if err == nil {
		p = int32(n)
	}
	return  PartData{brand, model, description, c, p, image, id}
}

