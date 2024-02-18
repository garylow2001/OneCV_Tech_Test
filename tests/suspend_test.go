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

func TestSuspendEndpoint(t *testing.T) {
	// Test valid suspension
	t.Run("Success", func(t *testing.T) {
		r, mock := setupTestRouter()
		defer mock.ExpectClose()
		testSuspendSuccess(r, t, "student1@gmail.com", mock)
	})

	// Test student not found
	t.Run("StudentNotFound", func(t *testing.T) {
		r, mock := setupTestRouter()
		defer mock.ExpectClose()
		testStudentNotFoundFailure(r, t, "nonexistent@gmail.com", mock)
	})

	// Test student already suspended
	t.Run("StudentAlreadySuspended", func(t *testing.T) {
		r, mock := setupTestRouter()
		defer mock.ExpectClose()
		testStudentAlreadySuspendedFailure(r, t, "student1@gmail.com", mock)
	})
}

func setupTestRouter() (*gin.Engine, sqlmock.Sqlmock) {
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	r := gin.Default()
	r.POST("/api/suspend", func(c *gin.Context) { api.SuspendHandler(db)(c) })

	return r, mock
}

func testSuspendSuccess(r *gin.Engine, t *testing.T, studentEmail string, mock sqlmock.Sqlmock) {
	// Set up expectations
	mock.ExpectQuery("^SELECT \\* FROM \"students\" WHERE email = \\$1 AND \"students\".\"deleted_at\" IS NULL ORDER BY \"students\".\"id\" LIMIT \\$2$").
		WithArgs(studentEmail, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "suspended"}).
			AddRow(1, studentEmail, false))

	mock.ExpectExec("^UPDATE students SET suspended = true WHERE email = (.+)$").
		WithArgs(studentEmail).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Test the request
	testRequest(r, t, studentEmail, http.StatusNoContent)
}

func testStudentNotFoundFailure(r *gin.Engine, t *testing.T, studentEmail string, mock sqlmock.Sqlmock) {
	// Set up expectations
	mock.ExpectExec("^UPDATE students SET suspended = true WHERE email = (.+)$").
		WithArgs(studentEmail).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Test the request
	testRequest(r, t, studentEmail, http.StatusNotFound)
}

func testStudentAlreadySuspendedFailure(r *gin.Engine, t *testing.T, studentEmail string, mock sqlmock.Sqlmock) {
	// Set up expectations
	mock.ExpectQuery("^SELECT \\* FROM \"students\" WHERE email = \\$1 AND \"students\".\"deleted_at\" IS NULL ORDER BY \"students\".\"id\" LIMIT \\$2$").
		WithArgs(studentEmail, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "suspended"}).
			AddRow(1, studentEmail, true)) // The student is already suspended

	// Test the request
	testRequest(r, t, studentEmail, http.StatusConflict)
}

func testRequest(r *gin.Engine, t *testing.T, studentEmail string, expectedStatus int) {
	requestBody, _ := json.Marshal(map[string]string{
		"student": studentEmail,
	})
	req, _ := http.NewRequest("POST", "/api/suspend", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, expectedStatus, w.Code)
}
