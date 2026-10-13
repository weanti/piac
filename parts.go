package main

import (
	"database/sql"
	"errors"
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

func createPartsTable(db *sql.DB) error {
	rows, err := db.Query(`CREATE TABLE Parts(
						Brand TEXT NOT NULL,
						Model TEXT NOT NULL, 
						Description TEXT DEFAULT(''),
						Count INTEGER DEFAULT(1),
						Price INTEGER DEFAULT(0),
						Image TEXT);`)
	defer rows.Close()
	return err
}

func sqlRowToPartData(rows *sql.Rows) (PartData, error) {
	var partData PartData
	var brand, model, description, image, id *string
	var count *int16
	var price *int32
	if rows.Next() {
		err := rows.Scan(&id,
			&brand,
			&model,
			&description,
			&count,
			&price,
			&image)
		if err != nil {
			fmt.Println(err)
			return PartData{}, err
		}
		if brand != nil {
			partData.Brand = *brand
		}
		if model != nil {
			partData.Model = *model
		}
		if description != nil {
			partData.Description = *description
		}
		if count != nil {
			partData.Count = *count
		}
		if price != nil {
			partData.Price = *price
		}
		if image != nil {
			partData.Image = *image
		}
		if id != nil {
			partData.Id = *id
		}
		return partData, nil
	}
	return partData, errors.New("No data available in query.")
}

func partsHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	parts := getParts(db)
	partsTemplate.Execute(writer, parts)
}

func getParts(db *sql.DB) []PartData {
	var parts []PartData
	rows, err := db.Query(`SELECT rowid,* from Parts;`)
	defer rows.Close()
	if err == nil {
		for err == nil {
			var partData PartData
			partData, err = sqlRowToPartData(rows)
			if err == nil {
				parts = append(parts, partData)
			}
		}
	} else {
		fmt.Println(err)
	}
	return parts
}

func getPart(db *sql.DB, id string) PartData {
	var partData PartData
	rows, err := db.Query("SELECT rowid,* FROM Parts WHERE rowid=?", id)
	defer rows.Close()
	if err == nil {
		partData, err = sqlRowToPartData(rows)
	}
	if err != nil {
		fmt.Println(err)
	}
	return partData
}

func updatePart(db *sql.DB, part PartData) error {
	rows, err := db.Query(	`UPDATE Parts SET (Brand, Model, Description, Count, Price, Image ) = ('` +
											part.Brand + `','` +
											part.Model + `','` +
											part.Description + `','` +
											strconv.Itoa( int(part.Count) ) + `','` +
											strconv.Itoa( int(part.Price) ) + `','` +
											part.Image + `' )` +
										` WHERE rowid = ` + part.Id + `;` )
	defer rows.Close()
	return err
}

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

func insertNewPart( request *http.Request, db *sql.DB ) PartData {
	partData := getPartDataFromRequest( request )
	rows, err := db.Query(`INSERT INTO Parts (
									Brand,
									Model,
									Description,
									Count,
									Price,
									Image )
							VALUES ('` + partData.Brand + `','` +
										partData.Model + `','` +
										partData.Description + `',` +
										strconv.Itoa( int(partData.Count) ) + `,` +
										strconv.Itoa( int(partData.Price) ) + `,'` +
										partData.Image + `');`)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	return partData
}

func editPartHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var id = request.URL.Query().Get("id")
	var partData PartData

	if id == "" {
		var brand = request.URL.Query().Get("brand")
		// let brand be the reference, if it's filled then the item shall be saved
		if brand != "" {
			partData = insertNewPart( request, db )
		} // else it's a completely empty entry, nothing to do
	} else {
		// load existing data
		partData = getPart( db, id )
		// if request contains data then it is modified
		partDataFromRequest := getPartDataFromRequest( request )
		if partDataFromRequest.Model != "" {
			err := updatePart( db, partDataFromRequest )
			if err == nil {
				partData = partDataFromRequest
			} else {
				fmt.Println(err)
			}
		}
	}
	edit_partTemplate.Execute(writer, partData)
}
