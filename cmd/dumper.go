package main

import (
	"database/sql"
	"fmt"
	"github.com/JamesStewy/go-mysqldump"
	"github.com/OnlyPiglet/fly/filetools"
	_ "github.com/OnlyPiglet/mysqldumper/config"
	fu "github.com/duke-git/lancet/v2/fileutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/viper"
	"log"
	"time"
)

func main() {
	inv := viper.GetDuration("db.interval")
	ticker := time.NewTicker(inv)
	defer ticker.Stop()
	for _ = range ticker.C {
		err := Dump()
		if err != nil {
			log.Println(err)
		}
	}
}

func Dump() error {
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	hostname := viper.GetString("db.host")
	port := viper.GetString("db.port")
	dbname := viper.GetString("db.name")

	maxKept := viper.GetInt("db.maxKept")

	dumpDir := "./dumps" // you should create this directory

	if !fu.IsDir(dumpDir) {
		err := fu.CreateDir(dumpDir)
		if err != nil {
			return err
		}
	}

	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", dbname) // accepts time layout string and add .sql at the end of file

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname))
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return err
	}

	// Register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return err
	}

	// Dump database to file
	resultFilename, err := dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", err)
		return err
	}
	log.Printf("File is saved to %s\n", resultFilename)

	// Close dumper and connected database
	err = dumper.Close()
	if err != nil {
		return err
	}
	return filetools.RecentFileMaxKept(dumpDir, maxKept)
}
