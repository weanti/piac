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
var save_configurationTemplate =  template.Must( template.ParseFiles("templates/save_configuration.html"))

func createDatabaseTables(db *sql.DB) error {
	err := createPartsTable( db )
	if err != nil {
		return err
	}
	err = createConfigurationsTable( db )
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
		fmt.Println("Creating new database.")
		err := createDatabaseTables( db )
		if err != nil {
			fmt.Println("Failed to create database: ", err)
			os.Exit(1)
		}
		// create file
	} else {
		fmt.Println("Opening existing database.")
	}
	fileServer := http.FileServer(http.Dir("./images"))
	http.Handle("/images/", http.StripPrefix("/images/", fileServer))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "templates/index.html") })
	http.HandleFunc("/templates/parts.html", func(w http.ResponseWriter, r *http.Request) { partsHandler(w, r, db) } )
	http.HandleFunc("/templates/edit_part.html", func( w http.ResponseWriter, r *http.Request) { editPartHandler(w, r, db) } )
	http.HandleFunc("/templates/configurations.html", func(w http.ResponseWriter, r *http.Request) {configurationHandler(w, r, db) } )
	http.HandleFunc("/templates/edit_configuration.html", func(w http.ResponseWriter, r *http.Request) { editConfigurationHandler(w, r, db) } )
	http.HandleFunc("/templates/save_configuration.html", func(w http.ResponseWriter, r *http.Request) { saveConfigurationHandler(w, r, db) } )
	http.ListenAndServe(":8080", nil)
}
