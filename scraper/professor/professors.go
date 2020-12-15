package professor

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func getProfessorsByDepartment(dep *string, page int, offset int) []entity.Professor {
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
		prof := entity.Professor{ID: i + offset, CodPes: int(codPes), Name: profName, Department: *dep}
		results = append(results, prof)
	})

	return results
}

func getProfessorByCode(codPes int) (entity.Professor, error) {
	infoURL := "https://uspdigital.usp.br/datausp/servicos/publico/indicadores_pos/perfil/docente/"
	infoURL += strconv.Itoa(codPes)
	resp, err := http.Get(infoURL)

	if err != nil {
		return entity.Professor{}, err
	} else if resp.StatusCode != http.StatusOK {
		return entity.Professor{}, fmt.Errorf("could not get professor %v", codPes)
	} else {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return entity.Professor{}, fmt.Errorf("error reading json with codPes %v", codPes)
	}

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	prof := entity.Professor{
		CodPes: codPes,
		Name:   fmt.Sprintf("%s", data["nompes"]),
	}

	return prof, nil
}

// GetProfessorHistory gets you the offerings since a given year for a given professor
func GetProfessorHistory(codPes, since int) ([]entity.Offering, error) {
	offerMask := "https://uspdigital.usp.br/datausp/servicos/publico/academico/aulas_ministradas/%d/%d/0/0/br"
	offerURL := fmt.Sprintf(offerMask, codPes, since)
	resp, err := http.Get(offerURL)

	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get professor %v", codPes)
	} else {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading json with codPes %v", codPes)
	}

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	history := data["aulasGradPorAno"].(map[string]interface{})
	results := make([]entity.Offering, 0, 100)

	for k, v := range history {
		offs := v.([]interface{})
		for _, subj := range offs {
			year, _ := strconv.Atoi(k)
			semester := int(fmt.Sprintf("%s", subj.(map[string]interface{})["codtur"])[4] - '0')
			subjName := fmt.Sprintf("%s", subj.(map[string]interface{})["coddis"])
			id := fmt.Sprintf("%v%v%v%v", subjName, codPes, year, semester)

			results = append(results, entity.Offering{
				HashID:    fmt.Sprintf("%x", md5.Sum([]byte(id))),
				Semester:  semester,
				Professor: codPes,
				Year:      year,
				Subject:   subjName,
			})
		}
	}

	return results, nil
}

// ScrapeDepartments scrapes the professors page
func ScrapeDepartments() []entity.Professor {
	deps := []string{"SCC", "SMA", "SME", "SSC"}
	results := make([]entity.Professor, 0, 1000)
	offset := 0

	for _, dep := range deps {
		i := 1
		for {
			profs := getProfessorsByDepartment(&dep, i, offset)
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
