package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CommonStudentsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		teachers := c.QueryArray("teacher")

		// Validate query parameters
		if len(teachers) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one teacher email must be provided"})
			return
		}

		// Retrieve common students
		var commonStudents []string
		for i, teacherEmail := range teachers {
			var teacherStudents []struct {
				StudentEmail string
			}
			if err := db.Table("teacher_students").Select("student_email").Where("teacher_email = ?", teacherEmail).Scan(&teacherStudents).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
				return
			}

			if i == 0 {
				// Initialize common students with students from the first teacher
				for _, student := range teacherStudents {
					commonStudents = append(commonStudents, student.StudentEmail)
				}
			} else {
				// Keep only common students
				var newCommonStudents []string
				for _, student := range teacherStudents {
					if contains(commonStudents, student.StudentEmail) {
						newCommonStudents = append(newCommonStudents, student.StudentEmail)
					}
				}
				commonStudents = newCommonStudents
			}
		}

		// Return response
		c.JSON(http.StatusOK, gin.H{"students": commonStudents})
	}
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
