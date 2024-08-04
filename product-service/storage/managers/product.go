package managers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pb "progress-service/genprotos"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductManager struct {
	Collection *mongo.Collection
	PgClient   *sql.DB
}

func NewProductManager(client *mongo.Client, dbName, collectionName string, pgClient *sql.DB) *ProductManager {
	collection := client.Database(dbName).Collection(collectionName)
	return &ProductManager{Collection: collection, PgClient: pgClient}
}

func (m *ProductManager) Create(req *pb.ProductCReq) (*pb.Void, error) {
	product := bson.M{
		"name":              req.Name,
		"category":          req.Category,
		"count":             req.Count,
		"description":       req.Description,
		"img_url":           req.ImgUrl,
		"weight":            req.Weight,
		"rating":            0,
		"seller":            req.Seller,
		"additionalDetails": req.AdditionalDetails,
	}
	_, err := m.Collection.InsertOne(context.Background(), product)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (m *ProductManager) Update(req *pb.ProductUReq) (*pb.Void, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"$set": bson.M{
			"name":               req.Name,
			"category":           req.Category,
			"description":        req.Description,
			"weight":             req.Weight,
			"seller":             req.Seller,
			"additional_details": req.AdditionalDetails,
		},
	}
	_, err = m.Collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (m *ProductManager) Delete(req *pb.ByID) (*pb.Void, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = m.Collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (m *ProductManager) Get(req *pb.ByID) (*pb.ProductGRes, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	var product pb.ProductGRes
	err = m.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (m *ProductManager) GetAll(req *pb.ProductGAReq) (*pb.ProductGARes, error) {
	filter := bson.M{}
	if req.Category != "" {
		filter["category"] = req.Category
	}
	if req.Count != "" {
		filter["count"] = req.Count
	}
	if req.Rating != "" {
		filter["rating"] = req.Rating
	}
	if req.Seller != "" {
		filter["seller"] = req.Seller
	}
	options := options.Find()
	if req.Pagination != nil {
		options.SetSkip(int64((req.Pagination.Offset - 1) * req.Pagination.Limit))
		options.SetLimit(int64(req.Pagination.Limit))
	}
	cursor, err := m.Collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var products []*pb.ProductGRes
	for cursor.Next(context.Background()) {
		var product pb.ProductGRes
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return &pb.ProductGARes{Products: products}, nil
}

func (m *ProductManager) UpdateRating(req *pb.ProductRatingUReq) (*pb.Void, error) {
	id, err := primitive.ObjectIDFromHex(req.ProductId)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"$set": bson.M{"rating": req.Rate},
	}
	_, err = m.Collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (m *ProductManager) UpdateCount(req *pb.ProductCountUReq) (*pb.Void, error) {
	id, err := primitive.ObjectIDFromHex(req.ProductId)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"$set": bson.M{"count": req.Count},
	}
	_, err = m.Collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (m *ProductManager) UpdateImg(req *pb.ProductImageUReq) (*pb.Void, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"$set": bson.M{"img_url": req.ImgUrl},
	}
	_, err = m.Collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
