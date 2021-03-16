package scraper

import (
	"errors"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"io"
	"net/http"
)

var (
	DefaultInstituteURLMask = "https://uspdigital.usp.br/jupiterweb/jupCursoLista?codcg=%s&tipo=N"
	DefaultCourseURLMask    = "https://uspdigital.usp.br/jupiterweb/listarGradeCurricular?codcg=%s&codcur=%s&codhab=%s&tipo=N"
	DefaultSubjectURLMask   = "https://uspdigital.usp.br/jupiterweb/obterDisciplina?sgldis=%s&codcur=%s&codhab=%s"
)

var (
	ErrorCourseNotExist   = errors.New("could not fetch course in institute page")
	ErrorCourseNoSubjects = errors.New("could not fetch subjects in course page")
)

type Starter interface {
	Start() (db.Manager, error)
}

type Scraper interface {
	Starter
	Scrape(reader io.Reader) (db.Manager, error)
}

func Start(scraper Scraper, startURL string) (db.Manager, error) {
	client := &http.Client{
		Timeout: 0,
	}

	var object db.Manager

	if resp, reader, err := utils.HTTPGetWithUTF8(client, startURL); err != nil {
		return nil, err
	} else {
		// successfully got start page
		defer resp.Body.Close()
		if object, err = scraper.Scrape(reader); err != nil {
			return nil, err
		}
	}

	return object, nil
}
