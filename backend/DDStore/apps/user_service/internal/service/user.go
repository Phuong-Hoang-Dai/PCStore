package service

import (
	"context"
	"net/mail"
	"strings"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
)

// CreateUser handles both ways an account comes to exist:
//   - with a password, it registers; if someone already ordered as a guest with this email,
//     that guest account is claimed (it gets the name and password and keeps its orders)
//   - without one, it is a guest placing an order: the account already using the email is
//     returned, or one is created without a password, which can't log in until registered
func CreateUser(uD UserCreateDTO, repos UserRepos, ctx context.Context) (int, error) {
	uD.Email = strings.ToLower(strings.TrimSpace(uD.Email))
	isGuest := uD.Password == ""
	if isGuest {
		uD.Name = strings.TrimSpace(uD.Name)
	}
	if strings.TrimSpace(uD.Name) == "" {
		return 0, model.ErrNameRequired
	}
	if (isGuest || uD.Email != "") && !isValidEmail(uD.Email) {
		return 0, model.ErrInvalidEmail
	}

	var u model.User
	MapUserCreateDTOtoUser(uD, &u)
	u.Role = "user"

	if isGuest {
		return repos.FindOrCreateGuest(&u, ctx)
	}
	if id, err := repos.CreateUser(&u, ctx); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

// isValidEmail accepts a bare address such as a@b.c, without a display name.
func isValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
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
