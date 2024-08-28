package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var visitedUrl = make(map[string]struct{})
var resultFile *os.File
var countChan chan struct{}
var wg sync.WaitGroup

func main() {
	countChan = make(chan struct{}, 10)
	var err error
	resultFile, err = os.OpenFile("result.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data, err := os.ReadFile("site.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	siteStr := strings.ReplaceAll(string(data), "\r", "")
	sites := strings.Split(siteStr, "\n")
	for _, site := range sites {
		u, err := url.Parse(site)
		if err != nil {
			fmt.Println(site + " url 格式不正确")
			continue
		}
		if _, ok := visitedUrl[u.Host]; ok {
			continue
		}
		visitedUrl[u.Host] = struct{}{}
		wg.Add(1)
		countChan <- struct{}{}
		go fetchCompanySite(site)
	}
	wg.Wait()
}
func fetchCompanySite(requestUrl string) error {
	defer func() {
		<-countChan
		wg.Done()
	}()
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	replacer := strings.NewReplacer(" ", "", "\n", "", " ", "", "\r", "")
	title := replacer.Replace(doc.Find("title").Text())
	fmt.Println(requestUrl, title)
	_, err = resultFile.WriteString(fmt.Sprintf("%s,%s\n", requestUrl, title))
	if err != nil {
		return err
	}
	return nil
}
