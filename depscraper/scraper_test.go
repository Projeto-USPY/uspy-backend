package depscraper

import (
    "fmt"
    "testing"
)

func TestScrape(t *testing.T) {
    result := Scrape()
    fmt.Println(result)
}
