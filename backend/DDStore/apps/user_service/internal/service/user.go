package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
)

func CreateUser(uD UserCreateDTO, repos UserRepos, ctx context.Context) (int, error) {
	var u model.User
	MapUserCreateDTOtoUser(uD, &u)
	u.Role = "user"

	if id, err := repos.CreateUser(&u, ctx); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func GetUserById(id int, repos UserRepos, ctx context.Context) (u UserResponeDTO, err error) {
	var userDAO model.User

	if err := repos.GetUserById(id, &userDAO, ctx); err != nil {
		return u, err
	} else {
		MapUsertoUserResponeDTO(userDAO, &u)
		return u, nil
	}
}

func UpdateUser(uD UserUpdateDTO, repos UserRepos, ctx context.Context) error {
	var u model.User
	MapUserUpdateDTOtoUser(uD, &u)

	if err := repos.UpdateUser(u, ctx); err != nil {
		return err
	} else {
		return nil
	}
}

func DeleteUser(id int, repos UserRepos, ctx context.Context) error {
	if err := repos.DeleteUser(id, ctx); err != nil {
		return err
	} else {
		return nil
	}
}

func GetUsers(p *model.Paging, repos UserRepos, ctx context.Context) ([]UserResponeDTO, error) {
	p.Process()
	var userDAO []model.User
	if err := repos.GetUsers(*p, &userDAO, ctx); err != nil {
		return nil, err
	}
	data := make([]UserResponeDTO, len(userDAO), cap(userDAO))
	for i := range userDAO {
		MapUsertoUserResponeDTO(userDAO[i], &data[i])
	}

	return data, nil
}

func VerifyPassword(data UserDTO, repos UserRepos, ctx context.Context) (UserDTO, error) {
	userDao := model.User{}
	if err := repos.GetUserByName(data.Name, &userDao, ctx); err != nil {
		return UserDTO{}, err
	}

	if err := userDao.VerifyPassword([]byte(data.Password)); err != nil {
		return UserDTO{}, model.ErrUserNameOrPasswordIncorrect
	}
	dataResp := UserDTO{}
	MapUsertoUserDTO(userDao, &dataResp)

	return dataResp, nil
}
