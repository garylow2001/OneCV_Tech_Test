package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Email     string     `json:"email" gorm:"unique;not null" binding:"required"`
	Suspended bool       `json:"suspended" gorm:"default:false"`
	Teachers  []*Teacher `gorm:"many2many:teacher_student_relations;joinForeignKey:StudentID;JoinReferences:TeacherID"`
}
