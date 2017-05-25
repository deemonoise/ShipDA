package main

import (
	"net/http"
	"log"
	//"flag"
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"ShipDA/sdt"
	"github.com/jmoiron/sqlx"
)

var Cfg Config

var Sdt sdt.Sdt

func init() {
	//var configFile string
	//flag.StringVar(&configFile, "config", "", "TOML configfile")
	//flag.Parse()
	//Cfg = ReadConfig(configFile)

	Sdt = sdt.Init(Cfg.Sdt.ApiUrl, Cfg.Sdt.PartnerId, Cfg.Sdt.Password)
}

var db *sqlx.DB

func dbInit() {
	dsn := fmt.Sprint(Cfg.Database.User, ":", Cfg.Database.Password, "@tcp(", Cfg.Database.Host, ":", Cfg.Database.Port, ")/", Cfg.Database.Db)

	var err error

	//db, err = sql.Open("mysql", dsn)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbInit()
	defer db.Close()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":" + Cfg.Port, router))
}


