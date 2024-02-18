package api

import (
	"net/http"

	"github.com/garylow2001/OneCV_Tech_Test/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Teacher  string   `json:"teacher" binding:"required"`
	Students []string `json:"students" binding:"required"`
}

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse request body
		var requestData RegisterRequest
		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request format"})
			return
		}

		// Check if teacher exists
		teacher, err := getOrCreateTeacher(db, requestData.Teacher)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create teacher"})
			return
		}

		// Register students
		for _, studentEmail := range requestData.Students {
			student, err := getOrCreateStudent(db, studentEmail)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create student"})
				return
			}
			db.Model(&teacher).Association("Students").Append(student)
		}

		c.Status(http.StatusNoContent)
	}
}

func getOrCreateTeacher(db *gorm.DB, teacherEmail string) (*models.Teacher, error) {
	var teacher models.Teacher
	result := db.Where("email = ?", teacherEmail).First(&teacher)
	if result.Error != nil {
		// Create teacher if not found
		teacher = models.Teacher{Email: teacherEmail}
		createTeacherError := db.Create(&teacher)
		if createTeacherError.Error != nil {
			return nil, createTeacherError.Error
		}
	}
	return &teacher, nil
}

func getOrCreateStudent(db *gorm.DB, studentEmail string) (*models.Student, error) {
	var student models.Student
	result := db.Where("email = ?", studentEmail).First(&student)
	if result.Error != nil {
		// Create student if not found
		student = models.Student{Email: studentEmail}
		createStudentError := db.Create(&student)
		if createStudentError.Error != nil {
			return nil, createStudentError.Error
		}
	}
	return &student, nil
}
