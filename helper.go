package main

import (
	"fmt"
)

// Target - Main scrape site
const Target = "https://scriptznull.nl/"

func formatTarget(s *Section, t *Thread) string {
	if s == nil {
		return Target
	}

	if s != nil && t == nil {
		return fmt.Sprintf("%s%spage-%d", Target, s.Href, s.Page)
	}

	return fmt.Sprintf("%s%s%spage-%d", Target, s.Href, t.Href, t.Page)
}
