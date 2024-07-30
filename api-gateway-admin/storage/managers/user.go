package managers

import (
	"auth-service/config"
	"auth-service/models"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserManager struct {
	PgClient    *sql.DB
	MongoClient *mongo.Collection
}

func NewUserManager(db *sql.DB, client *mongo.Client, dbName, collectionName string) *UserManager {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserManager{PgClient: db, MongoClient: collection}
}

func (m *UserManager) IsEmailExists(email string) error {
	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	var count int
	err := m.PgClient.QueryRow(query, email).Scan(&count)
	if err != nil {
		return errors.New("server error: " + err.Error())
	}
	if count > 0 {
		return errors.New("email already registered: " + email)
	}
	return nil
}

func (m *UserManager) BanUser(req models.BanUserReq) error {
	var query string
	if req.ID != "" {
		query = "UPDATE users SET role = 'banned' WHERE id = $1 and role = 'user'"
		res, err := m.PgClient.Exec(query, req.ID)
		if err != nil {
			return err
		}
		if err := config.CheckRowsAffected(res, "user"); err != nil {
			return err
		}
	} else if req.Email != "" {
		query = "UPDATE users SET role = 'banned' WHERE email = $1 and role = 'user'"
		res, err := m.PgClient.Exec(query, req.Email)
		if err != nil {
			return err
		}
		if err := config.CheckRowsAffected(res, "user"); err != nil {
			return err
		}
	}
	return nil
}

func (m *UserManager) UnbanUser(req models.UnbanUserReq) error {
	var query string
	if req.ID != "" {
		query = "UPDATE users SET role = 'user' WHERE id = $1 and role = 'banned'"
		res, err := m.PgClient.Exec(query, req.ID)
		if err != nil {
			return err
		}
		if err := config.CheckRowsAffected(res, "user"); err != nil {
			return err
		}
	} else if req.Email != "" {
		query = "UPDATE users SET role = 'user' WHERE email = $1 and role = 'banned'"
		res, err := m.PgClient.Exec(query, req.Email)
		if err != nil {
			return err
		}
		if err := config.CheckRowsAffected(res, "user"); err != nil {
			return err
		}
	}
	return nil
}

func (m *UserManager) AddCourier(courier *models.AddCourierReq) error {
	query := "INSERT INTO users (id, email, password, role) VALUES ($1, $2, $3, $4)"
	_, err := m.PgClient.Exec(query, uuid.NewString(), courier.Email, courier.Password, "courier")
	if err != nil {
		return err
	}
	return nil
}

func (m *UserManager) DeleteCourier(req models.DeleteCourierReq) error {
	var query string
	if req.ID != "" {
		query = "DELETE FROM users WHERE id = $1 and role = 'courier'"
		if uuid.Validate(req.ID) != nil {
			return errors.New("invalid user uuid")
		}
		res, err := m.PgClient.Exec(query, req.ID)
		if err != nil {
			return err
		}
		if err := config.CheckRowsAffected(res, "courier"); err != nil {
			return err
		}
	} else if req.Email != "" {
		query = "DELETE FROM users WHERE email = $1 and role = 'courier'"
		res, err := m.PgClient.Exec(query, req.Email)
		if err != nil {
			return err
		}
		if err := config.CheckRowsAffected(res, "courier"); err != nil {
			return err
		}
	}
	return nil
}

func (m *UserManager) AddProductManager(productManager *models.AddProductManagerReq) error {
	query := "INSERT INTO users (id, email, password, role) VALUES ($1, $2, $3, $4)"
	_, err := m.PgClient.Exec(query, uuid.NewString(), productManager.Email, productManager.Password, "manager")
	if err != nil {
		return err
	}
	return nil
}

func (m *UserManager) DeleteProductManager(req models.DeleteProductManagerReq) error {
	var query string
	if req.ID != "" {
		query = "DELETE FROM users WHERE id = $1 and role = 'manager'"
		if uuid.Validate(req.ID) != nil {
			return errors.New("invalid user uuid")
		}
		res, err := m.PgClient.Exec(query, req.ID)
		if err != nil {
			return err
		}
		if err := config.CheckRowsAffected(res, "manager"); err != nil {
			return err
		}
	} else if req.Email != "" {
		query = "DELETE FROM users WHERE email = $1 and role = 'manager'"
		res, err := m.PgClient.Exec(query, req.Email)
		if err != nil {
			return err
		}
		if err := config.CheckRowsAffected(res, "manager"); err != nil {
			return err
		}
	}
	return nil
}
