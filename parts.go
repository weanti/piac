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
			fmt.Println(err)
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
	}
	return tempPartData, errors.New("No data available in query.")
}

func partsHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var parts []PartData
	rows, err := db.Query(`SELECT rowid,* from Parts;`)
	if err != nil {
		fmt.Println(err)
	} else {
		err = nil
		for err == nil {
			var tempPartData PartData
			tempPartData, err = sqlRowToPartData(rows)
			if err == nil {
				parts = append(parts, tempPartData)
			}
		}
	}
	partsTemplate.Execute(writer, parts)
}

func editPartHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var id = request.URL.Query().Get("id")
	tempPartData := PartData{}
	if id == "" {
		// can be a completely new item or an unsaved item
		// let's examine other params
		var brand = request.URL.Query().Get("brand")
		var model = request.URL.Query().Get("model")
		var description = request.URL.Query().Get("description")
		var count = request.URL.Query().Get("count")
		var price = request.URL.Query().Get("price")
		var image = request.URL.Query().Get("image")
		// let brand be the reference, if it's filled then the item shall be saved
		if brand != "" {
			// save mode
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
			tempPartData = PartData{brand, model, description, c, p, image, 0}
		}
	} else {
		rows, err := db.Query("SELECT rowid,* FROM Parts WHERE rowid=?", id)
		if err == nil {
			tempPartData, err = sqlRowToPartData(rows)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}
	edit_partTemplate.Execute(writer, tempPartData)
}
