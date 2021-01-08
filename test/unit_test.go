package test

import (
	"fmt"
	"testing"
	"os"
	"github.com/PrinceOfPuppers/wikiweb/data"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
)
func setLogging(){
	formatter := runtime.Formatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp: true,
	}}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.WithFields(log.Fields{
		"file": "main.go",
	}).Info("Running Tests...")
}

func makeBlankDb(path string) *data.DataBase{
	db:=data.StartDb(path)
	db.Initalize()
	return db
}
func closeAndDeleteDb(db *data.DataBase, path string){
	db.StopDb()
	err := os.Remove(path) 
	if err!=nil{
		log.Fatal("Cannot Delete Database")
	}
}

func TestCreateDb(t *testing.T) {


	dbPath := "./testing.db"
	db := makeBlankDb(dbPath)
	defer closeAndDeleteDb(db,dbPath)

	db.NewUser("jimbo")
	db.NewUser("Robert'); DROP TABLE auth;--")//bobby tables test

	db.AddScore("jimbo",1,10,300)
	db.AddScore("jimbo",4,20,500)
	db.AddScore("Robert'); DROP TABLE auth;--",4,3,200)

	authTable := db.GetTableContents("auth")
	week1ScoreTable := db.GetTableContents("w4 ")

	for _,row := range table {
		fmt.Printf("%#v\n", row)
	}
	t.Fatal("")
}
