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
)

func getPage(link string) {
	resp, err := http.Get(link)
	if err!=nil {
		println(err)
		return
	}
	defer resp.Body.Close()


	bodyBytes, _ := ioutil.ReadAll(resp.Body)


	links := getLinks(bodyBytes)



	for i, link := range links.links {
		fmt.Printf("%v: %v\n",i,link.url)
	}
}

func main(){
	//db := startDb()
	//db.initalize()
	//db.newUser("asdf")
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
	//getPage("https://en.wikipedia.org/wiki/Electron")
	//db.stopDb()
}