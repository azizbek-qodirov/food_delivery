package handlers

import (
	pb "gateway-admin/genprotos"

	"google.golang.org/grpc"
)

type HTTPHandler struct {
	ProductManager pb.ProductServiceClient
}

func NewHandler(connP *grpc.ClientConn) *HTTPHandler {
	return &HTTPHandler{
		ProductManager: pb.NewProductServiceClient(connP),
	}
}
