package models

type Users struct {
	Id       uint   `gorm:"primary_key;column:id" json:"id"`
	Email    string `gorm:"column:email" json:"email"`
	Fullname string `gorm:"column:fullname" json:"fullname"`
}
type APImessage struct {
	Message string
}
