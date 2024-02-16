package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Email    string     `json:"email" gorm:"primaryKey"`
	Teachers []*Teacher `gorm:"many2many:teacher_students;"` // Many-to-many relationship with teachers
}
