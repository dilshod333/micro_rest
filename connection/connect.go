package connection

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)


const (
	host="localhost"
	password="dilshod"
	port=5432
	dbname="practise"
	user="postgres"
)

func Initialize() (*sql.DB, error){
	dbInfo := fmt.Sprintf("host=%s password=%s port=%d dbname=%s user=%s sslmode=disable", host, password, port, dbname, user)
	db, err := sql.Open("postgres", dbInfo)

	if err != nil{
		log.Fatal("Error while oppennning...", err)
	}

	err = db.Ping()
	if err != nil{
		log.Fatal("Error...", err)
	}

	return db, nil 
}