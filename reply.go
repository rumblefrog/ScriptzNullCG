package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

// Reply ...
type Reply struct {
	ID       string `json:"id"`
	Author   string `json:"author"`
	Content  string `json:"content"`
	LikeHref string `json:"likehref"`
	Date     string `json:"date"`
}

func fetchReply(s *Section, t *Thread) {
	if isInCache(t) {
		log.Println("Hit Cache")
		t.Page = t.Pages
		process(s, t)
	}

	ReplyCollector.SetRequestTimeout(60000000000)
	ReplyCollector.OnRequest(onRequest)
	ReplyCollector.OnError(onError)

	ReplyCollector.OnHTML("div[id=headerMover]", func(e *colly.HTMLElement) {
		e.ForEach("li.message[data-author]", func(_ int, m *colly.HTMLElement) {
			if HasID(t, m.Attr("id")) {
				return
			}

			if len(m.ChildText("a[href].like > span.LikeLabel")) <= 0 {
				return
			}

			t.Replies = append(t.Replies, &Reply{
				ID:       m.Attr("id"),
				Author:   m.Attr("data-author"),
				Content:  m.ChildText("article > blockquote.messageText"),
				LikeHref: m.ChildAttr("a[href].LikeLinkHide", "href"),
				Date:     m.ChildText("a[href].datePermalink > span.DateTime"),
			})

			log.Printf(fmt.Sprintf("Inserted %s into %s", m.Attr("id"), t.Name))
		})

		t.XFToken = e.ChildAttr("input[name=_xfToken]:first-of-type", "value")
	})

	ReplyCollector.OnScraped(func(r *colly.Response) {
		log.Print("Memory: ")
		log.Println(t.Replies)
		for _, v := range t.Replies {
			log.Println("Loop: " + v.ID)
		}
		process(s, t)
	})

	log.Println("ReplyCollector: " + formatTarget(s, t))

	ReplyCollector.Visit(formatTarget(s, t))
}

// HasID - Checks thread replies array for matching ID
func HasID(s *Thread, ID string) bool {
	for _, v := range s.Replies {
		if v.ID == ID {
			return true
		}
	}
	return false
}
