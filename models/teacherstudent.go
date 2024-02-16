package models

import "gorm.io/gorm"

type TeacherStudent struct {
	gorm.Model
	TeacherEmail string `gorm:"primaryKey"`
	StudentEmail string `gorm:"primaryKey"`
}
