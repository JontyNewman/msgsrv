module github.com/jontynewman/msgsrv

require internal/http v0.1.0

require internal/repo v0.1.0

require internal/repo/redis v0.1.0

require (
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/redis/go-redis/v9 v9.5.3
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace internal/http => ./internal/http

replace internal/repo => ./internal/repo

replace internal/repo/redis => ./internal/repo/redis

go 1.22.4
