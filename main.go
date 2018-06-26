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
	UA               = ""
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
		cli.StringFlag{
			Name:        "auth",
			Value:       "",
			Usage:       "The authentication cookie header",
			Destination: &Cookie,
		},
		cli.StringFlag{
			Name:        "ua",
			Value:       "",
			Usage:       "The user-agent you authenticated with",
			Destination: &UA,
		},
	}

	app.Action = func(c *cli.Context) error {
		if Cookie == "" {
			log.Fatal("Cookie (auth) is not provided")
		}

		if UA == "" {
			log.Fatal("User-agent (ua) was not provided")
		}

		Progress = pb.StartNew(Credit)

		Tracker = make(map[*Thread]int)

		Headers = http.Header{
			"User-Agent":   {UA},
			"Accept":       {"application/json, text/javascript, */*; q=0.01"},
			"Content-Type": {"application/x-www-form-urlencoded"},
			"Cookie":       {Cookie},
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
