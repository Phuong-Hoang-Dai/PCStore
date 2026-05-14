package repos

import (
	"context"
	"errors"
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

var notDeleted = bson.M{"deleted_at": nil}

func (m mognoProductRepo) CreateProduct(ctx context.Context, data model.Product) (bson.ObjectID, error) {
	data.Id = bson.NewObjectID()
	data.CreatedAt = time.Now()
	_, err := m.col.InsertOne(ctx, data)
	return data.Id, err
}

func (m mognoProductRepo) GetProductById(ctx context.Context, id bson.ObjectID) (model.Product, error) {
	var data model.Product
	filter := bson.M{"_id": id, "deleted_at": nil}
	err := m.col.FindOne(ctx, filter).Decode(&data)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Product{}, model.ErrNotFound
	}
	return data, err
}

func (m mognoProductRepo) UpdateProduct(ctx context.Context, data model.Product) error {
	filter := bson.M{"_id": data.Id, "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"name":        data.Name,
		"description": data.Desc,
		"brand":       data.Brand,
		"cate":        data.Cate,
		"type":        data.Type,
		"images":      data.Images,
		"updated_at":  time.Now(),
	}}
	result, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (m mognoProductRepo) UpdateProducts(ctx context.Context, data []model.Product) error {
	session, err := m.col.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		for i := range data {
			filter := bson.M{"_id": data[i].Id, "deleted_at": nil}
			update := bson.M{"$set": bson.M{
				"name":        data[i].Name,
				"description": data[i].Desc,
				"brand":       data[i].Brand,
				"cate":        data[i].Cate,
				"type":        data[i].Type,
				"images":      data[i].Images,
				"updated_at":  time.Now(),
			}}
			result, err := m.col.UpdateOne(ctx, filter, update)
			if err != nil {
				return nil, err
			}
			if result.MatchedCount == 0 {
				return nil, model.ErrNotFound
			}
		}
		return nil, nil
	})
	return err
}

func (m mognoProductRepo) DeleteProduct(ctx context.Context, id bson.ObjectID) error {
	filter := bson.M{"_id": id, "deleted_at": nil}
	now := time.Now()
	update := bson.M{"$set": bson.M{"deleted_at": now}}
	result, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (m mognoProductRepo) GetProducts(ctx context.Context, p model.Paging) ([]model.Product, error) {
	opts := options.Find().SetSkip(int64(p.Offset)).SetLimit(int64(p.Limit))
	cursor, err := m.col.Find(ctx, notDeleted, opts)
	if err != nil {
		return nil, err
	}
	var data []model.Product
	err = cursor.All(ctx, &data)
	return data, err
}

func (m mognoProductRepo) GetProductsByCate(ctx context.Context, p model.Paging, cate model.Category) ([]model.Product, error) {
	filter := bson.M{"cate.id": cate.Id, "deleted_at": nil}
	opts := options.Find().SetSkip(int64(p.Offset)).SetLimit(int64(p.Limit))
	cursor, err := m.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var data []model.Product
	err = cursor.All(ctx, &data)
	return data, err
}
