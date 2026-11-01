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

const (
	CategoryCollection = "categories"
	CounterCollection  = "counters"
	// categoryCounterId is the counters document that hands out category ids.
	categoryCounterId = "categories"
)

type mongoCateRepo struct {
	col      *mongo.Collection
	counters *mongo.Collection
}

func NewMongoCateRepo(db *mongo.Database) service.CateRepos {
	if db == nil {
		panic("mongo database must not be nil")
	}
	return mongoCateRepo{col: db.Collection(CategoryCollection), counters: db.Collection(CounterCollection)}
}

func (m mongoCateRepo) nextId(ctx context.Context) (int, error) {
	var counter struct {
		Seq int `bson:"seq"`
	}
	err := m.counters.FindOneAndUpdate(ctx,
		bson.M{"_id": categoryCounterId},
		bson.M{"$inc": bson.M{"seq": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&counter)
	return counter.Seq, err
}

func (m mongoCateRepo) CreateCate(ctx context.Context, data model.Category) (int, error) {
	id, err := m.nextId(ctx)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	data.Id = id
	data.CreatedAt = now
	data.UpdatedAt = now
	data.DeletedAt = nil
	if _, err := m.col.InsertOne(ctx, data); err != nil {
		return 0, err
	}
	return id, nil
}

func (m mongoCateRepo) GetCateById(ctx context.Context, id int) (model.Category, error) {
	var data model.Category
	err := m.col.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&data)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Category{}, model.ErrNotFound
	}
	return data, err
}

func (m mongoCateRepo) GetCates(ctx context.Context) ([]model.Category, error) {
	cursor, err := m.col.Find(ctx, bson.M{"deleted_at": nil}, options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	data := []model.Category{}
	err = cursor.All(ctx, &data)
	return data, err
}

func (m mongoCateRepo) UpdateCate(ctx context.Context, data model.Category) error {
	result, err := m.col.UpdateOne(ctx,
		bson.M{"_id": data.Id, "deleted_at": nil},
		bson.M{"$set": bson.M{"name": data.Name, "updated_at": time.Now()}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (m mongoCateRepo) DeleteCate(ctx context.Context, id int) error {
	result, err := m.col.UpdateOne(ctx,
		bson.M{"_id": id, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return model.ErrNotFound
	}
	return nil
}
