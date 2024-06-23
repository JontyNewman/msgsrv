# Message Service (`msgsrv`) &ndash; Good Growth Technical Test

For the Good Growth technical test, I have developed a simple web service which allows users to store and retrieve plain text
messages. I have chosen to use Go in order to implement it.

An overview for non-technical audiences is provided as a [video](https://raw.githubusercontent.com/JontyNewman/msgsrv/v0.1.0/overview.mp4).

## Quickstart

### Redis via Docker Compose

Docker Compose will start a container for the message service as well as for Redis.

```sh
docker compose up
```

### Redis (`msgsrv-redis`)

- Replace `localhost:80` with the TCP network address for the message service to listen on.
- Replace `redis://localhost:6379/0` with the [Redis URI](https://github.com/redis/go-redis/tree/v9.5.3#connecting-via-a-redis-url) of your instance.

```sh
cd cmd/msgsrv-redis
go run . -addr=localhost:80 redis://localhost:6379/0 
```

### SQLite (`msgsrv-sqlite`)

- Replace `localhost:80` with the TCP network address for the message service to listen on.
- Replace `/path/to/sqlite.db` with an [SQLite connection string](https://github.com/mattn/go-sqlite3/tree/v2.0.6#connection-string) (e.g. `:memory:` for an in-memory database).

```sh
cd cmd/msgsrv-sqlite
go run . -addr=localhost:80 /path/to/sqlite.db
```

### Standalone (`msgsrv`)

Replace `localhost:80` with the TCP network address for the message service to listen on.

```sh
cd cmd/msgsrv
go run . localhost:80
```

### Web Application (`msgsrv-web`)

- Replace `localhost:8080` with the TCP network address for the web application to listen on.
- Replace `http://localhost:80/messages/` with the base address of the message service.

```sh
cd cmd/msgsrv-web
go run . -addr=localhost:8080 -repo=http://localhost:80/messages/
```

## Command Overview

All commands are located in the `cmd/` folder. For usage instructions, refer to [quickstart](#quickstart).

### `msgsrv`

Provides the message service without any persistence. Intended for development purposes.

### `msgsrv-redis`

Provides the message service by storing and retrieving messages using a given Redis instance, specified as a [Redis URI](https://github.com/redis/go-redis/tree/v9.5.3#connecting-via-a-redis-url).

### `msgsrv-sqlite`

Provides the message service by storing and retrieving messages using an SQLite database, specified as a [connection string](https://github.com/mattn/go-sqlite3/tree/v2.0.6#connection-string). A `message` table will be created (if it does not already exist) for persisting messages.

### `msgsrv-web`

Provides a web application for interfacing with the message service. Used to demonstrate the service to non-technical audiences.

## Requirements

### Reliability

The service relies on either Redis or SQLite to store and retrieve messages, both of which are in common use throughout the industry.

It is possible to have multiple instances of the service running, with each instance using the same persistence method, hence providing availability and redundancy.

### Performance

The service does minimal processing when storing and retrieving messages, instead relying on the performance of the persistence layer (either Redis or SQLite).

The Redis persistence layer creates unique identifiers by using the atomic `INCR` operation, which is a [common design pattern](https://redis.io/docs/latest/develop/use/patterns/twitter-clone/#data-layout).

The SQLite persistence layer uses an optimal schema (with the identifier of the message being automatically incremented and indexed as a primary key). It also uses optimal queries (with messages being retrieved using a primary key).

### Maintainability

The codebase is structured in such a way so as to facilitate the creation of other persistence layers (referred to as &ldquo;repositories&rdquo;) without needing to modify the core functionality of the service.

For example, a PostgreSQL persistence layer could be created based on the existing `SqlMessageRepository` type.

```go
func InitPostgresqlMessageRepository(db *sql.DB) (*SqlMessageRepository, error) {

	initStatement := `
	CREATE TABLE IF NOT EXISTS message (
		id SERIAL PRIMARY KEY,
		body TEXT
	);`

	return InitSqlMessageRepository(db, initStatement)
}
```

### Deployability

The service is containerized for Redis using a Dockerfile and is configurable using command-line arguments. It is therefore possible to deploy the application to services such as Google Cloud Run and Google Kubernetes Engine.