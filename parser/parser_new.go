package parser

import (
	"GetBingPictures/fetcher"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	Path                = "wallpapers"
	bingSrc             = "https://bing.ioliu.cn/"
	bingGetLatestNum    = "body > div.page > span"
	bingTarget          = "https://cn.bing.com/th?id=OHR."
	selectorDate        = "body > div.container > div:nth-child(ReplaceHere) > div > div.description > p.calendar > em"
	selectorName        = "body > div.container > div:nth-child(ReplaceHere) > div > a"
	selectorDescription = "body > div.container > div:nth-child(ReplaceHere) > div > div.description > h3"
)

type BingPic struct {
	Date        string
	Name        string
	Description string
	Url         string
	OldUrl      string
}

func Parser(tid int, fp *os.File) {
	var picName string
	var picUrl string
	bingPicList := make([]BingPic, 12)
	log.SetOutput(fp)
	log.SetPrefix("[GetBingWallpaperTool]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	result, _ := fetcher.Fetch(bingSrc + "?p=" + strconv.Itoa(tid))
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(result)))

	for i := 1; i <= 12; i++ {
		selectors := changeSelectors(i, selectorDate, selectorName, selectorDescription)
		bingPicList[i-1].getSelectors(dom, selectors...)
	}

	for _, b := range bingPicList {
		picName = b.Name + "_UHD.jpg"
		picUrl = b.Url
		resp, err := fetchBody(picUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Get Image Error %d\n", err)
			continue
		}
		if resp.StatusCode == 404 {
			picName = b.Name + "_1920x1080.jpg"
			picUrl = b.OldUrl
			resp, err = fetchBody(picUrl)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Get Image Error %d\n", err)
				continue
			}
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Image Read Error %s\n", err)
			resp.Body.Close()
			continue
		}
		err = ioutil.WriteFile(Path+"/"+b.Date+"_"+picName, data, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "IO Write Error %s\n", err)
			resp.Body.Close()
			continue
		}
		fmt.Printf("%s download completed\n", b.Name)
		resp.Body.Close()
	}

}

func FetchLatestPageNum() (int, error) {
	result, _ := fetcher.Fetch(bingSrc + "1")
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(result)))
	lastNum := selectorParser(bingGetLatestNum, dom)
	lastNum = strings.Replace(lastNum, "1 / ", "", -1)
	return strconv.Atoi(lastNum)
}

func (b *BingPic) getSelectors(dom *goquery.Document, selectors ...string) {
	b.Date = selectorParser(selectors[0], dom)
	b.Name = selectorParser(selectors[1], dom)
	b.Description = selectorParser(selectors[2], dom)
	b.Url = bingTarget + b.Name + "_UHD.jpg"
	b.OldUrl = bingSrc + "photo/" + b.Name + "?force=download"
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

func fetchBody(link string) (*http.Response, error) {
	<-fetcher.RateLimiter
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	transport := &http.Transport{
		Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse(fetcher.ProxyAdd)
		},
	}
	client := &http.Client{Transport: transport}
	random := browser.Random()
	req.Header.Set("User-Agent", random)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return resp, nil
}
