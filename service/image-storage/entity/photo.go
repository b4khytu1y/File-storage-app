package entity

type Photo struct {
	ID   uint `gorm:"primary_key"`
	Name string
	Data []byte
}
