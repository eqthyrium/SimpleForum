package main

import (
	"SimpleForum/internal/config"
	"SimpleForum/internal/repository/sqllite"
	"SimpleForum/internal/service/usecase"
	"SimpleForum/internal/transport/customHttp"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conf := config.NewConfiguration()

	db, err := openDb(*conf.Dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if runMigration(db); err != nil {
		log.Fatal(err)
	}
	log.Println("Migration applied succesfully")
	repositoryObject := sqllite.NewRepository(db)
	serviceObject := usecase.NewUseCase(repositoryObject)
	httpTransport := customHttp.NewTransportHttpHandler(serviceObject)

	router := httpTransport.Routering()
	message := fmt.Sprintf("The server is running at: http://localhost%s/\n", *conf.Addr)
	log.Print(message)
	log.Fatalln(http.ListenAndServe(*conf.Addr, router))

}

func runMigration(db *sql.DB) error {
	sqlFil, err := os.Open("../migration/mydatabase.sql")
	if err != nil {
		return fmt.Errorf("Error opening mifration file: %v", err)
	}
	content, err := io.ReadAll(sqlFil)
	if err != nil {
		return fmt.Errorf("Error reading megration file:%v", err)
	}
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("Error exec data in db: %v", err)
	}
	return nil
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
