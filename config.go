package main

import (
	"github.com/BurntSushi/toml"
	"os"
	"log"
)

type Config struct {
	Port      string
	LogToFile bool
	LogFile   string
	Database  Database
	Sdt       SDT
}

type SDT struct {
	ApiUrl    string
	PartnerId string
	Password  string
}

type Database struct {
	Host     string
	Port     string
	Db       string
	User     string
	Password string
}

/*
read config from TOML file
 */
func ReadConfig(configFile string) Config {

	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Config file is missing: ", err, " : ", configFile)
	}

	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
