package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/garylow2001/OneCV_Tech_Test/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCommonStudentsEndpoint(t *testing.T) {
	// Test common students between two teachers
	t.Run("TwoTeacherSuccess", func(t *testing.T) {
		teacherStringGroup := "teacherken@gmail.com&teacher=teacherjoe@gmail.com"
		r, mock := setUpRouters("/api/commonstudents", api.CommonStudentsHandler)
		defer mock.ExpectClose()
		testCommonStudentsBetweenTwoTeachersSuccess(r, t,
			teacherStringGroup,
			[]string{"commonstudent1@gmail.com", "commonstudent2@gmail.com"},
			mock,
		)
	})

	// Test common students of one teacher
	t.Run("OneTeacherSuccess", func(t *testing.T) {
		r, mock := setUpRouters("/api/commonstudents", api.CommonStudentsHandler)
		defer mock.ExpectClose()
		testCommonStudentsOneTeacherSuccess(r, t,
			"teacherken@gmail.com",
			[]string{"commonstudent1@gmail.com", "commonstudent2@gmail.com", "student_only_under_teacher_ken@gmail.com"},
			mock,
		)
	})

	// Test no teacher provided
	t.Run("NoTeacherFailure", func(t *testing.T) {
		r, mock := setUpRouters("/api/commonstudents", api.CommonStudentsHandler)
		defer mock.ExpectClose()
		testNoTeacherFailure(r, t, mock)
	})
}

func testNoTeacherFailure(r *gin.Engine, t *testing.T, mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT DISTINCT students.email FROM "teachers" 
		inner join teacher_student_relations on teachers.id = teacher_student_relations.teacher_id 
		inner join students on teacher_student_relations.student_id = students.id 
		WHERE teachers.email = $1`,
	)).WithArgs("").WillReturnError(errors.New("At least one teacher email is required"))
	testRequestForCommonStudentsFailure(r, t, "", http.StatusBadRequest, mock)
}

func testCommonStudentsBetweenTwoTeachersSuccess(r *gin.Engine, t *testing.T, teacherStringGroup string, expectedStudents []string, mock sqlmock.Sqlmock) {
	// Set up expectations
	mock.ExpectQuery("^SELECT DISTINCT (.+) FROM \"teachers\" inner join teacher_student_relations on teachers.id = teacher_student_relations.teacher_id inner join students on teacher_student_relations.student_id = students.id WHERE teachers.email = (.+)$").
		WithArgs("teacherken@gmail.com").
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow("commonstudent1@gmail.com").
			AddRow("commonstudent2@gmail.com").
			AddRow("student_only_under_teacher_ken@gmail.com"))
	mock.ExpectQuery("^SELECT DISTINCT (.+) FROM \"teachers\" inner join teacher_student_relations on teachers.id = teacher_student_relations.teacher_id inner join students on teacher_student_relations.student_id = students.id WHERE teachers.email = (.+)$").
		WithArgs("teacherjoe@gmail.com").
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow("commonstudent1@gmail.com").
			AddRow("commonstudent2@gmail.com"))

	// Test the request
	testRequestForCommonStudentsSuccess(r, t, teacherStringGroup, expectedStudents, http.StatusOK, mock)
}

func testCommonStudentsOneTeacherSuccess(r *gin.Engine, t *testing.T, teacherEmail string, expectedStudents []string, mock sqlmock.Sqlmock) {
	// Set up expectations
	mock.ExpectQuery("^SELECT DISTINCT (.+) FROM \"teachers\" inner join teacher_student_relations on teachers.id = teacher_student_relations.teacher_id inner join students on teacher_student_relations.student_id = students.id WHERE teachers.email = (.+)$").
		WithArgs("teacherken@gmail.com").
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow("commonstudent1@gmail.com").
			AddRow("commonstudent2@gmail.com").
			AddRow("student_only_under_teacher_ken@gmail.com"))

	// Test the request
	testRequestForCommonStudentsSuccess(r, t, teacherEmail, expectedStudents, http.StatusOK, mock)
}

func testRequestForCommonStudentsSuccess(r *gin.Engine, t *testing.T, teacherStringGroup string,
	expectedStudents []string, expectedStatus int, mock sqlmock.Sqlmock) {
	req, _ := http.NewRequest("GET", "/api/commonstudents?teacher="+teacherStringGroup, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Students []string `json:"students"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, expectedStatus, w.Code)
	assert.Equal(t, expectedStudents, response.Students)
}

func testRequestForCommonStudentsFailure(r *gin.Engine, t *testing.T, teacherStringGroup string,
	expectedStatus int, mock sqlmock.Sqlmock) {
	req, _ := http.NewRequest("GET", "/api/commonstudents?teacher="+teacherStringGroup, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response struct {
		Students []string `json:"students"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, expectedStatus, w.Code)
}
