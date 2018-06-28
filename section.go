package main

import (
	"math"
	"math/rand"
	"strconv"
	"time"

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
	Search       bool
}

var (
	tc  uint64
	mc  uint64
	err error
)

func fetchSections() {
	SectionCollector := colly.NewCollector()

	SectionCollector.SetRequestTimeout(time.Second * 60)
	SectionCollector.OnRequest(onRequest)
	SectionCollector.OnError(onError)

	SectionCollector.OnHTML("li.node > div.nodeInfo > div.nodeText", func(e *colly.HTMLElement) {
		if tc, err = strconv.ParseUint(e.ChildText("div.nodeStats > dl:first-child > dd"), 10, 64); err != nil {
			return
		}

		if mc, err = strconv.ParseUint(e.ChildText("div.nodeStats > dl:last-child > dd"), 10, 64); err != nil {
			return
		}

		if e.ChildText("h3.nodeTitle > a[href]") == "Announcements" {
			return
		}

		Sections = append(Sections, &Section{
			Name:         e.ChildText("h3.nodeTitle > a[href]"),
			Href:         e.ChildAttr("h3.nodeTitle > a[href]", "href"),
			Page:         1,
			Pages:        uint64(math.Ceil(float64(tc) / float64(20))),
			ThreadCount:  tc,
			MessageCount: mc,
			Search:       false,
		})
	})

	SectionCollector.OnScraped(func(r *colly.Response) {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(Sections), func(i, j int) { Sections[i], Sections[j] = Sections[j], Sections[i] })

		Sections = append(Sections, Sections[0])
		Sections[0] = &Section{
			Name:         "New Posts",
			Href:         "find-new/posts",
			Page:         1,
			Pages:        10,
			ThreadCount:  200,
			MessageCount: 0,
			Search:       true,
		}

		Progress.Prefix("SectionCollector: Done")

		fetchThreads(Sections[0])
	})

	SectionCollector.Visit(formatTarget(nil, nil))
}
