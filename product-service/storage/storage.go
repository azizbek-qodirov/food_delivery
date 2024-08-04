package storage

import (
	pb "progress-service/genprotos"
)

type StorageI interface {
	Product() ProductI
}

type ProductI interface {
	Create(*pb.ProductCReq) (*pb.Void, error)
	Update(*pb.ProductUReq) (*pb.Void, error)
	Delete(*pb.ByID) (*pb.Void, error)
	Get(*pb.ByID) (*pb.ProductGRes, error)
	GetAll(*pb.ProductGAReq) (*pb.ProductGARes, error)
	UpdateRating(*pb.ProductRatingUReq) (*pb.Void, error)
	UpdateCount(*pb.ProductCountUReq) (*pb.Void, error)
}
