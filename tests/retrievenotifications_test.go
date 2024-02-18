package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/garylow2001/OneCV_Tech_Test/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRetrieveForNotificationsHandler(t *testing.T) {
	// Test valid notification retrieval
	t.Run("Success", func(t *testing.T) {
		r, mock := setUpRouters("/api/retrievefornotifications", api.RetrieveForNotificationsHandler)
		defer mock.ExpectClose()
		testMessage := "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com"
		testRetrieveForNotificationsSuccess(r, t, "teacherken@gmail.com", testMessage, mock)
	})

	// Test invalid notification format (missing @)
	t.Run("InvalidNotificationFormat", func(t *testing.T) {
		r, mock := setUpRouters("/api/retrievefornotifications", api.RetrieveForNotificationsHandler)
		defer mock.ExpectClose()
		testMessageMissingSymbol := "Hello students studentagnes@gmail.com @studentmiche@gmail.com"
		testRetrieveForNotificationsInvalidFormat(r, t, "teacherken@gmail.com", testMessageMissingSymbol, mock)
	})

	// Test student not found
	t.Run("StudentNotFound", func(t *testing.T) {
		r, mock := setUpRouters("/api/retrievefornotifications", api.RetrieveForNotificationsHandler)
		defer mock.ExpectClose()
		testMessage := "Hello students! @nonexistent@gmail.com"
		testRetrieveForNotificationsStudentNotFound(r, t, "teacherken@gmail.com", testMessage, mock)
	})
}

func testRetrieveForNotificationsSuccess(r *gin.Engine, t *testing.T, teacherEmail, notification string, mock sqlmock.Sqlmock) {
	// Set up expectations
	mock.ExpectQuery("^SELECT \\* FROM \"students\" WHERE email = \\$1 ORDER BY \"students\".\"suspended\" LIMIT \\$2$").
		WithArgs("studentagnes@gmail.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "suspended"}).
			AddRow(1, "studentagnes@gmail.com", false))
	mock.ExpectQuery("^SELECT \\* FROM \"students\" WHERE email = \\$1 ORDER BY \"students\".\"suspended\" LIMIT \\$2$").
		WithArgs("studentmiche@gmail.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "suspended"}).
			AddRow(2, "studentmiche@gmail.com", false))
	mock.ExpectQuery("^SELECT DISTINCT students.email FROM \"teachers\" inner join teacher_student_relations on teachers.id = teacher_student_relations.teacher_id inner join students on teacher_student_relations.student_id = students.id WHERE teachers.email = \\$1 AND students.suspended = \\$2$").
		WithArgs("teacherken@gmail.com", false).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow("studentbob@gmail.com"))

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
	testRequestForNotifications(r, t, teacherEmail, notification, http.StatusBadRequest, nil, mock)
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

		sort.Strings(expectedRecipients)
		sort.Strings(response.Recipients)
		assert.Equal(t, expectedRecipients, response.Recipients)
	}
}
