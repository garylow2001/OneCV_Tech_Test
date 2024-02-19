package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/garylow2001/OneCV_Tech_Test/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterEndpoint(t *testing.T) {
	// Test valid registration
	t.Run("Success", func(t *testing.T) {
		r, mock := setUpRouters("/api/register", api.RegisterHandler)
		defer mock.ExpectClose()
		testValidRegistrationSuccess(r, t, mock)
	})

	// Test missing teacher input
	t.Run("MissingTeacherInput", func(t *testing.T) {
		r, mock := setUpRouters("/api/register", api.RegisterHandler)
		defer mock.ExpectClose()
		testMissingTeacherInputFailure(r, t)
	})

	// Test missing student input
	t.Run("MissingStudentInput", func(t *testing.T) {
		r, mock := setUpRouters("/api/register", api.RegisterHandler)
		defer mock.ExpectClose()
		testMissingStudentInputFailure(r, t, mock)
	})
}

func testValidRegistrationSuccess(r *gin.Engine, t *testing.T, mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`^SELECT \* FROM "teachers" WHERE email = \$1 AND "teachers"."deleted_at" IS NULL ORDER BY "teachers"."id" LIMIT \$2`).
		WithArgs("teacherken@gmail.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "teacherken@gmail.com"))
	mock.ExpectQuery(`^SELECT \* FROM "students" WHERE email = \$1 AND "students"."deleted_at" IS NULL ORDER BY "students"."id" LIMIT \$2`).
		WithArgs("studentjon@gmail.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "studentjon@gmail.com"))
	mock.ExpectQuery(`^SELECT \* FROM "students" WHERE email = \$1 AND "students"."deleted_at" IS NULL ORDER BY "students"."id" LIMIT \$2`).
		WithArgs("studenthon@gmail.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "studenthon@gmail.com"))

	testRequestForRegister(r, t, "teacherken@gmail.com", []string{"studentjon@gmail.com", "studenthon@gmail.com"}, http.StatusNoContent, mock)
}

func testMissingTeacherInputFailure(r *gin.Engine, t *testing.T) {
	testRequestForRegister(r, t, "", []string{"studentjon@gmail.com"}, http.StatusBadRequest, nil)
}

func testMissingStudentInputFailure(r *gin.Engine, t *testing.T, mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`^SELECT \* FROM "teachers" WHERE email = \$1 AND "teachers"."deleted_at" IS NULL ORDER BY "teachers"."id" LIMIT \$2`).
		WithArgs("teacherken@gmail.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "teacherken@gmail.com"))
	testRequestForRegister(r, t, "teacherken@gmail.com", []string{}, http.StatusBadRequest, nil)
}

func testRequestForRegister(r *gin.Engine, t *testing.T, teacherEmail string, students []string, expectedStatus int, mock sqlmock.Sqlmock) {
	payload := map[string]interface{}{
		"teacher":  teacherEmail,
		"students": students,
	}
	reqBody, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, expectedStatus, w.Code)
}
