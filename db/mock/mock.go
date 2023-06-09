package mock

import (
	"fmt"
	"sync"
	"time"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/iddigital"
	"github.com/Projeto-USPY/uspy-backend/server/models/account"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// Mock Dat
var (
	MockSubjects = []models.Subject{
		{
			Code:           "SCC0230",
			CourseCode:     "55090",
			Specialization: "0",
			Name:           "Inteligência Artificial",
			Semester:       6,
			Description:    "Apresentar ao aluno as idéias fundamentais da Inteligência Artificial e algumas características relacionadas à implementação desse tipo de sistemas.",
			ClassCredits:   4,
			AssignCredits:  1,
			TotalHours:     "90 h",
			Stats:          map[string]int{"total": 0, "worth_it": 0},
			Optional:       false,
		},
		{
			Code:           "SCC0222",
			CourseCode:     "55041",
			Specialization: "0",
			Name:           "Laboratório de Introdução à Ciência de Computação I",
			Semester:       2,
			Description:    "Implementar em laboratório as técnicas de programação apresentadas em Introdução à Ciência da Computação I, utilizando uma linguagem de programação estruturada.",
			ClassCredits:   2,
			AssignCredits:  2,
			TotalHours:     "90 h",
			Stats:          map[string]int{"total": 0, "worth_it": 0},
			Optional:       true,
		},
		{
			Code:           "SCC0217",
			CourseCode:     "55041",
			Specialization: "0",
			Name:           "Linguagens de Programação e Compiladores",
			Semester:       6,
			Description:    "Dar ao aluno as noções básicas sobre linguagens de programação e técnicas de construção de compiladores para linguagens de programação de alto nível.",
			ClassCredits:   4,
			AssignCredits:  2,
			TotalHours:     "120 h",
			Stats:          map[string]int{"total": 0, "worth_it": 0},
			Optional:       false,
		},
	}

	MockInstitutes = []models.Institute{
		{
			Name: "Instituto de Ciências Matemáticas e de Computação",
			Code: "55",
		},
	}

	MockCourses = []models.Course{
		{
			Name:           "Bacharelado em Ciência de Dados",
			Code:           "55090",
			Specialization: "0",
			Shift:          "integral",
			SubjectCodes: map[string]string{
				"SCC0230": "Inteligência Artificial",
			},
		},
		{
			Name:           "Bacharelado em Ciências de Computação",
			Code:           "55041",
			Specialization: "0",
			Shift:          "integral",
			SubjectCodes: map[string]string{
				"SCC0222": "Laboratório de Introdução à Ciência de Computação I",
				"SCC0217": "Linguages de Programação e Compiladores",
			},
		},
	}

	MockTranscript = iddigital.Transcript{
		Name: "Usuário teste",
		Nusp: "123456789",
		Grades: []models.Record{
			{
				Grade:          9.0,
				Frequency:      100,
				Status:         "A",
				Subject:        "SCC0217",
				Course:         "55041",
				Specialization: "0",
				Semester:       1,
				Year:           2018,
			},
			{
				Grade:          9.0,
				Frequency:      60,
				Status:         "RF",
				Subject:        "SCC0217",
				Course:         "55041",
				Specialization: "0",
				Semester:       1,
				Year:           2017,
			},
			{
				Grade:          4.0,
				Frequency:      90,
				Status:         "RN",
				Subject:        "SCC0217",
				Course:         "55041",
				Specialization: "0",
				Semester:       1,
				Year:           2016,
			},
			{
				Grade:          8.0,
				Frequency:      95,
				Status:         "A",
				Subject:        "SCC0222",
				Course:         "55041",
				Specialization: "0",
				Semester:       2,
				Year:           2018,
			},
			{
				Grade:          4.0,
				Frequency:      93,
				Status:         "A",
				Subject:        "SCC0222",
				Course:         "55041",
				Specialization: "0",
				Semester:       2,
				Year:           2017,
			},
		},
		Course:         "55041",
		Specialization: "0",
	}
)

// SetupMockData inputs some mocked data into the emulator database
func SetupMockData(
	DB db.Database,
) error {
	config.TestSetup()

	timezone, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return err
	}

	errChannel := make(chan error, 100)
	var wg sync.WaitGroup
	for _, v := range MockSubjects {
		wg.Add(1)
		go func(v models.Subject) {
			defer wg.Done()
			errChannel <- DB.Insert(v, "subjects")
		}(v)
	}

	for _, i := range MockInstitutes {
		wg.Add(1)
		go func(i models.Institute) {
			defer wg.Done()
			errChannel <- DB.Insert(i, "institutes")
		}(i)
	}

	for _, c := range MockCourses {
		wg.Add(1)
		go func(c models.Course) {
			defer wg.Done()
			errChannel <- DB.Insert(c, fmt.Sprintf("institutes/%s/courses", utils.SHA256("55")))
		}(c)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		user, userErr := models.NewUser(
			"123456789",
			"Usuário teste",
			time.Date(2020, time.January, 0, 0, 0, 0, 0, timezone),
			map[string][]int{
				"2018": {1, 2},
			},
		)

		errChannel <- userErr
		errChannel <- account.InsertUser(DB, user, &MockTranscript)
		errChannel <- models.CompleteSignup(DB, user.Hash(), "users", "email_test@usp.br", "r4nd0mpass123!@#")
	}()

	wg.Wait()
	close(errChannel)

	var jointErr error
	for err := range errChannel {
		if err != nil && jointErr == nil {
			jointErr = err
		}
	}

	return jointErr
}
