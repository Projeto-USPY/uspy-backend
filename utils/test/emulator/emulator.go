package emulator

import (
	"context"
	"time"

	"sync"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/Projeto-USPY/uspy-backend/iddigital"
	"github.com/Projeto-USPY/uspy-backend/server/models/account"
)

var (
	insertOnce sync.Once
	insertCnt  int = 0
)

var testSubjects = []entity.Subject{
	{
		Code:           "SCC0230",
		CourseCode:     "55090",
		Specialization: "0",
		Name:           "Inteligência Artificial",
		Description:    "Apresentar ao aluno as idéias fundamentais da Inteligência Artificial e algumas características relacionadas à implementação desse tipo de sistemas.",
		ClassCredits:   4,
		AssignCredits:  1,
		TotalHours:     "90 h",
		Stats:          map[string]int{"total": 0, "worth_it": 0},
	},
	{
		Code:           "SCC0222",
		CourseCode:     "55041",
		Specialization: "0",
		Name:           "Laboratório de Introdução à Ciência de Computação I",
		Description:    "Implementar em laboratório as técnicas de programação apresentadas em Introdução à Ciência da Computação I, utilizando uma linguagem de programação estruturada.",
		ClassCredits:   2,
		AssignCredits:  2,
		TotalHours:     "90 h",
		Stats:          map[string]int{"total": 0, "worth_it": 0},
	},
	{
		Code:           "SCC0217",
		CourseCode:     "55041",
		Specialization: "0",
		Name:           "Linguagens de Programação e Compiladores",
		Description:    "Dar ao aluno as noções básicas sobre linguagens de programação e técnicas de construção de compiladores para linguagens de programação de alto nível.",
		ClassCredits:   4,
		AssignCredits:  2,
		TotalHours:     "120 h",
		Stats:          map[string]int{"total": 0, "worth_it": 0},
	},
}

var testCourses = []entity.Course{
	{
		Name:           "Bacharelado em Ciência de Dados",
		Code:           "55090",
		Specialization: "0",
		SubjectCodes: map[string]string{
			"SCC0230": "Inteligência Artificial",
		},
	},
	{
		Name:           "Bacharelado em Ciências de Computação",
		Code:           "55041",
		Specialization: "0",
		SubjectCodes: map[string]string{
			"SCC0222": "Laboratório de Introdução à Ciência de Computação I",
			"SCC0217": "Linguages de Programação e Compiladores",
		},
	},
}

func Setup(DB db.Env) error {
	config.TestSetup()

	errChannel := make(chan error, 100)
	var wg sync.WaitGroup
	for _, v := range testSubjects {
		wg.Add(1)
		go func(v entity.Subject) {
			defer wg.Done()
			errChannel <- DB.Insert(v, "subjects")
		}(v)
	}

	for _, c := range testCourses {
		wg.Add(1)
		go func(c entity.Course) {
			defer wg.Done()
			errChannel <- DB.Insert(c, "courses")
		}(c)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		user, userErr := entity.NewUserWithOptions(
			"123456789",
			"r4nd0mpass123!@#",
			"Usuário teste",
			time.Now(),
			entity.WithPasswordHash{},
			entity.WithNameHash{},
		)

		errChannel <- userErr

		recs := iddigital.Records{
			Name: user.Name,
			Nusp: user.Login,
			Grades: []entity.Grade{
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
		}

		errChannel <- account.Signup(DB, user, recs)
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

func MustGet() db.Env {
	if emu, err := Get(); err != nil {
		panic("failed to get emulator while running MustGet")
	} else {
		return emu
	}
}

func Get() (testDB db.Env, getError error) {
	testDB = db.Env{Ctx: context.Background()}

	if client, err := firestore.NewClient(testDB.Ctx, "test"); err != nil {
		return db.Env{}, err
	} else {
		testDB.Client = client
	}

	if err := Setup(testDB); err != nil {
		getError = err
		return
	}

	return
}
