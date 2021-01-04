package main

import (
	"bytes"
	//"unicode/utf8"
	//"fmt"
	//"net/http"
	//"io/ioutil"
	"regexp"
	"net/url"
	"sort"
	//"container/list"
	//"os"
	//"golang.org/x/net/html"
	//"golang.org/net/html/atom"
	//"strings"
	//log "github.com/llimllib/loglevel"
)

const linkBase = "https://en.wikipedia.org"

var filterStrings = []string{"Wikipedia:","Special:","Help:","Portal:","Talk:","Category:","(identifier)","Book:","Main_Page",".svg",".png",".jpg",".jpeg","(disambiguation)","Template:","Template_talk:"}

var filters = loadFilters(filterStrings)

func loadFilters(reStrings []string) []*regexp.Regexp{
	filters := make([]*regexp.Regexp,len(reStrings))

	for i,str := range reStrings{
		filters[i],_ = regexp.Compile(str)
	}
	return filters
}


type link struct{
	linkEnd []byte

	url string
	name string
}
func newLink(linkEnd []byte) *link {
	linkEndFormatted,_ := url.QueryUnescape(string(linkEnd))
	name := linkEndFormatted[6:]
	url := linkBase+linkEndFormatted
	l := link{linkEnd: linkEnd, url:url, name: name}

	return &l
}



type linkList struct{
	links []*link 
}
func newLinkList() *linkList{
	links := make([]*link,0)
	linkList := linkList{links:links}
	return &linkList
}
func (ll *linkList) appendNew(linkEnd []byte){
	ll.links = append(ll.links,newLink(linkEnd))
}
func (ll *linkList) inLinks(linkEnd []byte) bool{
	inLinks := false
	for _,otherLink := range ll.links{
		if bytes.Equal(otherLink.linkEnd,linkEnd) {
			inLinks = true
			continue
		}
	}
	return inLinks
}
// comparison used to alphabetize links based on the name of the artilce they link to
func (ll  *linkList) Less(i, j int) bool{

	l1 := ll.links[i]
	l2 := ll.links[j]

	minLen := 0

	l1Len := len(l1.name)
	l2Len := len(l2.name)

	if l1Len < l2Len {
		minLen = l1Len
	}else{
		minLen = l2Len
	}
	
	for i := 0; i < minLen; i++ {
		if l1.name[i] == l2.name[i]{
			 continue
		}
		if l1.name[i] < l2.name[i]{
			return true
		}
		return false
	}

	// return smaller name if all letters overlap
	if l1Len < l2Len {
		return true
	}
	return false	
}
func (ll *linkList) Len() int  {
	return len(ll.links)
}
func  (ll *linkList) Swap(i,j int)  {
	ll.links[i],ll.links[j] = ll.links[j],ll.links[i]
}


func inFilters(filters []*regexp.Regexp, linkEnd []byte) bool{
	
	for _,filter := range filters {
		if filter.Match(linkEnd) {
			return true
		}
	}
	return false
}



func getLinks(html []byte) *linkList{
	re,_ := regexp.Compile("href=\"/wiki/.*?\"")

	matches := re.FindAll(html,-1)


	links := newLinkList()

	var linkEnd []byte

	for _,match := range matches {

		linkEnd = match[6:len(match)-1]

		// remove subsections (everything past #)
		re,_ = regexp.Compile("#.*")
		linkEnd = re.ReplaceAllLiteral(linkEnd,[]byte(""))
		
		if inFilters(filters,linkEnd){continue}
		
		// checks if already in links
		if links.inLinks(linkEnd){continue}

		links.appendNew(linkEnd)
	}
	sort.Sort(links)
	return links
}