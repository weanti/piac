package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type PartData struct {
	Brand       string
	Model       string
	Description string
	Count       int16
	Price       int32
	Image       string
	Id		string	
}

// Main function for displaying parts list in parts.html
func partsHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var action string = request.URL.Query().Get("action")
	switch action {
	case "show":
		showPartsAction( writer, db )
	case "save":
		savePartAction( writer, request, db )
	case "delete":
		deletePartAction( request, db )
	default:
		fmt.Println("Unknown action")
		os.Exit(1)
	}
}

func showPartsAction( writer http.ResponseWriter, db *sql.DB ) {
	parts := getParts(db)
	partsTemplate.Execute(writer, parts)
}

func savePartAction( writer http.ResponseWriter, request *http.Request, db *sql.DB ) {
	var id = request.URL.Query().Get("id")
	partDataFromRequest, err := getPartDataFromRequest(request)
	if err != nil {
		failureMessage := "Failed to get data from form: " + err.Error()
		saveResultTemplate.Execute(writer, failureMessage)
		return
	} else {
		if partDataFromRequest.Brand != "" {
			if id == "" {
				// saving new item
				err := insertPart(db, partDataFromRequest)
				if err != nil {
					failureMessage := "Failed to save new item: " + err.Error()
					saveResultTemplate.Execute(writer, failureMessage)
					return
				}
			} else {
				// update existing item
				err := updatePart(db, partDataFromRequest)
				if err != nil {
					failureMessage := "Failed to save item: " + err.Error()
					saveResultTemplate.Execute(writer, failureMessage)
					return
				}
			}
		}
	}
	// if everything is OK then display the Parts list
	showPartsAction( writer, db )
}

func deletePartAction( request *http.Request, db *sql.DB ) {
	var id string = request.URL.Query().Get("id")
	if id == "" {
		fmt.Println("Deleting Part with empty id is forbidden")
	} else {
		err := deletePart( db, id )
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Helper function that convert GET method params to PartData
func getPartDataFromRequest( request *http.Request ) (PartData, error) {
	query := request.URL.Query()
	var brand = query.Get("brand")
	var model = query.Get("model")
	var description = query.Get("description")
	var count = query.Get("count")
	var price = query.Get("price")
	var image = query.Get("image")
	var id = query.Get("id")
	var n int64
	var c int16 = 0
	var err error
	n, err = strconv.ParseInt(count, 10, 16)
	if err != nil {
		return PartData{}, err
	}
	c = int16(n)
	var p int32 = 0
	n, err = strconv.ParseInt(price, 10, 32)
	if err != nil {
		return PartData{}, err
	}
	p = int32(n)
	return  PartData{brand, model, description, c, p, image, id}, nil
}

