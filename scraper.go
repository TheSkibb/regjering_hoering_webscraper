package main

import (
	"fmt"

    "github.com/gocolly/colly"
)

type horingResult struct{
    result []horing
}

type horing struct {
    date string
    deadline string
    department string
    excerpt string
    horingsbrec string
    horingsnotat_url string
    horingssvar []horingssvar
    id string
    status string
    horings_type string
    url string
}

type horingssvar struct {
    header string
    horingsTitle string
    pdf_link string
    text string
}


func main() {
    links := []string{}
    scrapeMainPage(&links)
    fmt.Println(links)
    fmt.Println(len(links))
}

func scrapeMainPage(links *[]string){
    c := colly.NewCollector(
        colly.AllowedDomains("www.regjeringen.no"),
    )

    // Find and print all links
    c.OnHTML("div.results a[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        *links = append(*links, link)
    })

    c.OnError(func(r *colly.Response, err error) {
        fmt.Printf("error while scraping: %s\n", err.Error())
    })
    c.Visit("https://www.regjeringen.no/no/dokument/hoyringar/id1763/?ownerid=750&term=")
}

func scrapHoringPage(url string){
    c := colly.NewCollector(
        colly.AllowedDomains("www.regjeringen.no"),
    )

    // Find and print all links
    c.OnHTML("div.results a[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        *links = append(*links, link)
    })

    c.OnError(func(r *colly.Response, err error) {
        fmt.Printf("error while scraping: %s\n", err.Error())
    })
    c.Visit("https://www.regjeringen.no/no/dokument/hoyringar/id1763/?ownerid=750&term=")
}
