package repos_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/repos"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mockGormDB *gorm.DB
var mockDB sqlmock.Sqlmock

func TestMain(m *testing.M) {
	db, mock, gormDB, err := setupMockDB()
	if err != nil {
		log.Fatalf("There is something wrong. Erorr: %v", err)
	}
	defer db.Close()

	mockGormDB = gormDB
	mockDB = mock

	m.Run()

}

func TestCreateProduct(t *testing.T) {
	data := model.Product{
		Name:     "Haha",
		Desc:     "hahahaha",
		Price:    99,
		Quantity: 6,
	}

	mockCreate(mockDB)

	t.Run("Test Create", func(t *testing.T) {
		product_Repos := repos.NewMysqlProductRepo(mockGormDB)
		_, err := product_Repos.CreateProduct(data)

		assert.NoError(t, err)
	})
}

func TestUpdateProduct(t *testing.T) {
	data := model.Product{
		Id:       1,
		Name:     "Haha",
		Desc:     "hahahaha",
		Price:    99,
		Quantity: 6,
	}

	mockUpdate(mockDB)

	t.Run("Test Update", func(t *testing.T) {
		product_Repos := repos.NewMysqlProductRepo(mockGormDB)
		err := product_Repos.UpdateProduct(data)
		assert.NoError(t, err)
	})
}

func TestUpdateProductNoRecord(t *testing.T) {
	data := model.Product{
		Id:       3,
		Name:     "Haha",
		Desc:     "hahahaha",
		Price:    99,
		Quantity: 6,
	}

	mockDB.ExpectBegin()
	mockDB.ExpectExec("UPDATE `products` SET .* WHERE .*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 3).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.ExpectCommit()

	t.Run("Test Update", func(t *testing.T) {
		product_Repos := repos.NewMysqlProductRepo(mockGormDB)
		err := product_Repos.UpdateProduct(data)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestDeleteProduct(t *testing.T) {
	data := model.Product{
		Id:       3,
		Name:     "Haha",
		Desc:     "hahahaha",
		Price:    99,
		Quantity: 6,
	}

	mockDelete(mockDB)

	t.Run("Test Delete", func(t *testing.T) {
		product_Repos := repos.NewMysqlProductRepo(mockGormDB)
		err := product_Repos.DeleteProduct(data.Id)
		assert.NoError(t, err)
	})
}

func TestUpdateProducts(t *testing.T) {
	data := []model.Product{
		{
			Id:       1,
			Name:     "Haha",
			Desc:     "hahahaha",
			Price:    99,
			Quantity: 6,
		},
		{
			Id:       2,
			Name:     "Haha",
			Desc:     "hahahaha",
			Price:    9,
			Quantity: 6,
		},
	}
	mockDB.ExpectBegin()
	for i := range data {
		mockDB.ExpectExec("UPDATE `products` SET .* WHERE .*").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), i+1).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mockDB.ExpectCommit()

	t.Run("Test UpdateProducts", func(t *testing.T) {
		product_Repos := repos.NewMysqlProductRepo(mockGormDB)
		err := product_Repos.UpdateProducts(data)
		assert.NoError(t, err)
	})
}

func TestGetProductById(t *testing.T) {
	data := model.Product{
		Id:       3,
		Name:     "Haha",
		Desc:     "hahahaha",
		Price:    99,
		Quantity: 6,
	}
	p := model.Product{
		Id: 3,
	}

	mockSelectProduct(mockDB, data)

	t.Run("Test GetProductById", func(t *testing.T) {
		product_Repos := repos.NewMysqlProductRepo(mockGormDB)
		p, err := product_Repos.GetProductById(p.Id)
		assert.NoError(t, err)
		assert.Equal(t, data.Name, p.Name)
		assert.Equal(t, data.Desc, p.Desc)
		assert.Equal(t, data.Price, p.Price)
		assert.Equal(t, data.Quantity, p.Quantity)
	})
}

func TestGetProducts(t *testing.T) {
	data := []model.Product{
		{
			Id:       3,
			Name:     "Milk",
			Desc:     "This is milk",
			Price:    99,
			Quantity: 68,
		},
		{
			Id:       9,
			Name:     "Milo",
			Desc:     "This is milo",
			Price:    55,
			Quantity: 6,
		},
	}

	p := model.Paging{
		Limit:  2,
		Offset: 1,
	}

	mockSelectProducts(mockDB, data, p)

	t.Run("Test GetProductById", func(t *testing.T) {
		product_Repos := repos.NewMysqlProductRepo(mockGormDB)
		p.Process()
		pList, err := product_Repos.GetProducts(p)
		assert.NoError(t, err)
		for i := range data {
			assert.Equal(t, data[i].Name, pList[i].Name)
			assert.Equal(t, data[i].Desc, pList[i].Desc)
			assert.Equal(t, data[i].Price, pList[i].Price)
			assert.Equal(t, data[i].Quantity, pList[i].Quantity)
		}
	})
}

func setupMockDB() (*sql.DB, sqlmock.Sqlmock, *gorm.DB, error) {
	db, mock, err := sqlmock.New()
	gormDB, err2 := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err2 != nil && err == nil {
		err = err2
	}

	return db, mock, gormDB, err
}

func mockCreate(m sqlmock.Sqlmock) {
	m.ExpectBegin()
	m.ExpectExec("INSERT INTO `products`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	m.ExpectCommit()
}

func mockUpdate(m sqlmock.Sqlmock) {
	m.ExpectBegin()
	m.ExpectExec("UPDATE `products` SET .* WHERE .*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	m.ExpectCommit()
}

func mockDelete(m sqlmock.Sqlmock) {
	m.ExpectBegin()
	m.ExpectExec("UPDATE `products` SET .*`deleted_at`=?.* WHERE .*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	m.ExpectCommit()
}

func mockSelectProduct(m sqlmock.Sqlmock, data model.Product) {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "quantity", "create_at", "update_at", "delete_at"}).
		AddRow(data.Id, data.Name, data.Desc, data.Price, data.Quantity, data.CreatedAt, data.UpdatedAt, data.DeletedAt)

	m.ExpectQuery(`SELECT .* FROM .*products.* WHERE .*id.*`).
		WithArgs(data.Id, 1).
		WillReturnRows(rows)
}

func mockSelectProducts(m sqlmock.Sqlmock, data []model.Product, paging model.Paging) {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "quantity", "create_at", "update_at", "delete_at"})
	for _, v := range data {
		rows.AddRow(v.Id, v.Name, v.Desc, v.Price, v.Quantity, v.CreatedAt, v.UpdatedAt, v.DeletedAt)

	}

	m.ExpectQuery("SELECT .* FROM `products` WHERE .* LIMIT .* OFFSET .*").
		WithArgs(paging.Limit, paging.Offset).
		WillReturnRows(rows)
}
