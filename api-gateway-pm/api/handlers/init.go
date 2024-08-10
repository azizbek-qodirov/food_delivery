package handlers

import (
	pb "gateway-admin/genprotos"

	"gateway-admin/drivers"

	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
)

type HTTPHandler struct {
	ProductManager pb.ProductServiceClient
	MinIO          *minio.Client
}

func NewHandler(connP *grpc.ClientConn) *HTTPHandler {
	minioClient, err := drivers.MinIOConnect()
	if err != nil {
		panic(err)
	}
	return &HTTPHandler{
		ProductManager: pb.NewProductServiceClient(connP),
		MinIO:          minioClient,
	}
}
