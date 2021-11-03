package storage

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"log"
	"os"
)


func CreateDb(host string) *sql.DB {
	log.Println("Connect to database: ", host)
	connect, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s", host))
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}


	return connect
}