package service

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/storage/managers"
	"database/sql"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	UM managers.UserManager
}

func NewUserService(PsqlConn *sql.DB, MongoConn *mongo.Client) *UserService {
	return &UserService{UM: *managers.NewUserManager(PsqlConn, MongoConn, config.Load().MONGO_DB_NAME, config.Load().MONGO_COLLECTION_NAME)}
}

func (u *UserService) Register(req *models.RegisterReq) error {
	req.ID = uuid.NewString()
	req.Role = "user"
	// err := u.UM.RegisterInMongo(req.ID)
	// if err != nil {
	// 	return err
	// }
	if err := u.UM.Register(*req); err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetProfile(req *models.GetProfileReq) (*models.GetProfileResp, error) {
	return u.UM.Profile(*req)
}

func (u *UserService) EmailExists(email string) (bool, error) {
	return u.UM.EmailExists(email)
}

func (u *UserService) GetByID(id *models.GetProfileByIdReq) (*models.GetProfileByIdResp, error) {
	return u.UM.GetByID(id)
}

func (u *UserService) ChangeRole(req *models.ChangeRoleReq) error {
	return u.UM.ChangeRole(*req)
}
