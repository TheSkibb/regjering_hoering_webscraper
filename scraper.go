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
    Url string
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
        for in, svar := range scrapeResult.Result[i].Horingssvar{
            if svar.Pdf_link == "" {
                scrapeHoringssvar(svar.Url, &scrapeResult, i, in)
                fmt.Println("scraping horingssvar")
            } else {
                fmt.Println("link already")
            }
            if svar.Pdf_link != "" && svar.Pdf_link[0] == 47 {
                svar.Pdf_link = "https://regjeringen.no" + svar.Pdf_link
            }
        }
    }
    fmt.Println("scrape complete")
    fmt.Println("******************")
    fmt.Println("results")

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
    c.OnHTML("div#horingssvar ul.link-list a[href]", func(e *colly.HTMLElement) {
        maxSvar := 10
        if len(scrapeResult.Result[index].Horingssvar) < maxSvar {
            url := e.Attr("href")
            svar := Horingssvar{
                Url: url,
            }

            if strings.Contains(url, ".pdf?uid"){
                svar.Pdf_link = "https://regjeringen.no" + url
                svar.Header = "Svar fra " + e.Text
            }

            scrapeResult.Result[index].Horingssvar = append(scrapeResult.Result[index].Horingssvar, svar)
        }
    })

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

func scrapeHoringssvar(url string, scrapeResult *HoringResult, horinIndex, svarIndex int){
    c := colly.NewCollector(
        colly.AllowedDomains("www.regjeringen.no"),
    )

    // Header
    c.OnHTML("header.article-header", func(e *colly.HTMLElement) {
        *&scrapeResult.Result[horinIndex].Horingssvar[svarIndex].Header = strings.TrimSpace(strings.ReplaceAll(e.Text, "\n", ""))
    })

    // Text
    c.OnHTML("div.article-body p", func(e *colly.HTMLElement) {
        *&scrapeResult.Result[horinIndex].Horingssvar[svarIndex].Text = e.Text
    })

    // Pdf_link
    c.OnHTML("div.hearing-answer ul.link-list a[href]", func(e *colly.HTMLElement) {
        *&scrapeResult.Result[horinIndex].Horingssvar[svarIndex].Pdf_link = e.Attr("href")
    })

    c.OnError(func(r *colly.Response, err error) {
        fmt.Printf("error while scraping: %s\n", err.Error())
    })

    c.Visit("https://www.regjeringen.no" + scrapeResult.Result[horinIndex].Url + url)
}
