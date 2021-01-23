package pdfparser

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PDF struct {
	Body *string
	Error error
	CreationDate time.Time
}

// Grade represents a Grade in jupiterweb
type Grade struct {
	Subject  string  `json:"subject"`
	Grade    float64 `json:"grade"`
	Status   string  `json:"status"`
	Semester int     `json:"semester"`
	Year     int     `json:"year"`
}

type Records struct {
	Grades []Grade `json:"grades"`
	Nusp   string  `json:"nusp"`
}

// NewPDF takes the Grades PDF response object and creates a new PDF object
func NewPDF(r *http.Response) (pdf PDF) {
	defer func() {
		if r := recover(); r != nil {
			pdf.Body = nil
			pdf.Error = r.(error)
			pdf.CreationDate = time.Now()
		}
	}()

	bodyPDF, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic("error converting response body to string")
	}

	parser := exec.Command("pdftotext", "-q", "-eol", "unix", "-enc", "UTF-8", "-layout", "-", "-")
	stdin, _ := parser.StdinPipe()
	_, _ = stdin.Write(bodyPDF)
	_ = stdin.Close()
	parsed, err := parser.Output()

	if err != nil {
		panic("an error occured while executing pdftotext")
	}

	body := string(parsed)

	dataExtractor := exec.Command("pdfinfo", "-", "-isodates")
	stdin, _ = dataExtractor.StdinPipe()
	_, _ = stdin.Write(bodyPDF)
	_ = stdin.Close()
	meta, err := dataExtractor.Output()

	if err != nil {
		panic("an error occured while executing pdftotext")
	}

	var creation time.Time
	lines := strings.Split(string(meta), "\n")
	for _, v := range lines {
		fields := strings.SplitN(v, ":", 2)
		fields[0] = strings.Trim(fields[0], " \n\t")
		fields[1] = strings.Trim(fields[1], " \n\t")
		if fields[0] == "CreationDate" {
			layout := "2006-01-02T15:04:05-0700"
			creation, err = time.Parse(layout, fields[1] + "00")
			break
		}
	}

	return PDF{
		Body: &body,
		Error: nil,
		CreationDate: creation,
	}
}

// ParsePDF takes the (already read) PDF and parses it into records
func (pdf PDF) ParsePDF() (rec Records, err error) {
	defer func() {
		if r := recover(); r != nil {
			rec = Records{nil, ""}
			err = r.(error)
		}
	}()

	i := 2
	strPDF := *pdf.Body
	rec.Nusp = ""
	semester, year := -1, -1

	for {
		// End of PDF
		if i == len(strPDF) {
			break
		}

		if idx, ok := nuspInRow(i, pdf.Body); ok && rec.Nusp == "" {
			nusp, err := parseNUSP((*pdf.Body)[i:idx])

			if err != nil {
				panic("couldnt parse nusp")
			}

			rec.Nusp = nusp[:len(nusp)-1]
		}

		s, y, foundSemester := semesterInRow(i, pdf.Body)
		if foundSemester {
			semester, year = s, y
		}

		// Found a subject
		if isSubject(i, pdf.Body) {
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
						Subject:  string(subjectCode),
						Grade:    gradeFloat,
						Status:   status,
						Semester: semester,
						Year:     year,
					}

					rec.Grades = append(rec.Grades, g)
				}

			}

			i = j
		} else {
			i++
		}
	}

	return
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

func semesterInRow(i int, body *string) (int, int, bool) {
	var j int = i
	for j < len(*body) && (*body)[j] != '\n' {
		j++
	}
	row := (*body)[i:j]
	cmp, err := regexp.Compile("\\d{4} [1-2]º\\. Semestre")
	if err != nil {
		return -1, -1, false
	}

	bytes := cmp.Find([]byte(row))
	if bytes == nil {
		return -1, -1, false
	}

	cmp, err = regexp.Compile("\\d+")
	if err != nil {
		return -1, -1, false
	}

	values := cmp.FindAllString(row, 2)
	parsedYear, _ := strconv.ParseInt(values[0], 10, 32)
	parsedSemester, _ := strconv.ParseInt(values[1], 10, 32)

	return int(parsedYear), int(parsedSemester), true
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
