package parser

import (
	"GetBingPictures/fetcher"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	BingSrc             = "https://bing.ioliu.cn/?p="
	BingGetLatestNum    = "body > div.page > span"
	BingTarget          = "https://cn.bing.com/th?id=OHR."
	selectorDate        = "body > div.container > div:nth-child(ReplaceHere) > div > div.description > p.calendar > em"
	selectorName        = "body > div.container > div:nth-child(ReplaceHere) > div > a"
	selectorDescription = "body > div.container > div:nth-child(ReplaceHere) > div > div.description > h3"
)

type BingPic struct {
	Date        string
	Name        string
	Description string
	Url         string
}

func Parser(pid, wid int, fp *os.File) {
	log.SetOutput(fp)
	log.SetPrefix("[GetBingTool]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	result, _ := fetcher.Fetch(BingSrc)

	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(result)))
	for i := 1; i <= 12; i++ {
		selectors := changeSelectors(i, selectorDate, selectorName, selectorDescription)
		getSelectors(dom, selectors...)
	}

}

func FetchLatestPageNum() (int, error) {
	result, _ := fetcher.Fetch(BingSrc + "1")
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(result)))
	lastNum := selectorParser(BingGetLatestNum, dom)
	fmt.Println(lastNum)
	lastNum = strings.TrimLeft(lastNum, "1 / ")
	return strconv.Atoi(lastNum)
}

func getSelectors(dom *goquery.Document, selectors ...string) {
	for _, selector := range selectors {
		fmt.Printf("%s\n", selectorParser(selector, dom))
	}
}

func changeSelectors(i int, selectors ...string) []string {
	for id, _ := range selectors {
		selectors[id] = strings.Replace(selectors[id], "ReplaceHere", strconv.Itoa(i), -1)
	}
	return selectors
}

func selectorParser(element string, dom *goquery.Document) string {
	var s string
	dom.Find(element).Each(func(i int, selection *goquery.Selection) {
		if selection.Is("a") {
			selection, _ := selection.Attr("href")
			f := func(c rune) bool {
				if c == '/' || c == '?' {
					return true
				} else {
					return false
				}
			}
			arr := strings.FieldsFunc(selection, f)
			s = arr[1]
		} else {
			s = selection.Text()
		}
	})
	return s
}
