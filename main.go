package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Ra struct {
	Id   int64  `json:"id"`
	Ra   string `json:"ra"`
	Geom string `json:"geom"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "geoapp"
	password = "geoapp"
	dbname   = "geoapp"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
		panic(err)

	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	rows, err := db.Query("SELECT id, ra, ST_AsGeoJSON(geom) FROM ras")
	if err != nil {
		log.Fatal(err)
	}

	var data []Ra

	for rows.Next() {
		var ra Ra
		rows.Scan(&ra.Id, &ra.Ra, &ra.Geom)
		data = append(data, ra)
	}

	dataBytes, _ := json.MarshalIndent(data, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(dataBytes)

	defer rows.Close()
	defer db.Close()
}

func main() {
	http.HandleFunc("/", GETHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
