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

func MySqlClient() *sql.DB {
	db, err := sql.Open("mysql", "root:Password.123@localhost/resik-arch")
	if err != nil {
		logger.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func MongoClient() (client *mongo.Client, cancel context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:Password.123@localhost:27017"))
	if err != nil {
		logger.Fatal(err)
	}

	return
}
