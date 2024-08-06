package main

import (
	"fmt"
	"gateway-admin/api"
	"gateway-admin/api/handlers"
	"gateway-admin/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cf := config.Load()
	em := config.NewErrorManager()

	ProductConn, err := grpc.NewClient(fmt.Sprintf("localhost%s", cf.PRODUCT_SERVICE_PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
	em.CheckErr(err)
	defer ProductConn.Close()

	handler := handlers.NewHandler(ProductConn)

	roter := api.NewRouter(handler)
	if err := roter.Run(cf.API_GATEWAY_PM_PORT); err != nil {
		panic(err)
	}
}
