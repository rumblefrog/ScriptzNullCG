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
		if s.Page == 1 {
			return fmt.Sprintf("%s%s", Target, s.Href)
		}
		return fmt.Sprintf("%s%spage-%d", Target, s.Href, s.Page)
	}

	if t.Page == 1 {
		return fmt.Sprintf("%s%s", Target, t.Href)
	}

	return fmt.Sprintf("%s%spage-%d", Target, t.Href, t.Page)
}
