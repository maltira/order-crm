package mongo

import (
	"context"
	"order-crm/pkg/database"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	ProductID int                    `json:"product_id" bson:"product_id"`
	Name      string                 `json:"name" bson:"name"`
	Price     float64                `json:"price" bson:"price"`
	Meta      map[string]interface{} `json:"meta" bson:"meta"`
}

func GetAllProducts() ([]Product, error) {
	ctx := context.TODO()
	cursor, err := database.MDB.Collection("products").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func GetProductById(productID int) (*Product, error) {
	var product Product

	filter := bson.M{"product_id": productID}

	err := database.MDB.Collection("products").FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
