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
	InternalId  int32
}

func createPartsTable(db *sql.DB) error {
	_, err := db.Query(`CREATE TABLE Parts(
						Brand TEXT NOT NULL,
						Model TEXT NOT NULL, 
						Description TEXT DEFAULT(''),
						Count INTEGER DEFAULT(1),
						Price INTEGER DEFAULT(0),
						Image TEXT);`)
	return err
}

func sqlRowToPartData(rows *sql.Rows) (PartData, error) {
	var partData PartData
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
		if internalId != nil {
			partData.InternalId = *internalId
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
	if err == nil {
		partData, err = sqlRowToPartData(rows)
	}
	if err != nil {
		fmt.Println(err)
	}
	return partData
}

func insertNewPart( request *http.Request, db *sql.DB ) PartData {
	var brand = request.URL.Query().Get("brand")
	var model = request.URL.Query().Get("model")
	var description = request.URL.Query().Get("description")
	var count = request.URL.Query().Get("count")
	var price = request.URL.Query().Get("price")
	var image = request.URL.Query().Get("image")
	_, err := db.Query(`INSERT INTO Parts(
									Brand,
									Model,
									Description,
									Count,
									Price,
									Image )
							VALUES ('` + brand + `','` +
										model + `','` +
										description + `',` +
										count + `,` +
										price + `,'` +
										image + `');`)
	if err != nil {
		fmt.Println(err)
	}
	var n int64
	var c int16
	n, err = strconv.ParseInt(count, 10, 16)
	if err != nil {
		c = 0
	} else {
		c = int16(n)
	}
	var p int32
	n, err = strconv.ParseInt(price, 10, 32)
	if err != nil {
		p = 0
	} else {
		p = int32(n)
	}
	return  PartData{brand, model, description, c, p, image, 0}
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
		partData = getPart( db, id )
		// TODO: update data
	}
	edit_partTemplate.Execute(writer, partData)
}
