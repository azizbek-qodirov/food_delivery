package storage

import (
	"auth-service/config"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(cf *config.Config) (*sql.DB, *mongo.Client, error) {
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cf.DB_USER, cf.DB_PASSWORD, cf.DB_HOST, cf.DB_PORT, cf.DB_NAME)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, nil, err
	}
	err = db.Ping()
	if err != nil {
		panic("Postgres not connected due to error: " + err.Error())
	}

	clientOptions := options.Client().ApplyURI(cf.MONGO_URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, nil, err
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		panic("MongoDB not connected due to error: " + err.Error())
	}
	fmt.Println("Successfully connected to database mongodb!")

	return db, client, nil
}
