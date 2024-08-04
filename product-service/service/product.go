package service

import (
	"context"
	pb "progress-service/genprotos"
	"progress-service/storage"
)

type ProductService struct {
	storage storage.StorageI
	pb.UnimplementedProductServiceServer
}

func NewProductService(storage storage.StorageI) *ProductService {
	return &ProductService{storage: storage}
}

func (s *ProductService) Create(ctx context.Context, req *pb.ProductCReq) (*pb.Void, error) {
	return s.storage.Product().Create(req)
}

func (s *ProductService) Update(ctx context.Context, req *pb.ProductUReq) (*pb.Void, error) {
	return s.storage.Product().Update(req)
}

func (s *ProductService) Delete(ctx context.Context, req *pb.ByID) (*pb.Void, error) {
	return s.storage.Product().Delete(req)
}

func (s *ProductService) GetAll(ctx context.Context, req *pb.ProductGAReq) (*pb.ProductGARes, error) {
	return s.storage.Product().GetAll(req)
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.ByID) (*pb.ProductGRes, error) {
	return s.storage.Product().Get(req)
}

func (s *ProductService) UpdateRating(ctx context.Context, req *pb.ProductRatingUReq) (*pb.Void, error) {
	return s.storage.Product().UpdateRating(req)
}

func (s *ProductService) UpdateCount(ctx context.Context, req *pb.ProductCountUReq) (*pb.Void, error) {
	return s.storage.Product().UpdateCount(req)
}
