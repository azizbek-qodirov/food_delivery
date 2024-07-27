package handlers

import (
	"auth-service/config"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0,
})

func (h *HTTPHandler) SendConfirmationCode(email string) error {
	code, err := generateConfirmationCode()
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Load().SENDER_EMAIL)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Recovery Code")
	m.SetBody("text/plain", fmt.Sprintf("Your password recovery code is: %d", code))

	d := gomail.NewDialer("smtp.gmail.com", 587, config.Load().SENDER_EMAIL, config.Load().APP_PASSWORD)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	err = rdb.Set(context.Background(), email, code, 3*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("server error storing confirmation code in Redis")
	}
	return nil
}

func generateConfirmationCode() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}
