package main

import (
	"strconv"

	"github.com/gocolly/colly"
)

// Thread ...
type Thread struct {
	Name       string
	Href       string
	ReplyCount uint64
	Views      uint64
	Page       uint64
	Pages      uint64
	Replies    []*Reply
}

func fetchThreads(s *Section) {
	ThreadCollector.OnRequest(onRequest)
	ThreadCollector.OnError(onError)

	var (
		r   uint64
		v   uint64
		ps  uint64
		err error
	)

	ThreadCollector.OnHTML("li.discussionListItem", func(e *colly.HTMLElement) {
		if r, err = strconv.ParseUint(e.ChildText("div.stats > dl.major > dd"), 10, 64); err != nil {
			return
		}

		if v, err = strconv.ParseUint(e.ChildText("div.stats > dl.minor > dd"), 10, 64); err != nil {
			return
		}

		if ps, err = strconv.ParseUint(e.ChildText("div.main > div.titleText > div.secondRow > span.itemPageNav > a[href]:last-child"), 10, 64); err != nil {
			return
		}

		s.Threads = append(s.Threads, &Thread{
			Name:       e.ChildText("div.main > div.titleText > h3.title > a[href].PreviewToolTip"),
			Href:       e.ChildAttr("div.main > div.titleText > h3.title > a[href].PreviewToolTip", "href"),
			ReplyCount: r,
			Views:      v,
			Page:       1,
			Pages:      ps,
		})
	})

	ThreadCollector.Visit(formatTarget(s, nil))
}
