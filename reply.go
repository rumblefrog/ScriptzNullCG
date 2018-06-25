package main

import "github.com/gocolly/colly"

// Reply ...
type Reply struct {
	Author string
	Liked  bool
	Date   string
}

func fetchReply(s *Section, t *Thread, p int) {
	ReplyCollector.OnRequest(onRequest)
	ReplyCollector.OnError(onError)

	ReplyCollector.OnHTML("li.message[data-author]", func(e *colly.HTMLElement) {
		t.Replies = append(t.Replies, &Reply{
			Author: e.Attr("data-author"),
			Liked:  e.ChildText("span.LikeLabel") == "Unlike",
			Date:   e.ChildText("a[href].datePermalink > span.DateTime"),
		})
	})
}
