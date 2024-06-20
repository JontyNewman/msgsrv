module github.com/jontynewman/msgsrv

require (
    github.com/jontynewman/msgsrv/internal/http v0.1.0
    github.com/jontynewman/msgsrv/internal/repo v0.1.0
)

replace github.com/jontynewman/msgsrv/internal/http => ./internal/http
replace github.com/jontynewman/msgsrv/internal/repo => ./internal/repo

go 1.22.4
