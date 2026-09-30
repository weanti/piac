package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type ConfigurationData struct {
	Name                string
	Description         string
	Price               int32
	AdditionalPrice     int32
	AdditionalPriceDesc string
	Parts               []PartData
	InternalId          int32
}

func createConfigurationsTable(db *sql.DB) error {
	_, err := db.Query(`CREATE TABLE Configurations(
						Name TEXT NOT NULL,
						Description TEXT DEFAULT(''),
						Price INTEGER DEFAULT(0),
						AdditionalPrice INTEGER DEFAULT(0),
						AdditionalPriceDesc TEXT DEFAULT(''));`)
	// TODO: somehow store the references here
	return err
}

func sqlRowToConfigurationData(rows *sql.Rows) (ConfigurationData, error) {
	var tempConfigurationData ConfigurationData
	var name, description, additionalPriceDesc *string
	var price, additionalPrice, internalId *int32
	if rows.Next() {
		// TODO: array of Parts internal IDs
		err := rows.Scan(&internalId,
			&name,
			&description,
			&price,
			&additionalPrice,
			&additionalPriceDesc)
		if err != nil {
			fmt.Println(err)
		} else {
			if name != nil {
				tempConfigurationData.Name = *name
			}
			if description != nil {
				tempConfigurationData.Description = *description
			}
			if additionalPriceDesc != nil {
				tempConfigurationData.AdditionalPriceDesc = *additionalPriceDesc
			}
			if price != nil {
				tempConfigurationData.Price = *price
			}
			if additionalPrice != nil {
				tempConfigurationData.AdditionalPrice = *additionalPrice
			}
			if internalId != nil {
				tempConfigurationData.InternalId = *internalId
			}
			return tempConfigurationData, nil
		}
	}
	return tempConfigurationData, errors.New("No data available in query.")
}

func configurationHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var configurations []ConfigurationData
	rows, err := db.Query(`SELECT rowid,* from Configurations;`)
	for rows.Next() {
		var tempConfigurationData ConfigurationData
		tempConfigurationData, err = sqlRowToConfigurationData(rows)
		if err == nil {
			configurations = append(configurations, tempConfigurationData)
		}
	}
	configurationTemplate.Execute(writer, configurations)
	// TODO: get ID from URL
	// if empty then all fields shall be empty and a STORE operation is to be performed
	// if not empty then query all fields (is this redundant? maybe recevie everything from query?)
	//	in this case an update need to be performed
}

func editConfigurationHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	
	fmt.Println("edit config")
	var id = request.URL.Query().Get("id")
	tempConfigurationData := ConfigurationData{}
	if id == "" {
		var name = request.URL.Query().Get("name")
		var description = request.URL.Query().Get("description")
		var price = request.URL.Query().Get("price")
		var additionalPrice = request.URL.Query().Get("additionalPrice")
		var additionalPriceDescription = request.URL.Query().Get("additionalPriceDescription")
		if name != "" {
			_, err := db.Query( `INSERT INTO Configurations(
												Name,
												Description,
												Price,
												AdditionalPrice,
												AdditionalPriceDesc )
											VALUES('` + name + `','` +
														description + `','` +	
														price + `','` +
														additionalPrice + `','` +
														additionalPriceDescription + `');'`)
			if err != nil {
				fmt.Println(err)
			}
			var n int64
			var p int32
			n, err = strconv.ParseInt( price, 10, 32 )
			if err != nil {
				p = int32(n)
			}
			var a int32
			n, err = strconv.ParseInt( additionalPrice, 10, 32 )
			if err != nil {
				a = int32(n)
			}
			tempConfigurationData = ConfigurationData{name, description, p, a, additionalPriceDescription, nil, 0 }
		}
	} else {
		rows, err := db.Query("SELECT rowid,* FROM Configurations WHERE rowid=?", id)
		if err == nil {
			tempConfigurationData, err = sqlRowToConfigurationData(rows)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}
	edit_configurationTemplate.Execute( writer, tempConfigurationData )
}
