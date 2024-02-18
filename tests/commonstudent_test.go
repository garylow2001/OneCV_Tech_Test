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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCommonStudentsEndpoint(t *testing.T) {
	// Create a mock database connection
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	// Set up a test router
	r := gin.Default()
	r.GET("/api/commonstudents", func(c *gin.Context) { api.CommonStudentsHandler(db)(c) })

	// Test common students between two teachers
	testCommonStudentsBetweenTwoTeachersSuccess(r, t,
		"teacherken@gmail.com&teacher=teacherjoe@gmail.com",
		[]string{"commonstudent1@gmail.com", "commonstudent2@gmail.com"},
		mock,
	)

	// Test common students of one teacher
	testCommonStudentsOneTeacherSuccess(r, t,
		"teacherken@gmail.com",
		[]string{"commonstudent1@gmail.com", "commonstudent2@gmail.com", "student_only_under_teacher_ken@gmail.com"},
		mock,
	)
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
