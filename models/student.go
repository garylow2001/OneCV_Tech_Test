package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Email     string `json:"email"`
	TeacherID uint
}
