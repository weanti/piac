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
	Id          		string
}
/*****************
*    Database    *
******************/
// Create Configurations database table
func createConfigurationsTable(db *sql.DB) error {
	rows, err := db.Query(`CREATE TABLE Configurations(
						Name TEXT NOT NULL,
						Description TEXT DEFAULT(''),
						Price INTEGER DEFAULT(0),
						AdditionalPrice INTEGER DEFAULT(0),
						AdditionalPriceDesc TEXT DEFAULT(''));`)
	defer rows.Close()
	// TODO: somehow store the references here
	return err
}

// Helper function to convert database query result to ConfigurationData
func sqlRowToConfigurationData(rows *sql.Rows) (ConfigurationData, error) {
	var configurationData ConfigurationData
	var name, description, additionalPriceDesc, id *string
	var price, additionalPrice *int32
	// TODO: array of Parts internal IDs
	err := rows.Scan(&id,
		&name,
		&description,
		&price,
		&additionalPrice,
		&additionalPriceDesc)
	if err != nil {
		fmt.Println(err)
	} else {
		if name != nil {
			configurationData.Name = *name
		}
		if description != nil {
			configurationData.Description = *description
		}
		if additionalPriceDesc != nil {
			configurationData.AdditionalPriceDesc = *additionalPriceDesc
		}
		if price != nil {
			configurationData.Price = *price
		}
		if additionalPrice != nil {
			configurationData.AdditionalPrice = *additionalPrice
		}
		if id != nil {
			configurationData.Id = *id
		}
		return configurationData, nil
	}

	return configurationData, errors.New("No data available in query.")
}

// Get all configurations from database
func getConfigurations( db *sql.DB) []ConfigurationData {
	var configurations []ConfigurationData
	rows, err := db.Query(`SELECT rowid,* from Configurations;`)
	defer rows.Close()
	for rows.Next() {
		var configurationData ConfigurationData
		configurationData, err = sqlRowToConfigurationData(rows)
		if err == nil {
			configurations = append(configurations, configurationData)
		}
	}
	return configurations
}

// Get specific configuration from database
func getConfiguration( db *sql.DB, id string ) ConfigurationData {
	rows, err := db.Query("SELECT rowid,* FROM Configurations WHERE rowid=?", id)
	defer rows.Close()
	var configurationData ConfigurationData
	if err == nil {
		if rows.Next() {
			configurationData, err = sqlRowToConfigurationData(rows)
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return configurationData
}

// Create new configuration in database
func insertNewConfiguration( configurationData ConfigurationData, db *sql.DB ) error {
	rows, err := db.Query( `INSERT INTO Configurations(
										Name,
										Description,
										Price,
										AdditionalPrice,
										AdditionalPriceDesc )
								VALUES('` + configurationData.Name + `','` +
											configurationData.Description + `','` +	
											strconv.Itoa( int(configurationData.Price) ) + `','` +
											strconv.Itoa( int(configurationData.AdditionalPrice) ) + `','` +
											configurationData.AdditionalPriceDesc + `');`)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}

	return err
}

// Update the configuration in database. Configuration is identified by ConfigurationData.Id, the rest is the new data.
func updateConfiguration( db *sql.DB, configuration ConfigurationData ) error {
	rows, err := db.Query(`UPDATE Configurations SET ( Name, Description, Price, AdditionalPrice, AdditionalPriceDesc ) = ('` +
								configuration.Name + `','` +
								configuration.Description + `','` +
								strconv.Itoa( int(configuration.Price) ) + `','` +
								strconv.Itoa( int(configuration.AdditionalPrice) ) + `','` +
								configuration.AdditionalPriceDesc + `' )` +
							` WHERE rowid = ` + configuration.Id + `;` )
	defer rows.Close()
	return err
}


/*****************
*      HTML      *
******************/

// Helper function that convert GET method params to ConfigurationData
func getConfigurationDataFromRequest( request *http.Request, db *sql.DB ) ConfigurationData {
	var values = request.URL.Query()
	var name = values.Get("name")
	var description = values.Get("description")
	var price = values.Get("price")
	var additionalPrice = values.Get("additionalPrice")
	var additionalPriceDescription = values.Get("additionalPriceDescription")
	var id = values.Get("id")
	var n int64
	var p int32 = 0
	var err error
	n, err = strconv.ParseInt( price, 10, 32 )
	if err == nil {
		p = int32(n)
	}
	var a int32 = 0
	n, err = strconv.ParseInt( additionalPrice, 10, 32 )
	if err == nil {
		a = int32(n)
	}
	
	var parts []PartData
	var partIds = values["part"]
	for partId := range partIds {
		parts = append( parts, getPart( db, partIds[partId] ) )
	}

	return ConfigurationData{ name, description, p, a, additionalPriceDescription, parts, id }
}

// Main function for displaying configurations in configurations.html.
func configurationHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	configurations := getConfigurations( db )
	configurationTemplate.Execute(writer, configurations)
}

// Main function for displaying edited configuration in edit_configuration.html.
func editConfigurationHandler(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	var id = request.URL.Query().Get("id")
	var configurationData ConfigurationData

	if id == "" {
		var name = request.URL.Query().Get("name")
		if name != "" {
			configurationData := getConfigurationDataFromRequest( request, db )
			insertNewConfiguration( configurationData, db )
		} // else it's a completely new empty entry, nothing to do
	} else {
		configurationData = getConfiguration( db, id )
		configurationDataFromRequest := getConfigurationDataFromRequest( request, db )
		if configurationDataFromRequest.Name != "" {
			err := updateConfiguration( db, configurationDataFromRequest )
			if err == nil {
				configurationData = configurationDataFromRequest
			} else {
				fmt.Println(err)
			}
		}
	}

	configurationData.Parts = getParts( db )
	edit_configurationTemplate.Execute( writer, configurationData )
}
