module github.com/jontynewman/msgsrv

require github.com/jontynewman/msgsrv/internal/http v0.1.0

require github.com/jontynewman/msgsrv/internal/repo v0.1.0

require github.com/mattn/go-sqlite3 v1.14.22

replace github.com/jontynewman/msgsrv/internal/http => ./internal/http

replace github.com/jontynewman/msgsrv/internal/repo => ./internal/repo

go 1.22.4
