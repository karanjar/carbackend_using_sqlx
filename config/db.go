package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Db *sqlx.DB
var Cache *redis.Client
var Client *mongo.Client

func ConnectDb() {
	const (
		user     = "dustin"
		password = "12345"
		dbname   = "cardb"
		host     = "pgdb"
		port     = "5432"
	)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=disable", host, port, user, password, dbname)
	//connecct to database using gorm
	db, err := sqlx.Connect("postgres", dsn)

	//db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Error openning database")
		panic(err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Error connecting to database")
		panic(err)
	}
	fmt.Println("Successfully connected to database")
	Db = db
}
func ConnectCache() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cmd := rdb.Ping(ctx)
	if cmd.Err() != nil {
		fmt.Println("Error connecting to caching database")
		panic(cmd.Err())
	}

	fmt.Println("Successfully connected to caching database")

	Cache = rdb
}
func Connectmongo() {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(opts)
	if err != nil {
		fmt.Println("Error connecting to mongo")
		panic(err)
	}

	Client = client

	fmt.Println("Successfully connected to mongo")
}
