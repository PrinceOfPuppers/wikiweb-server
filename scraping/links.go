package scraping

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


type Link struct{
	linkEnd []byte

	URL string
	Name string
}
func newLink(linkEnd []byte) *Link {
	linkEndFormatted,_ := url.QueryUnescape(string(linkEnd))
	name := linkEndFormatted[6:]
	url := linkBase+linkEndFormatted
	l := Link{linkEnd: linkEnd, URL:url, Name: name}

	return &l
}



type LinkList struct{
	Links []*Link 
}
func NewLinkList() *LinkList{
	links := make([]*Link,0)
	linkList := LinkList{Links:links}
	return &linkList
}
func (ll *LinkList) AppendNew(linkEnd []byte){
	ll.Links = append(ll.Links,newLink(linkEnd))
}
func (ll *LinkList) InLinks(linkEnd []byte) bool{
	inLinks := false
	for _,otherLink := range ll.Links{
		if bytes.Equal(otherLink.linkEnd,linkEnd) {
			inLinks = true
			continue
		}
	}
	return inLinks
}
// comparison used to alphabetize links based on the name of the artilce they link to
func (ll  *LinkList) Less(i, j int) bool{

	l1 := ll.Links[i]
	l2 := ll.Links[j]

	minLen := 0

	l1Len := len(l1.Name)
	l2Len := len(l2.Name)

	if l1Len < l2Len {
		minLen = l1Len
	}else{
		minLen = l2Len
	}
	
	for i := 0; i < minLen; i++ {
		if l1.Name[i] == l2.Name[i]{
			 continue
		}
		if l1.Name[i] < l2.Name[i]{
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
func (ll *LinkList) Len() int  {
	return len(ll.Links)
}
func  (ll *LinkList) Swap(i,j int)  {
	ll.Links[i],ll.Links[j] = ll.Links[j],ll.Links[i]
}

func inFilters(filters []*regexp.Regexp, linkEnd []byte) bool{
	
	for _,filter := range filters {
		if filter.Match(linkEnd) {
			return true
		}
	}
	return false
}



func GetLinks(html []byte) *LinkList{
	re,_ := regexp.Compile("href=\"/wiki/.*?\"")

	matches := re.FindAll(html,-1)


	links := NewLinkList()

	var linkEnd []byte

	for _,match := range matches {

		linkEnd = match[6:len(match)-1]

		// remove subsections (everything past #)
		re,_ = regexp.Compile("#.*")
		linkEnd = re.ReplaceAllLiteral(linkEnd,[]byte(""))
		
		if inFilters(filters,linkEnd){continue}
		
		// checks if already in links
		if links.InLinks(linkEnd){continue}

		links.AppendNew(linkEnd)
	}
	sort.Sort(links)
	return links
}