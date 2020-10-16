package pdfparser

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// ReadPDFFile takes the filename of a  PDF and returns its string
func ReadPDFFile(file string) (body *string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recovered] Couldnt read PDF: %v\n", r)
			body = nil
			ok = false
		}
	}()

	ch := make(chan *string, 1)
	go func() {
		out, err := exec.Command("pdftotext", file, "-q", "-eol", "unix", "-layout", "-").Output()
		if err != nil {
			panic("An error occured while reading the PDF")
		}

		str := string(out)
		ch <- &str
	}()

	return <-ch, true
}

// ReadPDFResponse takes the Grades PDF response object and reads it into a string
func ReadPDFResponse(r *http.Response) (body *string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recovered] Couldnt read PDF: %v\n", r)
			body = nil
			ok = false
		}
	}()

	bodyPDF, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic("error converting response body to string")
	}

	ch := make(chan *string, 1)

	go func() {
		parser := exec.Command("pdftotext", "-q", "-eol", "unix", "-layout", "-", "-")

		stdin, _ := parser.StdinPipe()
		stdin.Write(bodyPDF)
		stdin.Close()

		out, err := parser.Output()

		if err != nil {
			log.Print(err)
			panic("an error occured while executing pdftotext")
		}

		str := string(out)
		ch <- &str
	}()

	return <-ch, true
}

// Grade represents a Grade in jupiterweb
type Grade struct {
	subject string
	grade   float64
	status  string
}

// Student represents an ICMC student
type Student struct {
	Grades []Grade
	Nusp   string
}

// ParsePDF takes the (already read) PDF string and parses it to a list of Grades
func ParsePDF(body *string) (st Student, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recovered] Couldnt parse PDF string: %v\n", r)
			st.Grades = nil
			st.Nusp = ""
			ok = false
		}
	}()

	i := 2
	strPDF := *body
	st.Nusp = ""

	for {
		// End of PDF
		if i == len(strPDF) {
			break
		}

		if idx, ok := nuspInRow(i, body); ok && st.Nusp == "" {
			nusp, err := parseNUSP((*body)[i:idx])

			if err != nil {
				panic("couldnt parse nusp")
			}

			st.Nusp = nusp[:len(nusp)-1]
		}

		// Found a subject
		if isSubject(i, body) {
			var j int = i

			// Get to end of subject code
			for strPDF[j] != ' ' {
				j++
			}

			// Copying subject code to new slice
			subjectCode := make([]byte, 10)
			copy(subjectCode, strPDF[i-2:j])

			// Get to the end of the line
			for strPDF[j] != '\n' {
				j++
			}

			reGrade, _ := regexp.Compile("[0-9][0-9]?\\.[0-9]")
			grade := reGrade.FindString(strPDF[i+3 : j])

			// If grade was found
			if grade != "" {
				reStatus, _ := regexp.Compile("[A-Z]{1,4}")
				status := reStatus.FindString(strPDF[j-8 : j])

				gradeFloat, err := strconv.ParseFloat(grade, 64)

				// if grade parse succeeded and there's a status code
				if err == nil && status != "" {
					g := Grade{
						subject: string(subjectCode),
						grade:   gradeFloat,
						status:  status,
					}

					st.Grades = append(st.Grades, g)
				}

			}

			i = j
		} else {
			i++
		}
	}

	return st, true
}

func parseNUSP(row string) (string, error) {
	r, err := regexp.Compile("\\d+\\/")
	if err != nil {
		return "", err
	}

	return r.FindString(row), nil
}

func nuspInRow(i int, body *string) (idx int, ok bool) {
	var j int = i
	var found bool = false
	for j < len(*body) && (*body)[j] != '\n' {
		if strings.HasPrefix((*body)[i:j], "Aluno:") {
			found = true
		}
		j++
	}

	if found {
		return j, true
	}

	return -1, false
}

func isSubject(i int, body *string) bool {
	subjects := [4]string{"SMA", "SME", "SCC", "SSC"}

	for _, sub := range subjects {
		// Get last  three characters
		if sub == (*body)[i-2:i+1] {
			return true
		}
	}

	return false
}
