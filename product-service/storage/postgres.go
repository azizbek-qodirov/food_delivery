package storage

import (
	"context"
	"database/sql"
	"fmt"

	"progress-service/config"
	"progress-service/storage/managers"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	PgClient    *sql.DB
	MongoClient *mongo.Client

	ProductS ProductI
}

func NewPostgresStorage(config config.Config) (*Storage, error) {
	// #################    POSTGRESQL CONNECTION     ###################### //
	conn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		config.DB_HOST, config.DB_USER, config.DB_NAME, config.DB_PASSWORD, config.DB_PORT)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to the database pgsql!!!")

	// #################     MONGODB CONNECTION     ###################### //
	clientOptions := options.Client().ApplyURI(config.MONGO_URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to the database mongodb!!!")

	pm := managers.NewProductManager(client, config.MONGO_DB_NAME, config.MONGO_COLLECTION_NAME, db)

	return &Storage{
		PgClient: db,
		ProductS: pm,
	}, nil
}

func (s *Storage) Product() ProductI {
	if s.ProductS == nil {
		s.ProductS = managers.NewProductManager(s.MongoClient, config.Load().MONGO_DB_NAME, config.Load().MONGO_COLLECTION_NAME, s.PgClient)
	}
	return s.ProductS
}
