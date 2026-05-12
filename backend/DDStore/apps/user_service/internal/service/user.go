package service

import "github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"

func CreateUser(uD UserCreateDTO, repos UserRepos) (int, error) {
	var u model.User
	MapUserCreateDTOtoUser(uD, &u)
	u.Role = "user"

	if id, err := repos.CreateUser(&u); err != nil {
		return 0, nil
	} else {
		return id, err
	}
}

func GetUserById(id int, repos UserRepos) (u UserResponeDTO, err error) {
	var userDAO model.User

	if err := repos.GetUserById(id, &userDAO); err != nil {
		return u, err
	} else {
		MapUsertoUserResponeDTO(userDAO, &u)
		return u, nil
	}
}

func UpdateUser(uD UserUpdateDTO, repos UserRepos) error {
	var u model.User
	MapUserUpdateDTOtoUser(uD, &u)

	if err := repos.UpdateUser(u); err != nil {
		return nil
	} else {
		return err
	}
}

func DeleteUser(id int, repos UserRepos) error {
	if err := repos.DeleteUser(id); err != nil {
		return err
	} else {
		return nil
	}
}

func GetUsers(p *model.Paging, repos UserRepos) ([]UserResponeDTO, error) {
	p.Process()
	var userDAO []model.User
	if err := repos.GetUsers(*p, &userDAO); err != nil {
		return nil, err
	}
	data := make([]UserResponeDTO, len(userDAO), cap(userDAO))
	for i := range userDAO {
		MapUsertoUserResponeDTO(userDAO[i], &data[i])
	}

	return data, nil
}

func VerifyPassword(data UserDTO, repos UserRepos) (UserDTO, error) {
	userDao := model.User{}
	if err := repos.GetUserByName(data.Name, &userDao); err != nil {
		return UserDTO{}, err
	}

	if err := userDao.VerifyPassword([]byte(data.Password)); err != nil {
		return UserDTO{}, model.ErrUserNameOrPasswordIncorrect
	}
	dataResp := UserDTO{}
	MapUsertoUserDTO(userDao, &dataResp)

	return dataResp, nil
}
