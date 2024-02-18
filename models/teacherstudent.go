package models

type TeacherStudentRelation struct {
	TeacherID uint `gorm:"primaryKey"`
	StudentID uint `gorm:"primaryKey"`
	Teacher   Teacher
	Student   Student
}
