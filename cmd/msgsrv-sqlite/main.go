package main

import (
	"database/sql"
	"flag"
	"log"

	s "internal/repo/sql"

	"github.com/jontynewman/msgsrv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	addr := flag.String("addr", "", "The TCP network address for this service to listen on.")
	flag.Parse()

	args := flag.Args()

	max := 1
	len := len(args)

	if len > max {
		log.Fatalf("expected no more than %d non-flag command-line argument(s) (the SQLite data source name) but got %d", max, len)
	}

	dataSourceName := ""

	if len > 0 {
		dataSourceName = args[0]
	}

	db, err := sql.Open("sqlite3", dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	repo, err := s.InitSqliteMessageRepository(db)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(msgsrv.Run(*addr, repo))
}
