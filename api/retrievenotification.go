package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RetrieveForNotificationsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse request body
		var requestData struct {
			TeacherEmail string `json:"teacher" binding:"required"`
			Notification string `json:"notification" binding:"required"`
		}
		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Get mentioned students in notification using regex
		mentionedStudentsEmails := getMentionedStudentsInNotification(requestData.Notification)

		// Validate mentioned students and remove suspended students from the list
		studentDoesNotExistErr := validateMentionedStudents(mentionedStudentsEmails, db, c)
		if studentDoesNotExistErr != nil {
			return
		}

		// Retrieve students who are registered with the teacher
		var studentsUnderTeacherEmails []string
		if err := db.Table("teachers").
			Select("students.email").
			Distinct().
			Joins("inner join teacher_student_relations on teachers.id = teacher_student_relations.teacher_id").
			Joins("inner join students on teacher_student_relations.student_id = students.id").
			Where("teachers.email = ? AND students.suspended = ?", requestData.TeacherEmail, false).
			Scan(&studentsUnderTeacherEmails).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve students"})
			return
		}

		// Concatenate mentioned students and students under teacher and return the unique recipients
		recipients := uniqueConcat(mentionedStudentsEmails, studentsUnderTeacherEmails)

		c.JSON(http.StatusOK, gin.H{"recipients": recipients})
	}
}

func getMentionedStudentsInNotification(notification string) []string {
	var mentionedStudents []string
	re := regexp.MustCompile(`@(\w+@\w+\.\w+)`)
	// matches is a 2D slice, where each element is a slice containing the full match(@eg@example.com) and the captured group
	matches := re.FindAllStringSubmatch(notification, -1)
	// Extract the email from each match (which is the captured group at index 1 in each match)
	for _, match := range matches {
		mentionedStudents = append(mentionedStudents, match[1])
	}
	return mentionedStudents
}

func uniqueConcat(mentionedStudentsEmails []string, studentsUnderTeacherEmails []string) []string {
	emailMap := make(map[string]bool) // Make a map of emails to handle duplicates
	for _, email := range mentionedStudentsEmails {
		emailMap[email] = true
	}
	for _, email := range studentsUnderTeacherEmails {
		emailMap[email] = true
	}

	var recipients []string
	for email := range emailMap {
		recipients = append(recipients, email)
	}
	return recipients
}

func validateMentionedStudents(mentionedStudentsEmails []string, db *gorm.DB, c *gin.Context) error {
	for i, studentEmail := range mentionedStudentsEmails {
		var student struct {
			Suspended bool
		}
		if err := db.Table("students").Where("email = ?", studentEmail).First(&student).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to retrieve student with email: %s", studentEmail)})
			return err
		}
		if student.Suspended {
			mentionedStudentsEmails = append(mentionedStudentsEmails[:i], mentionedStudentsEmails[i+1:]...)
		}
	}
	return nil
}
