package main

import (
    "database/sql"
	"errors"
	"fmt"
	"strconv"
)
// Create Configurations database table
func createConfigurationsTable(db *sql.DB) error {
	rows, err := db.Query(`CREATE TABLE Configurations(
						Name TEXT NOT NULL,
						Description TEXT DEFAULT(''),
						Price INTEGER DEFAULT(0),
						AdditionalPrice INTEGER DEFAULT(0),
						AdditionalPriceDesc TEXT DEFAULT(''),
						Parts BLOB);`)
	defer rows.Close()
	return err
}

// Get all configurations from database
func getConfigurations(db *sql.DB) []ConfigurationData {
	var configurations []ConfigurationData
	rows, err := db.Query(`SELECT rowid,* from Configurations;`)
	defer rows.Close()
	for rows.Next() {
		var configurationData ConfigurationData
		configurationData, err = sqlRowToConfigurationData(db, rows)
		if err == nil {
			configurations = append(configurations, configurationData)
		}
	}
	return configurations
}

// Get specific configuration from database
func getConfiguration(db *sql.DB, id string) ConfigurationData {
	rows, err := db.Query("SELECT rowid,* FROM Configurations WHERE rowid=?", id)
	defer rows.Close()
	var configurationData ConfigurationData
	if err == nil {
		if rows.Next() {
			configurationData, err = sqlRowToConfigurationData(db, rows)
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return configurationData
}

// Create new configuration in database
func insertNewConfiguration(db *sql.DB, configurationData ConfigurationData) error {
	var ids []uint16
	for i := range configurationData.Parts {
		n, err := strconv.ParseInt(configurationData.Parts[i].Id, 10, 16)
		if err == nil {
			ids = append(ids, uint16(n))
		} else {
			fmt.Println(err)
		}
	}
	var bytes []byte = uint16ArrayToByteArray(ids)
	_, err := db.Exec(`INSERT INTO Configurations(
										Name,
										Description,
										Price,
										AdditionalPrice,
										AdditionalPriceDesc,
										Parts	
										)
								VALUES(?, ?, ?, ?, ?, ?)`, configurationData.Name,
		configurationData.Description,
		strconv.Itoa(int(configurationData.Price)),
		strconv.Itoa(int(configurationData.AdditionalPrice)),
		configurationData.AdditionalPriceDesc,
		bytes)
	fmt.Println("name ", configurationData.Name )
	return err
}

// Update the configuration in database. Configuration is identified by ConfigurationData.Id, the rest is the new data.
func updateConfiguration(db *sql.DB, configuration ConfigurationData) error {
	var ids []uint16
	for i := range configuration.Parts {
		n, err := strconv.ParseInt(configuration.Parts[i].Id, 10, 16)
		if err == nil {
			ids = append(ids, uint16(n))
		} else {
			fmt.Println(err)
		}
	}
	var bytes []byte = uint16ArrayToByteArray(ids)
	_, err := db.Exec(`UPDATE Configurations 
								SET Name = ?,
									Description = ?,
									Price = ?,
									AdditionalPrice = ?,
									AdditionalPriceDesc = ?,
									Parts = ?
								WHERE rowid = ?`,
		configuration.Name,
		configuration.Description,
		strconv.Itoa(int(configuration.Price)),
		strconv.Itoa(int(configuration.AdditionalPrice)),
		configuration.AdditionalPriceDesc,
		bytes,
		configuration.Id)
	return err
}

// Helper function to convert database query result to ConfigurationData
func sqlRowToConfigurationData(db *sql.DB, rows *sql.Rows) (ConfigurationData, error) {
	var configurationData ConfigurationData
	var name, description, additionalPriceDesc, id *string
	var price, additionalPrice *int32
	var bytes []byte
	err := rows.Scan(&id,
		&name,
		&description,
		&price,
		&additionalPrice,
		&additionalPriceDesc,
		&bytes)
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
		if bytes != nil {
			// TODO: maybe do this separately: store only part ids here and fill parts in another function
			// db shoult not be part of this function
			var ids []uint16 = byteArrayToUint16Array(bytes)
			for i := range ids {
				var s string = strconv.Itoa(int(ids[i]))
				var part PartData = getPart(db, s)
				// TODO: how to detect error
				configurationData.Parts = append(configurationData.Parts, part)
			}
		}
		return configurationData, nil
	}

	return configurationData, errors.New("No data available in query.")
}
