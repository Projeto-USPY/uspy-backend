// package professor contains useful functions require to scrape professor data from icmcpessoas
package professor

import (
	"fmt"
	"testing"
)

func TestScrapeDepartments(t *testing.T) {
	result := ScrapeDepartments()
	fmt.Println(result)
}
