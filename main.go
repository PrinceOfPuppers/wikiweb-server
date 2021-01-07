package main

import (
	//"bytes"
	//"unicode/utf8"
	"fmt"
	"net/http"
	"io/ioutil"
	//"regexp"
	//"net/url"
	//"container/list"
	//"os"
	//"golang.org/x/net/html"
	//"golang.org/net/html/atom"
	//"strings"
	//log "github.com/llimllib/loglevel"
	"github.com/PrinceOfPuppers/wikiweb/wikiweb"
	"github.com/PrinceOfPuppers/wikiweb/wikiweb/data"
)

func getPage(link string) {
	resp, err := http.Get(link)
	if err!=nil {
		println(err)
		return
	}
	defer resp.Body.Close()


	bodyBytes, _ := ioutil.ReadAll(resp.Body)


	links := wikiweb.GetLinks(bodyBytes)



	for i, link := range links.Links {
		fmt.Printf("%v: %v\n",i,link.URL)
	}
}

func main(){
	//db := startDb()
	//db.initalize()
	//db.newUser("testing123")
//
	//rows,_ := db.db.Query("SELECT * FROM auth")
	//defer rows.Close()
	//var userID string
	//var username string
	//var time int
//
	//for rows.Next(){
	//	rows.Scan(&userID,&username,&time)
	//	println("test",userID,username,time)
	//}
	getPage("https://en.wikipedia.org/wiki/Philosophy")
	db := data.StartDb()
	//getPage("https://en.wikipedia.org/wiki/Electron")
	//db.stopDb()
}