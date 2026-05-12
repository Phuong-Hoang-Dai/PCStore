package service

import (
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
	"gorm.io/gorm"
)

type UserDTO struct {
	Id       int    `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"column:role"`
}

type UserResponeDTO struct {
	Id        int            `json:"id" gorm:"column:id"`
	Name      string         `json:"name" gorm:"column:name"`
	Email     string         `json:"email" gorm:"column:email"`
	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `json:"delete_at" gorm:"index"`
	Role      string         `json:"roleId" gorm:"column:roleId"`
}

type UserCreateDTO struct {
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password"`
}

type UserUpdateDTO struct {
	Id    int    `json:"id" gorm:"column:id"`
	Name  string `json:"name" gorm:"column:name"`
	Email string `json:"email" gorm:"column:email"`
}

func MapUserCreateDTOtoUser(uD UserCreateDTO, u *model.User) {
	u.Name = uD.Name
	u.Email = uD.Email
	u.SetHashPassword([]byte(uD.Password))
}

func MapUserUpdateDTOtoUser(uD UserUpdateDTO, u *model.User) {
	u.Id = uD.Id
	u.Name = uD.Name
	u.Email = uD.Email
}
func MapUsertoUserDTO(u model.User, uD *UserDTO) {
	uD.Id = u.Id
	uD.Name = u.Name
	uD.Email = u.Email
	uD.Role = u.Role
}

func MapUsertoUserResponeDTO(u model.User, uD *UserResponeDTO) {
	uD.Id = u.Id
	uD.Name = u.Name
	uD.Email = u.Email
	uD.CreatedAt = u.CreatedAt
	uD.UpdatedAt = u.UpdatedAt
	uD.DeletedAt = u.DeletedAt
	uD.Role = u.Role
}
