package repos

import (
	"context"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mognoProductRepo struct {
	col *mongo.Collection
}

func NewMongoProductRepo(col *mongo.Collection) service.ProductRepos {
	if col == nil {
		panic("mongo collection must not be nil")
	}
	return mognoProductRepo{col: col}
}

var notDeleted = bson.M{"deletedat.valid": bson.M{"$ne": true}}

func (m mognoProductRepo) CreateProduct(data model.Product) (int, error) {
	_, err := m.col.InsertOne(context.Background(), data)
	return data.Id, err
}

func (m mognoProductRepo) GetProductById(id int) (model.Product, error) {
	var data model.Product
	filter := bson.M{"id": id, "deletedat.valid": bson.M{"$ne": true}}
	err := m.col.FindOne(context.Background(), filter).Decode(&data)
	return data, err
}

func (m mognoProductRepo) UpdateProduct(data model.Product) error {
	filter := bson.M{"id": data.Id, "deletedat.valid": bson.M{"$ne": true}}
	result, err := m.col.UpdateOne(context.Background(), filter, bson.M{"$set": data})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (m mognoProductRepo) UpdateProducts(data []model.Product) error {
	session, err := m.col.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(ctx context.Context) (any, error) {
		for i := range data {
			filter := bson.M{"id": data[i].Id, "deletedat.valid": bson.M{"$ne": true}}
			result, err := m.col.UpdateOne(ctx, filter, bson.M{"$set": data[i]})
			if err != nil {
				return nil, err
			}
			if result.MatchedCount == 0 {
				return nil, mongo.ErrNoDocuments
			}
		}
		return nil, nil
	})
	return err
}

func (m mognoProductRepo) DeleteProduct(id int) error {
	filter := bson.M{"id": id, "deletedat.valid": bson.M{"$ne": true}}
	softDelete := bson.M{"$set": bson.M{"deletedat": bson.M{"time": time.Now(), "valid": true}}}
	result, err := m.col.UpdateOne(context.Background(), filter, softDelete)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (m mognoProductRepo) GetProducts(p model.Paging) ([]model.Product, error) {
	opts := options.Find().SetSkip(int64(p.Offset)).SetLimit(int64(p.Limit))
	cursor, err := m.col.Find(context.Background(), notDeleted, opts)
	if err != nil {
		return nil, err
	}
	var data []model.Product
	err = cursor.All(context.Background(), &data)
	return data, err
}

func (m mognoProductRepo) GetProductsByCate(p model.Paging, cate model.Category) ([]model.Product, error) {
	filter := bson.M{"cate.id": cate.Id, "deletedat.valid": bson.M{"$ne": true}}
	opts := options.Find().SetSkip(int64(p.Offset)).SetLimit(int64(p.Limit))
	cursor, err := m.col.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	var data []model.Product
	err = cursor.All(context.Background(), &data)
	return data, err
}
