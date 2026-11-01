package db

import (
	"context"
	"log"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/repos"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// SeedProductIds are fixed so order_service can seed stock for the same products
// (see order_service/db/pg.go). Keep the two lists in sync.
var SeedProductIds = []string{
	"665000000000000000000001",
	"665000000000000000000002",
	"665000000000000000000003",
	"665000000000000000000004",
	"665000000000000000000005",
}

const seedImage = "https://cdn.hstatic.net/products/1000288298/11498_dsc05320_47c7abb602c949308c1f2bc50c3af657_master.jpg"

var seedCategories = []model.Category{
	{Id: 1, Name: "PCGAMING"},
	{Id: 2, Name: "PCDESIGN"},
	{Id: 3, Name: "PCWORK"},
}

// SeedData seeds categories and products into empty collections and returns the
// products it inserted, so the caller can index them for search.
func SeedData(ctx context.Context, database *mongo.Database) ([]model.Product, error) {
	if err := seedCategoriesData(ctx, database); err != nil {
		return nil, err
	}
	return seedProducts(ctx, database.Collection("products"))
}

func seedCategoriesData(ctx context.Context, database *mongo.Database) error {
	col := database.Collection(repos.CategoryCollection)
	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil || count > 0 {
		return err
	}

	now := time.Now()
	docs := make([]any, len(seedCategories))
	for i, c := range seedCategories {
		c.CreatedAt, c.UpdatedAt = now, now
		docs[i] = c
	}
	if _, err := col.InsertMany(ctx, docs); err != nil {
		return err
	}

	// New categories must get ids after the seeded ones.
	_, err = database.Collection(repos.CounterCollection).UpdateOne(ctx,
		bson.M{"_id": "categories"},
		bson.M{"$max": bson.M{"seq": len(seedCategories)}},
		options.UpdateOne().SetUpsert(true))
	return err
}

func seedProducts(ctx context.Context, col *mongo.Collection) ([]model.Product, error) {
	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		log.Println("Data is seeded")
		return nil, nil
	}

	gaming := model.ProductCategory{Id: 1, Name: "PCGAMING"}
	design := model.ProductCategory{Id: 2, Name: "PCDESIGN"}
	work := model.ProductCategory{Id: 3, Name: "PCWORK"}

	products := []model.Product{
		{Name: "PC BEST FOR GAMING", Desc: "i5 11400F - GTX 1660 Super 6GB", Cate: gaming, Price: 15990000},
		{Name: "PC CHƠI GAME HIỆU SUẤT CAO", Desc: "RTX 3060 12GB - 12400F", Cate: gaming, Price: 21490000},
		{Name: "PC DESIGNER - 3D RENDER - EDIT VIDEO", Desc: "i7 14700KF - RTX 5050 8GB", Cate: design, Price: 38990000},
		{Name: "PC DESIGNER EDIT VIDEO", Desc: "AMD Ryzen 9 9950X3D - RTX 5090 32GB OC", Cate: design, Price: 159990000},
		{Name: "PC HOME OFFICE", Desc: "Core i3 10105 - RAM 8GB - SSD 256GB - Kèm Màn Hình", Cate: work, Price: 8490000},
	}

	docs := make([]any, len(products))
	now := time.Now()
	for i := range products {
		id, err := bson.ObjectIDFromHex(SeedProductIds[i])
		if err != nil {
			return nil, err
		}
		products[i].Id = id
		products[i].State = model.IsAvailable
		products[i].Type = model.ProductTypeStandard
		products[i].Brand = "DDStore"
		products[i].Images = []model.Media{{Url: seedImage, Role: "thumbnail"}}
		products[i].CreatedAt = now
		docs[i] = products[i]
	}

	if _, err := col.InsertMany(ctx, docs); err != nil {
		return nil, err
	}
	return products, nil
}
