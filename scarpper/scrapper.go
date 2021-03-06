package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id          string
	title       string
	companyName string
	location    string
}

//Scrape Indeed by a term
func Scrape(term string) {
	var baseURL string = "https://www.indeed.com/jobs?q=" + term + "&limit=50"
	fmt.Println("start")
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseURL, 0)
	for i := 0; i < totalPages; i++ {
		go getpage(i, baseURL, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJob := <-c
		jobs = append(jobs, extractedJob...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func getpage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := url + "&start=" + strconv.Itoa(page*50)
	fmt.Println(pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := CleanString(card.Find("h2>span").Text())
	companyName := CleanString(card.Find(".companyName").Text())
	location := CleanString(card.Find("div pre").Text())
	c <- extractedJob{
		id:          id,
		title:       title,
		companyName: companyName,
		location:    location}
}

//CleanString Cleans a string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string, prevLast int) (lastPage int) {
	baseURL := url
	if prevLast != 0 {
		url = baseURL + "&start=" + strconv.Itoa((prevLast-1)*50)
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

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)

	defer w.Flush()

	headers := []string{"ID", "Title", "CompanyName", "Location"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://www.indeed.com/viewjob?jk=" + job.id, job.title, job.companyName, job.location}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
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
