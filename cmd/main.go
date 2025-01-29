package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"SimpleForum/internal/config"
	"SimpleForum/internal/repository/sqllite"
	"SimpleForum/internal/service/usecase"
	"SimpleForum/internal/transport/customHttp"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("You are in the directory:", dir)

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
	message := fmt.Sprintf("The server is running at: https://localhost%s/\n", *config.Config.Addr)
	log.Print(message)
	// log.Fatalln(http.ListenAndServe(*config.Config.Addr, router))

	ch := make(chan int)
	go func() {
		log.Fatalln(http.ListenAndServeTLS(*config.Config.Addr, "./tls/cert.pem", "./tls/key.pem", router))
		ch <- 1
	}()
	<-ch
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
