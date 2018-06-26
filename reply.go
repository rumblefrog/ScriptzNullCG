package main

import (
	"log"

	"github.com/gocolly/colly"
)

// Reply ...
type Reply struct {
	Author string `json:"author"`
	Liked  bool   `json:"liked"`
	Date   string `json:"date"`
}

func fetchReply(s *Section, t *Thread) {
	if isInCache(t) {
		t.Page = t.Pages
		process(s, t)
	}

	ReplyCollector.SetRequestTimeout(30000000000)
	ReplyCollector.OnRequest(onRequest)
	ReplyCollector.OnError(onError)

	ReplyCollector.OnHTML("li.message[data-author]", func(e *colly.HTMLElement) {
		t.Replies = append(t.Replies, &Reply{
			Author: e.Attr("data-author"),
			Liked:  e.ChildText("span.LikeLabel") == "Unlike",
			Date:   e.ChildText("a[href].datePermalink > span.DateTime"),
		})
	})

	ReplyCollector.OnHTML("input[type=_xfToken]", func(e *colly.HTMLElement) {
		t.XFToken = e.Attr("value")
	})

	ReplyCollector.OnScraped(func(r *colly.Response) {
		process(s, t)
	})

	log.Println("ReplyCollector: " + formatTarget(s, t))

	ReplyCollector.Visit(formatTarget(s, t))
}
