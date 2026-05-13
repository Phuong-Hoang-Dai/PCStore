package service

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
)

func CreateProduct(data model.Product, repos ProductRepos) (int, error) {
	if id, err := repos.CreateProduct(data); err != nil {
		return 0, err
	} else {
		return id, err
	}
}

func UpdateProduct(data model.Product, repos ProductRepos) error {
	if err := repos.UpdateProduct(data); err != nil {
		return err
	} else {
		return nil
	}
}

func GetProducts(p *model.Paging, repos ProductRepos) (data []model.Product, err error) {
	p.Process()
	if data, err = repos.GetProducts(*p); err != nil {
		return nil, err
	}
	return data, nil
}

func GetProductsByCate(p *model.Paging, repos ProductRepos, cate model.Category) (data []model.Product, err error) {
	p.Process()
	if data, err = repos.GetProductsByCate(*p, cate); err != nil {
		return nil, err
	}
	return data, nil
}

func GetProductById(id int, repos ProductRepos) (data model.Product, err error) {
	if data, err := repos.GetProductById(id); err != nil {
		return data, err
	} else {
		return data, nil
	}
}

func DeleteProduct(id int, repos ProductRepos) error {
	if err := repos.DeleteProduct(id); err != nil {
		return err
	} else {
		return nil
	}
}

func GetStock(data []OrderItemsDto, repos ProductRepos) error {
	stock, err := getProducts(data, repos)
	if err != nil {
		return err
	}

	for i := range data {
		if canGetStock := data[i].Quantity <= stock[i].Quantity; canGetStock {
			stock[i].Quantity -= data[i].Quantity
		} else {
			return model.ErrOutOfStock
		}
	}

	if err := repos.UpdateProducts(stock); err != nil {
		return err
	}

	return nil
}

func RestoreStock(data []OrderItemsDto, repos ProductRepos) error {
	stock, err := getProducts(data, repos)
	if err != nil {
		return err
	}

	for i := range data {
		stock[i].Quantity += (data)[i].Quantity
	}

	if err := repos.UpdateProducts(stock); err != nil {
		return err
	}

	return nil
}

func GetPriceProduct(data *[]OrderItemsDto, repos ProductRepos) (err error) {
	stock, err := getProducts(*data, repos)
	if err != nil {
		return err
	}

	for i := range *data {
		(*data)[i].Price = stock[i].Price
	}

	return nil
}

func getProducts(items []OrderItemsDto, repos ProductRepos) (stock []model.Product, err error) {
	stock = make([]model.Product, len(items), cap(items))

	for i := range items {
		if stock[i], err = repos.GetProductById(items[i].ProductId); err != nil {
			return nil, err
		}
	}

	return stock, err
}
