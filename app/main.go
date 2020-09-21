package main

import (
	"fmt"

	"github.com/tpreischadt/ProjetoJupiter/depscraper"
)

func main() {
	results := *depscraper.Scrape()

	for k, v := range results {
		fmt.Println(k, v)
	}
}
