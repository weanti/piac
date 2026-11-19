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
	var action string = request.URL.Query().Get("action")
	switch action {
	case "show":
		showPartsAction( writer, db )
	case "save":
		savePartAction( request, db )
	case "delete":
		deletePartAction( request, db )
	default:
		// TODO: exit here!
		fmt.Println("Unknown action")
	}
}

func showPartsAction( writer http.ResponseWriter, db *sql.DB ) {
	parts := getParts(db)
	partsTemplate.Execute(writer, parts)
}

func savePartAction(request *http.Request, db *sql.DB ) {
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

func deletePartAction( request *http.Request, db *sql.DB ) {
	// TODO: check if id is not empty
	var id string = request.URL.Query().Get("id")
	err := deletePart( db, id )
	if err != nil {
		fmt.Println(err)
	}
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

