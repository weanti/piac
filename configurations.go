package main

import( "database/sql"
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
						Name TEXT,
						Description TEXT,
						Price INTEGER,
						AdditionalPrice INTEGER,
						AdditionalPriceDesc TEXT,
						InternalId INTEGER)` )
						// TODO: somehow store the references here
	return err
}

func configurationHandler( writer http.ResponseWriter, request *http.Request ) {

	examplePartData := PartData {
	    Brand : "Shimano",
	    Model : "FD-A070",
	    Description : "black, max tooth 50",
	    Count : 1,
	    Price : 4000,
	    Image : "",
	    InternalId : 1,
	}
	exampleConfigurationData := ConfigurationData {
		Name : "Csepel Rapid",
		Description : "flatbar, 2x8, integrated crankset, 50/34, 35-622 Schwalbe",
		Price : 162000,
		AdditionalPrice : 0,
		AdditionalPriceDesc : "",
		Parts : []PartData{ examplePartData },
		InternalId : 1,
	}
	configurationTemplate.Execute( writer, []ConfigurationData { exampleConfigurationData } )
	// TODO: get ID from URL
	// if empty then all fields shall be empty and a STORE operation is to be performed
	// if not empty then query all fields (is this redundant? maybe recevie everything from query?)
	//	in this case an update need to be performed
//	configurationTemplate.Execute(writer, examplePartData)
}
