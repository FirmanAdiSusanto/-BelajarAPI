package data

import "clean1/features/todo/data"

type User struct {
	Nama     string
	Hp       string `gorm:"type:varchar(13);primaryKey"`
	Password string
	Todos    []data.Todo `gorm:"foreignKey:Pemilik;references:Hp"`
}
