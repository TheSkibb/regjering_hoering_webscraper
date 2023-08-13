package main

import (
	"fmt"
    "github.com/gocolly/colly"
    "strings"
    "encoding/json"
    //"os"
    //"log"
    "io/ioutil"
)

type HoringResult struct{
    Result []Horing
}

type Horing struct {
    Date string
    Deadline string
    Department string
    Excerpt string
    Horingsbrev string
    Horingsnotat_url string
    Horingssvar []Horingssvar
    Id string
    Status string
    Title string
    Horings_type string
    Url string
}

type Horingssvar struct {
    Header string
    HoringsTitle string
    Pdf_link string
    Text string
}

func main() {
    scrapeResult := HoringResult{
        Result: []Horing{},
    }
    fmt.Println("starting scrape........")
    scrapeMainPage(&scrapeResult)
    for i, s := range scrapeResult.Result{
        fmt.Println(i, s.Url)
        scrapHoringPage(s.Url, i, &scrapeResult)
    }
    fmt.Println("scrape complete")
    fmt.Println("******************")
    fmt.Println("results")
    /*for i, s := range scrapeResult.result{
        fmt.Println(i, s)
    }*/

    file, err := json.MarshalIndent(scrapeResult, "", " ")

    fmt.Println(scrapeResult)
    if err != nil {
        fmt.Println(err.Error())
    }

	_ = ioutil.WriteFile("output.json", file, 0644)
}

func scrapeMainPage(scrapeResult *HoringResult){
    c := colly.NewCollector(
        colly.AllowedDomains("www.regjeringen.no"),
    )

    // Find and print all links
    c.OnHTML("div.results ul.listing a[href]", func(e *colly.HTMLElement) {
        scrapedHoring := Horing{}
        link := e.Attr("href")
        scrapedHoring.Url = link
        scrapeResult.Result = append(scrapeResult.Result, scrapedHoring)
    })

    c.OnError(func(r *colly.Response, err error) {
        fmt.Printf("error while scraping: %s\n", err.Error())
    })
    c.Visit("https://www.regjeringen.no/no/dokument/hoyringar/id1763/?ownerid=750&term=")
}

func scrapHoringPage(url string, index int, scrapeResult *HoringResult){
    c := colly.NewCollector(
        colly.AllowedDomains("www.regjeringen.no"),
    )

    //date
    c.OnHTML("span.date", func(e *colly.HTMLElement) {
        scrapeResult.Result[index].Date = strings.ReplaceAll(e.Text, "Dato: ", "")
    })

    //department string
    c.OnHTML("div.content-owner-dep a[href]", func(e *colly.HTMLElement) {
        scrapeResult.Result[index].Department = e.Text
    })

    //excerpt string
    c.OnHTML("div.article-ingress", func(e *colly.HTMLElement) {
        scrapeResult.Result[index].Excerpt = e.Text
    })

    //horingsbrec string
    c.OnHTML("div.factbox div#horingsbrev", func(e *colly.HTMLElement) {
        scrapeResult.Result[index].Horingsbrev = e.Text
    })

    //horingsnotat_url string
    c.OnHTML("div#horingsnotater a[href]", func(e *colly.HTMLElement) {
        scrapeResult.Result[index].Horingsnotat_url = e.Attr("href")
    })

    //horingssvar []horingssvar

    //id string
    c.OnRequest(func(r *colly.Request) {
        scrapeResult.Result[index].Id = strings.Split(url, "/")[4]
    })

    //title string
    c.OnHTML("h1", func(e *colly.HTMLElement) {
        scrapeResult.Result[index].Title = e.Text
    })

    //horings_type string
    c.OnHTML("div.article-info span.type", func(e *colly.HTMLElement) {
        scrapeResult.Result[index].Horings_type = e.Text
    })

    //deadline and status
    c.OnHTML("div.horing-meta p", func(e *colly.HTMLElement) {
        // check if element is status or deadline
        if strings.Contains(e.Text, "Status") {
            scrapeResult.Result[index].Status = strings.ReplaceAll(e.Text, "Status:", "")
        } else {
            scrapeResult.Result[index].Deadline = strings.ReplaceAll(e.Text, "Høringsfrist: ", "")
        }
    })

    c.OnError(func(r *colly.Response, err error) {
        fmt.Printf("error while scraping: %s\n", err.Error())
    })
    c.Visit("https://www.regjeringen.no" + url)
}
