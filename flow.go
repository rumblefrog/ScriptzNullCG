package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// SectionIndex - Per Page of Threads Indexing
// ThreadIndex - Per Section Page of Threads Indexing
var (
	SectionIndex = 1 // Skip Announcements and start with General Discussion
	ThreadIndex  int
)

// Tracker - Keep track of index processing
var Tracker map[*Thread]int

func process(s *Section, t *Thread) {
	// Process all replies, check if we have enough, if not, fetch more
	index, ok := Tracker[t]

	if ok == false {
		Tracker[t] = 0
	}

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

	Tracker[t] = index

	if t.Page >= t.Pages {
		addToCache(t)
		if ThreadIndex+1 < len(s.Threads) {
			ThreadIndex++
			fetchReply(s, s.Threads[ThreadIndex])
		} else {
			if s.Page >= s.Pages {
				ThreadIndex = 0
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
