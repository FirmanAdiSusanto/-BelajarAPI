package data

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Kegiatan string
	Pemilik  string `gorm:"type:varchar(13);"`
}
