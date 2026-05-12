package service

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
)

func CreateCate(data model.Category, repos CateRepos) (int, error) {
	if id, err := repos.CreateCate(data); err != nil {
		return 0, err
	} else {
		return id, err
	}
}

func UpdateCate(data model.Category, repos CateRepos) error {
	if err := repos.UpdateCate(data); err != nil {
		return err
	} else {
		return nil
	}
}

func GetCates(repos CateRepos) (data []model.Category, err error) {
	if data, err = repos.GetCates(); err != nil {
		return nil, err
	}
	return data, nil
}

func GetCateById(id int, repos CateRepos) (data model.Category, err error) {
	if data, err := repos.GetCateById(id); err != nil {
		return data, err
	} else {
		return data, nil
	}
}

func DeleteCate(id int, repos CateRepos) error {
	if err := repos.DeleteCate(id); err != nil {
		return err
	} else {
		return nil
	}
}
