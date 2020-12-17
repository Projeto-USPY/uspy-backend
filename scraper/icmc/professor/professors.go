package professor

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/utils"
	"net/url"
	"strconv"
	"strings"
)

func getProfessorsByDepartment(dep *string, page int) []entity.Professor {
	icmcURL := "https://www.icmc.usp.br/templates/icmc2015/php/pessoas.php"
	formData := url.Values{
		"grupo":  {"Docente"},
		"depto":  {*dep},
		"nome":   {""},
		"pagina": {strconv.Itoa(page)},
	}

	response, body := utils.HTTPPostWithUTF8(icmcURL, formData)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	utils.CheckPanic(err)

	results := make([]entity.Professor, 0, 1000)

	doc.Find(".caption").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		codPesStr := strings.Split(href, "=")[1]
		codPes, _ := strconv.ParseInt(codPesStr, 10, 32)
		codPes = (codPes - 3) / 2
		profName := strings.TrimSpace(s.Text())
		prof := entity.Professor{CodPes: int(codPes), Name: profName, Department: *dep}
		results = append(results, prof)
	})

	return results
}

// ScrapeDepartments scrapes the professors page
func ScrapeDepartments() []entity.Professor {
	deps := []string{"SCC", "SMA", "SME", "SSC"}
	results := make([]entity.Professor, 0, 1000)
	offset := 0

	for _, dep := range deps {
		i := 1
		for {
			profs := getProfessorsByDepartment(&dep, i)
			offset += len(profs)

			if len(profs) == 0 {
				break
			} else {
				results = append(results, profs...)
				i++
			}
		}
	}

	return results
}
