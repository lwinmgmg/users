package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lwinmgmg/user/env"
	"github.com/redis/go-redis/v9"
)

var (
	Env = env.GetEnv()
)

func init() {
	var err error
	// PgDB
	if PgDb == nil {
		var pgDns string = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=UTC",
			Env.Settings.Postgres.Host,
			Env.Settings.Postgres.Port,
			Env.Settings.Postgres.Login,
			Env.Settings.Postgres.Password,
			Env.Settings.Postgres.DB,
		)
		PgDb, err = GetPgConn(pgDns)
		if err != nil {
			panic(err)
		}
	}
	// Redis
	if UserRedis == nil {
		UserRedis = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", Env.Settings.Redis.Host, Env.Settings.Redis.Port),
			Username: Env.Settings.Redis.Login,
			Password: Env.Settings.Redis.Password,
			DB:       Env.Settings.Redis.DB,
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	if mesg, err := UserRedis.Ping(ctx).Result(); err != nil {
		log.Fatalf("%v - %v", mesg, err)
		panic(err)
	}
	// Mail Server
	if MailSender == nil {
		MailSender = NewMailService(Env.Settings.Mail.Login, Env.Settings.Mail.Password, Env.Settings.Mail.Host, Env.Settings.Mail.Port)
	}
	if err := MailSender.Send("Hello", []string{Env.Settings.Mail.Login}); err != nil {
		panic(err)
	}
}

type MessageSender interface {
	Send(string, []string) error
}
