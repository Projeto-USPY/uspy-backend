package emulator

import (
	"context"
	"time"

	"sync"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
)

var (
	emuDB      db.Env
	insertOnce sync.Once
	insertCnt  int = 0
)

func getInsertables() (objects []db.Object, getError error) {
	if user, err := entity.NewUserWithOptions(
		"123456789",
		"r4nd0mpass123!@#",
		"Usuário teste",
		time.Now(),
		entity.WithNameHash{},
		entity.WithPasswordHash{},
	); err != nil {
		objects, getError = []db.Object{}, err
		return
	} else {
		objects = append(objects, db.Object{Collection: "users", Doc: user.Hash(), Data: user})
	}

	subjects := []entity.Subject{
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

	courses := []entity.Course{
		{
			Name:           "Bacharelado em Ciência de Dados",
			Code:           "55090",
			Specialization: "0",
			Subjects:       []entity.Subject{subjects[0]},
		},
		{
			Name:           "Bacharelado em Ciências de Computação",
			Code:           "55041",
			Specialization: "0",
			Subjects:       []entity.Subject{subjects[1], subjects[2]},
		},
	}

	for _, s := range subjects {
		objects = append(objects, db.Object{Collection: "subjects", Doc: s.Hash(), Data: s})
	}

	for _, c := range courses {
		objects = append(objects, db.Object{Collection: "courses", Doc: c.Hash(), Data: c})
	}

	return objects, nil
}

func setup(DB db.Env) error {
	config.TestSetup()
	if objects, err := getInsertables(); err != nil {
		return err
	} else {
		if err := DB.BatchWrite(objects); err != nil {
			return err
		}
	}

	return nil
}

func MustGet() db.Env {
	if emu, err := Get(); err != nil {
		panic("failed to get emulator while running MustGet")
	} else {
		return emu
	}
}

func Get() (testDB db.Env, getError error) {
	if emuDB != (db.Env{}) {
		return emuDB, nil
	}

	testDB = db.Env{Ctx: context.Background()}

	if client, err := firestore.NewClient(testDB.Ctx, "test"); err != nil {
		return db.Env{}, err
	} else {
		testDB.Client = client
	}

	insertOnce.Do(
		func() {
			if err := setup(testDB); err != nil {
				getError = err
				return
			}
		},
	)

	emuDB = testDB
	return
}
