package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gocolly/colly"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
)

// UA - Useragent to use for request
// Cookie - Authentication payload
// Credit - Amount of credit to process
// TotalProcessed - Total amount processed
// Progress - Pointer to progress bar
// Header - HTTP request header
// Sections - Forum categories
// SectionCollector - Forum categories collector
// ThreadCollector - Threads collector
// ReplyCollector - Reply collector
var (
	UA               = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36"
	Cookie           string
	Credit           int
	TotalProcessed   uint64
	Progress         *pb.ProgressBar
	Headers          http.Header
	Sections         []*Section
	SectionCollector = colly.NewCollector()
	ThreadCollector  = colly.NewCollector()
	ReplyCollector   = colly.NewCollector()
)

func main() {
	app := cli.NewApp()

	app.Name = "ScriptzNullCG"
	app.Usage = "ScriptzNull Credit Generator"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "c",
			Value:       100,
			Usage:       "Desired amount of credits",
			Destination: &Credit,
		},
	}

	app.Action = func(c *cli.Context) error {
		if Cookie = c.Args().Get(0); Cookie == "" {
			log.Fatal("Cookie is not provided")
		}

		Progress = pb.StartNew(Credit)

		Tracker = make(map[*Thread]int)

		Headers = http.Header{
			"User-Agent": {UA},
			"Accept":     {"application/json, text/javascript, */*; q=0.01"},
			"Cookie":     {Cookie},
		}

		loadCache()
		fetchSections()

		Cache = append(Cache, &Thread{
			Name: "test",
		})

		saveCache()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func onRequest(r *colly.Request) {
	r.Headers.Set("User-Agent", UA)
	r.Headers.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	r.Headers.Set("Cookie", Cookie)
}

func onError(r *colly.Response, e error) {
	log.Fatal("Try a fresh token perhaps?: ", e)
}

func onResponse(r *colly.Response) {
	log.Println(r.Body)
}
