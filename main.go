package main

import (
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/net/html"
)

const (
	target = "https://scriptznull.nl/"
)

var (
	cookie   string
	credit   uint64
	client   = &http.Client{}
	sections []string
	header   http.Header
)

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

		header = http.Header{
			"Accept":          {"application/json, text/javascript, */*; q=0.01"},
			"Accept-Language": {"Accept-Language: en-US,en;q=0.9"},
			"Accept-Encoding": {"gzip, deflate"},
			"Cookie":          {cookie},
			"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36"},
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
	req, err := http.NewRequest("GET", target, nil)

	if err != nil {
		log.Panicln(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Panicln(err)
	}

	root, err := html.Parse(resp.Body)

	if err != nil {
		log.Panicln(err)
	}

	if resp.StatusCode != 200 {
		log.Panicln("Try entering a fresh cookie")
	}

	log.Println("UH?")

}
