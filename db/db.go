package db

import (
	"fmt"
	"log"

	"github.com/erfanfs10/Blog-Echo/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func ConnectDb() {

	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbUser := utils.GetEnv("DB_USER", "root")
	dbPassword := utils.GetEnv("DB_PASSWORD", "12345")
	dbName := utils.GetEnv("DB_NAME", "blog")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("Can not connect to DB")
		log.Fatalln(err)
	}
	fmt.Println("Connected to DB")
	DB = db

	DB.MustExec(schemas)
}
