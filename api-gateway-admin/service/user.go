package service

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/storage/managers"
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	UM managers.UserManager
}

func NewUserService(PsqlConn *sql.DB, MongoConn *mongo.Client) *UserService {
	return &UserService{UM: *managers.NewUserManager(PsqlConn, MongoConn, config.Load().MONGO_DB_NAME, config.Load().MONGO_COLLECTION_NAME)}
}

func (u *UserService) IsEmailExists(email string) error {
	return u.UM.IsEmailExists(email)
}

func (u *UserService) BanUser(req *models.BanUserReq) error {
	return u.UM.BanUser(*req)
}

func (u *UserService) UnbanUser(req *models.UnbanUserReq) error {
	return u.UM.UnbanUser(*req)
}

func (u *UserService) AddCourier(req *models.AddCourierReq) error {
	return u.UM.AddCourier(req)
}

func (u *UserService) DeleteCourier(req *models.DeleteCourierReq) error {
	return u.UM.DeleteCourier(*req)
}

func (u *UserService) AddProductManager(req *models.AddProductManagerReq) error {
	return u.UM.AddProductManager(req)
}

func (u *UserService) DeleteProductManager(req *models.DeleteProductManagerReq) error {
	return u.UM.DeleteProductManager(*req)
}
