package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "root"
	hostname = "localhost"
	port     = "3307"
	dbname   = "asetmanajemen"
)

var OpenCon = 0

func DbClose(con *sql.DB) {
	con.Close()
	OpenCon -= 1
	fmt.Printf("%c%c", 13, 13)
	fmt.Print("Open Con <" + strconv.Itoa(OpenCon) + ">")
}

func DbConnection() (*sql.DB, error) {
	connectionString := username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + dbname + "?parseTime=true"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}
	err = db.Ping()

	if err != nil {
		fmt.Println("Err :", err)
		fmt.Println("titdak dapat membuka koneksi")
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 1)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	//log.Printf("Connected to DB %s successfully\n", dbname)
	OpenCon += 1
	fmt.Printf("%c%c", 13, 13)
	fmt.Print("Open Con <" + strconv.Itoa(OpenCon) + ">")
	return db, nil
}
