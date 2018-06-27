package main

import (
	"fmt"
	"time"

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
	ReplyCollector := colly.NewCollector()

	ReplyCollector.SetRequestTimeout(time.Second * 60)
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
		})

		t.XFToken = e.ChildAttr("input[name=_xfToken]:first-of-type", "value")
	})

	ReplyCollector.OnScraped(func(r *colly.Response) {
		Progress.Prefix(fmt.Sprintf("ReplyCollector: %s", formatTarget(s, t)))
		process(s, t)
	})

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
