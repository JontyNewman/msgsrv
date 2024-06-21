package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/jontynewman/msgsrv/internal/http"
	"github.com/jontynewman/msgsrv/internal/repo"
	s "github.com/jontynewman/msgsrv/internal/repo/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	addr := flag.String("addr", "", "The TCP network address to listen on.")
	repoType := flag.String("repo", "", "The type of repository in which to store messages.\nOne of: sqlite, runtime")
	source := flag.String("source", "", "The driver-specific data source name (for SQL-based repositories).")
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
	case "sqlite":
		repo, err := initSqliteMessageRepository(source)
		if err != nil {
			return nil, false, err
		}
		return repo, false, nil
	case "runtime":
		if source != "" {
			warnFlagIsNotApplicable("source", "runtime")
		}
		return new(repo.RuntimeMessageRepository), false, nil
	case "":
		return nil, true, fmt.Errorf("no repo specified")
	default:
		return nil, true, fmt.Errorf("unknown repo \"%s\" specified", repoType)
	}
}

func initSqliteMessageRepository(source string) (*s.SqlMessageRepository, error) {

	db, err := sql.Open("sqlite3", source)

	if err != nil {
		return nil, err
	}

	return s.InitSqliteMessageRepository(db)
}

func warnFlagIsNotApplicable(flag string, repoType string) {
	fmt.Fprintf(os.Stderr, "%s specified, but is not applicable to the current repo (%s)\n", flag, repoType)
}

func reportBadUsageAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	flag.Usage()
	os.Exit(1)
}
