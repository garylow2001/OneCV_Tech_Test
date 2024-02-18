package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/garylow2001/OneCV_Tech_Test/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCommonStudentsEndpoint(t *testing.T) {
	// Test common students between two teachers
	t.Run("Success", func(t *testing.T) {
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
	t.Run("Success2", func(t *testing.T) {
		r, mock := setUpRouters("/api/commonstudents", api.CommonStudentsHandler)
		defer mock.ExpectClose()
		testCommonStudentsOneTeacherSuccess(r, t,
			"teacherken@gmail.com",
			[]string{"commonstudent1@gmail.com", "commonstudent2@gmail.com", "student_only_under_teacher_ken@gmail.com"},
			mock,
		)
	})
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
	req, _ := http.NewRequest("GET", "/api/commonstudents?teacher="+teacherStringGroup, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Students []string `json:"students"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, expectedStudents, response.Students)
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
	req, _ := http.NewRequest("GET", "/api/commonstudents?teacher="+teacherEmail, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Students []string `json:"students"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, expectedStudents, response.Students)
}
