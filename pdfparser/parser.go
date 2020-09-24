package pdfparser

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

// ReadPDF takes the filename of a  PDF and returns its string
func ReadPDF(file string) (body *string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recovered] Coudlnt read PDF: %v\n", r)
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

// Grade represents a Grade in jupiterweb
type Grade struct {
	subject string
	grade   float64
	status  string
}

// ParsePDF takes the (already read) PDF string and parses it to a list of Grades
func ParsePDF(body *string) (grades []Grade, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recovered] Couldnt parse PDF string: %v\n", r)
			grades = nil
			ok = false
		}
	}()

	i := 2
	strPDF := *body

	for {
		// End of PDF
		if i == len(strPDF) {
			break
		}

		// Found a subject
		if isSubject(i, body) == true {
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

					grades = append(grades, g)
				}

			}

			i = j
		} else {
			i++
		}
	}

	return grades, true
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
