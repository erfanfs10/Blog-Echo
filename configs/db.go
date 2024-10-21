package configs

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var DB *sqlx.DB

func ConnectDb() {
	db, err := sqlx.Connect("mysql", "root:12345@/sn?parseTime=true")
	if err != nil {
		fmt.Println("Can not connect to DB")
		log.Fatalln(err)
	}
	fmt.Println("Connected to DB")
	DB = db
}
