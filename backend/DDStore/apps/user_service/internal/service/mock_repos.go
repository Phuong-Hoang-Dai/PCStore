package service

import (
	"context"
	"errors"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
)

type MockRepos struct{}

var userList []model.User

func (MockRepos) Init() {
	userList = []model.User{}
}

func (MockRepos) CreateUser(data *model.User, ctx context.Context) (int, error) {
	userList = append(userList, *data)
	userList[len(userList)-1].Id = len(userList)

	return userList[len(userList)-1].Id, nil
}

func (m MockRepos) FindOrCreateGuest(data *model.User, ctx context.Context) (int, error) {
	for _, v := range userList {
		if v.Email == data.Email {
			return v.Id, nil
		}
	}
	return m.CreateUser(data, ctx)
}

func (MockRepos) UpdateUser(data model.User, ctx context.Context) error {
	userList[data.Id-1] = data

	return nil
}

func (MockRepos) GetUserById(id int, data *model.User, ctx context.Context) error {
	*data = userList[id-1]
	return nil
}

func (MockRepos) GetUserByName(name string, data *model.User, ctx context.Context) error {
	for _, v := range userList {
		if v.Name == name && v.HashedPassword != nil {
			*data = v
			return nil
		}
	}
	return errors.New("can't find user")
}

func (MockRepos) GetUsers(p model.Paging, data *[]model.User, ctx context.Context) error {
	for i := p.Offset; i < len(userList); i++ {
		if i-p.Offset < p.Limit {
			*data = append(*data, userList[i])
		}
	}

	return nil
}

func (MockRepos) DeleteUser(id int, ctx context.Context) error {
	return nil
}
