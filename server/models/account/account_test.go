package account_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/models/account"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/Projeto-USPY/uspy-backend/utils/test"
	"github.com/Projeto-USPY/uspy-backend/utils/test/emulator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AccountSuite struct {
	suite.Suite
	DB          db.Env
	router      *gin.Engine
	accessToken *http.Cookie
}

// SetupTest runs before every test
func (s *AccountSuite) SetupTest() {
	s.DB, s.router, s.accessToken = test.MustGetEnvironment(s.Suite)
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountSuite))
}

func (s *AccountSuite) TestUpdateUser() {
	timezone, err := time.LoadLocation("America/Sao_Paulo")
	s.Require().NoError(err)
	creationDate := time.Date(2021, time.January, 1, 0, 0, 0, 0, timezone)

	var data = emulator.Transcript
	newRecord := models.Record{
		Subject:        "SCC0230",
		Course:         "55090",
		Specialization: "0",
		Year:           2021,
		Semester:       1,
		Grade:          7.5,
		Status:         "A",
		Frequency:      90,
	}

	// add new grade and change major to simulate new transcript
	data.Grades = append(data.Grades, newRecord)
	data.Course = "55090"
	data.Specialization = "0"

	err = account.UpdateUser(s.DB.Ctx, s.DB, &data, "123456789", creationDate)
	s.Require().NoError(err, "failed to update user")

	// ensure user had its last_update property updated
	var storedUser models.User
	snap, err := s.DB.Restore("users/" + utils.SHA256("123456789"))
	s.Require().NoError(err)

	err = snap.DataTo(&storedUser)
	s.Require().NoError(err)
	s.Require().Equal(creationDate, storedUser.LastUpdate.In(timezone))

	subHash := models.Subject{Code: "SCC0230", CourseCode: "55090", Specialization: "0"}.Hash()

	// ensure record was created
	expectedRecord := models.Record{
		Frequency: 90,
		Grade:     7.5,
		Status:    "A",
	}

	var storedRecord models.Record
	snap, err = s.DB.Restore(
		fmt.Sprintf(
			"users/%s/final_scores/%s/records/%s",
			utils.SHA256("123456789"),
			subHash,
			newRecord.Hash(),
		),
	)
	s.Require().NoError(err)

	err = snap.DataTo(&storedRecord)
	s.Require().NoError(err)

	s.Require().Equal(expectedRecord, storedRecord)

	// ensure subject grade was created
	snaps, err := s.DB.RestoreCollection("subjects/" + subHash + "/grades")
	s.Require().NoError(err)
	s.Require().Len(snaps, 1)

	// ensure major was created
	expectedMajor := models.Major{
		Code:           "55090",
		Specialization: "0",
	}

	var storedMajor models.Major
	snap, err = s.DB.Restore(
		fmt.Sprintf(
			"users/%s/majors/%s",
			utils.SHA256("123456789"),
			expectedMajor.Hash(),
		),
	)
	s.Require().NoError(err)

	err = snap.DataTo(&storedMajor)
	s.Require().NoError(err)

	s.Require().Equal(expectedMajor, storedMajor)
}
