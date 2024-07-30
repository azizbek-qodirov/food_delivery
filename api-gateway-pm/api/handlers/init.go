package handlers

import (
	pb "gateway-admin/genprotos"

	"google.golang.org/grpc"
)

type HTTPHandler struct {
	ProductManager pb.ProductManagerServiceClient
}

func NewHandler(connP *grpc.ClientConn) *HTTPHandler {
	return &HTTPHandler{
		ProductManager: pb.NewProductManagerServiceClient(connP),
	}
}
