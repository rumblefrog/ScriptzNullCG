package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// SectionIndex - Per Page of Threads Indexing
var SectionIndex = 1 // Skip Announcements and start with General Discussion

// ThreadTracker - Keep track of thread processing
var ThreadTracker map[*Section]int

// ReplyTracker - Keep track of index processing
var ReplyTracker map[*Thread]int

func process(s *Section, t *Thread) {
	// Process all replies, check if we have enough, if not, fetch more
	index := ReplyTracker[t]

	for ; index < len(t.Replies); index++ {
		if len(t.Replies[index].Content) > 20 && CreatePayload(t, t.Replies[index]) {
			if Credit <= 0 {
				saveCache()

				byteValues, _ := json.Marshal(&Sections)

				ioutil.WriteFile("forum.json", byteValues, 0777)

				Progress.FinishPrint(fmt.Sprintf("Processed %d credits", TotalProcessed))

				os.Exit(0)
			}
			Credit -= 3
			TotalProcessed += 3
			Progress.Add(3)
		}
	}

	ReplyTracker[t] = index

	ThreadIndex := ThreadTracker[s]

	if t.Page >= t.Pages {
		go addToCache(t)
		if ThreadIndex+1 < len(s.Threads) {
			ThreadIndex++
			ThreadTracker[s] = ThreadIndex
			fetchReply(s, s.Threads[ThreadIndex])
		} else {
			if s.Page >= s.Pages {
				SectionIndex++
				fetchThreads(Sections[SectionIndex])
			} else {
				s.Page++
				fetchThreads(s)
			}
		}
	} else {
		t.Page++
		fetchReply(s, t)
	}
}
