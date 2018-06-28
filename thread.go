package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

// Thread ...
type Thread struct {
	Name       string   `json:"name"`
	Href       string   `json:"href"`
	ReplyCount uint64   `json:"replycount"`
	Views      uint64   `json:"views"`
	Page       uint64   `json:"page"`
	Pages      uint64   `json:"pages"`
	Replies    []*Reply `json:"replies"`
	XFToken    string   `json:"xftoken"`
}

func fetchThreads(s *Section) {
	ThreadCollector := colly.NewCollector()

	ThreadCollector.SetRequestTimeout(time.Second * 60)
	ThreadCollector.OnRequest(onRequest)
	ThreadCollector.OnError(onError)

	ThreadCollector.RedirectHandler = onRedirect

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

		if ps, err = strconv.ParseUint(e.ChildText("div.main > div.titleText > div.secondRow > div.posterDate > span.itemPageNav > a[href]:last-child"), 10, 64); err != nil {
			ps = 1
		}

		if isInCache(
			e.ChildText("div.main > div.titleText > h3.title > a[href].PreviewTooltip"),
			e.ChildAttr("div.main > div.titleText > h3.title > a[href].PreviewTooltip", "href"),
		) {
			return
		}

		s.Threads = append(s.Threads, &Thread{
			Name:       e.ChildText("div.main > div.titleText > h3.title > a[href].PreviewTooltip"),
			Href:       e.ChildAttr("div.main > div.titleText > h3.title > a[href].PreviewTooltip", "href"),
			ReplyCount: r,
			Views:      v,
			Page:       1,
			Pages:      ps,
		})
	})

	ThreadCollector.OnScraped(func(r *colly.Response) {
		Progress.Prefix(fmt.Sprintf("ThreadCollector: %s | %d threads", formatTarget(s, nil), len(s.Threads)))

		if len(s.Threads) == 0 {
			if s.Page >= s.Pages {
				SectionIndex++
				fetchThreads(Sections[SectionIndex])
			} else {
				s.Page++
				fetchThreads(s)
			}
		} else {
			fetchReply(s, s.Threads[ThreadTracker[s]])
		}
	})

	ThreadCollector.Visit(formatTarget(s, nil))
}
