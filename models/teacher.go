package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Email    string     `json:"email" gorm:"primaryKey"`
	Students []*Student `gorm:"many2many:teacher_students;"` // Many-to-many relationship with students
}
