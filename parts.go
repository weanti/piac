package main

import ( "net/http" )

type PartData struct {
    Brand string
    Model string
    Description string
    Count int16
    Price int32
    Image string
    InternalId int32
}

func partsHandler( writer http.ResponseWriter, request *http.Request ) {
	examplePartData := PartData {
	    Brand : "Shimano",
	    Model : "FD-A070",
	    Description : "black, max tooth 50",
	    Count : 1,
	    Price : 4000,
	    Image : "",
	    InternalId : 1,
	}
	examplePartData2 := PartData {
	    Brand : "Sunrace",
	    Model : "RD90",
	    Description : "max tooth 28",
	    Count : 1,
	    Price : 4000,
	    Image : "",
	    InternalId : 2,
	}
	// TODO: perform an SQL query and map the result to PartData array
	parts := []PartData{ examplePartData, examplePartData2 }
	partsTemplate.Execute(writer, parts)
}

func editPartHandler( writer http.ResponseWriter, request *http.Request ) {
	examplePartData := PartData {
	    Brand : "Shimano",
	    Model : "FD-A070",
	    Description : "black, max tooth 50",
	    Count : 1,
	    Price : 4000,
	    Image : "",
	    InternalId : 1,
	}
	// TODO: get ID from URL
	// if empty then all fields shall be empty and a STORE operation is to be performed
	// if not empty then query all fields (is this redundant? maybe recevie everything from query?)
	//	in this case an update need to be performed
	edit_partTemplate.Execute(writer, examplePartData)
}
