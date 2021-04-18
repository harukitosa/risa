package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

type ScrapingObject struct {
	URL        string
	AbstractEN string
	AbstractJP string
}

func main() {
	url := "https://scholar.google.com/scholar?hl=ja&as_sdt=0%2C5&q=vqd&btnG="
	// ScrapingGoogleScholar(url)
	var fetcher FetchURLList
	var abstfetcher FetchAbstruct
	fetcher = GoogleScholar{}
	abstfetcher = Arxiv{}
	str, err := fetcher.Get(url)
	if err != nil {
		panic(err)
	}
	// log.Println(str)
	for _, s := range str {
		str, _ := abstfetcher.Get(s)
		fmt.Printf("URL: %s\n %s\n", s, str)
	}
}

type FetchURLList interface {
	Get(url string) ([]string, error)
}

type GoogleScholar struct{}

func (g GoogleScholar) Get(url string) ([]string, error) {
	doc := ScrapingGetContent(url)
	var str []string
	doc.Find(".gs_rt").Each(func(i int, s *goquery.Selection) {
		band := s.Find("a")
		herf, ok := band.Attr("href")
		if !ok {
			return
		}
		// title := s.Find("i").Text()
		str = append(str, herf)
		// if strings.Contains(herf, "https://arxiv.org/") {
		// 	log.Println("this is arxiv")
		// 	ScrapingArxiv(herf)
		// }
	})
	return str, nil
}

type FetchAbstruct interface {
	Get(url string) (string, error)
}

type Arxiv struct{}

func (a Arxiv) Get(url string) (string, error) {
	doc := ScrapingGetContent(url)
	var abst string
	doc.Find(".abstract").Each(func(i int, s *goquery.Selection) {
		abst = s.Text()
	})
	return abst, nil
}

func ScrapingGetContent(url string) *goquery.Document {
	// Getリクエスト
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	det := chardet.NewTextDetector()
	detRslt, _ := det.DetectBest(buf)
	// fmt.Println(detRslt.Charset)

	bReader := bytes.NewReader(buf)
	reader, err := charset.NewReaderLabel(detRslt.Charset, bReader)
	if err != nil {
		log.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
	}
	return doc
}
