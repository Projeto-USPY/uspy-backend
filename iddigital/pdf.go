/*Package iddigital contains all logic necessary to interact with uspdigital */
package iddigital

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// PDF represents the pdf file retrieved from uspdigital
// See PostAuthCode for more info.
type PDF struct {
	Body         string
	Error        error
	CreationDate time.Time
}

// Transcript represents the parsed data retrieved from the user's PDF file
type Transcript struct {
	Grades          []models.Record  `json:"grades"`
	TranscriptYears map[string][]int `json:"transcript_years"`

	Name string `json:"name"`
	Nusp string `json:"nusp"`

	Course         string `json:"course"`
	Specialization string `json:"specialization"`
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

// Parse takes the (already read) PDF and parses it into a Transcript object
func (pdf PDF) Parse(DB db.Env) (rec Transcript, err error) {
	defer func() {
		if r := recover(); r != nil {
			rec = Transcript{Grades: nil, TranscriptYears: nil, Name: "", Nusp: "", Course: "", Specialization: ""}
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
	// this does not scale well for multiple institutes
	// TODO: change field subjects in courses to a subcollection to change it to a collection group query
	snaps, err := DB.Client.CollectionGroup("courses").Documents(DB.Ctx).GetAll()

	if err != nil {
		panic(errors.New("could not fetch courses from firestore"))
	}

	rec.Course, rec.Specialization = course, specialization

	// Divide records data into each semester/year
	pairs := regexp.MustCompile(`\s+\d{4} [1-2]º\. Semestre\s+`).FindAllStringIndex(pdf.Body, -1)

	for i := 0; i < len(pairs); i++ {
		l := pairs[i][0]

		var r int
		if i+1 < len(pairs) {
			r = pairs[i+1][0]
		} else {
			r = len(pdf.Body)
		}

		// get current year and semester
		info := regexp.MustCompile(`(\d{4}) ([1-2])º\. Semestre`).FindStringSubmatch(pdf.Body[pairs[i][0]:pairs[i][1]])

		year := utils.MustAtoi(info[1])
		semester := utils.MustAtoi(info[2])

		// get all subjects in current year and semester
		subRXP := regexp.MustCompile(`([0-9A-Z]{5,10}).*`)
		gradeRows := subRXP.FindAllStringSubmatch(pdf.Body[l:r], -1)

		for _, match := range gradeRows {
			row, subCode := match[0], match[1]

			// get subject values (grade, frequency and status)
			gradeRXP := regexp.MustCompile(`(\d{1,3})\s+(\d{1,2}\.\d{1,2}) ([A-Z]+)`)
			values := gradeRXP.FindStringSubmatch(row)

			if len(values) < 4 { // array must be [whole string, grade, frequency, status]
				continue
			}

			freq := utils.MustAtoi(values[1])
			grade, _ := strconv.ParseFloat(values[2], 64)
			status := values[3]
			subCourse := ""
			subSpecialization := ""

			// determine subject course origin
			for _, s := range snaps {
				c := models.Course{}
				_ = s.DataTo(&c)
				_, exists := c.SubjectCodes[subCode]

				if exists {
					if c.Code == course && c.Specialization == specialization { // perfect match, there's a subject with exact course and specialization
						subCourse = c.Code
						subSpecialization = c.Specialization
						break
					} else if subCourse == "" && subSpecialization == "" { // get from any major that contains this subject code
						subCourse = c.Code
						subSpecialization = c.Specialization
					} else if c.Code == course { // replace any by a more likely match (this is useful for majors that have "ciclos básicos")
						subCode = c.Code
						subSpecialization = c.Specialization
					}
				}
			}

			rec.Grades = append(rec.Grades, models.Record{
				Subject:        subCode,
				Grade:          grade,
				Frequency:      freq,
				Status:         status,
				Course:         subCourse,
				Specialization: subSpecialization,
				Semester:       semester,
				Year:           year,
			})

			yearStr := strconv.Itoa(year)
			// append to transcript years it necessary
			if _, ok := rec.TranscriptYears[yearStr]; !ok {
				if rec.TranscriptYears == nil {
					rec.TranscriptYears = make(map[string][]int)
				}

				rec.TranscriptYears[yearStr] = make([]int, 0, 2)
			}

			found := false
			for i := 0; i < utils.Min(2, len(rec.TranscriptYears[yearStr])); i++ {
				if rec.TranscriptYears[yearStr][i] == semester { // already appended
					found = true
					break
				}
			}

			if !found { // if semester was not added yet
				rec.TranscriptYears[yearStr] = append(rec.TranscriptYears[yearStr], semester)
			}

		}
	}

	return
}
