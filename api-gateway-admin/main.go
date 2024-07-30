package main

import (
	"auth-service/api"
	"auth-service/api/handlers"
	"auth-service/config"
	"auth-service/service"
	"auth-service/storage"
)

func main() {
	cf := config.Load()
	em := config.NewErrorManager()

	pgsql, mongo, err := storage.ConnectDB(&cf)
	em.CheckErr(err)
	defer pgsql.Close()

	us := service.NewUserService(pgsql, mongo)
	handler := handlers.NewHandler(us)

	roter := api.NewRouter(handler)
	if err := roter.Run(cf.API_GATEWAY_ADMIN_PORT); err != nil {
		panic(err)
	}
}
