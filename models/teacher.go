package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Email    string    `json:"email"`
	Students []Student `json:"students" gorm:"foreignKey:TeacherID"`
}
