package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"github.com/dhaskew/rx/internal/films"
	"github.com/dhaskew/rx/internal/server"
)

func main() {

	envfile := flag.String("envfile", server.DefaultEnvPath, "an environment config file path")
	flag.Parse()

	cfg := server.Config{}
	cfg.Logger = server.NewLogger()
	ops, err := cfg.Load(*envfile)
	if err != nil {
		panic(err)
	}

	db_host := ops["DB_HOST"]
	db_port, err := strconv.Atoi(ops["DB_PORT"])
	if err != nil {
		panic(err)
	}
	db_user := ops["DB_USER"]
	db_password := ops["DB_PASSWORD"]
	db_name := ops["DB_NAME"]

	// open database
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// check db
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	rep := films.NewPostgresFilmRepository(db)

	http_port := ops["HTTP_PORT"]

	server.NewServer(
		server.WithLogger(cfg.Logger),

		server.WithRouterFunc(chi.NewRouter),
		server.WithFilmRepository(&rep),
		server.WithPort(http_port)).Start()

}

//NOTE: probably should move eval slog / replace zap, now part of standard lib
//NOTE: telementry and jaeger? tracing can be crazy useful
//TODO: caching
//TODO: code coverage
//TODO: integration tests
