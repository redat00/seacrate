package models

type Meta struct {
	Key   string `gorm:"primaryKey"`
	Value string
}
