package main

import( "database/sql"
		"fmt"
		"net/http" )

type ConfigurationData struct {
    Name string
    Description string
    Price int32
    AdditionalPrice int32
    AdditionalPriceDesc string
    Parts []PartData
    InternalId int32
}

func createConfigurationsTable( db *sql.DB ) error {
	_, err := db.Query(`CREATE TABLE Configurations(
						Name TEXT NOT NULL,
						Description TEXT DEFAULT(''),
						Price INTEGER DEFAULT(0),
						AdditionalPrice INTEGER DEFAULT(0),
						AdditionalPriceDesc TEXT DEFAULT(''),
						InternalId INTEGER);` )
						// TODO: somehow store the references here
	return err
}

func configurationHandler( writer http.ResponseWriter, request *http.Request, db *sql.DB ) {
	var configurations []ConfigurationData
	rows, err := db.Query(`SELECT * from Configurations;`)
	for rows.Next() {
		var tempConfigurationData ConfigurationData
		var name, description, additionalPriceDesc *string
		var price, additionalPrice, internalId *int32
		// TODO: array of Parts internal IDs
		err = rows.Scan(&name,
						&description,
						&price,
						&additionalPrice,
						&additionalPriceDesc,
						&internalId )
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
			configurations = append( configurations, tempConfigurationData )
		}
	}
	configurationTemplate.Execute( writer, configurations )
	// TODO: get ID from URL
	// if empty then all fields shall be empty and a STORE operation is to be performed
	// if not empty then query all fields (is this redundant? maybe recevie everything from query?)
	//	in this case an update need to be performed
//	configurationTemplate.Execute(writer, examplePartData)
}
