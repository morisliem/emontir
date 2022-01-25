package model

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var (
	postgreDBOnce sync.Once
	postgreDB     *sql.DB
)

func NewSQLDB() *sqlx.DB {
	// dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	// if err != nil {
	// 	log.Fatal().Err(fmt.Errorf("error when convert db_port: %w", err)).Send()
	// }
	postgreDBOnce.Do(func() {
		conCfg, err := pgx.ParseURI(os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatal().Err(fmt.Errorf("failed to connect to db: %w", err)).Send()
		}
		// conCfg := pgx.ConnConfig{
		// 	Host:     os.Getenv("DB_HOST"),
		// 	Port:     uint16(dbPort),
		// 	Database: os.Getenv("DB_DATABASE"),
		// 	User:     os.Getenv("DB_USERNAME"),
		// 	Password: os.Getenv("DB_PASSWORD"),
		// }
		db := stdlib.OpenDB(conCfg)
		postgreDB = db
	})

	err := postgreDB.Ping()
	if err != nil {
		log.Fatal().Err(fmt.Errorf("error when ping database: %w", err)).Send()
	}

	return sqlx.NewDb(postgreDB, "postgres")
}

// func NewRedis() *redis.Client {
// 	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
// 	if err != nil {
// 		log.Fatal().Err(fmt.Errorf("error when convert redis_port: %w", err)).Send()
// 	}

// 	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
// 	if err != nil {
// 		log.Fatal().Err(fmt.Errorf("error when convert redis_db: %w", err)).Send()
// 	}

// 	client := redis.NewClient(&redis.Options{
// 		Addr:     fmt.Sprintf("%s:%d", os.Getenv("REDIS_HOST"), redisPort),
// 		Password: os.Getenv("REDIS_PASSWORD"),
// 		DB:       redisDB,
// 	})

// 	_, err = client.Ping().Result()
// 	if err != nil {
// 		log.Fatal().Err(fmt.Errorf("error when pinging redis: %w", err)).Send()
// 		return nil
// 	}
// 	return client
// }

type Manager interface {
	User() User
	Service() Service
	Timeslot() Timeslot
	Cart() Cart
	Order() Order
	Review() Review
}

type manager struct {
	SQLDB *sqlx.DB
	Redis *redis.Client
}

func NewManager() Manager {
	sm := &manager{
		SQLDB: NewSQLDB(),
		// Redis: NewRedis(),
	}
	return sm
}

var (
	userModelOnce sync.Once
	userModel     User
)

func (c *manager) User() User {
	userModelOnce.Do(func() {
		userModel = NewUser(c.SQLDB)
	})
	return userModel
}

var (
	serviceModelOnce sync.Once
	serviceModel     Service
)

func (c *manager) Service() Service {
	serviceModelOnce.Do(func() {
		serviceModel = NewService(c.SQLDB)
	})
	return serviceModel
}

var (
	timeslotModelOnce sync.Once
	timeslotModel     Timeslot
)

func (c *manager) Timeslot() Timeslot {
	timeslotModelOnce.Do(func() {
		timeslotModel = NewTimeslot(c.SQLDB)
	})
	return timeslotModel
}

var (
	cartModelOnce sync.Once
	cartModel     Cart
)

func (c *manager) Cart() Cart {
	cartModelOnce.Do(func() {
		cartModel = NewCart(c.SQLDB)
	})
	return cartModel
}

var (
	orderModelOnce sync.Once
	orderModel     Order
)

func (c *manager) Order() Order {
	orderModelOnce.Do(func() {
		orderModel = NewOrder(c.SQLDB)
	})
	return orderModel
}

var (
	reviewModelOnce sync.Once
	reviewModel     Review
)

func (c *manager) Review() Review {
	reviewModelOnce.Do(func() {
		reviewModel = NewReview(c.SQLDB)
	})
	return reviewModel
}
