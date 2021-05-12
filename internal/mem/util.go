package mem

import (
	"context"
	"fmt"
	"os"

	redis "github.com/go-redis/redis/v8"
)

func newClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	pass := os.Getenv("REDIS_PASS")
	fmt.Println(host, port, pass)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
		DB:       0,
	})
	pong, err := client.Ping(context.TODO()).Result()

	if err != nil {
		fmt.Println("PINGING ERROR " + err.Error())
	}

	fmt.Println(pong)

	return client

}
