package main

import ( "net/http"
	"database/sql"
		"fmt")

type PartData struct {
    Brand string
    Model string
    Description string
    Count int16
    Price int32
    Image string
    InternalId int32
}

func partsHandler( writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var parts []PartData
	rows, err := db.Query( `SELECT * from Parts` )
	for rows.Next() {
		var tempPartData PartData 
		var brand, model, description, image *string
		var count *int16
		var price, internalId *int32

		err = rows.Scan( 	&brand,
							&model,
							&description,
							&count,
							&price,
							&image,
							&internalId )
		if err != nil {
			fmt.Println( err )
		} else {
			if brand != nil {
				tempPartData.Brand = *brand
			}
			if model != nil {
				tempPartData.Model = *model
			}
			if description != nil {
				tempPartData.Description = *description
			}
			if count != nil {
				tempPartData.Count = *count
			}
			if price != nil {
				tempPartData.Price = *price
			}
			if image != nil {
				tempPartData.Image = *image
			}
			if internalId != nil {
				tempPartData.InternalId = *internalId
			}
			parts = append( parts, tempPartData )
		}
	}
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
