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

func TestRetrieveForNotificationsHandler(t *testing.T) {
	// Create a mock database connection
	mockDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer mockDB.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	// Set up a test router
	r := gin.Default()
	r.POST("/api/retrievefornotifications", api.RetrieveForNotificationsHandler(db))

	// Test valid notification retrieval
	testMessage := "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com"
	testRetrieveForNotificationsSuccess(r, t, "teacherken@gmail.com", testMessage, mock)

	// Test invalid notification format (missing @)
	testMessageMissingSymbol := "Hello students studentagnes@gmail.com @studentmiche@gmail.com"
	testRetrieveForNotificationsInvalidFormat(r, t, "teacherken@gmail.com", testMessageMissingSymbol, mock)

	// Test student not found
	testRetrieveForNotificationsStudentNotFound(r, t, "teacherken@gmail.com", testMessage, mock)
}

func testRetrieveForNotificationsSuccess(r *gin.Engine, t *testing.T, teacherEmail, notification string, mock sqlmock.Sqlmock) {
	// Set up expectations
	mock.ExpectQuery("^SELECT \\* FROM \"students\" WHERE email = (.+) ORDER BY \"students\"\\.\"suspended\" LIMIT \\d+$").
		WillReturnRows(sqlmock.NewRows([]string{"suspended"}).
			AddRow(false).
			AddRow(false))

	mock.ExpectQuery("^SELECT DISTINCT (.+) FROM \"teachers\" INNER JOIN teacher_student_relations ON teachers.id = teacher_student_relations.teacher_id INNER JOIN students ON teacher_student_relations.student_id = students.id WHERE teachers.email = (.+)$").
		WithArgs("teacherken@gmail.com").
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow("studentbob@gmail.com").
			AddRow("studentagnes@gmail.com").
			AddRow("studentmiche@gmail.com"))

	// Test the request
	expectedRecipients := []string{"studentbob@gmail.com", "studentagnes@gmail.com", "studentmiche@gmail.com"}
	testRequestForNotifications(r, t, teacherEmail, notification, http.StatusOK, expectedRecipients, mock)
}

func testRetrieveForNotificationsInvalidFormat(r *gin.Engine, t *testing.T, teacherEmail, notification string, mock sqlmock.Sqlmock) {
	// Test the request with invalid notification format
	testRequestForNotifications(r, t, teacherEmail, notification, http.StatusBadRequest, nil, mock)
}

func testRetrieveForNotificationsStudentNotFound(r *gin.Engine, t *testing.T, teacherEmail, notification string, mock sqlmock.Sqlmock) {
	// Set up expectations for student not found scenario
	mock.ExpectQuery("^SELECT \\* FROM \"students\" WHERE email = (.+)$").
		WillReturnError(gorm.ErrRecordNotFound)

	// Test the request
	testRequestForNotifications(r, t, teacherEmail, notification, http.StatusInternalServerError, nil, mock)
}

func testRequestForNotifications(r *gin.Engine, t *testing.T, teacherEmail, notification string, expectedStatus int, expectedRecipients []string, mock sqlmock.Sqlmock) {
	requestBody, _ := json.Marshal(map[string]string{
		"teacher":      teacherEmail,
		"notification": notification,
	})
	req, _ := http.NewRequest("POST", "/api/retrievefornotifications", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, expectedStatus, w.Code)

	if expectedRecipients != nil {
		var response struct {
			Recipients []string `json:"recipients"`
		}
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, expectedRecipients, response.Recipients)
	}
}
