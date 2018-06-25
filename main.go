package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/urfave/cli"
)

const (
	target = "https://scriptznull.nl/"
)

var (
	ua       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36"
	cookie   string
	credit   uint64
	client   = &http.Client{}
	sections []*Section
	header   http.Header
)

// Section ...
type Section struct {
	Name         string
	Href         string
	Threads      []*Thread
	ThreadCount  uint64
	MessageCount uint64
}

// Thread ...
type Thread struct {
	Name    string
	Href    string
	Replies uint64
	Views   uint64
}

// History ..
type History struct {
	Threads []string
	PostIDs []uint32
}

func main() {
	app := cli.NewApp()

	app.Name = "ScriptzNullCG"
	app.Usage = "ScriptzNull Credit Generator"

	app.Flags = []cli.Flag{
		cli.Uint64Flag{
			Name:        "c",
			Value:       100,
			Usage:       "Desired amount of credits",
			Destination: &credit,
		},
	}

	app.Action = func(c *cli.Context) error {
		if cookie = c.Args().Get(0); cookie == "" {
			log.Fatal("Cookie is not provided")
		}

		fetchSections()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchSections() {
	c := colly.NewCollector()

	c.OnRequest(onRequest)
	c.OnError(onError)

	var (
		tc  uint64
		mc  uint64
		err error
	)

	c.OnHTML("li.node > div.nodeInfo > div.nodeText", func(e *colly.HTMLElement) {
		if tc, err = strconv.ParseUint(e.ChildText("div.nodeStats > dl:first-child > dd"), 10, 64); err != nil {
			return
		}

		if mc, err = strconv.ParseUint(e.ChildText("div.nodeStats > dl:last-child > dd"), 10, 64); err != nil {
			return
		}

		sections = append(sections, &Section{
			Name:         e.ChildText("h3.nodeTitle > a[href]"),
			Href:         e.ChildAttr("h3.nodeTitle > a[href]", "href"),
			ThreadCount:  tc,
			MessageCount: mc,
		})
	})

	c.OnScraped(func(r *colly.Response) {
		for _, s := range sections {
			log.Printf("Name: %s | Directive: %s | Threads: %d | Replies: %d", s.Name, s.Href, s.ThreadCount, s.MessageCount)
		}
	})

	c.Visit(target)
}

func fetchThreads(s *Section) {
	c := colly.NewCollector()

	c.OnRequest(onRequest)
	c.OnError(onError)

	var (
		r   uint64
		v   uint64
		err error
	)

	c.OnHTML("li.discussionListItem", func(e *colly.HTMLElement) {
		if r, err = strconv.ParseUint(e.ChildText("div.stats > dl.major > dd"), 10, 64); err != nil {
			return
		}

		if v, err = strconv.ParseUint(e.ChildText("div.stats > dl.minor > dd"), 10, 64); err != nil {
			return
		}

		s.Threads = append(s.Threads, &Thread{
			Name:    e.ChildText("div.main > div.titleText > h3.title > a[href].PreviewToolTip"),
			Href:    e.ChildAttr("div.main > div.titleText > h3.title > a[href].PreviewToolTip", "href"),
			Replies: r,
			Views:   v,
		})
	})

	c.OnHTML("div.PageNav nav > a[href].text", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
}

func onRequest(r *colly.Request) {
	r.Headers.Set("User-Agent", ua)
	r.Headers.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	r.Headers.Set("Cookie", cookie)
}

func onError(r *colly.Response, e error) {
	log.Println("Try a fresh token perhaps?: ", e)
}
