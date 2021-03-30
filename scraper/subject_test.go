package scraper

import (
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/go-playground/assert/v2"
	"testing"
)

type TestInput struct {
	subject        string
	course         string
	specialization string
}

var testCases = []struct {
	input    TestInput
	expected entity.Subject
}{
	{
		input: TestInput{
			subject:        "SCC0230",
			course:         "55090",
			specialization: "0",
		},
		expected: entity.Subject{
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
	},
	{
		input: TestInput{
			subject:        "SCC0222",
			course:         "55041",
			specialization: "0",
		},
		expected: entity.Subject{
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
	},
}

func TestNewSubjectScraper(t *testing.T) {
	for _, tc := range testCases {
		sc := NewSubjectScraper(tc.input.subject, tc.input.course, tc.input.specialization)

		if sub, err := sc.Start(); err != nil {
			t.Fatal(err)
		} else {
			assert.Equal(t, sub.(entity.Subject), tc.expected)
		}

	}
}
