package main

import (
	"database/sql"
	"fmt"
	"html/template"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
)

var partsTemplate = template.Must( template.ParseFiles("templates/parts.html") )
var edit_partTemplate = template.Must( template.ParseFiles("templates/edit_part.html"))
var configurationTemplate = template.Must( template.ParseFiles( "templates/configurations.html"))
var edit_configurationTemplate =  template.Must( template.ParseFiles("templates/edit_configuration.html"))

func createDatabaseTables(db *sql.DB) error {
	_, err := db.Query(`CREATE TABLE Parts(
						Brand TEXT,
						Model TEXT,
						Description TEXT,
						Count INTEGER,
						Price INTEGER,
						Image BLOB,
						InternalId INTEGER )` )
	if err != nil {
		return err
	}
	_, err = db.Query(`CREATE TABLE Configurations(
						Name TEXT,
						Description TEXT,
						Price INTEGER,
						AdditionalPrice INTEGER,
						AdditionalPriceDesc TEXT,
						InternalId INTEGER)` )
						// TODO: somehow store the references here
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("use <dbfile>")
		os.Exit(1)
	}
	var dbFile = os.Args[1]
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	_, err = os.Stat( dbFile )
	if os.IsNotExist(err) {
		fmt.Println("Creating new database")
		err := createDatabaseTables( db )
		if err != nil {
			fmt.Println("Failed to create database: ", err)
			os.Exit(1)
		}
		// create file
	} else {
		fmt.Println("Opening existing database.")
		_, err = db.Query("INSERT INTO Parts ( Brand, Model, Count, Price ) VALUES ('Shimano', 'A070',1,4000);")
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Query("INSERT INTO Configurations ( Name, Description, Price ) VALUES ('Csepel Rapid', '50/34, prowheel', 160000)")
		if err != nil {
			fmt.Println(err)
		}
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	http.HandleFunc("/templates/parts.html", partsHandler)
	http.HandleFunc("/templates/edit_part.html", editPartHandler)
	http.HandleFunc("/templates/configurations.html", configurationHandler)
	http.ListenAndServe(":8080", nil)
}
