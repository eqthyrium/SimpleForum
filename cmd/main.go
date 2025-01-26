package main

import (
	"SimpleForum/internal/config"
	"SimpleForum/internal/repository/sqllite"
	"SimpleForum/internal/service/usecase"
	"SimpleForum/internal/transport/customHttp"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {

	config.Config = config.NewConfiguration()

	db, err := openDb(*config.Config.Dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	repositoryObject := sqllite.NewRepository(db)
	serviceObject := usecase.NewUseCase(repositoryObject)
	httpTransport := customHttp.NewTransportHttpHandler(serviceObject)

	router := httpTransport.Routering()
	message := fmt.Sprintf("The server is running at: http://localhost%s/\n", *config.Config.Addr)
	log.Print(message)
	log.Fatalln(http.ListenAndServe(*config.Config.Addr, router))

}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
