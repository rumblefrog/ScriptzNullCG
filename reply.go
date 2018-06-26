package main

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// Reply ...
type Reply struct {
	Author   string `json:"author"`
	Content  string `json:"content"`
	Liked    bool   `json:"liked"`
	LikeHref string `json:"likehref"`
	Date     string `json:"date"`
}

func fetchReply(s *Section, t *Thread) {
	if isInCache(t) {
		t.Page = t.Pages
		process(s, t)
	}

	ReplyCollector.SetRequestTimeout(60000000000)
	ReplyCollector.OnRequest(onRequest)
	ReplyCollector.OnError(onError)

	var l bool

	ReplyCollector.OnHTML("li.message[data-author]", func(e *colly.HTMLElement) {

		if len(e.ChildText("span.LikeLabel")) > 0 {
			l = strings.Contains(e.ChildText("span.LikeLabel"), "Unlike")
		} else {
			l = true
		}

		t.Replies = append(t.Replies, &Reply{
			Author:   e.Attr("data-author"),
			Content:  e.ChildText("article > blockquote.messageText"),
			Liked:    l,
			LikeHref: e.ChildAttr("a[href].LikeLinkHide", "href"),
			Date:     e.ChildText("a[href].datePermalink > span.DateTime"),
		})
	})

	ReplyCollector.OnHTML("input[name=_xfToken]:first-of-type", func(e *colly.HTMLElement) {
		t.XFToken = e.Attr("value")
	})

	ReplyCollector.OnScraped(func(r *colly.Response) {
		process(s, t)
	})

	log.Println("ReplyCollector: " + formatTarget(s, t))

	ReplyCollector.Visit(formatTarget(s, t))
}
