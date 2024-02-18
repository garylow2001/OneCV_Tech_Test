package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Email    string     `json:"email" gorm:"unique;not null" binding:"required"`
	Students []*Student `gorm:"many2many:teacher_student_relations;joinForeignKey:TeacherID;JoinReferences:StudentID"`
}
