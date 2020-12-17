package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"io/ioutil"
	"net/http"
	"strconv"
)

func getProfessorByCode(codPes int) (entity.Professor, int, error) {
	infoURL := "https://uspdigital.usp.br/datausp/servicos/publico/indicadores_pos/perfil/docente/"
	infoURL += strconv.Itoa(codPes)
	resp, err := http.Get(infoURL)

	if err != nil {
		return entity.Professor{}, -1, err
	} else if resp.StatusCode != http.StatusOK {
		return entity.Professor{}, resp.StatusCode, nil
	} else {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return entity.Professor{}, -1, fmt.Errorf("error reading json with codPes %v", codPes)
	}

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	prof := entity.Professor{
		CodPes: codPes,
		Name:   fmt.Sprintf("%s", data["nompes"]),
	}

	return prof, http.StatusOK, nil
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

			results = append(results, entity.Offering{
				Semester:  semester,
				Professor: codPes,
				Year:      year,
				Subject:   subjName,
			})
		}
	}

	return results, nil
}
