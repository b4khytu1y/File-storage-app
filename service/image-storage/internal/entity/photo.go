package entity

type Photo struct {
	ID     uint `gorm:"primary_key"`
	UserID int
	Name   string
	Data   []byte
}
