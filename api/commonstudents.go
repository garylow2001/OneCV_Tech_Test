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
			c.JSON(http.StatusBadRequest, gin.H{"message": "At least one teacher email must be provided"})
			return
		}

		// Retrieve common students
		var commonStudents []string
		for i, teacherEmail := range teachers {
			// Students under teacher represent the students that are registered with the currently selected teacher
			var studentsUnderTeacher []struct {
				Email string
			}
			if err := db.Table("teachers").
				Select("students.email").
				Distinct().
				Joins("inner join teacher_student_relations on teachers.id = teacher_student_relations.teacher_id").
				Joins("inner join students on teacher_student_relations.student_id = students.id").
				Where("teachers.email = ?", teacherEmail).
				Scan(&studentsUnderTeacher).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve students"})
				return
			}

			if i == 0 {
				// Initialize common students with students from the first teacher
				for _, student := range studentsUnderTeacher {
					commonStudents = append(commonStudents, student.Email)
				}
			} else {
				// Keep only common students
				var newCommonStudents []string
				for _, student := range studentsUnderTeacher {
					if contains(commonStudents, student.Email) {
						newCommonStudents = append(newCommonStudents, student.Email)
					}
				}
				commonStudents = newCommonStudents
			}
		}

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
