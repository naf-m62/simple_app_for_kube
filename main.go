package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"simple_app_for_kube/cmd/apiserver"
	"simple_app_for_kube/cmd/database"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "cfg/config.yaml", "path to config file")
	flag.Parse()

	var config = new(apiserver.Config)
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("error reading config file: ", err)
	}

	if err := yaml.Unmarshal(configData, config); err != nil {
		log.Fatal("Unmarshal error:", err)
	}

	var db *sql.DB
	if db, err = connDB(config.Db); err != nil {
		log.Fatal("connection to db failed, error:", err)
	}

	apiserver.Start(config, db)
}

func connDB(config *database.Config) (db *sql.DB, err error) {
	var dsn = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.Username, config.Password, config.DBName)

	if db, err = sql.Open("postgres", dsn); err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}
