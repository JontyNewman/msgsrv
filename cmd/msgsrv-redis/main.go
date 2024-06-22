package main

import (
	"context"
	"flag"
	"log"

	r "internal/repo/redis"

	"github.com/jontynewman/msgsrv"
	"github.com/redis/go-redis/v9"
)

func main() {

	addr := flag.String("addr", "", "The TCP network address for this service to listen on.")
	flag.Parse()

	args := flag.Args()

	expected := 1
	actual := len(args)

	if expected != actual {
		log.Fatalf("expected %d non-flag command-line argument(s) (the Redis URI) but got %d", expected, actual)
	}

	client, err := initRedisClient(args[0])

	if err != nil {
		log.Fatal(err)
	}

	repo := r.NewRedisMessageRepository(context.Background(), client)

	log.Fatal(msgsrv.Run(*addr, &repo))
}

func initRedisClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	return client, nil
}
