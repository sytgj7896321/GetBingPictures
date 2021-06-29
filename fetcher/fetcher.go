package fetcher

import (
	"bufio"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	rateLimiter = time.Tick(200 * time.Millisecond)
	ProxyAdd    string
)

func Fetch(link string) ([]byte, error) {
	<-rateLimiter
	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	if ProxyAdd != "" {
		transport := &http.Transport{
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse(ProxyAdd)
			},
		}
		client = &http.Client{Transport: transport}
	}
	random := browser.Random()
	req.Header.Set("User-Agent", random)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}
	newReader := bufio.NewReader(resp.Body)
	e := determineEncoding(newReader)
	utf8Reader := transform.NewReader(newReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
