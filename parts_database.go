package main

import (
    "database/sql"
	"fmt"
	"strconv"
)

// Create Parts database table
func createPartsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE Parts(	Id INTEGER PRIMARY KEY,
						Brand TEXT NOT NULL,
						Model TEXT NOT NULL, 
						Description TEXT DEFAULT(''),
						Count INTEGER DEFAULT(1),
						Price INTEGER DEFAULT(0),
						Image TEXT);`)
	return err
}

// Get all parts from database
func getParts(db *sql.DB) []PartData {
	var parts []PartData
	rows, err := db.Query(`SELECT * from Parts;`)
	defer rows.Close()
	if err == nil {
		for rows.Next() {
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

// Get specific part from database
func getPart(db *sql.DB, id string) PartData {
	var partData PartData
	rows, err := db.Query("SELECT * FROM Parts WHERE Id=?", id)
	defer rows.Close()
	if err == nil {
		if rows.Next() {
			partData, err = sqlRowToPartData(rows)
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return partData
}

// Update the part in database. Part is identified by PartData.Id, the rest is the new data.
func updatePart(db *sql.DB, part PartData) error {
	_, err := db.Exec(	`UPDATE Parts SET (Brand, Model, Description, Count, Price, Image ) = ('` +
											part.Brand + `','` +
											part.Model + `','` +
											part.Description + `','` +
											strconv.Itoa( int(part.Count) ) + `','` +
											strconv.Itoa( int(part.Price) ) + `','` +
											part.Image + `' )` +
										` WHERE Id = ` + part.Id + `;` )
	return err
}
// Create new part in database.
func insertPart(db *sql.DB, partData PartData ) error {
	_, err := db.Exec(`INSERT INTO Parts (
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
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func deletePart(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE from Parts WHERE Id = ?", id )
	return err
 }

// Helper function to convert database query result to PartData
func sqlRowToPartData(rows *sql.Rows) (PartData, error) {
	var partData PartData
	var brand, model, description, image, id *string
	var count *int16
	var price *int32

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
