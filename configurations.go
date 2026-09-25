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
						AdditionalPriceDesc TEXT DEFAULT(''));` )
						// TODO: somehow store the references here
	return err
}

func fromQueryResult( rows *sql.Rows ) (ConfigurationData, error) {
	var tempConfigurationData ConfigurationData
	var name, description, additionalPriceDesc *string
	var price, additionalPrice, internalId *int32
	if rows.Next() {
	// TODO: array of Parts internal IDs
		err = rows.Scan(&internalId,
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

func configurationHandler( writer http.ResponseWriter, request *http.Request, db *sql.DB ) {
	var configurations []ConfigurationData
	rows, err := db.Query(`SELECT rowid,* from Configurations;`)
	for rows.Next() {
		var tempConfigurationData ConfigurationData
		tempConfigurationData, err = fromQueryResult( rows )
		if err == nil {
			configurations = append( configurations, tempConfigurationData )
		}
	}
	configurationTemplate.Execute( writer, configurations )
	// TODO: get ID from URL
	// if empty then all fields shall be empty and a STORE operation is to be performed
	// if not empty then query all fields (is this redundant? maybe recevie everything from query?)
	//	in this case an update need to be performed
}
