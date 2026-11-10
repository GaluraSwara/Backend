package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB menghubungkan ke MongoDB
func ConnectDB() {
	mongoURI := os.Getenv("MONGOSTRING")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DBNAME")
	if dbName == "" {
		dbName = "galuraswara"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Gagal konek ke MongoDB: ", err)
	}

	// Ping untuk memastikan koneksi berhasil
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Gagal ping MongoDB: ", err)
	}

	DB = client.Database(dbName)
	log.Println("✅ Terhubung ke MongoDB:", dbName)
}

// GetCollection mengembalikan collection dari database
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}
