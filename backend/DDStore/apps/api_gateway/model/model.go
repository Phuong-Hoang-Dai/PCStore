package model

type User struct {
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"column:role"`
}
type UseResp struct {
	Id       int    `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"column:role"`
}
