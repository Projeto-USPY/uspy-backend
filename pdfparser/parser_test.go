package pdfparser

import (
	"fmt"
	"testing"
)

func TestReadPDF(t *testing.T) {
	body, ok := ReadPDFFile("historico.pdf")

	if ok != true {
		t.Error("An error ocurred")
	} else {
		fmt.Println(*body)
	}
}

func TestParsePDF(t *testing.T) {
	body, ok := ReadPDFFile("historico.pdf")

	if ok != true {
		t.Error("Error reading PDF")
	}

	results, ok := ParsePDF(body)

	if ok != true {
		t.Error("Error parsing PDF")
	} else {
		fmt.Println(results)
	}
}
