package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/jontynewman/msgsrv/internal/http"
	"github.com/jontynewman/msgsrv/internal/repo"
	r "github.com/jontynewman/msgsrv/internal/repo/redis"
	s "github.com/jontynewman/msgsrv/internal/repo/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

func main() {

	addr := flag.String("addr", "", "The TCP network address to listen on.")
	repoType := flag.String("repo", "", "The type of repository in which to store messages.\nOne of: redis, runtime, sqlite")
	source := flag.String("source", "", "The driver-specific data source name (for Redis and SQL repositories).")
	flag.Parse()

	repo, badUsage, err := initMessageRepository(*repoType, *source)

	if err != nil {
		if badUsage {
			reportBadUsageAndExit(err)
		} else {
			log.Fatal(err)
		}
	}

	add := h.NewAddMessageHandler(repo.Add)
	fetch := h.NewFetchMessageHandler(repo.Fetch)

	http.Handle("POST /messages/{$}", &add)
	http.Handle("GET /messages/{id}", &fetch)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func initMessageRepository(repoType string, source string) (repo.MessageRepository, bool, error) {
	switch repoType {
	case "redis":
		repo, err := initRedisMessageRepository(source)
		if err != nil {
			return nil, false, err
		}
		return repo, false, nil
	case "runtime":
		if source != "" {
			return nil, true, toFlagIsNotApplicableError("source", "runtime")
		}
		return new(repo.RuntimeMessageRepository), false, nil
	case "sqlite":
		repo, err := initSqliteMessageRepository(source)
		if err != nil {
			return nil, false, err
		}
		return repo, false, nil
	case "":
		return nil, true, fmt.Errorf("no repo specified")
	default:
		return nil, true, fmt.Errorf("invalid repo \"%s\"", repoType)
	}
}

func initRedisMessageRepository(url string) (*r.RedisMessageRepository, error) {

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client := redis.NewClient(opts)
	repo := r.NewRedisMessageRepository(ctx, client)

	return &repo, nil
}

func initSqliteMessageRepository(source string) (*s.SqlMessageRepository, error) {

	db, err := sql.Open("sqlite3", source)

	if err != nil {
		return nil, err
	}

	return s.InitSqliteMessageRepository(db)
}

func toFlagIsNotApplicableError(flag string, repoType string) error {
	return fmt.Errorf("%s specified, but is not applicable to the current repo (%s)", flag, repoType)
}

func reportBadUsageAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	flag.Usage()
	os.Exit(1)
}
