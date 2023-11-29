package helper

import (
	"context"
	"database/sql"
	"time"

	"github.com/ariefsn/go-resik/logger"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MySqlClient(address string) *sql.DB {
	db, err := sql.Open("mysql", address)
	if err != nil {
		logger.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func MongoClient(address string) (client *mongo.Client, cancel context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(address))
	if err != nil {
		logger.Fatal(err)
	}

	return
}
