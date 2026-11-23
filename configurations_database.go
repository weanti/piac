package main

import (
    "database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
)
// Create Configurations database table
func createConfigurationsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE Configurations(
						Id Integer PRIMARY KEY,
						Name TEXT NOT NULL,
						Description TEXT DEFAULT(''),
						Price INTEGER DEFAULT(0),
						AdditionalPrice INTEGER DEFAULT(0),
						AdditionalPriceDesc TEXT DEFAULT(''));`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`PRAGMA foreign_keys = ON;
				CREATE TABLE ConfigurationPart (
					Configuration_Id INTEGER NOT NULL,
					Part_Id INTEGER NOT NULL,
					PRIMARY KEY (Configuration_Id, Part_Id),
					FOREIGN KEY (Configuration_Id)
						REFERENCES Configuration(Id)
						ON DELETE CASCADE,
					FOREIGN KEY (Part_Id)
						REFERENCES Part(Id)
						ON DELETE CASCADE );` )
	return err
}

// Get all configurations from database
func getConfigurations(db *sql.DB) []ConfigurationData {
	var configurations []ConfigurationData
	rows, err := db.Query(`SELECT * from Configurations;`)
	defer rows.Close()
	for rows.Next() {
		var configurationData ConfigurationData
		configurationData, err = sqlRowToConfigurationData(db, rows)
		if err == nil {
			configurations = append(configurations, configurationData)
		} else {
			fmt.Println(err)
		}
	}
	return configurations
}

// Get specific configuration from database
func getConfiguration(db *sql.DB, id string) ConfigurationData {
	rows, err := db.Query("SELECT * FROM Configurations WHERE Id=?", id)
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
func insertConfiguration(db *sql.DB, configurationData ConfigurationData) error {
	result, err := db.Exec(`INSERT INTO Configurations(
						Name,
						Description,
						Price,
						AdditionalPrice,
						AdditionalPriceDesc
						)
				VALUES(?, ?, ?, ?, ?)`, configurationData.Name,
		configurationData.Description,
		strconv.Itoa(int(configurationData.Price)),
		strconv.Itoa(int(configurationData.AdditionalPrice)),
		configurationData.AdditionalPriceDesc)
	if err != nil {
		return err
	}
	// TODO: this does not seem to contain the stored ID
	var id int64
	id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	configurationData.Id = strconv.Itoa(int(id))
	for _, part := range configurationData.Parts {
		_, err = db.Exec(`INSERT INTO ConfigurationPart(Configuration_Id, Part_Id) VALUES(?,?)`, configurationData.Id, part.Id )
		if err != nil {
			return err
		}
	}
	return err
}

// Update the configuration in database. Configuration is identified by ConfigurationData.Id, the rest is the new data.
func updateConfiguration(db *sql.DB, configuration ConfigurationData) error {
	rows, err := db.Query(`SELECT Part_Id FROM ConfigurationPart WHERE Configuration_Id = ?`, configuration.Id )
	if err != nil {
		return err
	}
	defer rows.Close()
	var allOldIds []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		allOldIds = append(allOldIds, id)
	}
	var updateStatement string = `UPDATE Configurations 
				SET 
					Name = '` + configuration.Name  + `',` +
					`Description = '` + configuration.Description + `',` +
					`Price = ` + strconv.Itoa(int(configuration.Price)) + `,` +
					`AdditionalPrice = ` + strconv.Itoa(int(configuration.AdditionalPrice)) + `,` +
					`AdditionalPriceDesc = '` + configuration.AdditionalPriceDesc + `' ` +
				` WHERE Id = '` + configuration.Id + `'`
	var allNewIds []string
	for _, part := range configuration.Parts {
		allNewIds = append(allNewIds, part.Id)
	}
	var insertStatement string = `INSERT INTO ConfigurationPart (Configuration_id, Part_Id) VALUES `
	var insertCount int = 0
	for _, id := range allNewIds {
		if ! slices.Contains( allOldIds, id ) {
			if insertCount > 0 {
				insertStatement = insertStatement + `,`
			}
			insertStatement = insertStatement + `(` + configuration.Id + `,` + id + `)`
			insertCount++
		}
	}
	var deleteStatement string = `DELETE FROM ConfigurationPart WHERE ( Configuration_Id, Part_Id ) IN (`
	var deleteCount int = 0 
	for _, id := range allOldIds {
		if ! slices.Contains(allNewIds, id) {
			if deleteCount > 0 {
				deleteStatement = deleteStatement + `,`
			}
			deleteStatement = deleteStatement + `(` + configuration.Id + `,` + id + `)`
		}
	}
	if deleteCount > 0 {
		deleteStatement = deleteStatement + `);`
	}

	var statement = `BEGIN TRANSACTION;`
	statement = statement + updateStatement + `;`
	if deleteCount > 0 {
		statement += deleteStatement + `;`
	}
	if insertCount > 0 {
		statement += insertStatement + `;`
	}
	statement += `COMMIT;`
	_, err = db.Exec( statement )
	return err
}

func deleteConfiguration(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM Configurations WHERE Id = ?", id )
	return err
}

// Helper function to convert database query result to ConfigurationData
func sqlRowToConfigurationData(db *sql.DB, rows *sql.Rows) (ConfigurationData, error) {
	var configurationData ConfigurationData
	var name, description, additionalPriceDesc, id *string
	var price, additionalPrice *int32
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
		var parts *sql.Rows
		parts, err =  db.Query("SELECT Part_Id FROM ConfigurationPart WHERE Configuration_Id = ?", configurationData.Id)
		if err != nil {
			return ConfigurationData{}, err
		}
		for parts.Next() {
			var partId string
			parts.Scan(&partId)
			part, err := getPart(db, partId)
			if err != nil {
				return ConfigurationData{}, err
			}
			configurationData.Parts = append(configurationData.Parts, part)
		}
		return configurationData, nil
	}

	return configurationData, errors.New("No data available in query.")
}
