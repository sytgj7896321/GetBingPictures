package parser

import (
	"GetBingPictures/fetcher/iii"
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	Path                      = "wallpapers"
	bingSrc                   = "https://bing.ioliu.cn/"
	bingGetLatestNum          = "body > div.page > span"
	bingTarget                = "https://cn.bing.com/th?id=OHR."
	selectorDate              = "body > div.container > div:nth-child(ReplaceHere) > div > div.description > p.calendar > em"
	selectorName              = "body > div.container > div:nth-child(ReplaceHere) > div > a"
	selectorArtistDescription = "body > div.container > div:nth-child(ReplaceHere) > div > div.description > h3"
)

var (
	LogScanner = map[string]bool{}
)

type BingPic struct {
	Date        string
	Name        string
	Artist      string
	Description string
	Url         string
}

func Parser(tid int, rp, fp *os.File) {
	var picName string
	var picUrl string
	bingPicList := make([]BingPic, 12)
	log.SetPrefix("[GetBingWallpaperTool]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	result, _ := iii.Fetch(bingSrc + "?p=" + strconv.Itoa(tid))
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(result)))

	for i := 1; i <= 12; i++ {
		selectors := changeSelectors(i, selectorDate, selectorName, selectorArtistDescription)
		bingPicList[i-1].getSelectors(dom, selectors...)
	}

	for _, b := range bingPicList {
		log.SetOutput(fp)
		if !LogScanner[b.Name] && b.Date >= "2018-09-11" {
			picName = b.Name + "_UHD.jpg"
			picUrl = b.Url + "_UHD.jpg"
			resp, err := http.Get(picUrl)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Get Image Error %s\n", err)
				log.Printf("%s\n", err)
				continue
			}
			if resp.StatusCode == 404 {
				picName = b.Name + "_1920x1080.jpg"
				picUrl = b.Url + "_1920x1080.jpg"
				resp, err = http.Get(picUrl)
				if err != nil || resp.StatusCode == 404 {
					continue
				}
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Image Read Error %s\n", err)
				log.Printf("%s\n", err)
				resp.Body.Close()
				continue
			}
			err = ioutil.WriteFile(Path+"/"+b.Date+"_"+picName, data, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "IO Write Error %s\n", err)
				log.Printf("%s\n", err)
				resp.Body.Close()
				continue
			}
			fmt.Printf("%s download completed\n", b.Name)
			io.WriteString(rp, b.Name+" "+b.Artist+" "+b.Description+"\n")
			resp.Body.Close()
		} else {
			fmt.Printf("%s has downloaded skip\n", b.Name)
		}
	}
}

func FetchLatestPageNum() (int, error) {
	result, _ := iii.Fetch(bingSrc + "?p=" + "1")
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(result)))
	lastNum, _ := selectorParser(bingGetLatestNum, dom)
	lastNum = strings.Replace(lastNum, "1 / ", "", -1)
	return strconv.Atoi(lastNum)
}

func ScannerRecord(rp *os.File) {
	scanner := bufio.NewScanner(rp)
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Read Record Error %s\n", err)
	} else {
		for scanner.Scan() {
			stringSlice := strings.Split(scanner.Text(), " ")
			LogScanner[stringSlice[0]] = true
		}
	}
}

func (b BingPic) getSelectors(dom *goquery.Document, selectors ...string) {
	b.Date, _ = selectorParser(selectors[0], dom)
	b.Name, _ = selectorParser(selectors[1], dom)
	_, bTmp := selectorParser(selectors[2], dom)
	b.Artist = bTmp[0]
	b.Description = bTmp[1]
	b.Url = bingTarget + b.Name
}

func changeSelectors(i int, selectors ...string) []string {
	for id, _ := range selectors {
		selectors[id] = strings.Replace(selectors[id], "ReplaceHere", strconv.Itoa(i), -1)
	}
	return selectors
}

func selectorParser(element string, dom *goquery.Document) (string, []string) {
	var s string
	var sSub []string
	dom.Find(element).Each(func(i int, selection *goquery.Selection) {
		if selection.Is("a") {
			selection, _ := selection.Attr("href")
			f := func(r rune) bool {
				if r == '/' || r == '?' {
					return true
				} else {
					return false
				}
			}
			arr := strings.FieldsFunc(selection, f)
			s = arr[1]
		} else if selection.Is("h3") {
			sSub = mySplit(selection.Text())
		} else {
			s = selection.Text()
		}
	})
	return s, sSub
}

func mySplit(str string) []string {
	var strSub []string
	strSub = strings.Split(str, " (© ")
	if len(strSub) != 2 {
		strSub = strings.Split(str, "©")
		for _, s := range strSub {
			s = strings.TrimRightFunc(s, func(r rune) bool {
				if r == '（' || r == '）' {
					return true
				} else {
					return false
				}
			})
		}
		return strSub
	}
	strSub[1] = strings.TrimRight(strSub[1], ")")
	return strSub
}
