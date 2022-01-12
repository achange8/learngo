package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var (
	limit   int    = 50
	baseURL string = fmt.Sprintf("https://www.indeed.com/jobs?q=python&limit=%d", limit)
)

func main() {
	fmt.Println("1234")
	totalPages := getPages(baseURL, 0)
	for i := 0; i < totalPages; i++ {
		getpage(i)
	}
}

func getpage(page int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println(pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")
		fmt.Println(id)
		title := card.Find("h2>span").Text()
		fmt.Println(title)

	})
}

func getPages(url string, prevLast int) (lastPage int) {
	if prevLast != 0 {
		url = baseURL + "&start=" + strconv.Itoa((prevLast-1)*limit)
	}

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages := s.Find("a")
		pageLen := pages.Length()

		if pageLen == 3 {
			lastPage = prevLast
		} else {
			nextLast := 0
			pages.Each(func(i int, s *goquery.Selection) {
				if i == pageLen-2 {
					nextLast, _ = strconv.Atoi(s.Text())
				}
			})
			lastPage = getPages(baseURL, nextLast)
		}
	})
	return
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
