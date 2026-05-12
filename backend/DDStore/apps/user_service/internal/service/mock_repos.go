package service

import (
	"errors"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
)

type MockRepos struct{}

var userList []model.User

func (MockRepos) Init() {
	userList = []model.User{}
}

func (MockRepos) CreateUser(data *model.User) (int, error) {
	userList = append(userList, *data)
	userList[len(userList)-1].Id = len(userList)

	return userList[len(userList)-1].Id, nil
}

func (MockRepos) UpdateUser(data model.User) error {
	userList[data.Id-1] = data

	return nil
}

func (MockRepos) GetUserById(id int, data *model.User) error {
	*data = userList[id-1]
	return nil
}

func (MockRepos) GetUserByName(name string, data *model.User) error {
	for _, v := range userList {
		if v.Name == name {
			*data = v
			return nil
		}
	}
	return errors.New("can't find user")
}

func (MockRepos) GetUsers(p model.Paging, data *[]model.User) error {
	for i := p.Offset; i < len(userList); i++ {
		if i-p.Offset < p.Limit {
			*data = append(*data, userList[i])
		}
	}

	return nil
}

func (MockRepos) DeleteUser(id int) error {
	return nil
}
