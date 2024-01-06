package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Movie struct {
	Title string
	Year  string
}

type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []Movie
}

func crawl(month, day int) {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)
	startURL := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	fmt.Println(startURL)
	c.Visit(startURL)
	
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response received", string(r.Body))
	})
	infoCollector := c.Clone()

	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileUrl := e.ChildAttr("div.lister-item-image > a", "href")
		profileUrl = e.Request.AbsoluteURL(profileUrl)
		infoCollector.Visit(profileUrl)
	})

	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(nextPage))
		c.Visit(nextPage)
	})

	infoCollector.OnHTML("ul.ipc-metadata-list", func(e *colly.HTMLElement) {
		tmpProfile := star{}
		tmpProfile.Name = e.ChildText("h1.ipc-title__text")
		tmpProfile.Photo = e.ChildAttr("#name-poster", "src")
		tmpProfile.JobTitle = e.ChildText("#name-job-categories > a > span.itemprop")
		tmpProfile.BirthDate = e.ChildAttr("#name-born-info time", "datetime")
		tmpProfile.Bio = strings.TrimSpace(e.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))
		e.ForEach(" div.knownfor-title", func(_ int, el *colly.HTMLElement) {
			tmpMovie := Movie{}
			tmpMovie.Title = el.ChildText("div.knowfor-title-role > a.knownfor-ellipsis")
			tmpMovie.Year = el.ChildText("div.knownfor-year > span.knownfor-ellipsis")
			tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
		})
		js, err := json.MarshalIndent(tmpProfile, "", "   ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		infoCollector.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting profile URL", r.URL.String())
		})
	})

}

func main() {
	month := flag.Int("month", 1, "Month to crawl")
	day := flag.Int("day", 1, "Day to crawl")
	flag.Parse()
	fmt.Println("Starting crawl for", *month, *day)
	crawl(*month, *day)
}
