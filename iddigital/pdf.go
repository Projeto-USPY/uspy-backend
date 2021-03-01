/* package iddigital contains all logic necessary to interact with uspdigital */
package iddigital

import (
	"errors"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// iddigital.PDF represents the pdf file retrieved from uspdigital
// See PostAuthCode for more info.
type PDF struct {
	Body         string
	Error        error
	CreationDate time.Time
}

// iddigital.Records represents the parsed data retrieved from the user's PDF file\
type Records struct {
	Grades []entity.Grade `json:"grades"`
	Name   string         `json:"name"`
	Nusp   string         `json:"nusp"`
}

// NewPDF takes the Grades PDF response object and creates a new PDF object
func NewPDF(r *http.Response) (pdf PDF) {
	defer func() {
		if r := recover(); r != nil {
			pdf.Body = ""
			pdf.Error = r.(error)
			pdf.CreationDate = time.Now()
		}
	}()

	bodyPDF, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(errors.New("error reading pdf response body: " + err.Error()))
	}

	// transform PDF to string
	parser := exec.Command("pdftotext", "-q", "-eol", "unix", "-enc", "UTF-8", "-layout", "-", "-")
	stdin, _ := parser.StdinPipe()
	_, _ = stdin.Write(bodyPDF)
	_ = stdin.Close()
	parsed, err := parser.Output()

	if err != nil {
		panic(errors.New("error parsing pdf: " + err.Error()))
	}

	body := string(parsed)

	// Get PDF CreationDate in ISO format
	dataExtractor := exec.Command("pdfinfo", "-isodates", "-")
	stdin, _ = dataExtractor.StdinPipe()
	_, _ = stdin.Write(bodyPDF)
	_ = stdin.Close()
	meta, err := dataExtractor.Output()

	if err != nil {
		panic(errors.New("error getting pdf info: " + err.Error()))
	}

	var creation time.Time
	lines := strings.Split(string(meta), "\n")
	for _, v := range lines {
		fields := strings.SplitN(v, ":", 2)
		fields[0] = strings.Trim(fields[0], " \n\t")
		fields[1] = strings.Trim(fields[1], " \n\t")
		if fields[0] == "CreationDate" {
			layout := "2006-01-02T15:04:05-0700"
			c, err := time.Parse(layout, fields[1]+"00") // must add 00 to adapt to timezone layout
			if err != nil {
				panic(errors.New("error parsing time: " + err.Error()))
			} else {
				creation = c
				break
			}
		}
	}

	return PDF{
		Body:         body,
		Error:        nil,
		CreationDate: creation,
	}
}

// Parse takes the (already read) PDF and parses it into a Records object
func (pdf PDF) Parse(DB db.Env) (rec Records, err error) {
	defer func() {
		if r := recover(); r != nil {
			rec = Records{nil, "", ""}
			err = r.(error)
		}
	}()

	// Look for User NUSP (user identifier) and User Name in PDF Header
	nuspMatches := regexp.MustCompile(`Aluno:\s+(\d+)/\d - (.+)`).FindStringSubmatch(pdf.Body)

	if nuspMatches == nil || len(nuspMatches) < 3 {
		panic(errors.New("could not parse user nusp and/or name"))
	}

	rec.Nusp, rec.Name = nuspMatches[1], nuspMatches[2]

	// Look for course code in PDF Header
	matches := regexp.MustCompile(`Curso:\s+(\d+)/(\d+) - .*`).FindStringSubmatch(pdf.Body)

	if matches == nil || len(matches) < 3 {
		panic(errors.New("could not parse user course code and/or specialization"))
	}

	course, specialization := matches[1], matches[2]
	snaps, err := DB.RestoreCollection("courses")
	if err != nil {
		panic(errors.New("could not fetch courses from firestore"))
	}

	// Divide records data into each semester/year
	pairs := regexp.MustCompile(`\s+\d{4} [1-2]º\. Semestre\s+`).FindAllStringIndex(pdf.Body, -1)

	for i := 0; i < len(pairs)-1; i++ {
		l, r := pairs[i][0], pairs[i+1][0]

		// get current year and semester
		info := regexp.MustCompile(`(\d{4}) ([1-2])º\. Semestre`).FindStringSubmatch(pdf.Body[pairs[i][0]:pairs[i][1]])

		year, _ := strconv.Atoi(info[1])
		semester, _ := strconv.Atoi(info[2])

		// get all subjects in current year and semester
		subRXP := regexp.MustCompile(`((?:SMA|SME|SSC|SCC)\d+).*`)
		gradeRows := subRXP.FindAllStringSubmatch(pdf.Body[l:r], -1)

		for _, match := range gradeRows {
			row, subCode := match[0], match[1]

			// get subject values (grade, frequency and status)
			gradeRXP := regexp.MustCompile(`(\d{1,3})\s+(\d{1,2}.\d{1,2}) ([A-Z]+)`)
			values := gradeRXP.FindStringSubmatch(row)

			freq, _ := strconv.Atoi(values[1])
			grade, _ := strconv.ParseFloat(values[2], 64)
			status := values[3]
			subCourse := ""
			subSpecialization := ""

			// determine subject course origin
			for _, s := range snaps {
				c := entity.Course{}
				_ = s.DataTo(&c)
				_, exists := c.SubjectCodes[subCode]

				if exists {
					if c.Code == course && c.Specialization == specialization { // if subject is from students course, then subject's course should be it
						subCourse = c.Code
						subSpecialization = c.Specialization
					} else if subCourse == "" && subSpecialization == "" { // otherwise choose any course to be the subject course
						subCourse = c.Code
						subSpecialization = c.Specialization
					}
				}
			}

			rec.Grades = append(rec.Grades, entity.Grade{
				Subject:        subCode,
				Grade:          grade,
				Frequency:      freq,
				Status:         status,
				Course:         subCourse,
				Specialization: subSpecialization,
				Semester:       semester,
				Year:           year,
			})
		}
	}

	return
}
