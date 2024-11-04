package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "asset_api"
	password = "-ETvslsBIda(3hTy"
	hostname = "leap.crossnet.co.id"
	port     = "3306"
	dbname   = "asset"
	// username = "asset_admin"
	// password = "JkRiCWLgkf(k.Qo2"
	// hostname = "leap.crossnet.co.id"
	// port     = "3306"
	// dbname   = "asset"
	// username = "root"
	// password = "root"
	// hostname = "localhost"
	// port     = "3307"
	// dbname   = "asset"
)

var OpenCon = 0

func DbClose(con *sql.DB) {
	con.Close()
	OpenCon -= 1
	fmt.Printf("%c%c", 13, 13)
	fmt.Print("Open Con <" + strconv.Itoa(OpenCon) + ">")
}

// func DbConnection() (*sql.DB, error) {
// 	connectionString := username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + dbname + "?parseTime=true"
// 	db, err := sql.Open("mysql", connectionString)
// 	if err != nil {
// 		log.Printf("Error %s when opening DB\n", err)
// 		return nil, err
// 	}
// 	err = db.Ping()

// 	if err != nil {
// 		fmt.Println("Err :", err)
// 		fmt.Println("titdak dapat membuka koneksi")
// 		return nil, err
// 	}

// 	db.SetMaxOpenConns(10)
// 	db.SetMaxIdleConns(10)
// 	db.SetConnMaxLifetime(time.Minute * 1)

// 	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelfunc()
// 	err = db.PingContext(ctx)
// 	if err != nil {
// 		log.Printf("Errors %s pinging DB", err)
// 		return nil, err
// 	}
// 	//log.Printf("Connected to DB %s successfully\n", dbname)
// 	OpenCon += 1
// 	fmt.Printf("%c%c", 13, 13)
// 	fmt.Print("Open Con <" + strconv.Itoa(OpenCon) + ">")
// 	return db, nil
// }

func DbConnection() (*sql.DB, error) {
	// connectionString := username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + dbname + "?parseTime=true"
	dsn := mysql.Config{
		User:                 username,
		Passwd:               password,
		AllowNativePasswords: true,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", hostname, port),
		DBName:               dbname,
		TLSConfig:            "skip-verify",
		MultiStatements:      false,
		Params: map[string]string{
			"charset": "utf8",
			// "parsetime": "true",
		},
	}
	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}
	err = db.Ping()

	if err != nil {
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
