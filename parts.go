package main

import ( "net/http"
	"database/sql"
	"errors"
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

func createPartsTable(db *sql.DB) error {
	_, err := db.Query(`CREATE TABLE Parts(
						Brand TEXT NOT NULL,
						Model TEXT NOT NULL, 
						Description TEXT DEFAULT(''),
						Count INTEGER DEFAULT(1),
						Price INTEGER DEFAULT(0),
						Image BLOB);` )
	return err
}

func fromQueryResult( rows *sql.Rows ) (PartData, error) {
	var tempPartData PartData 
	var brand, model, description, image *string
	var count *int16
	var price, internalId *int32
	if rows.Next() {
		err := rows.Scan(&internalId,
						&brand,
						&model,
						&description,
						&count,
						&price,
						&image)
		if err != nil {
			fmt.Println( err )
			return PartData{}, err 
		}
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
		return tempPartData, nil
	} else {
		return PartData{}, errors.New( "no data available" )
	}
}

func partsHandler( writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var parts []PartData
	rows, err := db.Query( `SELECT rowid,* from Parts;` )
	if err != nil {
		fmt.Println( err )
	} else {
		for {
			tempPartData, err := fromQueryResult( rows )
			if err == nil {
				parts = append( parts, tempPartData )
			} else {
				break
			}
		}
	}
	partsTemplate.Execute(writer, parts)
}

func editPartHandler( writer http.ResponseWriter, request *http.Request, db *sql.DB ) {
	var id = request.URL.Query().Get("id")
	var tempPartData PartData
	if id == "" {
		fmt.Println("new item")
		tempPartData = PartData{}
	} else {
		rows, err := db.Query("SELECT rowid,* FROM Parts WHERE rowid=?", id )
		if err != nil {
			fmt.Println(err)
		} else {
			tempPartData, err = fromQueryResult( rows )
			if err != nil {
			} else {
				// TODO: get ID from URL
				// if empty then all fields shall be empty and a STORE operation is to be performed
				// if not empty then query all fields (is this redundant? maybe recevie everything from query?)
				//	in this case an update need to be performed
			}
		}
	}
	edit_partTemplate.Execute(writer, tempPartData)
}
