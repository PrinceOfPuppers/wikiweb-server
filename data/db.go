package data
import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"crypto/sha1"
	"time"
	"github.com/PrinceOfPuppers/wikiweb/helpers"

)

//GetUserScoreTable gets a unique, sql safe name for the user score table
func GetUserScoreTable(username string) string {
	return fmt.Sprintf("a%x",sha1.Sum([]byte(username)))
}

//DataBase is a wrapper object for an sql database
type DataBase struct {
	db *sql.DB 
	// TODO add channel for sending queries

	// common strings
	weeklyScoreTable string
}

// StartDb opens the database and returns the database object
func StartDb(path string) *DataBase{
	db, err := sql.Open("sqlite3",path)
	
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	dataBase := DataBase{ db:db, weeklyScoreTable: "w%d scores" }
	return &dataBase
}

//StopDb Closes the database
func (dataBase *DataBase) StopDb(){
	dataBase.db.Close()
}


func (dataBase *DataBase) tableExists(tableName string) bool{
	s := "SELECT count(*) FROM sqlite_master WHERE type='table' AND name= $1;"

	var count int

	rows,err := dataBase.db.Query(s,tableName)
	if err != nil {
        log.Fatal(err)
	}
	defer rows.Close()

	rows.Next()
	rows.Scan(&count)

	return count!=0
}

func (dataBase *DataBase) usernameExists(username string) bool{
	s := "SELECT COUNT(1) FROM auth WHERE username = $1;"
	var count int

	rows,err := dataBase.db.Query(s,username)
	if err != nil {
        log.Fatal(err)
	}
	defer rows.Close()

	rows.Next()
	rows.Scan(&count)
	
	return count!=0
}

//Initalize creates the inital tables that need to be in the database (should be run on first time db creation)
func (dataBase *DataBase) Initalize() {
	if !dataBase.tableExists("auth"){
		s := "CREATE TABLE auth(userID VARCHAR(32) PRIMARY KEY, username VARCHAR(255) NOT NULL, lastOnline INT);"
		_,err := dataBase.db.Exec(s)
		if err != nil {
			log.Fatal(err)
		}
	}
	// add other tables if need be
}

//NewUser adds a new user to the auth table and returns (true,userID) on success
func (dataBase *DataBase) NewUser(username string) (bool,string){

	// TODO setup rollback thing here for adding user to auth and adding user score table
	if dataBase.usernameExists(username){
		log.Fatal("Username ",username," Already Exists")
		return false,""
	}
	n := 32
	
	//b := make([]byte, n)
	//_,err := rand.Read(b)
	userID := helpers.RandString(n)
	//if err != nil {
	//	log.Fatal("Unable to generate userID: ", err)
	//	return false,""
	//}
	//
	//userID := string(b)
	uTime := int(time.Now().Unix())
	s := "INSERT INTO auth VALUES($1, $2, $3);"
	_,err := dataBase.db.Exec(s,userID, username, uTime)
	if err != nil {
		log.Fatal(err)
		return false,""
	}

	// create user score table
	uScoreTable :=GetUserScoreTable(username)
	s = fmt.Sprintf("CREATE TABLE %v(tableIndex IDENTITY(1,1) PRIMARY KEY, week INT, scoreIndex INT);",uScoreTable)

	_,err = dataBase.db.Exec(s)
	if err != nil {
		log.Fatal("Failed to create user score table: ",err)
		return false,""
	}
	
	return true,userID
}

//AddScore adds score to user score table and to master score table
func (dataBase *DataBase) AddScore(username string, week int, numLinks int, runTime int){
	// main score entry
	wScoreTable := fmt.Sprintf(dataBase.weeklyScoreTable, week)

	if !dataBase.tableExists(wScoreTable){
		s :=fmt.Sprintf("CREATE TABLE %v(scoreIndex IDENTITY(1,1) PRIMARY KEY, username VARCHAR(255), " +
			 "numLinks INT, time INT, timeSubmitted INT);",wScoreTable)
		 _,err := dataBase.db.Exec(s)
		 if err != nil {
			 log.Fatal(err)
		 }
	}

	s := fmt.Sprintf("INSERT INTO %v(username, numLinks, time, timeSubmitted) VALUES($2, $3, $4, $5);",wScoreTable)
	uTime := int(time.Now().Unix())
	res,err := dataBase.db.Exec(s,username,numLinks,runTime,uTime)
	if err != nil {
		log.Fatal(err)
	}

	scoreIndex,err:=res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	// user score table
	uScoreTable := GetUserScoreTable(username)
	s = fmt.Sprintf("INSERT INTO %v(week, scoreIndex) VALUES($2, $3);",uScoreTable)
	_,err = dataBase.db.Exec(s,week,scoreIndex)

	if err != nil {
		log.Fatal(err)
	}
}

// GetScore returns the username, number of links, time, and time submitted for a given run
func (dataBase *DataBase) GetScore(week, scoreIndex int) (string, int, int, int){
	wScoreTable := fmt.Sprintf(dataBase.weeklyScoreTable, week)

	s := "SELECT * FROM $1 WHERE scoreIndex=$2;"
	rows,err := dataBase.db.Query(s,wScoreTable,scoreIndex)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var username string
	var numLinks int
	var time int
	var timeSubmitted int
	rows.Next()
	rows.Scan(&scoreIndex, &username, &numLinks, &time, &timeSubmitted)
	
	return username,numLinks,time,timeSubmitted
}

// ValidateUser returns true if the username, userID pair match
func (dataBase *DataBase) ValidateUser(username, userID string) bool {
	s := "SELECT userID FROM auth WHERE username=$1;"
	rows,err := dataBase.db.Query(s,username)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer rows.Close()

	var realUserID string

	rows.Next()
	rows.Scan(&realUserID)

	return userID==realUserID
}

// GetTableContents is used for debugging and testing, it returns an array of all rows
func (dataBase *DataBase) GetTableContents(tableName string) [][]string{

	var tableRows [][]string
	s := fmt.Sprintf("SELECT * FROM %v;",tableName)
	rows,err := dataBase.db.Query(s)
	if err != nil {
		log.Fatal(err)
		return tableRows
	}
	defer rows.Close()
	
	cols,_:=rows.Columns()



    rowBytes := make([][]byte, len(cols))
    

	pointers := make([]interface{}, len(cols))
	
	for i := 0; i < len(rowBytes); i++ {
        pointers[i] = &rowBytes[i]
	}

    for rows.Next() {
        err = rows.Scan(pointers...)
        if err != nil {
            fmt.Println("Failed to scan row", err)
		}
		
		rowStrings := make([]string, len(cols))
        for i, b := range rowBytes {
            if b == nil {
                rowStrings[i] = "\\N"
            } else {
                rowStrings[i] = string(b)
            }
        }
		tableRows = append(tableRows,rowStrings)
    }

	return tableRows
}