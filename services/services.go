package services

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func init() {
	var err error
	if PgDb == nil {
		PgDb, err = GetPgConn(PgDsn)
		if err != nil {
			panic(err)
		}
	}
	if UserRedis == nil {
		UserRedis = redis.NewClient(&redis.Options{
			Addr:     "0.0.0.0:6379",
			Password: "",
			DB:       0,
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	if mesg, err := UserRedis.Ping(ctx).Result(); err != nil {
		log.Fatalf("%v - %v", mesg, err)
	}
}
