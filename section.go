package main

import (
	"math"
	"strconv"

	"github.com/gocolly/colly"
)

// Section ...
type Section struct {
	Name         string    `json:"name"`
	Href         string    `json:"href"`
	Page         uint64    `json:"page"`
	Pages        uint64    `json:"pages"`
	Threads      []*Thread `json:"threads"`
	ThreadCount  uint64    `json:"threadcount"`
	MessageCount uint64    `json:"messagecount"`
}

var (
	tc  uint64
	mc  uint64
	err error
)

func fetchSections() {
	SectionCollector.SetRequestTimeout(30000000000)
	SectionCollector.OnRequest(onRequest)
	SectionCollector.OnError(onError)

	SectionCollector.OnHTML("li.node > div.nodeInfo > div.nodeText", func(e *colly.HTMLElement) {
		if tc, err = strconv.ParseUint(e.ChildText("div.nodeStats > dl:first-child > dd"), 10, 64); err != nil {
			return
		}

		if mc, err = strconv.ParseUint(e.ChildText("div.nodeStats > dl:last-child > dd"), 10, 64); err != nil {
			return
		}

		Sections = append(Sections, &Section{
			Name:         e.ChildText("h3.nodeTitle > a[href]"),
			Href:         e.ChildAttr("h3.nodeTitle > a[href]", "href"),
			Page:         1,
			Pages:        uint64(math.Ceil(float64(tc) / float64(20))),
			ThreadCount:  tc,
			MessageCount: mc,
		})
	})

	SectionCollector.OnScraped(func(r *colly.Response) {
		fetchThreads(Sections[0])
	})

	SectionCollector.Visit(formatTarget(nil, nil))
}
