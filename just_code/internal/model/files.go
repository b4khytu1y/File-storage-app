package model

type FileModel struct {
	Id          int    `gorm:"type:int;primary_key"`
	UserID      int    `gorm:"index"`
	User        Users  `gorm:"foreignkey:UserID"`
	Name        string `gorm:"type:varchar(255)"`
	ContentType string `gorm:"type:varchar(100)"`
	Size        int64  `gorm:"type:bigint"`
	Content     []byte `gorm:"type:bytea"`
}
