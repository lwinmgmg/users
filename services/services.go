package services

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func init() {
	var err error
	// PgDB
	if PgDb == nil {
		PgDb, err = GetPgConn(PgDsn)
		if err != nil {
			panic(err)
		}
	}
	// Redis
	if UserRedis == nil {
		UserRedis = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	if mesg, err := UserRedis.Ping(ctx).Result(); err != nil {
		log.Fatalf("%v - %v", mesg, err)
	}
	// Mail Server
	email := os.Getenv("GO_EMAIL")
	password := os.Getenv("GO_EMAIL_PASSWORD")
	if MailSender == nil {
		MailSender = NewMailService(email, password, "smtp.gmail.com", "587")
	}
	if err := MailSender.Send("Hello", []string{email}); err != nil {
		log.Println("No email provided")
		// panic(err)
	}
}

type MessageSender interface {
	Send(string, []string) error
}
