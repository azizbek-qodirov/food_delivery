package managers

import (
	"auth-service/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserManager struct {
	PgClient    *sql.DB
	MongoClient *mongo.Collection
}

func NewUserManager(db *sql.DB, client *mongo.Client, dbName, collectionName string) *UserManager {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserManager{PgClient: db, MongoClient: collection}
}

func (m *UserManager) Register(req models.RegisterReq, mongoUser *models.MongoUser) error {
	query := "INSERT INTO users (id, email, password, role) VALUES ($1, $2, $3, $4)"
	tx, err := m.PgClient.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction in PostgreSQL: %s", err.Error())
	}

	_, err = tx.Exec(query, req.ID, req.Email, req.Password, req.Role)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to register user in PostgreSQL: %s", err.Error())
	}

	_, err = m.MongoClient.InsertOne(context.Background(), mongoUser)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to register user in MongoDB: %s", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction in PostgreSQL: %s", err.Error())
	}

	return nil
}

func (m *UserManager) ConfirmUser(req *models.ConfirmUserReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "UPDATE users SET is_confirmed = true, confirmed_at = $1 WHERE email = $2"
	_, err := m.PgClient.ExecContext(ctx, query, time.Now(), req.Email)
	return err
}

func (m *UserManager) Profile(req models.GetProfileReq) (*models.GetProfileResp, error) {
	query := "SELECT id, email, password, role, is_confirmed FROM users WHERE email = $1"
	row := m.PgClient.QueryRow(query, req.Email)
	var user models.GetProfileResp
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.IsConfirmed)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserManager) UpdatePassword(req *models.UpdatePasswordReq) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "UPDATE users SET password = $1 WHERE email = $2"
	_, err = m.PgClient.Exec(query, string(hashedPassword), req.Email)
	return err
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

func (m *UserManager) GetByID(id *models.GetProfileByIdReq) (*models.GetProfileByIdResp, error) {
	query := "SELECT id, email, role FROM users WHERE id = $1"
	user := &models.GetProfileByIdResp{}
	err := m.PgClient.QueryRow(query, id.ID).Scan(&user.ID, &user, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}
