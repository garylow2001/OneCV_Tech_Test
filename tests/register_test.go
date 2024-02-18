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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRegisterEndpoint(t *testing.T) {
	// Create a mock database connection
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	// Set up a test router
	r := gin.Default()
	r.POST("/api/register", func(c *gin.Context) { api.RegisterHandler(db)(c) })

	// Test valid registration
	testValidRegistrationSuccess(r, t, mock)

	// Test missing teacher input
	testMissingTeacherInputFailure(r, t)

	// Add more test cases as needed (e.g., invalid student emails, database errors, etc.)
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
	payload := map[string]interface{}{
		"teacher": "teacherken@gmail.com",
		"students": []string{
			"studentjon@gmail.com",
			"studenthon@gmail.com",
		},
	}
	reqBody, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func testMissingTeacherInputFailure(r *gin.Engine, t *testing.T) {
	payload := map[string]interface{}{
		"students": []string{
			"studentjon@gmail.com",
		},
	}
	reqBody, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
